package test

import (
	"testing"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/stretchr/testify/assert"

	"github.com/ARCodington/vcnfts-onflow/test/test"
	"github.com/ARCodington/vcnfts-onflow/test/victory_market"
	"github.com/ARCodington/vcnfts-onflow/test/victory_items"
	"github.com/ARCodington/vcnfts-onflow/test/fusd"
)

const mtTypeID = 0
const mtBrandID = 0
const mtDropID = 0
const mtContentHash = "88232f58db5d619497e852dd8ebf3ef6971ec94422d9ae673db53ed2e0f390dc"
const mtStartIssueNum = 0
const mtMaxIssueNum = 5

func TestVictoryMarketDeployContracts(t *testing.T) {
	b := test.NewBlockchain()
	victory_market.DeployContracts(t, b)
}

func TestVictoryMarketSetupAccount(t *testing.T) {
	b := test.NewBlockchain()

	contracts := victory_market.DeployContracts(t, b)

	t.Run("Should be able to create an empty Market", func(t *testing.T) {
		userAddress, userSigner, _ := test.CreateAccount(t, b)
		victory_market.SetupAccount(t, b, userAddress, userSigner, contracts.VictoryMarketAddress)
	})
}

func TestVictoryMarketPrimaryAndSecondaryOffer(t *testing.T) {
	b := test.NewBlockchain()

	contracts := victory_market.DeployContracts(t, b)
	sellerAddress, sellerSigner := victory_market.CreateAccount(t, b, contracts)

	victory_items.MintItem(t, b, 
		contracts.NonFungibleTokenAddress, 
		contracts.VictoryItemsAddress, 
		contracts.VictoryItemsSigner,
		sellerAddress, 
		mtTypeID, mtBrandID, mtDropID, mtContentHash, 
		mtStartIssueNum, mtMaxIssueNum, mtMaxIssueNum)

	bundleIDs := []uint64{0, 3}

	victory_items.CreateBundle(t, b, 
		contracts.NonFungibleTokenAddress, 
		contracts.VictoryItemsAddress,
		sellerAddress,
		sellerSigner, 
		bundleIDs,
		false,
	)

	// bundleID should be 0
	check_ids := test.ExecuteScriptAndCheckUInt64Array(
		t, b,
		victory_items.GetCollectionBundleIDsScript(contracts.NonFungibleTokenAddress.String(), contracts.VictoryItemsAddress.String()),
		[][]byte{jsoncdc.MustEncode(cadence.NewAddress(sellerAddress)), 
					jsoncdc.MustEncode(cadence.NewUInt64(0))},
	)
	// should have bundled 2 items
	assert.EqualValues(t, 2, len(check_ids))

	// create account to receive royalty payments
	royaltyAddress, _ := victory_market.CreateAccount(t, b, contracts)

	// bundleID should be 0
	t.Run("Should be able to put a bundle up for sale", func(t *testing.T) {
		victory_market.CreateOffer(
			t, b,
			contracts,
			royaltyAddress,
			sellerAddress,
			sellerSigner,
			0, 0, 100.0, 0, 0, 150.0, 0.4, false,
		)

		// assert that the price is correct
		price := test.ExecuteScriptAndCheck(
			t, b,
			victory_market.GetBundlePriceScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(sellerAddress)),
					 jsoncdc.MustEncode(cadence.NewUInt64(0))},
		)
		assert.EqualValues(t, test.CadenceUFix64("100.0"), price)
	})

	// bundleID should be 0 and will be removed from collection as well
	t.Run("Should be able to remove a bundle up for sale", func(t *testing.T) {
		victory_market.RemoveOffer(
			t, b,
			contracts,
			sellerAddress,
			sellerSigner,
			0,
		)
	})

	// new bundle with ID 1
	victory_items.CreateBundle(t, b, 
		contracts.NonFungibleTokenAddress, 
		contracts.VictoryItemsAddress,
		sellerAddress,
		sellerSigner, 
		bundleIDs,
		false,
	)

	t.Run("Should not be able to set royalty above 100%", func(t *testing.T) {
		victory_market.CreateOffer(
			t, b,
			contracts,
			royaltyAddress,
			sellerAddress,
			sellerSigner,
			1, 0, 100.0, 0, 0, 150.0, 1.5, 
			// should fail
			true,
		)	
	})

	// assume this will work again
	victory_market.CreateOffer(
		t, b,
		contracts,
		royaltyAddress,
		sellerAddress,
		sellerSigner,
		1, 0, 100.0, 0, 0, 150.0, 0.4, false,
	)

	// create a buyer/re-seller account
	buyerAddress, buyerSigner := victory_market.CreateAccount(t, b, contracts)

	t.Run("Should not be able to transfer an item that is for sale", func(t *testing.T) {
		victory_items.TransferItem(
			t, b,
			contracts.NonFungibleTokenAddress, 
			contracts.VictoryItemsAddress,
			sellerAddress, 
			sellerSigner,
			0, buyerAddress, true,
		)
	})

	t.Run("Should not be able to purchase a bundle with insufficient funds", func(t *testing.T) {
		// assert that the account collection is empty
		length := test.ExecuteScriptAndCheck(
			t, b,
			victory_items.GetCollectionLengthScript(contracts.NonFungibleTokenAddress.String(), contracts.VictoryItemsAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(buyerAddress))},
		)
		assert.EqualValues(t, cadence.NewInt(0), length)

		// fund the buyer account
		fusd.Mint(
			t, b,
			contracts.FungibleTokenAddress,
			contracts.FUSDAddress,
			contracts.FUSDSigner,
			buyerAddress,
			"90.0",
			false,
		)
		// verify the balance
		userBalance := test.ExecuteScriptAndCheck(
			t, b,
			fusd.GetBalanceScript(contracts.FungibleTokenAddress, contracts.FUSDAddress),
			[][]byte{jsoncdc.MustEncode(cadence.Address(buyerAddress))},
		)
		assert.EqualValues(t, test.CadenceUFix64("90.0"), userBalance)

		// insufficient funds
		victory_market.BuyOffer(
			t, b,
			contracts,
			buyerAddress,
			buyerSigner,
			1,
			sellerAddress,
			true,
		)
	})

	t.Run("Should be able to purchase a bundle for sale", func(t *testing.T) {
		// assert that the account collection is empty
		length := test.ExecuteScriptAndCheck(
			t, b,
			victory_items.GetCollectionLengthScript(contracts.NonFungibleTokenAddress.String(), contracts.VictoryItemsAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(buyerAddress))},
		)
		assert.EqualValues(t, cadence.NewInt(0), length)

		// fund the buyer account to the amount of the offer
		fusd.Mint(
			t, b,
			contracts.FungibleTokenAddress,
			contracts.FUSDAddress,
			contracts.FUSDSigner,
			buyerAddress,
			"10.0",
			false,
		)
		// now we should be able to buy
		victory_market.BuyOffer(
			t, b,
			contracts,
			buyerAddress,
			buyerSigner,
			1,
			sellerAddress,
			false,
		)

		// Item ID 0 should no longer be for sale
		for_sale := test.ExecuteScriptAndCheck(
			t, b,
			victory_items.GetCollectionForSaleScript(contracts.NonFungibleTokenAddress.String(), contracts.VictoryItemsAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(sellerAddress)),
					 jsoncdc.MustEncode(cadence.NewUInt64(0))},
		)
		assert.EqualValues(t, cadence.NewBool(false), for_sale)

		// assert that the buyer has 2 NFTs now
		length = test.ExecuteScriptAndCheck(
			t, b,
			victory_items.GetCollectionLengthScript(contracts.NonFungibleTokenAddress.String(), contracts.VictoryItemsAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(buyerAddress))},
		)
		assert.EqualValues(t, cadence.NewInt(2), length)

		// verify the balances
		// buyer should have 0 FUSD left
		userBalance := test.ExecuteScriptAndCheck(
			t, b,
			fusd.GetBalanceScript(contracts.FungibleTokenAddress, contracts.FUSDAddress),
			[][]byte{jsoncdc.MustEncode(cadence.Address(buyerAddress))},
		)
		assert.EqualValues(t, test.CadenceUFix64("0.0"), userBalance)

		// seller should have 60 FUSD (100 FUSD - 40% royalty)
		userBalance = test.ExecuteScriptAndCheck(
			t, b,
			fusd.GetBalanceScript(contracts.FungibleTokenAddress, contracts.FUSDAddress),
			[][]byte{jsoncdc.MustEncode(cadence.Address(sellerAddress))},
		)
		assert.EqualValues(t, test.CadenceUFix64("60.0"), userBalance)

		// should have received 40 FUSD in royalties
		userBalance = test.ExecuteScriptAndCheck(
			t, b,
			fusd.GetBalanceScript(contracts.FungibleTokenAddress, contracts.FUSDAddress),
			[][]byte{jsoncdc.MustEncode(cadence.Address(royaltyAddress))},
		)
		assert.EqualValues(t, test.CadenceUFix64("40.0"), userBalance)
	})

	t.Run("Should be able to re-sell an item on secondary market", func(t *testing.T) {
		resellIDs := []uint64{0}

		victory_items.CreateBundle(t, b, 
			contracts.NonFungibleTokenAddress, 
			contracts.VictoryItemsAddress,
			buyerAddress,
			buyerSigner, 
			resellIDs,
			false,
		)

		// bundleID should be 0 since this is the first listing for buyer
		check_ids := test.ExecuteScriptAndCheckUInt64Array(
			t, b,
			victory_items.GetCollectionBundleIDsScript(contracts.NonFungibleTokenAddress.String(), contracts.VictoryItemsAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(buyerAddress)), 
						jsoncdc.MustEncode(cadence.NewUInt64(0))},
		)
		// should have bundled 1 item
		assert.EqualValues(t, 1, len(check_ids))

		// list it - royalty is 5% on re-sale for both the original owner and marketplace
		victory_market.CreateOffer(
			t, b,
			contracts,
			royaltyAddress,
			buyerAddress,
			buyerSigner, 
			0, 0, 50.0, 0, 0, 60.0, 0.05, false,
		)

		// buy it back for a loss - not recommended ;)
		// also note this arbitrary scenario means the original seller will receive
		// a royalty for re-purchasing what they previously owned - weird, but correct
		victory_market.BuyOffer(
			t, b,
			contracts,
			sellerAddress,
			sellerSigner,
			0,
			buyerAddress,
			false,
		)

		// verify the balances
		// buyer/re-seller should have 45 FUSD left (50 FUSD - 10% royalty)
		userBalance := test.ExecuteScriptAndCheck(
			t, b,
			fusd.GetBalanceScript(contracts.FungibleTokenAddress, contracts.FUSDAddress),
			[][]byte{jsoncdc.MustEncode(cadence.Address(buyerAddress))},
		)
		assert.EqualValues(t, test.CadenceUFix64("45.0"), userBalance)

		// original seller should have 12.5 FUSD (60 FUSD - 50 USD sale price + 2.5 FUSD royalty)
		userBalance = test.ExecuteScriptAndCheck(
			t, b,
			fusd.GetBalanceScript(contracts.FungibleTokenAddress, contracts.FUSDAddress),
			[][]byte{jsoncdc.MustEncode(cadence.Address(sellerAddress))},
		)
		assert.EqualValues(t, test.CadenceUFix64("12.5"), userBalance)

		// should have received 42.5 FUSD in royalties (previous 40 FUSD + 2.5 FUSD royalty from this sale)
		userBalance = test.ExecuteScriptAndCheck(
			t, b,
			fusd.GetBalanceScript(contracts.FungibleTokenAddress, contracts.FUSDAddress),
			[][]byte{jsoncdc.MustEncode(cadence.Address(royaltyAddress))},
		)
		assert.EqualValues(t, test.CadenceUFix64("42.5"), userBalance)
	})
}

