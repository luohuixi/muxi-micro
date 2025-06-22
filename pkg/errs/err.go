package errs

import (
	"errors"
	"fmt"
)

type Err struct {
	//错误的名称，方便打印日志的时候区分
	name string
	//错误的描述信息
	msg string
	//元数据，方便携带一些额外的信息
	metaData map[string]interface{}
	//为了支持errors.Is和errors.As
	cause error
}

func (e *Err) Error() string {
	msg := fmt.Sprintf("[%s] %s", e.name, e.msg)
	if e.metaData != nil {
		msg += fmt.Sprintf(" %v", e.metaData)
	}
	if e.cause != nil {
		msg += " => " + e.cause.Error()
	}
	return msg
}

func (e *Err) Name() string {
	return e.name
}
func (e *Err) Message() string {
	return e.msg
}
func (e *Err) MetaData() map[string]interface{} {
	return e.metaData
}
func (e *Err) Unwrap() error {
	return e.cause
}
func (e *Err) Is(target error) bool {
	var t *Err
	if errors.As(target, &t) {
		return e.name == t.name
	}
	return false
}

func (e *Err) WithCause(cause error) *Err {
	return &Err{
		name:     e.name,
		msg:      e.msg,
		metaData: e.metaData,
		cause:    cause,
	}
}

func (e *Err) WithMeta(meta map[string]interface{}) *Err {
	return withMetaData(e, meta) // 复用现有逻辑
}

func NewErr(name, msg string) *Err {
	return &Err{
		name: name,
		msg:  msg,
	}
}

func withMetaData(e *Err, meta map[string]interface{}) *Err {
	if e == nil {
		return nil
	}
	ne := &Err{
		name:     e.name,
		msg:      e.msg,
		metaData: meta,    // 仅更新元数据
		cause:    e.cause, // 保持原始原因链不变
	}
	return ne
}
