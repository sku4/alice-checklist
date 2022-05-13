package boltdb

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"hash/crc32"
	"strconv"
)

const (
	NotifyErrors Bucket = "notify_errors"
)

var (
	BucketsNotify = []Bucket{
		NotifyErrors,
	}
)

type ErrorsRepository struct {
	db Storage
}

func NewErrorsRepository(db Storage) *ErrorsRepository {
	return &ErrorsRepository{db: db}
}

func (r *ErrorsRepository) Add(errInput error) (err error) {
	err = r.db.Update(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket([]byte(NotifyErrors))
		errJson, err := json.Marshal(errInput.Error())
		if err != nil {
			return
		}
		k := genHashId(errJson)
		return b.Put([]byte(k), errJson)
	})

	return
}

func (r *ErrorsRepository) Get() (errs []error, err error) {
	err = r.db.View(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket([]byte(NotifyErrors))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var e error
			err = json.Unmarshal(v, &e)
			if err != nil {
				return
			}
			errs = append(errs, e)
		}
		return
	})

	return
}

func (r *ErrorsRepository) Truncate() (err error) {
	err = r.db.Update(func(tx *bolt.Tx) error {
		for _, b := range BucketsNotify {
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

func genHashId(b []byte) string {
	crc32q := crc32.MakeTable(crc32.IEEE)
	k := crc32.Checksum(b, crc32q)
	return strconv.FormatUint(uint64(k), 10)
}
