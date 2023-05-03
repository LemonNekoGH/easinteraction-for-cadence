package cmd_shared

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/LemonNekoGH/easinteraction-for-cadence/cmd/easi-gen/internal/gen"
	"github.com/onflow/cadence/runtime/parser"
)

func DoProcess(source io.Reader, output io.Writer, pkgName string, ignoreContractGeneration bool) error {
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
	g := gen.NewGenerator(pkgName, ignoreContractGeneration)
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

func GetOutputWriter(output string) (io.WriteCloser, error) {
	var outputWriter io.WriteCloser
	// fallback to stdout
	if output == "" {
		outputWriter = os.Stdout
		return outputWriter, nil
	}
	// check source exists
	if of, err := os.Stat(output); err != nil {
		// check parent dir exists
		outDir := filepath.Dir(output)
		if baseInfo, err2 := os.Stat(outDir); err2 != nil {
			// create
			err2 = os.MkdirAll(outDir, 0755)
			if err2 != nil {
				return nil, err2
			}
		} else if !baseInfo.IsDir() {
			return nil, errors.New("the parent path of the output should be a directory, not a file")
		}
		var err2 error
		outputWriter, err2 = os.Create(output)
		if err2 != nil {
			return nil, err2
		}
		return outputWriter, nil
	} else if of.IsDir() {
		return nil, errors.New("the path of the output should be a file, not a directory")
	}
	// open file as r/w mode
	outputWriter, err := os.OpenFile(output, os.O_RDWR, 0755)
	if err != nil {
		return nil, err
	}
	return outputWriter, nil
}

func GetSourceReader(source string) (io.ReadCloser, error) {
	// check source exists
	if _, err := os.Stat(source); err != nil {
		return nil, err
	}
	sourceReader, err := os.Open(source)
	if err != nil {
		return nil, err
	}
	return sourceReader, nil
}
