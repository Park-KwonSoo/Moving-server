package models

type UserPlaylist struct {
	baseType
	User     User `db:"user_id int references user(id)"`
	Playlist []Playlist
}
