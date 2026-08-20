package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeAssetType_ProductionOrdinalsPassThroughUnchanged(t *testing.T) {
	// Production / custom ordinals (1..6) must not be remapped
	require.Equal(t, uint8(AssetTypeCustom), NormalizeAssetType(uint8(AssetTypeCustom)))
	require.Equal(t, uint8(AssetTypeERC20), NormalizeAssetType(uint8(AssetTypeERC20)))
	require.Equal(t, uint8(AssetTypeERC721), NormalizeAssetType(uint8(AssetTypeERC721)))
	require.Equal(t, uint8(AssetTypeERC1155), NormalizeAssetType(uint8(AssetTypeERC1155)))
	require.Equal(t, uint8(AssetTypeEnygma), NormalizeAssetType(uint8(AssetTypeEnygma)))
	require.Equal(t, uint8(AssetTypeDvpERC721), NormalizeAssetType(uint8(AssetTypeDvpERC721)))
	require.Equal(t, uint8(AssetTypeDvpERC1155), NormalizeAssetType(uint8(AssetTypeDvpERC1155)))
}

func TestNormalizeAssetType_TestOrdinalsFoldBackToProductionBase(t *testing.T) {
	// Each test-only ordinal folds back to the production token it mirrors
	require.Equal(t, uint8(AssetTypeERC20), NormalizeAssetType(uint8(AssetTypeERC20Test)))
	require.Equal(t, uint8(AssetTypeERC721), NormalizeAssetType(uint8(AssetTypeERC721Test)))
	require.Equal(t, uint8(AssetTypeERC1155), NormalizeAssetType(uint8(AssetTypeERC1155Test)))
	require.Equal(t, uint8(AssetTypeEnygma), NormalizeAssetType(uint8(AssetTypeEnygmaTest)))
	require.Equal(t, uint8(AssetTypeDvpERC721), NormalizeAssetType(uint8(AssetTypeDvpERC721Test)))
	require.Equal(t, uint8(AssetTypeDvpERC1155), NormalizeAssetType(uint8(AssetTypeDvpERC1155Test)))
}

func TestNormalizeAssetType_OutOfRangeOrdinalPassesThroughUnchanged(t *testing.T) {
	// Ordinal 13 (beyond DvpERC1155Test) must not be remapped
	require.Equal(t, uint8(13), NormalizeAssetType(13))
}

func TestNormalizeAssetType_EveryOrdinalResolvesToMappedString(t *testing.T) {
	// Every normalized ordinal (1..12) must resolve to a non-empty ercStandard string
	for in := uint8(AssetTypeERC20); in <= uint8(AssetTypeDvpERC1155Test); in++ {
		require.NotEmpty(
			t,
			AssetTypeToString[NormalizeAssetType(in)],
			"ordinal %d normalized to an unmapped ercStandard string",
			in,
		)
	}
}
