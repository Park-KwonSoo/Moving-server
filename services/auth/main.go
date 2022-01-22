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

	member := &db.Member{
		MemId: sql.NullString{
			String: req.MemId,
			Valid:  true,
		},
		Password: req.Password,
		MemType:  req.RegisterType,
	}

	err := db.CreateNewMember(member)
	if err != nil {
		e := errHandler.ConflictErr("Register : Conflict")

		return &authpb.RegisterRes{
			RsltMsg: e.RsltMsg,
			RsltCd:  e.RsltCd,
		}, nil
	}

	member, err = db.FindOneMemberByMemId(req.MemId)

	profile := &db.Profile{
		Member:       *member,
		Name:         req.Name,
		Birth:        req.Birth,
		Gender:       req.Gender,
		ProfileImage: req.ProfileImg,
	}

	err = db.CreateNewProfile(profile)
	if err != nil {
		e := errHandler.ConflictErr("Register : Conflict")

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

	memId := req.MemId

	member, err := db.FindOneMemberByMemId(memId)
	if err != nil {
		e := errHandler.NotFoundErr("Login : Not Found User")

		return &authpb.LoginRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
			Token:   "",
		}, nil
	}

	isValidatePw, err := member.ValidatePassword(req.Password)
	if !isValidatePw {
		e := errHandler.AuthorizedErr("Login : Failed Password")

		return &authpb.LoginRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
			Token:   "",
		}, nil
	}

	token, err := jwtUtil.GetJwtToken(memId)

	if err != nil {
		e := errHandler.AuthorizedErr("Login : Token Authorized")

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

func (s *LoginServer) PasswordCheck(ctx context.Context, req *authpb.PasswordCheckReq) (*authpb.PasswordCheckRes, error) {

	memId, err := jwtUtil.ValidateToken(req.Token)
	if err != nil {
		e := errHandler.AuthorizedErr("PasswordCheck : Validate Token Error")

		return &authpb.PasswordCheckRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, nil
	}

	member, err := db.FindOneMemberByMemId(memId)
	if err != nil {
		e := errHandler.NotFoundErr("PasswordCheck : Not Found User")

		return &authpb.PasswordCheckRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, nil
	}

	isChecked, err := member.ValidatePassword(req.OldPassword)

	return &authpb.PasswordCheckRes{
		RsltCd:    "00",
		RsltMsg:   "Success",
		IsChecked: isChecked,
	}, nil
}

func (s *LoginServer) PasswordChange(ctx context.Context, req *authpb.PasswordChangeReq) (*authpb.PasswordChangeRes, error) {

	memId, err := jwtUtil.ValidateToken(req.Token)
	if err != nil {
		e := errHandler.AuthorizedErr("PasswordChange : Validate Token Error")

		return &authpb.PasswordChangeRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, nil
	}

	member, err := db.FindOneMemberByMemId(memId)
	if err != nil {
		e := errHandler.NotFoundErr("PasswordChange : Not Found User")

		return &authpb.PasswordChangeRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, nil
	}

	err = member.ChangePassword(req.NewPassword)
	if err != nil {
		e := errHandler.ForbiddenErr("PasswordChange : Forbidden")

		return &authpb.PasswordChangeRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, nil
	}

	return &authpb.PasswordChangeRes{
		RsltCd:  "00",
		RsltMsg: "Success",
	}, nil
}
