package pprofutil

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"runtime"

	"github.com/labstack/echo/v4"
)

// RegisterPprofForEcho 注册pprof路由
// 如果有全局路由拦截器需要把下面所有路由放在白名单中，或者不进行用户/权限校验
func RegisterPprofForEcho(e *echo.Echo) {
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)
	router := e.Group("/debug/pprof")
	router.GET("", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	router.GET("/allocs", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	router.GET("/block", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	router.GET("/goroutine", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	router.GET("/heap", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	router.GET("/mutex", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	router.GET("/threadcreate", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	router.GET("/cmdline", echo.WrapHandler(http.HandlerFunc(pprof.Cmdline)))
	router.GET("/profile", echo.WrapHandler(http.HandlerFunc(pprof.Profile)))
	router.GET("/symbol", echo.WrapHandler(http.HandlerFunc(pprof.Symbol)))
	router.GET("/trace", echo.WrapHandler(http.HandlerFunc(pprof.Trace)))
}

// StartPprofServer 启动pprof服务
// port: 端口号
// 针对本身没有提供http服务的服务，可以通过此方法启动pprof服务
// 如果外部包含原生http服务，不用传port参数
func StartPprofServer(port int) {
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)
	if port > 0 {
		go func() {
			_ = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		}()
	}
}
