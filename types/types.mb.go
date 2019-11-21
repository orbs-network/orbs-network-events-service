// AUTO GENERATED FILE (by membufc proto compiler v0.4.0)
package types

import (
	"bytes"
	"context"
	"fmt"
	"github.com/orbs-network/membuffers/go"
)

/////////////////////////////////////////////////////////////////////////////
// service Indexer

type Indexer interface {
	GetEvents(ctx context.Context, input *IndexerRequest) (*IndexerResponse, error)
}

/////////////////////////////////////////////////////////////////////////////
// message IndexedEvent

// reader

type IndexedEvent struct {
	// ContractName string
	// EventName string
	// BlockHeight uint64
	// Timestamp uint64
	// Txhash []byte
	// ExecutionResult ExecutionResult
	// Index uint32
	// Arguments []byte

	// internal
	// implements membuffers.Message
	_message membuffers.InternalMessage
}

func (x *IndexedEvent) String() string {
	if x == nil {
		return "<nil>"
	}
	return fmt.Sprintf("{ContractName:%s,EventName:%s,BlockHeight:%s,Timestamp:%s,Txhash:%s,ExecutionResult:%s,Index:%s,Arguments:%s,}", x.StringContractName(), x.StringEventName(), x.StringBlockHeight(), x.StringTimestamp(), x.StringTxhash(), x.StringExecutionResult(), x.StringIndex(), x.StringArguments())
}

var _IndexedEvent_Scheme = []membuffers.FieldType{membuffers.TypeString,membuffers.TypeString,membuffers.TypeUint64,membuffers.TypeUint64,membuffers.TypeBytes,membuffers.TypeUint16,membuffers.TypeUint32,membuffers.TypeBytes,}
var _IndexedEvent_Unions = [][]membuffers.FieldType{}

func IndexedEventReader(buf []byte) *IndexedEvent {
	x := &IndexedEvent{}
	x._message.Init(buf, membuffers.Offset(len(buf)), _IndexedEvent_Scheme, _IndexedEvent_Unions)
	return x
}

func (x *IndexedEvent) IsValid() bool {
	return x._message.IsValid()
}

func (x *IndexedEvent) Raw() []byte {
	return x._message.RawBuffer()
}

func (x *IndexedEvent) Equal(y *IndexedEvent) bool {
  if x == nil && y == nil {
    return true
  }
  if x == nil || y == nil {
    return false
  }
  return bytes.Equal(x.Raw(), y.Raw())
}

func (x *IndexedEvent) ContractName() string {
	return x._message.GetString(0)
}

func (x *IndexedEvent) RawContractName() []byte {
	return x._message.RawBufferForField(0, 0)
}

func (x *IndexedEvent) RawContractNameWithHeader() []byte {
	return x._message.RawBufferWithHeaderForField(0, 0)
}

func (x *IndexedEvent) MutateContractName(v string) error {
	return x._message.SetString(0, v)
}

func (x *IndexedEvent) StringContractName() string {
return fmt.Sprintf("%s", x.ContractName())}

func (x *IndexedEvent) EventName() string {
	return x._message.GetString(1)
}

func (x *IndexedEvent) RawEventName() []byte {
	return x._message.RawBufferForField(1, 0)
}

func (x *IndexedEvent) RawEventNameWithHeader() []byte {
	return x._message.RawBufferWithHeaderForField(1, 0)
}

func (x *IndexedEvent) MutateEventName(v string) error {
	return x._message.SetString(1, v)
}

func (x *IndexedEvent) StringEventName() string {
return fmt.Sprintf("%s", x.EventName())}

func (x *IndexedEvent) BlockHeight() uint64 {
	return x._message.GetUint64(2)
}

func (x *IndexedEvent) RawBlockHeight() []byte {
	return x._message.RawBufferForField(2, 0)
}

func (x *IndexedEvent) MutateBlockHeight(v uint64) error {
	return x._message.SetUint64(2, v)
}

