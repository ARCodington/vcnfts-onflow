package test

import (
	"io/ioutil"
	"regexp"
	"strings"
	"testing"

	"github.com/onflow/cadence"
	emulator "github.com/onflow/flow-emulator"
	"github.com/onflow/flow-emulator/types"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Contracts struct {
	FungibleTokenAddress    flow.Address
	FUSDAddress           	flow.Address
	FUSDSigner            	crypto.Signer
	NonFungibleTokenAddress	flow.Address
	VictoryItemsAddress     flow.Address
	VictoryItemsSigner      crypto.Signer
	VictoryMarketAddress    flow.Address
	VictoryMarketSigner     crypto.Signer
}

var (
	FungibleTokenAddressPlaceholder    	= regexp.MustCompile(`"[^"\s].*/FungibleToken.cdc"`)
	FUSDAddressPlaceholder           	= regexp.MustCompile(`"[^"\s].*/FUSD.cdc"`)
	NonFungibleTokenAddressPlaceholder	= regexp.MustCompile(`"[^"\s].*/NonFungibleToken.cdc"`)
	VictoryItemsAddressPlaceholder      = regexp.MustCompile(`"[^"\s].*/VictoryCollectible.cdc"`)
	VictoryMarketAddressPlaceholder    	= regexp.MustCompile(`"[^"\s].*/VictoryCollectibleSaleOffer.cdc"`)
)

// NewBlockchain returns a new emulated blockchain.
func NewBlockchain() *emulator.Blockchain {
	b, err := emulator.NewBlockchain()
	if err != nil {
		panic(err)
	}
	return b
}

// SignAndSubmit signs a transaction with an array of signers and adds their signatures to the transaction
// Then submits the transaction to the emulator. If the private keys don't match up with the addresses,
// the transaction will not succeed.
// shouldRevert parameter indicates whether the transaction should fail or not
// This function asserts the correct result and commits the block if it passed
func SignAndSubmit(
	t *testing.T,
	b *emulator.Blockchain,
	tx *flow.Transaction,
	signerAddresses []flow.Address,
	signers []crypto.Signer,
	shouldRevert bool,
) *types.TransactionResult {
	// sign transaction with each signer
	for i := len(signerAddresses) - 1; i >= 0; i-- {
		signerAddress := signerAddresses[i]
		signer := signers[i]

		if i == 0 {
			err := tx.SignEnvelope(signerAddress, 0, signer)
			assert.NoError(t, err)
		} else {
			err := tx.SignPayload(signerAddress, 0, signer)
			assert.NoError(t, err)
		}
	}

	return Submit(t, b, tx, shouldRevert)
}

// Submit submits a transaction and checks if it succeeds.
func Submit(
	t *testing.T,
	b *emulator.Blockchain,
	tx *flow.Transaction,
	shouldRevert bool,
) *types.TransactionResult {
	// submit the signed transaction
	err := b.AddTransaction(*tx)
	require.NoError(t, err)

	result, err := b.ExecuteNextTransaction()
	require.NoError(t, err)

	if shouldRevert {
		assert.True(t, result.Reverted())
	} else {
		if !assert.True(t, result.Succeeded()) {
			t.Log(result.Error.Error())
		}
	}

	_, err = b.CommitBlock()
	assert.NoError(t, err)

	return result
}

// ExecuteScriptAndCheck executes a script and checks if it succeeds.
func ExecuteScriptAndCheck(t *testing.T, b *emulator.Blockchain, script []byte, arguments [][]byte) cadence.Value {
	result, err := b.ExecuteScript(script, arguments)
	require.NoError(t, err)

	assert.NoError(t, result.Error)

	return result.Value
}

// ExecuteScriptAndCheck executes a script and checks if it succeeds.
func ExecuteScriptAndCheckUInt64Array(t *testing.T, b *emulator.Blockchain, script []byte, arguments [][]byte) []cadence.Value {
	result, err := b.ExecuteScript(script, arguments)
	require.NoError(t, err)

	assert.NoError(t, result.Error)

	values, err := cadence.CastToArray(result.Value)
	require.NoError(t, err)

	cvalues := make([]cadence.Value, len(values))
	for i, v := range values {
		cvalues[i] = cadence.NewUInt64(v.(uint64))
	}

	return cvalues
}

// ExecuteScriptAndCheck executes a script and checks if it succeeds.
func ExecuteScriptAndCheckByteArray(t *testing.T, b *emulator.Blockchain, script []byte, arguments [][]byte) []byte {
	result, err := b.ExecuteScript(script, arguments)
	require.NoError(t, err)

	assert.NoError(t, result.Error)

	values, err := cadence.CastToArray(result.Value)
	require.NoError(t, err)

	cvalues := make([]byte, len(values))
	for i, v := range values {
		cvalues[i] = v.(byte)
	}

	return cvalues
}

// ReadFile reads a file from the filesystem.
func ReadFile(path string) []byte {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return contents
}

// CadenceUFix64 returns a Cadence UFix64 value.
func CadenceUFix64(value string) cadence.Value {
	newValue, err := cadence.NewUFix64(value)

	if err != nil {
		panic(err)
	}

	return newValue
}

func ReplaceImports(
	code string,
	importReplacements map[string]*regexp.Regexp,
) string {
	for address, find := range importReplacements {
		if !strings.Contains(address, "0x") {
			address = "0x" + address
		}

		code = find.ReplaceAllString(code, address)
	}
	return code
}

func CreateAccount(t *testing.T, b *emulator.Blockchain) (flow.Address, crypto.Signer, *flow.AccountKey) {
	accountKeys := test.AccountKeyGenerator()
	accountKey, signer := accountKeys.NewWithSigner()
	address, err := b.CreateAccount([]*flow.AccountKey{accountKey}, nil)
	require.NoError(t, err)
	return address, signer, accountKey
}
