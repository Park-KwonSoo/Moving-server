package playlist_service

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"

	playlistpb "github.com/Park-Kwonsoo/moving-server/api/protos/v1/playlist"

	"github.com/Park-Kwonsoo/moving-server/pkg/database/nosql"
	errHandler "github.com/Park-Kwonsoo/moving-server/pkg/err-handler"

	nosqlModel "github.com/Park-Kwonsoo/moving-server/internal/models/nosql"
	sqlModel "github.com/Park-Kwonsoo/moving-server/internal/models/sql"
)

type PlaylistServer struct {
	playlistpb.PlaylistServiceServer
}

//playlist return type
func getPlaylistReturnType(e errHandler.ErrorRslt, code error, playlist []*nosqlModel.Playlist) (*playlistpb.GetMyPlaylistRes, error) {

	if playlist == nil {
		return &playlistpb.GetMyPlaylistRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, code
	}

	myPlayList := make([]*playlistpb.SimplePlaylist, 0)
	for i := 0; i < len(playlist); i++ {
		myPlayList = append(myPlayList, &playlistpb.SimplePlaylist{
			Id:           playlist[i].ID.String(),
			CreatedAt:    playlist[i].CreatedAt.String(),
			UpdatedAt:    playlist[i].UpdatedAt.String(),
			PlaylistName: playlist[i].PlaylistName,
			NumOfMusics:  uint64(len(playlist[i].Music)),
		})
	}

	return &playlistpb.GetMyPlaylistRes{
		RsltCd:     "00",
		RsltMsg:    "Success",
		MyPlaylist: myPlayList,
	}, nil
}

//specific playlist return type
func getSpecificPlaylistReturnType(e errHandler.ErrorRslt, code error, myPlaylist *nosqlModel.Playlist, numberOfLike int) (*playlistpb.GetSpecificPlaylistRes, error) {

	if myPlaylist == nil {
		return &playlistpb.GetSpecificPlaylistRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, code
	}

	playlist := &playlistpb.SpecificPlaylist{
		Id:           myPlaylist.ID.String(),
		CreatedAt:    myPlaylist.CreatedAt.String(),
		UpdatedAt:    myPlaylist.UpdatedAt.String(),
		PlaylistName: myPlaylist.PlaylistName,
		NumOfLike:    uint64(numberOfLike),
		MusicList:    make([]*playlistpb.Music, 0),
	}

	for _, music := range myPlaylist.Music {
		playlist.MusicList = append(playlist.MusicList, &playlistpb.Music{
			MusicId:     music.ID.String(),
			TrackNumber: uint64(music.TrackNumber),

			Title:  music.Title,
			Artist: music.Artist,
			Album:  music.Album,
			Genre:  music.Genre,

			AlbumImg: music.AlbumCoverUrl,
			MusicUrl: music.MusicUrl,

			IsTitle: music.IsTitle,
		})
	}

	return &playlistpb.GetSpecificPlaylistRes{
		RsltCd:   "00",
		RsltMsg:  "Success",
		Playlist: playlist,
	}, nil
}

/*
* Get My Playlist
 */
func (s *PlaylistServer) GetMyPlaylist(ctx context.Context, req *playlistpb.GetMyPlaylistReq) (*playlistpb.GetMyPlaylistRes, error) {

	memId := ctx.Value("memId").(string)
	if len(memId) == 0 {
		e, code := errHandler.AuthorizedErr("GetMyPlaylist : Validate Token Error")
		return getPlaylistReturnType(e, code, nil)
	}

	myPlaylist, err := nosqlModel.FindAllPlaylistByMemberMemId(memId)
	if err != nil {
		e, code := errHandler.NotFoundErr("GetMyPlaylist : Not Found User's Playlist")
		return getPlaylistReturnType(e, code, nil)
	}

	return getPlaylistReturnType(errHandler.ErrorRslt{}, nil, myPlaylist)
}

/**
*	Get Specifin Playlist By Platlist Id
 */