func (x *IndexedEvent) StringBlockHeight() string {
return fmt.Sprintf("%v", x.BlockHeight())}

func (x *IndexedEvent) Timestamp() uint64 {
	return x._message.GetUint64(3)
}

func (x *IndexedEvent) RawTimestamp() []byte {
	return x._message.RawBufferForField(3, 0)
}

func (x *IndexedEvent) MutateTimestamp(v uint64) error {
	return x._message.SetUint64(3, v)
}

func (x *IndexedEvent) StringTimestamp() string {
return fmt.Sprintf("%v", x.Timestamp())}

func (x *IndexedEvent) Txhash() []byte {
	return x._message.GetBytes(4)
}

func (x *IndexedEvent) RawTxhash() []byte {
	return x._message.RawBufferForField(4, 0)
}

func (x *IndexedEvent) RawTxhashWithHeader() []byte {
	return x._message.RawBufferWithHeaderForField(4, 0)
}

func (x *IndexedEvent) MutateTxhash(v []byte) error {
	return x._message.SetBytes(4, v)
}

func (x *IndexedEvent) StringTxhash() string {
return fmt.Sprintf("%x", x.Txhash())}

func (x *IndexedEvent) ExecutionResult() ExecutionResult {
	return ExecutionResult(x._message.GetUint16(5))
}

func (x *IndexedEvent) RawExecutionResult() []byte {
	return x._message.RawBufferForField(5, 0)
}

func (x *IndexedEvent) MutateExecutionResult(v ExecutionResult) error {
	return x._message.SetUint16(5, uint16(v))
}

func (x *IndexedEvent) StringExecutionResult() string {
return x.ExecutionResult().String()}

func (x *IndexedEvent) Index() uint32 {
	return x._message.GetUint32(6)
}

func (x *IndexedEvent) RawIndex() []byte {
	return x._message.RawBufferForField(6, 0)
}

func (x *IndexedEvent) MutateIndex(v uint32) error {
	return x._message.SetUint32(6, v)
}

func (x *IndexedEvent) StringIndex() string {
return fmt.Sprintf("%v", x.Index())}

func (x *IndexedEvent) Arguments() []byte {
	return x._message.GetBytes(7)
}

func (x *IndexedEvent) RawArguments() []byte {
	return x._message.RawBufferForField(7, 0)
}

func (x *IndexedEvent) RawArgumentsWithHeader() []byte {
	return x._message.RawBufferWithHeaderForField(7, 0)
}

func (x *IndexedEvent) MutateArguments(v []byte) error {
	return x._message.SetBytes(7, v)
}

func (x *IndexedEvent) StringArguments() string {
return fmt.Sprintf("%x", x.Arguments())}

// builder

type IndexedEventBuilder struct {
	ContractName string
	EventName string
	BlockHeight uint64
	Timestamp uint64
	Txhash []byte
	ExecutionResult ExecutionResult
	Index uint32
	Arguments []byte

	// internal
	// implements membuffers.Builder
	_builder membuffers.InternalBuilder
	_overrideWithRawBuffer []byte
}

func (w *IndexedEventBuilder) Write(buf []byte) (err error) {
	if w == nil {
		return
	}
	w._builder.NotifyBuildStart()
	defer w._builder.NotifyBuildEnd()
	defer func() {
		if r := recover(); r != nil {
			err = &membuffers.ErrBufferOverrun{}
		}
	}()
	if w._overrideWithRawBuffer != nil {
		return w._builder.WriteOverrideWithRawBuffer(buf, w._overrideWithRawBuffer)
	}
	w._builder.Reset()
	w._builder.WriteString(buf, w.ContractName)
	w._builder.WriteString(buf, w.EventName)
	w._builder.WriteUint64(buf, w.BlockHeight)
	w._builder.WriteUint64(buf, w.Timestamp)
	w._builder.WriteBytes(buf, w.Txhash)
	w._builder.WriteUint16(buf, uint16(w.ExecutionResult))
	w._builder.WriteUint32(buf, w.Index)
	w._builder.WriteBytes(buf, w.Arguments)
	return nil
}

