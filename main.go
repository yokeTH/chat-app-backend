package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/contrib/websocket"
	"github.com/yokeTH/gofiber-template/internal/adaptor/handler"
	"github.com/yokeTH/gofiber-template/internal/adaptor/middleware"
	"github.com/yokeTH/gofiber-template/internal/adaptor/repository"
	"github.com/yokeTH/gofiber-template/internal/config"
	"github.com/yokeTH/gofiber-template/internal/server"
	"github.com/yokeTH/gofiber-template/internal/usecase/book"
	"github.com/yokeTH/gofiber-template/internal/usecase/conversation"
	"github.com/yokeTH/gofiber-template/internal/usecase/file"
	"github.com/yokeTH/gofiber-template/internal/usecase/message"
	"github.com/yokeTH/gofiber-template/internal/usecase/user"
	"github.com/yokeTH/gofiber-template/pkg/db"
	"github.com/yokeTH/gofiber-template/pkg/storage"
)

// @title GO-FIBER-TEMPLATE API
// @version 1.0
// @servers https http
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Bearer token authentication
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	config := config.Load()

	// Setup infrastructure
	db, err := db.New(config.PSQL)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	publicBucket, err := storage.New(config.PublicBucket)
	if err != nil {
		log.Fatalf("failed to create public bucket instance: %v", err)
	}

	privateBucket, err := storage.New(config.PrivateBucket)
	if err != nil {
		log.Fatalf("failed to create private bucket instance: %v", err)
	}

	// Setup repository
	bookRepo := repository.NewBookRepository(db)
	fileRepo := repository.NewFileRepository(db)
	userRepo := repository.NewUserRepository(db)
	conversationRepo := repository.NewConversationRepository(db)

	msgServer := message.NewMessageServer(userRepo)
	go msgServer.Start(ctx, stop)

	// Setup use cases
	bookUC := book.NewBookUseCase(bookRepo)
	fileUC := file.NewFileUseCase(fileRepo, publicBucket, privateBucket)
	msgUC := message.NewMessageUseCase(msgServer)
	userUC := user.NewUserUseCase(userRepo)
	conversationUC := conversation.NewConversationUseCase(conversationRepo)

	// Setup handlers
	authHandler := handler.NewAuthHandler(userUC)
	bookHandler := handler.NewBookHandler(bookUC)
	fileHandler := handler.NewFileHandler(fileUC, privateBucket, publicBucket)
	msgHandler := handler.NewMessageHandler(msgUC)
	conversationHandler := handler.NewConversationHandler(conversationUC)

	// Setup middleware
	authMiddleware := middleware.NewAuthMiddleware(userUC)
	wsMiddleware := middleware.NewWebsocketMiddleware()

	// Setup server
	s := server.New(
		server.WithName(config.Server.Name),
		server.WithBodyLimitMB(config.Server.BodyLimitMB),
		server.WithPort(config.Server.Port),
		server.WithEnv(config.Server.Env),
		server.WithSwaggerProtection(config.Server.SwaggerUser, config.Server.SwaggerPass),
	)

	// Setup routes
	{
		auth := s.Group("/auth")
		{
			auth.Post("/google", authMiddleware.Auth, authHandler.HandleGoogleLogin)
		}
	}
	{
		book := s.Group("/books")
		{
			book.Get("", bookHandler.GetBooks)
			book.Get("/:id", bookHandler.GetBook)
			book.Post("", bookHandler.CreateBook)
			book.Patch("/:id", bookHandler.UpdateBook)
			book.Delete("/:id", bookHandler.DeleteBook)
		}
	}
	{
		file := s.Group("/files")
		{
			file.Get("/", fileHandler.List)
			file.Get("/:id", fileHandler.GetInfo)
			file.Post("/private", fileHandler.CreatePrivateFile)
			file.Post("/public", fileHandler.CreatePublicFile)
		}
	}
	{
		message := s.Group("/message")
		{
			message.Use("/ws", wsMiddleware.RequiredUpgradeProtocol)
			message.Get("/ws", websocket.New(msgHandler.HandleMessage))
		}
	}
	{
		conversation := s.Group("/conversation", authMiddleware.Auth)
		{
			conversation.Get("/", conversationHandler.HandleListConversation)
		}
	}

	// Start the server
	s.Start(ctx, stop)
}
