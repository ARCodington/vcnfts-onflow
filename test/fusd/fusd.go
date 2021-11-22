package fusd

import (
	"regexp"
	"testing"
	"github.com/onflow/cadence"
	emulator "github.com/onflow/flow-emulator"
	ftcontracts "github.com/onflow/flow-ft/lib/go/contracts"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/templates"
	sdktest "github.com/onflow/flow-go-sdk/test"
	"github.com/stretchr/testify/assert"

	"github.com/ARCodington/vcnfts-onflow/test/test"
)

const (
	fusdRootPath           = ".."
	fusdContractPath       = fusdRootPath + "/contracts/FUSD.cdc"
	fusdSetupAccountPath   = fusdRootPath + "/transactions/fusd/setup_account.cdc"
	fusdTransferTokensPath = fusdRootPath + "/transactions/fusd/transfer_tokens.cdc"
	fusdMintTokensPath     = fusdRootPath + "/transactions/fusd/mint_tokens.cdc"
	fusdGetBalancePath     = fusdRootPath + "/scripts/fusd/get_balance.cdc"
	fusdGetSupplyPath      = fusdRootPath + "/scripts/fusd/get_supply.cdc"
)

func DeployContracts(t *testing.T, b *emulator.Blockchain) (flow.Address, flow.Address, crypto.Signer) {
	accountKeys := sdktest.AccountKeyGenerator()

	// Should be able to deploy a contract as a new account with no keys.
	fungibleTokenCode := ftcontracts.FungibleToken()
	fungibleTokenAddress, err := b.CreateAccount(
		[]*flow.AccountKey{},
		[]templates.Contract{templates.Contract{
			Name:   "FungibleToken",
			Source: string(fungibleTokenCode),
		}},
	)
	assert.NoError(t, err)

	_, err = b.CommitBlock()
	assert.NoError(t, err)

	fusdAccountKey, fusdSigner := accountKeys.NewWithSigner()
	fusdCode := LoadFUSD(fungibleTokenAddress)

	fusdAddress, err := b.CreateAccount(
		[]*flow.AccountKey{fusdAccountKey},
		[]templates.Contract{templates.Contract{
			Name:   "FUSD",
			Source: string(fusdCode),
		}},
	)
	assert.NoError(t, err)

	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// Simplify testing by having the contract address also be our initial Vault.
	SetupAccount(t, b, fusdAddress, fusdSigner, fungibleTokenAddress, fusdAddress)

	return fungibleTokenAddress, fusdAddress, fusdSigner
}

func SetupAccount(
	t *testing.T,
	b *emulator.Blockchain,
	userAddress flow.Address,
	userSigner crypto.Signer,
	fungibleAddress flow.Address,
	fusdAddress flow.Address,
) {
	tx := flow.NewTransaction().
		SetScript(SetupFUSDAccountTransaction(fungibleAddress, fusdAddress)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(userAddress)

	test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, userAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), userSigner},
		false,
	)
}

func CreateAccount(
	t *testing.T,
	b *emulator.Blockchain,
	fungibleTokenAddress flow.Address,
	fusdAddress flow.Address,
) (flow.Address, crypto.Signer) {
	userAddress, userSigner, _ := test.CreateAccount(t, b)
	SetupAccount(t, b, userAddress, userSigner, fungibleTokenAddress, fusdAddress)
	return userAddress, userSigner
}

func Mint(
	t *testing.T,
	b *emulator.Blockchain,
	fungibleTokenAddress flow.Address,
	fusdAddress flow.Address,
	fusdSigner crypto.Signer,
	recipientAddress flow.Address,
	amount string,
	shouldRevert bool,
) {
	tx := flow.NewTransaction().
		SetScript(MintFUSDTransaction(fungibleTokenAddress, fusdAddress)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(fusdAddress)

	_ = tx.AddArgument(cadence.NewAddress(recipientAddress))
	_ = tx.AddArgument(test.CadenceUFix64(amount))

	test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, fusdAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), fusdSigner},
		shouldRevert,
	)

}

func replaceAddressPlaceholders(code string, fungibleTokenAddress, fusdAddress string) []byte {
	return []byte(test.ReplaceImports(
		code,
		map[string]*regexp.Regexp{
			fungibleTokenAddress: test.FungibleTokenAddressPlaceholder,
			fusdAddress:        test.FUSDAddressPlaceholder,
		},
	))
}

func LoadFUSD(fungibleTokenAddress flow.Address) []byte {
	return []byte(test.ReplaceImports(
		string(test.ReadFile(fusdContractPath)),
		map[string]*regexp.Regexp{
			fungibleTokenAddress.String(): test.FungibleTokenAddressPlaceholder,
		},
	))
}

func GetSupplyScript(fungibleTokenAddress, fusdAddress flow.Address) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(fusdGetSupplyPath)),
		fungibleTokenAddress.String(),
		fusdAddress.String(),
	)
}

func GetBalanceScript(fungibleTokenAddress, fusdAddress flow.Address) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(fusdGetBalancePath)),
		fungibleTokenAddress.String(),
		fusdAddress.String(),
	)
}
func TransferVaultScript(fungibleTokenAddress, fusdAddress flow.Address) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(fusdTransferTokensPath)),
		fungibleTokenAddress.String(),
		fusdAddress.String(),
	)
}

func SetupFUSDAccountTransaction(fungibleTokenAddress, fusdAddress flow.Address) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(fusdSetupAccountPath)),
		fungibleTokenAddress.String(),
		fusdAddress.String(),
	)
}

func MintFUSDTransaction(fungibleTokenAddress, fusdAddress flow.Address) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(fusdMintTokensPath)),
		fungibleTokenAddress.String(),
		fusdAddress.String(),
	)
}
