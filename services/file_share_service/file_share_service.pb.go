// Code generated by protoc-gen-go. DO NOT EDIT.
// source: services/file_share_service/file_share_service.proto

package chord

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

func init() {
	proto.RegisterFile("services/file_share_service/file_share_service.proto", fileDescriptor_06417f63079a9eef)
}

var fileDescriptor_06417f63079a9eef = []byte{
	// 133 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x32, 0x29, 0x4e, 0x2d, 0x2a,
	0xcb, 0x4c, 0x4e, 0x2d, 0xd6, 0x4f, 0xcb, 0xcc, 0x49, 0x8d, 0x2f, 0xce, 0x48, 0x2c, 0x4a, 0x8d,
	0x87, 0x8a, 0x61, 0x11, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x4d, 0xce, 0xc8, 0x2f,
	0x4a, 0x91, 0x32, 0x25, 0x52, 0x73, 0x6e, 0x6a, 0x71, 0x71, 0x62, 0x7a, 0x6a, 0x31, 0x44, 0xb7,
	0x91, 0x3b, 0x97, 0x80, 0x5b, 0x66, 0x4e, 0x6a, 0x30, 0x48, 0x2e, 0x18, 0xa2, 0x5a, 0xc8, 0x98,
	0x8b, 0x27, 0xa4, 0x28, 0x31, 0xaf, 0x38, 0x2d, 0xb5, 0x08, 0x24, 0x27, 0xc4, 0xaf, 0x07, 0xb6,
	0x42, 0x0f, 0xc4, 0xf1, 0xcc, 0x4b, 0xcb, 0x97, 0x12, 0x40, 0x12, 0x70, 0xce, 0x28, 0xcd, 0xcb,
	0x36, 0x60, 0x4c, 0x62, 0x03, 0x9b, 0x67, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0xae, 0xae, 0xdc,
	0x9d, 0xc5, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// FileShareServiceClient is the client API for FileShareService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type FileShareServiceClient interface {
	// It should return error if the File is not stored by FileInfo.id
	// or a stream with FileChunks if it is
	TransferFile(ctx context.Context, in *FileInfo, opts ...grpc.CallOption) (FileShareService_TransferFileClient, error)
}

type fileShareServiceClient struct {
	cc *grpc.ClientConn
}

func NewFileShareServiceClient(cc *grpc.ClientConn) FileShareServiceClient {
	return &fileShareServiceClient{cc}
}

func (c *fileShareServiceClient) TransferFile(ctx context.Context, in *FileInfo, opts ...grpc.CallOption) (FileShareService_TransferFileClient, error) {
	stream, err := c.cc.NewStream(ctx, &_FileShareService_serviceDesc.Streams[0], "/chord.FileShareService/TransferFile", opts...)
	if err != nil {
		return nil, err
	}
	x := &fileShareServiceTransferFileClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type FileShareService_TransferFileClient interface {
	Recv() (*FileChunk, error)
	grpc.ClientStream
}

type fileShareServiceTransferFileClient struct {
	grpc.ClientStream
}

func (x *fileShareServiceTransferFileClient) Recv() (*FileChunk, error) {
	m := new(FileChunk)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// FileShareServiceServer is the server API for FileShareService service.
type FileShareServiceServer interface {
	// It should return error if the File is not stored by FileInfo.id
	// or a stream with FileChunks if it is
	TransferFile(*FileInfo, FileShareService_TransferFileServer) error
}

// UnimplementedFileShareServiceServer can be embedded to have forward compatible implementations.
type UnimplementedFileShareServiceServer struct {
}

func (*UnimplementedFileShareServiceServer) TransferFile(req *FileInfo, srv FileShareService_TransferFileServer) error {
	return status.Errorf(codes.Unimplemented, "method TransferFile not implemented")
}

func RegisterFileShareServiceServer(s *grpc.Server, srv FileShareServiceServer) {
	s.RegisterService(&_FileShareService_serviceDesc, srv)
}

func _FileShareService_TransferFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(FileInfo)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(FileShareServiceServer).TransferFile(m, &fileShareServiceTransferFileServer{stream})
}

type FileShareService_TransferFileServer interface {
	Send(*FileChunk) error
	grpc.ServerStream
}

type fileShareServiceTransferFileServer struct {
	grpc.ServerStream
}

func (x *fileShareServiceTransferFileServer) Send(m *FileChunk) error {
	return x.ServerStream.SendMsg(m)
}

var _FileShareService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "chord.FileShareService",
	HandlerType: (*FileShareServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "TransferFile",
			Handler:       _FileShareService_TransferFile_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "services/file_share_service/file_share_service.proto",
}
