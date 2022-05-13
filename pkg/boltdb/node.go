package boltdb

import (
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	"github.com/sku4/alice-checklist/models/googlekeep"
	"regexp"
	"strings"
)

type Bucket string

const (
	NoteTextIds Bucket = "note_text_ids"
	NoteIdsBody Bucket = "note_ids_body"
)

var (
	ErrorNodeNotFound    = errors.New("node not found")
	ErrorEmptyBucketExit = errors.New("bucket exit")
	ErrorClearTextEmpty  = errors.New("text is empty after clear")
	BucketsNote          = []Bucket{
		NoteTextIds,
		NoteIdsBody,
	}
)

type NodeCacheRepository struct {
	db Storage
}

func NewNodeCacheRepository(db Storage) *NodeCacheRepository {
	return &NodeCacheRepository{db: db}
}

func (r *NodeCacheRepository) Save(node googlekeep.Node) (err error) {
	err = r.db.Update(func(tx *bolt.Tx) error {
		b1 := tx.Bucket([]byte(NoteTextIds))
		cl, err := ClearText(node.Text)
		if err != nil {
			return err
		}
		err = b1.Put([]byte(cl), []byte(node.Id))
		if err != nil {
			return err
		}

		b2 := tx.Bucket([]byte(NoteIdsBody))
		nodeJson, err := json.Marshal(node)
		return b2.Put([]byte(node.Id), nodeJson)
	})

	return
}

func (r *NodeCacheRepository) Get(name string) (node *googlekeep.Node, err error) {
	err = r.db.View(func(tx *bolt.Tx) (err error) {
		b1 := tx.Bucket([]byte(NoteTextIds))
		cl, err := ClearText(name)
		if err != nil {
			return
		}
		id := string(b1.Get([]byte(cl)))
		if id == "" {
			err = ErrorNodeNotFound
			return
		}

		b2 := tx.Bucket([]byte(NoteIdsBody))
		nodeJson := b2.Get([]byte(id))
		return json.Unmarshal(nodeJson, &node)
	})

	return
}

func (r *NodeCacheRepository) List() (nodes []googlekeep.Node, err error) {
	err = r.db.View(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket([]byte(NoteIdsBody))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var node googlekeep.Node
			err = json.Unmarshal(v, &node)
			if err != nil {
				return
			}
			nodes = append(nodes, node)
		}
		return
	})

	return
}

func (r *NodeCacheRepository) IsEmptyBucket() (empty bool, err error) {
	empty = true
	err = r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(NoteTextIds))

		err := b.ForEach(func(k, v []byte) error {
			empty = false
			return ErrorEmptyBucketExit
		})
		if err != nil && err != ErrorEmptyBucketExit {
			return err
		}
		return nil
	})
	return
}

func (r *NodeCacheRepository) Truncate() (err error) {
	err = r.db.Update(func(tx *bolt.Tx) error {
		for _, b := range BucketsNote {
			err := tx.DeleteBucket([]byte(b))
			if err != nil {
				return err
			}
			_, err = tx.CreateBucketIfNotExists([]byte(b))
			if err != nil {
				return err
			}
		}
		return nil
	})

	return
}

func ClearText(s string) (string, error) {
	s = strings.ToUpper(s)
	reg, err := regexp.Compile("[^a-zA-Zа-яА-Я0-9]+")
	if err != nil {
		return "", err
	}
	s = reg.ReplaceAllString(s, "")
	if s == "" {
		return s, ErrorClearTextEmpty
	}
	return s, nil
}
