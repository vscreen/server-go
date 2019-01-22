// Code generated by protoc-gen-go. DO NOT EDIT.
// source: vscreen.proto

package server

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
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

type StatusCode int32

const (
	StatusCode_OK               StatusCode = 0
	StatusCode_OPERATION_FAILED StatusCode = 1
)

var StatusCode_name = map[int32]string{
	0: "OK",
	1: "OPERATION_FAILED",
}

var StatusCode_value = map[string]int32{
	"OK":               0,
	"OPERATION_FAILED": 1,
}

func (x StatusCode) String() string {
	return proto.EnumName(StatusCode_name, int32(x))
}

func (StatusCode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_fc3c777aeb3c5845, []int{0}
}

type Info_State int32

const (
	Info_PLAYING Info_State = 0
	Info_PAUSED  Info_State = 1
	Info_STOPPED Info_State = 2
)

var Info_State_name = map[int32]string{
	0: "PLAYING",
	1: "PAUSED",
	2: "STOPPED",
}

var Info_State_value = map[string]int32{
	"PLAYING": 0,
	"PAUSED":  1,
	"STOPPED": 2,
}

func (x Info_State) String() string {
	return proto.EnumName(Info_State_name, int32(x))
}

func (Info_State) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_fc3c777aeb3c5845, []int{6, 0}
}

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_fc3c777aeb3c5845, []int{0}
}

func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (m *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(m, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

type Status struct {
	Code                 StatusCode `protobuf:"varint,1,opt,name=code,proto3,enum=StatusCode" json:"code,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Status) Reset()         { *m = Status{} }
func (m *Status) String() string { return proto.CompactTextString(m) }
func (*Status) ProtoMessage()    {}
func (*Status) Descriptor() ([]byte, []int) {
	return fileDescriptor_fc3c777aeb3c5845, []int{1}
}

func (m *Status) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Status.Unmarshal(m, b)
}
func (m *Status) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Status.Marshal(b, m, deterministic)
}
func (m *Status) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Status.Merge(m, src)
}
func (m *Status) XXX_Size() int {
	return xxx_messageInfo_Status.Size(m)
}
func (m *Status) XXX_DiscardUnknown() {
	xxx_messageInfo_Status.DiscardUnknown(m)
}

var xxx_messageInfo_Status proto.InternalMessageInfo

func (m *Status) GetCode() StatusCode {
	if m != nil {
		return m.Code
	}
	return StatusCode_OK
}

type Credential struct {
	Password             string   `protobuf:"bytes,1,opt,name=password,proto3" json:"password,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Credential) Reset()         { *m = Credential{} }
func (m *Credential) String() string { return proto.CompactTextString(m) }
func (*Credential) ProtoMessage()    {}
func (*Credential) Descriptor() ([]byte, []int) {
	return fileDescriptor_fc3c777aeb3c5845, []int{2}
}

func (m *Credential) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Credential.Unmarshal(m, b)
}
func (m *Credential) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Credential.Marshal(b, m, deterministic)
}
func (m *Credential) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Credential.Merge(m, src)
}
func (m *Credential) XXX_Size() int {
	return xxx_messageInfo_Credential.Size(m)
}
func (m *Credential) XXX_DiscardUnknown() {
	xxx_messageInfo_Credential.DiscardUnknown(m)
}

var xxx_messageInfo_Credential proto.InternalMessageInfo

func (m *Credential) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type Source struct {
	Url                  string   `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Source) Reset()         { *m = Source{} }
func (m *Source) String() string { return proto.CompactTextString(m) }
func (*Source) ProtoMessage()    {}
func (*Source) Descriptor() ([]byte, []int) {
	return fileDescriptor_fc3c777aeb3c5845, []int{3}
}

func (m *Source) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Source.Unmarshal(m, b)
}
func (m *Source) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Source.Marshal(b, m, deterministic)
}
func (m *Source) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Source.Merge(m, src)
}
func (m *Source) XXX_Size() int {
	return xxx_messageInfo_Source.Size(m)
}
func (m *Source) XXX_DiscardUnknown() {
	xxx_messageInfo_Source.DiscardUnknown(m)
}

var xxx_messageInfo_Source proto.InternalMessageInfo

func (m *Source) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

type Position struct {
	Value                float32  `protobuf:"fixed32,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Position) Reset()         { *m = Position{} }
func (m *Position) String() string { return proto.CompactTextString(m) }
func (*Position) ProtoMessage()    {}
func (*Position) Descriptor() ([]byte, []int) {
	return fileDescriptor_fc3c777aeb3c5845, []int{4}
}

func (m *Position) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Position.Unmarshal(m, b)
}
func (m *Position) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Position.Marshal(b, m, deterministic)
}
func (m *Position) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Position.Merge(m, src)
}
func (m *Position) XXX_Size() int {
	return xxx_messageInfo_Position.Size(m)
}
func (m *Position) XXX_DiscardUnknown() {
	xxx_messageInfo_Position.DiscardUnknown(m)
}

