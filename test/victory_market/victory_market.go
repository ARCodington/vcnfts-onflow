package victory_market

import (
	"regexp"
	"testing"
	"fmt"

	"github.com/onflow/cadence"
	emulator "github.com/onflow/flow-emulator"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	sdktemplates "github.com/onflow/flow-go-sdk/templates"
	sdktest "github.com/onflow/flow-go-sdk/test"
	"github.com/stretchr/testify/assert"

	"github.com/ARCodington/vcnfts-onflow/test/fusd"
	"github.com/ARCodington/vcnfts-onflow/test/victory_items"
	"github.com/ARCodington/vcnfts-onflow/test/test"
)

const (
	rootPath           					= ".."
	victoryMarketTransactionsRootPath 	= rootPath + "/transactions/VictoryNFTCollectionStorefront"
	victoryMarketScriptsRootPath      	= rootPath + "/scripts/VictoryNFTCollectionStorefront"

	victoryMarketContractPath           = rootPath + "/contracts/VictoryCollectibleSaleOffer.cdc"
	victoryMarketSetupAccountPath       = victoryMarketTransactionsRootPath + "/setup_account.cdc"
	victoryMarketSellPath       		= victoryMarketTransactionsRootPath + "/sell_market_bundle.cdc"
	victoryMarketRemovePath 			= victoryMarketTransactionsRootPath + "/remove_market_item.cdc"
	victoryMarketBuyPath 				= victoryMarketTransactionsRootPath + "/buy_market_item.cdc"
	victoryMarketPlaceBidPath 			= victoryMarketTransactionsRootPath + "/place_bid_market_item.cdc"
	
	victoryMarketGetBundlePricePath			= victoryMarketScriptsRootPath + "/read_bundle_price.cdc"
	victoryMarketGetVictoryIdsPath			= victoryMarketScriptsRootPath + "/read_collection_ids.cdc"
	victoryMarketGetCollectionLengthPath	= victoryMarketScriptsRootPath + "/read_collection_length.cdc"
)

func DeployContracts(t *testing.T, b *emulator.Blockchain) test.Contracts {
	accountKeys := sdktest.AccountKeyGenerator()

	fungibleTokenAddress, fusdAddress, fusdSigner := fusd.DeployContracts(t, b)
	nonFungibleTokenAddress, victoryItemsAddress, victoryItemsSigner := victory_items.DeployContracts(t, b)

	// should be able to deploy a contract as a new account with one key
	victoryMarketAccountKey, victoryMarketSigner := accountKeys.NewWithSigner()
	victoryMarketCode := loadVictoryMarket(
		fungibleTokenAddress, 
		nonFungibleTokenAddress, 
		fusdAddress, 
		victoryItemsAddress,
	)

	victoryMarketAddress, err := b.CreateAccount(
		[]*flow.AccountKey{victoryMarketAccountKey},
		[]sdktemplates.Contract{
			{
				Name:   "VictoryCollectibleSaleOffer",
				Source: string(victoryMarketCode),
			},
		},
	)
	assert.NoError(t, err)

	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// simplify the workflow by having the contract address also be our initial test collection
	victory_items.SetupAccount(t, b, victoryItemsAddress, victoryItemsSigner, nonFungibleTokenAddress, victoryItemsAddress)
	SetupAccount(t, b, victoryMarketAddress, victoryMarketSigner, victoryMarketAddress)

	return test.Contracts{
		FungibleTokenAddress:    	fungibleTokenAddress,
		FUSDAddress:           		fusdAddress,
		FUSDSigner:            		fusdSigner,
		NonFungibleTokenAddress:	nonFungibleTokenAddress,
		VictoryItemsAddress:       	victoryItemsAddress,
		VictoryItemsSigner:        	victoryItemsSigner,
		VictoryMarketAddress:    	victoryMarketAddress,
		VictoryMarketSigner:     	victoryMarketSigner,
	}
}

