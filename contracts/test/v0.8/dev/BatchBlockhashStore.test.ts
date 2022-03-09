import { assert, expect } from 'chai'
import { Contract, Signer } from 'ethers'
import { ethers } from 'hardhat'
import * as rlp from 'rlp'

function range(size: number, startAt = 0) {
  return [...Array(size).keys()].map((i) => i + startAt)
}

describe('BatchBlockhashStore', () => {
  let blockhashStore: Contract
  let batchBHS: Contract
  let owner: Signer

  beforeEach(async () => {
    const accounts = await ethers.getSigners()
    owner = accounts[0]

    const bhFactory = await ethers.getContractFactory(
      'src/v0.6/dev/BlockhashStore.sol:BlockhashStore',
      accounts[0],
    )

    blockhashStore = await bhFactory.deploy()

    const batchBHSFactory = await ethers.getContractFactory(
      'src/v0.8/dev/BatchBlockhashStore.sol:BatchBlockhashStore',
      accounts[0],
    )

    batchBHS = await batchBHSFactory.deploy(blockhashStore.address)

    // Mine some blocks so that we have some blockhashes to store.
    for (let i = 0; i < 10; i++) {
      await ethers.provider.send('evm_mine', [])
    }
  })

  describe('#store', () => {
    it('stores batches of blocknumbers available through the blockhash instruction', async () => {
      const latestBlock = await ethers.provider.send('eth_blockNumber', [])
      const bottomBlock = latestBlock - 5
      const numBlocks = 3
      await batchBHS.connect(owner).store(range(numBlocks, bottomBlock))

      // Mine some blocks to confirm the store batch tx above.
      for (let i = 0; i < 2; i++) {
        await ethers.provider.send('evm_mine', [])
      }

      // check the bhs if it was stored
      for (let i = bottomBlock; i < bottomBlock + numBlocks; i++) {
        const actualBh = await blockhashStore.connect(owner).getBlockhash(i)
        const expectedBh = (await ethers.provider.getBlock(i)).hash
        expect(expectedBh).to.equal(actualBh)
      }
    })

    it('reverts if block number is too far back for blockhash instruction', async () => {
      // blockhash(n) fails if n is more than 256 blocks behind the current block in which
      // the instruction is executing.
      for (let i = 0; i < 256; i++) {
        await ethers.provider.send('evm_mine', [])
      }

      expect(batchBHS.connect(owner).store([1, 2, 3])).to.be.revertedWith(
        'blockhash(n) failed',
      )
    })
  })

  describe('#getBlockhashes', () => {
    it('fetches blockhashes of a batch of block numbers', async () => {
      // Store a bunch of block hashes
      const latestBlock = await ethers.provider.send('eth_blockNumber', [])
      const bottomBlock = latestBlock - 5
      const numBlocks = 3
      await batchBHS.connect(owner).store(range(numBlocks, bottomBlock))

      // Mine some blocks to confirm the store batch tx above.
      for (let i = 0; i < 2; i++) {
        await ethers.provider.send('evm_mine', [])
      }

      // fetch the blocks in a batch
      const actualBlockhashes = await batchBHS
        .connect(owner)
        .getBlockhashes(range(numBlocks, bottomBlock))
      let expectedBlockhashes = []
      for (let i = bottomBlock; i < bottomBlock + numBlocks; i++) {
        const block = await ethers.provider.send('eth_getBlockByNumber', [
          '0x' + i.toString(16),
          false,
        ])
        expectedBlockhashes.push(block.hash)
      }
      assert.deepEqual(actualBlockhashes, expectedBlockhashes)
    })

    it('reverts if one of the block numbers does not have an associated blockhash', async () => {
      const latestBlock = await ethers.provider.send('eth_blockNumber', [])
      const bottomBlock = latestBlock - 5
      const numBlocks = 3
      expect(
        batchBHS.connect(owner).getBlockhashes(range(numBlocks, bottomBlock)),
      ).to.be.revertedWith('blockhash not found in store')
    })
  })

  describe('#storeVerifyHeader', () => {
    it('stores batches of blocknumbers using storeVerifyHeader', async () => {
      // Store a single blockhash and go backwards from there using storeVerifyHeader
      const latestBlock = await ethers.provider.send('eth_blockNumber', [])
      await batchBHS.connect(owner).store([latestBlock])

      await ethers.provider.send('evm_mine', [])

      const numBlocks = 3
      const startBlock = latestBlock - 1
      const blockNumbers = range(
        numBlocks + 1,
        startBlock - numBlocks,
      ).reverse()
      let blockHeaders = []
      let expectedBlockhashes = []
      for (let i of blockNumbers) {
        const block = await ethers.provider.send('eth_getBlockByNumber', [
          '0x' + (i + 1).toString(16),
          false,
        ])
        // eip 1559 header - switch to this if we upgrade hardhat
        // and use post-london forks of ethereum.
        // const encodedHeader = rlp.encode([
        // 	block.parentHash,
        // 	block.sha3Uncles,
        // 	ethers.utils.arrayify(block.miner),
        // 	block.stateRoot,
        // 	block.transactionsRoot,
        // 	block.receiptsRoot,
        // 	block.logsBloom,
        // 	block.difficulty,
        // 	block.number,
        // 	block.gasLimit,
        // 	block.gasUsed == '0x0' ? '0x' : block.gasUsed,
        // 	block.timestamp,
        // 	block.extraData,
        // 	block.mixHash,
        // 	block.nonce,
        // 	block.baseFeePerGas,
        // ])
        // pre-london block header serialization
        const encodedHeader = rlp.encode([
          block.parentHash,
          block.sha3Uncles,
          ethers.utils.arrayify(block.miner),
          block.stateRoot,
          block.transactionsRoot,
          block.receiptsRoot,
          block.logsBloom,
          block.difficulty,
          block.number,
          block.gasLimit,
          block.gasUsed == '0x0' ? '0x' : block.gasUsed,
          block.timestamp,
          block.extraData,
          block.mixHash,
          block.nonce,
        ])
        blockHeaders.push('0x' + encodedHeader.toString('hex'))
        expectedBlockhashes.push(
          (
            await ethers.provider.send('eth_getBlockByNumber', [
              '0x' + i.toString(16),
              false,
            ])
          ).hash,
        )
      }
      await batchBHS
        .connect(owner)
        .storeVerifyHeader(blockNumbers, blockHeaders)

      // fetch blocks that were just stored and assert correctness
      const actualBlockhashes = await batchBHS
        .connect(owner)
        .getBlockhashes(blockNumbers)
      assert.deepEqual(actualBlockhashes, expectedBlockhashes)
    })

    describe('bad input', () => {
      it('reverts on mismatched input array sizes', async () => {
        // Store a single blockhash and go backwards from there using storeVerifyHeader
        const latestBlock = await ethers.provider.send('eth_blockNumber', [])
        await batchBHS.connect(owner).store([latestBlock])

        await ethers.provider.send('evm_mine', [])

        const numBlocks = 3
        const startBlock = latestBlock - 1
        const blockNumbers = range(
          numBlocks + 1,
          startBlock - numBlocks,
        ).reverse()
        let blockHeaders = []
        let expectedBlockhashes = []
        for (let i of blockNumbers) {
          const block = await ethers.provider.send('eth_getBlockByNumber', [
            '0x' + (i + 1).toString(16),
            false,
          ])
          const encodedHeader = rlp.encode([
            block.parentHash,
            block.sha3Uncles,
            ethers.utils.arrayify(block.miner),
            block.stateRoot,
            block.transactionsRoot,
            block.receiptsRoot,
            block.logsBloom,
            block.difficulty,
            block.number,
            block.gasLimit,
            block.gasUsed == '0x0' ? '0x' : block.gasUsed,
            block.timestamp,
            block.extraData,
            block.mixHash,
            block.nonce,
            block.baseFeePerGas,
          ])
          blockHeaders.push('0x' + encodedHeader.toString('hex'))
          expectedBlockhashes.push(
            (
              await ethers.provider.send('eth_getBlockByNumber', [
                '0x' + i.toString(16),
                false,
              ])
            ).hash,
          )
        }
        // remove last element to simulate different input array sizes
        blockHeaders.pop()
        expect(
          batchBHS.connect(owner).storeVerifyHeader(blockNumbers, blockHeaders),
        ).to.be.revertedWith('input array arg lengths mismatch')
      })

      it('reverts on bad block header input', async () => {
        // Store a single blockhash and go backwards from there using storeVerifyHeader
        const latestBlock = await ethers.provider.send('eth_blockNumber', [])
        await batchBHS.connect(owner).store([latestBlock])

        await ethers.provider.send('evm_mine', [])

        const numBlocks = 3
        const startBlock = latestBlock - 1
        const blockNumbers = range(
          numBlocks + 1,
          startBlock - numBlocks,
        ).reverse()
        let blockHeaders = []
        let expectedBlockhashes = []
        for (let i of blockNumbers) {
          const block = await ethers.provider.send('eth_getBlockByNumber', [
            '0x' + (i + 1).toString(16),
            false,
          ])
          const encodedHeader = rlp.encode([
            block.parentHash,
            block.sha3Uncles,
            ethers.utils.arrayify(block.miner),
            block.stateRoot,
            block.transactionsRoot,
            block.receiptsRoot,
            block.logsBloom,
            block.difficulty,
            block.number,
            block.gasLimit,
            block.gasUsed, // incorrect: in cases where it's 0x0 it should be 0x instead.
            block.timestamp,
            block.extraData,
            block.mixHash,
            block.nonce,
            block.baseFeePerGas,
          ])
          blockHeaders.push('0x' + encodedHeader.toString('hex'))
          expectedBlockhashes.push(
            (
              await ethers.provider.send('eth_getBlockByNumber', [
                '0x' + i.toString(16),
                false,
              ])
            ).hash,
          )
        }
        expect(
          batchBHS.connect(owner).storeVerifyHeader(blockNumbers, blockHeaders),
        ).to.be.revertedWith('header has unknown blockhash')
      })
    })
  })
})