var xxx_messageInfo_Position proto.InternalMessageInfo

func (m *Position) GetValue() float32 {
	if m != nil {
		return m.Value
	}
	return 0
}

type User struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *User) Reset()         { *m = User{} }
func (m *User) String() string { return proto.CompactTextString(m) }
func (*User) ProtoMessage()    {}
func (*User) Descriptor() ([]byte, []int) {
	return fileDescriptor_fc3c777aeb3c5845, []int{5}
}

func (m *User) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_User.Unmarshal(m, b)
}
func (m *User) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_User.Marshal(b, m, deterministic)
}
func (m *User) XXX_Merge(src proto.Message) {
	xxx_messageInfo_User.Merge(m, src)
}
func (m *User) XXX_Size() int {
	return xxx_messageInfo_User.Size(m)
}
func (m *User) XXX_DiscardUnknown() {
	xxx_messageInfo_User.DiscardUnknown(m)
}

var xxx_messageInfo_User proto.InternalMessageInfo

func (m *User) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type Info struct {
	Title                string     `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	ThumbnailURL         string     `protobuf:"bytes,2,opt,name=thumbnailURL,proto3" json:"thumbnailURL,omitempty"`
	Volume               float64    `protobuf:"fixed64,3,opt,name=volume,proto3" json:"volume,omitempty"`
	Position             float64    `protobuf:"fixed64,4,opt,name=position,proto3" json:"position,omitempty"`
	State                Info_State `protobuf:"varint,5,opt,name=state,proto3,enum=Info_State" json:"state,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Info) Reset()         { *m = Info{} }
func (m *Info) String() string { return proto.CompactTextString(m) }
func (*Info) ProtoMessage()    {}
func (*Info) Descriptor() ([]byte, []int) {
	return fileDescriptor_fc3c777aeb3c5845, []int{6}
}

func (m *Info) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Info.Unmarshal(m, b)
}
func (m *Info) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Info.Marshal(b, m, deterministic)
}
func (m *Info) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Info.Merge(m, src)
}
func (m *Info) XXX_Size() int {
	return xxx_messageInfo_Info.Size(m)
}
func (m *Info) XXX_DiscardUnknown() {
	xxx_messageInfo_Info.DiscardUnknown(m)
}

var xxx_messageInfo_Info proto.InternalMessageInfo

func (m *Info) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Info) GetThumbnailURL() string {
	if m != nil {
		return m.ThumbnailURL
	}
	return ""
}

func (m *Info) GetVolume() float64 {
	if m != nil {
		return m.Volume
	}
	return 0
}

func (m *Info) GetPosition() float64 {
	if m != nil {
		return m.Position
	}
	return 0
}

func (m *Info) GetState() Info_State {
	if m != nil {
		return m.State
	}
	return Info_PLAYING
}

func init() {
	proto.RegisterEnum("StatusCode", StatusCode_name, StatusCode_value)
	proto.RegisterEnum("Info_State", Info_State_name, Info_State_value)
	proto.RegisterType((*Empty)(nil), "Empty")
	proto.RegisterType((*Status)(nil), "Status")
	proto.RegisterType((*Credential)(nil), "Credential")
	proto.RegisterType((*Source)(nil), "Source")
	proto.RegisterType((*Position)(nil), "Position")
	proto.RegisterType((*User)(nil), "User")
	proto.RegisterType((*Info)(nil), "Info")
}

