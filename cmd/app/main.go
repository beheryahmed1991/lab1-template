package main

import (
	"context"
	"log"

	"github.com/bmstu-rsoi/lab1-template.git/internal/handler"
	"github.com/bmstu-rsoi/lab1-template.git/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	db, err := pgxpool.New(
		ctx,
		"postgres://program:test@localhost:5432/persons",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewPersonRepository(db)
	svc := service.NewPersonService(repo)
	h := handler.NewPersonHandler(svc)

	router := gin.Default()

	api := router.Group("/api/v1")
	{
		api.GET("/persons", h.GetAll)
		api.GET("/persons/:id", h.GetByID)
		api.POST("/persons", h.Create)
	}

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
