package models

import (
	"strings"

	getTag "github.com/Park-Kwonsoo/moving-server/pkg/get-struct-info"
	qb "github.com/Park-Kwonsoo/moving-server/pkg/query-builder"
)

type MemberPlaylist struct {
	baseType
	Member   Member `db:"member_mem_id varchar(255) references member(mem_id)" mapping:"one2one member"`
	Playlist []Playlist
}

//member play list migrate
func memberPlaylistMigrate() error {

	column := make([]string, 0)
	column = append(column, strings.Join(getCreatedTableColumn(), ", "))
	column = append(column, getTag.GetStructInfoByTag("db", &MemberPlaylist{})...)

	query := qb.CreateTable("member_playlist").TableComlumn(
		column...,
	).ToString()

	_, err := psql.db.Exec(query)

	return err

}