func (w *IndexedEventBuilder) HexDump(prefix string, offsetFromStart membuffers.Offset) (err error) {
	if w == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			err = &membuffers.ErrBufferOverrun{}
		}
	}()
	w._builder.Reset()
	w._builder.HexDumpString(prefix, offsetFromStart, "IndexedEvent.ContractName", w.ContractName)
	w._builder.HexDumpString(prefix, offsetFromStart, "IndexedEvent.EventName", w.EventName)
	w._builder.HexDumpUint64(prefix, offsetFromStart, "IndexedEvent.BlockHeight", w.BlockHeight)
	w._builder.HexDumpUint64(prefix, offsetFromStart, "IndexedEvent.Timestamp", w.Timestamp)
	w._builder.HexDumpBytes(prefix, offsetFromStart, "IndexedEvent.Txhash", w.Txhash)
	w._builder.HexDumpUint16(prefix, offsetFromStart, "IndexedEvent.ExecutionResult", uint16(w.ExecutionResult))
	w._builder.HexDumpUint32(prefix, offsetFromStart, "IndexedEvent.Index", w.Index)
	w._builder.HexDumpBytes(prefix, offsetFromStart, "IndexedEvent.Arguments", w.Arguments)
	return nil
}

func (w *IndexedEventBuilder) GetSize() membuffers.Offset {
	if w == nil {
		return 0
	}
	return w._builder.GetSize()
}

func (w *IndexedEventBuilder) CalcRequiredSize() membuffers.Offset {
	if w == nil {
		return 0
	}
	w.Write(nil)
	return w._builder.GetSize()
}

func (w *IndexedEventBuilder) Build() *IndexedEvent {
	buf := make([]byte, w.CalcRequiredSize())
	if w.Write(buf) != nil {
		return nil
	}
	return IndexedEventReader(buf)
}

func IndexedEventBuilderFromRaw(raw []byte) *IndexedEventBuilder {
	return &IndexedEventBuilder{_overrideWithRawBuffer: raw}
}

/////////////////////////////////////////////////////////////////////////////
// message Filter

// reader

type Filter struct {
	// Argument []string

	// internal
	// implements membuffers.Message
	_message membuffers.InternalMessage
}

func (x *Filter) String() string {
	if x == nil {
		return "<nil>"
	}
	return fmt.Sprintf("{Argument:%s,}", x.StringArgument())
}

var _Filter_Scheme = []membuffers.FieldType{membuffers.TypeStringArray,}
var _Filter_Unions = [][]membuffers.FieldType{}

func FilterReader(buf []byte) *Filter {
	x := &Filter{}
	x._message.Init(buf, membuffers.Offset(len(buf)), _Filter_Scheme, _Filter_Unions)
	return x
}

func (x *Filter) IsValid() bool {
	return x._message.IsValid()
}

func (x *Filter) Raw() []byte {
	return x._message.RawBuffer()
}

func (x *Filter) Equal(y *Filter) bool {
  if x == nil && y == nil {
    return true
  }
  if x == nil || y == nil {
    return false
  }
  return bytes.Equal(x.Raw(), y.Raw())
}

func (x *Filter) ArgumentIterator() *FilterArgumentIterator {
	return &FilterArgumentIterator{iterator: x._message.GetStringArrayIterator(0)}
}

type FilterArgumentIterator struct {
	iterator *membuffers.Iterator
}

func (i *FilterArgumentIterator) HasNext() bool {
	return i.iterator.HasNext()
}

func (i *FilterArgumentIterator) NextArgument() string {
	return i.iterator.NextString()
}

func (x *Filter) RawArgumentArray() []byte {
	return x._message.RawBufferForField(0, 0)
}

func (x *Filter) RawArgumentArrayWithHeader() []byte {
	return x._message.RawBufferWithHeaderForField(0, 0)
}

func (x *Filter) StringArgument() (res string) {
	res = "["
	for i := x.ArgumentIterator(); i.HasNext(); {
		res += fmt.Sprintf(i.NextArgument()) + ","
	}
	res += "]"
	return
}

