package usecase

import (
	"errors"
	"project/domain/entity"
	mock "project/mock/mockrepo"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListUsers(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	adminRepo := mock.NewMockAdminRepository(ctrl)
	adminUseCase := NewAdmin(adminRepo)

	type args struct {
		pages int
		limit int
	}

	tests := []struct {
		name       string
		input      args
		beforeTest func(*mock.MockAdminRepository)
		want       []entity.User
		wantErr    error
	}{
		{
			name:  "Listing users",
			input: args{pages: 1, limit: 1},
			beforeTest: func(mar *mock.MockAdminRepository) {
				mar.EXPECT().GetAllUsers(0, 1).Times(1).Return([]entity.User{{
					Id:         1,
					Name:       "Abhi",
					Email:      "abhishek@gmail.com",
					Phone:      "1234567890",
					Permission: true,
				}}, nil)
			},
			want: []entity.User{{
				Id:         1,
				Name:       "Abhi",
				Email:      "abhishek@gmail.com",
				Phone:      "1234567890",
				Permission: true,
			}},
			wantErr: nil,
		},
		{
			name:  "Error in listing",
			input: args{pages: 1, limit: 1},
			beforeTest: func(mar *mock.MockAdminRepository) {
				mar.EXPECT().GetAllUsers(0, 1).Times(1).Return(nil, errors.New("error in fetching user list"))
			},
			want:    nil,
			wantErr: errors.New("error in fetching user list"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.beforeTest(adminRepo)

			got, err := adminUseCase.ExecuteUsersList(tt.input.pages, tt.input.limit)
			assert.Equal(t, err, tt.wantErr)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("adminUseCase_ExecutUsersList()=%v want %v", got, tt.want)
			}
		})
	}
}

func TestTooglePermission(t *testing.T) {
	ctrl := gomock.NewController(t)

	adminRepo := mock.NewMockAdminRepository(ctrl)
	adminUseCase := NewAdmin(adminRepo)

	tests := []struct {
		name       string
		input      int
		beforeTest func(*mock.MockAdminRepository, int)
		wantErr    error
	}{
		{
			name:  "block user",
			input: 1,
			beforeTest: func(mar *mock.MockAdminRepository, id int) {

				gomock.InOrder(
					mar.EXPECT().GetById(id).Times(1).Return(&entity.User{Id: id, Permission: true}, nil),
					mar.EXPECT().Update(gomock.Any()).Times(1).Return(nil),
				)
			},
			wantErr: nil,
		},
		{
			name:  "getBYID error",
			input: 2,
			beforeTest: func(mar *mock.MockAdminRepository, id int) {
				gomock.InOrder(
					mar.EXPECT().GetById(id).Times(1).Return(nil, errors.New("fetch error")),
				)
			},
			wantErr: errors.New("fetch error"),
		},
		{
			name:  "update error",
			input: 3,
			beforeTest: func(mar *mock.MockAdminRepository, id int) {
				gomock.InOrder(
					mar.EXPECT().GetById(id).Times(1).Return(&entity.User{Id: id, Permission: true}, nil),
					mar.EXPECT().Update(gomock.Any()).Times(1).Return(errors.New("user permission toggle failed")),
				)
			},
			wantErr: errors.New("user permission toggle failed"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.beforeTest(adminRepo, tt.input)

			err := adminUseCase.ExecuteTogglePermission(tt.input)

			assert.Equal(t, err, tt.wantErr)
		})
	}
}

func TestUserSearch(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	adminRepo := mock.NewMockAdminRepository(ctrl)
	adminUseCase := NewAdmin(adminRepo)

	tests := []struct {
		name        string
		page, limit int
		search      string
		beforeTest  func(mock.MockAdminRepository, int, int, string)
		want        []entity.User
		wantErr     error
	}{
		{
			name:   "succesful user search",
			page:   1,
			limit:  1,
			search: "Abhi",
			beforeTest: func(mar mock.MockAdminRepository, offset, limit int, search string) {
				mar.EXPECT().GetUsersBySearch(offset, limit, search).Times(1).Return([]entity.User{
					{Id: 1, Name: "Abhishek"},
				}, nil)
			},
			want: []entity.User{
				{Id: 1, Name: "Abhishek"},
			},
			wantErr: nil,
		},
		{
			name:   "error in search",
			page:   1,
			limit:  1,
			search: "Abhi",
			beforeTest: func(mar mock.MockAdminRepository, offset, limit int, search string) {
				mar.EXPECT().GetUsersBySearch(offset, limit, search).Times(1).Return(nil, errors.New("error in user search"))
			},
			want:    nil,
			wantErr: errors.New("error in user search"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			offset := (tt.page - 1) * tt.limit
			tt.beforeTest(*adminRepo, offset, tt.limit, tt.search)

			users, err := adminUseCase.ExecuteUserSearch(tt.page, tt.limit, tt.search)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, users)
		})
	}
}

func TestExecuteStocklessPr(t *testing.T){
	ctrl:=gomock.NewController(t)
	defer ctrl.Finish()

	adminRepo:=mock.NewMockAdminRepository(ctrl)
	adminUseCase:=NewAdmin(adminRepo)

	tests:=[]struct{
		name string
		beforeTest func(mock.MockAdminRepository)
		want *[]entity.Inventory
		wantErr error
	}{
		{
			name:"Succes",
			beforeTest: func(mar mock.MockAdminRepository){
				mar.EXPECT().GetstocklessProducts().Times(1).Return(&[]entity.Inventory{
					{ProductId: 1,Quantity: 0},
				},nil)
			},
			want: &[]entity.Inventory{
				{ProductId: 1,Quantity: 0},
			},
			wantErr: nil,
		},
		{
			name: "error",
			beforeTest: func(mar mock.MockAdminRepository){
				mar.EXPECT().GetstocklessProducts().Times(1).Return(nil,errors.New("error fetching"))
			},
			want: nil,
			wantErr: errors.New("error fetching"),
		},
	}
	for _,tt:=range tests{
		t.Run(tt.name,func(t *testing.T){
			tt.beforeTest(*adminRepo)

			product,err:=adminUseCase.ExecuteStocklessProducts()
			assert.Equal(t,tt.wantErr,err)
			assert.Equal(t,tt.want,product)
		})
	}
}