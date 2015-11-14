package argvreader

import (
	"io"
	"os"
)

type reader struct {
	current io.Reader
	args    []string
}

func New() io.Reader {
	args := os.Args[1:]
	if len(args) == 0 {
		return os.Stdin
	}
	return &reader{
		current: nil,
		args:    args,
	}
}

func (r *reader) Read(p []byte) (n int, err error) {
	if r.current == nil {
		err = r.openNext()
		if err != nil {
			return 0, err
		}
	}
	for {
		n, err = r.current.Read(p)
		if err == nil {
			return n, err
		}
		if err == io.EOF {
			err := r.openNext()
			if err == io.EOF {
				return 0, io.EOF
			} else if err != nil {
				return 0, err
			}
			continue
		}
		return n, err
	}
}

func (r *reader) openNext() error {
	r.current = nil
	if len(r.args) == 0 {
		return io.EOF
	}
	file, err := os.Open(r.args[0])
	if err != nil {
		return err
	}
	r.current = file
	r.args = r.args[1:]
	return nil
}
