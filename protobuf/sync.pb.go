// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protobuf/sync.proto

package protobuf

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type Payload struct {
	Ack                  *Ack       `protobuf:"bytes,1,opt,name=ack,proto3" json:"ack,omitempty"`
	Offer                *Offer     `protobuf:"bytes,2,opt,name=offer,proto3" json:"offer,omitempty"`
	Request              *Request   `protobuf:"bytes,3,opt,name=request,proto3" json:"request,omitempty"`
	Messages             []*Message `protobuf:"bytes,4,rep,name=messages,proto3" json:"messages,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Payload) Reset()         { *m = Payload{} }
func (m *Payload) String() string { return proto.CompactTextString(m) }
func (*Payload) ProtoMessage()    {}
func (*Payload) Descriptor() ([]byte, []int) {
	return fileDescriptor_2dca527c092c79d7, []int{0}
}

func (m *Payload) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Payload.Unmarshal(m, b)
}
func (m *Payload) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Payload.Marshal(b, m, deterministic)
}
func (m *Payload) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Payload.Merge(m, src)
}
func (m *Payload) XXX_Size() int {
	return xxx_messageInfo_Payload.Size(m)
}
func (m *Payload) XXX_DiscardUnknown() {
	xxx_messageInfo_Payload.DiscardUnknown(m)
}

var xxx_messageInfo_Payload proto.InternalMessageInfo

func (m *Payload) GetAck() *Ack {
	if m != nil {
		return m.Ack
	}
	return nil
}

func (m *Payload) GetOffer() *Offer {
	if m != nil {
		return m.Offer
	}
	return nil
}

func (m *Payload) GetRequest() *Request {
	if m != nil {
		return m.Request
	}
	return nil
}

func (m *Payload) GetMessages() []*Message {
	if m != nil {
		return m.Messages
	}
	return nil
}

type Ack struct {
	Id                   [][]byte `protobuf:"bytes,1,rep,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Ack) Reset()         { *m = Ack{} }
func (m *Ack) String() string { return proto.CompactTextString(m) }
func (*Ack) ProtoMessage()    {}
func (*Ack) Descriptor() ([]byte, []int) {
	return fileDescriptor_2dca527c092c79d7, []int{1}
}

func (m *Ack) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Ack.Unmarshal(m, b)
}
func (m *Ack) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Ack.Marshal(b, m, deterministic)
}
func (m *Ack) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Ack.Merge(m, src)
}
func (m *Ack) XXX_Size() int {
	return xxx_messageInfo_Ack.Size(m)
}
func (m *Ack) XXX_DiscardUnknown() {
	xxx_messageInfo_Ack.DiscardUnknown(m)
}

var xxx_messageInfo_Ack proto.InternalMessageInfo

func (m *Ack) GetId() [][]byte {
	if m != nil {
		return m.Id
	}
	return nil
}

type Message struct {
	GroupId              []byte   `protobuf:"bytes,1,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
	Timestamp            int64    `protobuf:"varint,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Body                 []byte   `protobuf:"bytes,3,opt,name=body,proto3" json:"body,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_2dca527c092c79d7, []int{2}
}

func (m *Message) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message.Unmarshal(m, b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message.Marshal(b, m, deterministic)
}
func (m *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(m, src)
}
func (m *Message) XXX_Size() int {
	return xxx_messageInfo_Message.Size(m)
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetGroupId() []byte {
	if m != nil {
		return m.GroupId
	}
	return nil
}

func (m *Message) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *Message) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

