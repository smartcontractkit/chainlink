package serialization

import (
	reflect "reflect"
	sync "sync"

	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

const _ = proto.ProtoPackageIsVersion4

type TelemetryMessageReceived struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConfigDigest []byte          `protobuf:"bytes,1,opt,name=configDigest,proto3" json:"configDigest,omitempty"`
	Msg          *MessageWrapper `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Sender       uint32          `protobuf:"varint,3,opt,name=sender,proto3" json:"sender,omitempty"`
}

func (x *TelemetryMessageReceived) Reset() {
	*x = TelemetryMessageReceived{}
	if protoimpl.UnsafeEnabled {
		mi := &file_telemetry_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TelemetryMessageReceived) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TelemetryMessageReceived) ProtoMessage() {}

func (x *TelemetryMessageReceived) ProtoReflect() protoreflect.Message {
	mi := &file_telemetry_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*TelemetryMessageReceived) Descriptor() ([]byte, []int) {
	return file_telemetry_proto_rawDescGZIP(), []int{0}
}

func (x *TelemetryMessageReceived) GetConfigDigest() []byte {
	if x != nil {
		return x.ConfigDigest
	}
	return nil
}

func (x *TelemetryMessageReceived) GetMsg() *MessageWrapper {
	if x != nil {
		return x.Msg
	}
	return nil
}

func (x *TelemetryMessageReceived) GetSender() uint32 {
	if x != nil {
		return x.Sender
	}
	return 0
}

type TelemetryMessageSent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConfigDigest []byte          `protobuf:"bytes,1,opt,name=configDigest,proto3" json:"configDigest,omitempty"`
	Msg          *MessageWrapper `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Receiver     uint32          `protobuf:"varint,3,opt,name=receiver,proto3" json:"receiver,omitempty"`
}

func (x *TelemetryMessageSent) Reset() {
	*x = TelemetryMessageSent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_telemetry_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TelemetryMessageSent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TelemetryMessageSent) ProtoMessage() {}

func (x *TelemetryMessageSent) ProtoReflect() protoreflect.Message {
	mi := &file_telemetry_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*TelemetryMessageSent) Descriptor() ([]byte, []int) {
	return file_telemetry_proto_rawDescGZIP(), []int{1}
}

func (x *TelemetryMessageSent) GetConfigDigest() []byte {
	if x != nil {
		return x.ConfigDigest
	}
	return nil
}

func (x *TelemetryMessageSent) GetMsg() *MessageWrapper {
	if x != nil {
		return x.Msg
	}
	return nil
}

func (x *TelemetryMessageSent) GetReceiver() uint32 {
	if x != nil {
		return x.Receiver
	}
	return 0
}

type TelemetryAssertionViolation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	E isTelemetryAssertionViolation_E `protobuf_oneof:"e"`
}

func (x *TelemetryAssertionViolation) Reset() {
	*x = TelemetryAssertionViolation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_telemetry_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TelemetryAssertionViolation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TelemetryAssertionViolation) ProtoMessage() {}

func (x *TelemetryAssertionViolation) ProtoReflect() protoreflect.Message {
	mi := &file_telemetry_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*TelemetryAssertionViolation) Descriptor() ([]byte, []int) {
	return file_telemetry_proto_rawDescGZIP(), []int{2}
}

func (m *TelemetryAssertionViolation) GetE() isTelemetryAssertionViolation_E {
	if m != nil {
		return m.E
	}
	return nil
}

func (x *TelemetryAssertionViolation) GetInvalidSignature() *TelemetryAssertionViolationInvalidSignature {
	if x, ok := x.GetE().(*TelemetryAssertionViolation_InvalidSignature); ok {
		return x.InvalidSignature
	}
	return nil
}

type isTelemetryAssertionViolation_E interface {
	isTelemetryAssertionViolation_E()
}

type TelemetryAssertionViolation_InvalidSignature struct {
	InvalidSignature *TelemetryAssertionViolationInvalidSignature `protobuf:"bytes,1,opt,name=invalidSignature,proto3,oneof"`
}

