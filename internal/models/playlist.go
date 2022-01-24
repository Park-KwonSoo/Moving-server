package models

import (
	"log"
	"strings"

	getTag "github.com/Park-Kwonsoo/moving-server/pkg/get-struct-info"
	qb "github.com/Park-Kwonsoo/moving-server/pkg/query-builder"
)

type Playlist struct {
	baseType
	Member       Member  `db:"member_mem_id varchar(255) references member(mem_id) on delete cascade" mapping:"many2one member"`
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

	if _, err := psql.Exec(query); err != nil {
		return err
	}

	if err := createUpdateTrigger("playlist"); err != nil {
		return err
	}
	//테이블을 매핑함
	err := tableMapping(&Playlist{})
	return err
}

/**
*	Create New Playlist
 */
func CreateNewPlaylist(playlist *Playlist) error {

	query := qb.Insert("playlist", "playlist_name, member_mem_id").Value(
		playlist.PlaylistName, playlist.Member.MemId.String,
	).ToString()

	err := psql.QueryRow(query).Scan(&playlist.ID)

	return err
}

/**
*	Add Music In Playlist
 */
func AddMusicInPlaylist(playlist *Playlist, musicIdList ...uint64) error {

	tx, err := psql.db.Begin()
	if err != nil {
		log.Panic(err)
	}
	defer tx.Rollback()

	for _, musicId := range musicIdList {
		query := qb.Insert("music_playlist", "music_id, playlist_id").Value(
			musicId, playlist.ID,
		).ToString()

		_, err := tx.Exec(query)
		if err != nil {
			log.Panic(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Panic(err)
	}

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

		query = qb.Select("music_id").From("music_playlist").Join("music_id", "music", "id").Where("playlist_id", playlist.ID).ToString()
		joinRows, err := psql.Query(query)
		if err != nil {
			return playlistList, err
		}
		defer joinRows.Close()

		for joinRows.Next() {
			music := &Music{}
			joinRows.Scan(
				&music.ID,
			)
			playlist.Music = append(playlist.Music, *music)
		}

		playlistList = append(playlistList, playlist)
	}

	return playlistList, nil
}

/**
*	playlist id로 playlist 하나 찾음
 */
func FindOnePlaylistById(id uint) (*Playlist, error) {

	playlist := &Playlist{}

	query := qb.Select("id, created_at, updated_at, playlist_name, member_mem_id").From("playlist").Where("id", id).ToString()
	if err := psql.QueryRow(query).Scan(
		&playlist.ID, &playlist.CreatedAt, &playlist.UpdatedAt, &playlist.PlaylistName, &playlist.Member.MemId.String,
	); err != nil {
		return nil, err
	}

	query = qb.Select("music_id, track_number, title, artist, album, genre, album_img, music_url, is_title").From("music_playlist").Join("music_id", "music", "id").Where("playlist_id", playlist.ID).ToString()
	rows, err := psql.Query(query)
	if err != nil {
		return playlist, err
	}
	defer rows.Close()

	for rows.Next() {
		music := &Music{}
		rows.Scan(
			&music.ID, &music.TrackNumber, &music.Title, &music.Artist, &music.Album, &music.Genre, &music.AlbumImg, &music.MusicUrl, &music.IsTitle,
		)
		playlist.Music = append(playlist.Music, *music)
	}

	return playlist, nil
}

/**
*	update One Playlist
 */
func UpdateOnePlaylist(playlist *Playlist) error {
	query := qb.Update("playlist").Set("playlist_name", []string{
		playlist.PlaylistName,
	}).Where("id", playlist.ID).ToString()

	_, err := psql.Exec(query)

	return err
}
