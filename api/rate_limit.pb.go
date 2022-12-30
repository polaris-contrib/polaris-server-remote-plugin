//*
// Tencent is pleased to support the open source community by making Polaris available.
//
// Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
//
// Licensed under the BSD 3-Clause License (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// https://opensource.org/licenses/BSD-3-Clause
//
// Unless required by applicable law or agreed to in writing, software distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.21.9
// source: rate_limit.proto

package api

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

// RatelimitType rate limit type
type RatelimitType int32

const (
	// IPRatelimit ip-based rate limit type.
	RatelimitType_IPRatelimit RatelimitType = 0
	// APIRatelimit api-based rate limit type.
	RatelimitType_APIRatelimit RatelimitType = 1
	// ServiceRatelimit service-based rate limit type.
	RatelimitType_ServiceRatelimit RatelimitType = 2
	// InstanceRatelimit instance-based rate limit type.
	RatelimitType_InstanceRatelimit RatelimitType = 3
)

// Enum value maps for RatelimitType.
var (
	RatelimitType_name = map[int32]string{
		0: "IPRatelimit",
		1: "APIRatelimit",
		2: "ServiceRatelimit",
		3: "InstanceRatelimit",
	}
	RatelimitType_value = map[string]int32{
		"IPRatelimit":       0,
		"APIRatelimit":      1,
		"ServiceRatelimit":  2,
		"InstanceRatelimit": 3,
	}
)

func (x RatelimitType) Enum() *RatelimitType {
	p := new(RatelimitType)
	*p = x
	return p
}

func (x RatelimitType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RatelimitType) Descriptor() protoreflect.EnumDescriptor {
	return file_rate_limit_proto_enumTypes[0].Descriptor()
}

func (RatelimitType) Type() protoreflect.EnumType {
	return &file_rate_limit_proto_enumTypes[0]
}

func (x RatelimitType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RatelimitType.Descriptor instead.
func (RatelimitType) EnumDescriptor() ([]byte, []int) {
	return file_rate_limit_proto_rawDescGZIP(), []int{0}
}

// RateLimitPluginRequest rate limit remote plugin request.
type RateLimitRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Type which rate limit type.
	Type RatelimitType `protobuf:"varint,1,opt,name=type,proto3,enum=api.RatelimitType" json:"type,omitempty"`
	// Key represents for api name or api identification.
	Key string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *RateLimitRequest) Reset() {
	*x = RateLimitRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rate_limit_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RateLimitRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RateLimitRequest) ProtoMessage() {}

