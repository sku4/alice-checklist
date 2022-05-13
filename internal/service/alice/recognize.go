package alice

import (
	"fmt"
	"github.com/sku4/alice-checklist/models/alice"
	"github.com/sku4/alice-checklist/models/googlekeep"
	"sort"
	"strconv"
	"strings"
)

const (
	ListIntentName      = "list"
	AddDeleteIntentName = "add_delete_products"
	AddPrefix           = "add"
	DeletePrefix        = "del"
)

//go:generate mockgen -source=recognize.go -destination=mocks/recognize.go

type Recognize interface {
	Recognize(chanId int, req alice.Request)
}

func (s *Service) Recognize(chanId int, req alice.Request) {
	command := req.ReqData.Command
	if command == "" {
		s.welcome(chanId)
		return
	}
	for i, intent := range req.ReqData.Nlu.Intents {
		switch i {
		case AddDeleteIntentName:
			var add, del []string
			addCheck := make(map[string]bool)
			delCheck := make(map[string]bool)
			addKeys := make([]int, 0, len(intent.Slots))
			for k, slot := range intent.Slots {
				if strings.HasPrefix(k, AddPrefix) {
					if _, exist := addCheck[slot.Value]; !exist {
						addCheck[slot.Value] = true
						r := []rune(k)
						c, err := strconv.Atoi(string(r[3:]))
						if err != nil {
							s.error(chanId, err)
							return
						}
						addKeys = append(addKeys, c)
					}
				} else if strings.HasPrefix(k, DeletePrefix) {
					if _, exist := delCheck[slot.Value]; !exist {
						delCheck[slot.Value] = true
						del = append(del, slot.Value)
					}
				}
			}

			sort.Slice(addKeys, func(i, j int) bool {
				return addKeys[i] > addKeys[j]
			})
			for _, k := range addKeys {
				if slot, exist := intent.Slots[AddPrefix+strconv.Itoa(k)]; exist {
					add = append(add, slot.Value)
				}
			}

			s.patch(chanId, add, del)
			return
		case ListIntentName:
			s.list(chanId)
			return
		}
	}
	s.repeat(chanId)
	return
}

func (s *Service) welcome(chanId int) {
	resp := *alice.NewResponse()
	resp.RespData.Text = fmt.Sprintf(s.loc.Translate("To complete the list, say add or remove"))
	resp.RespData.Tts = resp.RespData.Text
	s.chanAnswer.HotAnswer(chanId, resp)
	s.chanAnswer.ColdAnswer(chanId, resp, nil)
}

func (s *Service) repeat(chanId int) {
	resp := *alice.NewResponse()
	resp.RespData.Text = fmt.Sprintf(s.loc.Translate("Repeat please"))
	resp.RespData.Tts = resp.RespData.Text
	s.chanAnswer.HotAnswer(chanId, resp)
	s.chanAnswer.ColdAnswer(chanId, resp, nil)
}

func (s *Service) list(chanId int) {
	resp := *alice.NewResponse()
	resp.EndSession = true

	nodes, err := s.checklist.CacheList()
	text, tts := s.listToText(nodes)
	resp.RespData.Text = text
	resp.RespData.Tts = tts
	s.chanAnswer.HotAnswer(chanId, resp)

	nodes, err = s.checklist.List()
	text, tts = s.listToText(nodes)
	resp.RespData.Text = text
	resp.RespData.Tts = tts
	s.chanAnswer.ColdAnswer(chanId, resp, err)
}

func (s *Service) patch(chanId int, add, del []string) {
	resp := *alice.NewResponse()
	resp.EndSession = true
	resp.RespData.Text = fmt.Sprintf(s.loc.Translate("Process request"))
	resp.RespData.Tts = fmt.Sprintf(s.loc.Translate("Process request"))
	s.chanAnswer.HotAnswer(chanId, resp)

	err := s.checklist.Patch(add, del)
	resp.RespData.Text = fmt.Sprintf(s.loc.Translate("Points success updated"))
	resp.RespData.Tts = fmt.Sprintf(s.loc.Translate("Points success updated"))
	s.chanAnswer.ColdAnswer(chanId, resp, err)
}

func (s *Service) listToText(nodes []googlekeep.Node) (text, tts string) {
	var listNodesText []string
	var listNodesTts []string
	sort.Slice(nodes, func(i, j int) bool {
		i64, _ := strconv.ParseInt(nodes[i].SortValue, 10, 0)
		j64, _ := strconv.ParseInt(nodes[j].SortValue, 10, 0)
		return i64 > j64
	})
	for _, node := range nodes {
		if !node.Checked {
			listNodesTts = append(listNodesTts, node.Text)
			listNodesText = append(listNodesText, "○ "+node.Text)
		}
	}
	text = s.loc.Translate("List is empty")
	tts = s.loc.Translate("List is empty")
	if len(listNodesText) > 0 {
		text = fmt.Sprintf(s.loc.Translate("Your list: %s"), "\n"+strings.Join(listNodesText, "\n"))
		tts = fmt.Sprintf(s.loc.Translate("Your list: %s"), " sil<[200]> "+
			strings.Join(listNodesTts, " sil<[100]>, "))
	}
	return
}

func (s *Service) error(chanId int, err error) {
	resp := *alice.NewResponse()
	resp.RespData.Text = fmt.Sprintf(s.loc.Translate(err.Error()))
	resp.RespData.Tts = resp.RespData.Text
	s.chanAnswer.HotAnswer(chanId, resp)
	s.chanAnswer.ColdAnswer(chanId, resp, nil)
}
