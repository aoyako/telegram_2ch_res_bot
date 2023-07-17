package dvach_test

import (
	"testing"

	"github.com/aoyako/telegram_2ch_res_bot/dvach"
	mock_dvach "github.com/aoyako/telegram_2ch_res_bot/dvach/mock/requester"
	"github.com/stretchr/testify/assert"

	"github.com/aoyako/telegram_2ch_res_bot/controller"
	mock_telegram "github.com/aoyako/telegram_2ch_res_bot/dvach/mock/sender"
	"github.com/aoyako/telegram_2ch_res_bot/logic"
	mock_controller "github.com/aoyako/telegram_2ch_res_bot/telegram/mock"
	"github.com/golang/mock/gomock"
)

func TestAPIWorkerDvach_InitiateSending(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		publications     []logic.Publication
		users            [][]logic.User
		requestURL       dvach.RequestURL
		boards           []string
		expectAllThreads []dvach.ListResponse
		expectThreadData [][]dvach.ThreadData
		lastTimestamp    uint64
		filesToSend      []string
		urlFilesToSend   []string
		threadsToProcess [][]string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Initiate test",
			args: args{
				publications: []logic.Publication{
					{ID: 1, Board: "a", Type: ".img", Tags: "\"abc\""},
				},
				users: [][]logic.User{
					{
						{ID: 1, ChatID: 123},
					},
				},
				requestURL: dvach.RequestURL{
					AllThreadsURL: "/board/%s",
					ThreadURL:     "/board/%s/thread/%s",
					ResourceURL:   "/res/%s",
				},
				boards: []string{"a"},
				expectAllThreads: []dvach.ListResponse{
					{
						Board: "a",
						Threads: []dvach.Thread{
							{
								Comment: "Default comment abc",
								ID:      123,
							},
						},
					},
				},
				expectThreadData: [][]dvach.ThreadData{
					[]dvach.ThreadData{
						dvach.ThreadData{
							ThreadPosts: []dvach.ThreadPost{
								dvach.ThreadPost{
									Posts: []dvach.Post{
										dvach.Post{
											Comment:   "Default post",
											Timestamp: 123,
										},
										dvach.Post{
											Comment: "Post with files",
											Files: []dvach.File{
												dvach.File{
													Name: "default_file.png",
													Path: "filepath.png",
													Size: 100,
												},
											},
											Timestamp: 124,
										},
									},
								},
							},
						},
					},
				},
				lastTimestamp:  124,
				filesToSend:    []string{"filepath.png"},
				urlFilesToSend: []string{"/res/filepath.png"},
				threadsToProcess: [][]string{
					[]string{"123"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := mock_telegram.NewMockSender(ctrl)
			cm := mock_controller.NewMockController(ctrl)
			sm := mock_dvach.NewMockRequester(ctrl)

			awdv := dvach.NewAPIWorkerDvach(&controller.Controller{
				User:         cm.MockUser,
				Subscription: cm.MockSubscription,
				Info:         cm.MockInfo,
			}, tm, sm)

			cm.MockSubscription.
				EXPECT().
				GetAllSubs().
				Return(tt.args.publications)

			for i := range tt.args.publications {
				cm.MockUser.
					EXPECT().
					GetUsersByPublication(gomock.Eq(&tt.args.publications[i])).
					Return(tt.args.users[i], nil)
			}

			for i := range tt.args.boards {
				sm.
					EXPECT().
					GetAllThreads(gomock.Eq(tt.args.boards[i])).
					Return(tt.args.expectAllThreads[i])
			}

			for i := range tt.args.expectThreadData {
				for j := range tt.args.threadsToProcess[i] {
					sm.
						EXPECT().
						GetThread(gomock.Eq(tt.args.boards[i]), gomock.Eq(tt.args.threadsToProcess[i][j])).
						Return(tt.args.expectThreadData[i][j])

					cm.MockInfo.
						EXPECT().
						GetLastTimestamp().
						Return(uint64(0))
				}
			}

			for i := range tt.args.filesToSend {
				sm.
					EXPECT().
					GetResourceURL(gomock.Eq(tt.args.filesToSend[i])).
					Return(tt.args.urlFilesToSend[i])

				receivers := make([]*logic.User, len(tt.args.users[i]))
				for j := range receivers {
					receivers[j] = &tt.args.users[i][j]
				}
				tm.
					EXPECT().
					Send(gomock.Eq(receivers),
						gomock.Eq(tt.args.urlFilesToSend[i]),
						gomock.Eq(tt.args.threadsToProcess[i][0]),
					).AnyTimes()
			}

			cm.MockInfo.
				EXPECT().
				SetLastTimestamp(gomock.Eq(tt.args.lastTimestamp))

			awdv.InitiateSending()
		})
	}
}

func Test_CheckFileExtension(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		filename string
		req      dvach.SourceType
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Have, pass",
			args: args{
				filename: "default.png",
				req: dvach.SourceType{
					Image: true,
				},
			},
			want: true,
		},
		{
			name: "Have, not pass",
			args: args{
				filename: "default.png",
				req: dvach.SourceType{
					Webm: true,
				},
			},
			want: false,
		},
		{
			name: "Not Have, not pass",
			args: args{
				filename: "default.png",
				req:      dvach.SourceType{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := dvach.CheckFileExtension(tt.args.filename, tt.args.req)
			assert.Equal(tt.want, result)
		})
	}
}

func Test_ParseKeywords(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		s      string
		inputs []string
	}
	tests := []struct {
		name string
		args args
		want []bool
	}{
		{
			name: "Test validation",
			args: args{
				s: `"ac"|"ab"&"bc"|"a"&!"z"`,
				inputs: []string{
					`acz`,
					`abc`,
					`ab`,
					`abz`,
					`bc`,
					`bcz`,
				},
			},
			want: []bool{
				true,
				true,
				true,
				false,
				false,
				false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := dvach.ParseKeywords(tt.args.s)
			for i := range tt.args.inputs {
				result := validator(tt.args.inputs[i])
				assert.Equal(tt.want[i], result, tt.args.inputs[i])
			}
		})
	}
}

func Test_ParseTypes(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want dvach.SourceType
	}{
		{
			name: "Image",
			args: args{
				".img",
			},
			want: dvach.SourceType{
				Image: true,
			},
		},
		{
			name: "Webm and gif",
			args: args{
				".webm.gif",
			},
			want: dvach.SourceType{
				Webm: true,
				Gif:  true,
			},
		},
		{
			name: "Gif and webm",
			args: args{
				".gif.webm",
			},
			want: dvach.SourceType{
				Webm: true,
				Gif:  true,
			},
		},
		{
			name: "All",
			args: args{
				".img.webm.gif",
			},
			want: dvach.SourceType{
				Image: true,
				Webm:  true,
				Gif:   true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := dvach.ParseTypes(tt.args.s)
			assert.Equal(tt.want, result, tt.args.s)
		})
	}
}
