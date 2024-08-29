// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.2
// source: common/proto_common.proto

package common

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

type ProtoLocation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value []byte `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *ProtoLocation) Reset() {
	*x = ProtoLocation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_common_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProtoLocation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProtoLocation) ProtoMessage() {}

func (x *ProtoLocation) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_common_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProtoLocation.ProtoReflect.Descriptor instead.
func (*ProtoLocation) Descriptor() ([]byte, []int) {
	return file_common_proto_common_proto_rawDescGZIP(), []int{0}
}

func (x *ProtoLocation) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

type ProtoHash struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value []byte `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *ProtoHash) Reset() {
	*x = ProtoHash{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_common_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProtoHash) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProtoHash) ProtoMessage() {}

func (x *ProtoHash) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_common_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProtoHash.ProtoReflect.Descriptor instead.
func (*ProtoHash) Descriptor() ([]byte, []int) {
	return file_common_proto_common_proto_rawDescGZIP(), []int{1}
}

func (x *ProtoHash) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

type ProtoHashes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hashes []*ProtoHash `protobuf:"bytes,1,rep,name=hashes,proto3" json:"hashes,omitempty"`
}

func (x *ProtoHashes) Reset() {
	*x = ProtoHashes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_common_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProtoHashes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProtoHashes) ProtoMessage() {}

func (x *ProtoHashes) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_common_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProtoHashes.ProtoReflect.Descriptor instead.
func (*ProtoHashes) Descriptor() ([]byte, []int) {
	return file_common_proto_common_proto_rawDescGZIP(), []int{2}
}

func (x *ProtoHashes) GetHashes() []*ProtoHash {
	if x != nil {
		return x.Hashes
	}
	return nil
}

type ProtoAddress struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value []byte `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *ProtoAddress) Reset() {
	*x = ProtoAddress{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_common_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProtoAddress) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProtoAddress) ProtoMessage() {}

func (x *ProtoAddress) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_common_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProtoAddress.ProtoReflect.Descriptor instead.
func (*ProtoAddress) Descriptor() ([]byte, []int) {
	return file_common_proto_common_proto_rawDescGZIP(), []int{3}
}

func (x *ProtoAddress) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

type ProtoNumber struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value uint64 `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *ProtoNumber) Reset() {
	*x = ProtoNumber{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_common_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProtoNumber) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProtoNumber) ProtoMessage() {}

func (x *ProtoNumber) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_common_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProtoNumber.ProtoReflect.Descriptor instead.
func (*ProtoNumber) Descriptor() ([]byte, []int) {
	return file_common_proto_common_proto_rawDescGZIP(), []int{4}
}

func (x *ProtoNumber) GetValue() uint64 {
	if x != nil {
		return x.Value
	}
	return 0
}

type ProtoLocations struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Locations []*ProtoLocation `protobuf:"bytes,1,rep,name=locations,proto3" json:"locations,omitempty"`
}

func (x *ProtoLocations) Reset() {
	*x = ProtoLocations{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_common_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProtoLocations) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProtoLocations) ProtoMessage() {}

func (x *ProtoLocations) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_common_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProtoLocations.ProtoReflect.Descriptor instead.
func (*ProtoLocations) Descriptor() ([]byte, []int) {
	return file_common_proto_common_proto_rawDescGZIP(), []int{5}
}

func (x *ProtoLocations) GetLocations() []*ProtoLocation {
	if x != nil {
		return x.Locations
	}
	return nil
}

var File_common_proto_common_proto protoreflect.FileDescriptor

var file_common_proto_common_proto_rawDesc = []byte{
	0x0a, 0x19, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x5f, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x22, 0x25, 0x0a, 0x0d, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x4c, 0x6f, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x21, 0x0a, 0x09, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x48, 0x61, 0x73, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x38, 0x0a,
	0x0b, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x48, 0x61, 0x73, 0x68, 0x65, 0x73, 0x12, 0x29, 0x0a, 0x06,
	0x68, 0x61, 0x73, 0x68, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x48, 0x61, 0x73, 0x68, 0x52,
	0x06, 0x68, 0x61, 0x73, 0x68, 0x65, 0x73, 0x22, 0x24, 0x0a, 0x0c, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x23, 0x0a,
	0x0b, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x22, 0x45, 0x0a, 0x0e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x4c, 0x6f, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x12, 0x33, 0x0a, 0x09, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x09,
	0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x42, 0x2f, 0x5a, 0x2d, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x6f, 0x6d, 0x69, 0x6e, 0x61, 0x6e, 0x74,
	0x2d, 0x73, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x69, 0x65, 0x73, 0x2f, 0x67, 0x6f, 0x2d, 0x71,
	0x75, 0x61, 0x69, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_common_proto_common_proto_rawDescOnce sync.Once
	file_common_proto_common_proto_rawDescData = file_common_proto_common_proto_rawDesc
)

func file_common_proto_common_proto_rawDescGZIP() []byte {
	file_common_proto_common_proto_rawDescOnce.Do(func() {
		file_common_proto_common_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_proto_common_proto_rawDescData)
	})
	return file_common_proto_common_proto_rawDescData
}

var file_common_proto_common_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_common_proto_common_proto_goTypes = []any{
	(*ProtoLocation)(nil),  // 0: common.ProtoLocation
	(*ProtoHash)(nil),      // 1: common.ProtoHash
	(*ProtoHashes)(nil),    // 2: common.ProtoHashes
	(*ProtoAddress)(nil),   // 3: common.ProtoAddress
	(*ProtoNumber)(nil),    // 4: common.ProtoNumber
	(*ProtoLocations)(nil), // 5: common.ProtoLocations
}
var file_common_proto_common_proto_depIdxs = []int32{
	1, // 0: common.ProtoHashes.hashes:type_name -> common.ProtoHash
	0, // 1: common.ProtoLocations.locations:type_name -> common.ProtoLocation
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_common_proto_common_proto_init() }
func file_common_proto_common_proto_init() {
	if File_common_proto_common_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_common_proto_common_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*ProtoLocation); i {
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
		file_common_proto_common_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*ProtoHash); i {
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
		file_common_proto_common_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*ProtoHashes); i {
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
		file_common_proto_common_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*ProtoAddress); i {
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
		file_common_proto_common_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*ProtoNumber); i {
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
		file_common_proto_common_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*ProtoLocations); i {
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
			RawDescriptor: file_common_proto_common_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_common_proto_common_proto_goTypes,
		DependencyIndexes: file_common_proto_common_proto_depIdxs,
		MessageInfos:      file_common_proto_common_proto_msgTypes,
	}.Build()
	File_common_proto_common_proto = out.File
	file_common_proto_common_proto_rawDesc = nil
	file_common_proto_common_proto_goTypes = nil
	file_common_proto_common_proto_depIdxs = nil
}
