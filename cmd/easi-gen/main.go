package main

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
)

// Version 用于注入版本号
var Version = "0.0.1"

func doProcess(source io.Reader, output io.Writer, pkgName string) error {
	// read cadence content
	content := bytes.NewBuffer([]byte{})
	_, err := io.Copy(content, source)
	if err != nil {
		return err
	}
	// output to writer
	_, err = io.Copy(output, content)
	if err != nil {
		return err
	}
	return nil
}

func runCommand(cmd *cobra.Command, args []string) {
	source, _ := cmd.Flags().GetString("source")
	output, _ := cmd.Flags().GetString("output")
	pkgName, _ := cmd.Flags().GetString("pkg-name")

	var sourceReader io.ReadCloser
	// fallback to stdin
	if source == "" {
		sourceReader = os.Stdin
	}
	defer sourceReader.Close()
	var outputWriter io.WriteCloser
	// fallback to stdout
	if output == "" {
		outputWriter = os.Stdout
	}
	defer outputWriter.Close()

	// do process
	err := doProcess(sourceReader, outputWriter, pkgName)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
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
