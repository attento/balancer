package source

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/attento/balancer/app"
	"github.com/attento/balancer/app/core"
	"github.com/stretchr/testify/assert"
)

type MyDaemonMock struct {
	app.DaemonInterface
}

func (d *MyDaemonMock) StartHttpServer(a core.Address, f core.Filter, us []*core.Upstream) error {
	return nil
}

func TestExtractWrongFile(t *testing.T) {
	d := MyDaemonMock{}
	extractor := FileExtractor{
		Daemon: &d,
		Path:   "nop",
	}
	err := extractor.Extract()
	assert.EqualError(t, err, "unable to open file")
}

func TestConfigServer(t *testing.T) {
	file, _ := ioutil.TempFile(os.TempDir(), "prefix")
	d1 := []byte("[{\"address\":\":80\"}]")
	ioutil.WriteFile(file.Name(), d1, 0644)
	defer os.Remove(file.Name())
	d := MyDaemonMock{}
	extractor := FileExtractor{
		Daemon: &d,
		Path:   file.Name(),
	}
	err := extractor.Extract()
	assert.Equal(t, err, nil)
}