func TestVictoryMarketAuction(t *testing.T) {
	b := test.NewBlockchain()

	contracts := victory_market.DeployContracts(t, b)
	sellerAddress, sellerSigner := victory_market.CreateAccount(t, b, contracts)

	victory_items.MintItem(t, b, 
		contracts.NonFungibleTokenAddress, 
		contracts.VictoryItemsAddress, 
		contracts.VictoryItemsSigner,
		sellerAddress, 
		mtTypeID, mtBrandID, mtDropID, mtContentHash, 
		mtStartIssueNum, mtMaxIssueNum, mtMaxIssueNum)

	bundleIDs := []uint64{0, 3}

	victory_items.CreateBundle(t, b, 
		contracts.NonFungibleTokenAddress, 
		contracts.VictoryItemsAddress,
		sellerAddress,
		sellerSigner, 
		bundleIDs,
		false,
	)

	// create account to receive royalty payments
	royaltyAddress, _ := victory_market.CreateAccount(t, b, contracts)

	// bundleID should be 0
	victory_market.CreateOffer(
		t, b,
		contracts,
		royaltyAddress,
		sellerAddress,
		sellerSigner,
		0, 0, 100.0, 0, 0, 150.0, 0.4, false,
	)

	// create bidder accounts
	buyerAddress, buyerSigner := victory_market.CreateAccount(t, b, contracts)
	buyer2Address, buyer2Signer := victory_market.CreateAccount(t, b, contracts)

	t.Run("Should be able to auction a bundle for sale", func(t *testing.T) {
		// should not be able to bid without funds
		victory_market.PlaceOfferBid(
			t, b,
			contracts,
			buyerAddress,
			buyerSigner,
			0,
			100.0,
			sellerAddress,
			true,
		)

		// fund the buyer accounts
		fusd.Mint(
			t, b,
			contracts.FungibleTokenAddress,
			contracts.FUSDAddress,
			contracts.FUSDSigner,
			buyerAddress,
			"150.0",
			false,
		)
		fusd.Mint(
			t, b,
			contracts.FungibleTokenAddress,
			contracts.FUSDAddress,
			contracts.FUSDSigner,
			buyer2Address,
			"150.0",
			false,
		)

		// should not be able to lower the price
		victory_market.PlaceOfferBid(
			t, b,
			contracts,
			buyerAddress,
			buyerSigner,
			0,
			90.0,
			sellerAddress,
			true,
		)

		// should be able to start the bidding at the initial price
		victory_market.PlaceOfferBid(
			t, b,
			contracts,
			buyerAddress,
			buyerSigner,
			0,
			100.0,
			sellerAddress,
			false,
		)

		// buyer should still have 150 FUSD left
		userBalance := test.ExecuteScriptAndCheck(
			t, b,
			fusd.GetBalanceScript(contracts.FungibleTokenAddress, contracts.FUSDAddress),
			[][]byte{jsoncdc.MustEncode(cadence.Address(buyerAddress))},
		)
		assert.EqualValues(t, test.CadenceUFix64("150.0"), userBalance)

		// should not be able to rebid at the same price after first bid
		victory_market.PlaceOfferBid(
			t, b,
			contracts,
			buyerAddress,
			buyerSigner,
			0,
			100.0,
			sellerAddress,
			true,
		)

		// place another bid
		victory_market.PlaceOfferBid(
			t, b,
			contracts,
			buyerAddress,
			buyerSigner,
			0,
			120.0,
			sellerAddress,
			false,
		)
		// assert that the price is correct
		price := test.ExecuteScriptAndCheck(
			t, b,
			victory_market.GetBundlePriceScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(sellerAddress)),
					 jsoncdc.MustEncode(cadence.NewUInt64(0))},
		)
		assert.EqualValues(t, test.CadenceUFix64("120.0"), price)

		// place another bid
		victory_market.PlaceOfferBid(
			t, b,
			contracts,
			buyer2Address,
			buyer2Signer,
			0,
			135.0,
			sellerAddress,
			false,
		)
		// assert that the price is correct
		price = test.ExecuteScriptAndCheck(
			t, b,
			victory_market.GetBundlePriceScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(sellerAddress)),
					 jsoncdc.MustEncode(cadence.NewUInt64(0))},
		)
		assert.EqualValues(t, test.CadenceUFix64("135.0"), price)

		// place the final bid
		victory_market.PlaceOfferBid(
			t, b,
			contracts,
			buyerAddress,
			buyerSigner,
			0,
			150.0,
			sellerAddress,
			false,
		)
		// assert that the price is correct
		price = test.ExecuteScriptAndCheck(
			t, b,
			victory_market.GetBundlePriceScript(contracts),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(sellerAddress)),
						jsoncdc.MustEncode(cadence.NewUInt64(0))},
		)
		assert.EqualValues(t, test.CadenceUFix64("150.0"), price)
		
		victory_market.BuyOffer(
			t, b,
			contracts,
			buyerAddress,
			buyerSigner,
			0,
			sellerAddress,
			false,
		)

		// verify the balances
		// buyer should have 0 FUSD left
		userBalance = test.ExecuteScriptAndCheck(
			t, b,
			fusd.GetBalanceScript(contracts.FungibleTokenAddress, contracts.FUSDAddress),
			[][]byte{jsoncdc.MustEncode(cadence.Address(buyerAddress))},
		)
		assert.EqualValues(t, test.CadenceUFix64("0.0"), userBalance)

		// seller should have 90 FUSD (150 FUSD - 40% royalty)
		userBalance = test.ExecuteScriptAndCheck(
			t, b,
			fusd.GetBalanceScript(contracts.FungibleTokenAddress, contracts.FUSDAddress),
			[][]byte{jsoncdc.MustEncode(cadence.Address(sellerAddress))},
		)
		assert.EqualValues(t, test.CadenceUFix64("90.0"), userBalance)

		// marketplace owner should have received 60 FUSD in royalties
		userBalance = test.ExecuteScriptAndCheck(
			t, b,
			fusd.GetBalanceScript(contracts.FungibleTokenAddress, contracts.FUSDAddress),
			[][]byte{jsoncdc.MustEncode(cadence.Address(royaltyAddress))},
		)
		assert.EqualValues(t, test.CadenceUFix64("60.0"), userBalance)
	})
}

