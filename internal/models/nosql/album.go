package nosql_model

import (
	"context"
	"time"

	noSql "github.com/Park-Kwonsoo/moving-server/pkg/database/nosql"
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

func init() {
}
