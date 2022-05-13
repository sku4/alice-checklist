package boltdb

import (
	"github.com/boltdb/bolt"
)

//go:generate mockgen -source=storage.go -destination=mocks/storage.go

type Storage interface {
	Close() error
	Update(fn func(*bolt.Tx) error) error
	View(fn func(*bolt.Tx) error) error
}
