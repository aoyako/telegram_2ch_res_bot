package storage

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserMock struct {
	storage *UserPostgres
	mock    sqlmock.Sqlmock
}

func (mock *UserMock) BeforeEach(t *testing.T) {
	var db *sql.DB
	var err error

	db, mocked, err := sqlmock.New()
	mock.mock = mocked
	assert.Nil(t, err)

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.Nil(t, err)

	mock.storage = NewUserPostgres(gdb, &InitDatabase{Admin: []int64{1}})
}

func (mock *UserMock) AfterEach(t *testing.T) {
	err := mock.mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestNewUserPostgres(t *testing.T) {
	assert := assert.New(t)

	db, _, err := sqlmock.New()
	assert.Nil(err)
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.Nil(err)

	config := InitDatabase{Admin: []int64{1}}

	infop := NewUserPostgres(gdb, &config)

	assert.Equal(gdb, infop.db, "Equal db instances")
	assert.Equal(config, *infop.cfg, "Equal config instances")
}

func TestUserPostgres_Register(t *testing.T) {
	assert := assert.New(t)
	dbmock := UserMock{}

	type args struct {
		user     *logic.User
		wantUser bool
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Create user",
			args: args{
				user: &logic.User{
					ID:     7,
					ChatID: 7,
				},
				wantUser: false,
			},
		},
		{
			name: "Do not create user",
			args: args{
				user: &logic.User{
					ID:     7,
					ChatID: 7,
				},
				wantUser: true,
			},
			err: errors.New("User already exists"),
		},
		{
			name: "Create admin",
			args: args{
				user: &logic.User{
					ID:     7,
					ChatID: 1,
				},
				wantUser: false,
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbmock.BeforeEach(t)

			userStorage := dbmock.storage
			userInst := tt.args.user

			rows := sqlmock.NewRows([]string{"count"})
			if tt.args.wantUser {
				rows.AddRow(1)
			} else {
				rows.AddRow(0)
			}

			const sqlSelectUser = `SELECT count(1) FROM "users" WHERE chat_id = $1`
			const sqlInsertUser = `INSERT INTO "users" ("chat_id","subs_count","id") VALUES ($1,$2,$3) RETURNING "id"`
			const sqlIsertAdmin = `INSERT INTO "admins" ("user_id") VALUES ($1) RETURNING "id"`

			dbmock.mock.ExpectQuery(regexp.QuoteMeta(sqlSelectUser)).
				WithArgs(userInst.ChatID).
				WillReturnRows(rows)

			if !tt.args.wantUser {
				dbmock.mock.ExpectQuery(regexp.QuoteMeta(sqlInsertUser)).
					WithArgs(userInst.ChatID, userInst.SubsCount, userInst.ID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tt.args.user.ID))
			}

			if userInst.ChatID == 1 {
				dbmock.mock.ExpectQuery(regexp.QuoteMeta(sqlIsertAdmin)).
					WithArgs(tt.args.user.ID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			}

			tstp := userStorage.Register(tt.args.user)
			assert.Equal(tt.err, tstp)

			dbmock.AfterEach(t)
		})
	}
}

func TestUserPostgres_Unregister(t *testing.T) {
	assert := assert.New(t)
	dbmock := UserMock{}

	type args struct {
		user *logic.User
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Delete user",
			args: args{
				user: &logic.User{
					ID:     7,
					ChatID: 7,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbmock.BeforeEach(t)

			userStorage := dbmock.storage
			userInst := tt.args.user

			const sqlDeleteUser = `DELETE FROM "users" WHERE "users"."id" = $1`

			dbmock.mock.ExpectExec(regexp.QuoteMeta(sqlDeleteUser)).
				WithArgs(userInst.ID).
				WillReturnResult(sqlmock.NewResult(int64(tt.args.user.ID), 1))

			tstp := userStorage.Unregister(tt.args.user)
			assert.Equal(tt.err, tstp)

			dbmock.AfterEach(t)
		})
	}
}

