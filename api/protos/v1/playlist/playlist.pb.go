// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.3
// source: api/protos/v1/playlist/playlist.proto

package playlist

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Playlist struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	NumOfMusics  uint64 `protobuf:"varint,2,opt,name=numOfMusics,proto3" json:"numOfMusics,omitempty"`
	CreatedAt    string `protobuf:"bytes,3,opt,name=createdAt,proto3" json:"createdAt,omitempty"`
	UpdatedAt    string `protobuf:"bytes,4,opt,name=updatedAt,proto3" json:"updatedAt,omitempty"`
	PlaylistName string `protobuf:"bytes,5,opt,name=playlistName,proto3" json:"playlistName,omitempty"`
}

func (x *Playlist) Reset() {
	*x = Playlist{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_protos_v1_playlist_playlist_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Playlist) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Playlist) ProtoMessage() {}

func (x *Playlist) ProtoReflect() protoreflect.Message {
	mi := &file_api_protos_v1_playlist_playlist_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Playlist.ProtoReflect.Descriptor instead.
func (*Playlist) Descriptor() ([]byte, []int) {
	return file_api_protos_v1_playlist_playlist_proto_rawDescGZIP(), []int{0}
}

func (x *Playlist) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Playlist) GetNumOfMusics() uint64 {
	if x != nil {
		return x.NumOfMusics
	}
	return 0
}

func (x *Playlist) GetCreatedAt() string {
	if x != nil {
		return x.CreatedAt
	}
	return ""
}

func (x *Playlist) GetUpdatedAt() string {
	if x != nil {
		return x.UpdatedAt
	}
	return ""
}

func (x *Playlist) GetPlaylistName() string {
	if x != nil {
		return x.PlaylistName
	}
	return ""
}

type GetMyPlaylistReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *GetMyPlaylistReq) Reset() {
	*x = GetMyPlaylistReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_protos_v1_playlist_playlist_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMyPlaylistReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMyPlaylistReq) ProtoMessage() {}

func (x *GetMyPlaylistReq) ProtoReflect() protoreflect.Message {
	mi := &file_api_protos_v1_playlist_playlist_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMyPlaylistReq.ProtoReflect.Descriptor instead.
func (*GetMyPlaylistReq) Descriptor() ([]byte, []int) {
	return file_api_protos_v1_playlist_playlist_proto_rawDescGZIP(), []int{1}
}

func (x *GetMyPlaylistReq) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type GetMyPlaylistRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RsltCd     string      `protobuf:"bytes,1,opt,name=rsltCd,proto3" json:"rsltCd,omitempty"`
	RsltMsg    string      `protobuf:"bytes,2,opt,name=rsltMsg,proto3" json:"rsltMsg,omitempty"`
	MyPlaylist []*Playlist `protobuf:"bytes,3,rep,name=myPlaylist,proto3" json:"myPlaylist,omitempty"`
}

func (x *GetMyPlaylistRes) Reset() {
	*x = GetMyPlaylistRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_protos_v1_playlist_playlist_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMyPlaylistRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMyPlaylistRes) ProtoMessage() {}

