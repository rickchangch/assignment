package service

import (
	"assignment-pe/internal/cx"
	"assignment-pe/internal/model"
	"assignment-pe/internal/repo"
	"assignment-pe/internal/util"
	"context"

	"github.com/go-redis/redis/v8"
)

// nolint: lll
type CampaignService interface {
	ListActiveCampaignsWithTasks(ctx context.Context, swapInfo model.SwapInfo) ([]model.ActiveCampaignWithTask, error)
	GetLeaderboard(ctx context.Context, campaignID string) ([]model.LeaderboardRow, error)
}

type campaignService struct {
	redis            *redis.Client
	campaignRepo     repo.CampaignRepo
	userCampaignRepo repo.UserCampaignRepo
}

func NewCampaignService(
	redis *redis.Client,
	campaignRepo repo.CampaignRepo,
	userCampaignRepo repo.UserCampaignRepo,
) CampaignService {
	return &campaignService{
		redis:            redis,
		campaignRepo:     campaignRepo,
		userCampaignRepo: userCampaignRepo,
	}
}

func (srv *campaignService) ListActiveCampaignsWithTasks(
	ctx context.Context,
	swapInfo model.SwapInfo,
) ([]model.ActiveCampaignWithTask, error) {
	result, err := srv.campaignRepo.ListActiveCampaignWithTask(ctx, swapInfo)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (srv *campaignService) GetLeaderboard(
	ctx context.Context,
	campaignID string,
) ([]model.LeaderboardRow, error) {
	logger := cx.GetLogger(ctx)
	var leaderboard []model.LeaderboardRow

	zsets, err := srv.redis.ZRevRangeWithScores(ctx, campaignID, 0, 10).Result()
	if err == nil {
		// Fetch from Redis
		leaderboard = util.Map(zsets, func(in redis.Z) model.LeaderboardRow {
			return model.LeaderboardRow{
				UserID: in.Member.(string),
				Points: int64(in.Score),
			}
		})
	} else {
		logger.Debugf("campaign's leaderboard not found in redis, key: %v", campaignID)

		// Fetch from RDB
		leaderboard, err = srv.userCampaignRepo.GetLeaderboard(ctx, campaignID)
		if err != nil {
			return nil, err
		}

		go func() {
			zsets := util.Map(leaderboard, func(in model.LeaderboardRow) *redis.Z {
				return &redis.Z{
					Score:  float64(in.Points),
					Member: in.UserID,
				}
			})
			ctx := context.Background()

			_, err := srv.redis.ZAdd(ctx, campaignID, zsets...).Result() // TODO: ttl
			if err != nil {
				logger.WithError(err).Errorf("redis zadd failed, key: %s", campaignID)
				_, err = srv.redis.Del(ctx, campaignID).Result()
				if err != nil {
					logger.WithError(err).Errorf("redis del failed, key: %s", campaignID)
				}
			}
		}()
	}

	return leaderboard, nil
}
