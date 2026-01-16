package logger

import (
	"bytes"
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

func getGID() uint64 {
	b := make([]byte, 64)
	n := runtime.Stack(b, false)
	b = b[:n]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	id, _ := strconv.ParseUint(string(b), 10, 64)
	return id
}

// custom encoder (non-structured, manually concatenate strings)
type formatterEncoder struct {
	buf    *bytes.Buffer
	fields []zapcore.Field // store structured fields
}

func NewFormatterEncoder() zapcore.Encoder {
	return &formatterEncoder{
		buf:    &bytes.Buffer{},
		fields: make([]zapcore.Field, 0),
	}
}

func (e *formatterEncoder) Clone() zapcore.Encoder {
	return &formatterEncoder{
		buf:    bytes.NewBuffer(e.buf.Bytes()),
		fields: append([]zapcore.Field{}, e.fields...),
	}
}

func (e *formatterEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	e.buf.Reset()

	timestamp := entry.Time.Format("2006-01-02 15:04:05.000")
	level := entry.Level.CapitalString()
	gid := getGID()

	// extract file name and function name
	file := "???"
	line := 0
	funcName := "???"
	if entry.Caller.Defined {
		file = filepath.Base(entry.Caller.File)
		line = entry.Caller.Line
		if fn := runtime.FuncForPC(entry.Caller.PC); fn != nil {
			parts := strings.Split(fn.Name(), "/")
			funcName = parts[len(parts)-1]
		}
	}

	msg := fmt.Sprintf(
		"%s[%s][gid-%d]%s[%s][-]%s:%d %s",
		timestamp,
		DefaultAppsName,
		gid,
		level,
		funcName,
		file,
		line,
		entry.Message,
	)

	// merge fields
	allFields := append(e.fields, fields...)

	// add fields
	for _, field := range allFields {
		msg += fmt.Sprintf(" %s=%v", field.Key, e.formatFieldValue(field))
	}

	e.buf.WriteString(msg)
	e.buf.WriteString("\n")

	// convert to zap buffer
	pool := buffer.NewPool()
	zapBuf := pool.Get()
	zapBuf.AppendString(e.buf.String())

	return zapBuf, nil
}

// format field value
func (e *formatterEncoder) formatFieldValue(field zapcore.Field) interface{} {
	switch field.Type {
	case zapcore.StringType:
		return field.String
	case zapcore.Int64Type, zapcore.Int32Type, zapcore.Int16Type, zapcore.Int8Type:
		return field.Integer
	case zapcore.Uint64Type, zapcore.Uint32Type, zapcore.Uint16Type, zapcore.Uint8Type, zapcore.UintptrType:
		return uint64(field.Integer)
	case zapcore.Float64Type, zapcore.Float32Type:
		return field.Interface
	case zapcore.BoolType:
		return field.Integer == 1
	case zapcore.TimeType:
		return time.Unix(0, field.Integer)
	case zapcore.DurationType:
		return time.Duration(field.Integer)
	default:
		return field.Interface
	}
}

func (e *formatterEncoder) addField(field zapcore.Field) {
	e.fields = append(e.fields, field)
}

func (e *formatterEncoder) AddArray(key string, marshaler zapcore.ArrayMarshaler) error {
	e.addField(zapcore.Field{Key: key, Type: zapcore.ArrayMarshalerType, Interface: marshaler})
	return nil
}

func (e *formatterEncoder) AddObject(key string, marshaler zapcore.ObjectMarshaler) error {
	e.addField(zapcore.Field{Key: key, Type: zapcore.ObjectMarshalerType, Interface: marshaler})
	return nil
}

func (e *formatterEncoder) AddBinary(key string, value []byte) {
	e.addField(zapcore.Field{Key: key, Type: zapcore.BinaryType, Interface: value})
}

func (e *formatterEncoder) AddByteString(key string, value []byte) {
	e.addField(zapcore.Field{Key: key, Type: zapcore.ByteStringType, Interface: value})
}

func (e *formatterEncoder) AddBool(key string, value bool) {
	var intVal int64
	if value {
		intVal = 1
	}
	e.addField(zapcore.Field{Key: key, Type: zapcore.BoolType, Integer: intVal})
}

func (e *formatterEncoder) AddComplex128(key string, value complex128) {
	e.addField(zapcore.Field{Key: key, Type: zapcore.Complex128Type, Interface: value})
}

func (e *formatterEncoder) AddComplex64(key string, value complex64) {
	e.addField(zapcore.Field{Key: key, Type: zapcore.Complex64Type, Interface: value})
}

func (e *formatterEncoder) AddDuration(key string, value time.Duration) {
	e.addField(zapcore.Field{Key: key, Type: zapcore.DurationType, Integer: int64(value)})
}

func (e *formatterEncoder) AddFloat64(key string, value float64) {
	e.addField(zapcore.Field{Key: key, Type: zapcore.Float64Type, Interface: value})
}

func (e *formatterEncoder) AddFloat32(key string, value float32) {
	e.addField(zapcore.Field{Key: key, Type: zapcore.Float32Type, Interface: value})
}

func (e *formatterEncoder) AddInt(key string, value int) {
	e.addField(zapcore.Field{Key: key, Type: zapcore.Int64Type, Integer: int64(value)})
}

func (e *formatterEncoder) AddInt64(key string, value int64) {
	e.addField(zapcore.Field{Key: key, Type: zapcore.Int64Type, Integer: value})
}

func (e *formatterEncoder) AddInt32(key string, value int32) {
	e.addField(zapcore.Field{Key: key, Type: zapcore.Int32Type, Integer: int64(value)})
}

func (e *formatterEncoder) AddInt16(key string, value int16) {
	e.addField(zapcore.Field{Key: key, Type: zapcore.Int16Type, Integer: int64(value)})
}

func (e *formatterEncoder) AddInt8(key string, value int8) {
	e.addField(zapcore.Field{Key: key, Type: zapcore.Int8Type, Integer: int64(value)})
}

func (e *formatterEncoder) AddString(key, value string) {
	e.addField(zapcore.Field{Key: key, Type: zapcore.StringType, String: value})
}

func (e *formatterEncoder) AddTime(key string, value time.Time) {
	e.addField(zapcore.Field{Key: key, Type: zapcore.TimeType, Integer: value.UnixNano()})
}

func (e *formatterEncoder) AddUint(key string, value uint) {
	e.addField(zapcore.Field{Key: key, Type: zapcore.Uint64Type, Integer: int64(value)})
}

func (e *formatterEncoder) AddUint64(key string, value uint64) {
	e.addField(zapcore.Field{Key: key, Type: zapcore.Uint64Type, Integer: int64(value)})
}

func (e *formatterEncoder) AddUint32(key string, value uint32) {
	e.addField(zapcore.Field{Key: key, Type: zapcore.Uint32Type, Integer: int64(value)})
}

func (e *formatterEncoder) AddUint16(key string, value uint16) {
	e.addField(zapcore.Field{Key: key, Type: zapcore.Uint16Type, Integer: int64(value)})
}

func (e *formatterEncoder) AddUint8(key string, value uint8) {
	e.addField(zapcore.Field{Key: key, Type: zapcore.Uint8Type, Integer: int64(value)})
}

func (e *formatterEncoder) AddUintptr(key string, value uintptr) {
	e.addField(zapcore.Field{Key: key, Type: zapcore.UintptrType, Integer: int64(value)})
}

func (e *formatterEncoder) AddReflected(key string, value interface{}) error {
	e.addField(zapcore.Field{Key: key, Type: zapcore.ReflectType, Interface: value})
	return nil
}

func (e *formatterEncoder) OpenNamespace(key string) {
	e.addField(zapcore.Field{Key: key, Type: zapcore.NamespaceType})
}