type Offer struct {
	Id                   [][]byte `protobuf:"bytes,1,rep,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Offer) Reset()         { *m = Offer{} }
func (m *Offer) String() string { return proto.CompactTextString(m) }
func (*Offer) ProtoMessage()    {}
func (*Offer) Descriptor() ([]byte, []int) {
	return fileDescriptor_2dca527c092c79d7, []int{3}
}

func (m *Offer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Offer.Unmarshal(m, b)
}
func (m *Offer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Offer.Marshal(b, m, deterministic)
}
func (m *Offer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Offer.Merge(m, src)
}
func (m *Offer) XXX_Size() int {
	return xxx_messageInfo_Offer.Size(m)
}
func (m *Offer) XXX_DiscardUnknown() {
	xxx_messageInfo_Offer.DiscardUnknown(m)
}

var xxx_messageInfo_Offer proto.InternalMessageInfo

func (m *Offer) GetId() [][]byte {
	if m != nil {
		return m.Id
	}
	return nil
}

type Request struct {
	Id                   [][]byte `protobuf:"bytes,1,rep,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Request) Reset()         { *m = Request{} }
func (m *Request) String() string { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()    {}
func (*Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_2dca527c092c79d7, []int{4}
}

func (m *Request) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Request.Unmarshal(m, b)
}
func (m *Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Request.Marshal(b, m, deterministic)
}
func (m *Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Request.Merge(m, src)
}
func (m *Request) XXX_Size() int {
	return xxx_messageInfo_Request.Size(m)
}
func (m *Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Request proto.InternalMessageInfo

func (m *Request) GetId() [][]byte {
	if m != nil {
		return m.Id
	}
	return nil
}

func init() {
	proto.RegisterType((*Payload)(nil), "mvds.Payload")
	proto.RegisterType((*Ack)(nil), "mvds.Ack")
	proto.RegisterType((*Message)(nil), "mvds.Message")
	proto.RegisterType((*Offer)(nil), "mvds.Offer")
	proto.RegisterType((*Request)(nil), "mvds.Request")
}

func init() { proto.RegisterFile("protobuf/sync.proto", fileDescriptor_2dca527c092c79d7) }

var fileDescriptor_2dca527c092c79d7 = []byte{
	// 258 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x90, 0x41, 0x4f, 0xb4, 0x30,
	0x10, 0x86, 0x03, 0x65, 0x3f, 0xd8, 0x59, 0x3e, 0x0f, 0x63, 0x8c, 0xdd, 0xe8, 0x01, 0xb9, 0x88,
	0x17, 0x4c, 0xf4, 0x17, 0xac, 0x37, 0x0f, 0x46, 0xd3, 0x83, 0x07, 0x2f, 0xa6, 0xd0, 0xb2, 0x21,
	0x88, 0x45, 0x0a, 0x26, 0xfc, 0x18, 0xff, 0xab, 0x61, 0xba, 0xeb, 0x26, 0x7a, 0x9b, 0x79, 0x9f,
	0x27, 0xe9, 0xbc, 0x85, 0xe3, 0xae, 0x37, 0x83, 0x29, 0xc6, 0xea, 0xda, 0x4e, 0xef, 0x65, 0x4e,
	0x1b, 0x06, 0xed, 0xa7, 0xb2, 0xe9, 0x97, 0x07, 0xe1, 0x93, 0x9c, 0xde, 0x8c, 0x54, 0x78, 0x06,
	0x4c, 0x96, 0x0d, 0xf7, 0x12, 0x2f, 0x5b, 0xdd, 0x2c, 0xf3, 0x99, 0xe7, 0x9b, 0xb2, 0x11, 0x73,
	0x8a, 0x17, 0xb0, 0x30, 0x55, 0xa5, 0x7b, 0xee, 0x13, 0x5e, 0x39, 0xfc, 0x38, 0x47, 0xc2, 0x11,
	0xbc, 0x84, 0xb0, 0xd7, 0x1f, 0xa3, 0xb6, 0x03, 0x67, 0x24, 0xfd, 0x77, 0x92, 0x70, 0xa1, 0xd8,
	0x53, 0xbc, 0x82, 0xa8, 0xd5, 0xd6, 0xca, 0xad, 0xb6, 0x3c, 0x48, 0xd8, 0xc1, 0x7c, 0x70, 0xa9,
	0xf8, 0xc1, 0xe9, 0x09, 0xb0, 0x4d, 0xd9, 0xe0, 0x11, 0xf8, 0xb5, 0xe2, 0x5e, 0xc2, 0xb2, 0x58,
	0xf8, 0xb5, 0x4a, 0x9f, 0x21, 0xdc, 0xb9, 0xb8, 0x86, 0x68, 0xdb, 0x9b, 0xb1, 0x7b, 0x25, 0xc1,
	0xcb, 0x62, 0x11, 0xd2, 0x7e, 0xaf, 0xf0, 0x1c, 0x96, 0x43, 0xdd, 0x6a, 0x3b, 0xc8, 0xb6, 0xa3,
	0xbb, 0x99, 0x38, 0x04, 0x88, 0x10, 0x14, 0x46, 0x4d, 0x74, 0x6b, 0x2c, 0x68, 0x4e, 0x4f, 0x61,
	0x41, 0x95, 0xfe, 0x3c, 0xb8, 0x86, 0x70, 0x57, 0xe3, 0x37, 0xba, 0x83, 0x97, 0x68, 0xff, 0xbf,
	0xc5, 0x3f, 0x9a, 0x6e, 0xbf, 0x03, 0x00, 0x00, 0xff, 0xff, 0xd9, 0xf0, 0x22, 0x6e, 0x72, 0x01,
	0x00, 0x00,
}