func TestUserPostgres_GetUserByChatID(t *testing.T) {
	assert := assert.New(t)
	dbmock := UserMock{}

	type args struct {
		ChatID int64
	}
	tests := []struct {
		name     string
		args     args
		wantUser *logic.User
		err      error
	}{
		{
			name: "Get user",
			args: args{
				ChatID: 7,
			},
			wantUser: &logic.User{
				ID:     8,
				ChatID: 7,
			},
			err: nil,
		},
		{
			name: "No user",
			args: args{
				ChatID: 7,
			},
			wantUser: nil,
			err:      errors.New("No user found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbmock.BeforeEach(t)

			userStorage := dbmock.storage
			wantUser := tt.wantUser

			const sqlCountUser = `SELECT count(1) FROM "users" WHERE chat_id = $1`
			const sqlSelectUser = `SELECT * FROM "users" WHERE chat_id = $1 ORDER BY "users"."id" LIMIT 1`

			userRows := sqlmock.NewRows([]string{"id", "chat_id", "subs_count"})
			countRows := sqlmock.NewRows([]string{"count"})
			if wantUser != nil {
				userRows.AddRow(wantUser.ID, wantUser.ChatID, wantUser.SubsCount)
				countRows.AddRow(1)
			} else {
				countRows.AddRow(0)
			}

			dbmock.mock.ExpectQuery(regexp.QuoteMeta(sqlCountUser)).
				WithArgs(tt.args.ChatID).
				WillReturnRows(countRows)

			if wantUser != nil {
				dbmock.mock.ExpectQuery(regexp.QuoteMeta(sqlSelectUser)).
					WithArgs(tt.args.ChatID).
					WillReturnRows(userRows)
			}

			tstp, err := userStorage.GetUserByChatID(tt.args.ChatID)

			assert.Equal(tt.wantUser, tstp)
			assert.Equal(tt.err, err)

			dbmock.AfterEach(t)
		})
	}
}

func TestUserPostgres_Update(t *testing.T) {
	assert := assert.New(t)
	dbmock := UserMock{}

	type args struct {
		user *logic.User
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Update user",
			args: args{
				user: &logic.User{
					ID:     8,
					ChatID: 7,
				},
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbmock.BeforeEach(t)

			userStorage := dbmock.storage
			userInst := tt.args.user

			const sqlDeleteUser = `UPDATE "users" SET "chat_id"=$1,"subs_count"=$2 WHERE "id" = $3`

			dbmock.mock.ExpectExec(regexp.QuoteMeta(sqlDeleteUser)).
				WithArgs(userInst.ChatID, userInst.SubsCount, userInst.ID).
				WillReturnResult(sqlmock.NewResult(int64(tt.args.user.ID), 1))

			tstp := userStorage.Update(tt.args.user)
			assert.Equal(tt.err, tstp)

			dbmock.AfterEach(t)
		})
	}
}

func TestUserPostgres_GetUserByID(t *testing.T) {
	assert := assert.New(t)
	dbmock := UserMock{}

	type args struct {
		ID int64
	}
	tests := []struct {
		name     string
		args     args
		wantUser *logic.User
		err      error
	}{
		{
			name: "Get user",
			args: args{
				ID: 8,
			},
			wantUser: &logic.User{
				ID:     8,
				ChatID: 7,
			},
			err: nil,
		},
		{
			name: "No user",
			args: args{
				ID: 7,
			},
			wantUser: nil,
			err:      errors.New("No user found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbmock.BeforeEach(t)

			userStorage := dbmock.storage
			wantUser := tt.wantUser

			const sqlCountUser = `SELECT count(1) FROM "users" WHERE id = $1`
			const sqlSelectUser = `SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT 1`

			userRows := sqlmock.NewRows([]string{"id", "chat_id", "subs_count"})
			countRows := sqlmock.NewRows([]string{"count"})
			if wantUser != nil {
				userRows.AddRow(wantUser.ID, wantUser.ChatID, wantUser.SubsCount)
				countRows.AddRow(1)
			} else {
				countRows.AddRow(0)
			}

			dbmock.mock.ExpectQuery(regexp.QuoteMeta(sqlCountUser)).
				WithArgs(tt.args.ID).
				WillReturnRows(countRows)

			if wantUser != nil {
				dbmock.mock.ExpectQuery(regexp.QuoteMeta(sqlSelectUser)).
					WithArgs(tt.args.ID).
					WillReturnRows(userRows)
			}

			tstp, err := userStorage.GetUserByID(tt.args.ID)

			assert.Equal(tt.wantUser, tstp)
			assert.Equal(tt.err, err)

			dbmock.AfterEach(t)
		})
	}
}

