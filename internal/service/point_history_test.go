package service

import (
	mockRepo "assignment-pe/generate/mock/repo"
	"assignment-pe/internal/model"
	"assignment-pe/internal/repo"
	"context"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_pointHistoryService_BatchInsert(t *testing.T) {
	ctrl := gomock.NewController(t)
	type fields struct {
		pointHistoryRepo *mockRepo.MockPointHistoryRepo
	}
	type args struct {
		ctx   context.Context
		batch []model.PointHistory
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "positive_test",
			fields: fields{
				pointHistoryRepo: func() *mockRepo.MockPointHistoryRepo {
					repo := mockRepo.NewMockPointHistoryRepo(ctrl)
					repo.EXPECT().
						BatchInsert(gomock.Any(), gomock.Any()).
						Return(nil)
					return repo
				}(),
			},
			args: args{
				ctx:   context.TODO(),
				batch: []model.PointHistory{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &pointHistoryService{
				pointHistoryRepo: tt.fields.pointHistoryRepo,
			}
			if err := srv.BatchInsert(tt.args.ctx, tt.args.batch); (err != nil) != tt.wantErr {
				t.Errorf("pointHistoryService.BatchInsert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_pointHistoryService_ListPagination(t *testing.T) {
	ctrl := gomock.NewController(t)
	type fields struct {
		pointHistoryRepo repo.PointHistoryRepo
	}
	type args struct {
		ctx       context.Context
		condition model.PointHistorySearchCondition
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.PointHistory
		wantErr bool
	}{
		{
			name: "positive_test",
			fields: fields{
				pointHistoryRepo: func() *mockRepo.MockPointHistoryRepo {
					repo := mockRepo.NewMockPointHistoryRepo(ctrl)
					repo.EXPECT().
						ListPagination(gomock.Any(), gomock.Any()).
						Return([]model.PointHistory{}, nil)
					return repo
				}(),
			},
			args: args{
				ctx:       context.TODO(),
				condition: model.PointHistorySearchCondition{},
			},
			want:    []model.PointHistory{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &pointHistoryService{
				pointHistoryRepo: tt.fields.pointHistoryRepo,
			}
			got, err := srv.ListPagination(tt.args.ctx, tt.args.condition)
			if (err != nil) != tt.wantErr {
				t.Errorf("pointHistoryService.ListPagination() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pointHistoryService.ListPagination() = %v, want %v", got, tt.want)
			}
		})
	}
}
