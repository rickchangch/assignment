package service

import (
	"assignment-pe/internal/model"
	"assignment-pe/internal/repo"
	"context"
)

// nolint: lll
type UserCampaignService interface {
	Get(ctx context.Context, userID string, campaignID string) (*model.UserCampaign, error)
	ListTasks(ctx context.Context, userID string, campaignID string) ([]*model.UserCampaignTask, error)
	UpsertAmount(ctx context.Context, swapInfo model.SwapInfo) ([]model.PointHistory, []int64, error)
}

type userCampaignService struct {
	userCampaignRepo repo.UserCampaignRepo
}

func NewUserCampaignService(
	userCampaignRepo repo.UserCampaignRepo,
) UserCampaignService {
	return &userCampaignService{
		userCampaignRepo: userCampaignRepo,
	}
}

func (srv *userCampaignService) Get(
	ctx context.Context,
	userID string,
	campaignID string,
) (*model.UserCampaign, error) {
	userCampaign, err := srv.userCampaignRepo.Get(ctx, userID, campaignID)
	if err != nil {
		return nil, err
	}
	return userCampaign, nil
}

func (srv *userCampaignService) ListTasks(
	ctx context.Context,
	userID string,
	campaignID string,
) ([]*model.UserCampaignTask, error) {
	userCampaignTasks, err := srv.userCampaignRepo.ListTasks(ctx, userID, campaignID)
	if err != nil {
		return nil, err
	}
	return userCampaignTasks, nil
}

func (srv *userCampaignService) UpsertAmount(
	ctx context.Context,
	swapInfo model.SwapInfo,
) ([]model.PointHistory, []int64, error) {
	pointHistories, userCampaignAmounts, err := srv.userCampaignRepo.UpsertAmountToCampaignAndTask(ctx, swapInfo)
	if err != nil {
		return nil, nil, err
	}
	return pointHistories, userCampaignAmounts, nil
}
