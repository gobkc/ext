package ext

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gobkc/to"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Context struct {
	*gin.Context
	handlerPFFunc    *func(*Context) any
	handlerOtherFunc *func(*Context, uint) any
}

type HandlerFunc func(*Context)

type Server struct {
	engine           *gin.Engine
	handlerPFFunc    *func(*Context) any
	handlerOtherFunc *func(*Context, uint) any
	mid
}

func NewServer() *Server {
	s := &Server{engine: gin.New()}
	s.mid.s = s
	return s
}

func (s *Server) GET(path string, handlers ...HandlerFunc) {
	s.engine.GET(path, s.ConvertHandler(handlers...))
}

func (s *Server) POST(path string, handlers ...HandlerFunc) {
	s.engine.POST(path, s.ConvertHandler(handlers...))
}

func (s *Server) PUT(path string, handlers ...HandlerFunc) {
	s.engine.PUT(path, s.ConvertHandler(handlers...))
}

func (s *Server) PATCH(path string, handlers ...HandlerFunc) {
	s.engine.PATCH(path, s.ConvertHandler(handlers...))
}

func (s *Server) HEAD(path string, handlers ...HandlerFunc) {
	s.engine.HEAD(path, s.ConvertHandler(handlers...))
}

func (s *Server) DELETE(path string, handlers ...HandlerFunc) {
	s.engine.DELETE(path, s.ConvertHandler(handlers...))
}

func (s *Server) OPTIONS(path string, handlers ...HandlerFunc) {
	s.engine.OPTIONS(path, s.ConvertHandler(handlers...))
}

func (s *Server) GROUP(path string, handlers ...HandlerFunc) *Group {
	g := s.engine.Group(path, s.ConvertHandler(handlers...))
	ng := NewGroup(g)
	ng.handlerPFFunc = s.handlerPFFunc
	ng.handlerOtherFunc = s.handlerOtherFunc
	ng.mid = s.mid
	ng.mid.g = ng
	ng.mid.g.handlerPFFunc = s.handlerPFFunc
	ng.mid.g.handlerOtherFunc = s.handlerOtherFunc
	return ng
}

type Group struct {
	group            *gin.RouterGroup
	handlerPFFunc    *func(*Context) any
	handlerOtherFunc *func(*Context, uint) any
	mid
}

func NewGroup(group *gin.RouterGroup) *Group {
	return &Group{group: group}
}

func (g *Group) GET(path string, handlers ...HandlerFunc) {
	g.group.GET(path, g.ConvertHandler(handlers...))
}

func (g *Group) POST(path string, handlers ...HandlerFunc) {
	g.group.POST(path, g.ConvertHandler(handlers...))
}

func (g *Group) PUT(path string, handlers ...HandlerFunc) {
	g.group.PUT(path, g.ConvertHandler(handlers...))
}

func (g *Group) PATCH(path string, handlers ...HandlerFunc) {
	g.group.PATCH(path, g.ConvertHandler(handlers...))
}

func (g *Group) HEAD(path string, handlers ...HandlerFunc) {
	g.group.HEAD(path, g.ConvertHandler(handlers...))
}

func (g *Group) DELETE(path string, handlers ...HandlerFunc) {
	g.group.DELETE(path, g.ConvertHandler(handlers...))
}

func (g *Group) OPTIONS(path string, handlers ...HandlerFunc) {
	g.group.OPTIONS(path, g.ConvertHandler(handlers...))
}

func (g *Group) GROUP(path string, handlers ...HandlerFunc) *Group {
	og := g.group.Group(path, g.ConvertHandler(handlers...))
	ng := NewGroup(og)
	ng.handlerPFFunc = g.handlerPFFunc
	ng.handlerOtherFunc = g.handlerOtherFunc
	return ng
}

func (s *Server) SetGroupRouters(name string, router func(g *Group), handlers ...HandlerFunc) {
	rs := s.GROUP(name, handlers...)
	router(rs)
}

func (s *Server) SetPublicRouters(router func(s *Server)) {
	router(s)
}

func (s *Server) RUN(addr string) {
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	s.engine.Run(addr)
	sign := <-signals
	slog.Default().Warn(`Stopping HTTP RESTFul Server`, slog.String(`Signal`, sign.String()))
}

var pfOnce sync.Once

func (s *Server) SetHandlerProfileFunc(handlerPFFunc func(*Context) any) {
	*s.handlerPFFunc = handlerPFFunc
}

func (s *Server) SetHandlerContextValueFunc(handlerContextValueFunc func(*Context, uint) any) {
	*s.handlerOtherFunc = handlerContextValueFunc
}

type mid struct {
	s *Server
	g *Group
}

func (m *mid) ConvertHandler(fs ...HandlerFunc) gin.HandlerFunc {
	return func(context *gin.Context) {
		c := &Context{
			Context:          context,
			handlerPFFunc:    m.s.handlerPFFunc,
			handlerOtherFunc: m.s.handlerOtherFunc,
		}
		for _, f := range fs {
			if c.IsAborted() {
				return
			}
			f(c)
		}
	}
}

func (c *Context) JSON(data any) {
	if c.Request == nil {
		return
	}
	mb, _ := MarshalGzipJson(data)
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Content-Encoding", "gzip")
	c.Writer.Header().Set("Vary", "Accept-Encoding")
	c.Writer.Header().Set("Content-Length", fmt.Sprintf("%v", len(mb)))
	c.Writer.Write(mb)
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Status()
}

func (c *Context) JSONWithHTTPCode(httpCode int, data any) {
	c.JSON(data)
	c.Writer.WriteHeader(httpCode)
	c.Writer.Status()
}

type PaginationParameters struct {
	Limit   int64  `json:"limit"`
	Offset  int64  `json:"offset"`
	Keyword string `json:"keyword"`
}

func (c *Context) GetPaginationParameters() *PaginationParameters {
	return &PaginationParameters{
		Limit:   to.Int[int64](c.DefaultQuery(`limit`, `10`)),
		Offset:  to.Int[int64](c.DefaultQuery(`offset`, `0`)),
		Keyword: c.Query(`keyword`),
	}
}

var GetUserProfileErrParamErr = errors.New(`parameter is not a pointer`)

func (c *Context) GetUserProfile() any {
	if c.handlerPFFunc == nil {
		return nil
	}
	data := (*c.handlerPFFunc)(c)
	return data
}

func (c *Context) GetContextWithKey(key uint) any {
	if c.handlerPFFunc == nil {
		return nil
	}
	data := (*c.handlerOtherFunc)(c, key)
	return data
}
