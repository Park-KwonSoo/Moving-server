package models

type UserPlaylist struct {
	baseType
	User     User
	Playlist []Playlist
}
