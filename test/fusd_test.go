package test

import (
	"testing"
	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/stretchr/testify/assert"

	"github.com/ARCodington/vcnfts-onflow/test/fusd"
	"github.com/ARCodington/vcnfts-onflow/test/test"
)

func TestFUSDSetupAccount(t *testing.T) {
	b := test.NewBlockchain()

	t.Run("Should be able to create empty vault that does not affect supply", func(t *testing.T) {
		fungibleAddr, fusdAddr, _ := fusd.DeployContracts(t, b)

		supply1 := test.ExecuteScriptAndCheck(
			t, b,
			fusd.GetSupplyScript(fungibleAddr, fusdAddr),
			nil,
		)

		userAddress, _ := fusd.CreateAccount(t, b, fungibleAddr, fusdAddr)

		userBalance := test.ExecuteScriptAndCheck(
			t, b,
			fusd.GetBalanceScript(fungibleAddr, fusdAddr),
			[][]byte{jsoncdc.MustEncode(cadence.Address(userAddress))},
		)
		assert.EqualValues(t, test.CadenceUFix64("0.0"), userBalance)

		supply2 := test.ExecuteScriptAndCheck(
			t, b,
			fusd.GetSupplyScript(fungibleAddr, fusdAddr),
			nil,
		)
		assert.EqualValues(t, supply1, supply2)
	})
}

func TestFUSDMinting(t *testing.T) {
	b := test.NewBlockchain()

	fungibleTokenAddress, fusdAddress, fusdSigner := fusd.DeployContracts(t, b)

	userAddress, _ := fusd.CreateAccount(t, b, fungibleTokenAddress, fusdAddress)

	t.Run("Should not be able to mint zero tokens", func(t *testing.T) {
		fusd.Mint(
			t, b,
			fungibleTokenAddress,
			fusdAddress,
			fusdSigner,
			userAddress,
			"0.0",
			true,
		)
	})

	t.Run("Should be able to mint tokens, deposit, update balance and total supply", func(t *testing.T) {

		// Assert that total supply is correct
		supply1 := test.ExecuteScriptAndCheck(
			t, b,
			fusd.GetSupplyScript(fungibleTokenAddress, fusdAddress),
			nil,
		)

		fusd.Mint(
			t, b,
			fungibleTokenAddress,
			fusdAddress,
			fusdSigner,
			userAddress,
			"50.0",
			false,
		)

		// Assert that vault balance is correct
		userBalance := test.ExecuteScriptAndCheck(
			t, b,
			fusd.GetBalanceScript(fungibleTokenAddress, fusdAddress),
			[][]byte{jsoncdc.MustEncode(cadence.Address(userAddress))},
		)

		assert.EqualValues(t, test.CadenceUFix64("50.0"), userBalance)

		// Assert that total supply is correct
		supply2 := test.ExecuteScriptAndCheck(
			t, b,
			fusd.GetSupplyScript(fungibleTokenAddress, fusdAddress),
			nil,
		)

		var target = supply1.ToGoValue().(uint64) + test.CadenceUFix64("50.0").ToGoValue().(uint64)
		var supply = supply2.ToGoValue().(uint64)
		
		assert.EqualValues(t, target, supply)
	})
}

func TestFUSDTransfers(t *testing.T) {
	b := test.NewBlockchain()

	fungibleTokenAddress, fusdAddress, fusdSigner := fusd.DeployContracts(t, b)

	userAddress, _ := fusd.CreateAccount(t, b, fungibleTokenAddress, fusdAddress)

	// Mint 1000 new FUSD into the FUSD contract account
	fusd.Mint(
		t, b,
		fungibleTokenAddress,
		fusdAddress,
		fusdSigner,
		fusdAddress,
		"1000.0",
		false,
	)

	t.Run("Should not be able to withdraw more than the balance of the vault", func(t *testing.T) {
		tx := flow.NewTransaction().
			SetScript(fusd.TransferVaultScript(fungibleTokenAddress, fusdAddress)).
			SetGasLimit(100).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(fusdAddress)

		_ = tx.AddArgument(test.CadenceUFix64("2000000.0"))
		_ = tx.AddArgument(cadence.NewAddress(userAddress))

		test.SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, fusdAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), fusdSigner},
			true,
		)

		// Assert that vault balances are correct

		fusdBalance := test.ExecuteScriptAndCheck(
			t, b,
			fusd.GetBalanceScript(fungibleTokenAddress, fusdAddress),
			[][]byte{jsoncdc.MustEncode(cadence.Address(fusdAddress))},
		)

		assert.EqualValues(t, test.CadenceUFix64("1001000.0"), fusdBalance)

		userBalance := test.ExecuteScriptAndCheck(
			t, b,
			fusd.GetBalanceScript(fungibleTokenAddress, fusdAddress),
			[][]byte{jsoncdc.MustEncode(cadence.Address(userAddress))},
		)

		assert.EqualValues(t, test.CadenceUFix64("0.0"), userBalance)
	})

	t.Run("Should be able to withdraw and deposit tokens from a vault", func(t *testing.T) {
		tx := flow.NewTransaction().
			SetScript(fusd.TransferVaultScript(fungibleTokenAddress, fusdAddress)).
			SetGasLimit(100).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(fusdAddress)

		_ = tx.AddArgument(test.CadenceUFix64("300.0"))
		_ = tx.AddArgument(cadence.NewAddress(userAddress))

		test.SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, fusdAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), fusdSigner},
			false,
		)

		// Assert that vault balances are correct

		fusdBalance := test.ExecuteScriptAndCheck(
			t, b,
			fusd.GetBalanceScript(fungibleTokenAddress, fusdAddress),
			[][]byte{jsoncdc.MustEncode(cadence.Address(fusdAddress))},
		)

		assert.EqualValues(t, test.CadenceUFix64("1000700.0"), fusdBalance)

		userBalance := test.ExecuteScriptAndCheck(
			t, b,
			fusd.GetBalanceScript(fungibleTokenAddress, fusdAddress),
			[][]byte{jsoncdc.MustEncode(cadence.Address(userAddress))},
		)

		assert.EqualValues(t, test.CadenceUFix64("300.0"), userBalance)

		supply := test.ExecuteScriptAndCheck(
			t, b,
			fusd.GetSupplyScript(fungibleTokenAddress, fusdAddress),
			nil,
		)
		assert.EqualValues(t, test.CadenceUFix64("1001000.0"), supply)
	})
}
