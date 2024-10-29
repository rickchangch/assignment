package controller

import (
	"assignment-pe/internal/errs"
	"assignment-pe/internal/model"
	"assignment-pe/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SwapController struct {
	swapSvc service.SwapService
}

func NewSwapController(
	swapSvc service.SwapService,
) *SwapController {
	return &SwapController{
		swapSvc: swapSvc,
	}
}

func (ctrl *SwapController) Swap(c *gin.Context) {
	type SwapDto struct {
		UserID string  `json:"userID"`
		Pair   string  `json:"pair"`
		Amount float64 `json:"amount"`
	}
	var body SwapDto
	if err := c.BindJSON(&body); err != nil {
		_ = c.Error(errs.ErrInvalidArgument.New())
		return
	}

	swapInfo := model.SwapInfo{
		UserID: body.UserID,
		Pair:   body.Pair,
		Amount: body.Amount,
	}
	err := ctrl.swapSvc.Swap(c.Request.Context(), swapInfo)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
