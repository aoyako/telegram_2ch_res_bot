package controller

import (
	"errors"
	"testing"

	mock_storage "github.com/aoyako/telegram_2ch_res_bot/controller/mock"
	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"github.com/aoyako/telegram_2ch_res_bot/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserController_Register(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name   string
		chatID int64
		want   error
	}{
		{
			"Register user",
			123,
			nil,
		},
		{
			"Do not register user",
			123,
			errors.New("User already exists"),
		},
	}

	for _, tt := range tests {
		m := mock_storage.NewMockStorage(ctrl)
		m.MockUser.
			EXPECT().
			Register(gomock.Eq(&logic.User{ChatID: tt.chatID})).
			Return(tt.want)

		ucon := NewUserController(&storage.Storage{
			Info:         m.MockInfo,
			User:         m.MockUser,
			Subscription: m.MockSubscription,
		})

		err := ucon.Register(tt.chatID)

		assert.Equal(tt.want, err)
	}
}

func TestUserController_Unregister(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name   string
		chatID int64
		want   error
	}{
		{
			"Unregister user",
			123,
			nil,
		},
	}

	for _, tt := range tests {
		m := mock_storage.NewMockStorage(ctrl)
		m.MockUser.
			EXPECT().
			Unregister(gomock.Eq(&logic.User{ChatID: tt.chatID})).
			Return(tt.want)

		ucon := NewUserController(&storage.Storage{
			Info:         m.MockInfo,
			User:         m.MockUser,
			Subscription: m.MockSubscription,
		})

		err := ucon.Unregister(tt.chatID)

		assert.Equal(tt.want, err)
	}
}

func TestUserController_GetUsersByPublication(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name string
		arg  *logic.Publication
		want []logic.User
	}{
		{
			name: "Get users",
			arg: &logic.Publication{
				ID: 7,
			},
			want: []logic.User{
				{ID: 1},
				{ID: 2},
				{ID: 3},
			},
		},
		{
			name: "Get no user",
			arg: &logic.Publication{
				ID: 7,
			},
			want: []logic.User{},
		},
	}

	for _, tt := range tests {
		m := mock_storage.NewMockStorage(ctrl)
		m.MockUser.
			EXPECT().
			GetUsersByPublication(gomock.Eq(tt.arg)).
			Return(tt.want, nil)

		ucon := NewUserController(&storage.Storage{
			Info:         m.MockInfo,
			User:         m.MockUser,
			Subscription: m.MockSubscription,
		})

		res, err := ucon.GetUsersByPublication(tt.arg)

		assert.Equal(tt.want, res)
		assert.Nil(err)
	}
}
