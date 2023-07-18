package telegram

import (
	"errors"
	"testing"

	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"github.com/stretchr/testify/assert"

	"github.com/aoyako/telegram_2ch_res_bot/controller"
	"github.com/aoyako/telegram_2ch_res_bot/downloader"
	mock_controller "github.com/aoyako/telegram_2ch_res_bot/telegram/mock"
	mock_sender "github.com/aoyako/telegram_2ch_res_bot/telegram/mock/bot"
	mock_downloader "github.com/aoyako/telegram_2ch_res_bot/telegram/mock/downloader"
	"github.com/golang/mock/gomock"
	telebot "gopkg.in/tucnak/telebot.v2"
)

func Test_start(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		chatID  int64
		willAdd bool
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Register User",
			args{
				123,
				true,
			},
			HelpMessage,
		},
		{
			"User exists",
			args{
				123,
				false,
			},
			HelpMessage,
		},
	}

	for _, tt := range tests {
		cm := mock_controller.NewMockController(ctrl)
		dm := mock_downloader.NewMockLoader(ctrl)
		sm := mock_sender.NewMockMessageSender(ctrl)

		controller := &controller.Controller{
			User:         cm.MockUser,
			Info:         cm.MockInfo,
			Subscription: cm.MockSubscription,
		}
		downloader := &downloader.Downloader{Loader: dm}
		bot := &TgBot{
			Controller: controller,
			Downloader: downloader,
			Bot:        sm,
		}

		handler := start(bot)

		returnMessage := errors.New("User already exists")
		if tt.args.willAdd {
			returnMessage = nil
		}
		cm.MockUser.
			EXPECT().
			Register(gomock.Eq(tt.args.chatID)).
			Return(returnMessage)
		sm.
			EXPECT().
			Send(nil, HelpMessage, telebot.ModeMarkdown).
			Return(&telebot.Message{
				Chat: &telebot.Chat{
					ID: int64(tt.args.chatID),
				},
				Text: tt.want,
			}, nil)

		message := telebot.Message{
			Chat: &telebot.Chat{
				ID: int64(tt.args.chatID),
			},
		}

		handler(&message)
	}
}

func Test_subs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		chatID int64
		subs   []logic.Publication
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "List subs",
			args: args{
				chatID: 1,
				subs: []logic.Publication{
					{ID: 1, Board: "a", Tags: "\"b\"", Type: ".g"},
					{ID: 2, Alias: "Default"},
				},
			},
			want: "Your subs:\n1: /a .g \"b\"\n2: Default",
		},
	}
	for _, tt := range tests {
		cm := mock_controller.NewMockController(ctrl)
		dm := mock_downloader.NewMockLoader(ctrl)
		sm := mock_sender.NewMockMessageSender(ctrl)

		controller := &controller.Controller{
			User:         cm.MockUser,
			Info:         cm.MockInfo,
			Subscription: cm.MockSubscription,
		}
		downloader := &downloader.Downloader{Loader: dm}
		bot := &TgBot{
			Controller: controller,
			Downloader: downloader,
			Bot:        sm,
		}

		handler := subs(bot)

		cm.MockSubscription.
			EXPECT().
			GetSubsByChatID(gomock.Eq(tt.args.chatID)).
			Return(tt.args.subs, nil)
		sm.
			EXPECT().
			Send(nil, tt.want).
			Return(&telebot.Message{
				Chat: &telebot.Chat{
					ID: int64(tt.args.chatID),
				},
				Text: tt.want,
			}, nil)

		message := telebot.Message{
			Chat: &telebot.Chat{
				ID: int64(tt.args.chatID),
			},
		}

		handler(&message)
	}
}

func Test_list(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		chatID int64
		subs   []logic.Publication
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "List subs",
			args: args{
				chatID: 1,
				subs: []logic.Publication{
					{ID: 1, Board: "a", Tags: "\"b\"", Type: ".g"},
					{ID: 2, Alias: "Default"},
				},
			},
			want: "Available subs:\n1: /a .g \"b\"\n2: Default",
		},
	}
	for _, tt := range tests {
		cm := mock_controller.NewMockController(ctrl)
		dm := mock_downloader.NewMockLoader(ctrl)
		sm := mock_sender.NewMockMessageSender(ctrl)

		controller := &controller.Controller{
			User:         cm.MockUser,
			Info:         cm.MockInfo,
			Subscription: cm.MockSubscription,
		}
		downloader := &downloader.Downloader{Loader: dm}
		bot := &TgBot{
			Controller: controller,
			Downloader: downloader,
			Bot:        sm,
		}

		handler := list(bot)

		cm.MockSubscription.
			EXPECT().
			GetAllDefaultSubs().
			Return(tt.args.subs)
		sm.
			EXPECT().
			Send(nil, tt.want).
			Return(&telebot.Message{
				Chat: &telebot.Chat{
					ID: int64(tt.args.chatID),
				},
				Text: tt.want,
			}, nil)

		message := telebot.Message{
			Chat: &telebot.Chat{
				ID: int64(tt.args.chatID),
			},
		}

		handler(&message)
	}
}

