package models

import (
	"strings"

	getTag "github.com/Park-Kwonsoo/moving-server/pkg/get-struct-info"
	qb "github.com/Park-Kwonsoo/moving-server/pkg/query-builder"
)

type Playlist struct {
	baseType
	Member       Member  `db:"member_mem_id varchar(255) references member(mem_id)" mapping:"many2one member"`
	PlaylistName string  `db:"playlist_name varchar(255) not null default '이름 없음'"`
	Music        []Music `mapping:"many2many music"`
}

//playlist migrate
func playlistMigrate() error {

	column := make([]string, 0)
	column = append(column, strings.Join(getCreatedTableColumn(), ", "))
	column = append(column, getTag.GetStructInfoByTag("db", &Playlist{})...)

	query := qb.CreateTable("playlist").TableComlumn(
		column...,
	).ToString()

	_, err := psql.Exec(query)
	if err != nil {
		return err
	}

	err = createUpdateTrigger("playlist")
	if err != nil {
		return err
	}
	//테이블을 매핑함
	err = tableMapping(&Playlist{})
	return err
}

/**
*	find all by member id
 */
func FindAllPlaylistByMemberMemId(memId string) ([]*Playlist, error) {

	playlistList := make([]*Playlist, 0)

	query := qb.Select("id, created_at, updated_at, playlist_name").From("playlist").Where("member_mem_id", memId).ToString()
	rows, err := psql.Query(query)
	if err != nil {
		return playlistList, err
	}
	defer rows.Close()

	for rows.Next() {
		playlist := &Playlist{}
		rows.Scan(&playlist.ID, &playlist.CreatedAt, &playlist.UpdatedAt, &playlist.PlaylistName)
		playlistList = append(playlistList, playlist)
	}

	return playlistList, nil
}
