package models

type UserPlaylist struct {
	basicType
	User     User `gorm:"constraint:OnDelete:CASCADE;"`
	Playlist []Playlist
}
