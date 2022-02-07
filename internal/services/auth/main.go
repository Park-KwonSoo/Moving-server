package auth_service

import (
	"context"
	"database/sql"

	authpb "github.com/Park-Kwonsoo/moving-server/api/protos/v1/auth"

	errHandler "github.com/Park-Kwonsoo/moving-server/pkg/err-handler"

	sqlModel "github.com/Park-Kwonsoo/moving-server/internal/models/sql"
	jwtUtil "github.com/Park-Kwonsoo/moving-server/pkg/jwt-utility"
)

type LoginServer struct {
	authpb.LoginServiceServer
}

type RegisterServer struct {
	authpb.RegisterServiceServer
}

func (s *RegisterServer) Register(ctx context.Context, req *authpb.RegisterReq) (*authpb.RegisterRes, error) {

	if req.Password != req.PasswordCheck {
		e, code := errHandler.BadRequestErr("Register : Password Not Equal")

		return &authpb.RegisterRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, code
	}

	member := &sqlModel.Member{
		MemId: sql.NullString{
			String: req.MemId,
			Valid:  true,
		},
		Password: req.Password,
		MemType:  req.RegisterType,
	}

	err := sqlModel.CreateNewMember(member)
	if err != nil {
		e, code := errHandler.ConflictErr("Register : Conflict")

		return &authpb.RegisterRes{
			RsltMsg: e.RsltMsg,
			RsltCd:  e.RsltCd,
		}, code
	}

	profile := &sqlModel.Profile{
		Member:       *member,
		Name:         req.Name,
		Birth:        req.Birth,
		Gender:       req.Gender,
		ProfileImage: req.ProfileImg,
	}

	err = sqlModel.CreateNewProfile(profile)
	if err != nil {
		e, code := errHandler.ConflictErr("Register : Conflict")

		return &authpb.RegisterRes{
			RsltMsg: e.RsltMsg,
			RsltCd:  e.RsltCd,
		}, code
	}

	return &authpb.RegisterRes{
		RsltMsg: "Success",
		RsltCd:  "00",
	}, nil

}

func (s *LoginServer) Login(ctx context.Context, req *authpb.LoginReq) (*authpb.LoginRes, error) {

	memId := req.MemId
	member, err := sqlModel.FindOneMemberByMemId(memId)
	if err != nil {
		e, code := errHandler.NotFoundErr("Login : Not Found User")

		return &authpb.LoginRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, code
	}

	isValidatePw, err := member.ValidatePassword(req.Password)
	if !isValidatePw {
		e, code := errHandler.AuthorizedErr("Login : Failed Password")

		return &authpb.LoginRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, code
	}

	token, err := jwtUtil.GenerateJwtToken(memId)

	if err != nil {
		e, code := errHandler.AuthorizedErr("Login : Token Authorized")

		return &authpb.LoginRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, code
	}

	return &authpb.LoginRes{
		RsltCd:  "00",
		RsltMsg: "Success",
		Token:   token,
	}, nil
}

func (s *LoginServer) PasswordCheck(ctx context.Context, req *authpb.PasswordCheckReq) (*authpb.PasswordCheckRes, error) {

	memId := ctx.Value("memId").(string)
	if len(memId) == 0 {
		e, code := errHandler.AuthorizedErr("JWT Token Validate : Authorized Error")

		return &authpb.PasswordCheckRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, code
	}

	member, err := sqlModel.FindOneMemberByMemId(memId)
	if err != nil {
		e, code := errHandler.NotFoundErr("PasswordCheck : Not Found User")

		return &authpb.PasswordCheckRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, code
	}

	isChecked, err := member.ValidatePassword(req.OldPassword)

	return &authpb.PasswordCheckRes{
		RsltCd:    "00",
		RsltMsg:   "Success",
		IsChecked: isChecked,
	}, nil
}

func (s *LoginServer) PasswordChange(ctx context.Context, req *authpb.PasswordChangeReq) (*authpb.PasswordChangeRes, error) {

	memId := ctx.Value("memId").(string)
	if len(memId) == 0 {
		e, code := errHandler.AuthorizedErr("JWT Token Validate : Authorized Error")

		return &authpb.PasswordChangeRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, code
	}

	member, err := sqlModel.FindOneMemberByMemId(memId)
	if err != nil {
		e, code := errHandler.NotFoundErr("PasswordChange : Not Found User")

		return &authpb.PasswordChangeRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, code
	}

	err = member.ChangePassword(req.NewPassword)
	if err != nil {
		e, code := errHandler.ForbiddenErr("PasswordChange : Forbidden")

		return &authpb.PasswordChangeRes{
			RsltCd:  e.RsltCd,
			RsltMsg: e.RsltMsg,
		}, code
	}

	return &authpb.PasswordChangeRes{
		RsltCd:  "00",
		RsltMsg: "Success",
	}, nil
}
