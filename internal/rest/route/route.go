package route

import (
	"assignment-pe/internal/errs"
	"assignment-pe/internal/rest/controller"
	"assignment-pe/internal/rest/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Route interface {
	Index()
}

type route struct {
	engine           *gin.Engine
	mw               middleware.Middleware
	campaignCtrl     *controller.CampaignController
	userCampaignCtrl *controller.UserCampaignController
	pointHistoryCtrl *controller.PointHistoryController
	swapCtrl         *controller.SwapController
}

func NewRoute(
	engine *gin.Engine,
	mw middleware.Middleware,
	campaignCtrl *controller.CampaignController,
	userCampaignCtrl *controller.UserCampaignController,
	pointHistoryCtrl *controller.PointHistoryController,
	swapCtrl *controller.SwapController,
) Route {
	return &route{
		engine:           engine,
		mw:               mw,
		campaignCtrl:     campaignCtrl,
		userCampaignCtrl: userCampaignCtrl,
		pointHistoryCtrl: pointHistoryCtrl,
		swapCtrl:         swapCtrl,
	}
}

func (r *route) Index() {
	r.engine.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, errs.ErrNotFound)
	})

	root := r.engine.Group("/")
	root.Use(r.mw.Logger(), r.mw.AccessLog(), r.mw.ErrorHandler(), r.mw.PanicRecovery(), r.mw.Tx())
	{
		campaign := root.Group("/")
		{
			campaign.GET("/v1/addresses/:address/campaigns/:campaignID", r.userCampaignCtrl.GetUserCampaign)
			campaign.GET("/v1/campaigns/:campaignID/leaderboard", r.campaignCtrl.GetCampaignLeaderboard)
		}

		pointHistory := root.Group("/")
		{
			pointHistory.GET("/v1/users/:userID/pointHistories", r.pointHistoryCtrl.ListPagination)
		}

		swap := root.Group("/")
		{
			swap.POST("/v1/swap", r.swapCtrl.Swap)
		}
	}
}
