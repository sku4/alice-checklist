package service

import (
	"github.com/sku4/alice-checklist/internal/repository"
	"github.com/sku4/alice-checklist/internal/service/alice"
	"github.com/sku4/alice-checklist/lang"
	model "github.com/sku4/alice-checklist/models/alice"
)

//go:generate mockgen -source=service.go -destination=mocks/service.go

type Alice interface {
	Command(model.Request) (model.Response, error)
}

type Service struct {
	Alice
}

func NewService(loc *lang.Localize, repos *repository.Repository) *Service {
	return &Service{
		Alice: alice.NewService(loc, repos.Checklist, repos.Notify),
	}
}
