package controller

import (
	"assignment-pe/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CampaignController struct {
	campaignSrv service.CampaignService
}

func NewCampaignController(
	campaignSrv service.CampaignService,
) *CampaignController {
	return &CampaignController{
		campaignSrv: campaignSrv,
	}
}

func (ctrl *CampaignController) GetCampaignLeaderboard(c *gin.Context) {
	campaignID := c.Param("campaignID")

	leaderBoard, err := ctrl.campaignSrv.GetLeaderboard(c.Request.Context(), campaignID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, leaderBoard)
}
