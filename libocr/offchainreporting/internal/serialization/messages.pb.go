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

type MessageNewEpoch struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Epoch uint64 `protobuf:"varint,1,opt,name=epoch,proto3" json:"epoch,omitempty"`
}

func (x *MessageNewEpoch) Reset() {
	*x = MessageNewEpoch{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageNewEpoch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageNewEpoch) ProtoMessage() {}

func (x *MessageNewEpoch) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*MessageNewEpoch) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{0}
}

func (x *MessageNewEpoch) GetEpoch() uint64 {
	if x != nil {
		return x.Epoch
	}
	return 0
}

type MessageObserveReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Round uint64 `protobuf:"varint,1,opt,name=round,proto3" json:"round,omitempty"`
	Epoch uint64 `protobuf:"varint,2,opt,name=epoch,proto3" json:"epoch,omitempty"`
}

func (x *MessageObserveReq) Reset() {
	*x = MessageObserveReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageObserveReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageObserveReq) ProtoMessage() {}

func (x *MessageObserveReq) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*MessageObserveReq) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{1}
}

func (x *MessageObserveReq) GetRound() uint64 {
	if x != nil {
		return x.Round
	}
	return 0
}

func (x *MessageObserveReq) GetEpoch() uint64 {
	if x != nil {
		return x.Epoch
	}
	return 0
}

type ReportingContext struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConfigDigest []byte `protobuf:"bytes,1,opt,name=configDigest,proto3" json:"configDigest,omitempty"`
	Epoch        uint64 `protobuf:"varint,2,opt,name=epoch,proto3" json:"epoch,omitempty"`
	Round        uint64 `protobuf:"varint,3,opt,name=round,proto3" json:"round,omitempty"`
}

func (x *ReportingContext) Reset() {
	*x = ReportingContext{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReportingContext) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReportingContext) ProtoMessage() {}

func (x *ReportingContext) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ReportingContext) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{2}
}

func (x *ReportingContext) GetConfigDigest() []byte {
	if x != nil {
		return x.ConfigDigest
	}
	return nil
}

func (x *ReportingContext) GetEpoch() uint64 {
	if x != nil {
		return x.Epoch
	}
	return 0
}

func (x *ReportingContext) GetRound() uint64 {
	if x != nil {
		return x.Round
	}
	return 0
}

type ObservationValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value []byte `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *ObservationValue) Reset() {
	*x = ObservationValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ObservationValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ObservationValue) ProtoMessage() {}

func (x *ObservationValue) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ObservationValue) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{3}
}

func (x *ObservationValue) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

type Observation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ctx       *ReportingContext `protobuf:"bytes,1,opt,name=ctx,proto3" json:"ctx,omitempty"`
	Value     *ObservationValue `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Signature []byte            `protobuf:"bytes,3,opt,name=signature,proto3" json:"signature,omitempty"`
	OracleID  uint32            `protobuf:"varint,4,opt,name=oracleID,proto3" json:"oracleID,omitempty"`
}

func (x *Observation) Reset() {
	*x = Observation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Observation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Observation) ProtoMessage() {}

func (x *Observation) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*Observation) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{4}
}

func (x *Observation) GetCtx() *ReportingContext {
	if x != nil {
		return x.Ctx
	}
	return nil
}

func (x *Observation) GetValue() *ObservationValue {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *Observation) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

func (x *Observation) GetOracleID() uint32 {
	if x != nil {
		return x.OracleID
	}
	return 0
}

type MessageObserve struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Epoch uint64       `protobuf:"varint,1,opt,name=epoch,proto3" json:"epoch,omitempty"`
	Round uint64       `protobuf:"varint,2,opt,name=round,proto3" json:"round,omitempty"`
	Obs   *Observation `protobuf:"bytes,3,opt,name=obs,proto3" json:"obs,omitempty"`
}

func (x *MessageObserve) Reset() {
	*x = MessageObserve{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageObserve) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageObserve) ProtoMessage() {}

func (x *MessageObserve) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*MessageObserve) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{5}
}

func (x *MessageObserve) GetEpoch() uint64 {
	if x != nil {
		return x.Epoch
	}
	return 0
}

func (x *MessageObserve) GetRound() uint64 {
	if x != nil {
		return x.Round
	}
	return 0
}