func SetupAccount(
	t *testing.T,
	b *emulator.Blockchain,
	userAddress flow.Address,
	userSigner crypto.Signer,
	victoryMarketAddress flow.Address,
) {
	tx := flow.NewTransaction().
		SetScript(SetupAccountScript(victoryMarketAddress.String())).
		SetGasLimit(9999).
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


// Create a new account with the Kibble and KittyItems resources set up BUT no NFTStorefront resource.
func CreatePurchaserAccount(
	t *testing.T,
	b *emulator.Blockchain,
	contracts test.Contracts,
) (flow.Address, crypto.Signer) {	
	userAddress, userSigner, _ := test.CreateAccount(t, b)
	fusd.SetupAccount(t, b, userAddress, userSigner, contracts.FungibleTokenAddress, contracts.FUSDAddress)
	victory_items.SetupAccount(t, b, userAddress, userSigner, 
		contracts.NonFungibleTokenAddress, contracts.VictoryItemsAddress)
	return userAddress, userSigner
}

func CreateAccount(
	t *testing.T,
	b *emulator.Blockchain,
	contracts test.Contracts,
) (flow.Address, crypto.Signer) {
	userAddress, userSigner := CreatePurchaserAccount(t, b, contracts)
	SetupAccount(t, b, userAddress, userSigner, contracts.VictoryMarketAddress)
	return userAddress, userSigner
}

func CreateOffer(
	t *testing.T, b *emulator.Blockchain,
	contracts test.Contracts,
	royaltyAddress flow.Address,
	userAddress flow.Address,
	userSigner crypto.Signer,
	bundleID uint64, offerType uint8, price float64, startTime uint32, endTime uint32, 
	targetPrice float64, royaltyFactor float64,
	shouldFail bool,
) {
	tx := flow.NewTransaction().
		SetScript(CreateOfferScript(contracts)).
		SetGasLimit(9999).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(userAddress)

	priceValue := test.CadenceUFix64(fmt.Sprintf("%f", price))
	targetValue := test.CadenceUFix64(fmt.Sprintf("%f", targetPrice))
	royaltyValue := test.CadenceUFix64(fmt.Sprintf("%f", royaltyFactor))

	_ = tx.AddArgument(cadence.NewAddress(royaltyAddress))
	_ = tx.AddArgument(cadence.NewUInt64(bundleID))
	_ = tx.AddArgument(cadence.NewUInt8(offerType))
	_ = tx.AddArgument(priceValue)
	_ = tx.AddArgument(cadence.NewUInt32(startTime))
	_ = tx.AddArgument(cadence.NewUInt32(endTime))
	_ = tx.AddArgument(targetValue)
	_ = tx.AddArgument(royaltyValue)

	result := test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, userAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), userSigner},
		shouldFail,
	)

	// confirm an event was raised
	if (!shouldFail) {
		eventType := fmt.Sprintf(
			"A.%s.VictoryCollectibleSaleOffer.CollectionInsertedSaleOffer",
			contracts.VictoryMarketAddress,
		)

		for _, event := range result.Events {
			if event.Type == eventType {
				return
			}
		}
		assert.Fail(t, "CollectionInsertedSaleOffer event was not emitted")
	}
}

func RemoveOffer(
	t *testing.T, b *emulator.Blockchain,
	contracts test.Contracts,
	userAddress flow.Address,
	userSigner crypto.Signer,
	bundleID uint64,
) {
	tx := flow.NewTransaction().
		SetScript(RemoveOfferScript(contracts)).
		SetGasLimit(9999).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(userAddress)

	_ = tx.AddArgument(cadence.NewUInt64(bundleID))

	result := test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, userAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), userSigner},
		false,
	)

	// confirm an event was raised
	eventType := fmt.Sprintf(
		"A.%s.VictoryCollectibleSaleOffer.CollectionRemovedSaleOffer",
		contracts.VictoryMarketAddress,
	)

	for _, event := range result.Events {
		if event.Type == eventType {
			return
		}
	}
	assert.Fail(t, "CollectionRemovedSaleOffer event was not emitted")
}

