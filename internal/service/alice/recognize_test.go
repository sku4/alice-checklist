package alice

import (
	"github.com/sku4/alice-checklist/lang"
	"github.com/sku4/alice-checklist/models/googlekeep"
	testify "github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestService_listToText(t *testing.T) {
	loc, err := lang.InitLocalize(lang.Ru)
	if err != nil {
		testify.Fail(t, err.Error())
		return
	}

	nodes := make(map[int]*googlekeep.Node)
	for i := 1; i <= 3; i++ {
		nodes[i] = googlekeep.NewNode()
		nodes[i].SortValue = strconv.Itoa(i * 10)
		nodes[i].Text = strconv.Itoa(i * 10)
	}

	tests := []struct {
		name     string
		nodes    []googlekeep.Node
		wantText string
		wantTts  string
	}{
		{
			name: "Sort by value",
			nodes: []googlekeep.Node{
				*nodes[3],
				*nodes[1],
				*nodes[2],
			},
			wantText: "Your list: \n○ 30\n○ 20\n○ 10",
			wantTts:  "Your list:  sil<[200]> 30 sil<[100]>, 20 sil<[100]>, 10",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{loc: loc}
			gotText, gotTts := s.listToText(tt.nodes)
			if gotText != tt.wantText {
				t.Errorf("listToText() gotText = %v, want %v", gotText, tt.wantText)
			}
			if gotTts != tt.wantTts {
				t.Errorf("listToText() gotTts = %v, want %v", gotTts, tt.wantTts)
			}
		})
	}
}