func (x *MessageObserve) GetObs() *Observation {
	if x != nil {
		return x.Obs
	}
	return nil
}

type MessageReportReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Round        uint64         `protobuf:"varint,1,opt,name=round,proto3" json:"round,omitempty"`
	Epoch        uint64         `protobuf:"varint,2,opt,name=epoch,proto3" json:"epoch,omitempty"`
	Observations []*Observation `protobuf:"bytes,3,rep,name=observations,proto3" json:"observations,omitempty"`
}

func (x *MessageReportReq) Reset() {
	*x = MessageReportReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageReportReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageReportReq) ProtoMessage() {}

func (x *MessageReportReq) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*MessageReportReq) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{6}
}

func (x *MessageReportReq) GetRound() uint64 {
	if x != nil {
		return x.Round
	}
	return 0
}

func (x *MessageReportReq) GetEpoch() uint64 {
	if x != nil {
		return x.Epoch
	}
	return 0
}

func (x *MessageReportReq) GetObservations() []*Observation {
	if x != nil {
		return x.Observations
	}
	return nil
}

type Signature struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Signature []byte `protobuf:"bytes,1,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (x *Signature) Reset() {
	*x = Signature{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Signature) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Signature) ProtoMessage() {}

func (x *Signature) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*Signature) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{7}
}

func (x *Signature) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

type OracleValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OracleID uint32            `protobuf:"varint,1,opt,name=oracleID,proto3" json:"oracleID,omitempty"`
	Value    *ObservationValue `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *OracleValue) Reset() {
	*x = OracleValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OracleValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OracleValue) ProtoMessage() {}

func (x *OracleValue) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*OracleValue) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{8}
}

func (x *OracleValue) GetOracleID() uint32 {
	if x != nil {
		return x.OracleID
	}
	return 0
}

func (x *OracleValue) GetValue() *ObservationValue {
	if x != nil {
		return x.Value
	}
	return nil
}

type ContractReport struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ctx    *ReportingContext `protobuf:"bytes,1,opt,name=ctx,proto3" json:"ctx,omitempty"`
	Values []*OracleValue    `protobuf:"bytes,2,rep,name=values,proto3" json:"values,omitempty"`
	Sig    []byte            `protobuf:"bytes,3,opt,name=sig,proto3" json:"sig,omitempty"`
}

func (x *ContractReport) Reset() {
	*x = ContractReport{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ContractReport) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContractReport) ProtoMessage() {}

func (x *ContractReport) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ContractReport) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{9}
}

func (x *ContractReport) GetCtx() *ReportingContext {
	if x != nil {
		return x.Ctx
	}
	return nil
}

func (x *ContractReport) GetValues() []*OracleValue {
	if x != nil {
		return x.Values
	}
	return nil
}

func (x *ContractReport) GetSig() []byte {
	if x != nil {
		return x.Sig
	}
	return nil
}

type MessageReport struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Epoch          uint64          `protobuf:"varint,1,opt,name=epoch,proto3" json:"epoch,omitempty"`
	Round          uint64          `protobuf:"varint,2,opt,name=round,proto3" json:"round,omitempty"`
	ContractReport *ContractReport `protobuf:"bytes,3,opt,name=contractReport,proto3" json:"contractReport,omitempty"`
}

func (x *MessageReport) Reset() {
	*x = MessageReport{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageReport) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageReport) ProtoMessage() {}

func (x *MessageReport) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*MessageReport) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{10}
}

func (x *MessageReport) GetEpoch() uint64 {
	if x != nil {
		return x.Epoch
	}
	return 0
}

func (x *MessageReport) GetRound() uint64 {
	if x != nil {
		return x.Round
	}
	return 0
}

func (x *MessageReport) GetContractReport() *ContractReport {
	if x != nil {
		return x.ContractReport
	}
	return nil
}

type ContractReportWithSignatures struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ContractReport *ContractReport `protobuf:"bytes,1,opt,name=contractReport,proto3" json:"contractReport,omitempty"`
	Signatures     []*Signature    `protobuf:"bytes,2,rep,name=signatures,proto3" json:"signatures,omitempty"`
}

func (x *ContractReportWithSignatures) Reset() {
	*x = ContractReportWithSignatures{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ContractReportWithSignatures) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContractReportWithSignatures) ProtoMessage() {}