// builder

type FilterBuilder struct {
	Argument []string

	// internal
	// implements membuffers.Builder
	_builder membuffers.InternalBuilder
	_overrideWithRawBuffer []byte
}

func (w *FilterBuilder) Write(buf []byte) (err error) {
	if w == nil {
		return
	}
	w._builder.NotifyBuildStart()
	defer w._builder.NotifyBuildEnd()
	defer func() {
		if r := recover(); r != nil {
			err = &membuffers.ErrBufferOverrun{}
		}
	}()
	if w._overrideWithRawBuffer != nil {
		return w._builder.WriteOverrideWithRawBuffer(buf, w._overrideWithRawBuffer)
	}
	w._builder.Reset()
	w._builder.WriteStringArray(buf, w.Argument)
	return nil
}

func (w *FilterBuilder) HexDump(prefix string, offsetFromStart membuffers.Offset) (err error) {
	if w == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			err = &membuffers.ErrBufferOverrun{}
		}
	}()
	w._builder.Reset()
	w._builder.HexDumpStringArray(prefix, offsetFromStart, "Filter.Argument", w.Argument)
	return nil
}

func (w *FilterBuilder) GetSize() membuffers.Offset {
	if w == nil {
		return 0
	}
	return w._builder.GetSize()
}

func (w *FilterBuilder) CalcRequiredSize() membuffers.Offset {
	if w == nil {
		return 0
	}
	w.Write(nil)
	return w._builder.GetSize()
}

func (w *FilterBuilder) Build() *Filter {
	buf := make([]byte, w.CalcRequiredSize())
	if w.Write(buf) != nil {
		return nil
	}
	return FilterReader(buf)
}

func FilterBuilderFromRaw(raw []byte) *FilterBuilder {
	return &FilterBuilder{_overrideWithRawBuffer: raw}
}

/////////////////////////////////////////////////////////////////////////////
// message IndexerRequest

// reader

type IndexerRequest struct {
	// ProtocolVersion uint32
	// VirtualChainId uint32
	// ContractName string
	// EventName []string
	// FromBlock uint64
	// ToBlock uint64
	// FromTime uint64
	// ToTime uint64
	// Filters []Filter

	// internal
	// implements membuffers.Message
	_message membuffers.InternalMessage
}

func (x *IndexerRequest) String() string {
	if x == nil {
		return "<nil>"
	}
	return fmt.Sprintf("{ProtocolVersion:%s,VirtualChainId:%s,ContractName:%s,EventName:%s,FromBlock:%s,ToBlock:%s,FromTime:%s,ToTime:%s,Filters:%s,}", x.StringProtocolVersion(), x.StringVirtualChainId(), x.StringContractName(), x.StringEventName(), x.StringFromBlock(), x.StringToBlock(), x.StringFromTime(), x.StringToTime(), x.StringFilters())
}

var _IndexerRequest_Scheme = []membuffers.FieldType{membuffers.TypeUint32,membuffers.TypeUint32,membuffers.TypeString,membuffers.TypeStringArray,membuffers.TypeUint64,membuffers.TypeUint64,membuffers.TypeUint64,membuffers.TypeUint64,membuffers.TypeMessageArray,}
var _IndexerRequest_Unions = [][]membuffers.FieldType{}

func IndexerRequestReader(buf []byte) *IndexerRequest {
	x := &IndexerRequest{}
	x._message.Init(buf, membuffers.Offset(len(buf)), _IndexerRequest_Scheme, _IndexerRequest_Unions)
	return x
}

func (x *IndexerRequest) IsValid() bool {
	return x._message.IsValid()
}

func (x *IndexerRequest) Raw() []byte {
	return x._message.RawBuffer()
}

func (x *IndexerRequest) Equal(y *IndexerRequest) bool {
  if x == nil && y == nil {
    return true
  }
  if x == nil || y == nil {
    return false
  }
  return bytes.Equal(x.Raw(), y.Raw())
}

