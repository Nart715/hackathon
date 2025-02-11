// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.3
// source: proto/account/account.proto

package account

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Action int32

const (
	Action_CREDIT Action = 0
	Action_DEBIT  Action = 1
)

// Enum value maps for Action.
var (
	Action_name = map[int32]string{
		0: "CREDIT",
		1: "DEBIT",
	}
	Action_value = map[string]int32{
		"CREDIT": 0,
		"DEBIT":  1,
	}
)

func (x Action) Enum() *Action {
	p := new(Action)
	*p = x
	return p
}

func (x Action) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Action) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_account_account_proto_enumTypes[0].Descriptor()
}

func (Action) Type() protoreflect.EnumType {
	return &file_proto_account_account_proto_enumTypes[0]
}

func (x Action) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Action.Descriptor instead.
func (Action) EnumDescriptor() ([]byte, []int) {
	return file_proto_account_account_proto_rawDescGZIP(), []int{0}
}

type CreateAccountRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	AccountId     int32                  `protobuf:"varint,1,opt,name=accountId,proto3" json:"accountId,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateAccountRequest) Reset() {
	*x = CreateAccountRequest{}
	mi := &file_proto_account_account_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateAccountRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateAccountRequest) ProtoMessage() {}

func (x *CreateAccountRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_account_account_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateAccountRequest.ProtoReflect.Descriptor instead.
func (*CreateAccountRequest) Descriptor() ([]byte, []int) {
	return file_proto_account_account_proto_rawDescGZIP(), []int{0}
}

func (x *CreateAccountRequest) GetAccountId() int32 {
	if x != nil {
		return x.AccountId
	}
	return 0
}

type CreateAccountResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Code          int32                  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Data          *AccountData           `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateAccountResponse) Reset() {
	*x = CreateAccountResponse{}
	mi := &file_proto_account_account_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateAccountResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateAccountResponse) ProtoMessage() {}

func (x *CreateAccountResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_account_account_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateAccountResponse.ProtoReflect.Descriptor instead.
func (*CreateAccountResponse) Descriptor() ([]byte, []int) {
	return file_proto_account_account_proto_rawDescGZIP(), []int{1}
}

func (x *CreateAccountResponse) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *CreateAccountResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *CreateAccountResponse) GetData() *AccountData {
	if x != nil {
		return x.Data
	}
	return nil
}

type AccountData struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Ac            int32                  `protobuf:"varint,1,opt,name=ac,proto3" json:"ac,omitempty"`
	Bc            int32                  `protobuf:"varint,2,opt,name=bc,proto3" json:"bc,omitempty"`
	Tx            int64                  `protobuf:"varint,3,opt,name=tx,proto3" json:"tx,omitempty"`
	Itx           int64                  `protobuf:"varint,4,opt,name=itx,proto3" json:"itx,omitempty"`
	Up            int64                  `protobuf:"varint,5,opt,name=up,proto3" json:"up,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AccountData) Reset() {
	*x = AccountData{}
	mi := &file_proto_account_account_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AccountData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccountData) ProtoMessage() {}

func (x *AccountData) ProtoReflect() protoreflect.Message {
	mi := &file_proto_account_account_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccountData.ProtoReflect.Descriptor instead.
func (*AccountData) Descriptor() ([]byte, []int) {
	return file_proto_account_account_proto_rawDescGZIP(), []int{2}
}

func (x *AccountData) GetAc() int32 {
	if x != nil {
		return x.Ac
	}
	return 0
}

func (x *AccountData) GetBc() int32 {
	if x != nil {
		return x.Bc
	}
	return 0
}

func (x *AccountData) GetTx() int64 {
	if x != nil {
		return x.Tx
	}
	return 0
}

func (x *AccountData) GetItx() int64 {
	if x != nil {
		return x.Itx
	}
	return 0
}

func (x *AccountData) GetUp() int64 {
	if x != nil {
		return x.Up
	}
	return 0
}

type DepositAccountRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	AccountId     int32                  `protobuf:"varint,1,opt,name=accountId,proto3" json:"accountId,omitempty"`
	Amount        int32                  `protobuf:"varint,2,opt,name=amount,proto3" json:"amount,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DepositAccountRequest) Reset() {
	*x = DepositAccountRequest{}
	mi := &file_proto_account_account_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DepositAccountRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DepositAccountRequest) ProtoMessage() {}

