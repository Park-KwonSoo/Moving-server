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

/**
*	Playlist Like DB Create
 */
func CreateNewPlaylistLike(playlistLike *PlaylistLike) error {

	query := qb.Insert("playlist_like", "member_mem_id, playlist_id").Value(
		playlistLike.MemId,
		playlistLike.PlaylistId,
	).ToString()

	err := sqlDB.SQL.QueryRow(query).Scan(&playlistLike.ID)

	return err
}

/**
*	Check playlistLike Table Has Db By MemId And PlaylistId
 */
func HavePlaylistLikeByMemIdAndPlaylistId(memId string, playlistId string) (bool, error) {

	var count int

	query := qb.Select("count(*)").
		From("playlist_like").
		Where("member_mem_id", memId).
		And("playlist_id", playlistId).
		ToString()

	err := sqlDB.SQL.QueryRow(query).Scan(&count)
	if count > 0 {
		return false, err
	}

	return true, err
}

/**
*	Count Db in playlistLike Table By MemId And PlaylistId
 */
func CountPlaylistLikePlaylistId(playlistId string) (int, error) {

	var count int

	query := qb.Select("count(*)").
		From("playlist_like").
		Where("playlist_id", playlistId).
		ToString()

	err := sqlDB.SQL.QueryRow(query).Scan(&count)

	return count, err
}

func init() {
	if err := playlistLikeMigrate(); err != nil {
		logrus.Error(err)
	}
}
