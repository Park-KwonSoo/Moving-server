package router

import (
	"log"
	"net"

	"github.com/sirupsen/logrus"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	authpb "github.com/Park-Kwonsoo/moving-server/api/protos/v1/auth"
	auth_service "github.com/Park-Kwonsoo/moving-server/services/auth"

	memberpb "github.com/Park-Kwonsoo/moving-server/api/protos/v1/member"
	member_service "github.com/Park-Kwonsoo/moving-server/services/member"

	playlistpb "github.com/Park-Kwonsoo/moving-server/api/protos/v1/playlist"
	playlist_service "github.com/Park-Kwonsoo/moving-server/services/playlist"

	musicpb "github.com/Park-Kwonsoo/moving-server/api/protos/v1/music"
	music_service "github.com/Park-Kwonsoo/moving-server/services/music"
)

const (
	port = ":9000"
)

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
			//logging interceptor
			grpc_logrus.UnaryServerInterceptor(logrusEntry),

			//recovery interceptor : panic 발생해도 프로그램 종료 안됨
			grpc_recovery.UnaryServerInterceptor(),
		)),

		//streaming server interceptor middleware
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer()),
	)
	//서비스 등록
	registerService(s)
	reflection.Register(s)

	log.Printf("Start gRPC Server on %s server", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalln("Failed to Open gRPC Server :", err)
	}
}
