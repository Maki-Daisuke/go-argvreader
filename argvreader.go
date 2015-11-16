package argvreader

import (
	"io"
	"os"
)

type Reader interface {
	io.Reader
	CurrentFileName() string
}

type stdinReader struct {
	*os.File // for os.Stdin only
}

func (r stdinReader) CurrentFileName() string {
	return "-"
}

type argvReader struct {
	current *os.File
	args    []string
}

func New() Reader {
	args := os.Args[1:]
	if len(args) == 0 {
		return stdinReader{os.Stdin}
	}
	return &argvReader{
		current: nil,
		args:    args,
	}
}

func (r *argvReader) CurrentFileName() string {
	if r.current == nil {
		return ""
	} else {
		return r.current.Name()
	}
}

func (r *argvReader) Read(p []byte) (n int, err error) {
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
		r.current.Close()
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

func (r *argvReader) openNext() error {
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