func (x *IndexerRequest) ProtocolVersion() uint32 {
	return x._message.GetUint32(0)
}

func (x *IndexerRequest) RawProtocolVersion() []byte {
	return x._message.RawBufferForField(0, 0)
}

func (x *IndexerRequest) MutateProtocolVersion(v uint32) error {
	return x._message.SetUint32(0, v)
}

func (x *IndexerRequest) StringProtocolVersion() string {
return fmt.Sprintf("%v", x.ProtocolVersion())}

func (x *IndexerRequest) VirtualChainId() uint32 {
	return x._message.GetUint32(1)
}

func (x *IndexerRequest) RawVirtualChainId() []byte {
	return x._message.RawBufferForField(1, 0)
}

func (x *IndexerRequest) MutateVirtualChainId(v uint32) error {
	return x._message.SetUint32(1, v)
}

func (x *IndexerRequest) StringVirtualChainId() string {
return fmt.Sprintf("%v", x.VirtualChainId())}

func (x *IndexerRequest) ContractName() string {
	return x._message.GetString(2)
}

func (x *IndexerRequest) RawContractName() []byte {
	return x._message.RawBufferForField(2, 0)
}

func (x *IndexerRequest) RawContractNameWithHeader() []byte {
	return x._message.RawBufferWithHeaderForField(2, 0)
}

func (x *IndexerRequest) MutateContractName(v string) error {
	return x._message.SetString(2, v)
}

func (x *IndexerRequest) StringContractName() string {
return fmt.Sprintf("%s", x.ContractName())}

func (x *IndexerRequest) EventNameIterator() *IndexerRequestEventNameIterator {
	return &IndexerRequestEventNameIterator{iterator: x._message.GetStringArrayIterator(3)}
}

type IndexerRequestEventNameIterator struct {
	iterator *membuffers.Iterator
}

func (i *IndexerRequestEventNameIterator) HasNext() bool {
	return i.iterator.HasNext()
}

func (i *IndexerRequestEventNameIterator) NextEventName() string {
	return i.iterator.NextString()
}

func (x *IndexerRequest) RawEventNameArray() []byte {
	return x._message.RawBufferForField(3, 0)
}

func (x *IndexerRequest) RawEventNameArrayWithHeader() []byte {
	return x._message.RawBufferWithHeaderForField(3, 0)
}

func (x *IndexerRequest) StringEventName() (res string) {
	res = "["
	for i := x.EventNameIterator(); i.HasNext(); {
		res += fmt.Sprintf(i.NextEventName()) + ","
	}
	res += "]"
	return
}

func (x *IndexerRequest) FromBlock() uint64 {
	return x._message.GetUint64(4)
}

func (x *IndexerRequest) RawFromBlock() []byte {
	return x._message.RawBufferForField(4, 0)
}

func (x *IndexerRequest) MutateFromBlock(v uint64) error {
	return x._message.SetUint64(4, v)
}

func (x *IndexerRequest) StringFromBlock() string {
return fmt.Sprintf("%v", x.FromBlock())}

func (x *IndexerRequest) ToBlock() uint64 {
	return x._message.GetUint64(5)
}

func (x *IndexerRequest) RawToBlock() []byte {
	return x._message.RawBufferForField(5, 0)
}

func (x *IndexerRequest) MutateToBlock(v uint64) error {
	return x._message.SetUint64(5, v)
}

func (x *IndexerRequest) StringToBlock() string {
return fmt.Sprintf("%v", x.ToBlock())}

func (x *IndexerRequest) FromTime() uint64 {
	return x._message.GetUint64(6)
}

func (x *IndexerRequest) RawFromTime() []byte {
	return x._message.RawBufferForField(6, 0)
}

func (x *IndexerRequest) MutateFromTime(v uint64) error {
	return x._message.SetUint64(6, v)
}

func (x *IndexerRequest) StringFromTime() string {
return fmt.Sprintf("%v", x.FromTime())}

