package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/sku4/alice-checklist/lang"
	"github.com/sku4/alice-checklist/models/alice"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}

func newErrorAlice(c *gin.Context, message string, loc lang.Localize) {
	logrus.Error(message)
	resp := alice.NewResponse()
	resp.Text = message
	resp.Tts = loc.Translate("Invalid response, open chat to find details")
	c.JSON(http.StatusOK, resp)
}
