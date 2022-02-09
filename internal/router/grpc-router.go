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

	sqlModel "github.com/Park-Kwonsoo/moving-server/internal/models/sql"
	redisClient "github.com/Park-Kwonsoo/moving-server/pkg/cache-server"
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
	grpcPort = ":9000"
)

//유저 인증 JWT 토큰 Interceptor : token값을 decode하여 memId를 전달
func authInterceptor(ctx context.Context) (context.Context, error) {

	token, _ := grpc_auth.AuthFromMD(ctx, "bearer")
	if len(token) == 0 {
		newCtx := context.WithValue(ctx, "memId", "")
		return newCtx, nil
	}

	//cache를 확인
	value, _ := redisClient.GetCache(token)
	if len(value) > 0 {
		//cache가 있으면 해당 cache value값으로 리턴
		newCtx := context.WithValue(ctx, "memId", string(value[:]))
		return newCtx, nil
	}

	//jwt로 토큰 검증
	memId, err := jwtUtil.ValidateToken(token)
	if err != nil {
		//검증 실패시 memId return
		newCtx := context.WithValue(ctx, "memId", "")
		return newCtx, nil
	}

	//현재 사용중인 유저인지 검증
	rslt, err := sqlModel.IsOneMemberExistByMemIdAndUseYn(memId, "Y")
	if !rslt {
		//탈퇴한 유저라면 memId로 등록하지 않는다.
		newCtx := context.WithValue(ctx, "memId", "")
		return newCtx, err
	}

	//만약 위의 과정 모두 통과한다면 redis에 key , value값 설정한 후 리턴
	redisClient.SetCache(token, memId)
	newCtx := context.WithValue(ctx, "memId", memId)
	return newCtx, err
}

//cache를 위한 custom unary interceptor
func customCacheUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		res, err := handler(ctx, req)
		return res, err
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
func SetupGRPCRouter() {
	lis, err := net.Listen("tcp", grpcPort)
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

	log.Printf("Start gRPC Server on %s server", grpcPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalln("Failed to Open gRPC Server :", err)
	}
}
