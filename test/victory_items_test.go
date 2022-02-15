package test

import (
	"testing"
	"encoding/hex"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/stretchr/testify/assert"

	"github.com/ARCodington/vcnfts-onflow/test/victory_items"
	"github.com/ARCodington/vcnfts-onflow/test/test"
)

const typeID = 0
const brandID = 0
const dropID = 0
const contentHash = "88232f58db5d619497e852dd8ebf3ef69712394422d9ae673db53ed2e0f390dc"
const startIssueNum = 0
const maxIssueNum = 5
const metaURL = "https://offchain.storage.com/"
const geoURL = "https://geolocation.com:443/"
const singleItemMintCount = 1
const newHash = "52e989867cf370697ac633856582c04b4568980fa88840dd5b8d5afb016b7a60"

func TestVictoryItemsDeployContracts(t *testing.T) {
	b := test.NewBlockchain()
	victory_items.DeployContracts(t, b)
}

func TestCreateVictoryItem(t *testing.T) {
	b := test.NewBlockchain()

	nftAddress, victoryItemsAddr, victoryItemsSigner := victory_items.DeployContracts(t, b)

	supply := test.ExecuteScriptAndCheck(
		t, b,
		victory_items.GetVictoryItemSupplyScript(nftAddress.String(), victoryItemsAddr.String()),
		nil,
	)
	assert.EqualValues(t, cadence.NewUInt64(0), supply)

	// assert that the account collection is empty
	length := test.ExecuteScriptAndCheck(
		t, b,
		victory_items.GetCollectionLengthScript(nftAddress.String(), victoryItemsAddr.String()),
		[][]byte{jsoncdc.MustEncode(cadence.NewAddress(victoryItemsAddr))},
	)
	assert.EqualValues(t, cadence.NewInt(0), length)

	t.Run("Should be able to mint a victoryItems", func(t *testing.T) {
		victory_items.MintItem(t, b, nftAddress, victoryItemsAddr, victoryItemsSigner, 
			victoryItemsAddr, typeID, brandID, dropID, contentHash, 
			startIssueNum, singleItemMintCount, singleItemMintCount, 
			metaURL, geoURL)

		// assert that the account collection is correct length
		length = test.ExecuteScriptAndCheck(
			t, b,
			victory_items.GetCollectionLengthScript(nftAddress.String(), victoryItemsAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(victoryItemsAddr))},
		)
		assert.EqualValues(t, cadence.NewInt(1), length)
	})
}

