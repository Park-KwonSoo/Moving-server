package nosql_model

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	noSql "github.com/Park-Kwonsoo/moving-server/pkg/database/nosql"
)

type Music struct {
	noSql.BaseType `bson:",inline"`
	TrackNumber    uint   `bson:"track_number"`
	Artist         string `bson:"artist"`
	Album          string `bson:"album"`
	Genre          string `bson:"genre"`
	Title          string `bson:"title"`
	MusicUrl       string `bson:"music_url"`
	AlbumCoverUrl  string `bson:"album_cover_url"`
	IsTitle        bool   `bson:"is_title"`
}

/**
*	Create New Music
 */
func CreateNewMusic(music *Music) error {

	music.ID = primitive.NewObjectID()
	music.CreatedAt = time.Now()
	music.UpdatedAt = time.Now()

	_, err := noSql.GetCollection("music").InsertOne(context.TODO(), music)

	return err
}

/**
*	id로 음악 찾기
 */
func FindOneMusicById(id primitive.ObjectID) (*Music, error) {

	music := &Music{}
	if err := noSql.GetCollection("music").FindOne(context.TODO(), bson.M{"_id": id}).Decode(music); err != nil {
		return nil, err
	}

	return music, nil
}

/**
*	toDo : 키워드로 음악 찾기
 */
func FindAllMusicByKeyword(keyword string) ([]*Music, error) {
	musicList := make([]*Music, 0)
	// query := qb.Select("id, create_at, tracknumber, title, artist, album, genre, album_img, music_url, is_title").From("music").ToString()

	return musicList, nil
}

func init() {
}
