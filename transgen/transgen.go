package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tpukep/bambam/parser"
)

func main() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	debug := flag.Bool("debug", false, "print lots of debug info as we process.")
	// verrequest := flag.Bool("version", false, "request git commit hash used to build this bambam")
	outdir := flag.String("o", "odir", "specify output directory")
	pkg := flag.String("p", "main", "specify package for generated code")
	privs := flag.Bool("X", false, "export private as well as public struct fields")
	overwrite := flag.Bool("OVERWRITE", false, "replace named .go files with capid tagged versions.")
	flag.Parse()

	if debug != nil {
		parser.Verbose = *debug
	}

	if outdir == nil || *outdir == "" {
		fmt.Fprintf(os.Stderr, "required -o option missing. Use bambam -o <dirname> myfile.go # to specify the output directory.\n")
		use()
	}

	// if verrequest != nil && *verrequest {
	// 	fmt.Printf("%s\n", parser.LASTGITCOMMITHASH)
	// 	os.Exit(0)
	// }

	if !parser.DirExists(*outdir) {
		err := os.MkdirAll(*outdir, 0755)
		if err != nil {
			panic(err)
		}
	}

	if pkg == nil || *pkg == "" {
		fmt.Fprintf(os.Stderr, "required -p option missing. Specify a package name for the generated go code with -p <pkgname>\n")
		use()
	}

	// all the rest are input .go files
	// inputFiles := flag.Args()

	// if len(inputFiles) == 0 {
	// 	fmt.Fprintf(os.Stderr, "bambam needs at least one .go golang source file to process specified on the command line.\n")
	// 	os.Exit(1)
	// }

	// for _, fn := range inputFiles {
	// 	if !strings.HasSuffix(fn, ".go") && !strings.HasSuffix(fn, ".go.txt") {
	// 		fmt.Fprintf(os.Stderr, "error: bambam input file '%s' did not end in '.go' or '.go.txt'.\n", fn)
	// 		os.Exit(1)
	// 	}
	// }

	x := parser.NewExtractor()
	x.FieldPrefix = "   "
	x.FieldSuffix = "\n"
	x.OutDir = *outdir
	if privs != nil {
		x.ExtractPrivate = *privs
	}
	if overwrite != nil {
		x.Overwrite = *overwrite
	}

	// for _, inFile := range inputFiles {
	_, err := x.ExtractStructsFromOneFile(nil, inFile)
	if err != nil {
		panic(err)
	}
	// }

	// get rid of default tmp dir
	x.CompileDir.Cleanup()

	x.CompileDir.DirPath = *outdir
	x.PkgName = *pkg

	// translator library of go functions is separate from the schema
	translateFn := x.CompileDir.DirPath + "/translateCapn.go"
	translatorFile, err := os.Create(translateFn)
	if err != nil {
		panic(err)
	}
	defer translatorFile.Close()
	fmt.Fprintf(translatorFile, `package %s

import (
  capn "github.com/glycerine/go-capnproto"
  "io"
)

`, x.PkgName)

	_, err = x.WriteToTranslators(translatorFile)
	if err != nil {
		panic(err)
	}
}