func (s *PlaylistServer) GetSpecificPlaylist(ctx context.Context, req *playlistpb.GetSpecificPlaylistReq) (*playlistpb.GetSpecificPlaylistRes, error) {

	playlistId, _ := primitive.ObjectIDFromHex(req.PlaylistId)

	//플레이리스트 가져오는 채널
	playListChan := make(chan *nosqlModel.Playlist)
	//좋아요 수 가져오는 채널
	countChan := make(chan int)

	chanErr := make(chan error)
	chanRslt := make(chan *playlistpb.GetSpecificPlaylistRes)

	go func() {
		playlist, err := nosqlModel.FindOnePlaylistById(playlistId)
		if err != nil {
			e, code := errHandler.NotFoundErr("GetSpecificPlaylist : Not Found Playlist")
			rslt, _ := getSpecificPlaylistReturnType(e, code, nil, 0)

			chanErr <- code
			chanRslt <- rslt
			playListChan <- nil

			return
		}
		chanErr <- nil
		chanRslt <- nil
		playListChan <- playlist
	}()

	go func() {
		count, err := sqlModel.CountPlaylistLikePlaylistId(req.PlaylistId)
		if err != nil {
			e, code := errHandler.BadRequestErr("GetSpecificPlaylist : SQL Query Error")
			rslt, _ := getSpecificPlaylistReturnType(e, code, nil, 0)

			chanErr <- code
			chanRslt <- rslt
			countChan <- 0

			return
		}

		chanErr <- nil
		chanRslt <- nil
		countChan <- count
	}()

	for i := 0; i < 2; i++ {
		if err := <-chanErr; err != nil {
			return <-chanRslt, err
		} else {
			<-chanRslt
		}
	}

	return getSpecificPlaylistReturnType(errHandler.ErrorRslt{}, nil, <-playListChan, <-countChan)

}

/**
*	Create New Playlist
 */
func (s *PlaylistServer) CreateNewPlaylist(ctx context.Context, req *playlistpb.CreateNewPlaylistReq) (*playlistpb.CreateNewPlaylistRes, error) {

	memId := ctx.Value("memId").(string)
	if len(memId) == 0 {
		e, code := errHandler.AuthorizedErr("CreateNewPlaylist : Authorized User")

		return &playlistpb.CreateNewPlaylistRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, code
	}

	playlist := &nosqlModel.Playlist{
		PlaylistName: req.PlaylistName,
		MemId:        memId,
		Music:        make([]nosqlModel.Music, 0),
		BaseType: nosql.BaseType{
			UseYn: "Y",
		},
	}

	err := nosqlModel.CreateNewPlaylist(playlist)
	if err != nil {
		e, code := errHandler.BadRequestErr("CreateNewPlaylist : Bad Request")

		return &playlistpb.CreateNewPlaylistRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, code
	}

	return &playlistpb.CreateNewPlaylistRes{
		RsltCd:  "00",
		RsltMsg: "Success",
	}, nil
}

/**
*	Update Playlist
 */
func (s *PlaylistServer) UpdatePlaylist(ctx context.Context, req *playlistpb.UpdatePlaylistReq) (*playlistpb.UpdatePlaylistRes, error) {

	memId := ctx.Value("memId").(string)
	if len(memId) == 0 {
		e, code := errHandler.AuthorizedErr("UpdatePlaylist : Authorized User")

		return &playlistpb.UpdatePlaylistRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, code
	}

	playlistId, _ := primitive.ObjectIDFromHex(req.PlaylistId)
	playlist, _ := nosqlModel.FindOnePlaylistById(playlistId)
	if playlist.MemId != memId {
		e, code := errHandler.ForbiddenErr("UpdatePlaylist : Forbidden User")

		return &playlistpb.UpdatePlaylistRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, code
	}

	if len(req.PlaylistName) > 0 {
		playlist.PlaylistName = req.PlaylistName
	}

	err := nosqlModel.UpdateOnePlaylist(playlist)
	if err != nil {
		e, code := errHandler.BadRequestErr("UpdatePlaylist : Update Failed")

		return &playlistpb.UpdatePlaylistRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, code
	}

	return &playlistpb.UpdatePlaylistRes{
		RsltCd:  "00",
		RsltMsg: "Success",
	}, nil
}

/**
*	Add Music in Playlist
 */
func (s *PlaylistServer) AddNewMusicInPlaylist(ctx context.Context, req *playlistpb.AddNewMusicInPlaylistReq) (*playlistpb.AddNewMusicInPlaylistRes, error) {

	memId := ctx.Value("memId").(string)
	if len(memId) == 0 {
		e, code := errHandler.AuthorizedErr("AddNewMusicInPlaylist : Authorized User")

		return &playlistpb.AddNewMusicInPlaylistRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, code
	}

	playlistId, _ := primitive.ObjectIDFromHex(req.PlaylistId)
	playlist, err := nosqlModel.FindOnePlaylistById(playlistId)
	if err != nil {
		e, code := errHandler.NotFoundErr("AddNewMusicInPlaylist : Not Found Playlist")

		return &playlistpb.AddNewMusicInPlaylistRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, code
	}

	if memId != playlist.MemId {
		e, code := errHandler.ForbiddenErr("AddNewMusicInPlaylist : Forbidden User")

		return &playlistpb.AddNewMusicInPlaylistRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, code
	}

	err = nosqlModel.AddMusicInPlaylist(playlist, req.MusicIdList...)
	if err != nil {
		e, code := errHandler.BadRequestErr("AddNewMusicInPlaylist : Bad Request")

		return &playlistpb.AddNewMusicInPlaylistRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, code
	}

	return &playlistpb.AddNewMusicInPlaylistRes{
		RsltCd:  "00",
		RsltMsg: "Success",
	}, nil
}

