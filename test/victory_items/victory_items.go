package victory_items

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/onflow/cadence"
	emulator "github.com/onflow/flow-emulator"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	sdktemplates "github.com/onflow/flow-go-sdk/templates"
	sdktest "github.com/onflow/flow-go-sdk/test"
	nftcontracts "github.com/onflow/flow-nft/lib/go/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ARCodington/vcnfts-onflow/test/test"
)

const (
	rootPath           					= ".."
	victoryItemsTransactionsRootPath 	= rootPath + "/transactions/VictoryNFTCollectionItem"
	victoryItemsScriptsRootPath      	= rootPath + "/scripts/VictoryNFTCollectionItem"

	victoryItemsContractPath            	= rootPath + "/contracts/VictoryNFTCollectionItem.cdc"

	victoryItemsSetupAccountPath        	= victoryItemsTransactionsRootPath + "/setup_account.cdc"
	victoryItemsMintVictoryItemPath       	= victoryItemsTransactionsRootPath + "/mint_collectible.cdc"
	victoryItemsTransferVictoryItemPath 	= victoryItemsTransactionsRootPath + "/transfer_collectible.cdc"
	victoryItemsUpdateMetaURLPath			= victoryItemsTransactionsRootPath + "/update_collectible_metaurl.cdc"
	victoryItemsUpdateGeoURLPath 			= victoryItemsTransactionsRootPath + "/update_collectible_geourl.cdc"
	victoryItemsUpdateHashVictoryItemPath 	= victoryItemsTransactionsRootPath + "/update_collectible_hash.cdc"
	victoryItemsCreateBundlePath	 		= victoryItemsTransactionsRootPath + "/create_collection_bundle.cdc"
	victoryItemsRemoveBundlePath	 		= victoryItemsTransactionsRootPath + "/remove_collection_bundle.cdc"
	victoryItemsRemoveAllBundlesPath 		= victoryItemsTransactionsRootPath + "/remove_all_bundles.cdc"

	victoryItemsGetVictoryItemSupplyPath	= victoryItemsScriptsRootPath + "/read_collectibles_supply.cdc"
	victoryItemsGetCollectibleMetaURLPath	= victoryItemsScriptsRootPath + "/read_collectible_metaurl.cdc"
	victoryItemsGetCollectibleGeoURLPath	= victoryItemsScriptsRootPath + "/read_collectible_geourl.cdc"
	victoryItemsGetCollectibleHashPath		= victoryItemsScriptsRootPath + "/read_collectible_hash.cdc"
	victoryItemsGetCollectionLengthPath		= victoryItemsScriptsRootPath + "/read_collection_length.cdc"
	victoryItemsGetCollectionIDsPath		= victoryItemsScriptsRootPath + "/read_collection_ids.cdc"
	victoryItemsGetCollectionBundleIDsPath	= victoryItemsScriptsRootPath + "/read_collection_bundle_ids.cdc"
	victoryItemsGetCollectionForSalePath	= victoryItemsScriptsRootPath + "/read_collection_item_for_sale.cdc"
	victoryItemsGetNextBundlePath			= victoryItemsScriptsRootPath + "/read_collection_next_bundle_id.cdc"
)

func DeployContracts(
	t *testing.T,
	b *emulator.Blockchain,
) (flow.Address, flow.Address, crypto.Signer) {
	accountKeys := sdktest.AccountKeyGenerator()

	// should be able to deploy a contract as a new account with no keys
	nftCode := nftcontracts.NonFungibleToken()
	nftAddress, err := b.CreateAccount(
		nil,
		[]sdktemplates.Contract{
			{
				Name:   "NonFungibleToken",
				Source: string(nftCode),
			},
		},
	)
	require.NoError(t, err)

	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// should be able to deploy a contract as a new account with one key
	victoryItemsAccountKey, victoryItemsSigner := accountKeys.NewWithSigner()
	victoryItemsCode := loadVictoryItems(nftAddress.String())
	victoryItemsAddr, err := b.CreateAccount(
		[]*flow.AccountKey{victoryItemsAccountKey},
		[]sdktemplates.Contract{
			{
				Name:   "VictoryNFTCollectionItem",
				Source: string(victoryItemsCode),
			},
		},
	)
	assert.NoError(t, err)

	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// simplify the workflow by having the contract address also be our initial test collection
	SetupAccount(t, b, victoryItemsAddr, victoryItemsSigner, nftAddress, victoryItemsAddr)

	return nftAddress, victoryItemsAddr, victoryItemsSigner
}

