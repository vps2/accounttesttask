// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: statistics.proto

package api

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_statistics_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_statistics_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_statistics_proto_rawDescGZIP(), []int{0}
}

var File_statistics_proto protoreflect.FileDescriptor

var file_statistics_proto_rawDesc = []byte{
	0x0a, 0x10, 0x73, 0x74, 0x61, 0x74, 0x69, 0x73, 0x74, 0x69, 0x63, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x03, 0x61, 0x70, 0x69, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x32, 0x36, 0x0a, 0x11, 0x53, 0x74, 0x61, 0x74, 0x69, 0x73, 0x74, 0x69, 0x63, 0x73, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x21, 0x0a, 0x05, 0x52, 0x65, 0x73, 0x65, 0x74, 0x12, 0x0a,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0a, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_statistics_proto_rawDescOnce sync.Once
	file_statistics_proto_rawDescData = file_statistics_proto_rawDesc
)

func file_statistics_proto_rawDescGZIP() []byte {
	file_statistics_proto_rawDescOnce.Do(func() {
		file_statistics_proto_rawDescData = protoimpl.X.CompressGZIP(file_statistics_proto_rawDescData)
	})
	return file_statistics_proto_rawDescData
}

var file_statistics_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_statistics_proto_goTypes = []interface{}{
	(*Empty)(nil), // 0: api.Empty
}
var file_statistics_proto_depIdxs = []int32{
	0, // 0: api.StatisticsService.Reset:input_type -> api.Empty
	0, // 1: api.StatisticsService.Reset:output_type -> api.Empty
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_statistics_proto_init() }
func file_statistics_proto_init() {
	if File_statistics_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_statistics_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
			RawDescriptor: file_statistics_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_statistics_proto_goTypes,
		DependencyIndexes: file_statistics_proto_depIdxs,
		MessageInfos:      file_statistics_proto_msgTypes,
	}.Build()
	File_statistics_proto = out.File
	file_statistics_proto_rawDesc = nil
	file_statistics_proto_goTypes = nil
	file_statistics_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// StatisticsServiceClient is the client API for StatisticsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type StatisticsServiceClient interface {
	Reset(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
}

type statisticsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStatisticsServiceClient(cc grpc.ClientConnInterface) StatisticsServiceClient {
	return &statisticsServiceClient{cc}
}

func (c *statisticsServiceClient) Reset(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/api.StatisticsService/Reset", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StatisticsServiceServer is the server API for StatisticsService service.
type StatisticsServiceServer interface {
	Reset(context.Context, *Empty) (*Empty, error)
}

// UnimplementedStatisticsServiceServer can be embedded to have forward compatible implementations.
type UnimplementedStatisticsServiceServer struct {
}

func (*UnimplementedStatisticsServiceServer) Reset(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Reset not implemented")
}

func RegisterStatisticsServiceServer(s *grpc.Server, srv StatisticsServiceServer) {
	s.RegisterService(&_StatisticsService_serviceDesc, srv)
}

func _StatisticsService_Reset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatisticsServiceServer).Reset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.StatisticsService/Reset",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatisticsServiceServer).Reset(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _StatisticsService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.StatisticsService",
	HandlerType: (*StatisticsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Reset",
			Handler:    _StatisticsService_Reset_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "statistics.proto",
}
