package alice

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/sku4/alice-checklist/internal/repository"
	"github.com/sku4/alice-checklist/lang"
	"github.com/sku4/alice-checklist/models/alice"
	"github.com/sku4/alice-checklist/models/alice/channel"
	"time"
)

const (
	RequestTimeOut = 2300 * time.Millisecond
)

//go:generate mockgen -source=alice.go -destination=mocks/alice.go

type ChanAnswer interface {
	HotAnswer(int, alice.Response)
	ColdAnswer(int, alice.Response, error)
	DropAnswer(int)
}

type Service struct {
	checklist      repository.Checklist
	notify         repository.Notify
	chanHotAnswer  map[int]chan alice.Response
	chanColdAnswer map[int]chan interface{}
	loc            *lang.Localize
	recognize      Recognize
	chanAnswer     ChanAnswer
	requestTimeOut time.Duration
}

func NewService(loc *lang.Localize, checklist repository.Checklist, notify repository.Notify) *Service {
	s := &Service{
		checklist:      checklist,
		notify:         notify,
		chanHotAnswer:  make(map[int]chan alice.Response),
		chanColdAnswer: make(map[int]chan interface{}),
		loc:            loc,
		requestTimeOut: RequestTimeOut,
	}
	s.recognize = s
	s.chanAnswer = s
	return s
}

func (s *Service) Command(req alice.Request) (alice.Response, error) {
	var err error
	chanId := genChanId()

	s.chanHotAnswer[chanId] = make(chan alice.Response)
	s.chanColdAnswer[chanId] = make(chan interface{})

	go s.recognize.Recognize(chanId, req)
	resp := <-s.chanHotAnswer[chanId]

	timeout := time.After(s.requestTimeOut)
	select {
	case i := <-s.chanColdAnswer[chanId]:
		switch i.(type) {
		case channel.Response:
			ret := i.(channel.Response)
			resp = ret.Response
			err = ret.Error
		}
	case <-timeout:
		go s.chanAnswer.DropAnswer(chanId)
	}

	return resp, err
}

func (s *Service) HotAnswer(chanId int, resp alice.Response) {
	if _, exist := s.chanHotAnswer[chanId]; exist {
		s.chanHotAnswer[chanId] <- resp
		close(s.chanHotAnswer[chanId])
		delete(s.chanHotAnswer, chanId)
	}
}

func (s *Service) ColdAnswer(chanId int, resp alice.Response, err error) {
	if _, exist := s.chanColdAnswer[chanId]; exist {
		s.processError(&resp, err)
		s.chanColdAnswer[chanId] <- channel.Response{
			Response: resp,
			Error:    err,
		}
		close(s.chanColdAnswer[chanId])
		delete(s.chanColdAnswer, chanId)
	}
}

func (s *Service) DropAnswer(chanId int) {
	var err error

	select {
	case i := <-s.chanColdAnswer[chanId]:
		switch i.(type) {
		case channel.Response:
			ret := i.(channel.Response)
			err = ret.Error
		}
	}

	if err != nil {
		_ = s.notify.Add(err)
	}
}

func (s *Service) processError(resp *alice.Response, err error) {
	if err != nil {
		logrus.Error(err)
		resp.Text = fmt.Sprintf(s.loc.Translate("Invalid response: %s"), err.Error())
		resp.Tts = s.loc.Translate("Invalid response, open chat to find details")
	}
}

var genChanId = func() func() int {
	c := -1
	return func() int {
		c++
		return c
	}
}()
