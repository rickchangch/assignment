package controller

import (
	"assignment-pe/internal/errs"
	"assignment-pe/internal/model"
	"assignment-pe/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PointHistoryController struct {
	pointHistorySrv service.PointHistoryService
}

func NewPointHistoryController(
	pointHistorySrv service.PointHistoryService,
) *PointHistoryController {
	return &PointHistoryController{
		pointHistorySrv: pointHistorySrv,
	}
}

func (ctrl *PointHistoryController) ListPagination(c *gin.Context) {
	userID := c.Param("userID")
	cursor := c.Query("cursor")

	var cursorPtr *string
	if len(cursor) != 0 {
		cursorPtr = &cursor
	}

	size, err := strconv.Atoi(c.Query("size"))
	if err != nil {
		_ = c.Error(errs.ErrInvalidArgument.New())
		return
	}
	if size == 0 {
		size = 100
	}

	pointHistories, err := ctrl.pointHistorySrv.ListPagination(
		c.Request.Context(),
		model.PointHistorySearchCondition{
			UserID: userID,
			Cursor: cursorPtr,
			Size:   size,
		})
	if err != nil {
		_ = c.Error(err)
		return
	}

	type pageVO struct {
		PointHistories []model.PointHistory `json:"pointHistories"`
		Cursor         string               `json:"cursor"`
	}
	result := pageVO{
		PointHistories: pointHistories,
	}
	if len(pointHistories) > 0 {
		result.Cursor = pointHistories[len(pointHistories)-1].ID
	}
	c.JSON(http.StatusOK, result)
}
