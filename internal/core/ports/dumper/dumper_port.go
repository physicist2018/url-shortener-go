package dumper

type URLDumper interface {
	Dump(filename string) error
	Restore(filename string) error
}
