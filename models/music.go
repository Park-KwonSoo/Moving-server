package models

import (
	"reflect"
	"strings"

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
	m := Music{}

	trackNumber, _ := reflect.TypeOf(m).FieldByName("TrackNumber")
	title, _ := reflect.TypeOf(m).FieldByName("Title")
	artist, _ := reflect.TypeOf(m).FieldByName("Artist")
	album, _ := reflect.TypeOf(m).FieldByName("Album")
	albumImg, _ := reflect.TypeOf(m).FieldByName("AlbumImg")
	genre, _ := reflect.TypeOf(m).FieldByName("Genre")
	musicUrl, _ := reflect.TypeOf(m).FieldByName("MusicUrl")
	isTitle, _ := reflect.TypeOf(m).FieldByName("IsTitle")

	query := qb.CreateTable("music").TableComlumn([]string{
		strings.Join(getCreatedTableColumn(), ", "),
		trackNumber.Tag.Get("db"),
		title.Tag.Get("db"),
		artist.Tag.Get("db"),
		album.Tag.Get("db"),
		albumImg.Tag.Get("db"),
		genre.Tag.Get("db"),
		musicUrl.Tag.Get("db"),
		isTitle.Tag.Get("db"),
	}).ToString()

	_, err := psql.db.Exec(query)

	return err
}

func FindMusicById(id uint) *Music {
	music := Music{}
	// db.db(&Music{}).First(music, "ID = ?", id)

	return &music
}
