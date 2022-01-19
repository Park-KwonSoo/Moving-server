package auth_service

import (
	"context"
	"database/sql"

	authpb "github.com/Park-Kwonsoo/moving-server/api/protos/v1/auth"

	errHandler "github.com/Park-Kwonsoo/moving-server/pkg/err-handler"

	db "github.com/Park-Kwonsoo/moving-server/models"
	jwtUtil "github.com/Park-Kwonsoo/moving-server/pkg/jwt-utility"
)

type LoginServer struct {
	authpb.LoginServiceServer
}

type RegisterServer struct {
	authpb.RegisterServiceServer
}

func (s *RegisterServer) Register(ctx context.Context, req *authpb.RegisterReq) (*authpb.RegisterRes, error) {

	user := db.User{
		UserId: sql.NullString{
			String: req.UserId,
			Valid:  true,
		},
		Password: req.Password,
		UserType: req.RegisterType,
	}

	err := db.CreateNewUser(user)
	if err != nil {
		e := errHandler.ConflictErr()

		return &authpb.RegisterRes{
			RsltMsg: e.RsltMsg,
			RsltCd:  e.RsltCd,
		}, nil
	}

	profile := db.Profile{
		User:         user,
		Name:         req.Name,
		Birth:        req.Birth,
		Gender:       req.Gender,
		ProfileImage: req.ProfileImg,
	}

	err = db.CreateNewProfile(profile)
	if err != nil {
		e := errHandler.ConflictErr()

		return &authpb.RegisterRes{
			RsltMsg: e.RsltMsg,
			RsltCd:  e.RsltCd,
		}, nil
	}

	return &authpb.RegisterRes{
		RsltMsg: "Success",
		RsltCd:  "00",
	}, nil

}

func (s *LoginServer) Login(ctx context.Context, req *authpb.LoginReq) (*authpb.LoginRes, error) {

	userId := req.UserId

	user, err := db.FindUserByUserId(userId)
	if err != nil {
		e := errHandler.NotFoundErr()

		return &authpb.LoginRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
			Token:   "",
		}, nil
	}

	isValidatePw, err := user.ValidatePassword(req.Password)
	if !isValidatePw {
		e := errHandler.AuthorizedErr()

		return &authpb.LoginRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
			Token:   "",
		}, nil
	}

	token, err := jwtUtil.GetJwtToken(userId)

	if err != nil {
		e := errHandler.AuthorizedErr()

		return &authpb.LoginRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
			Token:   "",
		}, nil
	}

	return &authpb.LoginRes{
		RsltCd:  "00",
		RsltMsg: "Success",
		Token:   token,
	}, nil
}
