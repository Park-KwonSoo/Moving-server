package models

type UserPlaylist struct {
	baseType
	User     User `gorm:"constraint:OnDelete:CASCADE;"`
	Playlist []Playlist
}