func (x *IndexerRequest) ToTime() uint64 {
	return x._message.GetUint64(7)
}

func (x *IndexerRequest) RawToTime() []byte {
	return x._message.RawBufferForField(7, 0)
}

func (x *IndexerRequest) MutateToTime(v uint64) error {
	return x._message.SetUint64(7, v)
}

func (x *IndexerRequest) StringToTime() string {
return fmt.Sprintf("%v", x.ToTime())}

func (x *IndexerRequest) FiltersIterator() *IndexerRequestFiltersIterator {
	return &IndexerRequestFiltersIterator{iterator: x._message.GetMessageArrayIterator(8)}
}

type IndexerRequestFiltersIterator struct {
	iterator *membuffers.Iterator
}

func (i *IndexerRequestFiltersIterator) HasNext() bool {
	return i.iterator.HasNext()
}

func (i *IndexerRequestFiltersIterator) NextFilters() *Filter {
	b, s := i.iterator.NextMessage()
	return FilterReader(b[:s])
}

func (x *IndexerRequest) RawFiltersArray() []byte {
	return x._message.RawBufferForField(8, 0)
}

func (x *IndexerRequest) RawFiltersArrayWithHeader() []byte {
	return x._message.RawBufferWithHeaderForField(8, 0)
}

func (x *IndexerRequest) StringFilters() (res string) {
	res = "["
	for i := x.FiltersIterator(); i.HasNext(); {
		res += i.NextFilters().String() + ","
	}
	res += "]"
	return
}

// builder

type IndexerRequestBuilder struct {
	ProtocolVersion uint32
	VirtualChainId uint32
	ContractName string
	EventName []string
	FromBlock uint64
	ToBlock uint64
	FromTime uint64
	ToTime uint64
	Filters []*FilterBuilder

	// internal
	// implements membuffers.Builder
	_builder membuffers.InternalBuilder
	_overrideWithRawBuffer []byte
}

func (w *IndexerRequestBuilder) arrayOfFilters() []membuffers.MessageWriter {
	res := make([]membuffers.MessageWriter, len(w.Filters))
	for i, v := range w.Filters {
		res[i] = v
	}
	return res
}

func (w *IndexerRequestBuilder) Write(buf []byte) (err error) {
	if w == nil {
		return
	}
	w._builder.NotifyBuildStart()
	defer w._builder.NotifyBuildEnd()
	defer func() {
		if r := recover(); r != nil {
			err = &membuffers.ErrBufferOverrun{}
		}
	}()
	if w._overrideWithRawBuffer != nil {
		return w._builder.WriteOverrideWithRawBuffer(buf, w._overrideWithRawBuffer)
	}
	w._builder.Reset()
	w._builder.WriteUint32(buf, w.ProtocolVersion)
	w._builder.WriteUint32(buf, w.VirtualChainId)
	w._builder.WriteString(buf, w.ContractName)
	w._builder.WriteStringArray(buf, w.EventName)
	w._builder.WriteUint64(buf, w.FromBlock)
	w._builder.WriteUint64(buf, w.ToBlock)
	w._builder.WriteUint64(buf, w.FromTime)
	w._builder.WriteUint64(buf, w.ToTime)
	err = w._builder.WriteMessageArray(buf, w.arrayOfFilters())
	if err != nil {
		return
	}
	return nil
}

