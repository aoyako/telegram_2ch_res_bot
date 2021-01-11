package storage

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type SubscribtionMock struct {
	storage *SubscriptionPostgres
	mock    sqlmock.Sqlmock
}

func (mock *SubscribtionMock) BeforeEach(t *testing.T) {
	var db *sql.DB
	var err error

	db, mocked, err := sqlmock.New()
	mock.mock = mocked
	assert.Nil(t, err)

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.Nil(t, err)

	mock.storage = NewSubscriptionPostgres(gdb)
}

func (mock *SubscribtionMock) AfterEach(t *testing.T) {
	err := mock.mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestNewSubscriptionPostgres(t *testing.T) {
	assert := assert.New(t)

	db, _, err := sqlmock.New()
	assert.Nil(err)
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.Nil(err)

	infop := NewSubscriptionPostgres(gdb)

	assert.Equal(gdb, infop.db, "Equal db instances")
}

func TestSubscriptionPostgres_Add(t *testing.T) {
	assert := assert.New(t)
	dbmock := SubscribtionMock{}

	type args struct {
		user        *logic.User
		publication *logic.Publication
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Insert publication with id",
			args: args{
				user: &logic.User{
					ID:     1,
					ChatID: 1,
				},
				publication: &logic.Publication{
					Board:     "abc",
					IsDefault: false,
					Type:      "",
					ID:        7,
				},
			},
			err: nil,
		},
		{
			name: "Insert publication with no id",
			args: args{
				user: &logic.User{
					ID:     1,
					ChatID: 1,
				},
				publication: &logic.Publication{
					Board:     "abc",
					IsDefault: false,
					Type:      "",
				},
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbmock.BeforeEach(t)

			subsStorage := dbmock.storage
			const sqlInsertPublicationWithID = `INSERT INTO "publications" ("board","tags","is_default","type","alias","id") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`
			const sqlInsertPublication = `INSERT INTO "publications" ("board","tags","is_default","type","alias") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`
			const sqlInsertUser = `INSERT INTO "users" ("chat_id","subs_count","id") VALUES ($1,$2,$3) ON CONFLICT DO NOTHING RETURNING "id"`
			const sqlInsertUserSubscribtion = `INSERT INTO "user_subscribtion" ("publication_id","user_id") VALUES ($1,$2) ON CONFLICT DO NOTHING`

			userInst := tt.args.user
			pubInst := tt.args.publication

			if pubInst.ID == 0 {
				dbmock.mock.ExpectQuery(regexp.QuoteMeta(sqlInsertPublication)).
					WithArgs(pubInst.Board, pubInst.Tags, pubInst.IsDefault, pubInst.Type, pubInst.Alias).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(1))
			} else {
				dbmock.mock.ExpectQuery(regexp.QuoteMeta(sqlInsertPublicationWithID)).
					WithArgs(pubInst.Board, pubInst.Tags, pubInst.IsDefault, pubInst.Type, pubInst.Alias, pubInst.ID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(pubInst.ID))
			}

			dbmock.mock.ExpectQuery(regexp.QuoteMeta(sqlInsertUser)).
				WithArgs(userInst.ChatID, userInst.SubsCount, userInst.ID).
				WillReturnRows(
					sqlmock.NewRows([]string{"id"}).
						AddRow(userInst.ID))

			if pubInst.ID == 0 {
				dbmock.mock.ExpectExec(regexp.QuoteMeta(sqlInsertUserSubscribtion)).
					WithArgs(1, userInst.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				dbmock.mock.ExpectExec(regexp.QuoteMeta(sqlInsertUserSubscribtion)).
					WithArgs(pubInst.ID, userInst.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			}

			tstp := subsStorage.Add(tt.args.user, tt.args.publication)
			assert.Equal(tt.err, tstp)

			dbmock.AfterEach(t)
		})
	}
}

func TestSubscriptionPostgres_AddDefault(t *testing.T) {
	assert := assert.New(t)
	dbmock := SubscribtionMock{}

	type args struct {
		publication *logic.Publication
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Insert publication with id",
			args: args{
				publication: &logic.Publication{
					Board: "abc",
					Type:  "",
				},
			},
			err: nil,
		},
		{
			name: "Insert publication with no id",
			args: args{
				publication: &logic.Publication{
					Board: "abc",
					Type:  "",
				},
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbmock.BeforeEach(t)

			subsStorage := dbmock.storage
			const sqlInsertPublicationWithID = `INSERT INTO "publications" ("board","tags","is_default","type","alias","id") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`
			const sqlInsertPublication = `INSERT INTO "publications" ("board","tags","is_default","type","alias") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`

			pubInst := tt.args.publication

			if pubInst.ID == 0 {
				dbmock.mock.ExpectQuery(regexp.QuoteMeta(sqlInsertPublication)).
					WithArgs(pubInst.Board, pubInst.Tags, true, pubInst.Type, pubInst.Alias).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(1))
			} else {
				dbmock.mock.ExpectQuery(regexp.QuoteMeta(sqlInsertPublicationWithID)).
					WithArgs(pubInst.Board, pubInst.Tags, true, pubInst.Type, pubInst.Alias, pubInst.ID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(pubInst.ID))
			}

			tstp := subsStorage.AddDefault(tt.args.publication)
			assert.Equal(tt.err, tstp)
			assert.True(tt.args.publication.IsDefault)

			dbmock.AfterEach(t)
		})
	}
}

