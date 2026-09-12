package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
)

var (
	infile  = flag.String("infile", "", "File to render")
	outfile = flag.String("outfile", "", "Output file")
)

func main() {
	flag.Parse()

	if infile == nil || *infile == "" {
		fmt.Fprintln(os.Stderr, "infile must be specified")
		os.Exit(1)
	}

	var out io.Writer

	if outfile == nil || *outfile == "" {
		fmt.Fprintln(os.Stderr, "outfile is not specified, using stdout")
		out = os.Stdout
	} else {
		f, err := os.Create(*outfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create outfile: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()

		out = f
	}

	tmp, err := template.ParseFiles(*infile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse file: %v\n", *infile)
		os.Exit(1)
	}

	env := map[string]string{}
	for _, kv := range os.Environ() {
		k, v, _ := strings.Cut(kv, "=")
		env[k] = v
	}

	if err := tmp.Execute(out, env); err != nil {
		fmt.Fprintf(os.Stderr, "failed to render outfile: %v\n", err)
		os.Exit(1)
	}
}
