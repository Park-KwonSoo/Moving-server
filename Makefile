PROTOC = protoc -I=. \
	--go_out . --go_opt paths=source_relative \
	--go-grpc_out . --go-grpc_opt paths=source_relative

PROTO_DIR = api/protos/v1

all : login register

login :
	$(PROTOC) $(PROTO_DIR)/auth/login.proto

register :
	${PROTOC} $(PROTO_DIR)/auth/register.proto