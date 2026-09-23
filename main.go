package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	infile    = flag.String("infile", "", "File to render")
	outfile   = flag.String("outfile", "", "Output file")
	secretDir = flag.String("secret-dir", "/run/secrets", "Where to look for secrets")
)

var (
	funcs = template.FuncMap{
		"extend": func(kv ...any) (map[string]any, error) {
			if len(kv)%2 != 0 {
				return nil, errors.New("extend: odd number of args")
			}
			m := make(map[string]any, len(kv)/2)
			for i := 0; i < len(kv); i += 2 {
				k, ok := kv[i].(string)
				if !ok {
					return nil, errors.New("extend: keys must be strings")
				}
				m[k] = kv[i+1]
			}
			return m, nil
		},
	}
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

	tmp, err := template.New(*infile).Funcs(funcs).ParseFiles(*infile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse file: %v\n", *infile)
		os.Exit(1)
	}

	env := map[string]string{}
	for _, kv := range os.Environ() {
		k, v, _ := strings.Cut(kv, "=")
		env[k] = v
	}

	if *secretDir != "" {
		entries, err := os.ReadDir(*secretDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read secrets dir (%v): %v\n", *secretDir, err)
			os.Exit(1)
		}

		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			b, err := os.ReadFile(filepath.Join(*secretDir, e.Name()))
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to read file: %v: %v\n", e.Name(), err)
				os.Exit(1)
			}

			key := fmt.Sprintf("SECRET_%v", strings.ToUpper(filepath.Base(e.Name())))
			env[key] = strings.TrimSpace(string(b))
		}
	}

	if err := tmp.Execute(out, env); err != nil {
		fmt.Fprintf(os.Stderr, "failed to render outfile: %v\n", err)
		os.Exit(1)
	}
}