func TestBatchCreateVictoryItem(t *testing.T) {
	b := test.NewBlockchain()

	nftAddress, victoryItemsAddr, victoryItemsSigner := victory_items.DeployContracts(t, b)

	supply := test.ExecuteScriptAndCheck(
		t, b,
		victory_items.GetVictoryItemSupplyScript(nftAddress.String(), victoryItemsAddr.String()),
		nil,
	)
	assert.EqualValues(t, cadence.NewUInt64(0), supply)

	// assert that the account collection is empty
	length := test.ExecuteScriptAndCheck(
		t, b,
		victory_items.GetCollectionLengthScript(nftAddress.String(), victoryItemsAddr.String()),
		[][]byte{jsoncdc.MustEncode(cadence.NewAddress(victoryItemsAddr))},
	)
	assert.EqualValues(t, cadence.NewInt(0), length)

	t.Run("Should be able to mint batches of victoryItems", func(t *testing.T) {
		victory_items.MintItem(t, b, nftAddress, victoryItemsAddr, victoryItemsSigner, 
			victoryItemsAddr, typeID, brandID, dropID, contentHash, 
			startIssueNum, singleItemMintCount, maxIssueNum, 
			metaURL, geoURL)

		// assert that the account collection is correct length
		length = test.ExecuteScriptAndCheck(
			t, b,
			victory_items.GetCollectionLengthScript(nftAddress.String(), victoryItemsAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(victoryItemsAddr))},
		)
		assert.EqualValues(t, cadence.NewInt(1), length)

		// assert that the maxIssueNum is as expected
		maxIssue := test.ExecuteScriptAndCheck(
			t, b,
			victory_items.GetCollectibleMaxIssueScript(nftAddress.String(), victoryItemsAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(victoryItemsAddr)), 
					 jsoncdc.MustEncode(cadence.NewUInt64(0))},
		)
		assert.EqualValues(t, cadence.NewUInt32(maxIssueNum), maxIssue)

		// mint another batch
		victory_items.MintItem(t, b, nftAddress, victoryItemsAddr, victoryItemsSigner, 
			victoryItemsAddr, typeID, brandID, dropID, contentHash, 
			startIssueNum+1, maxIssueNum-1, maxIssueNum, 
			metaURL, geoURL)

		// assert that the account collection is correct length
		length = test.ExecuteScriptAndCheck(
			t, b,
			victory_items.GetCollectionLengthScript(nftAddress.String(), victoryItemsAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(victoryItemsAddr))},
		)
		assert.EqualValues(t, cadence.NewInt(maxIssueNum), length)

		// assert that the issue numbers are as expected
		// first item in the second batch should have issue 1 not 0
		issueNum := test.ExecuteScriptAndCheck(
			t, b,
			victory_items.GetCollectibleIssueNumScript(nftAddress.String(), victoryItemsAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(victoryItemsAddr)), 
					 jsoncdc.MustEncode(cadence.NewUInt64(1))},
		)
		assert.EqualValues(t, cadence.NewUInt32(1), issueNum)

		// max issue should be consistent
		maxIssue = test.ExecuteScriptAndCheck(
			t, b,
			victory_items.GetCollectibleMaxIssueScript(nftAddress.String(), victoryItemsAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(victoryItemsAddr)), 
					 jsoncdc.MustEncode(cadence.NewUInt64(1))},
		)
		assert.EqualValues(t, cadence.NewUInt32(maxIssueNum), maxIssue)

		// mint a third batch
		victory_items.MintItem(t, b, nftAddress, victoryItemsAddr, victoryItemsSigner, 
			victoryItemsAddr, typeID, brandID, dropID, contentHash, 
			startIssueNum, maxIssueNum, maxIssueNum, 
			metaURL, geoURL)

		// assert that the account collection is correct length
		length = test.ExecuteScriptAndCheck(
			t, b,
			victory_items.GetCollectionLengthScript(nftAddress.String(), victoryItemsAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(victoryItemsAddr))},
		)
		assert.EqualValues(t, cadence.NewInt(maxIssueNum*2), length)

		// assert that the issue numbers are as expected
		// first item in the third batch should have issue 0
		issueNum = test.ExecuteScriptAndCheck(
			t, b,
			victory_items.GetCollectibleIssueNumScript(nftAddress.String(), victoryItemsAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(victoryItemsAddr)), 
					 jsoncdc.MustEncode(cadence.NewUInt64(maxIssueNum))},
		)
		assert.EqualValues(t, cadence.NewUInt32(0), issueNum)

		// max issue should be consistent
		maxIssue = test.ExecuteScriptAndCheck(
			t, b,
			victory_items.GetCollectibleMaxIssueScript(nftAddress.String(), victoryItemsAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(victoryItemsAddr)), 
					 jsoncdc.MustEncode(cadence.NewUInt64(maxIssueNum))},
		)
		assert.EqualValues(t, cadence.NewUInt32(maxIssueNum), maxIssue)
	})
}

func TestTransferNFT(t *testing.T) {
	b := test.NewBlockchain()

	nftAddress, victoryItemsAddr, victoryItemsSigner := victory_items.DeployContracts(t, b)

	userAddress, userSigner, _ := test.CreateAccount(t, b)

	// create a new Collection
	t.Run("Should be able to create a new empty NFT Collection", func(t *testing.T) {
		victory_items.SetupAccount(t, b, userAddress, userSigner, nftAddress, victoryItemsAddr)

		length := test.ExecuteScriptAndCheck(
			t, b,
			victory_items.GetCollectionLengthScript(nftAddress.String(), victoryItemsAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(userAddress))},
		)
		assert.EqualValues(t, cadence.NewInt(0), length)
	})

	t.Run("Should not be able to withdraw an NFT that does not exist in a collection", func(t *testing.T) {
		nonExistentID := uint64(3333333)

		victory_items.TransferItem(
			t, b,
			nftAddress, victoryItemsAddr, 
			victoryItemsAddr, victoryItemsSigner,
			nonExistentID, userAddress, true,
		)
	})

	// transfer an NFT
	t.Run("Should be able to withdraw an NFT and deposit to another accounts collection", func(t *testing.T) {
		victory_items.MintItem(t, b, nftAddress, victoryItemsAddr, victoryItemsSigner, 
			victoryItemsAddr, typeID, brandID, dropID, contentHash, 
			startIssueNum, singleItemMintCount, singleItemMintCount, 
			metaURL, geoURL)

		// Cheat: we have minted at least one item, ID zero is valid
		victory_items.TransferItem(
			t, b,
			nftAddress, victoryItemsAddr, 
			victoryItemsAddr, victoryItemsSigner,
			0, userAddress, false,
		)
	})
}