func (*TelemetryAssertionViolation_InvalidSignature) isTelemetryAssertionViolation_E() {}

type TelemetryAssertionViolationInvalidSignature struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConfigDigest []byte          `protobuf:"bytes,1,opt,name=configDigest,proto3" json:"configDigest,omitempty"`
	Msg          *MessageWrapper `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Sender       uint32          `protobuf:"varint,3,opt,name=sender,proto3" json:"sender,omitempty"`
}

func (x *TelemetryAssertionViolationInvalidSignature) Reset() {
	*x = TelemetryAssertionViolationInvalidSignature{}
	if protoimpl.UnsafeEnabled {
		mi := &file_telemetry_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TelemetryAssertionViolationInvalidSignature) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TelemetryAssertionViolationInvalidSignature) ProtoMessage() {}

func (x *TelemetryAssertionViolationInvalidSignature) ProtoReflect() protoreflect.Message {
	mi := &file_telemetry_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*TelemetryAssertionViolationInvalidSignature) Descriptor() ([]byte, []int) {
	return file_telemetry_proto_rawDescGZIP(), []int{3}
}

func (x *TelemetryAssertionViolationInvalidSignature) GetConfigDigest() []byte {
	if x != nil {
		return x.ConfigDigest
	}
	return nil
}

func (x *TelemetryAssertionViolationInvalidSignature) GetMsg() *MessageWrapper {
	if x != nil {
		return x.Msg
	}
	return nil
}

func (x *TelemetryAssertionViolationInvalidSignature) GetSender() uint32 {
	if x != nil {
		return x.Sender
	}
	return 0
}

type TelemetryStateUpdate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConfigDigest []byte `protobuf:"bytes,1,opt,name=configDigest,proto3" json:"configDigest,omitempty"`
	Epoch        uint64 `protobuf:"varint,2,opt,name=epoch,proto3" json:"epoch,omitempty"`
	Round        uint64 `protobuf:"varint,3,opt,name=round,proto3" json:"round,omitempty"`
	Time         uint64 `protobuf:"varint,4,opt,name=time,proto3" json:"time,omitempty"`
}

func (x *TelemetryStateUpdate) Reset() {
	*x = TelemetryStateUpdate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_telemetry_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TelemetryStateUpdate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TelemetryStateUpdate) ProtoMessage() {}

func (x *TelemetryStateUpdate) ProtoReflect() protoreflect.Message {
	mi := &file_telemetry_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*TelemetryStateUpdate) Descriptor() ([]byte, []int) {
	return file_telemetry_proto_rawDescGZIP(), []int{4}
}

func (x *TelemetryStateUpdate) GetConfigDigest() []byte {
	if x != nil {
		return x.ConfigDigest
	}
	return nil
}

func (x *TelemetryStateUpdate) GetEpoch() uint64 {
	if x != nil {
		return x.Epoch
	}
	return 0
}

func (x *TelemetryStateUpdate) GetRound() uint64 {
	if x != nil {
		return x.Round
	}
	return 0
}

func (x *TelemetryStateUpdate) GetTime() uint64 {
	if x != nil {
		return x.Time
	}
	return 0
}

var File_telemetry_proto protoreflect.FileDescriptor

var file_telemetry_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x0d, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x1a, 0x0e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x87, 0x01, 0x0a, 0x18, 0x54, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x64, 0x12, 0x22, 0x0a,
	0x0c, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x44, 0x69, 0x67, 0x65, 0x73, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x0c, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x44, 0x69, 0x67, 0x65, 0x73,
	0x74, 0x12, 0x2f, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d,
	0x2e, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x52, 0x03, 0x6d,
	0x73, 0x67, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x22, 0x87, 0x01, 0x0a, 0x14, 0x54,
	0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x53,
	0x65, 0x6e, 0x74, 0x12, 0x22, 0x0a, 0x0c, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x44, 0x69, 0x67,
	0x65, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x44, 0x69, 0x67, 0x65, 0x73, 0x74, 0x12, 0x2f, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x57, 0x72, 0x61, 0x70,
	0x70, 0x65, 0x72, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x63, 0x65,
	0x69, 0x76, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x72, 0x65, 0x63, 0x65,
	0x69, 0x76, 0x65, 0x72, 0x22, 0x8c, 0x01, 0x0a, 0x1b, 0x54, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74,
	0x72, 0x79, 0x41, 0x73, 0x73, 0x65, 0x72, 0x74, 0x69, 0x6f, 0x6e, 0x56, 0x69, 0x6f, 0x6c, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x68, 0x0a, 0x10, 0x69, 0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x53,
	0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x3a,
	0x2e, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x54,
	0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x41, 0x73, 0x73, 0x65, 0x72, 0x74, 0x69, 0x6f,
	0x6e, 0x56, 0x69, 0x6f, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x48, 0x00, 0x52, 0x10, 0x69, 0x6e,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x42, 0x03,
	0x0a, 0x01, 0x65, 0x22, 0x9a, 0x01, 0x0a, 0x2b, 0x54, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72,
	0x79, 0x41, 0x73, 0x73, 0x65, 0x72, 0x74, 0x69, 0x6f, 0x6e, 0x56, 0x69, 0x6f, 0x6c, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74,
	0x75, 0x72, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x44, 0x69, 0x67,
	0x65, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x44, 0x69, 0x67, 0x65, 0x73, 0x74, 0x12, 0x2f, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x57, 0x72, 0x61, 0x70,
	0x70, 0x65, 0x72, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x6e, 0x64,
	0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72,
	0x22, 0x7a, 0x0a, 0x14, 0x54, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x53, 0x74, 0x61,
	0x74, 0x65, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x44, 0x69, 0x67, 0x65, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x44, 0x69, 0x67, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05,
	0x65, 0x70, 0x6f, 0x63, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x65, 0x70, 0x6f,
	0x63, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x42, 0x11, 0x5a, 0x0f,
	0x2e, 0x3b, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_telemetry_proto_rawDescOnce sync.Once
	file_telemetry_proto_rawDescData = file_telemetry_proto_rawDesc
)

func file_telemetry_proto_rawDescGZIP() []byte {
	file_telemetry_proto_rawDescOnce.Do(func() {
		file_telemetry_proto_rawDescData = protoimpl.X.CompressGZIP(file_telemetry_proto_rawDescData)
	})
	return file_telemetry_proto_rawDescData
}

var file_telemetry_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_telemetry_proto_goTypes = []interface{}{
	(*TelemetryMessageReceived)(nil), (*TelemetryMessageSent)(nil), (*TelemetryAssertionViolation)(nil), (*TelemetryAssertionViolationInvalidSignature)(nil), (*TelemetryStateUpdate)(nil), (*MessageWrapper)(nil)}
var file_telemetry_proto_depIdxs = []int32{
	5, 5, 3, 5, 4, 4, 4, 4, 0}

func init() { file_telemetry_proto_init() }
func file_telemetry_proto_init() {
	if File_telemetry_proto != nil {
		return
	}
	file_messages_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_telemetry_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TelemetryMessageReceived); i {
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
		file_telemetry_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TelemetryMessageSent); i {
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
		file_telemetry_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TelemetryAssertionViolation); i {
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
		file_telemetry_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TelemetryAssertionViolationInvalidSignature); i {
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
		file_telemetry_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TelemetryStateUpdate); i {
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
	file_telemetry_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*TelemetryAssertionViolation_InvalidSignature)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_telemetry_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_telemetry_proto_goTypes,
		DependencyIndexes: file_telemetry_proto_depIdxs,
		MessageInfos:      file_telemetry_proto_msgTypes,
	}.Build()
	File_telemetry_proto = out.File
	file_telemetry_proto_rawDesc = nil
	file_telemetry_proto_goTypes = nil
	file_telemetry_proto_depIdxs = nil
}
