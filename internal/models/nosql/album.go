package nosql_model

import (
	"context"
	"time"

	noSql "github.com/Park-Kwonsoo/moving-server/pkg/database/nosql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Album struct {
	noSql.BaseType `bson:",inline"`
	Album          string  `bson:"album"`
	Artist         string  `bson:"artist"`
	Genre          string  `bson:"genre"`
	Description    string  `bson:"description"`
	AlbumCoverUrl  string  `bson:"album_cover_url"`
	Music          []Music `bson:"music"`
}

/**
*	CreateNewAlbum
 */
func CreateNewAlbum(album *Album) error {

	album.ID = primitive.NewObjectID()
	album.CreatedAt = time.Now()
	album.UpdatedAt = time.Now()

	_, err := noSql.GetCollection("album").InsertOne(context.TODO(), album)

	return err
}

/**
*	Add Music In Album
 */
func AddMusicInAlbum(album *Album, music ...Music) error {

	update := bson.M{
		"$push": bson.M{
			"music": bson.M{
				"$each": append(make([]Music, 0), music...),
			},
		},
	}

	if _, err := noSql.GetCollection("album").UpdateOne(context.TODO(), bson.M{
		"_id": album.ID,
	}, update); err != nil {
		return err
	}

	return nil
}

/**
*	Find One Album By Id
*	Return : Album
 */
func FindOneAlbumById(id primitive.ObjectID) (*Album, error) {

	album := &Album{}
	if err := noSql.GetCollection("album").FindOne(context.TODO(), bson.M{
		"_id": id,
	}).Decode(album); err != nil {
		return nil, err
	}

	return album, nil
}

func init() {
}
