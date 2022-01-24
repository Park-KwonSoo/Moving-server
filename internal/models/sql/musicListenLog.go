package sql_model

import (
	"strings"

	"github.com/sirupsen/logrus"

	sqlDB "github.com/Park-Kwonsoo/moving-server/pkg/database/sql"
	getTag "github.com/Park-Kwonsoo/moving-server/pkg/get-struct-info"
	qb "github.com/Park-Kwonsoo/moving-server/pkg/query-builder"
)

type MusicListenLog struct {
	sqlDB.BaseType
	MemId   string `db:"member_mem_id varchar(255) references member(mem_id) on delete set null" mapping:"many2one member"`
	MusicId string `db:"music_id varchar(255)"`
}

//music_listen_log table migrate
func musicListenLogMigrate() error {

	column := make([]string, 0)
	column = append(column, strings.Join(sqlDB.GetCreatedTableColumn(), ", "))
	column = append(column, getTag.GetStructInfoByTag("db", &MusicListenLog{})...)

	query := qb.CreateTable("music_listen_log").TableComlumn(
		column...,
	).ToString()

	if _, err := sqlDB.SQL.Exec(query); err != nil {
		return err
	}

	if err := sqlDB.CreateUpdateTrigger("music_listen_log"); err != nil {
		return err
	}

	err := sqlDB.TableMapping(&MusicListenLog{})
	return err
}

func init() {
	if err := musicListenLogMigrate(); err != nil {
		logrus.Error(err)
	}
}