func (x *ContractReportWithSignatures) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ContractReportWithSignatures) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{11}
}

func (x *ContractReportWithSignatures) GetContractReport() *ContractReport {
	if x != nil {
		return x.ContractReport
	}
	return nil
}

func (x *ContractReportWithSignatures) GetSignatures() []*Signature {
	if x != nil {
		return x.Signatures
	}
	return nil
}

type MessageFinal struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Epoch  uint64                        `protobuf:"varint,1,opt,name=epoch,proto3" json:"epoch,omitempty"`
	Round  uint64                        `protobuf:"varint,2,opt,name=round,proto3" json:"round,omitempty"`
	Report *ContractReportWithSignatures `protobuf:"bytes,3,opt,name=report,proto3" json:"report,omitempty"`
}

func (x *MessageFinal) Reset() {
	*x = MessageFinal{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageFinal) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageFinal) ProtoMessage() {}

func (x *MessageFinal) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*MessageFinal) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{12}
}

func (x *MessageFinal) GetEpoch() uint64 {
	if x != nil {
		return x.Epoch
	}
	return 0
}

func (x *MessageFinal) GetRound() uint64 {
	if x != nil {
		return x.Round
	}
	return 0
}

func (x *MessageFinal) GetReport() *ContractReportWithSignatures {
	if x != nil {
		return x.Report
	}
	return nil
}

type MessageFinalEcho struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Final *MessageFinal `protobuf:"bytes,1,opt,name=Final,proto3" json:"Final,omitempty"`
}

func (x *MessageFinalEcho) Reset() {
	*x = MessageFinalEcho{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[13]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageFinalEcho) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageFinalEcho) ProtoMessage() {}

func (x *MessageFinalEcho) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[13]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*MessageFinalEcho) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{13}
}

func (x *MessageFinalEcho) GetFinal() *MessageFinal {
	if x != nil {
		return x.Final
	}
	return nil
}

type MessageWrapper struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg isMessageWrapper_Msg `protobuf_oneof:"msg"`
}

func (x *MessageWrapper) Reset() {
	*x = MessageWrapper{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[14]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageWrapper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageWrapper) ProtoMessage() {}

func (x *MessageWrapper) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[14]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*MessageWrapper) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{14}
}

func (m *MessageWrapper) GetMsg() isMessageWrapper_Msg {
	if m != nil {
		return m.Msg
	}
	return nil
}

func (x *MessageWrapper) GetMessageNewEpoch() *MessageNewEpoch {
	if x, ok := x.GetMsg().(*MessageWrapper_MessageNewEpoch); ok {
		return x.MessageNewEpoch
	}
	return nil
}

func (x *MessageWrapper) GetMessageObserveReq() *MessageObserveReq {
	if x, ok := x.GetMsg().(*MessageWrapper_MessageObserveReq); ok {
		return x.MessageObserveReq
	}
	return nil
}

func (x *MessageWrapper) GetMessageObserve() *MessageObserve {
	if x, ok := x.GetMsg().(*MessageWrapper_MessageObserve); ok {
		return x.MessageObserve
	}
	return nil
}

func (x *MessageWrapper) GetMessageReportReq() *MessageReportReq {
	if x, ok := x.GetMsg().(*MessageWrapper_MessageReportReq); ok {
		return x.MessageReportReq
	}
	return nil
}

func (x *MessageWrapper) GetMessageReport() *MessageReport {
	if x, ok := x.GetMsg().(*MessageWrapper_MessageReport); ok {
		return x.MessageReport
	}
	return nil
}

func (x *MessageWrapper) GetMessageFinal() *MessageFinal {
	if x, ok := x.GetMsg().(*MessageWrapper_MessageFinal); ok {
		return x.MessageFinal
	}
	return nil
}

func (x *MessageWrapper) GetMessageFinalEcho() *MessageFinalEcho {
	if x, ok := x.GetMsg().(*MessageWrapper_MessageFinalEcho); ok {
		return x.MessageFinalEcho
	}
	return nil
}

type isMessageWrapper_Msg interface {
	isMessageWrapper_Msg()
}

type MessageWrapper_MessageNewEpoch struct {
	MessageNewEpoch *MessageNewEpoch `protobuf:"bytes,2,opt,name=messageNewEpoch,proto3,oneof"`
}

