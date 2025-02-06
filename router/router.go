package router

import (
	"context"
	"service1/controllers"

	"github.com/omniful/go_commons/http"
)

func Initialize(ctx context.Context, s *http.Server) (err error) {
	OrderV1 := s.Engine.Group("/api/V1")
	OrderV1.POST("/order/bulk", controllers.UploadOrders)
	return
}
