package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/tpukep/bambam/bam"
)

var (
	outdir   = flag.String("o", ".", "specify output directory")
	capnpGen = flag.Bool("capnp", true, "generate Go Capn'p code")
	msgpGen  = flag.Bool("msgp", true, "generate Msgp code")
	source   = flag.String("source", "", "source schema file")
	verbose  = flag.Bool("verbose", false, "verbose mode")
)

func use() {
	fmt.Fprintf(os.Stderr, "\nuse: go-capnproto -o <outdir> -capnp -msgp -source=<model.capnp>\n")
	fmt.Fprintf(os.Stderr, "     # Tool reads .capnp files and writes: go structs with json tags, capn'proto code, translation code, msgp code.\n")
	fmt.Fprintf(os.Stderr, "     # options:\n")
	fmt.Fprintf(os.Stderr, "     #   -o=\"outdir\" specifies the directory to write to (created if need be).\n")
	fmt.Fprintf(os.Stderr, "     #   -capnp=true specifies generate Capn'proto code or not\n")
	fmt.Fprintf(os.Stderr, "     #   -msgp=true specifies generate Msgp code or not\n")
	fmt.Fprintf(os.Stderr, "     #   -verbose=true enables verbose mode \n")
	fmt.Fprintf(os.Stderr, "     # required:\n")
	fmt.Fprintf(os.Stderr, "     #   -source=model.capnp specifies input schema file\n")
	fmt.Fprintf(os.Stderr, "     #\n")
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func main() {
	flag.Parse()

	flag.Usage = use
	if *source == "" {
		use()
	}

	source := *source
	gopath := os.Getenv("GOPATH")

	checkSchemaPath := fmt.Sprintf("-I%s/src/github.com/tpukep/go-capnproto/check", gopath)
	jsonTagSchemaPath := fmt.Sprintf("-I%s/src/github.com/tpukep/go-capnproto/jsontag", gopath)
	msgpTagSchemaPath := fmt.Sprintf("-I%s/src/github.com/tpukep/go-capnproto/msgptag", gopath)
	goSchemaPath := fmt.Sprintf("-I%s/src/github.com/tpukep/go-capnproto/vendor/github.com/glycerine/go-capnproto", gopath)

	capnpArgs := []string{checkSchemaPath, jsonTagSchemaPath, msgpTagSchemaPath, goSchemaPath, "compile", "-opgo"}
	if *capnpGen {
		capnpArgs = append(capnpArgs, "-ogo")
	}

	if *outdir != "." {
		if !bam.DirExists(*outdir) {
			err := os.MkdirAll(*outdir, 0755)
			if err != nil {
				log.Fatalln("Failed to create output dir:", err)
			}
		}

		for i := 0; i < len(capnpArgs); i++ {
			capnpArgs[i] += *outdir
		}
	}
	capnpArgs = append(capnpArgs, source)

	// Generate plain Go code
	cmd := exec.Command("capnp", capnpArgs...)
	cmd.Stderr = os.Stderr

	if *verbose {
		fmt.Printf("Executing: %q\n", strings.Join(cmd.Args, " "))
	}

	err := cmd.Run()
	if err != nil {
		log.Fatalln("Failed to run Plain go code generator:", err)
	}

	if *capnpGen {
		// Add suffix "Capn" to Capn'proto structs
		outFilename := filepath.Join(*outdir, source[:strings.LastIndex(source, ".")]+".capnp.go")
		sourceData, err := ioutil.ReadFile(source)
		if err != nil {
			log.Fatalln("Failed to read schema file:", err)
		}

		data, err := ioutil.ReadFile(outFilename)
		if err != nil {
			log.Fatalln("Failed to read Capn'proto out file:", err)
		}

		content := string(data)
		re := regexp.MustCompile("struct ([[:alpha:]]+)")
		matches := re.FindAllStringSubmatch(string(sourceData), -1)
		for _, match := range matches {
			structType := match[1]
			if structType != "" {
				content = strings.Replace(content, structType, structType+"Capn", -1)
			}
		}

		err = ioutil.WriteFile(outFilename, []byte(content), 0644)
		if err != nil {
			log.Fatalln("Failed to run write replaced Capn'proto code file:", err)
		}
	}

	// Generate Msgp code
	if *msgpGen {
		baseFilename := filepath.Base(source)
		inFilename := filepath.Join(*outdir, source[:strings.LastIndex(source, ".")]+".go")
		outFilename := filepath.Join(*outdir, baseFilename[:strings.LastIndex(baseFilename, ".")]+".msgp.go")
		cmd = exec.Command("msgp", "-o="+outFilename, "-tests=false", "-file="+inFilename)
		cmd.Stderr = os.Stderr

		if *verbose {
			fmt.Printf("Executing: %q\n", strings.Join(cmd.Args, " "))
		}

		err = cmd.Run()
		if err != nil {
			log.Fatalln("Failed to run Msgp go code generator:", err)
		}
	}
}

func copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}
