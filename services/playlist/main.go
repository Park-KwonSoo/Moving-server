package playlist_service

import (
	"context"

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

	myPlayList := make([]*playlistpb.Playlist, 0)
	for i := 0; i < len(playlist); i++ {
		myPlayList = append(myPlayList, &playlistpb.Playlist{
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