func (w *IndexerRequestBuilder) HexDump(prefix string, offsetFromStart membuffers.Offset) (err error) {
	if w == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			err = &membuffers.ErrBufferOverrun{}
		}
	}()
	w._builder.Reset()
	w._builder.HexDumpUint32(prefix, offsetFromStart, "IndexerRequest.ProtocolVersion", w.ProtocolVersion)
	w._builder.HexDumpUint32(prefix, offsetFromStart, "IndexerRequest.VirtualChainId", w.VirtualChainId)
	w._builder.HexDumpString(prefix, offsetFromStart, "IndexerRequest.ContractName", w.ContractName)
	w._builder.HexDumpStringArray(prefix, offsetFromStart, "IndexerRequest.EventName", w.EventName)
	w._builder.HexDumpUint64(prefix, offsetFromStart, "IndexerRequest.FromBlock", w.FromBlock)
	w._builder.HexDumpUint64(prefix, offsetFromStart, "IndexerRequest.ToBlock", w.ToBlock)
	w._builder.HexDumpUint64(prefix, offsetFromStart, "IndexerRequest.FromTime", w.FromTime)
	w._builder.HexDumpUint64(prefix, offsetFromStart, "IndexerRequest.ToTime", w.ToTime)
	err = w._builder.HexDumpMessageArray(prefix, offsetFromStart, "IndexerRequest.Filters", w.arrayOfFilters())
	if err != nil {
		return
	}
	return nil
}

func (w *IndexerRequestBuilder) GetSize() membuffers.Offset {
	if w == nil {
		return 0
	}
	return w._builder.GetSize()
}

func (w *IndexerRequestBuilder) CalcRequiredSize() membuffers.Offset {
	if w == nil {
		return 0
	}
	w.Write(nil)
	return w._builder.GetSize()
}

func (w *IndexerRequestBuilder) Build() *IndexerRequest {
	buf := make([]byte, w.CalcRequiredSize())
	if w.Write(buf) != nil {
		return nil
	}
	return IndexerRequestReader(buf)
}

func IndexerRequestBuilderFromRaw(raw []byte) *IndexerRequestBuilder {
	return &IndexerRequestBuilder{_overrideWithRawBuffer: raw}
}

/////////////////////////////////////////////////////////////////////////////
// message IndexerResponse

// reader

type IndexerResponse struct {
	// Events []IndexedEvent

	// internal
	// implements membuffers.Message
	_message membuffers.InternalMessage
}

func (x *IndexerResponse) String() string {
	if x == nil {
		return "<nil>"
	}
	return fmt.Sprintf("{Events:%s,}", x.StringEvents())
}

var _IndexerResponse_Scheme = []membuffers.FieldType{membuffers.TypeMessageArray,}
var _IndexerResponse_Unions = [][]membuffers.FieldType{}

func IndexerResponseReader(buf []byte) *IndexerResponse {
	x := &IndexerResponse{}
	x._message.Init(buf, membuffers.Offset(len(buf)), _IndexerResponse_Scheme, _IndexerResponse_Unions)
	return x
}

func (x *IndexerResponse) IsValid() bool {
	return x._message.IsValid()
}

func (x *IndexerResponse) Raw() []byte {
	return x._message.RawBuffer()
}

func (x *IndexerResponse) Equal(y *IndexerResponse) bool {
  if x == nil && y == nil {
    return true
  }
  if x == nil || y == nil {
    return false
  }
  return bytes.Equal(x.Raw(), y.Raw())
}

func (x *IndexerResponse) EventsIterator() *IndexerResponseEventsIterator {
	return &IndexerResponseEventsIterator{iterator: x._message.GetMessageArrayIterator(0)}
}

type IndexerResponseEventsIterator struct {
	iterator *membuffers.Iterator
}

func (i *IndexerResponseEventsIterator) HasNext() bool {
	return i.iterator.HasNext()
}

func (i *IndexerResponseEventsIterator) NextEvents() *IndexedEvent {
	b, s := i.iterator.NextMessage()
	return IndexedEventReader(b[:s])
}

func (x *IndexerResponse) RawEventsArray() []byte {
	return x._message.RawBufferForField(0, 0)
}

func (x *IndexerResponse) RawEventsArrayWithHeader() []byte {
	return x._message.RawBufferWithHeaderForField(0, 0)
}

func (x *IndexerResponse) StringEvents() (res string) {
	res = "["
	for i := x.EventsIterator(); i.HasNext(); {
		res += i.NextEvents().String() + ","
	}
	res += "]"
	return
}

// builder

type IndexerResponseBuilder struct {
	Events []*IndexedEventBuilder

	// internal
	// implements membuffers.Builder
	_builder membuffers.InternalBuilder
	_overrideWithRawBuffer []byte
}

