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

func (server *GinServer) handleFunctions(method, path string, handlers ...func(c *Gtx)) gin.IRoutes {
	functions := make([]gin.HandlerFunc, len(handlers))
	for _, handler := range handlers {
		functions = append(functions, func(c *gin.Context) {
			handler(&Gtx{Context: c})
		})
	}
	return server.Engine.Handle(method, path, functions...)
}

func (server *GinServer) GET(relativePath string, handlers ...func(c *Gtx)) gin.IRoutes {
	return server.handleFunctions(http.MethodGet, relativePath, handlers...)
}

func (server *GinServer) POST(relativePath string, handlers ...func(c *Gtx)) gin.IRoutes {
	return server.handleFunctions(http.MethodPost, relativePath, handlers...)
}

func (server *GinServer) PATCH(relativePath string, handlers ...func(c *Gtx)) gin.IRoutes {
	return server.handleFunctions(http.MethodPatch, relativePath, handlers...)
}

func (server *GinServer) PUT(relativePath string, handlers ...func(c *Gtx)) gin.IRoutes {
	return server.handleFunctions(http.MethodPut, relativePath, handlers...)
}
func (server *GinServer) DELETE(relativePath string, handlers ...func(c *Gtx)) gin.IRoutes {
	return server.handleFunctions(http.MethodDelete, relativePath, handlers...)
}

func (server *GinServer) HEAD(relativePath string, handlers ...func(c *Gtx)) gin.IRoutes {
	return server.handleFunctions(http.MethodHead, relativePath, handlers...)
}

func (server *GinServer) Group(relativePath string, handlers ...func(c *Gtx)) *GinGroup {
	RHandles := make([]gin.HandlerFunc, 0)
	for _, handle := range handlers {
		RHandles = append(RHandles, func(c *gin.Context) {
			handle(&Gtx{Context: c})
		})
	}
	return &GinGroup{server.Engine.Group(relativePath, RHandles...)}
}

func (r *GinGroup) handleGroups(method, path string, handlers ...func(c *Gtx)) gin.IRoutes {
	functions := make([]gin.HandlerFunc, len(handlers))
	for _, handler := range handlers {
		functions = append(functions, func(c *gin.Context) {
			handler(&Gtx{Context: c})
		})
	}
	return r.RouterGroup.Handle(method, path, functions...)
}

func (r *GinGroup) GET(relativePath string, handlers ...func(c *Gtx)) gin.IRoutes {
	return r.handleGroups(http.MethodGet, relativePath, handlers...)
}

func (r *GinGroup) POST(relativePath string, handlers ...func(c *Gtx)) gin.IRoutes {
	return r.handleGroups(http.MethodPost, relativePath, handlers...)
}

func (r *GinGroup) PUT(relativePath string, handlers ...func(c *Gtx)) gin.IRoutes {
	return r.handleGroups(http.MethodPut, relativePath, handlers...)
}

func (r *GinGroup) PATCH(relativePath string, handlers ...func(c *Gtx)) gin.IRoutes {
	return r.handleGroups(http.MethodPatch, relativePath, handlers...)
}

func (r *GinGroup) DELETE(relativePath string, handlers ...func(c *Gtx)) gin.IRoutes {
	return r.handleGroups(http.MethodDelete, relativePath, handlers...)
}

func (r *GinGroup) HEAD(relativePath string, handlers ...func(c *Gtx)) gin.IRoutes {
	return r.handleGroups(http.MethodHead, relativePath, handlers...)
}

func (r *GinGroup) Use(middlewares ...func(c *Gtx)) gin.IRoutes {
	rMiddlewares := make([]gin.HandlerFunc, 0)
	for _, middleware := range middlewares {
		rMiddlewares = append(rMiddlewares, func(c *gin.Context) {
			middleware(&Gtx{Context: c})
		})
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
