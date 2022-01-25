package music_service

import (
	"context"
	"fmt"
	"io"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"

	musicpb "github.com/Park-Kwonsoo/moving-server/api/protos/v1/music"

	"github.com/Park-Kwonsoo/moving-server/pkg/database/nosql"
	errHandler "github.com/Park-Kwonsoo/moving-server/pkg/err-handler"

	nosqlModel "github.com/Park-Kwonsoo/moving-server/internal/models/nosql"
	sqlModel "github.com/Park-Kwonsoo/moving-server/internal/models/sql"
)

type MusicServer struct {
	musicpb.MusicServiceServer
}

/**
*	GetMusicDetail의 리턴 타입을 가져옴
 */
func getMusicDetailReturnType(e errHandler.ErrorRslt, code error, music *nosqlModel.Music) (*musicpb.GetMusicDetailRes, error) {

	if music == nil {
		return &musicpb.GetMusicDetailRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, code
	}

	return &musicpb.GetMusicDetailRes{
		RsltCd:  "00",
		RsltMsg: "Success",
		Music: &musicpb.Music{
			TrackNumber: uint64(music.TrackNumber),

			MusicId:       music.ID.String(),
			Title:         music.Title,
			Artist:        music.Artist,
			Album:         music.Album,
			Genre:         music.AlbumCoverUrl,
			MusicUrl:      music.MusicUrl,
			AlbumCoverUrl: music.AlbumCoverUrl,

			IsTitle: music.IsTitle,
		},
	}, nil
}

/**
* GetMusicByKeyword의 리턴 타입을 가져옴
 */
func getMusicByKeywordReturnType(e errHandler.ErrorRslt, code error, musicList []*nosqlModel.Music) (*musicpb.GetMusicByKeywordRes, error) {

	if code != nil {
		return &musicpb.GetMusicByKeywordRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, code
	}

	rsltList := make([]*musicpb.Music, 0)
	for _, music := range musicList {
		rsltList = append(rsltList, &musicpb.Music{
			MusicId:     music.ID.String(),
			TrackNumber: uint64(music.TrackNumber),

			Title:  music.Title,
			Artist: music.Artist,
			Album:  music.Album,
			Genre:  music.Genre,

			MusicUrl:      music.MusicUrl,
			AlbumCoverUrl: music.AlbumCoverUrl,

			IsTitle: music.IsTitle,
		})
	}

	return &musicpb.GetMusicByKeywordRes{
		RsltCd:       "00",
		RsltMsg:      "Success",
		SearchResult: rsltList,
	}, nil
}

func (s *MusicServer) GetMusicDetail(ctx context.Context, req *musicpb.GetMusicDetailReq) (*musicpb.GetMusicDetailRes, error) {

	musicId, _ := primitive.ObjectIDFromHex(req.MusicId)
	music, err := nosqlModel.FindOneMusicById(musicId)
	if err != nil {
		e, code := errHandler.NotFoundErr("GetMusicDetail : Not Found Music")
		return getMusicDetailReturnType(e, code, nil)
	}

	return getMusicDetailReturnType(errHandler.ErrorRslt{}, nil, music)
}

func (s *MusicServer) GetMusicByKeyword(ctx context.Context, req *musicpb.GetMusicByKeywordReq) (*musicpb.GetMusicByKeywordRes, error) {

	musicList, err := nosqlModel.FindAllMusicByKeyword(req.Keyword)
	if err != nil {
		e, code := errHandler.NotFoundErr("GetMusicByKeyword : Not Found Music")
		return getMusicByKeywordReturnType(e, code, nil)
	}

	return getMusicByKeywordReturnType(errHandler.ErrorRslt{}, nil, musicList)
}

