package gext

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GinServer struct {
	*gin.Engine
}

type Gtx struct {
	*gin.Context
	ID int
}

type GinGroup struct {
	*gin.RouterGroup
}

func NewGinServer() *GinServer {
	gServer := &GinServer{Engine: gin.Default()}
	gServer.Engine.GET(`/ping`, Cors(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{`msg`: `pong`})
	})
	return gServer
}

func handleFunc(handler func(c *Gtx)) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		handler(&Gtx{Context: c})
	}
}

func (server *GinServer) Group(relativePath string, handlers ...func(c *Gtx)) *GinGroup {
	RHandles := make([]gin.HandlerFunc, 0)
	for _, handle := range handlers {
		RHandles = append(RHandles, handleFunc(handle))
	}
	return &GinGroup{server.Engine.Group(relativePath, RHandles...)}
}

func (server *GinServer) GET(relativePath string, handlers ...func(c *Gtx)) gin.IRoutes {
	RHandles := make([]gin.HandlerFunc, 0)
	for _, handle := range handlers {
		RHandles = append(RHandles, handleFunc(handle))
	}
	return server.Engine.GET(relativePath, RHandles...)
}

func (server *GinServer) POST(relativePath string, handlers ...func(c *Gtx)) gin.IRoutes {
	RHandles := make([]gin.HandlerFunc, 0)
	for _, handle := range handlers {
		RHandles = append(RHandles, handleFunc(handle))
	}
	return server.Engine.POST(relativePath, RHandles...)
}

func (server *GinServer) PATCH(relativePath string, handlers ...func(c *Gtx)) gin.IRoutes {
	RHandles := make([]gin.HandlerFunc, 0)
	for _, handle := range handlers {
		RHandles = append(RHandles, handleFunc(handle))
	}
	return server.Engine.PATCH(relativePath, RHandles...)
}

func (server *GinServer) PUT(relativePath string, handlers ...func(c *Gtx)) gin.IRoutes {
	RHandles := make([]gin.HandlerFunc, 0)
	for _, handle := range handlers {
		RHandles = append(RHandles, handleFunc(handle))
	}
	return server.Engine.PUT(relativePath, RHandles...)
}
func (server *GinServer) DELETE(relativePath string, handlers ...func(c *Gtx)) gin.IRoutes {
	RHandles := make([]gin.HandlerFunc, 0)
	for _, handle := range handlers {
		RHandles = append(RHandles, handleFunc(handle))
	}
	return server.Engine.DELETE(relativePath, RHandles...)
}

func (server *GinServer) HEAD(relativePath string, handlers ...func(c *Gtx)) gin.IRoutes {
	RHandles := make([]gin.HandlerFunc, 0)
	for _, handle := range handlers {
		RHandles = append(RHandles, handleFunc(handle))
	}
	return server.Engine.HEAD(relativePath, RHandles...)
}

func (r *GinGroup) GET(relativePath string, handlers ...func(c *Gtx)) gin.IRoutes {
	rHandles := make([]gin.HandlerFunc, 0)
	for _, handle := range handlers {
		rHandles = append(rHandles, handleFunc(handle))
	}
	return r.RouterGroup.GET(relativePath, rHandles...)
}

func (r *GinGroup) POST(relativePath string, handlers ...func(c *Gtx)) gin.IRoutes {
	rHandles := make([]gin.HandlerFunc, 0)
	for _, handle := range handlers {
		rHandles = append(rHandles, handleFunc(handle))
	}
	return r.RouterGroup.POST(relativePath, rHandles...)
}

func (r *GinGroup) PUT(relativePath string, handlers ...func(c *Gtx)) gin.IRoutes {
	rHandles := make([]gin.HandlerFunc, 0)
	for _, handle := range handlers {
		rHandles = append(rHandles, handleFunc(handle))
	}
	return r.RouterGroup.PUT(relativePath, rHandles...)
}

func (r *GinGroup) PATCH(relativePath string, handlers ...func(c *Gtx)) gin.IRoutes {
	rHandles := make([]gin.HandlerFunc, 0)
	for _, handle := range handlers {
		rHandles = append(rHandles, handleFunc(handle))
	}
	return r.RouterGroup.PATCH(relativePath, rHandles...)
}

func (r *GinGroup) DELETE(relativePath string, handlers ...func(c *Gtx)) gin.IRoutes {
	rHandles := make([]gin.HandlerFunc, 0)
	for _, handle := range handlers {
		rHandles = append(rHandles, handleFunc(handle))
	}
	return r.RouterGroup.DELETE(relativePath, rHandles...)
}

func (r *GinGroup) HEAD(relativePath string, handlers ...func(c *Gtx)) gin.IRoutes {
	rHandles := make([]gin.HandlerFunc, 0)
	for _, handle := range handlers {
		rHandles = append(rHandles, handleFunc(handle))
	}
	return r.RouterGroup.HEAD(relativePath, rHandles...)
}

func (r *GinGroup) Use(middlewares ...func(c *Gtx)) gin.IRoutes {
	rMiddlewares := make([]gin.HandlerFunc, 0)
	for _, middleware := range middlewares {
		rMiddlewares = append(rMiddlewares, handleFunc(middleware))
	}
	return r.RouterGroup.Use(rMiddlewares...)
}

func (g *Gtx) JSON(data any) {
	if g.Request == nil {
		return
	}
	mb, _ := MarshalGzipJson(data)
	g.Writer.Header().Set("Content-Type", "application/json")
	g.Writer.Header().Set("Content-Encoding", "gzip")
	g.Writer.Header().Set("Vary", "Accept-Encoding")
	g.Writer.Header().Set("Content-Length", fmt.Sprintf("%v", len(mb)))
	g.Writer.Write(mb)
	g.Writer.WriteHeader(http.StatusOK)
	g.Writer.Status()
}

func (g *Gtx) JSONWithHTTPCode(httpCode int, data any) {
	if g.Request == nil {
		return
	}
	mb, _ := MarshalGzipJson(data)
	g.Writer.Header().Set("Content-Type", "application/json")
	g.Writer.Header().Set("Content-Encoding", "gzip")
	g.Writer.Header().Set("Vary", "Accept-Encoding")
	g.Writer.Header().Set("Content-Length", fmt.Sprintf("%v", len(mb)))
	g.Writer.Write(mb)
	g.Writer.WriteHeader(httpCode)
	g.Writer.Status()
}

type UserInfoKey struct{}

var UserInfoKeyVal = UserInfoKey{}

func GetProfile[U any](c *Gtx) (userInfo *U) {
	userInfo = new(U)
	profile := c.Request.Context().Value(UserInfoKeyVal)
	query, ok := profile.(U)
	if ok {
		userInfo = &query
	}
	return
}
