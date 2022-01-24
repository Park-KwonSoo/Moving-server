package nosql_model

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	noSql "github.com/Park-Kwonsoo/moving-server/pkg/database/nosql"
)

type Playlist struct {
	noSql.BaseType `bson:",inline"`
	MemId          string  `bson:"memId"`
	PlaylistName   string  `bson:"playlist_name"`
	Music          []Music `bson:"music"`
}

/**
*	Create New Playlist
 */
func CreateNewPlaylist(playlist *Playlist) error {

	playlist.ID = primitive.NewObjectID()
	playlist.CreatedAt = time.Now()
	playlist.UpdatedAt = time.Now()

	_, err := noSql.GetCollection("playlist").InsertOne(context.TODO(), playlist)

	return err
}

/**
*	Add Music In Playlist
 */
func AddMusicInPlaylist(playlist *Playlist, musicIdList ...string) error {

	//음악 찾기 쿼리
	musicObjectIdList := make([]primitive.ObjectID, 0)
	for _, musicId := range musicIdList {
		objId, _ := primitive.ObjectIDFromHex(musicId)
		musicObjectIdList = append(musicObjectIdList, objId)
	}

	filter := make([]bson.M, 0)
	for i := 0; i < len(musicObjectIdList); i++ {
		filter = append(filter, bson.M{
			"_id": musicObjectIdList[i],
		})
	}

	findQuery := bson.M{
		"$or": filter,
	}

	//음악들을 가져옴
	musicList := make([]Music, 0)
	res, err := noSql.GetCollection("music").Find(context.TODO(), findQuery)
	if err != nil {
		return err
	}

	if err := res.All(context.TODO(), &musicList); err != nil {
		return err
	}

	//슬라이스에서 추가 : push
	update := bson.M{
		"$push": bson.M{
			"music": bson.M{
				"$each": append(make([]Music, 0), musicList...),
			},
		},
	}

	if _, err := noSql.GetCollection("playlist").UpdateOne(context.TODO(), bson.M{
		"_id": playlist.ID,
	}, update); err != nil {
		return err
	}

	return nil
}

/**
*	remove music in playlist
 */
func RemoveMusicInPlaylist(playlist *Playlist, musicIdList ...string) error {
	//음악 찾기 쿼리
	musicObjectIdList := make([]primitive.ObjectID, 0)
	for _, musicId := range musicIdList {
		objId, _ := primitive.ObjectIDFromHex(musicId)
		musicObjectIdList = append(musicObjectIdList, objId)
	}

	filter := make([]bson.M, 0)
	for i := 0; i < len(musicObjectIdList); i++ {
		filter = append(filter, bson.M{
			"_id": musicObjectIdList[i],
		})
	}

	findQuery := bson.M{
		"$or": filter,
	}

	//음악들을 가져옴
	musicList := make([]Music, 0)
	res, err := noSql.GetCollection("music").Find(context.TODO(), findQuery)
	if err != nil {
		return err
	}

	if err := res.All(context.TODO(), &musicList); err != nil {
		return err
	}

	//슬라이스에서 제거 : pull
	update := bson.M{
		"$pull": bson.M{
			"music": bson.M{
				"$in": append(make([]Music, 0), musicList...),
			},
		},
	}

	if _, err := noSql.GetCollection("playlist").UpdateOne(context.TODO(), bson.M{
		"_id": playlist.ID,
	}, update); err != nil {
		return err
	}

	return nil
}

/**
*	find all by member id
 */
func FindAllPlaylistByMemberMemId(memId string) ([]*Playlist, error) {

	playlist := make([]*Playlist, 0)
	query := bson.M{
		"memId": memId,
	}

	res, err := noSql.GetCollection("playlist").Find(context.TODO(), query)
	if err != nil {
		return playlist, err
	}

	if err := res.All(context.TODO(), &playlist); err != nil {
		return playlist, err
	}

	return playlist, nil

}

/**
*	playlist id로 playlist 하나 찾음
 */
func FindOnePlaylistById(id primitive.ObjectID) (*Playlist, error) {

	playlist := &Playlist{}
	if err := noSql.GetCollection("playlist").FindOne(context.TODO(), bson.M{"_id": id}).Decode(playlist); err != nil {
		return nil, err
	}

	return playlist, nil
}

/**
*	update One Playlist
 */
func UpdateOnePlaylist(playlist *Playlist) error {
	query := bson.M{
		"_id": playlist.ID,
	}

	if _, err := noSql.GetCollection("playlist").UpdateOne(context.TODO(), query, playlist); err != nil {
		return err
	}

	return nil
}

func init() {
}
