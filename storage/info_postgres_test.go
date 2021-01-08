package storage

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type StorageMock struct {
	storage *InfoPostgres
	mock    sqlmock.Sqlmock
}

func (mock *StorageMock) BeforeEach(t *testing.T) {
	var db *sql.DB
	var err error

	db, mocked, err := sqlmock.New()
	mock.mock = mocked
	assert.Nil(t, err)

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.Nil(t, err)

	mock.storage = NewInfoPostgres(gdb)
}

func (mock *StorageMock) AfterEach(t *testing.T) {
	err := mock.mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestNewInfoPostgres(t *testing.T) {
	assert := assert.New(t)

	db, _, err := sqlmock.New()
	assert.Nil(err)
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.Nil(err)

	infop := NewInfoPostgres(gdb)

	assert.Equal(gdb, infop.db, "Equal db instances")
}

func TestInfoPostgres_GetLastTimestamp(t *testing.T) {
	assert := assert.New(t)
	dbmock := StorageMock{}

	type fields struct {
		timestamp uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   uint64
	}{
		{"Get Timestamp", fields{123}, 123},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbmock.BeforeEach(t)

			infoStorage := dbmock.storage
			rows := sqlmock.
				NewRows([]string{"id", "last_post"}).
				AddRow(1, tt.fields.timestamp)
			const sqlSelectOne = `SELECT * FROM "infos" ORDER BY "infos"."id" LIMIT 1`
			dbmock.mock.ExpectQuery(regexp.QuoteMeta(sqlSelectOne)).WillReturnRows(rows)

			tstp := infoStorage.GetLastTimestamp()
			assert.Equal(tt.want, tstp)

			dbmock.AfterEach(t)
		})
	}
}

func TestInfoPostgres_SetLastTimestamp(t *testing.T) {
	dbmock := StorageMock{}
	type fields struct {
		timestamp uint64
	}
	type args struct {
		timestamp uint64
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"Set Timestamp", fields{123}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbmock.BeforeEach(t)

			infoStorage := dbmock.storage

			const sqlUpdate = `UPDATE "infos" SET "last_post"=$1 WHERE "id" = $2`
			const sqlSelect = `SELECT * FROM "infos" ORDER BY "infos"."id" LIMIT 1`

			dbmock.mock.ExpectQuery(regexp.QuoteMeta(sqlSelect)).
				WillReturnRows(sqlmock.
					NewRows([]string{"id", "last_post"}).
					AddRow(1, 1))

			dbmock.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "infos" SET "last_post"=$1 WHERE "id" = $2`)).
				WithArgs(tt.fields.timestamp, 1).
				WillReturnResult(sqlmock.NewResult(1, 1))

			infoStorage.SetLastTimestamp(tt.fields.timestamp)

			dbmock.AfterEach(t)
		})
	}
}
