// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.6.1
// source: protobuf/system-monitor.proto

package protobuf

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

type CPUstats struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	La   float64 `protobuf:"fixed64,1,opt,name=la,proto3" json:"la,omitempty"`
	Usr  float64 `protobuf:"fixed64,2,opt,name=usr,proto3" json:"usr,omitempty"`
	Sys  float64 `protobuf:"fixed64,3,opt,name=sys,proto3" json:"sys,omitempty"`
	Idle float64 `protobuf:"fixed64,4,opt,name=idle,proto3" json:"idle,omitempty"`
}

func (x *CPUstats) Reset() {
	*x = CPUstats{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_system_monitor_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CPUstats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CPUstats) ProtoMessage() {}

func (x *CPUstats) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_system_monitor_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CPUstats.ProtoReflect.Descriptor instead.
func (*CPUstats) Descriptor() ([]byte, []int) {
	return file_protobuf_system_monitor_proto_rawDescGZIP(), []int{0}
}

func (x *CPUstats) GetLa() float64 {
	if x != nil {
		return x.La
	}
	return 0
}

func (x *CPUstats) GetUsr() float64 {
	if x != nil {
		return x.Usr
	}
	return 0
}

func (x *CPUstats) GetSys() float64 {
	if x != nil {
		return x.Sys
	}
	return 0
}

func (x *CPUstats) GetIdle() float64 {
	if x != nil {
		return x.Idle
	}
	return 0
}

type DevStats struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name  string  `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Tps   float64 `protobuf:"fixed64,2,opt,name=tps,proto3" json:"tps,omitempty"`
	Read  float64 `protobuf:"fixed64,3,opt,name=read,proto3" json:"read,omitempty"`
	Write float64 `protobuf:"fixed64,4,opt,name=write,proto3" json:"write,omitempty"`
}

func (x *DevStats) Reset() {
	*x = DevStats{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_system_monitor_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DevStats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DevStats) ProtoMessage() {}

func (x *DevStats) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_system_monitor_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DevStats.ProtoReflect.Descriptor instead.
func (*DevStats) Descriptor() ([]byte, []int) {
	return file_protobuf_system_monitor_proto_rawDescGZIP(), []int{1}
}

func (x *DevStats) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *DevStats) GetTps() float64 {
	if x != nil {
		return x.Tps
	}
	return 0
}

func (x *DevStats) GetRead() float64 {
	if x != nil {
		return x.Read
	}
	return 0
}

func (x *DevStats) GetWrite() float64 {
	if x != nil {
		return x.Write
	}
	return 0
}

type FsStats struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name         string  `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Bytes        uint64  `protobuf:"varint,2,opt,name=bytes,proto3" json:"bytes,omitempty"`
	BytesPercent float64 `protobuf:"fixed64,3,opt,name=bytesPercent,proto3" json:"bytesPercent,omitempty"`
	Inode        uint64  `protobuf:"varint,4,opt,name=inode,proto3" json:"inode,omitempty"`
	InodePercent float64 `protobuf:"fixed64,5,opt,name=inodePercent,proto3" json:"inodePercent,omitempty"`
}

func (x *FsStats) Reset() {
	*x = FsStats{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_system_monitor_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FsStats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FsStats) ProtoMessage() {}

func (x *FsStats) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_system_monitor_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FsStats.ProtoReflect.Descriptor instead.
func (*FsStats) Descriptor() ([]byte, []int) {
	return file_protobuf_system_monitor_proto_rawDescGZIP(), []int{2}
}

func (x *FsStats) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *FsStats) GetBytes() uint64 {
	if x != nil {
		return x.Bytes
	}
	return 0
}

func (x *FsStats) GetBytesPercent() float64 {
	if x != nil {
		return x.BytesPercent
	}
	return 0
}

func (x *FsStats) GetInode() uint64 {
	if x != nil {
		return x.Inode
	}
	return 0
}

func (x *FsStats) GetInodePercent() float64 {
	if x != nil {
		return x.InodePercent
	}
	return 0
}

type Stats struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CPUstats *CPUstats   `protobuf:"bytes,1,opt,name=CPUstats,proto3" json:"CPUstats,omitempty"`
	DevStats []*DevStats `protobuf:"bytes,2,rep,name=DevStats,proto3" json:"DevStats,omitempty"`
	FsStats  []*FsStats  `protobuf:"bytes,3,rep,name=FsStats,proto3" json:"FsStats,omitempty"`
}

