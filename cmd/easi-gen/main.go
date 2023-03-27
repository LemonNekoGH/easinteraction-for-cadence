package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/LemonNekoGH/easiteraction-for-cadence/cmd/easi-gen/internal/gen"
	"github.com/onflow/cadence/runtime/parser"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path/filepath"
)

// Version 用于注入版本号
var Version = "0.0.1"

func doProcess(source io.Reader, output io.Writer, pkgName string) error {
	// read cadence content
	sInput := bytes.NewBuffer([]byte{})
	_, err := io.Copy(sInput, source)
	if err != nil {
		return err
	}
	// parse cadence content
	cdc, err := parser.ParseProgram(nil, sInput.Bytes(), parser.Config{})
	if err != nil {
		return err
	}
	// gen golang code
	g := gen.NewGenerator(pkgName)
	if err = g.Gen(cdc); err != nil {
		return err
	}
	// output to writer
	_, err = io.Copy(output, g.GetOutput())
	if err != nil {
		return err
	}
	return nil
}

func runCommand(cmd *cobra.Command, _ []string) {
	source, _ := cmd.Flags().GetString("source")
	output, _ := cmd.Flags().GetString("output")
	pkgName, _ := cmd.Flags().GetString("pkg-name")
	err := runCommand0(source, output, pkgName)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func runCommand0(source, output, pkgName string) error {
	var sourceReader io.ReadCloser
	// fallback to stdin
	if source == "" {
		sourceReader = os.Stdin
	} else {
		// check source exists
		if _, err := os.Stat(source); err != nil {
			return err
		}
		var err error
		sourceReader, err = os.Open(source)
		if err != nil {
			return err
		}
	}
	if sourceReader == nil {
		return errors.New("get source reader error")
	}
	defer sourceReader.Close()
	var outputWriter io.WriteCloser
	// fallback to stdout
	if output == "" {
		outputWriter = os.Stdout
	} else {
		// check source exists
		if of, err := os.Stat(output); err != nil {
			// check parent dir exists
			outDir := filepath.Dir(output)
			if baseInfo, err2 := os.Stat(outDir); err2 != nil {
				// create
				err2 = os.MkdirAll(outDir, 0755)
				if err2 != nil {
					return err2
				}
				outputWriter, err2 = os.Create(output)
				if err2 != nil {
					return err2
				}
			} else if !baseInfo.IsDir() {
				return errors.New("the parent path of the output should be a directory, not a file")
			}
		} else if of.IsDir() {
			return errors.New("the path of the output should be a file, not a directory")
		} else {
			var err2 error
			// open file as r/w mode
			outputWriter, err2 = os.OpenFile(output, os.O_RDWR, 0755)
			if err2 != nil {
				return err2
			}
		}
	}
	if outputWriter == nil {
		return errors.New("get output reader error")
	}
	defer outputWriter.Close()

	// do process
	return doProcess(sourceReader, outputWriter, pkgName)
}

func main() {
	cmd := cobra.Command{
		Use:     "easi-gen",
		Version: Version,
		Run:     runCommand,
	}
	cmd.Flags().StringP("source", "s", "", "Specify Cadence source file path.")
	cmd.Flags().StringP("output", "o", "", "Specify output go source file path.")
	cmd.Flags().StringP("pkg-name", "p", "mypackage", "Specify go package name.")

	if err := cmd.Execute(); err != nil {
		_, _ = os.Stderr.Write([]byte(err.Error()))
		os.Exit(1)
	}
}
