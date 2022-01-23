package playlist_service

import (
	"context"
	"database/sql"

	playlistpb "github.com/Park-Kwonsoo/moving-server/api/protos/v1/playlist"

	errHandler "github.com/Park-Kwonsoo/moving-server/pkg/err-handler"

	db "github.com/Park-Kwonsoo/moving-server/models"
	jwtUtil "github.com/Park-Kwonsoo/moving-server/pkg/jwt-utility"
)

type PlaylistServer struct {
	playlistpb.PlaylistServiceServer
}

//playlist return type
func getPlaylistReturnType(e errHandler.ErrorRslt, playlist []*db.Playlist) (*playlistpb.GetMyPlaylistRes, error) {

	if playlist == nil {
		return &playlistpb.GetMyPlaylistRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, nil
	}

	myPlayList := make([]*playlistpb.SimplePlaylist, 0)
	for i := 0; i < len(playlist); i++ {
		myPlayList = append(myPlayList, &playlistpb.SimplePlaylist{
			Id:           uint64(playlist[i].ID),
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
func getSpecificPlaylistReturnType(e errHandler.ErrorRslt, myPlaylist *db.Playlist) (*playlistpb.GetSpecificPlaylistRes, error) {

	if myPlaylist == nil {
		return &playlistpb.GetSpecificPlaylistRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, nil
	}

	playlist := &playlistpb.SpecificPlaylist{
		Id:           uint64(myPlaylist.ID),
		CreatedAt:    myPlaylist.CreatedAt.String(),
		UpdatedAt:    myPlaylist.UpdatedAt.String(),
		PlaylistName: myPlaylist.PlaylistName,
		MusicList:    make([]*playlistpb.Music, 0),
	}

	for _, music := range myPlaylist.Music {
		playlist.MusicList = append(playlist.MusicList, &playlistpb.Music{
			MusicId:     uint64(music.ID),
			TrackNumber: uint64(music.TrackNumber),

			Title:  music.Title,
			Artist: music.Artist,
			Album:  music.Album,
			Genre:  music.Genre,

			AlbumImg: music.AlbumImg,
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

	memId, err := jwtUtil.ValidateToken(req.Token)
	if err != nil {
		return getPlaylistReturnType(errHandler.AuthorizedErr("GetMyPlaylist : Validate Token Error"), nil)
	}

	myPlaylist, err := db.FindAllPlaylistByMemberMemId(memId)
	if err != nil {
		return getPlaylistReturnType(errHandler.NotFoundErr("GetMyPlaylist : Not Found User's Playlist"), nil)
	}

	return getPlaylistReturnType(errHandler.ErrorRslt{}, myPlaylist)
}

/**
*	Get Specifin Playlist By Platlist Id
 */
func (s *PlaylistServer) GetSpecificPlaylist(ctx context.Context, req *playlistpb.GetSpecificPlaylistReq) (*playlistpb.GetSpecificPlaylistRes, error) {

	playlist, err := db.FindOnePlaylistById(uint(req.PlaylistId))
	if err != nil {
		return getSpecificPlaylistReturnType(errHandler.NotFoundErr("GetSpecificPlaylist : Not Found Playlist"), nil)
	}

	return getSpecificPlaylistReturnType(errHandler.ErrorRslt{}, playlist)
}

/**
*	Create New Playlist
 */
func (s *PlaylistServer) CreateNewPlaylist(ctx context.Context, req *playlistpb.CreateNewPlaylistReq) (*playlistpb.CreateNewPlaylistRes, error) {

	memId, err := jwtUtil.ValidateToken(req.Token)
	if err != nil {
		e := errHandler.AuthorizedErr("CreateNewPlaylist : Authorized User")

		return &playlistpb.CreateNewPlaylistRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, nil
	}

	playlist := &db.Playlist{
		PlaylistName: req.PlaylistName,
		Member: db.Member{
			MemId: sql.NullString{
				String: memId,
			},
		},
	}

	err = db.CreateNewPlaylist(playlist)
	if err != nil {
		e := errHandler.BadRequestErr("CreateNewPlaylist : Bad Request")

		return &playlistpb.CreateNewPlaylistRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, nil
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
	memId, err := jwtUtil.ValidateToken(req.Token)
	if err != nil {
		e := errHandler.AuthorizedErr("UpdatePlaylist : Authorized User")

		return &playlistpb.UpdatePlaylistRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, nil
	}

	playlist, err := db.FindOnePlaylistById(uint(req.PlaylistId))
	if playlist.Member.MemId.String != memId {
		e := errHandler.ForbiddenErr("UpdatePlaylist : Forbidden User")

		return &playlistpb.UpdatePlaylistRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, nil
	}

	if len(req.PlaylistName) > 0 {
		playlist.PlaylistName = req.PlaylistName
	}

	err = db.UpdateOnePlaylist(playlist)
	if err != nil {
		e := errHandler.BadRequestErr("UpdatePlaylist : Update Failed")

		return &playlistpb.UpdatePlaylistRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, nil
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

	playlist, err := db.FindOnePlaylistById(uint(req.PlayListId))
	if err != nil {
		e := errHandler.NotFoundErr("AddNewMusicInPlaylist : Not Found Playlist")

		return &playlistpb.AddNewMusicInPlaylistRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, nil
	}

	err = db.AddMusicInPlaylist(playlist, req.MusicIdList...)
	if err != nil {
		e := errHandler.BadRequestErr("AddNewMusicInPlaylist : Bad Request")

		return &playlistpb.AddNewMusicInPlaylistRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, nil
	}

	return &playlistpb.AddNewMusicInPlaylistRes{
		RsltCd:  "00",
		RsltMsg: "Success",
	}, nil
}

/**
*	Remove MusicId
 */
func (s *PlaylistServer) RemoveMusicInPlaylist(ctx context.Context, req *playlistpb.RemoveMusicInPlaylistReq) (*playlistpb.RemoveMusicInPlaylistRes, error) {
	return &playlistpb.RemoveMusicInPlaylistRes{
		RsltCd:  "00",
		RsltMsg: "Success",
	}, nil
}
