package protobuf

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)

	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type OffchainConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EncryptionPKs   [][]byte `protobuf:"bytes,2,rep,name=encryptionPKs,proto3" json:"encryptionPKs,omitempty"`
	SignaturePKs    [][]byte `protobuf:"bytes,3,rep,name=signaturePKs,proto3" json:"signaturePKs,omitempty"`
	EncryptionGroup string   `protobuf:"bytes,4,opt,name=encryptionGroup,proto3" json:"encryptionGroup,omitempty"`
	Translator      string   `protobuf:"bytes,5,opt,name=translator,proto3" json:"translator,omitempty"`
}

func (x *OffchainConfig) Reset() {
	*x = OffchainConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_offchain_config_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OffchainConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OffchainConfig) ProtoMessage() {}

func (x *OffchainConfig) ProtoReflect() protoreflect.Message {
	mi := &file_offchain_config_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*OffchainConfig) Descriptor() ([]byte, []int) {
	return file_offchain_config_proto_rawDescGZIP(), []int{0}
}

func (x *OffchainConfig) GetEncryptionPKs() [][]byte {
	if x != nil {
		return x.EncryptionPKs
	}
	return nil
}

func (x *OffchainConfig) GetSignaturePKs() [][]byte {
	if x != nil {
		return x.SignaturePKs
	}
	return nil
}

func (x *OffchainConfig) GetEncryptionGroup() string {
	if x != nil {
		return x.EncryptionGroup
	}
	return ""
}

func (x *OffchainConfig) GetTranslator() string {
	if x != nil {
		return x.Translator
	}
	return ""
}

var File_offchain_config_proto protoreflect.FileDescriptor

var file_offchain_config_proto_rawDesc = []byte{
	0x0a, 0x15, 0x6f, 0x66, 0x66, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x74, 0x79, 0x70, 0x65, 0x73, 0x22, 0xa4,
	0x01, 0x0a, 0x0e, 0x6f, 0x66, 0x66, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x12, 0x24, 0x0a, 0x0d, 0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x50,
	0x4b, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x0d, 0x65, 0x6e, 0x63, 0x72, 0x79, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x50, 0x4b, 0x73, 0x12, 0x22, 0x0a, 0x0c, 0x73, 0x69, 0x67, 0x6e, 0x61,
	0x74, 0x75, 0x72, 0x65, 0x50, 0x4b, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x0c, 0x73,
	0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x50, 0x4b, 0x73, 0x12, 0x28, 0x0a, 0x0f, 0x65,
	0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x47, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x1e, 0x0a, 0x0a, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61,
	0x74, 0x6f, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x74, 0x72, 0x61, 0x6e, 0x73,
	0x6c, 0x61, 0x74, 0x6f, 0x72, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_offchain_config_proto_rawDescOnce sync.Once
	file_offchain_config_proto_rawDescData = file_offchain_config_proto_rawDesc
)

func file_offchain_config_proto_rawDescGZIP() []byte {
	file_offchain_config_proto_rawDescOnce.Do(func() {
		file_offchain_config_proto_rawDescData = protoimpl.X.CompressGZIP(file_offchain_config_proto_rawDescData)
	})
	return file_offchain_config_proto_rawDescData
}

var file_offchain_config_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_offchain_config_proto_goTypes = []interface{}{
	(*OffchainConfig)(nil),
}
var file_offchain_config_proto_depIdxs = []int32{
	0,
	0,
	0,
	0,
	0,
}

func init() { file_offchain_config_proto_init() }
func file_offchain_config_proto_init() {
	if File_offchain_config_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_offchain_config_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OffchainConfig); i {
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
			RawDescriptor: file_offchain_config_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_offchain_config_proto_goTypes,
		DependencyIndexes: file_offchain_config_proto_depIdxs,
		MessageInfos:      file_offchain_config_proto_msgTypes,
	}.Build()
	File_offchain_config_proto = out.File
	file_offchain_config_proto_rawDesc = nil
	file_offchain_config_proto_goTypes = nil
	file_offchain_config_proto_depIdxs = nil
}
