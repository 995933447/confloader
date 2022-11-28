package confloader

import (
	"encoding/json"
	"errors"
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
	"time"
)

type Loader struct {
	loadFile string
	reloadInterval time.Duration
	cancelLoopSignCh chan struct{}
	conf interface{}
}

func NewLoader(loadConfigFile string, reloadInterval time.Duration, conf interface{}) *Loader {
	loader := &Loader{
		loadFile: loadConfigFile,
		reloadInterval: reloadInterval,
		cancelLoopSignCh: make(chan struct{}),
		conf: conf,
	}
	return loader
}

func (l *Loader) WatchToLoad(errCh chan error) {
	doLoad := func() {
		if err := l.Load(); err != nil && errCh != nil {
			errCh <- err
		}
	}

	doLoad()

	var (
		reloadTicker = time.NewTicker(l.reloadInterval)
		isCanceledLoop bool
	)
	for {
		select {
		case _ = <- reloadTicker.C:
			doLoad()
		case _ = <- l.cancelLoopSignCh:
			reloadTicker.Stop()
			isCanceledLoop = true
		}

		if isCanceledLoop {
			break
		}
	}
}

func (l *Loader) Load() error {
	buf, err := os.ReadFile(l.loadFile)
	if err != nil {
		return err
	}

	switch filepath.Ext(l.loadFile) {
	case ".json":
		err = json.Unmarshal(buf, l.conf)
		if err != nil {
			return err
		}
	case ".toml":
		_, err = toml.Decode(string(buf), l.conf)
		if err != nil {
			return err
		}
	default:
		return errors.New("either TOML or JSON is supported")
	}

	return nil
}

func (l *Loader) CancelLoop() {
	l.cancelLoopSignCh <- struct{}{}
}