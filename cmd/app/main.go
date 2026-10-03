package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/bmstu-rsoi/lab1-template.git/internal/handler"
	"github.com/bmstu-rsoi/lab1-template.git/internal/repository"
	"github.com/bmstu-rsoi/lab1-template.git/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://program:test@localhost:5432/persons"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db, err := pgxpool.New(
		ctx,
		databaseURL,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := db.Ping(pingCtx); err != nil {
		log.Fatal("database connection failed: ", err)
	}

	schemaCtx, cancelSchema := context.WithTimeout(ctx, 30*time.Second)
	defer cancelSchema()
	_, err = db.Exec(schemaCtx, `CREATE TABLE IF NOT EXISTS persons (
		id BIGSERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		age INTEGER,
		address VARCHAR(255),
		work VARCHAR(255)
	)`)
	if err != nil {
		log.Fatal("database initialization failed: ", err)
	}

	repo := repository.NewPersonRepository(db)
	svc := service.NewPersonService(repo)
	h := handler.NewPersonHandler(svc)

	router := gin.Default()

	api := router.Group("/api/v1")
	{
		api.GET("/persons", h.GetAll)
		api.GET("/persons/:id", h.GetByID)
		api.POST("/persons", h.Create)
		api.PATCH("/persons/:id", h.Update)
		api.DELETE("/persons/:id", h.Delete)
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