func TestVictoryMarketMintOnDemand(t *testing.T) {
	b := test.NewBlockchain()

	contracts := victory_market.DeployContracts(t, b)
	sellerAddress, sellerSigner := victory_market.CreateAccount(t, b, contracts)

	victory_items.MintItem(t, b, 
		contracts.NonFungibleTokenAddress, 
		contracts.VictoryItemsAddress, 
		contracts.VictoryItemsSigner,
		sellerAddress, 
		mtTypeID, mtBrandID, mtDropID, mtContentHash, 
		mtStartIssueNum, 1, mtMaxIssueNum)

	// mint one more on demand based on the first item as reference
	victory_items.MintItemOnDemand(t, b, 
		contracts.NonFungibleTokenAddress, 
		contracts.VictoryItemsAddress, 
		contracts.VictoryItemsSigner,
		sellerAddress, 0, 1)

	bundleIDs := []uint64{1}

	victory_items.CreateBundle(t, b, 
		contracts.NonFungibleTokenAddress, 
		contracts.VictoryItemsAddress,
		sellerAddress,
		sellerSigner, 
		bundleIDs,
		false,
	)

	// create account to receive royalty payments
	royaltyAddress, _ := victory_market.CreateAccount(t, b, contracts)

	// bundleID should be 0
	victory_market.CreateOffer(
		t, b,
		contracts,
		royaltyAddress,
		sellerAddress,
		sellerSigner,
		0, 0, 100.0, 0, 0, 150.0, 0.4, false,
	)

	// create a buyer/re-seller account
	buyerAddress, buyerSigner := victory_market.CreateAccount(t, b, contracts)

	t.Run("Should be able to purchase a mint on demand bundle for sale", func(t *testing.T) {
		// assert that the account collection is empty
		length := test.ExecuteScriptAndCheck(
			t, b,
			victory_items.GetCollectionLengthScript(contracts.NonFungibleTokenAddress.String(), contracts.VictoryItemsAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(buyerAddress))},
		)
		assert.EqualValues(t, cadence.NewInt(0), length)

		// fund the buyer account to the amount of the offer
		fusd.Mint(
			t, b,
			contracts.FungibleTokenAddress,
			contracts.FUSDAddress,
			contracts.FUSDSigner,
			buyerAddress,
			"100.0",
			false,
		)
		// now we should be able to buy
		victory_market.BuyOffer(
			t, b,
			contracts,
			buyerAddress,
			buyerSigner,
			0,
			sellerAddress,
			false,
		)

		// Item ID 0 should no longer be for sale
		for_sale := test.ExecuteScriptAndCheck(
			t, b,
			victory_items.GetCollectionForSaleScript(contracts.NonFungibleTokenAddress.String(), contracts.VictoryItemsAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(sellerAddress)),
					 jsoncdc.MustEncode(cadence.NewUInt64(1))},
		)
		assert.EqualValues(t, cadence.NewBool(false), for_sale)

		// assert that the buyer has 1 NFT now
		length = test.ExecuteScriptAndCheck(
			t, b,
			victory_items.GetCollectionLengthScript(contracts.NonFungibleTokenAddress.String(), contracts.VictoryItemsAddress.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(buyerAddress))},
		)
		assert.EqualValues(t, cadence.NewInt(1), length)

		// verify the balances
		// buyer should have 0 FUSD left
		userBalance := test.ExecuteScriptAndCheck(
			t, b,
			fusd.GetBalanceScript(contracts.FungibleTokenAddress, contracts.FUSDAddress),
			[][]byte{jsoncdc.MustEncode(cadence.Address(buyerAddress))},
		)
		assert.EqualValues(t, test.CadenceUFix64("0.0"), userBalance)

		// seller should have 60 FUSD (100 FUSD - 40% royalty)
		userBalance = test.ExecuteScriptAndCheck(
			t, b,
			fusd.GetBalanceScript(contracts.FungibleTokenAddress, contracts.FUSDAddress),
			[][]byte{jsoncdc.MustEncode(cadence.Address(sellerAddress))},
		)
		assert.EqualValues(t, test.CadenceUFix64("60.0"), userBalance)

		// should have received 40 FUSD in royalties
		userBalance = test.ExecuteScriptAndCheck(
			t, b,
			fusd.GetBalanceScript(contracts.FungibleTokenAddress, contracts.FUSDAddress),
			[][]byte{jsoncdc.MustEncode(cadence.Address(royaltyAddress))},
		)
		assert.EqualValues(t, test.CadenceUFix64("40.0"), userBalance)
	})
}