type MessageWrapper_MessageObserveReq struct {
	MessageObserveReq *MessageObserveReq `protobuf:"bytes,3,opt,name=messageObserveReq,proto3,oneof"`
}

type MessageWrapper_MessageObserve struct {
	MessageObserve *MessageObserve `protobuf:"bytes,4,opt,name=messageObserve,proto3,oneof"`
}

type MessageWrapper_MessageReportReq struct {
	MessageReportReq *MessageReportReq `protobuf:"bytes,5,opt,name=messageReportReq,proto3,oneof"`
}

type MessageWrapper_MessageReport struct {
	MessageReport *MessageReport `protobuf:"bytes,6,opt,name=messageReport,proto3,oneof"`
}

type MessageWrapper_MessageFinal struct {
	MessageFinal *MessageFinal `protobuf:"bytes,7,opt,name=messageFinal,proto3,oneof"`
}

type MessageWrapper_MessageFinalEcho struct {
	MessageFinalEcho *MessageFinalEcho `protobuf:"bytes,8,opt,name=messageFinalEcho,proto3,oneof"`
}

func (*MessageWrapper_MessageNewEpoch) isMessageWrapper_Msg() {}

func (*MessageWrapper_MessageObserveReq) isMessageWrapper_Msg() {}

func (*MessageWrapper_MessageObserve) isMessageWrapper_Msg() {}

func (*MessageWrapper_MessageReportReq) isMessageWrapper_Msg() {}

func (*MessageWrapper_MessageReport) isMessageWrapper_Msg() {}

func (*MessageWrapper_MessageFinal) isMessageWrapper_Msg() {}

func (*MessageWrapper_MessageFinalEcho) isMessageWrapper_Msg() {}

var File_messages_proto protoreflect.FileDescriptor