func TestSubscriptionPostgres_Remove(t *testing.T) {
	assert := assert.New(t)
	dbmock := SubscribtionMock{}

	type args struct {
		publication *logic.Publication
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Delete publication",
			args: args{
				publication: &logic.Publication{
					Board: "abc",
					Type:  "",
					ID:    2,
				},
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbmock.BeforeEach(t)

			subsStorage := dbmock.storage
			const sqlDelete = `DELETE FROM "publications" WHERE "publications"."id" = $1`

			dbmock.mock.ExpectExec(regexp.QuoteMeta(sqlDelete)).
				WithArgs(tt.args.publication.ID).
				WillReturnResult(sqlmock.NewResult(int64(tt.args.publication.ID), 1))

			tstp := subsStorage.Remove(tt.args.publication)
			assert.Equal(tt.err, tstp)

			dbmock.AfterEach(t)
		})
	}
}

func TestSubscriptionPostgres_Connect(t *testing.T) {
	assert := assert.New(t)
	dbmock := SubscribtionMock{}

	type args struct {
		publication *logic.Publication
		user        *logic.User
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Connect publication",
			args: args{
				user: &logic.User{
					ID:     7,
					ChatID: 1,
				},
				publication: &logic.Publication{
					Board:     "abc",
					IsDefault: false,
					Type:      "",
					ID:        8,
				},
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbmock.BeforeEach(t)

			subsStorage := dbmock.storage
			userInst := tt.args.user

			const sqlInsertUser = `INSERT INTO "users" ("chat_id","subs_count","id") VALUES ($1,$2,$3) ON CONFLICT DO NOTHING RETURNING "id"`
			const sqlInsertUserSubscribtion = `INSERT INTO "user_subscribtion" ("publication_id","user_id") VALUES ($1,$2) ON CONFLICT DO NOTHING`

			dbmock.mock.ExpectQuery(regexp.QuoteMeta(sqlInsertUser)).
				WithArgs(userInst.ChatID, userInst.SubsCount, userInst.ID).
				WillReturnRows(
					sqlmock.NewRows([]string{"id"}).
						AddRow(userInst.ID))

			dbmock.mock.ExpectExec(regexp.QuoteMeta(sqlInsertUserSubscribtion)).
				WithArgs(tt.args.publication.ID, tt.args.user.ID).
				WillReturnResult(sqlmock.NewResult(1, 1))

			tstp := subsStorage.Connect(tt.args.user, tt.args.publication)
			assert.Equal(tt.err, tstp)

			dbmock.AfterEach(t)
		})
	}
}

func TestSubscriptionPostgres_Update(t *testing.T) {
	assert := assert.New(t)
	dbmock := SubscribtionMock{}

	type args struct {
		publication *logic.Publication
		user        *logic.User
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Update Publication",
			args: args{
				publication: &logic.Publication{
					Board: "abc",
					Type:  "",
					ID:    20,
				},
				user: nil,
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbmock.BeforeEach(t)

			subsStorage := dbmock.storage
			const sqlInsertPublication = `UPDATE "publications" SET "board"=$1,"tags"=$2,"is_default"=$3,"type"=$4,"alias"=$5 WHERE "id" = $6`
			pubInst := tt.args.publication

			dbmock.mock.ExpectExec(regexp.QuoteMeta(sqlInsertPublication)).
				WithArgs(pubInst.Board, pubInst.Tags, pubInst.IsDefault, pubInst.Type, pubInst.Alias, pubInst.ID).
				WillReturnResult(sqlmock.NewResult(1, 1))

			tstp := subsStorage.Update(tt.args.user, tt.args.publication)
			assert.Equal(tt.err, tstp)

			dbmock.AfterEach(t)
		})
	}
}

