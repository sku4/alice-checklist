package handler

import (
	"github.com/gin-contrib/pprof"
	"github.com/sku4/alice-checklist/internal/service"
	"github.com/sku4/alice-checklist/lang"

	"github.com/gin-gonic/gin"

	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "github.com/sku4/alice-checklist/docs"
)

type Handler struct {
	services service.Service
	loc      lang.Localize
}

func NewHandler(loc *lang.Localize, services *service.Service) *Handler {
	return &Handler{
		services: *services,
		loc:      *loc,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	list := router.Group("/cmd")
	{
		list.POST("/", h.aliceRequest)
	}

	pprof.Register(router)

	return router
}
