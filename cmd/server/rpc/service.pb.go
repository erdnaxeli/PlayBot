// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: cmd/server/service.proto

package rpc

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type TextMessage struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ChannelName   string                 `protobuf:"bytes,1,opt,name=channelName,proto3" json:"channelName,omitempty"`
	PersonName    string                 `protobuf:"bytes,2,opt,name=personName,proto3" json:"personName,omitempty"`
	Msg           string                 `protobuf:"bytes,3,opt,name=msg,proto3" json:"msg,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TextMessage) Reset() {
	*x = TextMessage{}
	mi := &file_cmd_server_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TextMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TextMessage) ProtoMessage() {}

func (x *TextMessage) ProtoReflect() protoreflect.Message {
	mi := &file_cmd_server_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TextMessage.ProtoReflect.Descriptor instead.
func (*TextMessage) Descriptor() ([]byte, []int) {
	return file_cmd_server_service_proto_rawDescGZIP(), []int{0}
}

func (x *TextMessage) GetChannelName() string {
	if x != nil {
		return x.ChannelName
	}
	return ""
}

func (x *TextMessage) GetPersonName() string {
	if x != nil {
		return x.PersonName
	}
	return ""
}

func (x *TextMessage) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

type Result struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Msg           []*IrcMessage          `protobuf:"bytes,1,rep,name=msg,proto3" json:"msg,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Result) Reset() {
	*x = Result{}
	mi := &file_cmd_server_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Result) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Result) ProtoMessage() {}

func (x *Result) ProtoReflect() protoreflect.Message {
	mi := &file_cmd_server_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Result.ProtoReflect.Descriptor instead.
func (*Result) Descriptor() ([]byte, []int) {
	return file_cmd_server_service_proto_rawDescGZIP(), []int{1}
}

func (x *Result) GetMsg() []*IrcMessage {
	if x != nil {
		return x.Msg
	}
	return nil
}

type IrcMessage struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Msg           string                 `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	To            string                 `protobuf:"bytes,2,opt,name=to,proto3" json:"to,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *IrcMessage) Reset() {
	*x = IrcMessage{}
	mi := &file_cmd_server_service_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *IrcMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IrcMessage) ProtoMessage() {}

func (x *IrcMessage) ProtoReflect() protoreflect.Message {
	mi := &file_cmd_server_service_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IrcMessage.ProtoReflect.Descriptor instead.
func (*IrcMessage) Descriptor() ([]byte, []int) {
	return file_cmd_server_service_proto_rawDescGZIP(), []int{2}
}

func (x *IrcMessage) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *IrcMessage) GetTo() string {
	if x != nil {
		return x.To
	}
	return ""
}

var File_cmd_server_service_proto protoreflect.FileDescriptor

const file_cmd_server_service_proto_rawDesc = "" +
	"\n" +
	"\x18cmd/server/service.proto\x12\x03rpc\"a\n" +
	"\vTextMessage\x12 \n" +
	"\vchannelName\x18\x01 \x01(\tR\vchannelName\x12\x1e\n" +
	"\n" +
	"personName\x18\x02 \x01(\tR\n" +
	"personName\x12\x10\n" +
	"\x03msg\x18\x03 \x01(\tR\x03msg\"+\n" +
	"\x06Result\x12!\n" +
	"\x03msg\x18\x01 \x03(\v2\x0f.rpc.IrcMessageR\x03msg\".\n" +
	"\n" +
	"IrcMessage\x12\x10\n" +
	"\x03msg\x18\x01 \x01(\tR\x03msg\x12\x0e\n" +
	"\x02to\x18\x02 \x01(\tR\x02to26\n" +
	"\n" +
	"PlaybotCli\x12(\n" +
	"\aExecute\x12\x10.rpc.TextMessage\x1a\v.rpc.ResultB*Z(github.com/erdnaxeli/playbot/cmd/cli/rpcb\x06proto3"

var (
	file_cmd_server_service_proto_rawDescOnce sync.Once
	file_cmd_server_service_proto_rawDescData []byte
)

func file_cmd_server_service_proto_rawDescGZIP() []byte {
	file_cmd_server_service_proto_rawDescOnce.Do(func() {
		file_cmd_server_service_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_cmd_server_service_proto_rawDesc), len(file_cmd_server_service_proto_rawDesc)))
	})
	return file_cmd_server_service_proto_rawDescData
}

var file_cmd_server_service_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_cmd_server_service_proto_goTypes = []any{
	(*TextMessage)(nil), // 0: rpc.TextMessage
	(*Result)(nil),      // 1: rpc.Result
	(*IrcMessage)(nil),  // 2: rpc.IrcMessage
}
var file_cmd_server_service_proto_depIdxs = []int32{
	2, // 0: rpc.Result.msg:type_name -> rpc.IrcMessage
	0, // 1: rpc.PlaybotCli.Execute:input_type -> rpc.TextMessage
	1, // 2: rpc.PlaybotCli.Execute:output_type -> rpc.Result
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_cmd_server_service_proto_init() }
func file_cmd_server_service_proto_init() {
	if File_cmd_server_service_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_cmd_server_service_proto_rawDesc), len(file_cmd_server_service_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_cmd_server_service_proto_goTypes,
		DependencyIndexes: file_cmd_server_service_proto_depIdxs,
		MessageInfos:      file_cmd_server_service_proto_msgTypes,
	}.Build()
	File_cmd_server_service_proto = out.File
	file_cmd_server_service_proto_goTypes = nil
	file_cmd_server_service_proto_depIdxs = nil
}
