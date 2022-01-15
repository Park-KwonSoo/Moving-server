package models

type Playlist struct {
	basicType
	UserPlaylistID uint
	PlaylistName   string
	Music          []Music `gorm:"many2many:playlist_musics;"`
}
