package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unacademy/kubernetes-pod-monitor/sessions"
)

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

func (h *Controller) HealthHandler(c *gin.Context) {
	sessions.HealthOrPanic()
	c.JSON(http.StatusOK, gin.H{
		"status": "up",
	})
}
