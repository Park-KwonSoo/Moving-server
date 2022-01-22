package models

import (
	"strings"

	getTag "github.com/Park-Kwonsoo/moving-server/pkg/get-struct-info"
	qb "github.com/Park-Kwonsoo/moving-server/pkg/query-builder"
)

type Music struct {
	baseType
	TrackNumber uint   `db:"track_number integer not null"`
	Title       string `db:"title varchar(200) not null"`
	Artist      string `db:"artist varchar(200) not null"`
	Album       string `db:"album varchar(200) not null"`
	AlbumImg    string `db:"album_img varchar(2000) not null"`
	Genre       string `db:"genre varchar(200) not null"`
	MusicUrl    string `db:"music_url varchar(2000) unique not null"`
	IsTitle     bool   `db:"is_title boolean not null"`
}

//music migrate
func musicMigrate() error {

	column := make([]string, 0)
	column = append(column, strings.Join(getCreatedTableColumn(), ", "))
	column = append(column, getTag.GetStructInfoByTag("db", &Music{})...)

	query := qb.CreateTable("music").TableComlumn(
		column...,
	).ToString()

	_, err := psql.Exec(query)
	if err != nil {
		return err
	}

	err = createUpdateTrigger("music")
	if err != nil {
		return err
	}

	err = tableMapping(&Music{})
	return err
}

/**
*	id로 음악 찾기
 */
func FindOneMusicById(id uint) (*Music, error) {
	music := &Music{}
	query := qb.Select("id, created_at, track_number, title, artist, album, album_img, genre, music_url, is_title").From("music").Where("id", id).ToString()

	if err := psql.QueryRow(query).Scan(
		&music.ID, &music.CreatedAt, &music.TrackNumber, &music.Title, &music.Artist, &music.Album, &music.AlbumImg, &music.Genre, &music.MusicUrl, &music.IsTitle,
	); err != nil {
		return nil, err
	}

	return music, nil
}