func SetupAccount(
	t *testing.T,
	b *emulator.Blockchain,
	userAddress flow.Address,
	userSigner crypto.Signer,
	nftAddress flow.Address,
	victoryItemsAddress flow.Address,
) {
	tx := flow.NewTransaction().
		SetScript(SetupAccountScript(nftAddress.String(), victoryItemsAddress.String())).
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

func MintItem(
	t *testing.T, b *emulator.Blockchain,
	nftAddress, 
	victoryItemsAddr flow.Address,
	victoryItemsSigner crypto.Signer, 
	recipientAddr flow.Address,
	typeID uint64, brandID uint64, dropID uint64, contentHash string, maxIssueNum uint32, metaURL string, geoURL string,
) {
	tx := flow.NewTransaction().
		SetScript(MintVictoryItemScript(nftAddress.String(), victoryItemsAddr.String())).
		SetGasLimit(9999).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(victoryItemsAddr)

	_ = tx.AddArgument(cadence.NewAddress(recipientAddr))
	_ = tx.AddArgument(cadence.NewUInt64(typeID))
	_ = tx.AddArgument(cadence.NewUInt64(brandID))
	_ = tx.AddArgument(cadence.NewUInt64(dropID))
	_ = tx.AddArgument(cadence.NewString(contentHash))
	_ = tx.AddArgument(cadence.NewUInt32(maxIssueNum))
	_ = tx.AddArgument(cadence.NewString(metaURL))
	_ = tx.AddArgument(cadence.NewString(geoURL))

	result := test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, victoryItemsAddr},
		[]crypto.Signer{b.ServiceKey().Signer(), victoryItemsSigner},
		false,
	)

	// confirm an event was raised
	eventType := fmt.Sprintf(
		"A.%s.VictoryNFTCollectionItem.Minted",
		victoryItemsAddr,
	)

	for _, event := range result.Events {
		if event.Type == eventType {
			return
		}
	}
	assert.Fail(t, "Minted event was not emitted")
}

func TransferItem(
	t *testing.T, b *emulator.Blockchain,
	nftAddress, 
	victoryItemsAddr flow.Address, 
	sellerAddr flow.Address, 
	sellerSigner crypto.Signer,
	withdrawID uint64, recipientAddr flow.Address, shouldFail bool,
) {

	tx := flow.NewTransaction().
		SetScript(TransferVictoryItemScript(nftAddress.String(), victoryItemsAddr.String())).
		SetGasLimit(9999).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(sellerAddr)

	_ = tx.AddArgument(cadence.NewAddress(recipientAddr))
	_ = tx.AddArgument(cadence.NewUInt64(withdrawID))

	result := test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, sellerAddr},
		[]crypto.Signer{b.ServiceKey().Signer(), sellerSigner},
		shouldFail,
	)

	if (!shouldFail) {
		// confirm an event was raised
		eventType := fmt.Sprintf(
			"A.%s.VictoryNFTCollectionItem.Deposit",
			victoryItemsAddr,
		)

		for _, event := range result.Events {
			if event.Type == eventType {
				return
			}
		}
		assert.Fail(t, "Deposit event was not emitted")
	}
}

func UpdateItemMeta(
	t *testing.T, b *emulator.Blockchain,
	nftAddress, victoryItemsAddr flow.Address, victoryItemsSigner crypto.Signer,
	updateID uint64, updateValue string,
) {

	tx := flow.NewTransaction().
		SetScript(UpdateMetaURLVictoryItemScript(nftAddress.String(), victoryItemsAddr.String())).
		SetGasLimit(9999).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(victoryItemsAddr)

	_ = tx.AddArgument(cadence.NewUInt64(updateID))
	_ = tx.AddArgument(cadence.NewString(updateValue))

	result := test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, victoryItemsAddr},
		[]crypto.Signer{b.ServiceKey().Signer(), victoryItemsSigner},
		false,
	)

	// confirm an event was raised
	eventType := fmt.Sprintf(
		"A.%s.VictoryNFTCollectionItem.MetaUpdated",
		victoryItemsAddr,
	)

	for _, event := range result.Events {
		if event.Type == eventType {
			return
		}
	}
	assert.Fail(t, "MetaUpdated event was not emitted")
}