func (s *MusicServer) AddNewMusic(stream musicpb.MusicService_AddNewMusicServer) error {

	memId := fmt.Sprintf("%v", stream.Context().Value("memId"))
	if len(memId) == 0 {
		_, code := errHandler.AuthorizedErr("AddNewMusic : Validate Token Error")
		return code
	}

	member, _ := sqlModel.FindOneMemberByMemId(memId)
	if member == nil {
		_, code := errHandler.NotFoundErr("AddNewMusic : Not Found User")
		return code
	}
	if member.MemPosition != "MANAGER" {
		_, code := errHandler.ForbiddenErr("AddNewMusic : Forbidden, Not A Manager")
		return code
	}

	album := &nosqlModel.Album{}
	music := &nosqlModel.Music{
		BaseType: nosql.BaseType{
			UseYn: "Y",
		},
	}
	//music file 생성
	musicFile, e := os.Create("temp_music")
	if e != nil {
		_, code := errHandler.BadRequestErr("AddNewMusic : Audio File Create Err")
		return code
	}

	for {
		res, e := stream.Recv()

		//더 이상 받아올 데이터가 없을 때 = 데이터를 전부 다 받아왔을 때 break
		if e == io.EOF {
			musicFile.Close()
			os.Remove("temp_music")
			break
		}
		if e != nil {
			musicFile.Close()
			os.Remove("temp_music")
			_, code := errHandler.BadRequestErr("AddNewMusic : Stream Receive Error")
			return code
		}

		musicFile.Write(res.Music)
		musicFile.Close()

		albumId, e := primitive.ObjectIDFromHex(res.AlbumId)
		if e != nil {
			musicFile.Close()
			os.Remove("temp_music")
			_, code := errHandler.BadRequestErr("AddNewMusic : Album Id")
			return code
		}
		album, e = nosqlModel.FindOneAlbumById(albumId)
		if e != nil {
			musicFile.Close()
			os.Remove("temp_music")
			_, code := errHandler.NotFoundErr("AddNewMusic : Not Found Album")
			return code
		}

		music.Title = res.Title
		music.Artist = res.Artist
		music.Genre = res.Genre
		music.TrackNumber = uint(res.TrackNumber)
		music.IsTitle = res.IsTitle

		os.Rename("temp_music", fmt.Sprintf("%s_%s_%s.wav", music.Title, music.Artist, music.Album))
	}

	err := nosqlModel.CreateNewMusic(music, album)
	if err != nil {
		_, code := errHandler.BadRequestErr("AddNewMusic : Bad Request")
		return code
	}

	return stream.SendAndClose(&musicpb.AddNewMusicRes{
		RsltCd:  "00",
		RsltMsg: "Success",
	})
}

func (s *MusicServer) AddNewAlbum(stream musicpb.MusicService_AddNewAlbumServer) error {

	memId := fmt.Sprintf("%v", stream.Context().Value("memId"))
	if len(memId) == 0 {
		_, code := errHandler.AuthorizedErr("AddNewAlbum : Validate Token Error")
		return code
	}

	member, _ := sqlModel.FindOneMemberByMemId(memId)
	if member == nil {
		_, code := errHandler.NotFoundErr("AddNewAlbum : Not Found User")
		return code
	}
	if member.MemPosition != "MANAGER" {
		_, code := errHandler.ForbiddenErr("AddNewAlbum : Forbidden, Not A Manager")
		return code
	}

	//toDo : GCS
	album := &nosqlModel.Album{
		Music: make([]nosqlModel.Music, 0),
		BaseType: nosql.BaseType{
			UseYn: "Y",
		},
	}

	for {
		res, e := stream.Recv()
		if e == io.EOF {
			break
		}
		if e != nil {
			_, code := errHandler.BadRequestErr("AddNewAlbum : Stream Receive Error")
			return code
		}

		album.Album = res.Album
		album.Artist = res.Artist
		album.Genre = res.Genre
		album.Description = res.Description
		album.AlbumCoverUrl = res.AlbumCoverUrl
	}

	err := nosqlModel.CreateNewAlbum(album)
	if err != nil {
		_, code := errHandler.BadRequestErr("AddNewAlbum : Bad Request")
		return code
	}

	return stream.SendAndClose(&musicpb.AddNewAlbumRes{
		RsltCd:  "00",
		RsltMsg: "Success",
	})
}

func (s *MusicServer) ListenMusic(ctx context.Context, req *musicpb.ListenMusicReq) (*musicpb.ListenMusicRes, error) {

	memId := fmt.Sprintf("%v", ctx.Value("memId"))
	if len(memId) == 0 {
		_, code := errHandler.AuthorizedErr("ListenMusic : Not Authorized")
		return &musicpb.ListenMusicRes{}, code
	}

	chanErr := make(chan error)

	go func() {
		member, _ := sqlModel.FindOneMemberByMemId(memId)
		if member == nil {
			_, code := errHandler.NotFoundErr("ListenMusic : Not Found User")
			chanErr <- code

			return
		}
		chanErr <- nil
	}()

	go func() {
		musicId, _ := primitive.ObjectIDFromHex(req.MusicId)
		music, err := nosqlModel.FindOneMusicById(musicId)
		if music == nil || err != nil {
			_, code := errHandler.NotFoundErr("ListenMusic : Not Found Music")
			chanErr <- code

			return
		}

		chanErr <- nil
	}()

	for i := 0; i < 2; i++ {
		if err := <-chanErr; err != nil {
			return &musicpb.ListenMusicRes{}, err
		}
	}

	musicListenLog := &sqlModel.MusicListenLog{
		MemId:   memId,
		MusicId: req.MusicId,
	}

	if err := sqlModel.CreateNewMusicListenLog(musicListenLog); err != nil {
		_, code := errHandler.BadRequestErr("ListenMusic : Bad Request")
		return &musicpb.ListenMusicRes{}, code
	}

	return &musicpb.ListenMusicRes{}, nil
}
