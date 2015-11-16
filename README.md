# go-argvreader

    import argvreader "github.com/Maki-Daisuke/go-argvreader"

Package argvreader makes it a bit less painful to write filter programs in Go
especially for Perl programmers.


### Description

The `argvreader.Reader` behaves just like Perl's `<ARGV>`, that reads lines from
one file after another if file list is given in command line arguments, reads
from STDIN otherwise. This may be useful to implement UNIX-style filter command,
like gzip or grep.

### Example

```go
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
```

Then, in your command line:

    $ ./main foo bar baz

This reads lines from foo, bar and baz one after another. If no argument is
given:

    $ ./main

This reads lines from STDIN instead.

You can pass file list manually, of course:

```go
    import flags "github.com/jessevdk/go-flags"

    var opts struct { ... }

    args, err := flags.ParseArgs(&opts, args)
    r := argvreader.NewReader(args)
```

## Usage

#### type Reader

```go
type Reader interface {
	io.Reader
	CurrentFileName() string
}
```

`Reader` provides functionality just like Perl's `ARGV` file-handle and `$ARGV`
variable.

`CurrnetFileName` returns a file name that is currently open and being read. It
returns `"-"` if the Reader is reading from STDIN.

#### func  New

```go
func New() Reader
```
New is just shorthand of `NewReader(os.Args[1:])`.

#### func  NewReader

```go
func NewReader(args []string) Reader
```
`NewReader` creates Reader by manually passing file list. If the list is empty or
nil, it returns a Reader reading from `os.Stdin`.


## License

The Simplified BSD License (2-clause).
See [LICENSE](LICENSE) file also.


## Author

Daisuke (yet another) Maki
