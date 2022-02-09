package router

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authpb "github.com/Park-Kwonsoo/moving-server/api/protos/v1/auth"
	memberpb "github.com/Park-Kwonsoo/moving-server/api/protos/v1/member"
	musicpb "github.com/Park-Kwonsoo/moving-server/api/protos/v1/music"
	playlistpb "github.com/Park-Kwonsoo/moving-server/api/protos/v1/playlist"
)

const (
	restPort = ":9001"
)

func registerServiceForRest(ctx context.Context, mux *runtime.ServeMux, opt []grpc.DialOption) {

	endpoint := "localhost" + grpcPort

	authpb.RegisterLoginServiceHandlerFromEndpoint(
		ctx, mux, endpoint, opt,
	)

	authpb.RegisterRegisterServiceHandlerFromEndpoint(
		ctx, mux, endpoint, opt,
	)

	memberpb.RegisterMemberServiceHandlerFromEndpoint(
		ctx, mux, endpoint, opt,
	)

	playlistpb.RegisterPlaylistServiceHandlerFromEndpoint(
		ctx, mux, endpoint, opt,
	)

	musicpb.RegisterMusicServiceHandlerFromEndpoint(
		ctx, mux, endpoint, opt,
	)

}

//grpc-gateway를 이용하여 9001번에 RESTful Server
func SetupRESTRouter() {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	opt := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	//http endpoint에 service를 등록
	registerServiceForRest(ctx, mux, opt)

	log.Printf("Start REST Server on %s server", restPort)
	if err := http.ListenAndServe(restPort, mux); err != nil {
		log.Println(err)
	}
}
