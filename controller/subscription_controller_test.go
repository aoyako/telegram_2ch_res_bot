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

func TestSubscriptionController_AddNew(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		chatID              int64
		request             string
		errGetUser          error
		errSaveUser         error
		errSaveSubscribtion error
	}

	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "Add sub",
			args: args{
				chatID:  1,
				request: "a .a \"a\"",
			},
			want: nil,
		},
		{
			name: "Error in getting user",
			args: args{
				chatID:     1,
				request:    "a .a \"a\"",
				errGetUser: errors.New("Error getting user by ChatID"),
			},
			want: errors.New("Error getting user by ChatID"),
		},
		{
			name: "Error in saving user",
			args: args{
				chatID:     1,
				request:    "a .a \"a\"",
				errGetUser: errors.New("Error saving user"),
			},
			want: errors.New("Error saving user"),
		},
		{
			name: "Error in adding subscribtion",
			args: args{
				chatID:     1,
				request:    "a .a \"a\"",
				errGetUser: errors.New("Error adding subscribtion"),
			},
			want: errors.New("Error adding subscribtion"),
		},
	}

	for _, tt := range tests {
		m := mock_storage.NewMockStorage(ctrl)

		scon := NewSubscriptionController(&storage.Storage{
			Info:         m.MockInfo,
			User:         m.MockUser,
			Subscription: m.MockSubscription,
		})

		m.MockUser.
			EXPECT().
			GetUserByChatID(gomock.Eq(tt.args.chatID)).
			Return(&logic.User{ID: 1, ChatID: tt.args.chatID}, tt.args.errGetUser)

		if tt.args.errGetUser == nil {
			if tt.args.errSaveUser == nil {
				m.MockSubscription.
					EXPECT().
					Add(gomock.Eq(&logic.User{ID: 1, ChatID: tt.args.chatID}), gomock.Eq(&logic.Publication{
						Board: "a",
						Type:  ".a",
						Tags:  "\"a\"",
					})).
					Return(tt.args.errSaveSubscribtion)

				if tt.args.errSaveSubscribtion == nil {
					m.MockUser.
						EXPECT().
						Update(gomock.Eq(&logic.User{ID: 1, ChatID: tt.args.chatID, SubsCount: 1})).
						Return(tt.args.errSaveUser)
				}

			}
		}

		err := scon.AddNew(tt.args.chatID, tt.args.request)
		assert.Equal(tt.want, err)
	}
}

func TestSubscriptionController_Create(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		chatID  int64
		request string
		isAdmin bool
	}

	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "Add default sub",
			args: args{
				chatID:  1,
				request: "a .a \"a\"",
				isAdmin: true,
			},
			want: nil,
		},
		{
			name: "User is not an admin",
			args: args{
				chatID:  1,
				request: "a .a \"a\"",
				isAdmin: false,
			},
			want: errors.New("access denied"),
		},
	}

	for _, tt := range tests {
		m := mock_storage.NewMockStorage(ctrl)

		scon := NewSubscriptionController(&storage.Storage{
			Info:         m.MockInfo,
			User:         m.MockUser,
			Subscription: m.MockSubscription,
		})

		m.MockUser.
			EXPECT().
			IsChatAdmin(gomock.Eq(tt.args.chatID)).
			Return(tt.args.isAdmin)

		if tt.args.isAdmin {
			m.MockSubscription.
				EXPECT().
				AddDefault(gomock.Eq(&logic.Publication{
					Board:     "a",
					Type:      ".a",
					Tags:      "\"a\"",
					IsDefault: true,    // Must be default
					Alias:     "\"a\"", // Must have alias
				})).
				Return(tt.want)
		}

		err := scon.Create(tt.args.chatID, tt.args.request)
		assert.Equal(tt.want, err)
	}
}

