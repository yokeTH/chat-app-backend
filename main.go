package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/contrib/websocket"
	"github.com/yokeTH/chat-app-backend/internal/adaptor/dto"
	"github.com/yokeTH/chat-app-backend/internal/adaptor/handler"
	"github.com/yokeTH/chat-app-backend/internal/adaptor/middleware"
	"github.com/yokeTH/chat-app-backend/internal/adaptor/repository"
	wsAdaptor "github.com/yokeTH/chat-app-backend/internal/adaptor/websocket"
	"github.com/yokeTH/chat-app-backend/internal/config"
	"github.com/yokeTH/chat-app-backend/internal/server"
	"github.com/yokeTH/chat-app-backend/internal/usecase/book"
	"github.com/yokeTH/chat-app-backend/internal/usecase/conversation"
	"github.com/yokeTH/chat-app-backend/internal/usecase/file"
	"github.com/yokeTH/chat-app-backend/internal/usecase/message"
	"github.com/yokeTH/chat-app-backend/internal/usecase/user"
	"github.com/yokeTH/chat-app-backend/pkg/db"
	"github.com/yokeTH/chat-app-backend/pkg/storage"
)

// @title GO-FIBER-TEMPLATE API
// @version 1.0
// @description Bearer token authentication
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @schemes http https
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

	// Setup Translator (Dto)
	fileDto := dto.NewFileDto(publicBucket)
	userDto := dto.NewUserDto()
	reactionDto := dto.NewReactionDto(userDto)
	messageDto := dto.NewMessageDto(fileDto, reactionDto, userDto)
	conversationDto := dto.NewConversationDto(userDto, messageDto)

	// Setup repository
	bookRepo := repository.NewBookRepository(db)
	fileRepo := repository.NewFileRepository(db)
	userRepo := repository.NewUserRepository(db)
	conversationRepo := repository.NewConversationRepository(db)
	messageRepo := repository.NewMessageRepository(db)

	// Setup use cases
	bookUC := book.NewBookUseCase(bookRepo)
	fileUC := file.NewFileUseCase(fileRepo, publicBucket)
	msgUC := message.NewMessageUseCase(messageRepo)
	userUC := user.NewUserUseCase(userRepo)
	conversationUC := conversation.NewConversationUseCase(conversationRepo)

	// Setup message server
	msgServer := wsAdaptor.NewMessageServer(userUC, msgUC, conversationUC, messageDto)
	go msgServer.Start(ctx, stop)

	// Setup handlers
	authHandler := handler.NewAuthHandler(userUC)
	bookHandler := handler.NewBookHandler(bookUC)
	fileHandler := handler.NewFileHandler(fileUC, fileDto, msgUC, messageDto, msgServer)
	msgHandler := handler.NewMessageHandler(msgUC, messageDto)
	conversationHandler := handler.NewConversationHandler(conversationUC, conversationDto, msgServer, msgUC, messageDto)
	userHandler := handler.NewUserHandler(userUC, userDto, msgServer)

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
		ws := s.Group("/ws", wsMiddleware.RequiredUpgradeProtocol)
		{
			ws.Get("/", websocket.New(msgServer.HandleWebsocket))
		}
	}
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
		message := s.Group("/messages", authMiddleware.Auth)
		{
			message.Post("/", msgHandler.HandleCreateMessage)
			message.Get("/:id", msgHandler.HandleGetMessage)
		}
	}
	{
		conversation := s.Group("/conversations", authMiddleware.Auth)
		{
			conversation.Get("/", conversationHandler.HandleListConversation)
			conversation.Post("/", conversationHandler.HandleCreateConversation)
			conversation.Get("/:conversationID/messages", msgHandler.HandleListMessagesByConversation)
			conversation.Get("/:id", conversationHandler.HandleGetConversation)
			conversation.Post("/:id/files", fileHandler.CreateFile)
		}
	}
	{
		user := s.Group("/users", authMiddleware.Auth)
		{
			user.Get("/", userHandler.HandleListUser)
			user.Get("/me", userHandler.HandleGetMe)
			user.Patch("/:id", userHandler.HandleUpdateUser)
		}
	}

	// Start the server
	s.Start(ctx, stop)
}
