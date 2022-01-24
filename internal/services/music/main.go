package music_service

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"

	musicpb "github.com/Park-Kwonsoo/moving-server/api/protos/v1/music"

	errHandler "github.com/Park-Kwonsoo/moving-server/pkg/err-handler"

	db "github.com/Park-Kwonsoo/moving-server/internal/models"
)

type MusicServer struct {
	musicpb.MusicServiceServer
}

/**
*	GetMusicDetail의 리턴 타입을 가져옴
 */
func getMusicDetailReturnType(e errHandler.ErrorRslt, code error, music *db.Music) (*musicpb.GetMusicDetailRes, error) {

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
			MusicId:     uint64(music.ID),
			TrackNumber: uint64(music.TrackNumber),

			Title:    music.Title,
			Artist:   music.Artist,
			Album:    music.Album,
			Genre:    music.Genre,
			AlbumImg: music.AlbumImg,
			MusicUrl: music.MusicUrl,

			IsTitle: music.IsTitle,
		},
	}, nil
}

/**
* GetMusicByKeyword의 리턴 타입을 가져옴
 */
func getMusicByKeywordReturnType(e errHandler.ErrorRslt, code error, musicList []*db.Music) (*musicpb.GetMusicByKeywordRes, error) {

	if len(musicList) == 0 {
		return &musicpb.GetMusicByKeywordRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, code
	}

	rsltList := make([]*musicpb.Music, 0)
	for _, music := range musicList {
		rsltList = append(rsltList, &musicpb.Music{
			MusicId:     uint64(music.ID),
			TrackNumber: uint64(music.TrackNumber),

			Title:    music.Title,
			Artist:   music.Artist,
			Album:    music.Album,
			Genre:    music.Genre,
			AlbumImg: music.AlbumImg,
			MusicUrl: music.MusicUrl,

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

	music, err := db.FindOneMusicById(uint(req.MusicId))
	if err != nil {
		e, code := errHandler.NotFoundErr("GetMusicDetail : Not Found Music")
		return getMusicDetailReturnType(e, code, nil)
	}

	return getMusicDetailReturnType(errHandler.ErrorRslt{}, nil, music)
}

func (s *MusicServer) GetMusicByKeyword(ctx context.Context, req *musicpb.GetMusicByKeywordReq) (*musicpb.GetMusicByKeywordRes, error) {

	musicList, err := db.FindAllMusicByKeyword(req.Keyword)
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

	member, _ := db.FindOneMemberByMemId(memId)
	if member == nil {
		_, code := errHandler.NotFoundErr("AddNewMusic : Not Found User")
		return code
	}
	// if member.MemPosition != "MANAGER" {
	// 	_, code := errHandler.ForbiddenErr("AddNewMusic : Forbidden, Not A Manager")
	// 	return code
	// }

	//music file 생성
	music, e := os.Create("temp_music")
	if e != nil {
		logrus.Error(e)
		_, code := errHandler.BadRequestErr("AddNewMusic : Audio File Create Err")
		return code
	}

	//img 파일 생성
	img, e := os.Create("temp_image")
	if e != nil {
		logrus.Error(e)
		_, code := errHandler.BadRequestErr("AddNewMusic : Iamge File Create Err")
		//rslt 채널에 bad request 에러 송신
		return code
	}

	for {
		res, e := stream.Recv()

		//더 이상 받아올 데이터가 없을 때 = 데이터를 전부 다 받아왔을 때 break
		if e == io.EOF {
			logrus.Info("AddNewMusic : Recieve Done")
			music.Close()
			os.Remove("temp_music")
			img.Close()
			os.Remove("temp_image")
			break
		}
		if e != nil {
			logrus.Error(e)
			music.Close()
			os.Remove("temp_music")
			img.Close()
			os.Remove("temp_image")
			_, code := errHandler.BadRequestErr("AddNewMusic : Stream Receive Err")
			return code
		}

		music.Write(res.Music)
		img.Write(res.Music)

		music.Close()
		img.Close()

		//toDo : File GCS Upload And Make DB
		os.Rename("temp_music", fmt.Sprintf("%s_%s_%s_%s.wav", res.Title, res.Artist, res.Album, res.Genre))
		os.Rename("temp_image", fmt.Sprintf("%s_%s_%s.png", res.Artist, res.Album, res.Genre))

	}

	return stream.SendAndClose(&musicpb.AddNewMusicRes{
		RsltCd:  "00",
		RsltMsg: "Success",
	})
}
