// Copyright 2015 Daisuke (yet another) Maki. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package argvreader makes it a bit less painful to write filter programs
in Go especially for Perl programmers.


Description

The `argvreader.Reader` behaves just like Perl's `<ARGV>`, that reads
lines from one file after another if file list is given in command line
arguments, reads from STDIN otherwise.
This may be useful to implement UNIX-style filter command, like gzip
or grep.


Example

		import argvreader "github.com/Maki-Daisuke/go-argvreader"

		func main(){
		  r := argvreader.New()
		  b := bufio.NewReader(r)
		  for{
		    line, err := b.ReadString('\n')
		    if err != nil {
		      if err == io.EOF {
		        break
		      }
		      panic(err)
		    }
		    do_something(line)
		  }
		}

Then, in your command line:

		$ ./main foo bar baz

This reads lines from foo, bar and baz one after another.
If no argument is given:

		$ ./main

This reads lines from STDIN instead.

You can pass file list manually, of course:

		import flags "github.com/jessevdk/go-flags"

		var opts struct { ... }

		args, err := flags.ParseArgs(&opts, args)
		r := argvreader.NewReader(args)
*/
package argvreader

import (
	"io"
	"os"
)

// Reader provides functionality just like Perl's `ARGV` file-handle and
// `$ARGV` variable.
//
// Name returns a file name that is currently open and being read.
// It returns "-" if the Reader is reading from STDIN.
type Reader interface {
	io.Reader
	Name() string
}

type stdinReader struct {
	*os.File // for os.Stdin only
}

func (r stdinReader) Name() string {
	return "-"
}

type argvReader struct {
	current Reader
	args    []string
}

// NewReader creates Reader by manually passing file list. If the list is empty
// or nil, it returns a Reader reading from os.Stdin.
func NewReader(args []string) Reader {
	if len(args) == 0 {
		return stdinReader{os.Stdin}
	}
	return &argvReader{
		current: nil,
		args:    args,
	}
}

// New is just shorthand of `NewReader(os.Args[1:])`.
func New() Reader {
	return NewReader(os.Args[1:])
}

func (r *argvReader) Name() string {
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
		if f, ok := r.current.(*os.File); ok {
			f.Close()
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

func (r *argvReader) openNext() error {
	r.current = nil
	if len(r.args) == 0 {
		return io.EOF
	}
	next := r.args[0]
	r.args = r.args[1:]
	if next == "-" {
		r.current = stdinReader{os.Stdin}
	} else {
		file, err := os.Open(next)
		if err != nil {
			return err
		}
		r.current = file
	}
	return nil
}
