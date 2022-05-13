package notify

import (
	"github.com/sku4/alice-checklist/configs"
	"github.com/sku4/alice-checklist/pkg/boltdb"
)

type Notify struct {
	errors *boltdb.ErrorsRepository
	config *configs.Config
}

func NewNotify(cfg *configs.Config, db boltdb.Storage) *Notify {
	return &Notify{
		errors: boltdb.NewErrorsRepository(db),
		config: cfg,
	}
}

func (c *Notify) Add(err error) error {
	return c.errors.Add(err)
}

func (c *Notify) Get() ([]error, error) {
	return c.errors.Get()
}

func (c *Notify) Truncate() error {
	return c.errors.Truncate()
}