func TestSubscriptionPostgres_GetSubsByUser(t *testing.T) {
	assert := assert.New(t)
	dbmock := SubscribtionMock{}

	type args struct {
		user *logic.User
	}
	tests := []struct {
		name string
		args args
		want []logic.Publication
	}{
		{
			name: "Get subs",
			args: args{
				user: &logic.User{
					ID:     7,
					ChatID: 1,
				},
			},
			want: []logic.Publication{
				logic.Publication{
					ID: 2,
				},
				logic.Publication{
					ID: 3,
				},
			},
		},
		{
			name: "No subs",
			args: args{
				user: &logic.User{
					ID:     7,
					ChatID: 1,
				},
			},
			want: []logic.Publication{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbmock.BeforeEach(t)

			subsStorage := dbmock.storage

			userInst := tt.args.user
			const sqlInsertUserSubscribtion = `SELECT 
				"publications"."id","publications"."board","publications"."tags","publications"."is_default","publications"."type","publications"."alias"
				FROM "publications" JOIN "user_subscribtion" ON "user_subscribtion"."publication_id" = "publications"."id" AND "user_subscribtion"."user_id" = $1`

			rows := sqlmock.NewRows([]string{"id", "board", "tags", "is_default", "type", "alias"})
			for i := range tt.want {
				pub := tt.want[i]
				rows.AddRow(pub.ID, pub.Board, pub.Tags, pub.IsDefault, pub.Type, pub.Alias)
			}

			dbmock.mock.ExpectQuery(regexp.QuoteMeta(sqlInsertUserSubscribtion)).
				WithArgs(userInst.ID).
				WillReturnRows(rows)

			tstp, err := subsStorage.GetSubsByUser(tt.args.user)
			assert.Equal(tt.want, tstp)
			assert.Nil(err)

			dbmock.AfterEach(t)
		})
	}
}

func TestSubscriptionPostgres_GetAllSubs(t *testing.T) {
	assert := assert.New(t)
	dbmock := SubscribtionMock{}

	tests := []struct {
		name string
		want []logic.Publication
	}{
		{
			name: "Get all subs",
			want: []logic.Publication{
				logic.Publication{
					ID: 2,
				},
				logic.Publication{
					ID: 3,
				},
			},
		},
		{
			name: "No subs",
			want: []logic.Publication{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbmock.BeforeEach(t)

			subsStorage := dbmock.storage

			const sqlInsertUserSubscribtion = `SELECT * FROM "publications"`

			rows := sqlmock.NewRows([]string{"id", "board", "tags", "is_default", "type", "alias"})
			for i := range tt.want {
				pub := tt.want[i]
				rows.AddRow(pub.ID, pub.Board, pub.Tags, pub.IsDefault, pub.Type, pub.Alias)
			}

			dbmock.mock.ExpectQuery(regexp.QuoteMeta(sqlInsertUserSubscribtion)).
				WillReturnRows(rows)

			tstp := subsStorage.GetAllSubs()
			assert.Equal(tt.want, tstp)

			dbmock.AfterEach(t)
		})
	}
}

func TestSubscriptionPostgres_GetAllDefaultSubs(t *testing.T) {
	assert := assert.New(t)
	dbmock := SubscribtionMock{}

	tests := []struct {
		name string
		want []logic.Publication
		err  error
	}{
		{
			name: "Get all subs",
			want: []logic.Publication{
				logic.Publication{
					ID:        2,
					IsDefault: true,
				},
				logic.Publication{
					ID:        3,
					IsDefault: true,
				},
			},
			err: nil,
		},
		{
			name: "No subs",
			want: []logic.Publication{},
			err:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbmock.BeforeEach(t)

			subsStorage := dbmock.storage

			const sqlInsertUserSubscribtion = `SELECT * FROM "publications" WHERE is_default = $1`

			rows := sqlmock.NewRows([]string{"id", "board", "tags", "is_default", "type", "alias"})
			for i := range tt.want {
				pub := tt.want[i]
				rows.AddRow(pub.ID, pub.Board, pub.Tags, true, pub.Type, pub.Alias)
			}

			dbmock.mock.ExpectQuery(regexp.QuoteMeta(sqlInsertUserSubscribtion)).
				WithArgs(true).
				WillReturnRows(rows)

			tstp := subsStorage.GetAllDefaultSubs()
			assert.Equal(tt.want, tstp)

			dbmock.AfterEach(t)
		})
	}
}
