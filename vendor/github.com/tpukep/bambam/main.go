package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/tpukep/bambam/bam"
)

func use() {
	fmt.Fprintf(os.Stderr, "\nuse: bambam -o outdir -p package myGoSourceFile.go myGoSourceFile2.go ...\n")
	fmt.Fprintf(os.Stderr, "     # Bambam makes it easy to use Capnproto serialization[1] from Go.\n")
	fmt.Fprintf(os.Stderr, "     # Bambam reads .go files and writes a .capnp schema and Go bindings.\n")
	fmt.Fprintf(os.Stderr, "     # options:\n")
	fmt.Fprintf(os.Stderr, "     #   -o=\"odir\" specifies the directory to write to (created if need be).\n")
	fmt.Fprintf(os.Stderr, "     #   -p=\"main\" specifies the package header to write (e.g. main, mypkg).\n")
	fmt.Fprintf(os.Stderr, "     #   -X exports private fields of Go structs. Default only maps public fields.\n")
	fmt.Fprintf(os.Stderr, "     #   -version   shows build version with git commit hash.\n")
	fmt.Fprintf(os.Stderr, "     #   -debug     print lots of debug info as we process.\n")
	fmt.Fprintf(os.Stderr, "     #   -OVERWRITE modify .go files in-place, adding capid tags (write to -o dir by default).\n")
	fmt.Fprintf(os.Stderr, "     # required: at least one .go source file for struct definitions. Must be last, after options.\n")
	fmt.Fprintf(os.Stderr, "     #\n")
	fmt.Fprintf(os.Stderr, "     # [1] https://github.com/glycerine/go-capnproto \n")
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func main() {
	MainArgs(os.Args)
}

// allow invocation from test
func MainArgs(args []string) {
	//fmt.Println(os.Args)
	os.Args = args

	flag.Usage = use
	if len(os.Args) < 2 {
		use()
	}

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	debug := flag.Bool("debug", false, "print lots of debug info as we process.")
	verrequest := flag.Bool("version", false, "request git commit hash used to build this bambam")
	outdir := flag.String("o", "odir", "specify output directory")
	pkg := flag.String("p", "main", "specify package for generated code")
	privs := flag.Bool("X", false, "export private as well as public struct fields")
	overwrite := flag.Bool("OVERWRITE", false, "replace named .go files with capid tagged versions.")
	flag.Parse()

	if debug != nil {
		bam.Verbose = *debug
	}

	if verrequest != nil && *verrequest {
		fmt.Printf("%s\n", LASTGITCOMMITHASH)
		os.Exit(0)
	}

	if outdir == nil || *outdir == "" {
		fmt.Fprintf(os.Stderr, "required -o option missing. Use bambam -o <dirname> myfile.go # to specify the output directory.\n")
		use()
	}

	if !bam.DirExists(*outdir) {
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
	inputFiles := flag.Args()

	if len(inputFiles) == 0 {
		fmt.Fprintf(os.Stderr, "bambam needs at least one .go golang source file to process specified on the command line.\n")
		os.Exit(1)
	}

	for _, fn := range inputFiles {
		if !strings.HasSuffix(fn, ".go") && !strings.HasSuffix(fn, ".go.txt") {
			fmt.Fprintf(os.Stderr, "error: bambam input file '%s' did not end in '.go' or '.go.txt'.\n", fn)
			os.Exit(1)
		}
	}

	x := bam.NewExtractor()
	x.FieldPrefix = "   "
	x.FieldSuffix = "\n"
	x.OutDir = *outdir
	if privs != nil {
		x.ExtractPrivate = *privs
	}
	if overwrite != nil {
		x.Overwrite = *overwrite
	}

	for _, inFile := range inputFiles {
		_, err := x.ExtractStructsFromOneFile(nil, inFile)
		if err != nil {
			panic(err)
		}
	}
	// get rid of default tmp dir
	x.CompileDir = bam.NewTempDir()
	x.CompileDir.Cleanup()

	x.CompileDir.DirPath = *outdir
	x.PkgName = *pkg

	schemaFN := x.CompileDir.DirPath + "/schema.capnp"
	schemaFile, err := os.Create(schemaFN)
	if err != nil {
		panic(err)
	}
	defer schemaFile.Close()

	by := x.GenCapnpHeader()
	schemaFile.Write(by.Bytes())

	_, err = x.WriteToSchema(schemaFile)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(schemaFile, "\n")
	fmt.Fprintf(schemaFile, "##compile with:\n\n##\n##\n##   capnp compile -ogo %s\n\n", schemaFN)

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

	err = x.CopySourceFilesAddCapidTag()
	if err != nil {
		panic(err)
	}

	exec.Command("cp", "-p", "go.capnp", x.CompileDir.DirPath).Run()
	fmt.Printf("generated files in '%s'\n", x.CompileDir.DirPath)
}