func TestSubscriptionController_Subscribe(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		chatID      int64
		request     string
		maxInd      uint
		ind         uint
		passConvert bool
	}

	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "Subscribe",
			args: args{
				chatID:      1,
				request:     "1",
				maxInd:      20,
				ind:         0,
				passConvert: true,
			},
			want: nil,
		},
		{
			name: "Request index out of range",
			args: args{
				chatID:      1,
				request:     "100",
				maxInd:      20,
				ind:         100,
				passConvert: true,
			},
			want: errors.New("bad index"),
		},
		{
			name: "Request index out of range - negative",
			args: args{
				chatID:      1,
				request:     "-1",
				maxInd:      20,
				passConvert: true,
			},
			want: errors.New("bad index"),
		},
		{
			name: "Bad request index",
			args: args{
				chatID:  1,
				request: "temp",
				maxInd:  20,
			},
			want: errors.New("bad index"),
		},
	}

	for _, tt := range tests {
		m := mock_storage.NewMockStorage(ctrl)

		scon := NewSubscriptionController(&storage.Storage{
			Info:         m.MockInfo,
			User:         m.MockUser,
			Subscription: m.MockSubscription,
		})

		m.MockUser.
			EXPECT().
			GetUserByChatID(gomock.Eq(tt.args.chatID)).
			Return(&logic.User{ID: 1}, nil)

		wantedPubs := make([]logic.Publication, tt.args.maxInd)
		if tt.args.passConvert {
			for i := range wantedPubs {
				wantedPubs[i].IsDefault = true
				wantedPubs[i].ID = i
			}
			m.MockSubscription.
				EXPECT().
				GetAllDefaultSubs().
				Return(wantedPubs)
		}

		if tt.want == nil {
			m.MockSubscription.
				EXPECT().
				Connect(gomock.Eq(&logic.User{ID: 1}), gomock.Eq(&wantedPubs[tt.args.ind])).
				Return(nil)

			m.MockUser.
				EXPECT().
				Update(gomock.Eq(&logic.User{ID: 1, SubsCount: 1})).
				Return(nil)
		}

		err := scon.Subscribe(tt.args.chatID, tt.args.request)
		assert.Equal(tt.want, err)
	}
}

func TestSubscriptionController_Remove(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		chatID      int64
		request     string
		maxInd      uint
		ind         uint
		passConvert bool
		isDefault   bool
	}

	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "Unsubscribe",
			args: args{
				chatID:      1,
				request:     "1",
				maxInd:      20,
				ind:         0,
				passConvert: true,
			},
			want: nil,
		},
		{
			name: "Request index out of range",
			args: args{
				chatID:      1,
				request:     "100",
				maxInd:      20,
				ind:         100,
				passConvert: true,
			},
			want: errors.New("bad index"),
		},
		{
			name: "Request index out of range - negative",
			args: args{
				chatID:      1,
				request:     "-1",
				maxInd:      20,
				passConvert: true,
			},
			want: errors.New("bad index"),
		},
		{
			name: "Bad request index",
			args: args{
				chatID:  1,
				request: "temp",
				maxInd:  20,
			},
			want: errors.New("bad index"),
		},
	}

	for _, tt := range tests {
		m := mock_storage.NewMockStorage(ctrl)

		scon := NewSubscriptionController(&storage.Storage{
			Info:         m.MockInfo,
			User:         m.MockUser,
			Subscription: m.MockSubscription,
		})

		m.MockUser.
			EXPECT().
			GetUserByChatID(gomock.Eq(tt.args.chatID)).
			Return(&logic.User{ID: 1, SubsCount: tt.args.maxInd}, nil)

		wantedPubs := make([]logic.Publication, tt.args.maxInd)
		for i := range wantedPubs {
			wantedPubs[i].ID = i
		}
		m.MockSubscription.
			EXPECT().
			GetSubsByUser(gomock.Eq(&logic.User{ID: 1, SubsCount: tt.args.maxInd})).
			Return(wantedPubs, nil)

		if tt.want == nil {
			m.MockSubscription.
				EXPECT().
				Disonnect(gomock.Eq(&logic.User{ID: 1, SubsCount: tt.args.maxInd}), gomock.Eq(&wantedPubs[tt.args.ind])).
				Return(nil)

			if !tt.args.isDefault {
				m.MockSubscription.
					EXPECT().
					Remove(gomock.Eq(&wantedPubs[tt.args.ind])).
					Return(nil)
			}

			m.MockUser.
				EXPECT().
				Update(gomock.Eq(&logic.User{ID: 1, SubsCount: tt.args.maxInd - 1})).
				Return(nil)
		}

		err := scon.Remove(tt.args.chatID, tt.args.request)
		assert.Equal(tt.want, err)
	}
}