func Test_cleverList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		chatID int64
		subs   []logic.Publication
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "List subs",
			args: args{
				chatID: 1,
				subs: []logic.Publication{
					{ID: 1, Board: "a", Tags: "\"b\"", Type: ".g"},
					{ID: 2, Alias: "Default"},
				},
			},
			want: "Available subs:\n1: /a .g \"b\"\n2: /  ",
		},
	}
	for _, tt := range tests {
		cm := mock_controller.NewMockController(ctrl)
		dm := mock_downloader.NewMockLoader(ctrl)
		sm := mock_sender.NewMockMessageSender(ctrl)

		controller := &controller.Controller{
			User:         cm.MockUser,
			Info:         cm.MockInfo,
			Subscription: cm.MockSubscription,
		}
		downloader := &downloader.Downloader{Loader: dm}
		bot := &TgBot{
			Controller: controller,
			Downloader: downloader,
			Bot:        sm,
		}

		handler := cleverList(bot)

		cm.MockSubscription.
			EXPECT().
			GetAllDefaultSubs().
			Return(tt.args.subs)
		sm.
			EXPECT().
			Send(nil, tt.want).
			Return(&telebot.Message{
				Chat: &telebot.Chat{
					ID: int64(tt.args.chatID),
				},
				Text: tt.want,
			}, nil)

		message := telebot.Message{
			Chat: &telebot.Chat{
				ID: int64(tt.args.chatID),
			},
		}

		handler(&message)
	}
}

func Test_help(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		chatID int64
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Help",
			args{
				123,
			},
			HelpMessage,
		},
	}

	for _, tt := range tests {
		cm := mock_controller.NewMockController(ctrl)
		dm := mock_downloader.NewMockLoader(ctrl)
		sm := mock_sender.NewMockMessageSender(ctrl)

		controller := &controller.Controller{
			User:         cm.MockUser,
			Info:         cm.MockInfo,
			Subscription: cm.MockSubscription,
		}
		downloader := &downloader.Downloader{Loader: dm}
		bot := &TgBot{
			Controller: controller,
			Downloader: downloader,
			Bot:        sm,
		}

		handler := help(bot)

		sm.
			EXPECT().
			Send(nil, HelpMessage, telebot.ModeMarkdown).
			Return(&telebot.Message{
				Chat: &telebot.Chat{
					ID: int64(tt.args.chatID),
				},
				Text: tt.want,
			}, nil)

		message := telebot.Message{
			Chat: &telebot.Chat{
				ID: int64(tt.args.chatID),
			},
		}

		handler(&message)
	}
}

func Test_create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		chatID        int64
		request       string
		subscribtion  logic.Publication
		arg           string
		failArgsCheck bool
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Create",
			args{
				chatID:  123,
				request: "/create_default a .b \"c\"",
				subscribtion: logic.Publication{
					Board:     "a",
					Type:      ".b",
					Tags:      "\"c\"",
					IsDefault: false,
				},
				arg: "a .b \"c\"",
			},
			"OK",
		},
		{
			"Do not create",
			args{
				chatID:  123,
				request: "/create_default",
				subscribtion: logic.Publication{
					Board:     "a",
					Type:      ".b",
					Tags:      "\"c\"",
					IsDefault: false,
				},
				arg:           "",
				failArgsCheck: true,
			},
			"Bad request",
		},
	}

	for _, tt := range tests {
		cm := mock_controller.NewMockController(ctrl)
		dm := mock_downloader.NewMockLoader(ctrl)
		sm := mock_sender.NewMockMessageSender(ctrl)

		controller := &controller.Controller{
			User:         cm.MockUser,
			Info:         cm.MockInfo,
			Subscription: cm.MockSubscription,
		}
		downloader := &downloader.Downloader{Loader: dm}
		bot := &TgBot{
			Controller: controller,
			Downloader: downloader,
			Bot:        sm,
		}

		handler := create(bot)

		if !tt.args.failArgsCheck {
			cm.MockSubscription.
				EXPECT().
				AddNew(gomock.Eq(tt.args.chatID), gomock.Eq(tt.args.arg)).
				Return(nil)
		}
		sm.
			EXPECT().
			Send(nil, tt.want).
			Return(&telebot.Message{
				Chat: &telebot.Chat{
					ID: int64(tt.args.chatID),
				},
				Text: tt.want,
			}, nil)

		message := telebot.Message{
			Chat: &telebot.Chat{
				ID: int64(tt.args.chatID),
			},
			Text: tt.args.request,
		}

		handler(&message)
	}
}

