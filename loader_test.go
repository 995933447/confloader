package confloader

import (
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func TestLoader(t *testing.T) {
	type Conf struct {
		Name string
		Desc string
	}
	var (
		c     Conf
		errCh = make(chan error)
		sc    = make(chan os.Signal)
	)
	signal.Notify(sc, syscall.SIGINT)
	loader := NewLoader("./test_json.json", time.Second /* 定时更新配置间隔 */, &c)
	go func() {
		tk := time.Tick(time.Second)
		for {
			select {
			case _ = <-tk:
				t.Logf("c:%+v", c)
			case err := <-errCh:
				t.Error(err)
			case <-sc:
				t.Log("cancel loop")
				// 推出循环阻塞
				loader.CancelWatch()
			}
		}
	}()

	// 阻塞定时加载更新配置
	loader.WatchToLoad(errCh)
}