func PlaceOfferBid(
	t *testing.T, b *emulator.Blockchain,
	contracts test.Contracts,
	userAddress flow.Address,
	userSigner crypto.Signer,
	bundleID uint64,
	bidPrice float64,
	sellerAddress flow.Address,
	shouldFail bool,
) {
	tx := flow.NewTransaction().
		SetScript(PlaceBidScript(contracts)).
		SetGasLimit(9999).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(userAddress)

	priceValue := test.CadenceUFix64(fmt.Sprintf("%f", bidPrice))

	_ = tx.AddArgument(cadence.NewUInt64(bundleID))
	_ = tx.AddArgument(priceValue)
	_ = tx.AddArgument(cadence.NewAddress(userAddress))
	_ = tx.AddArgument(cadence.NewAddress(sellerAddress))

	result := test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, userAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), userSigner},
		shouldFail,
	)

	// confirm an event was raised
	if (!shouldFail) {
		eventType := fmt.Sprintf(
			"A.%s.VictoryCollectibleSaleOffer.CollectionPriceRaised",
			contracts.VictoryMarketAddress,
		)

		for _, event := range result.Events {
			if event.Type == eventType {
				return
			}
		}
		assert.Fail(t, "CollectionPriceRaised event was not emitted")
	}
}

func BuyOffer(
	t *testing.T, b *emulator.Blockchain,
	contracts test.Contracts,
	userAddress flow.Address,
	userSigner crypto.Signer,
	bundleID uint64,
	sellerAddress flow.Address,
	shouldFail bool,
) {
	tx := flow.NewTransaction().
		SetScript(BuyOfferScript(contracts)).
		SetGasLimit(9999).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(userAddress)

	_ = tx.AddArgument(cadence.NewUInt64(bundleID))
	_ = tx.AddArgument(cadence.NewAddress(sellerAddress))

	result := test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, userAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), userSigner},
		shouldFail,
	)

	// confirm an event was raised
	if (!shouldFail) {
		eventType := fmt.Sprintf(
			"A.%s.VictoryCollectible.Withdraw",
			contracts.VictoryItemsAddress,
		)

		for _, event := range result.Events {
			if event.Type == eventType {
				return
			}
		}
		assert.Fail(t, "Withdraw event was not emitted")
	}
}

func loadVictoryMarket(fungibleTokenAddress, 
	nonFungibleTokenAddress, 
	fusdAddress, 
	victoryItemsAddress flow.Address,
) []byte {
	return replaceAddressPlaceholders(
		test.ReadFile(victoryMarketContractPath),
		test.Contracts{
			FungibleTokenAddress:    	fungibleTokenAddress,
			FUSDAddress:	         	fusdAddress,
			NonFungibleTokenAddress:	nonFungibleTokenAddress,
			VictoryItemsAddress:       	victoryItemsAddress,
		},
	)
}

func replaceAddressPlaceholders(codeBytes []byte, contracts test.Contracts) []byte {
	return []byte(test.ReplaceImports(
		string(codeBytes),
		map[string]*regexp.Regexp{
			contracts.FungibleTokenAddress.String():    test.FungibleTokenAddressPlaceholder,
			contracts.FUSDAddress.String():           	test.FUSDAddressPlaceholder,
			contracts.NonFungibleTokenAddress.String():	test.NonFungibleTokenAddressPlaceholder,
			contracts.VictoryItemsAddress.String():     test.VictoryItemsAddressPlaceholder,
			contracts.VictoryMarketAddress.String():    test.VictoryMarketAddressPlaceholder,
		},
	))
}

func SetupAccountScript(victoryMarketAddress string) []byte {
	return []byte(test.ReplaceImports(
		string(test.ReadFile(victoryMarketSetupAccountPath)),
		map[string]*regexp.Regexp{
			victoryMarketAddress: test.VictoryMarketAddressPlaceholder,
		},
	))
}

func CreateOfferScript(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		test.ReadFile(victoryMarketSellPath),
		contracts,
	)
}

func RemoveOfferScript(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		test.ReadFile(victoryMarketRemovePath),
		contracts,
	)
}

func BuyOfferScript(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		test.ReadFile(victoryMarketBuyPath),
		contracts,
	)
}

func PlaceBidScript(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		test.ReadFile(victoryMarketPlaceBidPath),
		contracts,
	)
}

func GetBundlePriceScript(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		test.ReadFile(victoryMarketGetBundlePricePath),
		contracts,
	)
}

func GetVictoryIdsScript(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		test.ReadFile(victoryMarketGetVictoryIdsPath),
		contracts,
	)
}

func GetCollectionLengthScript(contracts test.Contracts) []byte {
	return replaceAddressPlaceholders(
		test.ReadFile(victoryMarketGetCollectionLengthPath),
		contracts,
	)
}