func init() { proto.RegisterFile("vscreen.proto", fileDescriptor_fc3c777aeb3c5845) }

var fileDescriptor_fc3c777aeb3c5845 = []byte{
	// 451 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0xdd, 0x6e, 0xd3, 0x40,
	0x10, 0x85, 0x63, 0xd7, 0x3f, 0xc9, 0x04, 0x2a, 0x6b, 0x54, 0x55, 0xae, 0x05, 0x34, 0xec, 0x55,
	0xa8, 0x84, 0x85, 0xca, 0x13, 0x98, 0x36, 0xa0, 0x88, 0x28, 0xb1, 0xec, 0x06, 0x09, 0x6e, 0x90,
	0x13, 0x0f, 0xaa, 0x85, 0xe3, 0x8d, 0x76, 0xd7, 0x81, 0x3e, 0x05, 0x6f, 0xc5, 0x73, 0x21, 0xaf,
	0xdd, 0x26, 0xa0, 0x72, 0xe7, 0xb3, 0xdf, 0x1c, 0x7b, 0x7c, 0xce, 0xc2, 0xd3, 0x9d, 0x5c, 0x0b,
	0xa2, 0x2a, 0xdc, 0x0a, 0xae, 0x38, 0x73, 0xc1, 0x9e, 0x6c, 0xb6, 0xea, 0x8e, 0xbd, 0x02, 0x27,
	0x55, 0x99, 0xaa, 0x25, 0x9e, 0x83, 0xb5, 0xe6, 0x39, 0xf9, 0xc6, 0xc8, 0x18, 0x1f, 0x5f, 0x0e,
	0xc3, 0xf6, 0xf8, 0x8a, 0xe7, 0x94, 0x68, 0xc0, 0xc6, 0x00, 0x57, 0x82, 0x72, 0xaa, 0x54, 0x91,
	0x95, 0x18, 0x40, 0x7f, 0x9b, 0x49, 0xf9, 0x83, 0x8b, 0x5c, 0x5b, 0x06, 0xc9, 0x83, 0x66, 0x01,
	0x38, 0x29, 0xaf, 0xc5, 0x9a, 0xd0, 0x83, 0xa3, 0x5a, 0x94, 0xdd, 0x40, 0xf3, 0xc8, 0x46, 0xd0,
	0x8f, 0xb9, 0x2c, 0x54, 0xc1, 0x2b, 0x3c, 0x01, 0x7b, 0x97, 0x95, 0x75, 0xfb, 0x4d, 0x33, 0x69,
	0x05, 0x3b, 0x05, 0x6b, 0x29, 0x49, 0xe0, 0x31, 0x98, 0xc5, 0xfd, 0xbb, 0xcd, 0x22, 0x67, 0xbf,
	0x0d, 0xb0, 0xa6, 0xd5, 0x37, 0xde, 0xd8, 0x54, 0xa1, 0x4a, 0xea, 0x58, 0x2b, 0x90, 0xc1, 0x13,
	0x75, 0x5b, 0x6f, 0x56, 0x55, 0x56, 0x94, 0xcb, 0x64, 0xe6, 0x9b, 0x1a, 0xfe, 0x75, 0x86, 0xa7,
	0xe0, 0xec, 0x78, 0x59, 0x6f, 0xc8, 0x3f, 0x1a, 0x19, 0x63, 0x23, 0xe9, 0x94, 0xfe, 0x99, 0x6e,
	0x29, 0xdf, 0xd2, 0xe4, 0x41, 0xe3, 0x4b, 0xb0, 0xa5, 0xca, 0x14, 0xf9, 0x76, 0x17, 0x4c, 0xb3,
	0x83, 0x4e, 0x87, 0x92, 0x96, 0xb0, 0xd7, 0x60, 0x6b, 0x8d, 0x43, 0x70, 0xe3, 0x59, 0xf4, 0x79,
	0x3a, 0xff, 0xe0, 0xf5, 0x10, 0xc0, 0x89, 0xa3, 0x65, 0x3a, 0xb9, 0xf6, 0x8c, 0x06, 0xa4, 0x37,
	0x8b, 0x38, 0x9e, 0x5c, 0x7b, 0xe6, 0xc5, 0x05, 0xc0, 0x3e, 0x5c, 0x74, 0xc0, 0x5c, 0x7c, 0xf4,
	0x7a, 0x78, 0x02, 0xde, 0x22, 0x9e, 0x24, 0xd1, 0xcd, 0x74, 0x31, 0xff, 0xfa, 0x3e, 0x9a, 0xce,
	0x1a, 0xe3, 0xe5, 0x2f, 0x13, 0xdc, 0x4f, 0xa9, 0xae, 0x0e, 0x5f, 0x80, 0x15, 0xd5, 0xea, 0x16,
	0x87, 0xe1, 0xbe, 0x87, 0xc0, 0xed, 0x8a, 0x62, 0x3d, 0x3c, 0x03, 0x2b, 0x2e, 0xb3, 0x3b, 0x74,
	0x42, 0xdd, 0xed, 0x21, 0x0a, 0xc0, 0x8e, 0xb3, 0x5a, 0xd2, 0x63, 0xec, 0x0c, 0xac, 0x54, 0xf1,
	0xed, 0x7f, 0xd0, 0x9c, 0x7e, 0xaa, 0xc7, 0xd1, 0x51, 0x94, 0xe7, 0xe8, 0x86, 0x6d, 0xd3, 0x87,
	0xe8, 0x19, 0x58, 0x29, 0xd1, 0x77, 0x1c, 0x84, 0xf7, 0x4d, 0x1f, 0xd2, 0xe7, 0x30, 0x48, 0xeb,
	0x95, 0x5c, 0x8b, 0x62, 0x45, 0x68, 0x87, 0x4d, 0xd5, 0x81, 0xad, 0x43, 0x65, 0xbd, 0x37, 0x06,
	0x9e, 0xc3, 0x70, 0x59, 0xc9, 0x7f, 0x07, 0xf6, 0xfe, 0x77, 0xfd, 0x2f, 0x8e, 0x24, 0xb1, 0x23,
	0xb1, 0x72, 0xf4, 0x5d, 0x7e, 0xfb, 0x27, 0x00, 0x00, 0xff, 0xff, 0x2d, 0x0c, 0x57, 0xf1, 0xdc,
	0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// VScreenClient is the client API for VScreen service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type VScreenClient interface {
	Auth(ctx context.Context, in *Credential, opts ...grpc.CallOption) (*Status, error)
	Play(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Status, error)
	Pause(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Status, error)
	Stop(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Status, error)
	Next(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Status, error)
	Add(ctx context.Context, in *Source, opts ...grpc.CallOption) (*Status, error)
	Seek(ctx context.Context, in *Position, opts ...grpc.CallOption) (*Status, error)
	Subscribe(ctx context.Context, in *User, opts ...grpc.CallOption) (VScreen_SubscribeClient, error)
	Unsubscribe(ctx context.Context, in *User, opts ...grpc.CallOption) (*Status, error)
}

