package boltdb

import (
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/golang/mock/gomock"
	"github.com/sku4/alice-checklist/models/googlekeep"
	mock_boltdb "github.com/sku4/alice-checklist/pkg/boltdb/mocks"
	"testing"
)

func TestNodeCacheRepository_Save(t *testing.T) {
	type mockBehavior func(r *mock_boltdb.MockStorage, fn txHandler)
	txFunc := func(tx *bolt.Tx) error {
		return nil
	}

	tests := []struct {
		name    string
		node    googlekeep.Node
		wantErr bool
		mockBehavior
	}{
		{
			name:    "Ok",
			node:    *googlekeep.NewNode(),
			wantErr: false,
			mockBehavior: func(r *mock_boltdb.MockStorage, fn txHandler) {
				r.EXPECT().Update(newHandlerMatcher(fn)).Return(nil)
			},
		},
		{
			name:    "Check error",
			node:    *googlekeep.NewNode(),
			wantErr: true,
			mockBehavior: func(r *mock_boltdb.MockStorage, fn txHandler) {
				r.EXPECT().Update(newHandlerMatcher(fn)).Return(errors.New("something went wrong"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			dbStorage := mock_boltdb.NewMockStorage(c)
			tt.mockBehavior(dbStorage, txFunc)

			r := &NodeCacheRepository{
				db: dbStorage,
			}
			if err := r.Save(tt.node); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type txHandler func(*bolt.Tx) error

// https://github.com/zachwalton/gomock-test
type handlerMatcher struct {
	handler txHandler
}

var _ gomock.Matcher = &handlerMatcher{}

// Matches implements gomock.Matcher
func (fm *handlerMatcher) Matches(x interface{}) bool {
	_, ok := x.(func(*bolt.Tx) error)
	if !ok {
		return false
	}

	// compare output of expected handler and handler invoked by gmt.Start()
	return true
}

// String implements gomock.Matcher
func (fm handlerMatcher) String() string {
	return fmt.Sprintf("handler: %+v", fm.handler)
}

func newHandlerMatcher(handler txHandler) *handlerMatcher {
	return &handlerMatcher{handler: handler}
}
