package models

type Music struct {
	basicType
	Title       string
	Artist      string
	Album       string
	AlbumImg    string
	Genre       string
	TrackNumber uint
	MusicUrl    string
	IsTitle     bool
}

func FindMusicById(id uint) *Music {
	music := Music{}
	// db.Model(&Music{}).First(music, "ID = ?", id)

	return &music
}
