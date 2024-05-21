import { ethers, network } from 'hardhat'
import { BigNumber, Contract } from 'ethers'
import { expect } from 'chai'
import { SignerWithAddress } from '@nomiclabs/hardhat-ethers/signers'

describe('ScrollSequencerUptimeFeed', () => {
  let l2CrossDomainMessenger: Contract
  let scrollUptimeFeed: Contract
  let uptimeFeedConsumer: Contract
  let deployer: SignerWithAddress
  let l1Owner: SignerWithAddress
  let l2Messenger: SignerWithAddress
  let dummy: SignerWithAddress
  const gasUsedDeviation = 100
  const initialStatus = 0

  before(async () => {
    const accounts = await ethers.getSigners()
    deployer = accounts[0]
    l1Owner = accounts[1]
    dummy = accounts[3]

    const l2CrossDomainMessengerFactory = await ethers.getContractFactory(
      'src/v0.8/l2ep/test/mocks/scroll/MockScrollL2CrossDomainMessenger.sol:MockScrollL2CrossDomainMessenger',
      deployer,
    )

    l2CrossDomainMessenger = await l2CrossDomainMessengerFactory.deploy()

    // Pretend we're on L2
    await network.provider.request({
      method: 'hardhat_impersonateAccount',
      params: [l2CrossDomainMessenger.address],
    })
    l2Messenger = await ethers.getSigner(l2CrossDomainMessenger.address)
    // Credit the L2 messenger with some ETH
    await dummy.sendTransaction({
      to: l2Messenger.address,
      value: ethers.utils.parseEther('10'),
    })
  })

  beforeEach(async () => {
    const scrollSequencerStatusRecorderFactory =
      await ethers.getContractFactory(
        'src/v0.8/l2ep/dev/scroll/ScrollSequencerUptimeFeed.sol:ScrollSequencerUptimeFeed',
        deployer,
      )
    scrollUptimeFeed = await scrollSequencerStatusRecorderFactory.deploy(
      l1Owner.address,
      l2CrossDomainMessenger.address,
      initialStatus,
    )

    // Set mock sender in mock L2 messenger contract
    await l2CrossDomainMessenger.setSender(l1Owner.address)

    // Mock consumer
    const statusFeedConsumerFactory = await ethers.getContractFactory(
      'src/v0.8/tests/FeedConsumer.sol:FeedConsumer',
      deployer,
    )
    uptimeFeedConsumer = await statusFeedConsumerFactory.deploy(
      scrollUptimeFeed.address,
    )
  })

  describe('constructor', () => {
    it('should have been deployed with the correct initial state', async () => {
      const l1Sender = await scrollUptimeFeed.l1Sender()
      expect(l1Sender).to.equal(l1Owner.address)
      const { roundId, answer } = await scrollUptimeFeed.latestRoundData()
      expect(roundId).to.equal(1)
      expect(answer).to.equal(initialStatus)
    })
  })

  describe('#updateStatus', () => {
    it('should revert if called by an address that is not the L2 Cross Domain Messenger', async () => {
      const timestamp = await scrollUptimeFeed.latestTimestamp()
      expect(
        scrollUptimeFeed.connect(dummy).updateStatus(true, timestamp),
      ).to.be.revertedWith('InvalidSender')
    })

    it('should revert if called by an address that is not the L2 Cross Domain Messenger and is not the L1 sender', async () => {
      const timestamp = await scrollUptimeFeed.latestTimestamp()
      await l2CrossDomainMessenger.setSender(dummy.address)
      expect(
        scrollUptimeFeed.connect(dummy).updateStatus(true, timestamp),
      ).to.be.revertedWith('InvalidSender')
    })

    it(`should update status when status has not changed and incoming timestamp is the same as latest`, async () => {
      const timestamp = await scrollUptimeFeed.latestTimestamp()
      let tx = await scrollUptimeFeed
        .connect(l2Messenger)
        .updateStatus(true, timestamp)
      await expect(tx)
        .to.emit(scrollUptimeFeed, 'AnswerUpdated')
        .withArgs(1, 2 /** roundId */, timestamp)
      expect(await scrollUptimeFeed.latestAnswer()).to.equal(1)

      const latestRoundBeforeUpdate = await scrollUptimeFeed.latestRoundData()

      tx = await scrollUptimeFeed
        .connect(l2Messenger)
        .updateStatus(true, timestamp.add(200))

      // Submit another status update with the same status
      const latestBlock = await ethers.provider.getBlock('latest')

      await expect(tx)
        .to.emit(scrollUptimeFeed, 'RoundUpdated')
        .withArgs(1, latestBlock.timestamp)
      expect(await scrollUptimeFeed.latestAnswer()).to.equal(1)
      expect(await scrollUptimeFeed.latestTimestamp()).to.equal(timestamp)

      // Verify that latest round has been properly updated
      const latestRoundDataAfterUpdate =
        await scrollUptimeFeed.latestRoundData()
      expect(latestRoundDataAfterUpdate.roundId).to.equal(
        latestRoundBeforeUpdate.roundId,
      )
      expect(latestRoundDataAfterUpdate.answer).to.equal(
        latestRoundBeforeUpdate.answer,
      )
      expect(latestRoundDataAfterUpdate.startedAt).to.equal(
        latestRoundBeforeUpdate.startedAt,
      )
      expect(latestRoundDataAfterUpdate.answeredInRound).to.equal(
        latestRoundBeforeUpdate.answeredInRound,
      )
      expect(latestRoundDataAfterUpdate.updatedAt).to.equal(
        latestBlock.timestamp,
      )
    })

    it(`should update status when status has changed and incoming timestamp is newer than the latest`, async () => {
      let timestamp = await scrollUptimeFeed.latestTimestamp()
      let tx = await scrollUptimeFeed
        .connect(l2Messenger)
        .updateStatus(true, timestamp)
      await expect(tx)
        .to.emit(scrollUptimeFeed, 'AnswerUpdated')
        .withArgs(1, 2 /** roundId */, timestamp)
      expect(await scrollUptimeFeed.latestAnswer()).to.equal(1)

      // Submit another status update, different status, newer timestamp should update
      timestamp = timestamp.add(2000)
      tx = await scrollUptimeFeed
        .connect(l2Messenger)
        .updateStatus(false, timestamp)
      await expect(tx)
        .to.emit(scrollUptimeFeed, 'AnswerUpdated')
        .withArgs(0, 3 /** roundId */, timestamp)
      expect(await scrollUptimeFeed.latestAnswer()).to.equal(0)
      expect(await scrollUptimeFeed.latestTimestamp()).to.equal(timestamp)
    })

    it(`should update status when status has changed and incoming timestamp is the same as latest`, async () => {
      const timestamp = await scrollUptimeFeed.latestTimestamp()
      let tx = await scrollUptimeFeed
        .connect(l2Messenger)
        .updateStatus(true, timestamp)
      await expect(tx)
        .to.emit(scrollUptimeFeed, 'AnswerUpdated')
        .withArgs(1, 2 /** roundId */, timestamp)
      expect(await scrollUptimeFeed.latestAnswer()).to.equal(1)

      // Submit another status update, different status, same timestamp should update
      tx = await scrollUptimeFeed
        .connect(l2Messenger)
        .updateStatus(false, timestamp)
      await expect(tx)
        .to.emit(scrollUptimeFeed, 'AnswerUpdated')
        .withArgs(0, 3 /** roundId */, timestamp)
      expect(await scrollUptimeFeed.latestAnswer()).to.equal(0)
      expect(await scrollUptimeFeed.latestTimestamp()).to.equal(timestamp)
    })

    it('should ignore out-of-order updates', async () => {
      const timestamp = (await scrollUptimeFeed.latestTimestamp()).add(10_000)
      // Update status
      let tx = await scrollUptimeFeed
        .connect(l2Messenger)
        .updateStatus(true, timestamp)
      await expect(tx)
        .to.emit(scrollUptimeFeed, 'AnswerUpdated')
        .withArgs(1, 2 /** roundId */, timestamp)
      expect(await scrollUptimeFeed.latestAnswer()).to.equal(1)

      // Update with different status, but stale timestamp, should be ignored
      const staleTimestamp = timestamp.sub(1000)
      tx = await scrollUptimeFeed
        .connect(l2Messenger)
        .updateStatus(false, staleTimestamp)
      await expect(tx).to.not.emit(scrollUptimeFeed, 'AnswerUpdated')
      await expect(tx).to.emit(scrollUptimeFeed, 'UpdateIgnored')
    })
  })

  describe('AggregatorV3Interface', () => {
    it('should return valid answer from getRoundData and latestRoundData', async () => {
      let [roundId, answer, startedAt, updatedAt, answeredInRound] =
        await scrollUptimeFeed.latestRoundData()
      expect(roundId).to.equal(1)
      expect(answer).to.equal(0)
      expect(answeredInRound).to.equal(roundId)
      expect(startedAt).to.equal(updatedAt)

      // Submit status update with different status and newer timestamp, should update
      const timestamp = (startedAt as BigNumber).add(1000)
      await scrollUptimeFeed.connect(l2Messenger).updateStatus(true, timestamp)
      ;[roundId, answer, startedAt, updatedAt, answeredInRound] =
        await scrollUptimeFeed.getRoundData(2)
      expect(roundId).to.equal(2)
      expect(answer).to.equal(1)
      expect(answeredInRound).to.equal(roundId)
      expect(startedAt).to.equal(timestamp)
      expect(updatedAt.lte(startedAt)).to.be.true

      // Check that last round is still returning the correct data
      ;[roundId, answer, startedAt, updatedAt, answeredInRound] =
        await scrollUptimeFeed.getRoundData(1)
      expect(roundId).to.equal(1)
      expect(answer).to.equal(0)
      expect(answeredInRound).to.equal(roundId)
      expect(startedAt).to.equal(updatedAt)

      // Assert latestRoundData corresponds to latest round id
      expect(await scrollUptimeFeed.getRoundData(2)).to.deep.equal(
        await scrollUptimeFeed.latestRoundData(),
      )
    })

    it('should revert from #getRoundData when round does not yet exist (future roundId)', async () => {
      expect(scrollUptimeFeed.getRoundData(2)).to.be.revertedWith(
        'NoDataPresent()',
      )
    })

    it('should revert from #getAnswer when round does not yet exist (future roundId)', async () => {
      expect(scrollUptimeFeed.getAnswer(2)).to.be.revertedWith(
        'NoDataPresent()',
      )
    })

    it('should revert from #getTimestamp when round does not yet exist (future roundId)', async () => {
      expect(scrollUptimeFeed.getTimestamp(2)).to.be.revertedWith(
        'NoDataPresent()',
      )
    })
  })

  describe('Protect reads on AggregatorV2V3Interface functions', () => {
    it('should disallow reads on AggregatorV2V3Interface functions when consuming contract is not whitelisted', async () => {
      // Sanity - consumer is not whitelisted
      expect(await scrollUptimeFeed.checkEnabled()).to.be.true
      expect(
        await scrollUptimeFeed.hasAccess(uptimeFeedConsumer.address, '0x00'),
      ).to.be.false

      // Assert reads are not possible from consuming contract
      await expect(uptimeFeedConsumer.latestAnswer()).to.be.revertedWith(
        'No access',
      )
      await expect(uptimeFeedConsumer.latestRoundData()).to.be.revertedWith(
        'No access',
      )
    })

    it('should allow reads on AggregatorV2V3Interface functions when consuming contract is whitelisted', async () => {
      // Whitelist consumer
      await scrollUptimeFeed.addAccess(uptimeFeedConsumer.address)
      // Sanity - consumer is whitelisted
      expect(await scrollUptimeFeed.checkEnabled()).to.be.true
      expect(
        await scrollUptimeFeed.hasAccess(uptimeFeedConsumer.address, '0x00'),
      ).to.be.true

      // Assert reads are possible from consuming contract
      expect(await uptimeFeedConsumer.latestAnswer()).to.be.equal('0')
      const [roundId, answer] = await uptimeFeedConsumer.latestRoundData()
      expect(roundId).to.equal(1)
      expect(answer).to.equal(0)
    })
  })

  describe('Gas costs', () => {
    it('should consume a known amount of gas for updates @skip-coverage', async () => {
      // Sanity - start at flag = 0 (`false`)
      expect(await scrollUptimeFeed.latestAnswer()).to.equal(0)
      let timestamp = await scrollUptimeFeed.latestTimestamp()

      // Gas for no update
      timestamp = timestamp.add(1000)
      const _noUpdateTx = await scrollUptimeFeed
        .connect(l2Messenger)
        .updateStatus(false, timestamp)
      const noUpdateTx = await _noUpdateTx.wait(1)
      // Assert no update
      expect(await scrollUptimeFeed.latestAnswer()).to.equal(0)
      expect(noUpdateTx.cumulativeGasUsed.toNumber()).to.be.closeTo(
        38594,
        gasUsedDeviation,
      )

      // Gas for update
      timestamp = timestamp.add(1000)
      const _updateTx = await scrollUptimeFeed
        .connect(l2Messenger)
        .updateStatus(true, timestamp)
      const updateTx = await _updateTx.wait(1)
      // Assert update
      expect(await scrollUptimeFeed.latestAnswer()).to.equal(1)
      expect(updateTx.cumulativeGasUsed.toNumber()).to.be.closeTo(
        58458,
        gasUsedDeviation,
      )
    })

    describe('Aggregator interface', () => {
      beforeEach(async () => {
        const timestamp = (await scrollUptimeFeed.latestTimestamp()).add(1000)
        // Initialise a round
        await scrollUptimeFeed
          .connect(l2Messenger)
          .updateStatus(true, timestamp)
      })

      it('should consume a known amount of gas for getRoundData(uint80) @skip-coverage', async () => {
        const _tx = await l2Messenger.sendTransaction(
          await scrollUptimeFeed
            .connect(l2Messenger)
            .populateTransaction.getRoundData(1),
        )
        const tx = await _tx.wait(1)
        expect(tx.cumulativeGasUsed.toNumber()).to.be.closeTo(
          30952,
          gasUsedDeviation,
        )
      })

      it('should consume a known amount of gas for latestRoundData() @skip-coverage', async () => {
        const _tx = await l2Messenger.sendTransaction(
          await scrollUptimeFeed
            .connect(l2Messenger)
            .populateTransaction.latestRoundData(),
        )
        const tx = await _tx.wait(1)
        expect(tx.cumulativeGasUsed.toNumber()).to.be.closeTo(
          28523,
          gasUsedDeviation,
        )
      })

      it('should consume a known amount of gas for latestAnswer() @skip-coverage', async () => {
        const _tx = await l2Messenger.sendTransaction(
          await scrollUptimeFeed
            .connect(l2Messenger)
            .populateTransaction.latestAnswer(),
        )
        const tx = await _tx.wait(1)
        expect(tx.cumulativeGasUsed.toNumber()).to.be.closeTo(
          28229,
          gasUsedDeviation,
        )
      })

      it('should consume a known amount of gas for latestTimestamp() @skip-coverage', async () => {
        const _tx = await l2Messenger.sendTransaction(
          await scrollUptimeFeed
            .connect(l2Messenger)
            .populateTransaction.latestTimestamp(),
        )
        const tx = await _tx.wait(1)
        expect(tx.cumulativeGasUsed.toNumber()).to.be.closeTo(
          28129,
          gasUsedDeviation,
        )
      })

      it('should consume a known amount of gas for latestRound() @skip-coverage', async () => {
        const _tx = await l2Messenger.sendTransaction(
          await scrollUptimeFeed
            .connect(l2Messenger)
            .populateTransaction.latestRound(),
        )
        const tx = await _tx.wait(1)
        expect(tx.cumulativeGasUsed.toNumber()).to.be.closeTo(
          28145,
          gasUsedDeviation,
        )
      })

      it('should consume a known amount of gas for getAnswer(roundId) @skip-coverage', async () => {
        const _tx = await l2Messenger.sendTransaction(
          await scrollUptimeFeed
            .connect(l2Messenger)
            .populateTransaction.getAnswer(1),
        )
        const tx = await _tx.wait(1)
        expect(tx.cumulativeGasUsed.toNumber()).to.be.closeTo(
          30682,
          gasUsedDeviation,
        )
      })

      it('should consume a known amount of gas for getTimestamp(roundId) @skip-coverage', async () => {
        const _tx = await l2Messenger.sendTransaction(
          await scrollUptimeFeed
            .connect(l2Messenger)
            .populateTransaction.getTimestamp(1),
        )
        const tx = await _tx.wait(1)
        expect(tx.cumulativeGasUsed.toNumber()).to.be.closeTo(
          30570,
          gasUsedDeviation,
        )
      })
    })
  })
})