func (w *IndexerResponseBuilder) arrayOfEvents() []membuffers.MessageWriter {
	res := make([]membuffers.MessageWriter, len(w.Events))
	for i, v := range w.Events {
		res[i] = v
	}
	return res
}

func (w *IndexerResponseBuilder) Write(buf []byte) (err error) {
	if w == nil {
		return
	}
	w._builder.NotifyBuildStart()
	defer w._builder.NotifyBuildEnd()
	defer func() {
		if r := recover(); r != nil {
			err = &membuffers.ErrBufferOverrun{}
		}
	}()
	if w._overrideWithRawBuffer != nil {
		return w._builder.WriteOverrideWithRawBuffer(buf, w._overrideWithRawBuffer)
	}
	w._builder.Reset()
	err = w._builder.WriteMessageArray(buf, w.arrayOfEvents())
	if err != nil {
		return
	}
	return nil
}

func (w *IndexerResponseBuilder) HexDump(prefix string, offsetFromStart membuffers.Offset) (err error) {
	if w == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			err = &membuffers.ErrBufferOverrun{}
		}
	}()
	w._builder.Reset()
	err = w._builder.HexDumpMessageArray(prefix, offsetFromStart, "IndexerResponse.Events", w.arrayOfEvents())
	if err != nil {
		return
	}
	return nil
}

func (w *IndexerResponseBuilder) GetSize() membuffers.Offset {
	if w == nil {
		return 0
	}
	return w._builder.GetSize()
}

func (w *IndexerResponseBuilder) CalcRequiredSize() membuffers.Offset {
	if w == nil {
		return 0
	}
	w.Write(nil)
	return w._builder.GetSize()
}

func (w *IndexerResponseBuilder) Build() *IndexerResponse {
	buf := make([]byte, w.CalcRequiredSize())
	if w.Write(buf) != nil {
		return nil
	}
	return IndexerResponseReader(buf)
}

func IndexerResponseBuilderFromRaw(raw []byte) *IndexerResponseBuilder {
	return &IndexerResponseBuilder{_overrideWithRawBuffer: raw}
}

/////////////////////////////////////////////////////////////////////////////
// enums

type ExecutionResult uint16

const (
	EXECUTION_RESULT_RESERVED ExecutionResult = 0
	EXECUTION_RESULT_SUCCESS ExecutionResult = 1
	EXECUTION_RESULT_ERROR_SMART_CONTRACT ExecutionResult = 2
	EXECUTION_RESULT_ERROR_INPUT ExecutionResult = 3
	EXECUTION_RESULT_ERROR_CONTRACT_NOT_DEPLOYED ExecutionResult = 4
	EXECUTION_RESULT_ERROR_UNEXPECTED ExecutionResult = 5
	EXECUTION_RESULT_NOT_EXECUTED ExecutionResult = 6
)

func (n ExecutionResult) String() string {
	switch n {
	case EXECUTION_RESULT_RESERVED:
		return "EXECUTION_RESULT_RESERVED"
	case EXECUTION_RESULT_SUCCESS:
		return "EXECUTION_RESULT_SUCCESS"
	case EXECUTION_RESULT_ERROR_SMART_CONTRACT:
		return "EXECUTION_RESULT_ERROR_SMART_CONTRACT"
	case EXECUTION_RESULT_ERROR_INPUT:
		return "EXECUTION_RESULT_ERROR_INPUT"
	case EXECUTION_RESULT_ERROR_CONTRACT_NOT_DEPLOYED:
		return "EXECUTION_RESULT_ERROR_CONTRACT_NOT_DEPLOYED"
	case EXECUTION_RESULT_ERROR_UNEXPECTED:
		return "EXECUTION_RESULT_ERROR_UNEXPECTED"
	case EXECUTION_RESULT_NOT_EXECUTED:
		return "EXECUTION_RESULT_NOT_EXECUTED"
	}
	return "UNKNOWN"
}

