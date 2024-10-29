package service

import (
	mockRepo "assignment-pe/generate/mock/repo"
	"assignment-pe/internal/errs"
	"assignment-pe/internal/model"
	"context"
	"reflect"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
)

func Test_campaignService_ListActiveCampaignsWithTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	type fields struct {
		redis            *redis.Client // TODO: mock
		campaignRepo     *mockRepo.MockCampaignRepo
		userCampaignRepo *mockRepo.MockUserCampaignRepo
	}
	type args struct {
		ctx      context.Context
		swapInfo model.SwapInfo
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.ActiveCampaignWithTask
		wantErr error
	}{
		{
			name: "positive_test",
			fields: fields{
				redis: nil,
				campaignRepo: func() *mockRepo.MockCampaignRepo {
					repo := mockRepo.NewMockCampaignRepo(ctrl)
					repo.EXPECT().
						ListActiveCampaignWithTask(gomock.Any(), gomock.Any()).
						Return([]model.ActiveCampaignWithTask{}, nil)
					return repo
				}(),
				userCampaignRepo: mockRepo.NewMockUserCampaignRepo(ctrl),
			},
			args: args{
				ctx: context.Background(),
				swapInfo: model.SwapInfo{
					UserID:                  "test_user",
					Pair:                    "USDC/ETH",
					Amount:                  100,
					ActiveCampaignWithTasks: nil,
				},
			},
			want:    []model.ActiveCampaignWithTask{},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &campaignService{
				redis:            tt.fields.redis,
				campaignRepo:     tt.fields.campaignRepo,
				userCampaignRepo: tt.fields.userCampaignRepo,
			}
			got, err := srv.ListActiveCampaignsWithTasks(tt.args.ctx, tt.args.swapInfo)
			if err != nil {
				wantErr, ok := tt.wantErr.(errs.AppError)
				if !ok || !wantErr.Is(err) {
					t.Errorf("campaignService.ListActiveCampaignsWithTasks() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("campaignService.ListActiveCampaignsWithTasks() = %v, want %v", got, tt.want)
			}
		})
	}
}
