import {
  contract,
  helpers as h,
  matchers,
  setup,
} from '@chainlink/test-helpers'
import { assert } from 'chai'
import { ethers } from 'ethers'
import { BlockhashStoreTestHelperFactory } from '../../ethers/v0.6/BlockhashStoreTestHelperFactory'

let personas: setup.Personas
const provider = setup.provider()
const blockhashStoreTestHelperFactory = new BlockhashStoreTestHelperFactory()

beforeAll(async () => {
  personas = await setup.users(provider).then((x) => x.personas)
})

describe('BlockhashStore', () => {
  let blockhashStoreTestHelper: contract.Instance<BlockhashStoreTestHelperFactory>

  const mainnetBlocks: {
    num: number
    rlpHeader: Uint8Array
    hash: string
  }[] = [
    {
      num: 10000467,
      rlpHeader: ethers.utils.arrayify(
        '0xf90215a058ee3c05e880cb25a3db92b9f1479c5453690ca97f9bcbb18d21965d3213578ea01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d4934794ea674fdde714fd979de3edf0f56aa9716b898ec8a0a448355652812a7d518b5c979a15bba02cfe4576d8eb61e8b5731ecc37f2bec6a0049f25ed97f9ed9a9c8521ab39cd2c48438d1d18c84dcab5bf494c19595bd462a0b1169f28bdbe5dd61ebc20b7a459be9d7fa898f5a3ba5fed6d502d94b9a8101bb901001000008180000210000080010001080310e004800c3040000060484000010804088050044302a500240041040010012120840002400005092000808000640012081000880010008040200208000004050800400002244044006041040040010890040504020040008004222502000800220000021800006400802036500000000400014640d00020002110000001440000001509543802080004210004100de04744a2810000000032250080810000502210c04289480800000423080800004000a020220030203000020001000000042c00420090000008003308459020e010a01000200190900040e81000040040000020000a8044001000202010000600c087086c49cadb1b57839898538398909483984b9e845eb02fbf94505059452d65746865726d696e652d6575312d34a06d0287c21536fac432714bd3f3712ff1a7e409faf1b10edac9b9547da1d4f7b188930531280477460c',
      ),
      hash:
        '0x4a65bcdf3466a16740b74849cc10fc57d4acb24cce148665482812699a400464',
    },
    {
      num: 10000468,
      rlpHeader: ethers.utils.arrayify(
        '0xf9020da04a65bcdf3466a16740b74849cc10fc57d4acb24cce148665482812699a400464a01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d493479404668ec2f57cc15c381b461b9fedab5d451c8f7fa0bcd4ddbb7125a5c06df96d862921dc0bba8664b3f759a233fe565a615c0ab3eaa0087ab379852c83e4b660de1668fc93333201ad0d233167ea6cef8bacaf5cba2aa0d81855037b2a6b56eba0c2ed129fb4102fb831b0baa187a0f6e0c155400f7855b9010080040040200000000010102081000000500040010408040800010110000000008000005808020000902021818000210000000000081100401000400014400001041008000020448800180128800008000200000420e01200000000000000011000001000020000208000b42200a0008000510200080200008c002018108010014030200000080000000002000010008000011008004003081000400080100803040080040300000002044080480000000000008080101000000050000000000840000002200040000a0080000442008006005502800000040008000890201002022402208002900020900000000080000100100201080000000003400000004887086d57541477ba839898548398968083989147845eb02fc28c73706964657230380b03ac53a076c676a0ab090b373b6242851a4beab7b8cdc9d3ebe211747a255b78c0278c42880ea13d40042dd1e6',
      ),
      hash:
        '0x00fd2589a272b85ffaf63223641571bf95891c936b7514ee4e87a593e52de7c9',
    },
    {
      num: 10000469,
      rlpHeader: ethers.utils.arrayify(
        '0xf90211a000fd2589a272b85ffaf63223641571bf95891c936b7514ee4e87a593e52de7c9a01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347945a0b54d5dc17e0aadc383d2db43b0a0d3e029c4ca01b28d3b4e4d3442a9e6caed9f80b6a639bce6f3a283d4e004e6bb44e483ceeeba067c00d9067bc023b8fab8e3afd1bc0f2470f08003bdf9f167fbfeede2422ac4ea09d8b344d9ab1b7288f8c5e644c53b1a5288da6d6ee0e388ec76f291def48da15b90100c462095870a26a0804132e208110329710d459054558159c103208d44820002016108136199200061063699d8400254a22828c11b5512e3303c98ec7747cc02d00161880c2f2c580e806bccc04805190265f096601342058020a8324c277735d8202220412f03303201252a3000038883a4bb0010e6b004408306232150a84d110100d0c4b9d228022812602c05c801d20500d4ed10010ce2400428a96950a98050c00e603292a806c4983b25814880000440a23821191121996410c5110c949616c2066a4a0488087d4c226c14208042c00d609b5cc44051400219d93626818728612a9b18690e03c902014a900e0018828011494b80d4708799b0d8a83cace87086e64fefefb48839898558398968083986664845eb02fc7906574682d70726f2d687a662d74303032a09f1918a362b55ebd072cc9548fb74f89301d41c2a1feb13c08a1c2c3cb0606d88810dfa530069367fb',
      ),
      hash:
        '0x325fde74e261fc483a16506bbc711b645b043ad24c7db0136845a1be262cf0c9',
    },
    {
      num: 10000470,
      rlpHeader: ethers.utils.arrayify(
        '0xf90215a0325fde74e261fc483a16506bbc711b645b043ad24c7db0136845a1be262cf0c9a01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d4934794ea674fdde714fd979de3edf0f56aa9716b898ec8a020647cfa35563093442a12d80bf2bacb83da1de8340366677f3822591a334ccea066ad7285f6c5b6407f62c6b65a83aeaaa71ad9a97c2bb15139140f2dbb60f7e0a0c0e633851d0b5ce661ecc054517425e82425fcc6170db9693e5b5a6dd5ef6d6bb90100c0c000c1520708182080c8e461891c2402800a80d44a00034259414012012a5006a1416331181504902044960808f1129018800311621e920886804693749b10542400142e984580ccba634881c4156962200ecfb004000005468db44842781c59923110262660802315006106388b028412c42c000820c508e66b7851fa68002008144cd7860cd884280802915163399c168d5a11b0649486084110149469a1e61c31134204b903206566885180bc0426c0c6c0a4d408e182242f08180d204c624a040248425041ac028010d088820402ba4bd38c2d1215829300543465603822110500811290490148049300040e000c280086a09e8100089818ce480a887e87086c4965bf3c8a839898568398705c839847d2845eb02fe994505059452d65746865726d696e652d6575312d35a09d8ae288d0eede524f3ef5e6cfcc5ba07f380bc695bb71578a7b91cfa517071b8859d0976006378e52',
      ),
      hash:
        '0x5cf096dfd1fc2d2947a96fdec5377ab7beaa0eb00c80728a3b96f0864cec506a',
    },
  ]

  const deployment = setup.snapshot(provider, async () => {
    blockhashStoreTestHelper = await blockhashStoreTestHelperFactory
      .connect(personas.Default)
      .deploy()
  })

  beforeEach(async () => {
    await deployment()

    const [lastBlock] = mainnetBlocks.slice(-1)
    await blockhashStoreTestHelper
      .connect(personas.Default)
      .godmodeSetHash(lastBlock.num, lastBlock.hash)
    assert.strictEqual(
      await blockhashStoreTestHelper.getBlockhash(lastBlock.num),
      lastBlock.hash,
    )
  })

  it('getBlockhash reverts for unknown blockhashes', async () => {
    await matchers.evmRevert(
      blockhashStoreTestHelper.getBlockhash(99999999),
      'blockhash not found in store',
    )
  })

  it('storeVerifyHeader records valid blockhashes', async () => {
    for (let i = mainnetBlocks.length - 2; i >= 0; i--) {
      await blockhashStoreTestHelper
        .connect(personas.Default)
        .storeVerifyHeader(mainnetBlocks[i].num, mainnetBlocks[i + 1].rlpHeader)
      assert.strictEqual(
        await blockhashStoreTestHelper.getBlockhash(mainnetBlocks[i].num),
        mainnetBlocks[i].hash,
      )
    }
  })

  it('storeVerifyHeader rejects unknown headers', async () => {
    const unknownBlock = mainnetBlocks[0]
    await matchers.evmRevert(
      blockhashStoreTestHelper
        .connect(personas.Default)
        .storeVerifyHeader(unknownBlock.num - 1, unknownBlock.rlpHeader),
      'header has unknown blockhash',
    )
  })

  it('storeVerifyHeader rejects corrupted headers', async () => {
    const [lastBlock] = mainnetBlocks.slice(-1)
    const modifiedHeader = new Uint8Array(lastBlock.rlpHeader)
    modifiedHeader[137] += 1
    await matchers.evmRevert(
      blockhashStoreTestHelper
        .connect(personas.Default)
        .storeVerifyHeader(lastBlock.num - 1, modifiedHeader),
      'header has unknown blockhash',
    )
  })

  it('store accepts recent block numbers', async () => {
    await h.mineBlock(provider)

    const n = (await provider.getBlockNumber()) - 1
    await blockhashStoreTestHelper.connect(personas.Default).store(n)

    assert.equal(
      await blockhashStoreTestHelper.getBlockhash(n),
      (await provider.getBlock(n)).hash,
    )
  })

  it('store rejects future block numbers', async () => {
    await matchers.evmRevert(
      blockhashStoreTestHelper.connect(personas.Default).store(99999999999),
      'blockhash(n) failed',
    )
  })

  it('store rejects old block numbers', async () => {
    for (let i = 0; i < 300; i++) {
      await h.mineBlock(provider)
    }

    await matchers.evmRevert(
      blockhashStoreTestHelper
        .connect(personas.Default)
        .store((await provider.getBlockNumber()) - 256),
      'blockhash(n) failed',
    )
  })

  it('storeEarliest works', async () => {
    for (let i = 0; i < 300; i++) {
      await h.mineBlock(provider)
    }

    await blockhashStoreTestHelper.connect(personas.Default).storeEarliest()

    const n = (await provider.getBlockNumber()) - 256
    assert.equal(
      await blockhashStoreTestHelper.getBlockhash(n),
      (await provider.getBlock(n)).hash,
    )
  })
})