type vScreenClient struct {
	cc *grpc.ClientConn
}

func NewVScreenClient(cc *grpc.ClientConn) VScreenClient {
	return &vScreenClient{cc}
}

func (c *vScreenClient) Auth(ctx context.Context, in *Credential, opts ...grpc.CallOption) (*Status, error) {
	out := new(Status)
	err := c.cc.Invoke(ctx, "/VScreen/Auth", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vScreenClient) Play(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Status, error) {
	out := new(Status)
	err := c.cc.Invoke(ctx, "/VScreen/Play", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vScreenClient) Pause(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Status, error) {
	out := new(Status)
	err := c.cc.Invoke(ctx, "/VScreen/Pause", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vScreenClient) Stop(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Status, error) {
	out := new(Status)
	err := c.cc.Invoke(ctx, "/VScreen/Stop", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vScreenClient) Next(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Status, error) {
	out := new(Status)
	err := c.cc.Invoke(ctx, "/VScreen/Next", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vScreenClient) Add(ctx context.Context, in *Source, opts ...grpc.CallOption) (*Status, error) {
	out := new(Status)
	err := c.cc.Invoke(ctx, "/VScreen/Add", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vScreenClient) Seek(ctx context.Context, in *Position, opts ...grpc.CallOption) (*Status, error) {
	out := new(Status)
	err := c.cc.Invoke(ctx, "/VScreen/Seek", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vScreenClient) Subscribe(ctx context.Context, in *User, opts ...grpc.CallOption) (VScreen_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &_VScreen_serviceDesc.Streams[0], "/VScreen/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &vScreenSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type VScreen_SubscribeClient interface {
	Recv() (*Info, error)
	grpc.ClientStream
}

type vScreenSubscribeClient struct {
	grpc.ClientStream
}

func (x *vScreenSubscribeClient) Recv() (*Info, error) {
	m := new(Info)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *vScreenClient) Unsubscribe(ctx context.Context, in *User, opts ...grpc.CallOption) (*Status, error) {
	out := new(Status)
	err := c.cc.Invoke(ctx, "/VScreen/Unsubscribe", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VScreenServer is the server API for VScreen service.
type VScreenServer interface {
	Auth(context.Context, *Credential) (*Status, error)
	Play(context.Context, *Empty) (*Status, error)
	Pause(context.Context, *Empty) (*Status, error)
	Stop(context.Context, *Empty) (*Status, error)
	Next(context.Context, *Empty) (*Status, error)
	Add(context.Context, *Source) (*Status, error)
	Seek(context.Context, *Position) (*Status, error)
	Subscribe(*User, VScreen_SubscribeServer) error
	Unsubscribe(context.Context, *User) (*Status, error)
}

func RegisterVScreenServer(s *grpc.Server, srv VScreenServer) {
	s.RegisterService(&_VScreen_serviceDesc, srv)
}

func _VScreen_Auth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Credential)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VScreenServer).Auth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/VScreen/Auth",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VScreenServer).Auth(ctx, req.(*Credential))
	}
	return interceptor(ctx, in, info, handler)
}

func _VScreen_Play_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VScreenServer).Play(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/VScreen/Play",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VScreenServer).Play(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _VScreen_Pause_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VScreenServer).Pause(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/VScreen/Pause",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VScreenServer).Pause(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _VScreen_Stop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VScreenServer).Stop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/VScreen/Stop",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VScreenServer).Stop(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _VScreen_Next_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VScreenServer).Next(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/VScreen/Next",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VScreenServer).Next(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _VScreen_Add_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Source)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VScreenServer).Add(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/VScreen/Add",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VScreenServer).Add(ctx, req.(*Source))
	}
	return interceptor(ctx, in, info, handler)
}

func _VScreen_Seek_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Position)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VScreenServer).Seek(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/VScreen/Seek",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VScreenServer).Seek(ctx, req.(*Position))
	}
	return interceptor(ctx, in, info, handler)
}

