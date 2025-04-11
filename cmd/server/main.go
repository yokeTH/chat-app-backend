package main

import (
	"context"
	"log"

	"github.com/gofiber/swagger"
	"github.com/swaggo/swag"
	"github.com/yokeTH/gofiber-template/docs"
	"github.com/yokeTH/gofiber-template/internal/adaptor/handler"
	"github.com/yokeTH/gofiber-template/internal/adaptor/repository"
	"github.com/yokeTH/gofiber-template/internal/config"
	"github.com/yokeTH/gofiber-template/internal/server"
	"github.com/yokeTH/gofiber-template/internal/usecase/book"
	"github.com/yokeTH/gofiber-template/internal/usecase/file"
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
	ctx, stop := context.WithCancel(context.Background())
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

	// Setup use cases
	bookUC := book.NewBookUseCase(bookRepo)
	fileUC := file.NewFileUseCase(fileRepo, publicBucket, privateBucket)

	// Setup handlers
	bookHandler := handler.NewBookHandler(bookUC)
	fileHandler := handler.NewFileHandler(fileUC, privateBucket, publicBucket)

	// Setup server
	s := server.New(
		server.WithName(config.Server.Name),
		server.WithBodyLimitMB(config.Server.BodyLimitMB),
		server.WithPort(config.Server.Port),
	)

	// Setup routes
	swag.Register(docs.SwaggerInfo.InstanceName(), docs.SwaggerInfo)
	s.Get("/swagger/*", swagger.HandlerDefault)

	{
		book := s.Group("/book")
		{
			book.Get("", bookHandler.GetBooks)
			book.Get("/:id", bookHandler.GetBook)
			book.Post("", bookHandler.CreateBook)
			book.Patch("/:id", bookHandler.UpdateBook)
			book.Delete("/:id", bookHandler.DeleteBook)
		}
	}
	{
		file := s.Group("/file")
		{
			file.Get("/", fileHandler.List)
			file.Get("/:id", fileHandler.GetInfo)
			file.Post("/private", fileHandler.CreatePrivateFile)
			file.Post("/public", fileHandler.CreatePublicFile)
		}
	}

	// Start the server
	s.Start(ctx, stop)
}