func (x *DepositAccountRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_account_account_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DepositAccountRequest.ProtoReflect.Descriptor instead.
func (*DepositAccountRequest) Descriptor() ([]byte, []int) {
	return file_proto_account_account_proto_rawDescGZIP(), []int{3}
}

func (x *DepositAccountRequest) GetAccountId() int32 {
	if x != nil {
		return x.AccountId
	}
	return 0
}

func (x *DepositAccountRequest) GetAmount() int32 {
	if x != nil {
		return x.Amount
	}
	return 0
}

type DepositAccountResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Code          int32                  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Data          *AccountData           `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DepositAccountResponse) Reset() {
	*x = DepositAccountResponse{}
	mi := &file_proto_account_account_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DepositAccountResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DepositAccountResponse) ProtoMessage() {}

func (x *DepositAccountResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_account_account_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DepositAccountResponse.ProtoReflect.Descriptor instead.
func (*DepositAccountResponse) Descriptor() ([]byte, []int) {
	return file_proto_account_account_proto_rawDescGZIP(), []int{4}
}

func (x *DepositAccountResponse) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *DepositAccountResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *DepositAccountResponse) GetData() *AccountData {
	if x != nil {
		return x.Data
	}
	return nil
}

type BalanceChangeRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Ac            int32                  `protobuf:"varint,1,opt,name=ac,proto3" json:"ac,omitempty"`
	Tx            int64                  `protobuf:"varint,2,opt,name=tx,proto3" json:"tx,omitempty"`
	Am            int32                  `protobuf:"varint,3,opt,name=am,proto3" json:"am,omitempty"`
	Act           int32                  `protobuf:"varint,4,opt,name=act,proto3" json:"act,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BalanceChangeRequest) Reset() {
	*x = BalanceChangeRequest{}
	mi := &file_proto_account_account_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BalanceChangeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BalanceChangeRequest) ProtoMessage() {}