var file_messages_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0d, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22,
	0x27, 0x0a, 0x0f, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4e, 0x65, 0x77, 0x45, 0x70, 0x6f,
	0x63, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x70, 0x6f, 0x63, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x05, 0x65, 0x70, 0x6f, 0x63, 0x68, 0x22, 0x3f, 0x0a, 0x11, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x52, 0x65, 0x71, 0x12, 0x14, 0x0a,
	0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x72, 0x6f,
	0x75, 0x6e, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x70, 0x6f, 0x63, 0x68, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x05, 0x65, 0x70, 0x6f, 0x63, 0x68, 0x22, 0x62, 0x0a, 0x10, 0x52, 0x65, 0x70,
	0x6f, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x12, 0x22, 0x0a,
	0x0c, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x44, 0x69, 0x67, 0x65, 0x73, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x0c, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x44, 0x69, 0x67, 0x65, 0x73,
	0x74, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x70, 0x6f, 0x63, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x05, 0x65, 0x70, 0x6f, 0x63, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x22, 0x28, 0x0a,
	0x10, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0xb1, 0x01, 0x0a, 0x0b, 0x4f, 0x62, 0x73, 0x65,
	0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x31, 0x0a, 0x03, 0x63, 0x74, 0x78, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x43, 0x6f,
	0x6e, 0x74, 0x65, 0x78, 0x74, 0x52, 0x03, 0x63, 0x74, 0x78, 0x12, 0x35, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x73, 0x65, 0x72, 0x69,
	0x61, 0x6c, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12,
	0x1a, 0x0a, 0x08, 0x6f, 0x72, 0x61, 0x63, 0x6c, 0x65, 0x49, 0x44, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x08, 0x6f, 0x72, 0x61, 0x63, 0x6c, 0x65, 0x49, 0x44, 0x22, 0x6a, 0x0a, 0x0e, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x65, 0x70, 0x6f, 0x63, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x65, 0x70,
	0x6f, 0x63, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x12, 0x2c, 0x0a, 0x03, 0x6f, 0x62, 0x73,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69,
	0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x03, 0x6f, 0x62, 0x73, 0x22, 0x7e, 0x0a, 0x10, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x71, 0x12, 0x14, 0x0a, 0x05, 0x72,
	0x6f, 0x75, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x72, 0x6f, 0x75, 0x6e,
	0x64, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x70, 0x6f, 0x63, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x05, 0x65, 0x70, 0x6f, 0x63, 0x68, 0x12, 0x3e, 0x0a, 0x0c, 0x6f, 0x62, 0x73, 0x65, 0x72,
	0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x4f, 0x62,
	0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0c, 0x6f, 0x62, 0x73, 0x65, 0x72,
	0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x29, 0x0a, 0x09, 0x53, 0x69, 0x67, 0x6e, 0x61,
	0x74, 0x75, 0x72, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75,
	0x72, 0x65, 0x22, 0x60, 0x0a, 0x0b, 0x4f, 0x72, 0x61, 0x63, 0x6c, 0x65, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6f, 0x72, 0x61, 0x63, 0x6c, 0x65, 0x49, 0x44, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x08, 0x6f, 0x72, 0x61, 0x63, 0x6c, 0x65, 0x49, 0x44, 0x12, 0x35, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x73,
	0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x4f, 0x62, 0x73,
	0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x22, 0x89, 0x01, 0x0a, 0x0e, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63,
	0x74, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x31, 0x0a, 0x03, 0x63, 0x74, 0x78, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x43, 0x6f,
	0x6e, 0x74, 0x65, 0x78, 0x74, 0x52, 0x03, 0x63, 0x74, 0x78, 0x12, 0x32, 0x0a, 0x06, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x73, 0x65, 0x72,
	0x69, 0x61, 0x6c, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x4f, 0x72, 0x61, 0x63, 0x6c,
	0x65, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x12, 0x10,
	0x0a, 0x03, 0x73, 0x69, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x73, 0x69, 0x67,
	0x22, 0x82, 0x01, 0x0a, 0x0d, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x70, 0x6f, 0x63, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x05, 0x65, 0x70, 0x6f, 0x63, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x6f, 0x75, 0x6e,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x12, 0x45,
	0x0a, 0x0e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69,
	0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x52,
	0x65, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x0e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x52,
	0x65, 0x70, 0x6f, 0x72, 0x74, 0x22, 0x9f, 0x01, 0x0a, 0x1c, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61,
	0x63, 0x74, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x57, 0x69, 0x74, 0x68, 0x53, 0x69, 0x67, 0x6e,
	0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x12, 0x45, 0x0a, 0x0e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61,
	0x63, 0x74, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d,
	0x2e, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x43,
	0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x0e, 0x63,
	0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x38, 0x0a,
	0x0a, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x18, 0x2e, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x52, 0x0a, 0x73, 0x69, 0x67,
	0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x22, 0x7f, 0x0a, 0x0c, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x46, 0x69, 0x6e, 0x61, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x70, 0x6f, 0x63, 0x68,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x65, 0x70, 0x6f, 0x63, 0x68, 0x12, 0x14, 0x0a,
	0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x72, 0x6f,
	0x75, 0x6e, 0x64, 0x12, 0x43, 0x0a, 0x06, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x52, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x57, 0x69, 0x74, 0x68, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73,
	0x52, 0x06, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x22, 0x45, 0x0a, 0x10, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x46, 0x69, 0x6e, 0x61, 0x6c, 0x45, 0x63, 0x68, 0x6f, 0x12, 0x31, 0x0a, 0x05,
	0x46, 0x69, 0x6e, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x73, 0x65,
	0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x46, 0x69, 0x6e, 0x61, 0x6c, 0x52, 0x05, 0x46, 0x69, 0x6e, 0x61, 0x6c, 0x22,
	0xa5, 0x04, 0x0a, 0x0e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x57, 0x72, 0x61, 0x70, 0x70,
	0x65, 0x72, 0x12, 0x4a, 0x0a, 0x0f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4e, 0x65, 0x77,
	0x45, 0x70, 0x6f, 0x63, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x73, 0x65,
	0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x4e, 0x65, 0x77, 0x45, 0x70, 0x6f, 0x63, 0x68, 0x48, 0x00, 0x52, 0x0f, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4e, 0x65, 0x77, 0x45, 0x70, 0x6f, 0x63, 0x68, 0x12, 0x50,
	0x0a, 0x11, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x52, 0x65, 0x71, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x73, 0x65, 0x72, 0x69,
	0x61, 0x6c, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x52, 0x65, 0x71, 0x48, 0x00, 0x52, 0x11, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x52, 0x65, 0x71,
	0x12, 0x47, 0x0a, 0x0e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x62, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x73, 0x65, 0x72, 0x69, 0x61,
	0x6c, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x48, 0x00, 0x52, 0x0e, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x12, 0x4d, 0x0a, 0x10, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x71, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x72,
	0x74, 0x52, 0x65, 0x71, 0x48, 0x00, 0x52, 0x10, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52,
	0x65, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x71, 0x12, 0x44, 0x0a, 0x0d, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1c, 0x2e, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x48, 0x00, 0x52,
	0x0d, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x41,
	0x0a, 0x0c, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x46, 0x69, 0x6e, 0x61, 0x6c, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x46, 0x69, 0x6e, 0x61,
	0x6c, 0x48, 0x00, 0x52, 0x0c, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x46, 0x69, 0x6e, 0x61,
	0x6c, 0x12, 0x4d, 0x0a, 0x10, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x46, 0x69, 0x6e, 0x61,
	0x6c, 0x45, 0x63, 0x68, 0x6f, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x73, 0x65,
	0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x46, 0x69, 0x6e, 0x61, 0x6c, 0x45, 0x63, 0x68, 0x6f, 0x48, 0x00, 0x52, 0x10,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x46, 0x69, 0x6e, 0x61, 0x6c, 0x45, 0x63, 0x68, 0x6f,
	0x42, 0x05, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x42, 0x11, 0x5a, 0x0f, 0x2e, 0x3b, 0x73, 0x65, 0x72,
	0x69, 0x61, 0x6c, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_messages_proto_rawDescOnce sync.Once
	file_messages_proto_rawDescData = file_messages_proto_rawDesc
)