func Test_subscribe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		chatID        int64
		request       string
		arg           string
		failArgsCheck bool
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Subscribe",
			args{
				chatID:  123,
				request: "/subscribe 1",
				arg:     "1",
			},
			"OK",
		},
		{
			"Do not subscribe",
			args{
				chatID:        123,
				request:       "/subscribe",
				arg:           "",
				failArgsCheck: true,
			},
			"Bad request",
		},
	}

	for _, tt := range tests {
		cm := mock_controller.NewMockController(ctrl)
		dm := mock_downloader.NewMockLoader(ctrl)
		sm := mock_sender.NewMockMessageSender(ctrl)

		controller := &controller.Controller{
			User:         cm.MockUser,
			Info:         cm.MockInfo,
			Subscription: cm.MockSubscription,
		}
		downloader := &downloader.Downloader{Loader: dm}
		bot := &TgBot{
			Controller: controller,
			Downloader: downloader,
			Bot:        sm,
		}

		handler := subscribe(bot)

		if !tt.args.failArgsCheck {
			cm.MockSubscription.
				EXPECT().
				Subscribe(gomock.Eq(tt.args.chatID), gomock.Eq(tt.args.arg)).
				Return(nil)
		}
		sm.
			EXPECT().
			Send(nil, tt.want).
			Return(&telebot.Message{
				Chat: &telebot.Chat{
					ID: int64(tt.args.chatID),
				},
				Text: tt.want,
			}, nil)

		message := telebot.Message{
			Chat: &telebot.Chat{
				ID: int64(tt.args.chatID),
			},
			Text: tt.args.request,
		}

		handler(&message)
	}
}

func Test_createDefault(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		chatID        int64
		request       string
		subscribtion  logic.Publication
		arg           string
		failArgsCheck bool
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Create",
			args{
				chatID:  123,
				request: "/create a .b \"c\" Default",
				subscribtion: logic.Publication{
					Board:     "a",
					Type:      ".b",
					Tags:      "\"c\"",
					IsDefault: true,
				},
				arg: "a .b \"c\" Default",
			},
			"OK",
		},
		{
			"Do not create",
			args{
				chatID:  123,
				request: "/create",
				subscribtion: logic.Publication{
					Board:     "a",
					Type:      ".b",
					Tags:      "\"c\"",
					IsDefault: true,
				},
				arg:           "",
				failArgsCheck: true,
			},
			"Bad request",
		},
	}

	for _, tt := range tests {
		cm := mock_controller.NewMockController(ctrl)
		dm := mock_downloader.NewMockLoader(ctrl)
		sm := mock_sender.NewMockMessageSender(ctrl)

		controller := &controller.Controller{
			User:         cm.MockUser,
			Info:         cm.MockInfo,
			Subscription: cm.MockSubscription,
		}
		downloader := &downloader.Downloader{Loader: dm}
		bot := &TgBot{
			Controller: controller,
			Downloader: downloader,
			Bot:        sm,
		}

		handler := createDefault(bot)

		if !tt.args.failArgsCheck {
			cm.MockSubscription.
				EXPECT().
				Create(gomock.Eq(tt.args.chatID), gomock.Eq(tt.args.arg)).
				Return(nil)
		}
		sm.
			EXPECT().
			Send(nil, tt.want).
			Return(&telebot.Message{
				Chat: &telebot.Chat{
					ID: int64(tt.args.chatID),
				},
				Text: tt.want,
			}, nil)

		message := telebot.Message{
			Chat: &telebot.Chat{
				ID: int64(tt.args.chatID),
			},
			Text: tt.args.request,
		}

		handler(&message)
	}
}

