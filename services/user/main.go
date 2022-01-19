package user_service

import (
	"context"

	userpb "github.com/Park-Kwonsoo/moving-server/api/protos/v1/user"

	errHandler "github.com/Park-Kwonsoo/moving-server/pkg/err-handler"

	db "github.com/Park-Kwonsoo/moving-server/models"
	jwtUtil "github.com/Park-Kwonsoo/moving-server/pkg/jwt-utility"
)

type UserServer struct {
	userpb.UserServiceServer
}

func (s *UserServer) GetMyProfile(ctx context.Context, req *userpb.GetMyProfileReq) (*userpb.GetMyProfileRes, error) {
	userId, err := jwtUtil.ValidateToken(req.Token)
	if err != nil {
		e := errHandler.AuthorizedErr()

		return &userpb.GetMyProfileRes{
			RsltCd:    e.RsltCd,
			RsltMsg:   e.RsltMsg,
			MyProfile: nil,
		}, nil
	}

	profile, err := db.FindProfileByUserId(userId)
	if profile == nil {
		e := errHandler.NotFoundErr()

		return &userpb.GetMyProfileRes{
			RsltCd:    e.RsltCd,
			RsltMsg:   e.RsltMsg,
			MyProfile: nil,
		}, nil
	}

	return &userpb.GetMyProfileRes{
		RsltCd:  "00",
		RsltMsg: "Success",
		MyProfile: &userpb.Profile{
			Id:        int32(profile.ID),
			CreatedAt: profile.CreatedAt.String(),
			UpdatedAt: profile.UpdatedAt.String(),
			DeletedAt: profile.DeletedAt.String(),

			User: &userpb.User{
				Id:        int32(profile.User.ID),
				CreatedAt: profile.User.CreatedAt.String(),
				UpdatedAt: profile.User.UpdatedAt.GoString(),
				DeletedAt: profile.User.DeletedAt.String(),

				UserId:   profile.User.UserId.String,
				UserType: profile.User.UserType,
			},

			Name:       profile.Name,
			Gender:     profile.Gender,
			Birth:      profile.Birth,
			ProfileImg: profile.ProfileImage,
		},
	}, nil
}
