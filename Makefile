PROTOC = protoc -I=. \
	--go_out . --go_opt paths=source_relative \
	--go-grpc_out . --go-grpc_opt paths=source_relative

PROTOC_GRPC_GATEWAY = protoc -I . --grpc-gateway_out . \
    --grpc-gateway_opt logtostderr=true \
    --grpc-gateway_opt paths=source_relative \
    --grpc-gateway_opt generate_unbound_methods=true

PROTO_DIR_V1 = api/protos/v1
PROTO_DIR_V2 = api/protos/v2

all : login login_gateway register register_gateway member member_gateway playlist playlist_gateway music music_gateway
gateway : login_gateway register_gateway member_gateway playlist_gateway music_gateway


login :
	$(PROTOC) $(PROTO_DIR_V1)/auth/login.proto

login_gateway :
	$(PROTOC_GRPC_GATEWAY) $(PROTO_DIR_V1)/auth/login.proto

register :
	${PROTOC} $(PROTO_DIR_V1)/auth/register.proto

register_gateway :
	$(PROTOC_GRPC_GATEWAY) $(PROTO_DIR_V1)/auth/register.proto

member :
	${PROTOC} ${PROTO_DIR_V1}/member/member.proto

member_gateway :
	${PROTOC_GRPC_GATEWAY} ${PROTO_DIR_V1}/member/member.proto

playlist :
	${PROTOC} ${PROTO_DIR_V1}/playlist/playlist.proto

playlist_gateway :
	${PROTOC_GRPC_GATEWAY} ${PROTO_DIR_V1}/playlist/playlist.proto

music :
	${PROTOC} ${PROTO_DIR_V1}/music/music.proto

music_gateway :
	${PROTOC_GRPC_GATEWAY} ${PROTO_DIR_V1}/music/music.proto