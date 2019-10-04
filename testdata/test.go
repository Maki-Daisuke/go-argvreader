package main

import (
	"fmt"
	"os"

	"github.com/Maki-Daisuke/go-argvreader"
	"github.com/mattn/go-forlines"
)

func main() {
	r := argvreader.New()
	for {
		err := forlines.Do(r, func(line string) error {
			fmt.Printf("%s: %s\n", r.Name(), line)
			return nil
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		break
	}
	os.Exit(0)
}