func _VScreen_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(User)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(VScreenServer).Subscribe(m, &vScreenSubscribeServer{stream})
}

type VScreen_SubscribeServer interface {
	Send(*Info) error
	grpc.ServerStream
}

type vScreenSubscribeServer struct {
	grpc.ServerStream
}

func (x *vScreenSubscribeServer) Send(m *Info) error {
	return x.ServerStream.SendMsg(m)
}

func _VScreen_Unsubscribe_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(User)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VScreenServer).Unsubscribe(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/VScreen/Unsubscribe",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VScreenServer).Unsubscribe(ctx, req.(*User))
	}
	return interceptor(ctx, in, info, handler)
}

var _VScreen_serviceDesc = grpc.ServiceDesc{
	ServiceName: "VScreen",
	HandlerType: (*VScreenServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Auth",
			Handler:    _VScreen_Auth_Handler,
		},
		{
			MethodName: "Play",
			Handler:    _VScreen_Play_Handler,
		},
		{
			MethodName: "Pause",
			Handler:    _VScreen_Pause_Handler,
		},
		{
			MethodName: "Stop",
			Handler:    _VScreen_Stop_Handler,
		},
		{
			MethodName: "Next",
			Handler:    _VScreen_Next_Handler,
		},
		{
			MethodName: "Add",
			Handler:    _VScreen_Add_Handler,
		},
		{
			MethodName: "Seek",
			Handler:    _VScreen_Seek_Handler,
		},
		{
			MethodName: "Unsubscribe",
			Handler:    _VScreen_Unsubscribe_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _VScreen_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "vscreen.proto",
}