func TestUpdateVictoryItem(t *testing.T) {
	b := test.NewBlockchain()

	nftAddress, victoryItemsAddr, victoryItemsSigner := victory_items.DeployContracts(t, b)

	t.Run("Should be able to update a victoryItem", func(t *testing.T) {
		victory_items.MintItem(t, b, nftAddress, victoryItemsAddr, victoryItemsSigner,
			victoryItemsAddr, typeID, brandID, dropID, contentHash, 
			startIssueNum, singleItemMintCount, singleItemMintCount, 
			metaURL, geoURL)

		// Cheat: we have minted at least one item, ID zero is valid
		victory_items.UpdateItemMeta(
			t, b,
			nftAddress, victoryItemsAddr, victoryItemsSigner,
			0, "http://someotherurl.com",
		)

		url := test.ExecuteScriptAndCheck(
			t, b,
			victory_items.GetCollectibleMetaURLScript(nftAddress.String(), victoryItemsAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(victoryItemsAddr)), 
					 jsoncdc.MustEncode(cadence.NewUInt64(0))},
		)
		assert.EqualValues(t, cadence.NewString("http://someotherurl.com"), url)

		victory_items.UpdateItemGeo(
			t, b,
			nftAddress, victoryItemsAddr, victoryItemsSigner,
			0, "http://someothergeourl.com",
		)

		url = test.ExecuteScriptAndCheck(
			t, b,
			victory_items.GetCollectibleGeoURLScript(nftAddress.String(), victoryItemsAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(victoryItemsAddr)), 
					 jsoncdc.MustEncode(cadence.NewUInt64(0))},
		)
		assert.EqualValues(t, cadence.NewString("http://someothergeourl.com"), url)

		victory_items.UpdateItemHash(
			t, b,
			nftAddress, victoryItemsAddr, victoryItemsSigner,
			0, newHash,
		)

		hash := test.ExecuteScriptAndCheckByteArray(
			t, b,
			victory_items.GetCollectibleHashScript(nftAddress.String(), victoryItemsAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(victoryItemsAddr)), 
					 jsoncdc.MustEncode(cadence.NewUInt64(0))},
		)
		assert.EqualValues(t, cadence.NewString(newHash), cadence.NewString(hex.EncodeToString(hash)))
	})
}

