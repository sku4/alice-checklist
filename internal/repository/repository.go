package repository

import (
	"github.com/sku4/alice-checklist/configs"
	"github.com/sku4/alice-checklist/internal/repository/googlekeep"
	"github.com/sku4/alice-checklist/internal/repository/notify"
	"github.com/sku4/alice-checklist/lang"
	model "github.com/sku4/alice-checklist/models/googlekeep"
	"github.com/sku4/alice-checklist/pkg/boltdb"
)

//go:generate mockgen -source=repository.go -destination=mocks/repository.go

type Checklist interface {
	Patch(add, delete []string) error
	List() ([]model.Node, error)
	CacheList() ([]model.Node, error)
}

type Notify interface {
	Add(error) error
	Get() ([]error, error)
	Truncate() error
}

type Repository struct {
	Checklist
	Notify
}

func NewRepository(loc *lang.Localize, cfg *configs.Config, db boltdb.Storage) *Repository {
	return &Repository{
		Checklist: googlekeep.NewChecklist(loc, cfg, db),
		Notify:    notify.NewNotify(cfg, db),
	}
}