func UpdateItemGeo(
	t *testing.T, b *emulator.Blockchain,
	nftAddress, victoryItemsAddr flow.Address, victoryItemsSigner crypto.Signer,
	updateID uint64, updateValue string,
) {

	tx := flow.NewTransaction().
		SetScript(UpdateGeoURLVictoryItemScript(nftAddress.String(), victoryItemsAddr.String())).
		SetGasLimit(9999).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(victoryItemsAddr)

	_ = tx.AddArgument(cadence.NewUInt64(updateID))
	_ = tx.AddArgument(cadence.NewString(updateValue))

	result := test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, victoryItemsAddr},
		[]crypto.Signer{b.ServiceKey().Signer(), victoryItemsSigner},
		false,
	)

	// confirm an event was raised
	eventType := fmt.Sprintf(
		"A.%s.VictoryNFTCollectionItem.Geolocated",
		victoryItemsAddr,
	)

	for _, event := range result.Events {
		if event.Type == eventType {
			return
		}
	}
	assert.Fail(t, "Geolocated event was not emitted")
}

func UpdateItemHash(
	t *testing.T, b *emulator.Blockchain,
	nftAddress, victoryItemsAddr flow.Address, victoryItemsSigner crypto.Signer,
	updateID uint64, updateValue string,
) {

	tx := flow.NewTransaction().
		SetScript(UpdateHashVictoryItemScript(nftAddress.String(), victoryItemsAddr.String())).
		SetGasLimit(9999).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(victoryItemsAddr)

	_ = tx.AddArgument(cadence.NewUInt64(updateID))
	_ = tx.AddArgument(cadence.NewString(updateValue))

	result := test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, victoryItemsAddr},
		[]crypto.Signer{b.ServiceKey().Signer(), victoryItemsSigner},
		false,
	)

	// confirm an event was raised
	eventType := fmt.Sprintf(
		"A.%s.VictoryNFTCollectionItem.HashUpdated",
		victoryItemsAddr,
	)

	for _, event := range result.Events {
		if event.Type == eventType {
			return
		}
	}
	assert.Fail(t, "HashUpdated event was not emitted")
}

func CreateBundle(
	t *testing.T, b *emulator.Blockchain,
	nftAddress, 
	victoryItemsAddr flow.Address, 
	sellerAddr flow.Address, 
	sellerSigner crypto.Signer,
	itemIDs []uint64,
	shouldFail bool,
) {

	values := make([]cadence.Value, len(itemIDs))
	for i, v := range itemIDs {
		values[i] = cadence.NewUInt64(v)
	}

	tx := flow.NewTransaction().
		SetScript(CreateBundleScript(nftAddress.String(), victoryItemsAddr.String())).
		SetGasLimit(9999).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(sellerAddr)

	_ = tx.AddArgument(cadence.NewArray(values))

	result := test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, sellerAddr},
		[]crypto.Signer{b.ServiceKey().Signer(), sellerSigner},
		shouldFail,
	)

	if (!shouldFail) {
		// confirm an event was raised
		eventType := fmt.Sprintf(
			"A.%s.VictoryNFTCollectionItem.BundleCreated",
			victoryItemsAddr,
		)

		for _, event := range result.Events {
			if event.Type == eventType {
				return
			}
		}
		assert.Fail(t, "BundleCreated event was not emitted")
	}
}

func RemoveBundle(
	t *testing.T, b *emulator.Blockchain,
	nftAddress, victoryItemsAddr flow.Address, victoryItemsSigner crypto.Signer,
	bundleID uint64,
	shouldFail bool,
) {

	tx := flow.NewTransaction().
		SetScript(RemoveBundleScript(nftAddress.String(), victoryItemsAddr.String())).
		SetGasLimit(9999).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(victoryItemsAddr)

	_ = tx.AddArgument(cadence.NewUInt64(bundleID))

	result := test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, victoryItemsAddr},
		[]crypto.Signer{b.ServiceKey().Signer(), victoryItemsSigner},
		shouldFail,
	)

	if (!shouldFail) {
		// confirm an event was raised
		eventType := fmt.Sprintf(
			"A.%s.VictoryNFTCollectionItem.BundleRemoved",
			victoryItemsAddr,
		)

		for _, event := range result.Events {
			if event.Type == eventType {
				return
			}
		}
		assert.Fail(t, "BundleRemoved event was not emitted")
	}
}

func RemoveAllBundles(
	t *testing.T, b *emulator.Blockchain,
	nftAddress, victoryItemsAddr flow.Address, victoryItemsSigner crypto.Signer,
	shouldFail bool,
) {

	tx := flow.NewTransaction().
		SetScript(RemoveAllBundlesScript(nftAddress.String(), victoryItemsAddr.String())).
		SetGasLimit(9999).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(victoryItemsAddr)

	result := test.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, victoryItemsAddr},
		[]crypto.Signer{b.ServiceKey().Signer(), victoryItemsSigner},
		shouldFail,
	)

	if (!shouldFail) {
		// confirm an event was raised
		eventType := fmt.Sprintf(
			"A.%s.VictoryNFTCollectionItem.AllBundlesRemoved",
			victoryItemsAddr,
		)

		for _, event := range result.Events {
			if event.Type == eventType {
				return
			}
		}
		assert.Fail(t, "AllBundlesRemoved event was not emitted")
	}
}