func TestBundleFunctions(t *testing.T) {
	b := test.NewBlockchain()

	nftAddress, victoryItemsAddr, victoryItemsSigner := victory_items.DeployContracts(t, b)

	t.Run("Should be able to mint multiple victoryItems", func(t *testing.T) {
		victory_items.MintItem(t, b, nftAddress, victoryItemsAddr, victoryItemsSigner,
			victoryItemsAddr, typeID, brandID, dropID, contentHash, 
			startIssueNum, maxIssueNum, maxIssueNum,
			metaURL, geoURL)

		// assert that the account collection is correct length
		length := test.ExecuteScriptAndCheck(
			t, b,
			victory_items.GetCollectionLengthScript(nftAddress.String(), victoryItemsAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(victoryItemsAddr))},
		)
		assert.EqualValues(t, cadence.NewInt(5), length)
	})

	t.Run("Should be able to get victoryItem IDs", func(t *testing.T) {
		ids := test.ExecuteScriptAndCheckUInt64Array(
			t, b,
			victory_items.GetCollectionIDsScript(nftAddress.String(), victoryItemsAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(victoryItemsAddr))},
		)
		// should have minted maxIssueNum items
		assert.EqualValues(t, maxIssueNum, len(ids))

		// IDs should be sequential integers
		for i, v := range ids {
			assert.EqualValues(t, i, v)
		}
	})

	t.Run("Should be able to create a bundle", func(t *testing.T) {
		bundleIDs := []uint64{0, 3}

		victory_items.CreateBundle(t, b, nftAddress, victoryItemsAddr,
			victoryItemsAddr, victoryItemsSigner, 
			bundleIDs,
			false,
		)

		// bundleID should be 0
		check_ids := test.ExecuteScriptAndCheckUInt64Array(
			t, b,
			victory_items.GetCollectionBundleIDsScript(nftAddress.String(), victoryItemsAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(victoryItemsAddr)), 
					 jsoncdc.MustEncode(cadence.NewUInt64(0))},
		)
		// should have bundled 2 items
		assert.EqualValues(t, 2, len(check_ids))

		// IDs should match those passed in
		for i, v := range check_ids {
			assert.EqualValues(t, v, bundleIDs[i])
		}
	})

	t.Run("Should be able to check if items are for sale", func(t *testing.T) {
		// Item ID 3 should be for sale
		for_sale := test.ExecuteScriptAndCheck(
			t, b,
			victory_items.GetCollectionForSaleScript(nftAddress.String(), victoryItemsAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(victoryItemsAddr)),
					 jsoncdc.MustEncode(cadence.NewUInt64(3))},
		)
		assert.EqualValues(t, cadence.NewBool(true), for_sale)

		// Item ID 1 should not be for sale
		for_sale = test.ExecuteScriptAndCheck(
			t, b,
			victory_items.GetCollectionForSaleScript(nftAddress.String(), victoryItemsAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(victoryItemsAddr)),
					 jsoncdc.MustEncode(cadence.NewUInt64(1))},
		)
		assert.EqualValues(t, cadence.NewBool(false), for_sale)
	})

	t.Run("Should be able to remove a bundle", func(t *testing.T) {
		// should not be able to remove a non-existent bundle
		nonExistentID := uint64(3333333)

		victory_items.RemoveBundle(t, b, nftAddress, victoryItemsAddr, victoryItemsSigner, 
									nonExistentID, true)

		// should be able to remove the bundle we created earlier
		victory_items.RemoveBundle(t, b, nftAddress, victoryItemsAddr, victoryItemsSigner, 
									uint64(0), false)

		// Item ID 3 should not be for sale
		for_sale := test.ExecuteScriptAndCheck(
			t, b,
			victory_items.GetCollectionForSaleScript(nftAddress.String(), victoryItemsAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(victoryItemsAddr)),
					 jsoncdc.MustEncode(cadence.NewUInt64(3))},
		)
		assert.EqualValues(t, cadence.NewBool(false), for_sale)
	})

	t.Run("Should be able to create multiple bundles", func(t *testing.T) {
		bundleIDs := []uint64{0, 3}

		victory_items.CreateBundle(t, b, nftAddress, victoryItemsAddr,
			victoryItemsAddr, victoryItemsSigner, 
			bundleIDs,
			false,
		)

		// bundleID should be 1
		check_ids := test.ExecuteScriptAndCheckUInt64Array(
			t, b,
			victory_items.GetCollectionBundleIDsScript(nftAddress.String(), victoryItemsAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(victoryItemsAddr)), 
					 jsoncdc.MustEncode(cadence.NewUInt64(1))},
		)
		// should have bundled 2 items
		assert.EqualValues(t, 2, len(check_ids))

		// IDs should match those passed in
		for i, v := range check_ids {
			assert.EqualValues(t, v, bundleIDs[i])
		}

		bundleIDs = []uint64{1, 2}

		victory_items.CreateBundle(t, b, nftAddress, victoryItemsAddr,
			victoryItemsAddr, victoryItemsSigner, 
			bundleIDs,
			false,
		)

		nextBundleID := test.ExecuteScriptAndCheck(
			t, b,
			victory_items.GetNextBundleIDScript(nftAddress.String(), victoryItemsAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(victoryItemsAddr))},
		)
		// should have created bundles with IDs 0, 1, 2 so next ID = 3
		assert.EqualValues(t, 3, nextBundleID)
	})

	t.Run("Should not be able to bundle the same item twice", func(t *testing.T) {
		bundleIDs := []uint64{0}

		victory_items.CreateBundle(t, b, nftAddress, victoryItemsAddr,
			victoryItemsAddr, victoryItemsSigner, 
			bundleIDs,
			true,
		)
	})

	t.Run("Should be able to remove all bundles", func(t *testing.T) {
		// should be able to remove the bundle we created earlier
		victory_items.RemoveAllBundles(t, b, nftAddress, victoryItemsAddr, victoryItemsSigner, false)

		// Item ID 0 should not be for sale
		for_sale := test.ExecuteScriptAndCheck(
			t, b,
			victory_items.GetCollectionForSaleScript(nftAddress.String(), victoryItemsAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(victoryItemsAddr)),
					 jsoncdc.MustEncode(cadence.NewUInt64(0))},
		)
		assert.EqualValues(t, cadence.NewBool(false), for_sale)

		// Item ID 1 should not be for sale
		for_sale = test.ExecuteScriptAndCheck(
			t, b,
			victory_items.GetCollectionForSaleScript(nftAddress.String(), victoryItemsAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(victoryItemsAddr)),
					 jsoncdc.MustEncode(cadence.NewUInt64(1))},
		)
		assert.EqualValues(t, cadence.NewBool(false), for_sale)

		nextBundleID := test.ExecuteScriptAndCheck(
			t, b,
			victory_items.GetNextBundleIDScript(nftAddress.String(), victoryItemsAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(victoryItemsAddr))},
		)
		// should not have reset the next ID (in case there are references to those IDs in the wild)
		assert.EqualValues(t, 3, nextBundleID)
	})

	t.Run("Should be able to create a bundle after removing all bundles", func(t *testing.T) {
		bundleIDs := []uint64{0, 3}

		victory_items.CreateBundle(t, b, nftAddress, victoryItemsAddr,
			victoryItemsAddr, victoryItemsSigner, 
			bundleIDs,
			false,
		)

		// bundleID should be 3
		check_ids := test.ExecuteScriptAndCheckUInt64Array(
			t, b,
			victory_items.GetCollectionBundleIDsScript(nftAddress.String(), victoryItemsAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(victoryItemsAddr)), 
					 jsoncdc.MustEncode(cadence.NewUInt64(3))},
		)
		// should have bundled 2 items
		assert.EqualValues(t, 2, len(check_ids))

		// IDs should match those passed in
		for i, v := range check_ids {
			assert.EqualValues(t, v, bundleIDs[i])
		}
	})
}