func (x *Stats) Reset() {
	*x = Stats{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_system_monitor_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Stats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Stats) ProtoMessage() {}

func (x *Stats) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_system_monitor_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Stats.ProtoReflect.Descriptor instead.
func (*Stats) Descriptor() ([]byte, []int) {
	return file_protobuf_system_monitor_proto_rawDescGZIP(), []int{3}
}

func (x *Stats) GetCPUstats() *CPUstats {
	if x != nil {
		return x.CPUstats
	}
	return nil
}

func (x *Stats) GetDevStats() []*DevStats {
	if x != nil {
		return x.DevStats
	}
	return nil
}

func (x *Stats) GetFsStats() []*FsStats {
	if x != nil {
		return x.FsStats
	}
	return nil
}

type Settings struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TimeBetweenTicks uint32 `protobuf:"varint,1,opt,name=timeBetweenTicks,proto3" json:"timeBetweenTicks,omitempty"`
	AveragingTime    uint32 `protobuf:"varint,2,opt,name=averagingTime,proto3" json:"averagingTime,omitempty"`
}

func (x *Settings) Reset() {
	*x = Settings{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_system_monitor_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Settings) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Settings) ProtoMessage() {}

func (x *Settings) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_system_monitor_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Settings.ProtoReflect.Descriptor instead.
func (*Settings) Descriptor() ([]byte, []int) {
	return file_protobuf_system_monitor_proto_rawDescGZIP(), []int{4}
}

func (x *Settings) GetTimeBetweenTicks() uint32 {
	if x != nil {
		return x.TimeBetweenTicks
	}
	return 0
}

func (x *Settings) GetAveragingTime() uint32 {
	if x != nil {
		return x.AveragingTime
	}
	return 0
}

var File_protobuf_system_monitor_proto protoreflect.FileDescriptor

var file_protobuf_system_monitor_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x79, 0x73, 0x74, 0x65,
	0x6d, 0x2d, 0x6d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0d, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x22, 0x52,
	0x0a, 0x08, 0x43, 0x50, 0x55, 0x73, 0x74, 0x61, 0x74, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x6c, 0x61,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x02, 0x6c, 0x61, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x73,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x75, 0x73, 0x72, 0x12, 0x10, 0x0a, 0x03,
	0x73, 0x79, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x73, 0x79, 0x73, 0x12, 0x12,
	0x0a, 0x04, 0x69, 0x64, 0x6c, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x04, 0x69, 0x64,
	0x6c, 0x65, 0x22, 0x5a, 0x0a, 0x08, 0x44, 0x65, 0x76, 0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x74, 0x70, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52,
	0x03, 0x74, 0x70, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x65, 0x61, 0x64, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x04, 0x72, 0x65, 0x61, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x77, 0x72, 0x69, 0x74,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x77, 0x72, 0x69, 0x74, 0x65, 0x22, 0x91,
	0x01, 0x0a, 0x07, 0x46, 0x73, 0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x62, 0x79, 0x74, 0x65, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x62,
	0x79, 0x74, 0x65, 0x73, 0x12, 0x22, 0x0a, 0x0c, 0x62, 0x79, 0x74, 0x65, 0x73, 0x50, 0x65, 0x72,
	0x63, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0c, 0x62, 0x79, 0x74, 0x65,
	0x73, 0x50, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x6f, 0x64,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x69, 0x6e, 0x6f, 0x64, 0x65, 0x12, 0x22,
	0x0a, 0x0c, 0x69, 0x6e, 0x6f, 0x64, 0x65, 0x50, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x0c, 0x69, 0x6e, 0x6f, 0x64, 0x65, 0x50, 0x65, 0x72, 0x63, 0x65,
	0x6e, 0x74, 0x22, 0xa3, 0x01, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x33, 0x0a, 0x08,
	0x43, 0x50, 0x55, 0x73, 0x74, 0x61, 0x74, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17,
	0x2e, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x43,
	0x50, 0x55, 0x73, 0x74, 0x61, 0x74, 0x73, 0x52, 0x08, 0x43, 0x50, 0x55, 0x73, 0x74, 0x61, 0x74,
	0x73, 0x12, 0x33, 0x0a, 0x08, 0x44, 0x65, 0x76, 0x53, 0x74, 0x61, 0x74, 0x73, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x4d, 0x6f, 0x6e, 0x69,
	0x74, 0x6f, 0x72, 0x2e, 0x44, 0x65, 0x76, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x08, 0x44, 0x65,
	0x76, 0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x30, 0x0a, 0x07, 0x46, 0x73, 0x53, 0x74, 0x61, 0x74,
	0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d,
	0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x46, 0x73, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52,
	0x07, 0x46, 0x73, 0x53, 0x74, 0x61, 0x74, 0x73, 0x22, 0x5c, 0x0a, 0x08, 0x53, 0x65, 0x74, 0x74,
	0x69, 0x6e, 0x67, 0x73, 0x12, 0x2a, 0x0a, 0x10, 0x74, 0x69, 0x6d, 0x65, 0x42, 0x65, 0x74, 0x77,
	0x65, 0x65, 0x6e, 0x54, 0x69, 0x63, 0x6b, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x10,
	0x74, 0x69, 0x6d, 0x65, 0x42, 0x65, 0x74, 0x77, 0x65, 0x65, 0x6e, 0x54, 0x69, 0x63, 0x6b, 0x73,
	0x12, 0x24, 0x0a, 0x0d, 0x61, 0x76, 0x65, 0x72, 0x61, 0x67, 0x69, 0x6e, 0x67, 0x54, 0x69, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0d, 0x61, 0x76, 0x65, 0x72, 0x61, 0x67, 0x69,
	0x6e, 0x67, 0x54, 0x69, 0x6d, 0x65, 0x32, 0x48, 0x0a, 0x07, 0x6d, 0x6f, 0x6e, 0x69, 0x74, 0x6f,
	0x72, 0x12, 0x3d, 0x0a, 0x08, 0x67, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x17, 0x2e,
	0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x53, 0x65,
	0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x1a, 0x14, 0x2e, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x4d,
	0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x73, 0x22, 0x00, 0x30, 0x01,
	0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protobuf_system_monitor_proto_rawDescOnce sync.Once
	file_protobuf_system_monitor_proto_rawDescData = file_protobuf_system_monitor_proto_rawDesc
)

