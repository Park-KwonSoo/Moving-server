package models

type Playlist struct {
	baseType
	UserPlaylistID uint   `db:"user_playlist_id int references user_playlist(id)"`
	PlaylistName   string `db:"playlist_name varchar(255) not null default '이름 없음'"`
	Music          []Music
}