func TestUserPostgres_GetUsersByPublication(t *testing.T) {
	assert := assert.New(t)
	dbmock := UserMock{}

	type args struct {
		publication *logic.Publication
	}
	tests := []struct {
		name      string
		args      args
		wantUsers []logic.User
		err       error
	}{
		{
			name: "Get users by pub",
			args: args{
				publication: &logic.Publication{
					ID: 77,
				},
			},
			wantUsers: []logic.User{
				logic.User{
					ID:     8,
					ChatID: 7,
				},
				logic.User{
					ID:     7,
					ChatID: 8,
				},
			},
			err: nil,
		},
		{
			name: "No users",
			args: args{
				publication: &logic.Publication{
					ID: 77,
				},
			},
			wantUsers: []logic.User{},
			err:       nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbmock.BeforeEach(t)

			userStorage := dbmock.storage
			wantUsers := tt.wantUsers

			const sqlSelectUsers = `SELECT 
				"users"."id","users"."chat_id","users"."subs_count" FROM "users" JOIN "user_subscribtion"
				ON "user_subscribtion"."user_id" = "users"."id" AND "user_subscribtion"."publication_id" = $1`

			userRows := sqlmock.NewRows([]string{"id", "chat_id", "subs_count"})
			for i := range wantUsers {
				user := wantUsers[i]
				userRows.AddRow(user.ID, user.ChatID, user.SubsCount)
			}

			dbmock.mock.ExpectQuery(regexp.QuoteMeta(sqlSelectUsers)).
				WithArgs(tt.args.publication.ID).
				WillReturnRows(userRows)

			tstp, err := userStorage.GetUsersByPublication(tt.args.publication)

			assert.Equal(tt.wantUsers, tstp)
			assert.Equal(tt.err, err)

			dbmock.AfterEach(t)
		})
	}
}

func TestUserPostgres_IsUserAdmin(t *testing.T) {
	assert := assert.New(t)
	dbmock := UserMock{}

	type args struct {
		user *logic.User
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Admin",
			args: args{
				user: &logic.User{
					ID:     8,
					ChatID: 7,
				},
			},
			want: true,
		},
		{
			name: "Not Admin",
			args: args{
				user: &logic.User{
					ID:     8,
					ChatID: 7,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbmock.BeforeEach(t)

			userStorage := dbmock.storage
			userInst := tt.args.user

			const sqlCountUser = `SELECT count(1) FROM "admins" WHERE user_id = $1`

			countRows := sqlmock.NewRows([]string{"count"})
			if tt.want {
				countRows.AddRow(1)
			} else {
				countRows.AddRow(0)
			}

			dbmock.mock.ExpectQuery(regexp.QuoteMeta(sqlCountUser)).
				WithArgs(userInst.ID).
				WillReturnRows(countRows)

			tstp := userStorage.IsUserAdmin(userInst)

			assert.Equal(tt.want, tstp)

			dbmock.AfterEach(t)
		})
	}
}

func TestUserPostgres_IsChatAdmin(t *testing.T) {
	assert := assert.New(t)
	dbmock := UserMock{}

	type args struct {
		ChatID int64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Admin",
			args: args{
				ChatID: 7,
			},
			want: true,
		},
		{
			name: "Not Admin",
			args: args{
				ChatID: 7,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbmock.BeforeEach(t)

			userStorage := dbmock.storage

			const sqlCountUser = `SELECT count(1) FROM "users" WHERE chat_id = $1`
			const sqlSelectUser = `SELECT * FROM "users" WHERE chat_id = $1 ORDER BY "users"."id" LIMIT 1`
			const sqlCountAdmin = `SELECT count(1) FROM "admins" WHERE user_id = $1`

			userRows := sqlmock.NewRows([]string{"id", "chat_id", "subs_count"})
			countRows := sqlmock.NewRows([]string{"count"})
			var userID uint
			if tt.want {
				userID = 1
			} else {
				userID = 2
			}
			userRows.AddRow(userID, tt.args.ChatID, 1)
			countRows.AddRow(1)

			dbmock.mock.ExpectQuery(regexp.QuoteMeta(sqlCountUser)).
				WithArgs(tt.args.ChatID).
				WillReturnRows(countRows)

			dbmock.mock.ExpectQuery(regexp.QuoteMeta(sqlSelectUser)).
				WithArgs(tt.args.ChatID).
				WillReturnRows(userRows)

			countAdminRows := sqlmock.NewRows([]string{"count"})
			if tt.want {
				countAdminRows.AddRow(1)
			} else {
				countAdminRows.AddRow(0)
			}

			dbmock.mock.ExpectQuery(regexp.QuoteMeta(sqlCountAdmin)).
				WithArgs(userID).
				WillReturnRows(countAdminRows)

			tstp := userStorage.IsChatAdmin(tt.args.ChatID)

			assert.Equal(tt.want, tstp)

			dbmock.AfterEach(t)
		})
	}
}
