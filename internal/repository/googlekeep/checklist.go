package googlekeep

import (
	"github.com/sku4/alice-checklist/configs"
	"github.com/sku4/alice-checklist/lang"
	model "github.com/sku4/alice-checklist/models/googlekeep"
	"github.com/sku4/alice-checklist/pkg/boltdb"
	"github.com/sku4/alice-checklist/pkg/googlekeep"
)

type Checklist struct {
	*googlekeep.Client
}

func NewChecklist(loc *lang.Localize, cfg *configs.Config, db boltdb.Storage) *Checklist {
	return &Checklist{
		Client: googlekeep.NewClient(loc, cfg, db),
	}
}

func (c *Checklist) Patch(add, delete []string) (err error) {
	err = c.Client.Patch(add, delete)
	if err != nil {
		return
	}

	go func() {
		_ = c.Client.Clean()
	}()

	return
}

func (c *Checklist) List() (nodes []model.Node, err error) {
	nodes, err = c.Client.List()
	if err != nil {
		return
	}

	return
}

func (c *Checklist) CacheList() (nodes []model.Node, err error) {
	nodes, err = c.Client.CacheList()
	if err != nil {
		return
	}

	return
}
