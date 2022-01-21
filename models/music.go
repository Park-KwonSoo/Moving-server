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

	_, err := psql.db.Exec(query)

	return err
}

func FindMusicById(id uint) *Music {
	music := Music{}
	// db.db(&Music{}).First(music, "ID = ?", id)

	return &music
}
