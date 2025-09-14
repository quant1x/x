package logger

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/quant1x/pkg/uuid"
	"github.com/quant1x/x/api"
	"github.com/quant1x/x/gls"
	"github.com/quant1x/x/mdc"
)

func TestGoId(t *testing.T) {
	buf := make([]byte, 1024)
	runtime.Stack(buf, false)
	text := api.Bytes2String(buf)
	fmt.Println(text)
}

func TestLogger(t *testing.T) {
	InitLogger("/opt/logs/test")
	u1 := uuid.NewV4()
	defer gls.DeleteGls(gls.GoID())
	mdc.Set(mdc.APP_TRACEID, u1.String())
	//logger := api.GetLogger("test1")
	//SetConsole()
	for i := 0; i < 200; i++ {
		Infof("info-%d", i)
		time.Sleep(time.Millisecond * 1)
	}
	Infof("测试中文\n")
	Debug("debug")
	Error("error")
	Warn("warn")
	Info("测试中文")
	fmt.Println("ok")
	Fatal("xxx")
	//logger.FlushLogger()
	FlushLogger()
	mdc.Remove(mdc.APP_TRACEID)
}
