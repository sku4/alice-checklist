package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sku4/alice-checklist/models/alice"
	"net/http"
)

// @Summary Webhook to Alice skill
// @Tags Alice
// @Description Get answer by webhook alice command
// @ID alice-request
// @Accept  json
// @Produce  json
// @Param request body alice.Request true "Body request"
// @Success 200 {object} alice.Response
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} alice.Response
// @Router /cmd [post]
func (h *Handler) aliceRequest(c *gin.Context) {
	var req alice.Request

	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.services.Alice.Command(req)
	if err != nil {
		newErrorAlice(c, err.Error(), h.loc)
		return
	}

	c.JSON(http.StatusOK, resp)
}
