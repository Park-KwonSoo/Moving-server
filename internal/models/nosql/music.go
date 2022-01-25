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
func CreateNewMusic(music *Music, album *Album) error {

	music.ID = primitive.NewObjectID()
	music.CreatedAt = time.Now()
	music.UpdatedAt = time.Now()

	//음악의 앨범과 커버는 자동으로 설정됨.
	music.Album = album.Album
	music.AlbumCoverUrl = album.AlbumCoverUrl
	//음악 추가시 자동으로 해당 앨범에 등록됨, 즉 앨범이 있어야 음악을 생성할 수 있음

	errChannelAddAlbum := make(chan error)
	errChannelCreateMusic := make(chan error)

	go func() {
		errChannelAddAlbum <- AddMusicInAlbum(album, *music)
	}()

	go func() {
		_, err := noSql.GetCollection("music").InsertOne(context.TODO(), music)
		errChannelCreateMusic <- err
	}()

	if err := <-errChannelAddAlbum; err != nil {
		return err
	}

	if err := <-errChannelCreateMusic; err != nil {
		return err
	}

	return nil

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
*	키워드 : 제목, 가수, 앨범, 장르
 */
func FindAllMusicByKeyword(keyword string) ([]*Music, error) {

	musicList := make([]*Music, 0)
	query := bson.M{
		"$or": append(make([]bson.M, 0),
			bson.M{
				"title": bson.M{
					"$regex": keyword,
				}},
			bson.M{
				"artist": bson.M{
					"$regex": keyword,
				}},
			bson.M{
				"album": bson.M{
					"$regex": keyword,
				}},
			bson.M{
				"genre": bson.M{
					"$regex": keyword,
				}},
		),
	}

	res, err := noSql.GetCollection("music").Find(context.TODO(), query)
	if err != nil {
		return musicList, err
	}
	if err := res.All(context.TODO(), &musicList); err != nil {
		return musicList, err
	}

	return musicList, nil
}

func init() {
}
