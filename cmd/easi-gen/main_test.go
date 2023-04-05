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

//go:embed internal/gen/UserProfiles.cdc
var userProfilesCdc []byte

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
	mainContent := strings.NewReader(`package main
import (
	"context"
	flowGrpc "github.com/onflow/flow-go-sdk/access/grpc"
	flowSdk "github.com/onflow/flow-go-sdk"
	flowCrypto "github.com/onflow/flow-go-sdk/crypto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"example/contracts"
	"time"
)

var flowCli *flowGrpc.Client
func initFlowClient() {
	client, err := flowGrpc.NewClient(
		"localhost:3569",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	err = client.Ping(context.Background())
	if err != nil {
		panic(err)
	}
	flowCli = client
}
func main() {
	initFlowClient()
	c, err := contracts.NewUserProfilesContract("0xf8d6e0586b0a20c7", flowCli)
	if err != nil {
		panic(err)
	}
	// get signer
	privateKey, err := flowCrypto.DecodePrivateKeyHex(flowCrypto.ECDSA_P256, "c47db93881bc34a6155192c2bec0d124731e08ff105672afdb09892e3dc9ccae")
	if err != nil {
		panic(err)
	}
	signer, err := flowCrypto.NewInMemorySigner(privateKey, flowCrypto.SHA3_256)
	if err != nil {
		panic(err)
	}
	// set name
	addr := flowSdk.HexToAddress("0xf8d6e0586b0a20c7")
	_, err = c.SetName("LemonNeko", addr, addr, addr, 0, 0, 0, signer, signer)
	if err != nil {
		panic(err)
	}
	time.Sleep(5*time.Second)
	// get name
	name, err := c.GetName("0xf8d6e0586b0a20c7")
	if err != nil {
		panic(err)
	}
	if name != "LemonNeko" {
		panic(err)
	}
	// set avatar
	_, err = c.SetAvatar("ForTwitter", "https://example.com/avatars/lemonneko", addr, addr, addr, 0, 0, 0, signer, signer)
	if err != nil {
		panic(err)
	}
	// get avatar
	avatars, err := c.GetAllAvatars("0xf8d6e0586b0a20c7")
	if err != nil {
		panic(err)
	}
	if avatars["ForTwitter"] != "https://example.com/avatars/lemonneko" {
		panic(err)
	}
	if len(avatars) != 1 {
		panic(err)
	}
	// get avatar names
	avatarNames, err := c.GetAllAvatarNames("0xf8d6e0586b0a20c7")
	if err != nil {
		panic(err)
	}
	if avatarNames[0] != "ForTwitter" {
		panic(err)
	}
	if len(avatarNames) != 1 {
		panic(err)
	}
	// get not exists avatar
	notExists, err := c.GetAvatarByName("0xf8d6e0586b0a20c7", "NotExists")
	if err != nil {
		panic(err)
	}
	if notExists != "" {
		panic(err)
	}
}
`)
	m, err := os.Create(dir + string(filepath.Separator) + "main.go")
	defer m.Close()
	r.Empty(err)
	_, err = io.Copy(m, mainContent)
	r.Empty(err)
	// write contract file
	c, err := os.Create(dir + string(filepath.Separator) + "UserProfiles.cdc")
	defer c.Close()
	r.Empty(err)
	_, err = io.Copy(c, bytes.NewBuffer(userProfilesCdc))
	r.Empty(err)
	// do generate
	err = runCommand0(c.Name(), dir+string(filepath.Separator)+"contracts"+string(filepath.Separator)+"user_profiles.go", "contracts")
	r.Empty(err)
	// write flow config
	fc, err := os.Create(dir + string(filepath.Separator) + "flow.json")
	defer fc.Close()
	r.Empty(err)
	cfg := strings.NewReader(`{
	"emulators": {
		"default": {
			"port": 3569,
			"serviceAccount": "emulator-account"
		}
	},
	"contracts": {
		"UserProfiles": "./UserProfiles.cdc"
	},
	"networks": {
		"emulator": "127.0.0.1:3569",
		"mainnet": "access.mainnet.nodes.onflow.org:9000",
		"testnet": "access.devnet.nodes.onflow.org:9000"
	},
	"accounts": {
		"emulator-account": {
			"address": "f8d6e0586b0a20c7",
			"key": "c47db93881bc34a6155192c2bec0d124731e08ff105672afdb09892e3dc9ccae"
		}
	},
	"deployments": {
		"emulator": {
			"emulator-account": ["UserProfiles"]
		}
	}
}`)
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
