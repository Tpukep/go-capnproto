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
	outdir       = flag.String("o", ".", "specify output directory")
	source       = flag.String("source", "", "specify output directory")
	verbose      = flag.Bool("verbose", false, "verbose mode")
	structTypeRe = regexp.MustCompile("struct ([[:alpha:]]+)")
)

const (
	CAPNP_CODEC_ANT = "$Codec.capnp;"
	MSGP_CODEC_ANT  = "$Codec.msgp;"
)

func use() {
	fmt.Fprintf(os.Stderr, "\nuse: caps -o <outdir> -source=<model.capnp>\n")
	fmt.Fprintf(os.Stderr, "     # Tool reads .capnp files and writes: go structs with json tags, capn'proto code, translation code, msgp code.\n")
	fmt.Fprintf(os.Stderr, "     # options:\n")
	fmt.Fprintf(os.Stderr, "     #   -o=\"outdir\" specifies the directory to write to (created if need be).\n")
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

	capsSchemaPath := fmt.Sprintf("-I%s/src/github.com/tpukep/go-capnproto", gopath)
	goSchemaPath := fmt.Sprintf("-I%s/src/github.com/tpukep/go-capnproto/vendor/github.com/glycerine/go-capnproto", gopath)

	capnpArgs := []string{capsSchemaPath, goSchemaPath, "compile", "-opgo"}

	sourceData, err := ioutil.ReadFile(source)
	if err != nil {
		log.Fatalln("Failed to read schema file:", err)
	}

	sourceContent := string(sourceData)
	// Remove comments
	re := regexp.MustCompile("(?m)[\r\n]+^.*#.*$")
	cleanContent := re.ReplaceAllString(sourceContent, "")

	// Find codec annotations
	capnpRe := regexp.MustCompile("(?m)[\r\n]+^.*" + regexp.QuoteMeta(CAPNP_CODEC_ANT) + ".*$")
	msgpRe := regexp.MustCompile("(?m)[\r\n]+^.*" + regexp.QuoteMeta(MSGP_CODEC_ANT) + ".*$")

	capnpGen := capnpRe.MatchString(cleanContent)
	msgpGen := msgpRe.MatchString(cleanContent)

	if capnpGen {
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

	err = cmd.Run()
	if err != nil {
		log.Fatalln("Failed to run Plain go code generator:", err)
	}

	if capnpGen {
		// Add suffix "Capn" to Capn'proto structs
		outFilename := filepath.Join(*outdir, source[:strings.LastIndex(source, ".")]+".capnp.go")

		data, err := ioutil.ReadFile(outFilename)
		if err != nil {
			log.Fatalln("Failed to read Capn'proto out file:", err)
		}

		content := string(data)
		matches := structTypeRe.FindAllStringSubmatch(sourceContent, -1)
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
	if msgpGen {
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
