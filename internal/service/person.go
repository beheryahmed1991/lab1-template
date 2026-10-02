package service

import (
	"context"

	"github.com/bmstu-rsoi/lab1-template.git/internal/model"
	"github.com/bmstu-rsoi/lab1-template.git/internal/repository"
)

type PersonService struct {
	repo *repository.PersonRepository
}

func (s *PersonService) Create(ctx context.Context, person model.PersonRequest) (int64, error) {
	return s.repo.Create(ctx, person)
}

func (s *PersonService) GetByID(ctx context.Context, id int64) (*model.Person, error) {
	return s.repo.GetByID(ctx, id)
}
