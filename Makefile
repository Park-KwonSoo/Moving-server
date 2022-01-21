PROTOC = protoc -I=. \
	--go_out . --go_opt paths=source_relative \
	--go-grpc_out . --go-grpc_opt paths=source_relative

PROTO_DIR_V1 = api/protos/v1
PROTO_DIR_V2 = api/protos/v2

all : login register member

login :
	$(PROTOC) $(PROTO_DIR_V1)/auth/login.proto

register :
	${PROTOC} $(PROTO_DIR_V1)/auth/register.proto

member :
	${PROTOC} ${PROTO_DIR_V1}/member/member.proto