func TestSubscriptionController_RemoveDefault(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		chatID      int64
		request     string
		maxInd      uint
		ind         uint
		passConvert bool
		isAdmin     bool
		subscribers uint
	}

	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "Remove Default",
			args: args{
				chatID:      1,
				request:     "1",
				maxInd:      20,
				ind:         0,
				passConvert: true,
				isAdmin:     true,
				subscribers: 2,
			},
			want: nil,
		},
		{
			name: "Request index out of range",
			args: args{
				chatID:      1,
				request:     "100",
				maxInd:      20,
				ind:         100,
				passConvert: false,
				isAdmin:     true,
			},
			want: errors.New("bad index"),
		},
		{
			name: "Request index out of range - negative",
			args: args{
				chatID:      1,
				request:     "-1",
				maxInd:      20,
				passConvert: false,
				isAdmin:     true,
			},
			want: errors.New("bad index"),
		},
		{
			name: "Bad request index",
			args: args{
				chatID:  1,
				request: "temp",
				maxInd:  20,
				isAdmin: true,
			},
			want: errors.New("bad index"),
		},
		{
			name: "Not an admin",
			args: args{
				chatID:      1,
				request:     "temp",
				maxInd:      20,
				isAdmin:     false,
				passConvert: true,
			},
			want: errors.New("access denied"),
		},
	}

	for _, tt := range tests {
		m := mock_storage.NewMockStorage(ctrl)

		scon := NewSubscriptionController(&storage.Storage{
			Info:         m.MockInfo,
			User:         m.MockUser,
			Subscription: m.MockSubscription,
		})

		m.MockUser.
			EXPECT().
			IsChatAdmin(gomock.Eq(tt.args.chatID)).
			Return(tt.args.isAdmin)

		if tt.args.isAdmin {
			wantedPubs := make([]logic.Publication, tt.args.maxInd)
			for i := range wantedPubs {
				wantedPubs[i].ID = i
				wantedPubs[i].IsDefault = true
			}
			m.MockSubscription.
				EXPECT().
				GetAllDefaultSubs().
				Return(wantedPubs)

			if tt.args.passConvert {
				subs := make([]logic.User, tt.args.subscribers)
				for i := range subs {
					subs[i].SubsCount = 1
				}
				m.MockUser.
					EXPECT().
					GetUsersByPublication(gomock.Eq(&wantedPubs[tt.args.ind])).
					Return(subs, nil)

				for i := range subs {
					subs[i].SubsCount--
					m.MockUser.
						EXPECT().
						Update(gomock.Eq(&subs[i])).
						Return(nil)
				}
				m.MockSubscription.
					EXPECT().
					Remove(gomock.Eq(&wantedPubs[tt.args.ind])).
					Return(nil)
			}
		}

		err := scon.RemoveDefault(tt.args.chatID, tt.args.request)
		assert.Equal(tt.want, err)
	}
}

func TestSubscriptionController_GetSubsByChatID(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		chatID int64
	}

	tests := []struct {
		name     string
		args     args
		wantErr  error
		wantSubs []logic.Publication
	}{
		{
			name: "Get user's subs",
			args: args{
				chatID: 1,
			},
			wantErr: nil,
			wantSubs: []logic.Publication{
				logic.Publication{
					ID: 1,
				},
				logic.Publication{
					ID: 2,
				},
			},
		},
	}

	for _, tt := range tests {
		m := mock_storage.NewMockStorage(ctrl)

		scon := NewSubscriptionController(&storage.Storage{
			Info:         m.MockInfo,
			User:         m.MockUser,
			Subscription: m.MockSubscription,
		})

		user := &logic.User{ID: 1}

		m.MockUser.
			EXPECT().
			GetUserByChatID(gomock.Eq(tt.args.chatID)).
			Return(user, nil)

		m.MockSubscription.
			EXPECT().
			GetSubsByUser(gomock.Eq(user)).
			Return(tt.wantSubs, nil)

		subs, err := scon.GetSubsByChatID(tt.args.chatID)
		assert.Equal(tt.wantErr, err)
		assert.Equal(tt.wantSubs, subs)
	}
}

func TestSubscriptionController_GetAllSubs(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name     string
		wantSubs []logic.Publication
	}{
		{
			name: "Get all subs",
			wantSubs: []logic.Publication{
				logic.Publication{
					ID: 1,
				},
				logic.Publication{
					ID: 2,
				},
			},
		},
	}

	for _, tt := range tests {
		m := mock_storage.NewMockStorage(ctrl)

		scon := NewSubscriptionController(&storage.Storage{
			Info:         m.MockInfo,
			User:         m.MockUser,
			Subscription: m.MockSubscription,
		})

		m.MockSubscription.
			EXPECT().
			GetAllSubs().
			Return(tt.wantSubs)

		subs := scon.GetAllSubs()
		assert.Equal(tt.wantSubs, subs)
	}
}

