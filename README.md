# Easinteraction For Cadence(WIP)
[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]

Easinteraction is a tool that help users to generate code for easier contract interaction.

This version is for `Cadence(Flow Blockchain)` and `Golang`.
## Get Started
### Installation
```shell
$ go install github.com/LemonNekoGH/easinteraction-for-cadence@latest
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
	"context"
	"fmt"
	"github.com/LemonNekoGH/godence"
	"github.com/onflow/cadence"
	flowSdk "github.com/onflow/flow-go-sdk"
	flowGrpc "github.com/onflow/flow-go-sdk/access/grpc"
	flowCrypto "github.com/onflow/flow-go-sdk/crypto"
)

type ContractUserProfiles struct {
	address string
	flowCli *flowGrpc.Client
}

// NewUserProfilesContract construct a new ContractUserProfiles instance.
func NewUserProfilesContract(address string, flowCli *flowGrpc.Client) (*ContractUserProfiles, error) {
	// prepare script
	script := fmt.Sprintf(`import UserProfiles from %s
pub fun main(){}
`, address)
	// send a script to ensure contract address is correct
	_, err := flowCli.ExecuteScriptAtLatestBlock(context.Background(), []byte(script), nil)
	if err != nil {
		return nil, err
	}
	return &ContractUserProfiles{
		address: address,
		flowCli: flowCli,
	}, nil
}

// SetName will change Flow blockchain's state.
// Signature: pub fun setName(user acc: AuthAccount,to name: String)
func (c *ContractUserProfiles) SetName(
	arg1 string,
	authorizer0, payer, proposer flowSdk.Address,
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
transaction(arg1: String) {
    prepare(arg0: AuthAccount) {
        UserProfiles.setName(user:arg0,to:arg1)
    }
}
`, c.address)
	// construct tx
	tx := flowSdk.NewTransaction().
		SetScript([]byte(script)).
		AddAuthorizer(authorizer0).
		SetPayer(payer).
		SetProposalKey(proposer, proposalKeyIndex, propAcct.Keys[proposalKeyIndex].SequenceNumber).
		SetReferenceBlockID(block.ID).
		SetGasLimit(9999)
	// add argument
	argCdc1, err := godence.ToCadence(
		arg1,
	)
	if err != nil {
		return nil, err
	}
	err = tx.AddArgument(argCdc1)
	if err != nil {
		return nil, err
	}

	if authorizer0.String() != payer.String() {
		// sign payload
		err = tx.SignPayload(authorizer0, authorizerKeyIndex, authorizerSigner)
		if err != nil {
			return nil, err
		}
	}
	// sign envelop
	err = tx.SignEnvelope(payer, payerKeyIndex, payerSigner)
	if err != nil {
		return nil, err
	}
	// send
	err = c.flowCli.SendTransaction(context.Background(), *tx)
	if err != nil {
		return nil, err
	}
	id := tx.ID()
	return &id, nil
}

// GetName will query from Flow blockchain.
// Signature: pub fun getName(_ addr: Address): String
func (c *ContractUserProfiles) GetName(
	arg0 string,

) (*string, error) {
	// gen script
	script := fmt.Sprintf(`import UserProfiles from %s
pub fun main(arg0: Address): String{
    return UserProfiles.getName(arg0)
}
`, c.address)
	// prepare args
	args := []cadence.Value{}
	argCdc0, err := godence.ToCadence(
		godence.Address(arg0),
	)
	if err != nil {
		return nil, err
	}
	args = append(args, argCdc0)
	// send query
	result, err := c.flowCli.ExecuteScriptAtLatestBlock(context.Background(), []byte(script), args)
	if err != nil {
		return nil, err
	}

	// covert result
	var ret string
	err = godence.ToGo(result, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil

}

```
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/LemonNekoGH/easinteraction-for-cadence.svg
[contributors-url]: https://github.com/LemonNekoGH/easinteraction-for-cadence/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/LemonNekoGH/easinteraction-for-cadence.svg
[forks-url]: https://github.com/LemonNekoGH/easinteraction-for-cadence/network/members
[stars-shield]: https://img.shields.io/github/stars/LemonNekoGH/easinteraction-for-cadence.svg
[stars-url]: https://github.com/LemonNekoGH/easinteraction-for-cadence/stargazers
[issues-shield]: https://img.shields.io/github/issues/LemonNekoGH/easinteraction-for-cadence.svg
[issues-url]: https://github.com/LemonNekoGH/easinteraction-for-cadence/issues
[license-shield]: https://img.shields.io/github/license/LemonNekoGH/easinteraction-for-cadence.svg
[license-url]: https://github.com/othneildrew/

### Testing
#### Requirements
- [Flow CLI](https://docs.onflow.org/flow-cli/): Use to emulate Flow blockchain network.
