package service

import (
	"assignment-pe/internal/cx"
	"assignment-pe/internal/model"
	"context"

	"github.com/go-redis/redis/v8"
)

type SwapService interface {
	Swap(ctx context.Context, swapInfo model.SwapInfo) error
}

type swapService struct {
	redis               *redis.Client
	campaignService     CampaignService
	userCampaignService UserCampaignService
	pointHistoryService PointHistoryService
}

func NewSwapService(
	redis *redis.Client,
	campaignService CampaignService,
	userCampaignService UserCampaignService,
	pointHistoryService PointHistoryService,
) SwapService {
	return &swapService{
		redis:               redis,
		campaignService:     campaignService,
		userCampaignService: userCampaignService,
		pointHistoryService: pointHistoryService,
	}
}

func (srv *swapService) Swap(
	ctx context.Context,
	swapInfo model.SwapInfo,
) error {
	logger := cx.GetLogger(ctx)

	activeCampaignTasks, err := srv.campaignService.ListActiveCampaignsWithTasks(ctx, swapInfo)
	if err != nil {
		return err
	}
	swapInfo.ActiveCampaignWithTasks = activeCampaignTasks

	pointHistories, userCampaignAmounts, err := srv.userCampaignService.UpsertAmount(ctx, swapInfo)
	if err != nil {
		return err
	}

	err = srv.pointHistoryService.BatchInsert(ctx, pointHistories)
	if err != nil {
		return err
	}

	go func() {
		ctx := context.Background()
		for i := range swapInfo.ActiveCampaignWithTasks {
			campaignID := swapInfo.ActiveCampaignWithTasks[i].CampaignID

			// TODO: ttl
			_, err := srv.redis.ZAdd(ctx, campaignID, &redis.Z{
				Score:  float64(userCampaignAmounts[i]),
				Member: swapInfo.UserID,
			}).Result()

			if err != nil {
				logger.WithError(err).Errorf("redis zadd failed, key: %s", campaignID)
				_, err = srv.redis.Del(ctx, campaignID).Result()
				if err != nil {
					logger.WithError(err).Errorf("redis del failed, key: %s", campaignID)
				}
			}
		}
	}()

	return nil
}
