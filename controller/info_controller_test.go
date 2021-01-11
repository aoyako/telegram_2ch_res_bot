package controller

import (
	"testing"

	mock_storage "github.com/aoyako/telegram_2ch_res_bot/controller/mock"
	"github.com/aoyako/telegram_2ch_res_bot/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestInfoController_GetLastTimestamp(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name string
		want uint64
	}{
		{
			"Get timestamp",
			123,
		},
	}

	for _, tt := range tests {
		m := mock_storage.NewMockStorage(ctrl)
		m.MockInfo.
			EXPECT().
			GetLastTimestamp().
			Return(tt.want)

		icon := NewInfoController(&storage.Storage{
			Info:         m.MockInfo,
			User:         m.MockUser,
			Subscription: m.MockSubscription,
		})

		tsmp := icon.GetLastTimestamp()

		assert.Equal(tt.want, tsmp)
	}
}

func TestInfoController_SetLastTimestamp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name         string
		before       uint64
		want         uint64
		shouldChange bool
	}{
		{
			name:         "Set timestamp",
			before:       123,
			want:         124,
			shouldChange: true,
		},
		{
			name:         "Do not set timestamp",
			before:       123,
			want:         122,
			shouldChange: false,
		},
	}

	for _, tt := range tests {
		m := mock_storage.NewMockStorage(ctrl)
		icon := NewInfoController(&storage.Storage{
			Info:         m.MockInfo,
			User:         m.MockUser,
			Subscription: m.MockSubscription,
		})

		m.MockInfo.
			EXPECT().
			GetLastTimestamp().
			Return(tt.before)

		if tt.shouldChange {
			m.MockInfo.
				EXPECT().
				SetLastTimestamp(gomock.Eq(tt.want))
		}

		icon.SetLastTimestamp(tt.want)
	}
}
