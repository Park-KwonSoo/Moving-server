package models

type Playlist struct {
	baseType
	UserPlaylistID uint
	PlaylistName   string
	Music          []Music `gorm:"many2many:playlist_musics;"`
}