func file_protobuf_system_monitor_proto_rawDescGZIP() []byte {
	file_protobuf_system_monitor_proto_rawDescOnce.Do(func() {
		file_protobuf_system_monitor_proto_rawDescData = protoimpl.X.CompressGZIP(file_protobuf_system_monitor_proto_rawDescData)
	})
	return file_protobuf_system_monitor_proto_rawDescData
}

var file_protobuf_system_monitor_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_protobuf_system_monitor_proto_goTypes = []interface{}{
	(*CPUstats)(nil), // 0: systemMonitor.CPUstats
	(*DevStats)(nil), // 1: systemMonitor.DevStats
	(*FsStats)(nil),  // 2: systemMonitor.FsStats
	(*Stats)(nil),    // 3: systemMonitor.Stats
	(*Settings)(nil), // 4: systemMonitor.Settings
}
var file_protobuf_system_monitor_proto_depIdxs = []int32{
	0, // 0: systemMonitor.Stats.CPUstats:type_name -> systemMonitor.CPUstats
	1, // 1: systemMonitor.Stats.DevStats:type_name -> systemMonitor.DevStats
	2, // 2: systemMonitor.Stats.FsStats:type_name -> systemMonitor.FsStats
	4, // 3: systemMonitor.monitor.getStats:input_type -> systemMonitor.Settings
	3, // 4: systemMonitor.monitor.getStats:output_type -> systemMonitor.Stats
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_protobuf_system_monitor_proto_init() }
func file_protobuf_system_monitor_proto_init() {
	if File_protobuf_system_monitor_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protobuf_system_monitor_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CPUstats); i {
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
		file_protobuf_system_monitor_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DevStats); i {
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
		file_protobuf_system_monitor_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FsStats); i {
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
		file_protobuf_system_monitor_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Stats); i {
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
		file_protobuf_system_monitor_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Settings); i {
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
			RawDescriptor: file_protobuf_system_monitor_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protobuf_system_monitor_proto_goTypes,
		DependencyIndexes: file_protobuf_system_monitor_proto_depIdxs,
		MessageInfos:      file_protobuf_system_monitor_proto_msgTypes,
	}.Build()
	File_protobuf_system_monitor_proto = out.File
	file_protobuf_system_monitor_proto_rawDesc = nil
	file_protobuf_system_monitor_proto_goTypes = nil
	file_protobuf_system_monitor_proto_depIdxs = nil
}