/**
*	toDo : Remove MusicId
 */
func (s *PlaylistServer) RemoveMusicInPlaylist(ctx context.Context, req *playlistpb.RemoveMusicInPlaylistReq) (*playlistpb.RemoveMusicInPlaylistRes, error) {

	memId := ctx.Value("memId").(string)
	if len(memId) == 0 {
		e, code := errHandler.AuthorizedErr("RemoveMusicInPlaylist : Authorized User")

		return &playlistpb.RemoveMusicInPlaylistRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, code
	}

	playlistId, _ := primitive.ObjectIDFromHex(req.PlaylistId)
	playlist, err := nosqlModel.FindOnePlaylistById(playlistId)
	if err != nil {
		e, code := errHandler.NotFoundErr("RemoveMusicInPlaylist : Not Found Playlist")

		return &playlistpb.RemoveMusicInPlaylistRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, code
	}

	if memId != playlist.MemId {
		e, code := errHandler.ForbiddenErr("RemoveMusicInPlaylist : Forbidden User")

		return &playlistpb.RemoveMusicInPlaylistRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, code
	}

	err = nosqlModel.RemoveMusicInPlaylist(playlist, req.MusicIdList...)
	if err != nil {
		log.Println(err)
		e, code := errHandler.BadRequestErr("RemoveMusicInPlaylist : Bad Request")

		return &playlistpb.RemoveMusicInPlaylistRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, code
	}

	return &playlistpb.RemoveMusicInPlaylistRes{
		RsltCd:  "00",
		RsltMsg: "Success",
	}, nil
}

type LikePlaylistStruct struct {
	LikePlaylistRes *playlistpb.LikePlaylistRes
	Err             error
}

func (s *PlaylistServer) LikePlaylist(ctx context.Context, req *playlistpb.LikePlaylistReq) (*playlistpb.LikePlaylistRes, error) {

	memId := ctx.Value("memId").(string)
	if len(memId) == 0 {
		e, code := errHandler.AuthorizedErr("LikePlaylist : Authorized User")

		return &playlistpb.LikePlaylistRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, code
	}

	//chanel
	chanRslt := make(chan *LikePlaylistStruct)

	go func() {
		playlistId, _ := primitive.ObjectIDFromHex(req.PlaylistId)
		playlist, err := nosqlModel.FindOnePlaylistById(playlistId)
		if err != nil {
			e, code := errHandler.BadRequestErr("LikePlaylist : Not Right PlaylistId")
			chanRslt <- &LikePlaylistStruct{
				LikePlaylistRes: &playlistpb.LikePlaylistRes{
					RsltCd:  e.RsltCd,
					RsltMsg: e.RsltMsg,
				},
				Err: code,
			}

			return

		}
		if playlist.MemId == memId {
			e, code := errHandler.BadRequestErr("LikePlaylist : Playlist Owner Can't Do like")
			chanRslt <- &LikePlaylistStruct{
				LikePlaylistRes: &playlistpb.LikePlaylistRes{
					RsltCd:  e.RsltCd,
					RsltMsg: e.RsltMsg,
				},
				Err: code,
			}

			return
		}

		chanRslt <- nil

	}()

	go func() {
		hadLike, err := sqlModel.HavePlaylistLikeByMemIdAndPlaylistId(memId, req.PlaylistId)
		if err != nil {
			e, code := errHandler.BadRequestErr("LikePlaylist : Find Like Log Error")
			chanRslt <- &LikePlaylistStruct{
				LikePlaylistRes: &playlistpb.LikePlaylistRes{
					RsltCd:  e.RsltCd,
					RsltMsg: e.RsltMsg,
				},
				Err: code,
			}

			return
		}
		if !hadLike {
			e, code := errHandler.BadRequestErr("LikePlaylist : You Already Like")
			chanRslt <- &LikePlaylistStruct{
				LikePlaylistRes: &playlistpb.LikePlaylistRes{
					RsltCd:  e.RsltCd,
					RsltMsg: e.RsltMsg,
				},
				Err: code,
			}

			return
		}

		chanRslt <- nil
	}()

	for i := 0; i < 2; i++ {
		if rslt := <-chanRslt; rslt != nil {
			return rslt.LikePlaylistRes, rslt.Err
		}
	}

	playlistLike := &sqlModel.PlaylistLike{
		MemId:      memId,
		PlaylistId: req.PlaylistId,
	}

	err := sqlModel.CreateNewPlaylistLike(playlistLike)
	if err != nil {
		e, code := errHandler.BadRequestErr("LikePlaylist : Create PlaylistLike error")

		return &playlistpb.LikePlaylistRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, code
	}

	return &playlistpb.LikePlaylistRes{
		RsltCd:  "00",
		RsltMsg: "Success",
	}, nil
}
