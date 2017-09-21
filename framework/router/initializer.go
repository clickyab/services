package router

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/middleware"
	"github.com/clickyab/services/initializer"
	"github.com/rs/xhandler"
	"github.com/rs/xmux"
	"github.com/sirupsen/logrus"
	onion "gopkg.in/fzerorubigd/onion.v3"
)

var (
	engine *xmux.Mux
	all    []Routes
	mid    []GlobalMiddleware

	// this is development mode
	mountPoint = config.RegisterString("services.framework.controller.mount_point", "/api", "http controller mount point")
	hammerTime = config.RegisterDuration("services.framework.controller.graceful_wait", 100*time.Millisecond, "the time for framework to stop for graceful shutdown")
	listen     onion.String
)

// GlobalMiddleware is the middleware that must be on all routes
type GlobalMiddleware interface {
	Handler(framework.Handler) framework.Handler
}

// Mux is the simple router interface
type Mux interface {

	// GET is a shortcut for mux.Handle("GET", path, handler)
	GET(string, framework.Handler)

	// HEAD is a shortcut for mux.Handle("HEAD", path, handler)
	HEAD(string, framework.Handler)

	// OPTIONS is a shortcut for mux.Handle("OPTIONS", path, handler)
	OPTIONS(string, framework.Handler)

	// POST is a shortcut for mux.Handle("POST", path, handler)
	POST(string, framework.Handler)

	// PUT is a shortcut for mux.Handle("PUT", path, handler)
	PUT(string, framework.Handler)

	// PATCH is a shortcut for mux.Handle("PATCH", path, handler)
	PATCH(string, framework.Handler)

	// DELETE is a shortcut for mux.Handle("DELETE", path, handler)
	DELETE(string, framework.Handler)

	// NewGroup creates a new routes group with the provided path prefix.
	// All routes added to the returned group will have the path prepended.
	NewGroup(string) Mux
}

// Routes the base rote structure
type Routes interface {
	// Routes is for adding new controller
	Routes(Mux)
}

type initer struct {
}

func (i *initer) Initialize(ctx context.Context) {
	engine = xmux.New()

	f := func(next framework.Handler) framework.Handler {
		for i := range mid {
			next = mid[i].Handler(next)
		}

		return next
	}

	xm := &xmuxer{
		middleware: f,
	}
	mp := mountPoint.String()
	if mp != "" {
		xm.group = engine.NewGroup(mp)
	} else {
		xm.engine = engine
	}

	for i := range all {
		all[i].Routes(xm)
	}
	// Append some generic middleware, to handle recovery and log
	handler := middleware.Recovery(
		middleware.Logger(
			xhandler.New(context.Background(), engine).ServeHTTP,
		),
	)
	server := &http.Server{Addr: listen.String(), Handler: handler}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			logrus.Error(err)
		}
	}()

	done := ctx.Done()
	assert.NotNil(done, "[BUG] the done channel is nil")
	go func() {
		<-done
		ctx, _ := context.WithTimeout(context.Background(), hammerTime.Duration())
		server.Shutdown(ctx)
	}()
}

// Register a new controller class
func Register(c ...Routes) {
	all = append(all, c...)
}

// RegisterGlobalMiddleware is a function to register a global middleware
func RegisterGlobalMiddleware(g GlobalMiddleware) {
	mid = append(mid, g)
}

func init() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}
	listen = config.RegisterString(
		"services.framework.listen",
		fmt.Sprintf(":%s", port),
		"address to listen for framework",
	)

	initializer.Register(&initer{}, 100)
}
