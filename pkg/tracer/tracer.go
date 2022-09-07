package tracer

import (
	"runtime"
)

type TraceCtx struct {
	FuncLink []string //函数的调用链路
}

func (ctx TraceCtx) ToString() string {
	var ans string
	for index, valueStr := range ctx.FuncLink { //下标和值
		if index != len(ctx.FuncLink)-1 {
			ans += valueStr + " --> "
		} else {
			ans += valueStr
		}
	}
	return ans
}

// TraceCaller 此方法记录调用该方法的函数信息并保存在FuncLink中
func (ctx *TraceCtx) TraceCaller() *TraceCtx {
	funcPc, _, _, ok := runtime.Caller(1)
	if !ok {
		return nil
	}
	callerName := runtime.FuncForPC(funcPc).Name()
	ctx.FuncLink = append(ctx.FuncLink, callerName)
	return ctx
}

// Background 返回空TraceContext
func Background() *TraceCtx {
	return &TraceCtx{}
}

// Clear 清空结构体信息并返回
func (ctx *TraceCtx) Clear() *TraceCtx {
	ctx.FuncLink = ctx.FuncLink[0:0]
	return ctx
}
