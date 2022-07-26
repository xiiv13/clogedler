package framework

import (
	"context"
	"net/http"
	"sync"
	"time"
)

var _ context.Context = new(Context)

type Context struct {
	request        *http.Request
	responseWriter http.ResponseWriter
	ctx            context.Context

	hasTimeout bool        // 是否超时标记位
	writerMux  *sync.Mutex // 写保护机制

	// 当前请求的handler链条
	handlerlist []ControllerHandler
	handlerIdx  int // 当前请求调用到调用链的哪个节点

	params map[string]string // url路由匹配的参数
}

func NewContext(r *http.Request, w http.ResponseWriter) *Context {
	return &Context{
		request:        r,
		responseWriter: w,
		ctx:            r.Context(),
		writerMux:      &sync.Mutex{},
		handlerIdx:     -1,
	}
}

// #region base function

func (ctx *Context) WriterMux() *sync.Mutex {
	return ctx.writerMux
}

func (ctx *Context) GetRequest() *http.Request {
	return ctx.request
}

func (ctx *Context) GetResponse() http.ResponseWriter {
	return ctx.responseWriter
}

func (ctx *Context) SetHasTimeout() {
	ctx.hasTimeout = true
}

func (ctx *Context) HasTimeout() bool {
	return ctx.hasTimeout
}

// SetHandlers 为context设置handlers
func (ctx *Context) SetHandlers(handlers []ControllerHandler) {
	ctx.handlerlist = handlers
}

func (ctx *Context) setParams(params map[string]string) {
	ctx.params = params
}

// Next 核心函数，调用context的下一个函数
func (ctx *Context) Next() error {
	ctx.handlerIdx++
	if ctx.handlerIdx < len(ctx.handlerlist) {
		if err := ctx.handlerlist[ctx.handlerIdx](ctx); err != nil {
			return err
		}
	}
	return nil
}

// #endregion

func (ctx *Context) BaseContext() context.Context {
	return ctx.request.Context()
}

// #region implement context.Context

func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return ctx.BaseContext().Deadline()
}

func (ctx *Context) Done() <-chan struct{} {
	return ctx.BaseContext().Done()
}

func (ctx *Context) Err() error {
	return ctx.BaseContext().Err()
}

func (ctx *Context) Value(key any) interface{} {
	return ctx.BaseContext().Value(key)
}

// #endregion
