// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.25.1
// source: trie/proto_trienode.proto

package trie

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

type ProtoTrieNode struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The serialized trie node data.
	ProtoNodeData []byte `protobuf:"bytes,1,opt,name=protoNodeData,proto3" json:"protoNodeData,omitempty"`
}

func (x *ProtoTrieNode) Reset() {
	*x = ProtoTrieNode{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trie_proto_trienode_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProtoTrieNode) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProtoTrieNode) ProtoMessage() {}

func (x *ProtoTrieNode) ProtoReflect() protoreflect.Message {
	mi := &file_trie_proto_trienode_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProtoTrieNode.ProtoReflect.Descriptor instead.
func (*ProtoTrieNode) Descriptor() ([]byte, []int) {
	return file_trie_proto_trienode_proto_rawDescGZIP(), []int{0}
}

func (x *ProtoTrieNode) GetProtoNodeData() []byte {
	if x != nil {
		return x.ProtoNodeData
	}
	return nil
}

var File_trie_proto_trienode_proto protoreflect.FileDescriptor

var file_trie_proto_trienode_proto_rawDesc = []byte{
	0x0a, 0x19, 0x74, 0x72, 0x69, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x5f, 0x74, 0x72, 0x69,
	0x65, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x74, 0x72, 0x69,
	0x65, 0x22, 0x35, 0x0a, 0x0d, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x54, 0x72, 0x69, 0x65, 0x4e, 0x6f,
	0x64, 0x65, 0x12, 0x24, 0x0a, 0x0d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x4e, 0x6f, 0x64, 0x65, 0x44,
	0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0d, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x4e, 0x6f, 0x64, 0x65, 0x44, 0x61, 0x74, 0x61, 0x42, 0x2d, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x6f, 0x6d, 0x69, 0x6e, 0x61, 0x6e, 0x74, 0x2d,
	0x73, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x69, 0x65, 0x73, 0x2f, 0x67, 0x6f, 0x2d, 0x71, 0x75,
	0x61, 0x69, 0x2f, 0x74, 0x72, 0x69, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_trie_proto_trienode_proto_rawDescOnce sync.Once
	file_trie_proto_trienode_proto_rawDescData = file_trie_proto_trienode_proto_rawDesc
)

func file_trie_proto_trienode_proto_rawDescGZIP() []byte {
	file_trie_proto_trienode_proto_rawDescOnce.Do(func() {
		file_trie_proto_trienode_proto_rawDescData = protoimpl.X.CompressGZIP(file_trie_proto_trienode_proto_rawDescData)
	})
	return file_trie_proto_trienode_proto_rawDescData
}

var file_trie_proto_trienode_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_trie_proto_trienode_proto_goTypes = []interface{}{
	(*ProtoTrieNode)(nil), // 0: trie.ProtoTrieNode
}
var file_trie_proto_trienode_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_trie_proto_trienode_proto_init() }
func file_trie_proto_trienode_proto_init() {
	if File_trie_proto_trienode_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_trie_proto_trienode_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProtoTrieNode); i {
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
			RawDescriptor: file_trie_proto_trienode_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_trie_proto_trienode_proto_goTypes,
		DependencyIndexes: file_trie_proto_trienode_proto_depIdxs,
		MessageInfos:      file_trie_proto_trienode_proto_msgTypes,
	}.Build()
	File_trie_proto_trienode_proto = out.File
	file_trie_proto_trienode_proto_rawDesc = nil
	file_trie_proto_trienode_proto_goTypes = nil
	file_trie_proto_trienode_proto_depIdxs = nil
}