func TestSubscriptionController_GetAllDefaultSubs(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name     string
		wantSubs []logic.Publication
	}{
		{
			name: "Get all default subs",
			wantSubs: []logic.Publication{
				logic.Publication{
					ID:        1,
					IsDefault: true,
				},
				logic.Publication{
					ID:        2,
					IsDefault: true,
				},
			},
		},
	}

	for _, tt := range tests {
		m := mock_storage.NewMockStorage(ctrl)

		scon := NewSubscriptionController(&storage.Storage{
			Info:         m.MockInfo,
			User:         m.MockUser,
			Subscription: m.MockSubscription,
		})

		m.MockSubscription.
			EXPECT().
			GetAllDefaultSubs().
			Return(tt.wantSubs)

		subs := scon.GetAllDefaultSubs()
		assert.Equal(tt.wantSubs, subs)
	}
}

func Test_parseRequest(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name            string
		request         string
		wantPublication *logic.Publication
		wantError       error
	}{
		{
			name:      "Normal",
			request:   "a .b.c \"D\"|\"e\"",
			wantError: nil,
			wantPublication: &logic.Publication{
				Board: "a",
				Type:  ".b.c",
				Tags:  "\"D\"|\"e\"",
			},
		},
		{
			name:      "Normal with one args",
			request:   "a .b \"C\"",
			wantError: nil,
			wantPublication: &logic.Publication{
				Board: "a",
				Type:  ".b",
				Tags:  "\"C\"",
			},
		},
		{
			name:      "Normal with many args",
			request:   "a .b \"C\"|\"De\"&\"F\"|!\"G\"",
			wantError: nil,
			wantPublication: &logic.Publication{
				Board: "a",
				Type:  ".b",
				Tags:  "\"C\"|\"De\"&\"F\"|!\"G\"",
			},
		},
		{
			name:      "No tags",
			request:   "a .b",
			wantError: errors.New("bad request"),
		},
		{
			name:      "No formats",
			request:   "a \"C\"",
			wantError: errors.New("bad request"),
		},
		{
			name:      "No board",
			request:   ".b \"C\"",
			wantError: errors.New("bad request"),
		},
		{
			name:      "Empty tags",
			request:   "a .b \"\"",
			wantError: errors.New("bad request"),
		},
		{
			name:      "Wrong type",
			request:   "a b.c \"D\"|\"e\"",
			wantError: errors.New("bad request"),
		},
	}

	for _, tt := range tests {
		res, err := parseRequest(tt.request)
		assert.Equal(tt.wantPublication, res)
		assert.Equal(tt.wantError, err)
	}
}

func Test_parseRequestAlias(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name            string
		request         string
		wantPublication *logic.Publication
		wantError       error
	}{
		{
			name:      "Normal",
			request:   "a .b.c \"D\"|\"e\" Default",
			wantError: nil,
			wantPublication: &logic.Publication{
				Board: "a",
				Type:  ".b.c",
				Tags:  "\"D\"|\"e\"",
				Alias: "Default",
			},
		},
		{
			name:      "Normal with long alias",
			request:   "a .b \"C\" Default name",
			wantError: nil,
			wantPublication: &logic.Publication{
				Board: "a",
				Type:  ".b",
				Tags:  "\"C\"",
				Alias: "Default name",
			},
		},
		{
			name:      "Normal with many args",
			request:   "a .b \"C\"|\"De\"&\"F\"|!\"G\"",
			wantError: nil,
			wantPublication: &logic.Publication{
				Board: "a",
				Type:  ".b",
				Tags:  "\"C\"|\"De\"&\"F\"|!\"G\"",
				Alias: "\"C\"|\"De\"&\"F\"|!\"G\"",
			},
		},
		{
			name:      "No tags",
			request:   "a .b Default",
			wantError: errors.New("bad request"),
		},
		{
			name:      "No formats",
			request:   "a \"C\" Default",
			wantError: errors.New("bad request"),
		},
		{
			name:      "No board",
			request:   ".b \"C\" Default",
			wantError: errors.New("bad request"),
		},
		{
			name:      "Empty tags",
			request:   "a .b \"\" Default",
			wantError: errors.New("bad request"),
		},
		{
			name:      "Wrong type",
			request:   "a b.c \"D\"|\"e\" Default",
			wantError: errors.New("bad request"),
		},
	}

	for _, tt := range tests {
		res, err := parseRequestAlias(tt.request)
		assert.Equal(tt.wantPublication, res)
		assert.Equal(tt.wantError, err)
	}
}
