package alice

import (
	"errors"
	"github.com/golang/mock/gomock"
	mock_repository "github.com/sku4/alice-checklist/internal/repository/mocks"
	"github.com/sku4/alice-checklist/lang"
	model "github.com/sku4/alice-checklist/models/alice"
	"github.com/stretchr/testify/assert"
	testify "github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type mockRecognize struct {
	chanAnswer   ChanAnswer
	hotResponse  model.Response
	coldResponse model.Response
	coldError    error
	sleep        time.Duration
}

func newMockRecognize(chanAnswer ChanAnswer) *mockRecognize {
	return &mockRecognize{chanAnswer: chanAnswer}
}

func (r *mockRecognize) Recognize(chanId int, _ model.Request) {
	r.chanAnswer.HotAnswer(chanId, r.hotResponse)
	time.Sleep(r.sleep)
	r.chanAnswer.ColdAnswer(chanId, r.coldResponse, r.coldError)
}

func TestService_Command(t *testing.T) {
	type recognizeBehavior func() (respHot, respCold model.Response, coldError error)

	loc, err := lang.InitLocalize(lang.Ru)
	if err != nil {
		testify.Fail(t, err.Error())
		return
	}

	tests := []struct {
		name              string
		chanId            int
		req               model.Request
		sleep             time.Duration
		requestTimeOut    time.Duration
		recognizeBehavior recognizeBehavior
		want              model.Response
		wantErr           assert.ErrorAssertionFunc
	}{
		{
			name:   "Hot response",
			chanId: 0,
			req:    model.Request{},
			sleep:  1 * time.Millisecond,
			recognizeBehavior: func() (hotResp, coldResp model.Response, coldError error) {
				hotResp = *model.NewResponse()
				hotResp.RespData.Text = "Hot response"
				hotResp.RespData.Tts = hotResp.RespData.Text
				return
			},
			want: model.Response{
				Version: model.Version,
				RespData: model.RespData{
					Text: "Hot response",
					Tts:  "Hot response",
				},
			},
		},
		{
			name:   "Cold response",
			chanId: 0,
			req:    model.Request{},
			sleep:  0 * time.Millisecond,
			recognizeBehavior: func() (hotResp, coldResp model.Response, coldError error) {
				coldResp = *model.NewResponse()
				coldResp.RespData.Text = "Cold response"
				coldResp.RespData.Tts = coldResp.RespData.Text
				return
			},
			want: model.Response{
				Version: model.Version,
				RespData: model.RespData{
					Text: "Cold response",
					Tts:  "Cold response",
				},
			},
		},
		{
			name:           "Cold response error",
			chanId:         0,
			req:            model.Request{},
			sleep:          0 * time.Millisecond,
			requestTimeOut: 1 * time.Millisecond,
			recognizeBehavior: func() (hotResp, coldResp model.Response, coldError error) {
				coldResp = *model.NewResponse()
				coldError = errors.New("something went wrong")
				return
			},
			want: model.Response{
				Version: model.Version,
				RespData: model.RespData{
					Text: loc.Translate("Invalid response: something went wrong"),
					Tts:  loc.Translate("Invalid response, open chat to find details"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repoChecklist := mock_repository.NewMockChecklist(c)
			repoNotify := mock_repository.NewMockNotify(c)
			s := NewService(loc, repoChecklist, repoNotify)
			s.requestTimeOut = tt.requestTimeOut

			mr := newMockRecognize(s.chanAnswer)
			mr.hotResponse, mr.coldResponse, mr.coldError = tt.recognizeBehavior()
			mr.sleep = tt.sleep
			s.recognize = mr

			got, _ := s.Command(tt.req)
			assert.Equalf(t, tt.want, got, "Command(%v)", tt.req)
		})
	}
}
