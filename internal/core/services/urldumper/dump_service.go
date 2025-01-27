package urldumper

import (
	"sync"

	"github.com/physicist2018/url-shortener-go/internal/core/ports/dumper"
)

type URLDumper struct {
	mu     sync.Mutex
	dumper dumper.URLDumper
	dbfile string
}

func NewURLDumper(dumper dumper.URLDumper, filename string) *URLDumper {
	return &URLDumper{
		dumper: dumper,
		dbfile: filename,
	}
}

func (d *URLDumper) Dump() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.dumper.Dump(d.dbfile)
}

func (d *URLDumper) Restore() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.dumper.Restore(d.dbfile)
}
