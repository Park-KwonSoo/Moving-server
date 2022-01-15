package auth_service

import (
	"context"
	"database/sql"
	"fmt"

	authpb "github.com/Park-Kwonsoo/moving-server/api/protos/v1/auth"

	db "github.com/Park-Kwonsoo/moving-server/models"
	jwtutility "github.com/Park-Kwonsoo/moving-server/pkg/jwt-utility"
)

type LoginServer struct {
	authpb.LoginServiceServer
}

type RegisterServer struct {
	authpb.RegisterServiceServer
}

func (s *RegisterServer) Register(ctx context.Context, req *authpb.RegisterReq) (*authpb.RegisterRes, error) {
	db.CreateNewUser(db.User{
		UserId: sql.NullString{
			String: req.UserId,
			Valid:  true,
		},
		Password: req.Password,
		UserType: req.RegisterType,
	})

	return &authpb.RegisterRes{
		RsltMsg: "Success",
		RsltCd:  "00",
	}, nil

}

func (s *LoginServer) Login(ctx context.Context, req *authpb.LoginReq) (*authpb.LoginRes, error) {

	userId := req.UserId

	user := db.FindUserByUserId(userId)
	if user == nil {
		return &authpb.LoginRes{
			RsltCd:  "88",
			RsltMsg: "Not Found User",
			Token:   "",
		}, nil
	}

	isValidatePw := user.ValidatePassword(req.Password)
	if !isValidatePw {
		return &authpb.LoginRes{
			RsltCd:  "77",
			RsltMsg: "Not Authorized",
			Token:   "",
		}, nil
	}

	token, err := jwtutility.GetJwtToken(userId)

	if err != nil {
		return &authpb.LoginRes{
			RsltCd:  "99",
			RsltMsg: "Login Failed",
			Token:   "",
		}, fmt.Errorf("login error")
	}

	return &authpb.LoginRes{
		RsltCd:  "00",
		RsltMsg: "Success",
		Token:   token,
	}, nil
}
