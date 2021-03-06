package err1

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/vmihailenco/msgpack"
)

//Error 错误接口
type Error interface {
	json.Marshaler
	Code() int64    //自定义错误编码
	Msg() string    //自定义错误消息
	Err() error     //具体的错误
	Caller() string //返回调用堆栈信息
	Error() string  //继承全局的error接口
}

var msgregister = false

func MsgPackRegister(id int8) {
	if msgregister {
		return
	}
	msgpack.RegisterExt(id, &err{})
	msgregister = true
}

//err 公用错误对象
type err struct {
	code int64
	msg  string
	e    error
}

//Error 错误
func (e *err) Error() string {
	if e.e != nil {
		return e.e.Error()
	}
	return e.msg
}

func (e *err) Caller() string {
	return CallInfo(3)
}

//Code 错误代码
func (e *err) Code() int64 {
	return e.code
}

//Msg 错误消息
func (e *err) Msg() string {
	return e.msg
}

//Err 原始错误
func (e *err) Err() error {
	if e.e == nil {
		return nil
	}
	return e.e
}

func (e *err) MarshalJSON() ([]byte, error) {
	if e == nil {
		return []byte{}, nil
	}
	return []byte(fmt.Sprintf("{\"code\":%d,\"msg\":\"%s\",\"errmsg\":\"%s\"}", e.code, e.msg, e.Error())), nil
}

func (e *err) UnmarshalJSON(b []byte) error {
	mp := map[string]interface{}{}
	err := json.Unmarshal(b, &mp)
	if err != nil {
		return err
	}
	e.code, _ = strconv.ParseInt(fmt.Sprintf("%d", mp["code"]), 10, 64)
	e.msg = mp["msg"].(string)
	if mp["errmsg"].(string) != "" && mp["errmsg"].(string) != e.msg {
		e.e = errors.New(mp["errmsg"].(string))
	}
	return nil
}

func (e *err) MarshalMsgpack() ([]byte, error) {
	return e.MarshalJSON()
}

func (e *err) UnmarshalMsgpack(b []byte) error {
	return e.UnmarshalJSON(b)
}

//NewError 新建错误
func NewError(code int64, msg string, errs ...error) Error {
	ret := &err{code: code, msg: msg}
	if errs != nil && len(errs) > 0 {
		ret.e = errs[0]
	}
	return ret
}
