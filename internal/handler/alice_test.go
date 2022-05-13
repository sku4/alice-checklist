package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
	"github.com/sku4/alice-checklist/internal/service"
	mock_service "github.com/sku4/alice-checklist/internal/service/mocks"
	"github.com/sku4/alice-checklist/lang"
	"github.com/sku4/alice-checklist/models/alice"
	testify "github.com/stretchr/testify/assert"
	"net/http/httptest"
	"os"
	"testing"
)

const (
	ListJsonPath = "./alice_test/list.json"
)

func TestHandler_aliceRequest(t *testing.T) {
	type mockBehavior func(r *mock_service.MockAlice, req alice.Request)

	loc, err := lang.InitLocalize(lang.Ru)
	if err != nil {
		testify.Fail(t, err.Error())
		return
	}

	tests := []struct {
		name                 string
		inputBody            string
		inputRequest         alice.Request
		inputJson            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputJson: ListJsonPath,
			mockBehavior: func(r *mock_service.MockAlice, req alice.Request) {
				resp := *alice.NewResponse()
				resp.EndSession = true
				resp.Text = "List is empty"
				resp.Tts = "List is empty"
				r.EXPECT().Command(req).Return(resp, nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: `{"version":"1.0","response":{"text":"List is empty","tts":"List is empty",` +
				`"end_session":true}}`,
		},
		{
			name:               "Wrong Input",
			inputBody:          `{"version": 0}`,
			inputRequest:       alice.Request{},
			mockBehavior:       func(r *mock_service.MockAlice, req alice.Request) {},
			expectedStatusCode: 400,
			expectedResponseBody: `{"message":"json: cannot unmarshal number into Go struct field ` +
				`Request.version of type string"}`,
		},
		{
			name:      "Service Error",
			inputJson: ListJsonPath,
			mockBehavior: func(r *mock_service.MockAlice, req alice.Request) {
				resp := *alice.NewResponse()
				r.EXPECT().Command(req).Return(resp, errors.New("something went wrong"))
			},
			expectedStatusCode: 200,
			expectedResponseBody: `{"version":"1.0","response":{"text":"something went wrong",` +
				`"tts":"Invalid response, open chat to find details","end_session":false}}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.inputJson != "" {
				var req alice.Request
				dataReq, err := os.ReadFile(test.inputJson)
				if err != nil {
					testify.Fail(t, err.Error())
					return
				}
				if err = json.Unmarshal(dataReq, &req); err != nil {
					testify.Fail(t, err.Error())
					return
				}
				test.inputBody = string(dataReq)
				test.inputRequest = req
			}

			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockAlice(c)
			test.mockBehavior(repo, test.inputRequest)

			services := &service.Service{Alice: repo}
			handler := Handler{
				services: *services,
				loc:      *loc,
			}

			// Init Endpoint
			r := gin.New()
			r.POST("/cmd", handler.aliceRequest)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/cmd",
				bytes.NewBufferString(test.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
