package pprofutil_test

import (
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/Jsharkc/mygopkg/pprofutil"
)

func TestRun(t *testing.T) {
	pprofutil.StartPprofServer(8080)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-signalChan
}
