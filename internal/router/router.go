package router

import (
	"context"
	"log"
	"net"

	"github.com/sirupsen/logrus"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	jwtUtil "github.com/Park-Kwonsoo/moving-server/pkg/jwt-utility"

	authpb "github.com/Park-Kwonsoo/moving-server/api/protos/v1/auth"
	auth_service "github.com/Park-Kwonsoo/moving-server/internal/services/auth"

	memberpb "github.com/Park-Kwonsoo/moving-server/api/protos/v1/member"
	member_service "github.com/Park-Kwonsoo/moving-server/internal/services/member"

	playlistpb "github.com/Park-Kwonsoo/moving-server/api/protos/v1/playlist"
	playlist_service "github.com/Park-Kwonsoo/moving-server/internal/services/playlist"

	musicpb "github.com/Park-Kwonsoo/moving-server/api/protos/v1/music"
	music_service "github.com/Park-Kwonsoo/moving-server/internal/services/music"
)

const (
	port = ":9000"
)

//유저 인증 JWT 토큰 Interceptor : token값을 decode하여 memId를 전달
func authInterceptor(ctx context.Context) (context.Context, error) {

	token, _ := grpc_auth.AuthFromMD(ctx, "bearer")
	if len(token) == 0 {
		newCtx := context.WithValue(ctx, "memId", "")
		return newCtx, nil
	}

	memId, _ := jwtUtil.ValidateToken(token)

	newCtx := context.WithValue(ctx, "memId", memId)
	return newCtx, nil
}

//cache를 위한 custom unary interceptor
func customCacheUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		return resp, err
	}
}

//cache custom streaming interceptor
func customCacheStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return handler(srv, ss)
	}
}

//service 등록
func registerService(s *grpc.Server) {
	authpb.RegisterRegisterServiceServer(s, &auth_service.RegisterServer{})
	authpb.RegisterLoginServiceServer(s, &auth_service.LoginServer{})

	memberpb.RegisterMemberServiceServer(s, &member_service.MemberServer{})

	playlistpb.RegisterPlaylistServiceServer(s, &playlist_service.PlaylistServer{})

	musicpb.RegisterMusicServiceServer(s, &music_service.MusicServer{})

}

//grpc Router 등록
func SetupRouter() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln("Failed to Listen :", err)
	}

	//log를 확인하기 위한 logrus
	logrus.ErrorKey = "grpc.error"
	logrusEntry := logrus.NewEntry(logrus.StandardLogger())

	s := grpc.NewServer(
		//unary server interceptor middleware
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			//auth interceptor
			grpc_auth.UnaryServerInterceptor(authInterceptor),

			//logging interceptor
			grpc_logrus.UnaryServerInterceptor(logrusEntry),

			//recovery interceptor : panic 발생해도 프로그램 종료 안됨
			grpc_recovery.UnaryServerInterceptor(),

			//cache interceptor : custom
			customCacheUnaryInterceptor(),
		)),

		//streaming server interceptor middleware
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			//auth_interceptor
			grpc_auth.StreamServerInterceptor(authInterceptor),

			//loggin_interceptor
			grpc_logrus.StreamServerInterceptor(logrusEntry),

			//recovery interceptor
			grpc_recovery.StreamServerInterceptor(),

			//cache interceptor : custom
			customCacheStreamInterceptor(),
		)),
	)
	//서비스 등록
	registerService(s)
	reflection.Register(s)

	log.Printf("Start gRPC Server on %s server", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalln("Failed to Open gRPC Server :", err)
	}
}
