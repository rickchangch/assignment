package service

import (
	"assignment-pe/internal/model"
	"context"
	"reflect"
	"testing"

	mockRepo "assignment-pe/generate/mock/repo"

	"github.com/golang/mock/gomock"
)

func Test_userService_GetUserByAddress(t *testing.T) {
	ctrl := gomock.NewController(t)
	type fields struct {
		userRepo *mockRepo.MockUserRepo
	}
	type args struct {
		ctx     context.Context
		address string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.User
		wantErr bool
	}{
		{
			name: "positive_test",
			fields: fields{
				userRepo: func() *mockRepo.MockUserRepo {
					repo := mockRepo.NewMockUserRepo(ctrl)
					repo.EXPECT().
						GetByAddress(gomock.Any(), gomock.Any()).
						Return(&model.User{}, nil)
					return repo
				}(),
			},
			args: args{
				ctx:     context.TODO(),
				address: "0xeeeeeeeeeeeeeeeeeeeee",
			},
			want:    &model.User{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &userService{
				userRepo: tt.fields.userRepo,
			}
			got, err := srv.GetUserByAddress(tt.args.ctx, tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("userService.GetUserByAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userService.GetUserByAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
