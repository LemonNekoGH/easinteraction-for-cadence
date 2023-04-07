package main

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed internal/gen/UserProfiles.cdc
	userProfilesCdc []byte
	//go:embed internal/gen/templates/main.gohtml
	mainGo string
	//go:embed internal/gen/templates/flow.json
	flowJson string
)

func Test_runCommand(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	dir, err := os.Getwd()
	r.Empty(err)
	dir += string(os.PathSeparator) + "testing"
	r.Empty(os.MkdirAll(dir, 0750))
	// clean up
	defer func() {
		err = os.RemoveAll(dir)
		r.Empty(err)
	}()
	// write sources
	mainContent := strings.NewReader(mainGo)
	m, err := os.Create(dir + string(filepath.Separator) + "main.go")
	defer func(m *os.File) {
		err := m.Close()
		r.Empty(err)
	}(m)
	r.Empty(err)
	_, err = io.Copy(m, mainContent)
	r.Empty(err)
	// write contract file
	c, err := os.Create(dir + string(filepath.Separator) + "UserProfiles.cdc")
	defer func(c *os.File) {
		err := c.Close()
		r.Empty(err)
	}(c)
	r.Empty(err)
	_, err = io.Copy(c, bytes.NewBuffer(userProfilesCdc))
	r.Empty(err)
	// do generate
	err = runCommand0(c.Name(), dir+string(filepath.Separator)+"contracts"+string(filepath.Separator)+"user_profiles.go", "contracts")
	r.Empty(err)
	// write flow config
	fc, err := os.Create(dir + string(filepath.Separator) + "flow.json")
	defer func(fc *os.File) {
		err := fc.Close()
		r.Empty(err)
	}(fc)
	r.Empty(err)
	cfg := strings.NewReader(flowJson)
	_, err = io.Copy(fc, cfg)
	r.Empty(err)
	// download dependencies
	out, err := exec.Command("bash", "-c", fmt.Sprintf("cd %s && go mod init example", dir)).CombinedOutput()
	fmt.Printf("%s\n", out)
	r.Empty(err)
	depMods := []string{
		"github.com/onflow/flow-go-sdk/access/grpc",
		"github.com/onflow/flow-go-sdk",
		"github.com/onflow/flow-go-sdk/crypto",
		"github.com/onflow/cadence",
		"github.com/lemonnekogh/godence",
		"google.golang.org/grpc",
	}
	for _, d := range depMods {
		out, err = exec.Command("bash", "-c", fmt.Sprintf("cd %s && go get %s", dir, d)).CombinedOutput()
		fmt.Printf("%s\n", out)
		r.Empty(err)
	}
	// start flow emulator
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = exec.CommandContext(ctx, "bash", "-c", fmt.Sprintf("cd %s && flow emulator", dir)).Start()
	r.Empty(err)
	time.Sleep(5 * time.Second)
	// deploy
	out, err = exec.Command("bash", "-c", fmt.Sprintf("cd %s && flow deploy --update", dir)).CombinedOutput() // if failed at this line, go to system process manager and kill it
	fmt.Printf("%s\n", out)
	r.Empty(err)
	// run test program
	out, err = exec.Command("bash", "-c", fmt.Sprintf("cd %s && go run main.go", dir)).CombinedOutput()
	fmt.Printf("%s\n", out)
	a.Empty(err)
}
