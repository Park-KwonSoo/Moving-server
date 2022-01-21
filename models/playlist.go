package models

import (
	"strings"

	getTag "github.com/Park-Kwonsoo/moving-server/pkg/get-struct-info"
	qb "github.com/Park-Kwonsoo/moving-server/pkg/query-builder"
)

type Playlist struct {
	baseType
	MemPlaylistID uint    `db:"mem_playlist_id int references member_playlist(id)" mapping:"many2one member_playlist"`
	PlaylistName  string  `db:"playlist_name varchar(255) not null default '이름 없음'"`
	Music         []Music `mapping:"many2many music"`
}

//playlist migrate
func playlistMigrate() error {

	column := make([]string, 0)
	column = append(column, strings.Join(getCreatedTableColumn(), ", "))
	column = append(column, getTag.GetStructInfoByTag("db", &Playlist{})...)

	query := qb.CreateTable("playlist").TableComlumn(
		column...,
	).ToString()

	_, err := psql.db.Exec(query)

	return err
}
