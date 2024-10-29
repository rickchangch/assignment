package controller

import (
	"assignment-pe/internal/model"
	"assignment-pe/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserCampaignController struct {
	userCampaignSvc service.UserCampaignService
	userSvc         service.UserService
}

func NewUserCampaignController(
	userCampaignSvc service.UserCampaignService,
	userSvc service.UserService,
) *UserCampaignController {
	return &UserCampaignController{
		userCampaignSvc: userCampaignSvc,
		userSvc:         userSvc,
	}
}

func (ctrl *UserCampaignController) GetUserCampaign(c *gin.Context) {
	address := c.Param("address")
	campaignID := c.Param("campaignID")

	ctx := c.Request.Context()
	user, err := ctrl.userSvc.GetUserByAddress(ctx, address)
	if err != nil {
		_ = c.Error(err)
		return
	}

	type GetUserCampaignVO struct {
		*model.UserCampaign
		Tasks []*model.UserCampaignTask `json:"tasks"`
	}

	userCampaign, err := ctrl.userCampaignSvc.Get(ctx, user.ID, campaignID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	userCampaignTasks, err := ctrl.userCampaignSvc.ListTasks(ctx, user.ID, campaignID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, GetUserCampaignVO{
		UserCampaign: userCampaign,
		Tasks:        userCampaignTasks,
	})
}