func (x *BalanceChangeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_account_account_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BalanceChangeRequest.ProtoReflect.Descriptor instead.
func (*BalanceChangeRequest) Descriptor() ([]byte, []int) {
	return file_proto_account_account_proto_rawDescGZIP(), []int{5}
}

func (x *BalanceChangeRequest) GetAc() int32 {
	if x != nil {
		return x.Ac
	}
	return 0
}

func (x *BalanceChangeRequest) GetTx() int64 {
	if x != nil {
		return x.Tx
	}
	return 0
}

func (x *BalanceChangeRequest) GetAm() int32 {
	if x != nil {
		return x.Am
	}
	return 0
}

func (x *BalanceChangeRequest) GetAct() int32 {
	if x != nil {
		return x.Act
	}
	return 0
}

type BalanceChangeResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Code          int32                  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Data          *AccountData           `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BalanceChangeResponse) Reset() {
	*x = BalanceChangeResponse{}
	mi := &file_proto_account_account_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BalanceChangeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BalanceChangeResponse) ProtoMessage() {}

func (x *BalanceChangeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_account_account_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BalanceChangeResponse.ProtoReflect.Descriptor instead.
func (*BalanceChangeResponse) Descriptor() ([]byte, []int) {
	return file_proto_account_account_proto_rawDescGZIP(), []int{6}
}

func (x *BalanceChangeResponse) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *BalanceChangeResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *BalanceChangeResponse) GetData() *AccountData {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_proto_account_account_proto protoreflect.FileDescriptor

var file_proto_account_account_proto_rawDesc = string([]byte{
	0x0a, 0x1b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2f,
	0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x61,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x34, 0x0a, 0x14, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c,
	0x0a, 0x09, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x09, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x22, 0x6f, 0x0a, 0x15,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x12, 0x28, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x14, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x41, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x5f, 0x0a,
	0x0b, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x44, 0x61, 0x74, 0x61, 0x12, 0x0e, 0x0a, 0x02,
	0x61, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x61, 0x63, 0x12, 0x0e, 0x0a, 0x02,
	0x62, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x62, 0x63, 0x12, 0x0e, 0x0a, 0x02,
	0x74, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x74, 0x78, 0x12, 0x10, 0x0a, 0x03,
	0x69, 0x74, 0x78, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x69, 0x74, 0x78, 0x12, 0x0e,
	0x0a, 0x02, 0x75, 0x70, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x75, 0x70, 0x22, 0x4d,
	0x0a, 0x15, 0x44, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x61, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x70, 0x0a,
	0x16, 0x44, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x28, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x41, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22,
	0x58, 0x0a, 0x14, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x61, 0x63, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x02, 0x61, 0x63, 0x12, 0x0e, 0x0a, 0x02, 0x74, 0x78, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x02, 0x74, 0x78, 0x12, 0x0e, 0x0a, 0x02, 0x61, 0x6d, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x02, 0x61, 0x6d, 0x12, 0x10, 0x0a, 0x03, 0x61, 0x63, 0x74, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x61, 0x63, 0x74, 0x22, 0x6f, 0x0a, 0x15, 0x42, 0x61, 0x6c,
	0x61, 0x6e, 0x63, 0x65, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x12, 0x28, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14,
	0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x2a, 0x1f, 0x0a, 0x06, 0x41, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0a, 0x0a, 0x06, 0x43, 0x52, 0x45, 0x44, 0x49, 0x54, 0x10, 0x00,
	0x12, 0x09, 0x0a, 0x05, 0x44, 0x45, 0x42, 0x49, 0x54, 0x10, 0x01, 0x32, 0xb0, 0x01, 0x0a, 0x0e,
	0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4e,
	0x0a, 0x0d, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12,
	0x1d, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e,
	0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4e,
	0x0a, 0x0d, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x12,
	0x1d, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63,
	0x65, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e,
	0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65,
	0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0f,
	0x5a, 0x0d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_proto_account_account_proto_rawDescOnce sync.Once
	file_proto_account_account_proto_rawDescData []byte
)

func file_proto_account_account_proto_rawDescGZIP() []byte {
	file_proto_account_account_proto_rawDescOnce.Do(func() {
		file_proto_account_account_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_account_account_proto_rawDesc), len(file_proto_account_account_proto_rawDesc)))
	})
	return file_proto_account_account_proto_rawDescData
}

var file_proto_account_account_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_account_account_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_proto_account_account_proto_goTypes = []any{
	(Action)(0),                    // 0: account.Action
	(*CreateAccountRequest)(nil),   // 1: account.CreateAccountRequest
	(*CreateAccountResponse)(nil),  // 2: account.CreateAccountResponse
	(*AccountData)(nil),            // 3: account.AccountData
	(*DepositAccountRequest)(nil),  // 4: account.DepositAccountRequest
	(*DepositAccountResponse)(nil), // 5: account.DepositAccountResponse
	(*BalanceChangeRequest)(nil),   // 6: account.BalanceChangeRequest
	(*BalanceChangeResponse)(nil),  // 7: account.BalanceChangeResponse
}
var file_proto_account_account_proto_depIdxs = []int32{
	3, // 0: account.CreateAccountResponse.data:type_name -> account.AccountData
	3, // 1: account.DepositAccountResponse.data:type_name -> account.AccountData
	3, // 2: account.BalanceChangeResponse.data:type_name -> account.AccountData
	1, // 3: account.AccountService.CreateAccount:input_type -> account.CreateAccountRequest
	6, // 4: account.AccountService.BalanceChange:input_type -> account.BalanceChangeRequest
	2, // 5: account.AccountService.CreateAccount:output_type -> account.CreateAccountResponse
	7, // 6: account.AccountService.BalanceChange:output_type -> account.BalanceChangeResponse
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_proto_account_account_proto_init() }
func file_proto_account_account_proto_init() {
	if File_proto_account_account_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_account_account_proto_rawDesc), len(file_proto_account_account_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_account_account_proto_goTypes,
		DependencyIndexes: file_proto_account_account_proto_depIdxs,
		EnumInfos:         file_proto_account_account_proto_enumTypes,
		MessageInfos:      file_proto_account_account_proto_msgTypes,
	}.Build()
	File_proto_account_account_proto = out.File
	file_proto_account_account_proto_goTypes = nil
	file_proto_account_account_proto_depIdxs = nil
}
