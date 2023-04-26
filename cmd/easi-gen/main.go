package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/LemonNekoGH/easinteraction-for-cadence/cmd/easi-gen/internal/cmd_shared"
	"github.com/LemonNekoGH/easinteraction-for-cadence/cmd/easi-gen/internal/types"
	"github.com/spf13/cobra"
	"io"
	"os"
)

// Version 用于注入版本号
var Version = "0.0.1"

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
		var err error
		sourceReader, err = cmd_shared.GetSourceReader(source)
		if err != nil {
			return err
		}
	}
	if sourceReader == nil {
		return errors.New("get source reader error")
	}
	defer sourceReader.Close()
	// check stdin is flow.json
	sourceBuffer := bytes.NewBuffer([]byte{})
	_, err := io.Copy(sourceBuffer, sourceReader)
	if err != nil {
		return err
	}

	var (
		flowJson types.FlowJson
	)
	err = json.Unmarshal(sourceBuffer.Bytes(), &flowJson)
	if err == nil {
		if source == "" {
			return errors.New("flow.json cannot read from stdin")
		}

		sourcesPath, outputsPath := flowJson.ResolvePath(source, pkgName, output)
		// get reader and writers
		for i, s := range sourcesPath {
			inputReader, err2 := cmd_shared.GetSourceReader(s)
			if err2 != nil {
				fmt.Println("get source reader failed, skipped: " + s)
				fmt.Println("	error: " + err2.Error())
				continue
			}

			o := outputsPath[i]
			outWriter, err2 := cmd_shared.GetOutputWriter(o)
			if err2 != nil {
				fmt.Println("get output writer failed, skipped: " + o)
				fmt.Println("	error: " + err2.Error())
				continue
			}

			if outWriter == nil {
				fmt.Println("get output writer failed, skipped: " + o)
				continue
			}

			err2 = cmd_shared.DoProcess(inputReader, outWriter, pkgName)
			if err2 != nil {
				fmt.Println("process failed, skipped: " + s)
				fmt.Println("	error: " + err2.Error())
			}

			inputReader.Close()
			outWriter.Close()
		}
		return nil
	} else {
		fmt.Println("json unmarshal failed, this is not a flow.json: " + source)
		fmt.Println("	error: " + err.Error())
	}
	// not flow json file
	outputWriter, err2 := cmd_shared.GetOutputWriter(output)
	defer outputWriter.Close()
	if err2 != nil {
		return err2
	}

	if outputWriter == nil {
		return errors.New("get output writer error")
	}

	// do process
	err = cmd_shared.DoProcess(sourceBuffer, outputWriter, pkgName) // sourceReader is EOF, use sourceBuffer instead, or use io.TeeReader. https://stackoverflow.com/questions/39791021/how-to-read-multiple-times-from-same-io-reader
	if err != nil {
		fmt.Println("process failed, skipped: " + source)
		fmt.Println("	error: " + err.Error())
	}
	return err
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
