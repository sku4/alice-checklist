package boltdb

import (
	"github.com/sku4/alice-checklist/models/googlekeep"
)

type CacheRepository interface {
	Save(node googlekeep.Node) error
	Get(name string) (*googlekeep.Node, error)
	Truncate() error
	IsEmptyBucket() (bool, error)
	List() ([]googlekeep.Node, error)
}