func TestMintOnDemandVictoryItem(t *testing.T) {
	b := test.NewBlockchain()

	nftAddress, victoryItemsAddr, victoryItemsSigner := victory_items.DeployContracts(t, b)

	supply := test.ExecuteScriptAndCheck(
		t, b,
		victory_items.GetVictoryItemSupplyScript(nftAddress.String(), victoryItemsAddr.String()),
		nil,
	)
	assert.EqualValues(t, cadence.NewUInt64(0), supply)

	// assert that the account collection is empty
	length := test.ExecuteScriptAndCheck(
		t, b,
		victory_items.GetCollectionLengthScript(nftAddress.String(), victoryItemsAddr.String()),
		[][]byte{jsoncdc.MustEncode(cadence.NewAddress(victoryItemsAddr))},
	)
	assert.EqualValues(t, cadence.NewInt(0), length)

	t.Run("Should be able to mint victoryItems on demand", func(t *testing.T) {
		victory_items.MintItem(t, b, nftAddress, victoryItemsAddr, victoryItemsSigner, 
			victoryItemsAddr, typeID, brandID, dropID, contentHash, 
			startIssueNum, singleItemMintCount, maxIssueNum, 
			metaURL, geoURL)

		// assert that the account collection is correct length
		length = test.ExecuteScriptAndCheck(
			t, b,
			victory_items.GetCollectionLengthScript(nftAddress.String(), victoryItemsAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(victoryItemsAddr))},
		)
		assert.EqualValues(t, cadence.NewInt(1), length)

		// assert that the maxIssueNum is as expected
		maxIssue := test.ExecuteScriptAndCheck(
			t, b,
			victory_items.GetCollectibleMaxIssueScript(nftAddress.String(), victoryItemsAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(victoryItemsAddr)), 
					 jsoncdc.MustEncode(cadence.NewUInt64(0))},
		)
		assert.EqualValues(t, cadence.NewUInt32(maxIssueNum), maxIssue)

		// mint another "on demand"
		victory_items.MintItemOnDemand(t, b, nftAddress, victoryItemsAddr, victoryItemsSigner, 
			victoryItemsAddr, 0, 1)

		// assert that the account collection is correct length
		length = test.ExecuteScriptAndCheck(
			t, b,
			victory_items.GetCollectionLengthScript(nftAddress.String(), victoryItemsAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(victoryItemsAddr))},
		)
		assert.EqualValues(t, cadence.NewInt(2), length)

		// assert that the issue numbers are as expected
		// first item in the second batch should have issue 1 not 0
		issueNum := test.ExecuteScriptAndCheck(
			t, b,
			victory_items.GetCollectibleIssueNumScript(nftAddress.String(), victoryItemsAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(victoryItemsAddr)), 
					 jsoncdc.MustEncode(cadence.NewUInt64(1))},
		)
		assert.EqualValues(t, cadence.NewUInt32(1), issueNum)

		// max issue should be consistent
		maxIssue = test.ExecuteScriptAndCheck(
			t, b,
			victory_items.GetCollectibleMaxIssueScript(nftAddress.String(), victoryItemsAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(victoryItemsAddr)), 
					 jsoncdc.MustEncode(cadence.NewUInt64(1))},
		)
		assert.EqualValues(t, cadence.NewUInt32(maxIssueNum), maxIssue)

		// assert that the metaURL was copied correctly
		url := test.ExecuteScriptAndCheck(
			t, b,
			victory_items.GetCollectibleMetaURLScript(nftAddress.String(), victoryItemsAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(victoryItemsAddr)), 
					 jsoncdc.MustEncode(cadence.NewUInt64(1))},
		)
		assert.EqualValues(t, cadence.NewString(metaURL), url)

	})
}
