package router

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	authpb "github.com/Park-Kwonsoo/moving-server/api/protos/v1/auth"
	auth_service "github.com/Park-Kwonsoo/moving-server/services/auth"

	memberpb "github.com/Park-Kwonsoo/moving-server/api/protos/v1/member"
	member_service "github.com/Park-Kwonsoo/moving-server/services/member"
)

const (
	port = ":9000"
)

//service 등록
func registerService(s *grpc.Server) {
	authpb.RegisterRegisterServiceServer(s, &auth_service.RegisterServer{})
	authpb.RegisterLoginServiceServer(s, &auth_service.LoginServer{})

	memberpb.RegisterMemberServiceServer(s, &member_service.MemberServer{})
}

func SetupRouter() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln("Failed to Listen : %v", err)
	}

	s := grpc.NewServer()
	registerService(s)
	reflection.Register(s)

	log.Printf("Start gRPC Server on %s server", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalln("Failed to Open gRPC Server : %v", err)
	}
}