func (x *GetMyPlaylistRes) ProtoReflect() protoreflect.Message {
	mi := &file_api_protos_v1_playlist_playlist_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMyPlaylistRes.ProtoReflect.Descriptor instead.
func (*GetMyPlaylistRes) Descriptor() ([]byte, []int) {
	return file_api_protos_v1_playlist_playlist_proto_rawDescGZIP(), []int{2}
}

func (x *GetMyPlaylistRes) GetRsltCd() string {
	if x != nil {
		return x.RsltCd
	}
	return ""
}

func (x *GetMyPlaylistRes) GetRsltMsg() string {
	if x != nil {
		return x.RsltMsg
	}
	return ""
}

func (x *GetMyPlaylistRes) GetMyPlaylist() []*Playlist {
	if x != nil {
		return x.MyPlaylist
	}
	return nil
}

var File_api_protos_v1_playlist_playlist_proto protoreflect.FileDescriptor

var file_api_protos_v1_playlist_playlist_proto_rawDesc = []byte{
	0x0a, 0x25, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x76, 0x31, 0x2f,
	0x70, 0x6c, 0x61, 0x79, 0x6c, 0x69, 0x73, 0x74, 0x2f, 0x70, 0x6c, 0x61, 0x79, 0x6c, 0x69, 0x73,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x11, 0x76, 0x31, 0x2e, 0x70, 0x6c, 0x61, 0x79,
	0x6c, 0x69, 0x73, 0x74, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9c, 0x01, 0x0a, 0x08, 0x50,
	0x6c, 0x61, 0x79, 0x6c, 0x69, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x6e, 0x75, 0x6d, 0x4f, 0x66,
	0x4d, 0x75, 0x73, 0x69, 0x63, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x6e, 0x75,
	0x6d, 0x4f, 0x66, 0x4d, 0x75, 0x73, 0x69, 0x63, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x64, 0x41, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x22, 0x0a, 0x0c, 0x70, 0x6c, 0x61, 0x79, 0x6c, 0x69, 0x73,
	0x74, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x70, 0x6c, 0x61,
	0x79, 0x6c, 0x69, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x28, 0x0a, 0x10, 0x47, 0x65, 0x74,
	0x4d, 0x79, 0x50, 0x6c, 0x61, 0x79, 0x6c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x12, 0x14, 0x0a,
	0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f,
	0x6b, 0x65, 0x6e, 0x22, 0x81, 0x01, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x4d, 0x79, 0x50, 0x6c, 0x61,
	0x79, 0x6c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x73, 0x6c, 0x74,
	0x43, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x73, 0x6c, 0x74, 0x43, 0x64,
	0x12, 0x18, 0x0a, 0x07, 0x72, 0x73, 0x6c, 0x74, 0x4d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x72, 0x73, 0x6c, 0x74, 0x4d, 0x73, 0x67, 0x12, 0x3b, 0x0a, 0x0a, 0x6d, 0x79,
	0x50, 0x6c, 0x61, 0x79, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b,
	0x2e, 0x76, 0x31, 0x2e, 0x70, 0x6c, 0x61, 0x79, 0x6c, 0x69, 0x73, 0x74, 0x5f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x6c, 0x69, 0x73, 0x74, 0x52, 0x0a, 0x6d, 0x79, 0x50,
	0x6c, 0x61, 0x79, 0x6c, 0x69, 0x73, 0x74, 0x32, 0x6e, 0x0a, 0x0f, 0x50, 0x6c, 0x61, 0x79, 0x6c,
	0x69, 0x73, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x5b, 0x0a, 0x0d, 0x47, 0x65,
	0x74, 0x4d, 0x79, 0x50, 0x6c, 0x61, 0x79, 0x6c, 0x69, 0x73, 0x74, 0x12, 0x23, 0x2e, 0x76, 0x31,
	0x2e, 0x70, 0x6c, 0x61, 0x79, 0x6c, 0x69, 0x73, 0x74, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x47, 0x65, 0x74, 0x4d, 0x79, 0x50, 0x6c, 0x61, 0x79, 0x6c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71,
	0x1a, 0x23, 0x2e, 0x76, 0x31, 0x2e, 0x70, 0x6c, 0x61, 0x79, 0x6c, 0x69, 0x73, 0x74, 0x5f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x79, 0x50, 0x6c, 0x61, 0x79, 0x6c, 0x69,
	0x73, 0x74, 0x52, 0x65, 0x73, 0x22, 0x00, 0x42, 0x3e, 0x5a, 0x3c, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x50, 0x61, 0x72, 0x6b, 0x2d, 0x4b, 0x77, 0x6f, 0x6e, 0x73,
	0x6f, 0x6f, 0x2f, 0x4d, 0x6f, 0x76, 0x69, 0x6e, 0x67, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x70,
	0x6c, 0x61, 0x79, 0x6c, 0x69, 0x73, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_protos_v1_playlist_playlist_proto_rawDescOnce sync.Once
	file_api_protos_v1_playlist_playlist_proto_rawDescData = file_api_protos_v1_playlist_playlist_proto_rawDesc
)

func file_api_protos_v1_playlist_playlist_proto_rawDescGZIP() []byte {
	file_api_protos_v1_playlist_playlist_proto_rawDescOnce.Do(func() {
		file_api_protos_v1_playlist_playlist_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_protos_v1_playlist_playlist_proto_rawDescData)
	})
	return file_api_protos_v1_playlist_playlist_proto_rawDescData
}

var file_api_protos_v1_playlist_playlist_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_api_protos_v1_playlist_playlist_proto_goTypes = []interface{}{
	(*Playlist)(nil),         // 0: v1.playlist_proto.Playlist
	(*GetMyPlaylistReq)(nil), // 1: v1.playlist_proto.GetMyPlaylistReq
	(*GetMyPlaylistRes)(nil), // 2: v1.playlist_proto.GetMyPlaylistRes
}
var file_api_protos_v1_playlist_playlist_proto_depIdxs = []int32{
	0, // 0: v1.playlist_proto.GetMyPlaylistRes.myPlaylist:type_name -> v1.playlist_proto.Playlist
	1, // 1: v1.playlist_proto.PlaylistService.GetMyPlaylist:input_type -> v1.playlist_proto.GetMyPlaylistReq
	2, // 2: v1.playlist_proto.PlaylistService.GetMyPlaylist:output_type -> v1.playlist_proto.GetMyPlaylistRes
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_api_protos_v1_playlist_playlist_proto_init() }
func file_api_protos_v1_playlist_playlist_proto_init() {
	if File_api_protos_v1_playlist_playlist_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_protos_v1_playlist_playlist_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Playlist); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_protos_v1_playlist_playlist_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetMyPlaylistReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_protos_v1_playlist_playlist_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetMyPlaylistRes); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_protos_v1_playlist_playlist_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_protos_v1_playlist_playlist_proto_goTypes,
		DependencyIndexes: file_api_protos_v1_playlist_playlist_proto_depIdxs,
		MessageInfos:      file_api_protos_v1_playlist_playlist_proto_msgTypes,
	}.Build()
	File_api_protos_v1_playlist_playlist_proto = out.File
	file_api_protos_v1_playlist_playlist_proto_rawDesc = nil
	file_api_protos_v1_playlist_playlist_proto_goTypes = nil
	file_api_protos_v1_playlist_playlist_proto_depIdxs = nil
}