func loadVictoryItems(nftAddress string) []byte {
	return []byte(test.ReplaceImports(
		string(test.ReadFile(victoryItemsContractPath)),
		map[string]*regexp.Regexp{
			nftAddress: test.NonFungibleTokenAddressPlaceholder,
		},
	))
}

func replaceAddressPlaceholders(code, nftAddress, victoryItemsAddress string) []byte {
	return []byte(test.ReplaceImports(
		code,
		map[string]*regexp.Regexp{
			nftAddress:        		test.NonFungibleTokenAddressPlaceholder,
			victoryItemsAddress:	test.VictoryItemsAddressPlaceholder,
		},
	))
}

func SetupAccountScript(nftAddress, victoryItemsAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(victoryItemsSetupAccountPath)),
		nftAddress,
		victoryItemsAddress,
	)
}

func MintVictoryItemScript(nftAddress, victoryItemsAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(victoryItemsMintVictoryItemPath)),
		nftAddress,
		victoryItemsAddress,
	)
}

func TransferVictoryItemScript(nftAddress, victoryItemsAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(victoryItemsTransferVictoryItemPath)),
		nftAddress,
		victoryItemsAddress,
	)
}

func UpdateMetaURLVictoryItemScript(nftAddress, victoryItemsAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(victoryItemsUpdateMetaURLPath)),
		nftAddress,
		victoryItemsAddress,
	)
}

func UpdateGeoURLVictoryItemScript(nftAddress, victoryItemsAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(victoryItemsUpdateGeoURLPath)),
		nftAddress,
		victoryItemsAddress,
	)
}

func UpdateHashVictoryItemScript(nftAddress, victoryItemsAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(victoryItemsUpdateHashVictoryItemPath)),
		nftAddress,
		victoryItemsAddress,
	)
}

func CreateBundleScript(nftAddress, victoryItemsAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(victoryItemsCreateBundlePath)),
		nftAddress,
		victoryItemsAddress,
	)
}

func RemoveBundleScript(nftAddress, victoryItemsAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(victoryItemsRemoveBundlePath)),
		nftAddress,
		victoryItemsAddress,
	)
}

func RemoveAllBundlesScript(nftAddress, victoryItemsAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(victoryItemsRemoveAllBundlesPath)),
		nftAddress,
		victoryItemsAddress,
	)
}

func GetVictoryItemSupplyScript(nftAddress, victoryItemsAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(victoryItemsGetVictoryItemSupplyPath)),
		nftAddress,
		victoryItemsAddress,
	)
}

func GetCollectionLengthScript(nftAddress, victoryItemsAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(victoryItemsGetCollectionLengthPath)),
		nftAddress,
		victoryItemsAddress,
	)
}

func GetCollectibleMetaURLScript(nftAddress, victoryItemsAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(victoryItemsGetCollectibleMetaURLPath)),
		nftAddress,
		victoryItemsAddress,
	)
}

func GetCollectibleGeoURLScript(nftAddress, victoryItemsAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(victoryItemsGetCollectibleGeoURLPath)),
		nftAddress,
		victoryItemsAddress,
	)
}

func GetCollectibleHashScript(nftAddress, victoryItemsAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(victoryItemsGetCollectibleHashPath)),
		nftAddress,
		victoryItemsAddress,
	)
}

func GetCollectionIDsScript(nftAddress, victoryItemsAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(victoryItemsGetCollectionIDsPath)),
		nftAddress,
		victoryItemsAddress,
	)
}

func GetCollectionBundleIDsScript(nftAddress, victoryItemsAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(victoryItemsGetCollectionBundleIDsPath)),
		nftAddress,
		victoryItemsAddress,
	)
}

func GetNextBundleIDScript(nftAddress, victoryItemsAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(victoryItemsGetNextBundlePath)),
		nftAddress,
		victoryItemsAddress,
	)
}

func GetCollectionForSaleScript(nftAddress, victoryItemsAddress string) []byte {
	return replaceAddressPlaceholders(
		string(test.ReadFile(victoryItemsGetCollectionForSalePath)),
		nftAddress,
		victoryItemsAddress,
	)
}