func (x *RateLimitRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rate_limit_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RateLimitRequest.ProtoReflect.Descriptor instead.
func (*RateLimitRequest) Descriptor() ([]byte, []int) {
	return file_rate_limit_proto_rawDescGZIP(), []int{0}
}

func (x *RateLimitRequest) GetType() RatelimitType {
	if x != nil {
		return x.Type
	}
	return RatelimitType_IPRatelimit
}

func (x *RateLimitRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

// RateLimitPluginResponse rate limit remote plugin response.
type RateLimitResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Allow if allow current req set true, else set false.
	Allow bool `protobuf:"varint,1,opt,name=allow,proto3" json:"allow,omitempty"`
}

func (x *RateLimitResponse) Reset() {
	*x = RateLimitResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rate_limit_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RateLimitResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RateLimitResponse) ProtoMessage() {}

func (x *RateLimitResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rate_limit_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RateLimitResponse.ProtoReflect.Descriptor instead.
func (*RateLimitResponse) Descriptor() ([]byte, []int) {
	return file_rate_limit_proto_rawDescGZIP(), []int{1}
}

func (x *RateLimitResponse) GetAllow() bool {
	if x != nil {
		return x.Allow
	}
	return false
}

var File_rate_limit_proto protoreflect.FileDescriptor

var file_rate_limit_proto_rawDesc = []byte{
	0x0a, 0x10, 0x72, 0x61, 0x74, 0x65, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x03, 0x61, 0x70, 0x69, 0x1a, 0x0c, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4c, 0x0a, 0x10, 0x52, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6d,
	0x69, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x26, 0x0a, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x52, 0x61,
	0x74, 0x65, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x22, 0x29, 0x0a, 0x11, 0x52, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x6c, 0x6c, 0x6f,
	0x77, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x2a, 0x5f,
	0x0a, 0x0d, 0x52, 0x61, 0x74, 0x65, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x0f, 0x0a, 0x0b, 0x49, 0x50, 0x52, 0x61, 0x74, 0x65, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x10, 0x00,
	0x12, 0x10, 0x0a, 0x0c, 0x41, 0x50, 0x49, 0x52, 0x61, 0x74, 0x65, 0x6c, 0x69, 0x6d, 0x69, 0x74,
	0x10, 0x01, 0x12, 0x14, 0x0a, 0x10, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x61, 0x74,
	0x65, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x10, 0x02, 0x12, 0x15, 0x0a, 0x11, 0x49, 0x6e, 0x73, 0x74,
	0x61, 0x6e, 0x63, 0x65, 0x52, 0x61, 0x74, 0x65, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x10, 0x03, 0x32,
	0x72, 0x0a, 0x0b, 0x52, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x12, 0x2b,
	0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x10, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x50, 0x69, 0x6e,
	0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x50,
	0x6f, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x36, 0x0a, 0x05, 0x41,
	0x6c, 0x6c, 0x6f, 0x77, 0x12, 0x15, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x52, 0x61, 0x74, 0x65, 0x4c,
	0x69, 0x6d, 0x69, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x52, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x42, 0x2c, 0x0a, 0x15, 0x63, 0x6e, 0x2e, 0x70, 0x6f, 0x6c, 0x61, 0x72, 0x69,
	0x73, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x61, 0x70, 0x69, 0x42, 0x09, 0x52, 0x61,
	0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x50, 0x00, 0x5a, 0x06, 0x2e, 0x2f, 0x3b, 0x61, 0x70,
	0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rate_limit_proto_rawDescOnce sync.Once
	file_rate_limit_proto_rawDescData = file_rate_limit_proto_rawDesc
)

func file_rate_limit_proto_rawDescGZIP() []byte {
	file_rate_limit_proto_rawDescOnce.Do(func() {
		file_rate_limit_proto_rawDescData = protoimpl.X.CompressGZIP(file_rate_limit_proto_rawDescData)
	})
	return file_rate_limit_proto_rawDescData
}

var file_rate_limit_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_rate_limit_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_rate_limit_proto_goTypes = []interface{}{
	(RatelimitType)(0),        // 0: api.RatelimitType
	(*RateLimitRequest)(nil),  // 1: api.RateLimitRequest
	(*RateLimitResponse)(nil), // 2: api.RateLimitResponse
	(*PingRequest)(nil),       // 3: api.PingRequest
	(*PongResponse)(nil),      // 4: api.PongResponse
}
var file_rate_limit_proto_depIdxs = []int32{
	0, // 0: api.RateLimitRequest.type:type_name -> api.RatelimitType
	3, // 1: api.RateLimiter.Ping:input_type -> api.PingRequest
	1, // 2: api.RateLimiter.Allow:input_type -> api.RateLimitRequest
	4, // 3: api.RateLimiter.Ping:output_type -> api.PongResponse
	2, // 4: api.RateLimiter.Allow:output_type -> api.RateLimitResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_rate_limit_proto_init() }
func file_rate_limit_proto_init() {
	if File_rate_limit_proto != nil {
		return
	}
	file_plugin_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_rate_limit_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RateLimitRequest); i {
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
		file_rate_limit_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RateLimitResponse); i {
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
			RawDescriptor: file_rate_limit_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_rate_limit_proto_goTypes,
		DependencyIndexes: file_rate_limit_proto_depIdxs,
		EnumInfos:         file_rate_limit_proto_enumTypes,
		MessageInfos:      file_rate_limit_proto_msgTypes,
	}.Build()
	File_rate_limit_proto = out.File
	file_rate_limit_proto_rawDesc = nil
	file_rate_limit_proto_goTypes = nil
	file_rate_limit_proto_depIdxs = nil
}