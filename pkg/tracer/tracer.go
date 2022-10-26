package tracer

import (
	"fmt"

	"github.com/pkg/errors"
)

type TraceCtx struct{}

// Background 返回空TraceContext
func Background() *TraceCtx {
	return &TraceCtx{}
}

// Clear 清空结构体信息并返回
func (ctx *TraceCtx) Clear() *TraceCtx {
	return &TraceCtx{}
}

func FormatParam(args ...interface{}) string {
	s := "Parameter: "
	for _, arg := range args {
		s = s + fmt.Sprintf("%+v   ", arg)
	}
	return fmt.Sprintf("{%s}", s)
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func FormatErr(err error) (s string) {
	e, ok := err.(stackTracer)
	if !ok {
		return err.Error()
	}
	//类型断言,检查error接口中的动态值是否实现了stackTracer接口。如果实现了返回类型为stackTracer的接口
	st := e.StackTrace()
	for i, frame := range st {
		if i == len(st)-1 {
			break
		}
		s = s + fmt.Sprintf("%n  <-- ", frame)
	}
	return s
}
