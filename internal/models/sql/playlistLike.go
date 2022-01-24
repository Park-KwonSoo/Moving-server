package sql_model

import (
	"strings"

	"github.com/sirupsen/logrus"

	sqlDB "github.com/Park-Kwonsoo/moving-server/pkg/database/sql"
	getTag "github.com/Park-Kwonsoo/moving-server/pkg/get-struct-info"
	qb "github.com/Park-Kwonsoo/moving-server/pkg/query-builder"
)

type PlaylistLike struct {
	sqlDB.BaseType
	MemId      string `db:"member_mem_id varchar(255) references member(mem_id) on delete cascade" mapping:"many2one member"`
	PlaylistId string `db:"playlist_id varchar(255)"`
}

//playlist_like table migrate
func playlistLikeMigrate() error {

	column := make([]string, 0)
	column = append(column, strings.Join(sqlDB.GetCreatedTableColumn(), ", "))
	column = append(column, getTag.GetStructInfoByTag("db", &PlaylistLike{})...)

	query := qb.CreateTable("playlist_like").TableComlumn(
		column...,
	).ToString()

	if _, err := sqlDB.SQL.Exec(query); err != nil {
		return err
	}

	if err := sqlDB.CreateUpdateTrigger("playlist_like"); err != nil {
		return err
	}

	err := sqlDB.TableMapping(&PlaylistLike{})
	return err
}

func init() {
	if err := playlistLikeMigrate(); err != nil {
		logrus.Error(err)
	}
}
