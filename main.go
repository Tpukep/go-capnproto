package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	filename := os.Args[1]

	gopath := os.Getenv("GOPATH")

	checkSchemaPath := fmt.Sprintf("-I%s/src/github.com/tpukep/go-capnproto/check", gopath)
	jsonTagSchemaPath := fmt.Sprintf("-I%s/src/github.com/tpukep/go-capnproto/jsontag", gopath)
	msgpTagSchemaPath := fmt.Sprintf("-I%s/src/github.com/tpukep/go-capnproto/msgptag", gopath)
	goSchemaPath := fmt.Sprintf("-I%s/src/github.com/tpukep/go-capnproto/vendor/github.com/glycerine/go-capnproto", gopath)

	// Generate plain Go code
	cmd := exec.Command("capnp", checkSchemaPath, jsonTagSchemaPath, msgpTagSchemaPath, goSchemaPath, "compile", "-opgo", filename)
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatal("Failed to run plain go code generator", err)
	}

	backupFilename := filename + ".orig"

	// Backup original schema
	if err = copy(filename, backupFilename); err != nil {
		log.Fatal("Failed to backup schema file: ", err)
	}

	// Add suffix "Capn" to Capn'proto structs
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("Failed to read schema file", err)
	}

	content := string(data)
	re := regexp.MustCompile("struct ([[:alpha:]]+)")
	matches := re.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if match[1] != "" {
			content = strings.Replace(content, match[1], match[1]+"Capn", -1)
		}
	}

	err = ioutil.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		log.Fatal("Failed to run capn'proto go code generator", err)
	}

	// Generate Capn'proto code
	cmd = exec.Command("capnp", checkSchemaPath, jsonTagSchemaPath, msgpTagSchemaPath, goSchemaPath, "compile", "-ogo", filename)
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Fatal("Failed to run capn'proto go code generator", err)
	}

	// Restore original schema file from backup
	err = os.Remove(filename)
	if err != nil {
		log.Fatal("Failed to remove modified file", err)
	}

	err = os.Rename(backupFilename, filename)
	if err != nil {
		log.Fatal("Failed to rename backup file", err)
	}

	// Generate Msgp code
	baseFilename := filename[:strings.LastIndex(filename, ".")]
	cmd = exec.Command("msgp", "-o="+baseFilename+".msgp.go", "-tests=false", "-file="+baseFilename+".go")
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Fatal("Failed to run capn'proto go code generator", err)
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