func file_messages_proto_rawDescGZIP() []byte {
	file_messages_proto_rawDescOnce.Do(func() {
		file_messages_proto_rawDescData = protoimpl.X.CompressGZIP(file_messages_proto_rawDescData)
	})
	return file_messages_proto_rawDescData
}

var file_messages_proto_msgTypes = make([]protoimpl.MessageInfo, 15)
var file_messages_proto_goTypes = []interface{}{
	(*MessageNewEpoch)(nil), (*MessageObserveReq)(nil), (*ReportingContext)(nil), (*ObservationValue)(nil), (*Observation)(nil), (*MessageObserve)(nil), (*MessageReportReq)(nil), (*Signature)(nil), (*OracleValue)(nil), (*ContractReport)(nil), (*MessageReport)(nil), (*ContractReportWithSignatures)(nil), (*MessageFinal)(nil), (*MessageFinalEcho)(nil), (*MessageWrapper)(nil)}
var file_messages_proto_depIdxs = []int32{
	2, 3, 4, 4, 3, 2, 8, 9, 9, 7, 11, 12, 0, 1, 5, 6, 10, 12, 13, 19, 19, 19, 19, 0}

func init() { file_messages_proto_init() }
func file_messages_proto_init() {
	if File_messages_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_messages_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageNewEpoch); i {
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
		file_messages_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageObserveReq); i {
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
		file_messages_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReportingContext); i {
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
		file_messages_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ObservationValue); i {
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
		file_messages_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Observation); i {
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
		file_messages_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageObserve); i {
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
		file_messages_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageReportReq); i {
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
		file_messages_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Signature); i {
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
		file_messages_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OracleValue); i {
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
		file_messages_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ContractReport); i {
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
		file_messages_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageReport); i {
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
		file_messages_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ContractReportWithSignatures); i {
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
		file_messages_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageFinal); i {
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
		file_messages_proto_msgTypes[13].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageFinalEcho); i {
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
		file_messages_proto_msgTypes[14].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageWrapper); i {
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
	file_messages_proto_msgTypes[14].OneofWrappers = []interface{}{
		(*MessageWrapper_MessageNewEpoch)(nil),
		(*MessageWrapper_MessageObserveReq)(nil),
		(*MessageWrapper_MessageObserve)(nil),
		(*MessageWrapper_MessageReportReq)(nil),
		(*MessageWrapper_MessageReport)(nil),
		(*MessageWrapper_MessageFinal)(nil),
		(*MessageWrapper_MessageFinalEcho)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_messages_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   15,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_messages_proto_goTypes,
		DependencyIndexes: file_messages_proto_depIdxs,
		MessageInfos:      file_messages_proto_msgTypes,
	}.Build()
	File_messages_proto = out.File
	file_messages_proto_rawDesc = nil
	file_messages_proto_goTypes = nil
	file_messages_proto_depIdxs = nil
}
