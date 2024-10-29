package service

import (
	mockRepo "assignment-pe/generate/mock/repo"
	"assignment-pe/internal/model"
	"context"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_userCampaignService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	type fields struct {
		userCampaignRepo *mockRepo.MockUserCampaignRepo
	}
	type args struct {
		ctx        context.Context
		userID     string
		campaignID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.UserCampaign
		wantErr bool
	}{
		{
			name: "positive_test",
			fields: fields{
				userCampaignRepo: func() *mockRepo.MockUserCampaignRepo {
					repo := mockRepo.NewMockUserCampaignRepo(ctrl)
					repo.EXPECT().
						Get(gomock.Any(), gomock.Any(), gomock.Any()).
						Return(&model.UserCampaign{}, nil)
					return repo
				}(),
			},
			args: args{
				ctx:        context.TODO(),
				userID:     "csf4eh04v3c4kmbslpt0",
				campaignID: "esf4eh04v3c4kmbslpt0",
			},
			want:    &model.UserCampaign{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &userCampaignService{
				userCampaignRepo: tt.fields.userCampaignRepo,
			}
			got, err := srv.Get(tt.args.ctx, tt.args.userID, tt.args.campaignID)
			if (err != nil) != tt.wantErr {
				t.Errorf("userCampaignService.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userCampaignService.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userCampaignService_ListTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	type fields struct {
		userCampaignRepo *mockRepo.MockUserCampaignRepo
	}
	type args struct {
		ctx        context.Context
		userID     string
		campaignID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*model.UserCampaignTask
		wantErr bool
	}{
		{
			name: "positive_test",
			fields: fields{
				userCampaignRepo: func() *mockRepo.MockUserCampaignRepo {
					repo := mockRepo.NewMockUserCampaignRepo(ctrl)
					repo.EXPECT().
						ListTasks(gomock.Any(), gomock.Any(), gomock.Any()).
						Return([]*model.UserCampaignTask{}, nil)
					return repo
				}(),
			},
			args: args{
				ctx:        context.TODO(),
				userID:     "csf4eh04v3c4kmbslpt0",
				campaignID: "esf4eh04v3c4kmbslpt0",
			},
			want:    []*model.UserCampaignTask{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &userCampaignService{
				userCampaignRepo: tt.fields.userCampaignRepo,
			}
			got, err := srv.ListTasks(tt.args.ctx, tt.args.userID, tt.args.campaignID)
			if (err != nil) != tt.wantErr {
				t.Errorf("userCampaignService.ListTasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userCampaignService.ListTasks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userCampaignService_UpsertAmount(t *testing.T) {
	ctrl := gomock.NewController(t)
	type fields struct {
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
		want    []model.PointHistory
		want1   []int64
		wantErr bool
	}{
		{
			name: "positive_test",
			fields: fields{
				userCampaignRepo: func() *mockRepo.MockUserCampaignRepo {
					repo := mockRepo.NewMockUserCampaignRepo(ctrl)
					repo.EXPECT().
						UpsertAmountToCampaignAndTask(gomock.Any(), gomock.Any()).
						Return([]model.PointHistory{}, []int64{}, nil)
					return repo
				}(),
			},
			args: args{
				ctx:      context.TODO(),
				swapInfo: model.SwapInfo{},
			},
			want:    []model.PointHistory{},
			want1:   []int64{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &userCampaignService{
				userCampaignRepo: tt.fields.userCampaignRepo,
			}
			got, got1, err := srv.UpsertAmount(tt.args.ctx, tt.args.swapInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("userCampaignService.UpsertAmount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userCampaignService.UpsertAmount() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("userCampaignService.UpsertAmount() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