func Test_deleleSub(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		chatID        int64
		request       string
		arg           string
		failArgsCheck bool
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Delete",
			args{
				chatID:  123,
				request: "/rm 1",
				arg:     "1",
			},
			"OK",
		},
		{
			"Do not delete",
			args{
				chatID:        123,
				request:       "/rm",
				arg:           "",
				failArgsCheck: true,
			},
			"Bad request",
		},
	}

	for _, tt := range tests {
		cm := mock_controller.NewMockController(ctrl)
		dm := mock_downloader.NewMockLoader(ctrl)
		sm := mock_sender.NewMockMessageSender(ctrl)

		controller := &controller.Controller{
			User:         cm.MockUser,
			Info:         cm.MockInfo,
			Subscription: cm.MockSubscription,
		}
		downloader := &downloader.Downloader{Loader: dm}
		bot := &TgBot{
			Controller: controller,
			Downloader: downloader,
			Bot:        sm,
		}

		handler := deleleSub(bot)

		if !tt.args.failArgsCheck {
			cm.MockSubscription.
				EXPECT().
				Remove(gomock.Eq(tt.args.chatID), gomock.Eq(tt.args.arg)).
				Return(nil)
		}
		sm.
			EXPECT().
			Send(nil, tt.want).
			Return(&telebot.Message{
				Chat: &telebot.Chat{
					ID: int64(tt.args.chatID),
				},
				Text: tt.want,
			}, nil)

		message := telebot.Message{
			Chat: &telebot.Chat{
				ID: int64(tt.args.chatID),
			},
			Text: tt.args.request,
		}

		handler(&message)
	}
}

func Test_removeDefault(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		chatID        int64
		request       string
		arg           string
		failArgsCheck bool
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Delete",
			args{
				chatID:  123,
				request: "/rm_default 1",
				arg:     "1",
			},
			"OK",
		},
		{
			"Do not delete",
			args{
				chatID:        123,
				request:       "/rm_default",
				arg:           "",
				failArgsCheck: true,
			},
			"Bad request",
		},
	}

	for _, tt := range tests {
		cm := mock_controller.NewMockController(ctrl)
		dm := mock_downloader.NewMockLoader(ctrl)
		sm := mock_sender.NewMockMessageSender(ctrl)

		controller := &controller.Controller{
			User:         cm.MockUser,
			Info:         cm.MockInfo,
			Subscription: cm.MockSubscription,
		}
		downloader := &downloader.Downloader{Loader: dm}
		bot := &TgBot{
			Controller: controller,
			Downloader: downloader,
			Bot:        sm,
		}

		handler := removeDefault(bot)

		if !tt.args.failArgsCheck {
			cm.MockSubscription.
				EXPECT().
				RemoveDefault(gomock.Eq(tt.args.chatID), gomock.Eq(tt.args.arg)).
				Return(nil)
		}
		sm.
			EXPECT().
			Send(nil, tt.want).
			Return(&telebot.Message{
				Chat: &telebot.Chat{
					ID: int64(tt.args.chatID),
				},
				Text: tt.want,
			}, nil)

		message := telebot.Message{
			Chat: &telebot.Chat{
				ID: int64(tt.args.chatID),
			},
			Text: tt.args.request,
		}

		handler(&message)
	}
}

func Test_parseCommand(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		cmd string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
	}{
		{
			name: "OK",
			args: args{
				cmd: "/command_name args",
			},
			want:    "args",
			wantErr: nil,
		},
		{
			name: "Bad request",
			args: args{
				cmd: "/command_name",
			},
			want:    "",
			wantErr: errors.New("bad request"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := parseCommand(tt.args.cmd)
			assert.Equal(tt.want, res)
			assert.Equal(tt.wantErr, err)
		})
	}
}

func Test_marshallSubs(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		subs         []logic.Publication
		displayAlias bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "List subs with alias",
			args: args{
				displayAlias: true,
				subs: []logic.Publication{
					{ID: 1, Board: "a", Tags: "\"b\"", Type: ".g"},
					{ID: 2, Alias: "Default"},
				},
			},
			want: "\n1: /a .g \"b\"\n2: Default",
		},
		{
			name: "List subs without alias",
			args: args{
				displayAlias: false,
				subs: []logic.Publication{
					{ID: 1, Board: "a", Tags: "\"b\"", Type: ".g"},
					{ID: 2, Alias: "Default"},
				},
			},
			want: "\n1: /a .g \"b\"\n2: /  ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := marshallSubs(tt.args.subs, tt.args.displayAlias)
			assert.Equal(tt.want, res)
		})
	}
}

func Test_marshallSub(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		sub logic.Publication
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "List sub",
			args: args{
				sub: logic.Publication{ID: 1, Board: "a", Tags: "\"b\"", Type: ".g"},
			},
			want: "/a .g \"b\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := marshallSub(tt.args.sub)
			assert.Equal(tt.want, res)
		})
	}
}
