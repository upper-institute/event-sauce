// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        (unknown)
// source: v1/paprika.proto

package apiv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Snapshot struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id               string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Version          string                 `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	NaturalTimestamp *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=natural_timestamp,json=naturalTimestamp,proto3" json:"natural_timestamp,omitempty"`
	StoreTimestamp   *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=store_timestamp,json=storeTimestamp,proto3" json:"store_timestamp,omitempty"`
	Label            string                 `protobuf:"bytes,5,opt,name=label,proto3" json:"label,omitempty"`
	Payload          *anypb.Any             `protobuf:"bytes,6,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *Snapshot) Reset() {
	*x = Snapshot{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_paprika_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Snapshot) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Snapshot) ProtoMessage() {}

func (x *Snapshot) ProtoReflect() protoreflect.Message {
	mi := &file_v1_paprika_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Snapshot.ProtoReflect.Descriptor instead.
func (*Snapshot) Descriptor() ([]byte, []int) {
	return file_v1_paprika_proto_rawDescGZIP(), []int{0}
}

func (x *Snapshot) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Snapshot) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *Snapshot) GetNaturalTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.NaturalTimestamp
	}
	return nil
}

func (x *Snapshot) GetStoreTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.StoreTimestamp
	}
	return nil
}

func (x *Snapshot) GetLabel() string {
	if x != nil {
		return x.Label
	}
	return ""
}

func (x *Snapshot) GetPayload() *anypb.Any {
	if x != nil {
		return x.Payload
	}
	return nil
}

type Snapshot_SetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Snapshot *Snapshot `protobuf:"bytes,1,opt,name=snapshot,proto3" json:"snapshot,omitempty"`
}

func (x *Snapshot_SetResponse) Reset() {
	*x = Snapshot_SetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_paprika_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Snapshot_SetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Snapshot_SetResponse) ProtoMessage() {}

func (x *Snapshot_SetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v1_paprika_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Snapshot_SetResponse.ProtoReflect.Descriptor instead.
func (*Snapshot_SetResponse) Descriptor() ([]byte, []int) {
	return file_v1_paprika_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Snapshot_SetResponse) GetSnapshot() *Snapshot {
	if x != nil {
		return x.Snapshot
	}
	return nil
}

type Snapshot_GetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *Snapshot_GetRequest) Reset() {
	*x = Snapshot_GetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_paprika_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Snapshot_GetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Snapshot_GetRequest) ProtoMessage() {}

func (x *Snapshot_GetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_paprika_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Snapshot_GetRequest.ProtoReflect.Descriptor instead.
func (*Snapshot_GetRequest) Descriptor() ([]byte, []int) {
	return file_v1_paprika_proto_rawDescGZIP(), []int{0, 1}
}

func (x *Snapshot_GetRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

var File_v1_paprika_proto protoreflect.FileDescriptor

var file_v1_paprika_proto_rawDesc = []byte{
	0x0a, 0x10, 0x76, 0x31, 0x2f, 0x70, 0x61, 0x70, 0x72, 0x69, 0x6b, 0x61, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x11, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x61, 0x75, 0x63, 0x65, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xee, 0x02, 0x0a, 0x08, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x18,
	0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x47, 0x0a, 0x11, 0x6e, 0x61, 0x74, 0x75,
	0x72, 0x61, 0x6c, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x10, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x61, 0x6c, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x12, 0x43, 0x0a, 0x0f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x12, 0x2e, 0x0a, 0x07,
	0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x41, 0x6e, 0x79, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x1a, 0x46, 0x0a, 0x0b,
	0x53, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x37, 0x0a, 0x08, 0x73,
	0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x61, 0x75, 0x63, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76,
	0x31, 0x2e, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x52, 0x08, 0x73, 0x6e, 0x61, 0x70,
	0x73, 0x68, 0x6f, 0x74, 0x1a, 0x1c, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x32, 0xad, 0x01, 0x0a, 0x0e, 0x50, 0x61, 0x70, 0x72, 0x69, 0x6b, 0x61, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4d, 0x0a, 0x03, 0x53, 0x65, 0x74, 0x12, 0x1b, 0x2e, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x73, 0x61, 0x75, 0x63, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31,
	0x2e, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x1a, 0x27, 0x2e, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x73, 0x61, 0x75, 0x63, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x6e,
	0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x2e, 0x53, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x12, 0x4c, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x26, 0x2e, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x73, 0x61, 0x75, 0x63, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e,
	0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x61, 0x75, 0x63, 0x65,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74,
	0x22, 0x00, 0x42, 0xc7, 0x01, 0x0a, 0x15, 0x63, 0x6f, 0x6d, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x73, 0x61, 0x75, 0x63, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x42, 0x0c, 0x50, 0x61,
	0x70, 0x72, 0x69, 0x6b, 0x61, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x3a, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x75, 0x70, 0x70, 0x65, 0x72, 0x2d, 0x69,
	0x6e, 0x73, 0x74, 0x69, 0x74, 0x75, 0x74, 0x65, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2d, 0x73,
	0x61, 0x75, 0x63, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f,
	0x76, 0x31, 0x3b, 0x61, 0x70, 0x69, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x45, 0x41, 0x58, 0xaa, 0x02,
	0x11, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x61, 0x75, 0x63, 0x65, 0x2e, 0x41, 0x70, 0x69, 0x2e,
	0x56, 0x31, 0xca, 0x02, 0x11, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x61, 0x75, 0x63, 0x65, 0x5c,
	0x41, 0x70, 0x69, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x1d, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x61,
	0x75, 0x63, 0x65, 0x5c, 0x41, 0x70, 0x69, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x13, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x61,
	0x75, 0x63, 0x65, 0x3a, 0x3a, 0x41, 0x70, 0x69, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_v1_paprika_proto_rawDescOnce sync.Once
	file_v1_paprika_proto_rawDescData = file_v1_paprika_proto_rawDesc
)

func file_v1_paprika_proto_rawDescGZIP() []byte {
	file_v1_paprika_proto_rawDescOnce.Do(func() {
		file_v1_paprika_proto_rawDescData = protoimpl.X.CompressGZIP(file_v1_paprika_proto_rawDescData)
	})
	return file_v1_paprika_proto_rawDescData
}

var file_v1_paprika_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_v1_paprika_proto_goTypes = []interface{}{
	(*Snapshot)(nil),              // 0: eventsauce.api.v1.Snapshot
	(*Snapshot_SetResponse)(nil),  // 1: eventsauce.api.v1.Snapshot.SetResponse
	(*Snapshot_GetRequest)(nil),   // 2: eventsauce.api.v1.Snapshot.GetRequest
	(*timestamppb.Timestamp)(nil), // 3: google.protobuf.Timestamp
	(*anypb.Any)(nil),             // 4: google.protobuf.Any
}
var file_v1_paprika_proto_depIdxs = []int32{
	3, // 0: eventsauce.api.v1.Snapshot.natural_timestamp:type_name -> google.protobuf.Timestamp
	3, // 1: eventsauce.api.v1.Snapshot.store_timestamp:type_name -> google.protobuf.Timestamp
	4, // 2: eventsauce.api.v1.Snapshot.payload:type_name -> google.protobuf.Any
	0, // 3: eventsauce.api.v1.Snapshot.SetResponse.snapshot:type_name -> eventsauce.api.v1.Snapshot
	0, // 4: eventsauce.api.v1.PaprikaService.Set:input_type -> eventsauce.api.v1.Snapshot
	2, // 5: eventsauce.api.v1.PaprikaService.Get:input_type -> eventsauce.api.v1.Snapshot.GetRequest
	1, // 6: eventsauce.api.v1.PaprikaService.Set:output_type -> eventsauce.api.v1.Snapshot.SetResponse
	0, // 7: eventsauce.api.v1.PaprikaService.Get:output_type -> eventsauce.api.v1.Snapshot
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_v1_paprika_proto_init() }
func file_v1_paprika_proto_init() {
	if File_v1_paprika_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_v1_paprika_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Snapshot); i {
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
		file_v1_paprika_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Snapshot_SetResponse); i {
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
		file_v1_paprika_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Snapshot_GetRequest); i {
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
			RawDescriptor: file_v1_paprika_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_v1_paprika_proto_goTypes,
		DependencyIndexes: file_v1_paprika_proto_depIdxs,
		MessageInfos:      file_v1_paprika_proto_msgTypes,
	}.Build()
	File_v1_paprika_proto = out.File
	file_v1_paprika_proto_rawDesc = nil
	file_v1_paprika_proto_goTypes = nil
	file_v1_paprika_proto_depIdxs = nil
}