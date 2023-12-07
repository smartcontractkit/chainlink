package utils

import (
	"encoding/binary"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var legacyV1FeedIDs = []FeedID{
	// Arbitrum mainnet (prod)
	mustHexToFeedID("0xb43dc495134fa357725f93539511c5a4febeadf56e7c29c96566c825094f0b20"),
	mustHexToFeedID("0xe65b31c6d5b9bdff43a8194dc5b2edc6914ddbc5e9f9e9521f605fc3738fabf5"),
	mustHexToFeedID("0x30f9926cdef3de98995fb38a100d5c582ae025ebbb8f9a931500596ce080280a"),
	mustHexToFeedID("0x0f49a4533a64c7f53bfdf5e86d791620d93afdec00cfe1896548397b0f4ec81c"),
	mustHexToFeedID("0x2cdd4aea8298f5d2e7f8505b91e3313e3aa04376a81f401b4a48c5aab78ee5cf"),
	mustHexToFeedID("0x5f82d154119f4251d83b2a58bf61c9483c84241053038a2883abf16ed4926433"),
	mustHexToFeedID("0x74aca63821bf7ead199e924d261d277cbec96d1026ab65267d655c51b4536914"),
	mustHexToFeedID("0x64ee16b94fdd72d0b3769955445cc82d6804573c22f0f49b67cd02edd07461e7"),
	mustHexToFeedID("0x95241f154d34539741b19ce4bae815473fd1b2a90ac3b4b023a692f31edfe90e"),
	mustHexToFeedID("0x297cc1e1ee5fc2f45dff1dd11a46694567904f4dbc596c7cc216d6c688605a1b"),
	// // Arbitrum mainnet (staging)
	mustHexToFeedID("0x62ce6a99c4bebb150191d7b72f7a0c0206af00baca480ab007caa4b5bf4bf02a"),
	mustHexToFeedID("0x984126712e6a8b5b4fe138c49b29483a12e77b5cb3213a0769252380c57480e4"),
	mustHexToFeedID("0xb74f650d9cae6259ab4212f76abe746600be3a4926947725ed107943915346c1"),
	mustHexToFeedID("0xa0098c4c06cbab05b2598aecad0cbf49d44780c56d40514e09fd7a9e76a2db00"),
	mustHexToFeedID("0x2206b467d04656a8a83af43a428d6b66f787162db629f9caed0c12b54a32998e"),
	mustHexToFeedID("0x55488e61b59ea629df66698c8eea1390f0aedc24942e074a6d565569fb90afde"),
	mustHexToFeedID("0x98d66aab30d62d044cc55ffccb79ae35151348f40ff06a98c92001ed6ec8e886"),
	mustHexToFeedID("0x2e768c0eca65d0449ee825b8a921349501339a2487c02146f77611ae01c31a50"),
	mustHexToFeedID("0xb29931d9fe1e9fc023b4d2f0f1789c8b5e21aabf389f86f9702241a0178345dd"),
	mustHexToFeedID("0xd8b8cfc1e2dd75116e5792d11810d830ef48843fd44e1633385e81157f8da6b5"),
	mustHexToFeedID("0x09f8d0caff8cecb7f5e493d4de2ab98b4392f6d07923cd19b2cb524779301b85"),
	mustHexToFeedID("0xe645924bbf507304dc4bd37f02c8dac73da3b7eb67378de98cfc59f17ba6774a"),
	// Arbitrum testnet (production)
	mustHexToFeedID("0x695be66b6a7979f2b3ed33a3d718eabebaf0a881f1f6598b5530875b7e8150ab"),
	mustHexToFeedID("0x259b566b9d3c64d1e4a8656e2d6fd4c08e19f9fa9637ae76d52e428d07cca8e9"),
	mustHexToFeedID("0x26c16f2054b7a1d77ae83a0429dace9f3000ba4dbf1690236e8f575742e98f66"),
	mustHexToFeedID("0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"),
	mustHexToFeedID("0xbf1febc8c335cb236c1995c1007a928a3f7ae8307a1a20cb31334e6d316c62d1"),
	mustHexToFeedID("0x4ce52cf28e49f4673198074968aeea280f13b5f897c687eb713bcfc1eeab89ba"),
	mustHexToFeedID("0xb21d58dccab05dcea22ab780ca010c4bec34e61ce7310e30f4ad0ff8c1621d27"),
	mustHexToFeedID("0x5ad0d18436dd95672e69903efe95bdfb43a05cb55e8965c5af93db8170c8820c"),
	mustHexToFeedID("0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"),
	mustHexToFeedID("0x14e044f932bb959cc2aa8dc1ba110c09224e639aae00264c1ffc2a0830904a3c"),
	mustHexToFeedID("0x555344432d5553442d415242495452554d2d544553544e455400000000000000"),
	mustHexToFeedID("0x12be1859ee43f46bab53750915f20855f54e891f88ddd524f26a72d6f4deed1d"),
	// // Arbitrum testnet (staging)
	mustHexToFeedID("0x8837f28f5172f18071f164b8540fe8c95162dc0051e31005023fadc1cd9c4b50"),
	mustHexToFeedID("0xd130b5acd88b47eb7c372611205d5a9ca474829a2719e396ab1eb4f956674e4e"),
	mustHexToFeedID("0x6d2f5a4b3ba6c1953b4bb636f6ad03aec01b6222274f8ca1e39e53ee12a8cdf3"),
	mustHexToFeedID("0x6962e629c3a0f5b7e3e9294b0c283c9b20f94f1c89c8ba8c1ee4650738f20fb2"),
	mustHexToFeedID("0x557b817c6be7392364cef0dd11007c43caea1de78ce42e4f1eadc383e7cb209c"),
	mustHexToFeedID("0x3250b5dd9491cb11138048d070b8636c35d96fff29671dc68b0723ad41f53433"),
	mustHexToFeedID("0x3781c2691f6980dc66a72c03a32edb769fe05a9c9cb729cd7e96ecfd89450a0a"),
	mustHexToFeedID("0xbbbf52c5797cc86d6bd9413d59ec624f07baf5045290ecd5ac6541d5a7ffd234"),
	mustHexToFeedID("0xf753e1201d54ac94dfd9334c542562ff7e42993419a661261d010af0cbfd4e34"),
	mustHexToFeedID("0x2489ce4577e814d6794218a13ef3c04cac976f991305400a4c0a1ddcffb90357"),
	mustHexToFeedID("0xa5b07943b89e2c278fc8a2754e2854316e03cb959f6d323c2d5da218fb6b0ff8"),
	mustHexToFeedID("0x1c2c0dfac0eb2aae2c05613f0d677daae164cdd406bd3dd6153d743302ce56e8"),
}

var legacyV1FeedIDM map[FeedID]struct{}

func init() {
	legacyV1FeedIDM = make(map[FeedID]struct{})
	for _, feedID := range legacyV1FeedIDs {
		legacyV1FeedIDM[feedID] = struct{}{}
	}
}

func mustHexToFeedID(s string) FeedID {
	f := new(FeedID)
	if err := f.UnmarshalText([]byte(s)); err != nil {
		panic(err)
	}
	return *f
}

type FeedVersion uint16

const (
	_ FeedVersion = iota
	REPORT_V1
	REPORT_V2
	REPORT_V3
	_
)

type FeedID [32]byte

func BytesToFeedID(b []byte) FeedID {
	return (FeedID)(utils.BytesToHash(b))
}

func (f FeedID) Hex() string { return (utils.Hash)(f).Hex() }

func (f FeedID) String() string { return (utils.Hash)(f).String() }

func (f *FeedID) UnmarshalText(input []byte) error {
	return (*utils.Hash)(f).UnmarshalText(input)
}

func (f FeedID) Version() FeedVersion {
	if _, exists := legacyV1FeedIDM[f]; exists {
		return REPORT_V1
	}
	return FeedVersion(binary.BigEndian.Uint16(f[:2]))
}

func (f FeedID) IsV1() bool { return f.Version() == REPORT_V1 }
func (f FeedID) IsV2() bool { return f.Version() == REPORT_V2 }
func (f FeedID) IsV3() bool { return f.Version() == REPORT_V3 }
