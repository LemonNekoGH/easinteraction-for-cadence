# Easinteraction For Cadence(WIP)
Easinteraction is a tool that help users to generate code for easier contract interaction.

This version is for `Cadence(Flow Blockchain)` and `Golang`.
## Get Started
### Installation
```shell
$ go install github.com/LemonNekoGH/easiteraction-for-cadence
```
### Usage
```shell
$ easi-gen --source NekosMerge.cdc --output nekosmerge.go --pkg-name nekosmerge
```
Now `easi-gen` will read your `Cadence` code from `NekosMerge.cdc`, then parse the code, then generate `Go` code with package name `nekosmerge` and output to `nekosmerge.go`. 

If no `--pkg-name` flag specified, it will use `mypackage` as default.  
If no `--source` flag specified, it will read from standard input.  
If no `--output` flag specified, it will output to standard output.

So you can also use it like this:
```shell
$ cat NekosMerge.cdc | easi-gen --pkg-name nekosmerge > nekosmerge.go
```

### Result Preview
Cadence code:
```cadence
pub contract UserProfiles {
    access(self) let usernames: {Address:String}

    pub fun setName(_ acc: AuthAccount, _ name: String) {
        self.usernames[acc.address] = name
    }

    pub fun getName(_ addr: Address): String {
        return self.usernames[addr] ?? ""
    }

    init() {
        self.usernames = {}
    }
}
```
Go code:

```go
package mypackage

import (
	"fmt"
	flowGrpc "github.com/onflow/flow-go-sdk/access/grpc"
	flowSdk "github.com/onflow/flow-go-sdk"
	flowCrypto "github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/cadence"
	"strings"
)

type ContractUserProfiles struct {
	address string
	flowCli flowGrpc.Client
}

// NewUserProfilesContract construct a new ContractUserProfiles instance.
func NewUserProfilesContract(address string, flowCli flowGrpc.Client) (*ContractUserProfiles, error) {
	// prepare script
	script := fmt.Sprintf(`import UserProfiles from %s
pub fun main(){}`, address)
	// send a script to ensure contract address is correct
	_, err := flowCli.ExecuteScriptAtLatestBlock(context.Background(), script, args)
	if err != nil {
		return nil, err
	}
	return &ContractUserProfiles{
		address: address,
		flowCli: flowCli,
	}, nil
}

// SetName will change Flow blockchain's state.
// Signature: pub fun setName(_ acc: AuthAccount, _ name: String)
func (c *ContractUserProfiles) SetName(
	name string,
	authorizer, payer, proposer flowSdk.Address,
	authorizerKeyIndex, payerKeyIndex, proposalKeyIndex int,
	authorizerSigner, payerSigner flowCrypto.Signer,
) (*flowSdk.Identifier, error) {
	// get reference id
	block, err := c.flowCli.GetLatestBlock(context.Background(), true)
	if err != nil {
		return nil, err
	}
	// get proposer sequence number
	propAcct, err := c.flowCli.GetAccount(context.Background(), proposer)
	if err != nil {
		return nil, err
	}
	// gen script
	script := fmt.Sprintf(`import UserProfiles from %s
transaction {
  prepare(acct: AuthAccount) {
    UserProfiles.setName(acct, %s)
  }
}`, c.address, name)
	// construct tx
	tx := flowSdk.NewTransaction().
		SetScript([]byte(script)).
		AddAuthorizer(authorizer).
		SetPayer(payer).
		SetProposalKey(proposer, proposalKeyIndex, propAcct.Keys[proposalKeyIndex].SequenceNumber).
		SetReferenceBlockID(block.ID).
		SetGasLimit(9999)
	// add argument
	err = tx.AddArgument(cadence.NewString(name))
	if err != nil {
		return nil, err
	}
	// sign payload
	err = tx.SignPayload(authorizer, authorizerKeyIndex, authorizerSigner)
	if err != nil {
		return nil, err
	}
	// sign envelop
	err = tx.SignEnvelop(payer, payerKeyIndex, payerSigner)
	if err != nil {
		return nil, err
	}
	// send
	err = c.flowCli.SendTransaction(context.Background(), *tx)
	if err != nil {
		return nil, err
	}
	id := tx.ID
	return &id, nil
}

// GetName will query from Flow blockchain.
// Signature: pub fun getName(_ addr: Address): String
func (c *ContractUserProfiles) GetName(addr flowSdk.Address) (string, error) {
	// gen script
	script := fmt.Sprintf(`import UserProfiles from %s
fun main(_ addr: Address): String {
    return UserProfiles.getName(addr)
}`, c.address)
	// prepare args
	arg0I, err := hex.DecodeString(strings.TrimPrefix(addr.String(), "0x"))
	if err != nil {
		return "", err
	}
	arg0 := cadence.BytesToAddress(arg0I)
	args := []cadence.Value{arg0}
	// send query
	result, err := c.flowCli.ExecuteScriptAtLatestBlock(context.Background(), script, args)
	if err != nil {
		return "", err
	}
	// covert result
	ret, ok := result.ToGoValue().(string)
	if !ok {
		return "", fmt.Errorf("type error, expect: string, but got: %T", result.ToGoValue())
	}
	return ret, nil
}
```
