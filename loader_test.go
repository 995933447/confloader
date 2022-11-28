package confloader

import (
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func TestLoader(t *testing.T) {
	type conf struct{
		Name string
		Desc string
	}
	var (
		c conf
		errCh = make(chan error)
		sc = make(chan os.Signal)
	)
	signal.Notify(sc, syscall.SIGINT)
	loader := NewLoader("./test_json.json", time.Second, &c)
	go func() {
		tk := time.Tick(time.Second)
		for {
			select {
			case _ = <- tk:
				t.Logf("c:%+v", c)
			case err := <- errCh:
				t.Error(err)
			case <-sc:
				t.Log("cancel loop")
				loader.CancelLoop()
			}
		}
	}()

	loader.WatchToLoad(errCh)
}
