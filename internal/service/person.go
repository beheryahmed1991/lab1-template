package service

import (
	"context"

	"github.com/bmstu-rsoi/lab1-template.git/internal/model"
	"github.com/bmstu-rsoi/lab1-template.git/internal/repository"
)

type PersonService struct {
	repo *repository.PersonRepository
}

func NewPersonService(repo *repository.PersonRepository) *PersonService {
	return &PersonService{repo: repo}
}

func (s *PersonService) Create(ctx context.Context, person model.PersonRequest) (int64, error) {
	return s.repo.Create(ctx, person)
}

func (s *PersonService) GetByID(ctx context.Context, id int64) (*model.Person, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PersonService) GetAll(ctx context.Context) ([]model.Person, error) {
	return s.repo.GetAll(ctx)
}

func (s *PersonService) Update(ctx context.Context, id int64, person model.PersonRequest) (*model.Person, error) {
	return s.repo.Update(ctx, id, person)
}

func (s *PersonService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
