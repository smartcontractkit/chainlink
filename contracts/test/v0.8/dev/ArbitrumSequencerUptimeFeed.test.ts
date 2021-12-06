import { ethers, network } from 'hardhat'
import { Contract, BigNumber } from 'ethers'
import { expect } from 'chai'
import { SignerWithAddress } from '@nomiclabs/hardhat-ethers/signers'

const now = () => BigNumber.from(Date.now()).div(1000)

describe('ArbitrumSequencerUptimeFeed', () => {
  let flags: Contract
  let arbitrumSequencerUptimeFeed: Contract
  let accessController: Contract
  let uptimeFeedConsumer: Contract
  let deployer: SignerWithAddress
  let l1Owner: SignerWithAddress
  let l2Messenger: SignerWithAddress
  before(async () => {
    const accounts = await ethers.getSigners()
    deployer = accounts[0]
    l1Owner = accounts[1]
    const dummy = accounts[2]
    const l2MessengerAddress = ethers.utils.getAddress(
      BigNumber.from(l1Owner.address)
        .add('0x1111000000000000000000000000000000001111')
        .toHexString(),
    )
    // Pretend we're on L2
    await network.provider.request({
      method: 'hardhat_impersonateAccount',
      params: [l2MessengerAddress],
    })
    l2Messenger = await ethers.getSigner(l2MessengerAddress)
    // Credit the L2 messenger with some ETH
    await dummy.sendTransaction({
      to: l2Messenger.address,
      value: (await dummy.getBalance()).sub(ethers.utils.parseEther('0.1')),
    })
  })

  beforeEach(async () => {
    const accessControllerFactory = await ethers.getContractFactory(
      'src/v0.8/SimpleWriteAccessController.sol:SimpleWriteAccessController',
      deployer,
    )
    accessController = await accessControllerFactory.deploy()

    const flagsHistoryFactory = await ethers.getContractFactory(
      'src/v0.8/dev/Flags.sol:Flags',
      deployer,
    )
    flags = await flagsHistoryFactory.deploy(
      accessController.address,
      accessController.address,
    )
    await accessController.addAccess(flags.address)

    const arbitrumSequencerStatusRecorderFactory =
      await ethers.getContractFactory(
        'src/v0.8/dev/ArbitrumSequencerUptimeFeed.sol:ArbitrumSequencerUptimeFeed',
        deployer,
      )
    arbitrumSequencerUptimeFeed =
      await arbitrumSequencerStatusRecorderFactory.deploy(
        flags.address,
        l1Owner.address,
      )
    // Required for ArbitrumSequencerUptimeFeed to raise/lower flags
    await accessController.addAccess(arbitrumSequencerUptimeFeed.address)
    // Required for ArbitrumSequencerUptimeFeed to read flags
    await flags.addAccess(arbitrumSequencerUptimeFeed.address)

    // Deployer requires access to invoke initialize
    await accessController.addAccess(deployer.address)
    // Once ArbitrumSequencerUptimeFeed has access, we can initialise the 0th aggregator round
    const initTx = await arbitrumSequencerUptimeFeed
      .connect(deployer)
      .initialize()
    await expect(initTx).to.emit(arbitrumSequencerUptimeFeed, 'Initialized')

    // Mock consumer
    const statusFeedConsumerFactory = await ethers.getContractFactory(
      'src/v0.8/tests/FeedConsumer.sol:FeedConsumer',
      deployer,
    )
    uptimeFeedConsumer = await statusFeedConsumerFactory.deploy(
      arbitrumSequencerUptimeFeed.address,
    )
  })

  describe('#updateStatus', () => {
    it('should only update status when status has changed', async () => {
      let timestamp = now()
      let tx = await arbitrumSequencerUptimeFeed
        .connect(l2Messenger)
        .updateStatus(true, timestamp)
      await expect(tx)
        .to.emit(arbitrumSequencerUptimeFeed, 'AnswerUpdated')
        .withArgs(1, 2 /** roundId */, timestamp)
      expect(await arbitrumSequencerUptimeFeed.latestAnswer()).to.equal(1)

      // Submit another status update, same status, should ignore
      timestamp = now()
      tx = await arbitrumSequencerUptimeFeed
        .connect(l2Messenger)
        .updateStatus(true, timestamp)
      await expect(tx).not.to.emit(arbitrumSequencerUptimeFeed, 'AnswerUpdated')
      expect(await arbitrumSequencerUptimeFeed.latestAnswer()).to.equal('1')
      expect(await arbitrumSequencerUptimeFeed.latestTimestamp()).to.equal(
        timestamp,
      )

      // Submit another status update, different status, should update
      timestamp = now()
      tx = await arbitrumSequencerUptimeFeed
        .connect(l2Messenger)
        .updateStatus(false, timestamp)
      await expect(tx)
        .to.emit(arbitrumSequencerUptimeFeed, 'AnswerUpdated')
        .withArgs(0, 3 /** roundId */, timestamp)
      expect(await arbitrumSequencerUptimeFeed.latestAnswer()).to.equal(0)
      expect(await arbitrumSequencerUptimeFeed.latestTimestamp()).to.equal(
        timestamp,
      )
    })
  })

  describe('AggregatorV3Interface', () => {
    it('should return valid answer from getRoundData and latestRoundData', async () => {
      let [roundId, answer, startedAt, updatedAt, answeredInRound] =
        await arbitrumSequencerUptimeFeed.getRoundData(1)
      expect(roundId).to.equal(1)
      expect(answer).to.equal(0)
      expect(answeredInRound).to.equal(roundId)
      expect(startedAt).to.equal(updatedAt)

      // Submit status update with different status, should update
      const timestamp = now()
      await arbitrumSequencerUptimeFeed
        .connect(l2Messenger)
        .updateStatus(true, timestamp)
      ;[roundId, answer, startedAt, updatedAt, answeredInRound] =
        await arbitrumSequencerUptimeFeed.getRoundData(2)
      expect(roundId).to.equal(2)
      expect(answer).to.equal(1)
      expect(answeredInRound).to.equal(roundId)
      expect(startedAt).to.equal(timestamp)
      expect(updatedAt).to.equal(startedAt)

      // Check that last round is still returning the correct data
      ;[roundId, answer, startedAt, updatedAt, answeredInRound] =
        await arbitrumSequencerUptimeFeed.getRoundData(1)
      expect(roundId).to.equal(1)
      expect(answer).to.equal(0)
      expect(answeredInRound).to.equal(roundId)
      expect(startedAt).to.equal(updatedAt)

      // Assert latestRoundData corresponds to latest round id
      expect(await arbitrumSequencerUptimeFeed.getRoundData(2)).to.deep.equal(
        await arbitrumSequencerUptimeFeed.latestRoundData(),
      )
    })

    it('should return 0 from #getRoundData when round does not yet exist (future roundId)', async () => {
      const [roundId, answer, startedAt, updatedAt, answeredInRound] =
        await arbitrumSequencerUptimeFeed.getRoundData(2)
      expect(roundId).to.equal(2)
      expect(answer).to.equal(0)
      expect(startedAt).to.equal(0)
      expect(updatedAt).to.equal(0)
      expect(answeredInRound).to.equal(2)
    })
  })

  describe('Protect reads on AggregatorV2V3Interface functions', () => {
    it('should disallow reads on AggregatorV2V3Interface functions when consuming contract is not whitelisted', async () => {
      // Sanity - consumer is not whitelisted
      expect(await arbitrumSequencerUptimeFeed.checkEnabled()).to.be.true
      expect(
        await arbitrumSequencerUptimeFeed.hasAccess(
          uptimeFeedConsumer.address,
          '0x00',
        ),
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
      await arbitrumSequencerUptimeFeed.addAccess(uptimeFeedConsumer.address)
      // Sanity - consumer is whitelisted
      expect(await arbitrumSequencerUptimeFeed.checkEnabled()).to.be.true
      expect(
        await arbitrumSequencerUptimeFeed.hasAccess(
          uptimeFeedConsumer.address,
          '0x00',
        ),
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
      expect(await arbitrumSequencerUptimeFeed.latestAnswer()).to.equal(0)

      // Gas for no update
      const _noUpdateTx = await arbitrumSequencerUptimeFeed
        .connect(l2Messenger)
        .updateStatus(false, now())
      const noUpdateTx = await _noUpdateTx.wait(1)
      // Assert no update
      expect(await arbitrumSequencerUptimeFeed.latestAnswer()).to.equal(0)
      expect(noUpdateTx.cumulativeGasUsed).to.equal(26307)

      // Gas for update
      const _updateTx = await arbitrumSequencerUptimeFeed
        .connect(l2Messenger)
        .updateStatus(true, now())
      const updateTx = await _updateTx.wait(1)
      // Assert update
      expect(await arbitrumSequencerUptimeFeed.latestAnswer()).to.equal(1)
      expect(updateTx.cumulativeGasUsed).to.equal(93066)
    })

    it('should consume a known amount of gas for getRoundData(uint80) @skip-coverage', async () => {
      // Initialise a round
      await arbitrumSequencerUptimeFeed
        .connect(l2Messenger)
        .updateStatus(true, now())

      const _tx = await l2Messenger.sendTransaction(
        await arbitrumSequencerUptimeFeed
          .connect(l2Messenger)
          .populateTransaction.getRoundData(1),
      )
      const tx = await _tx.wait(1)
      expect(tx.cumulativeGasUsed).to.equal(31135)
    })

    it('should consume a known amount of gas for latestRoundData() @skip-coverage', async () => {
      // Initialise a round
      await arbitrumSequencerUptimeFeed
        .connect(l2Messenger)
        .updateStatus(true, now())

      const _tx = await l2Messenger.sendTransaction(
        await arbitrumSequencerUptimeFeed
          .connect(l2Messenger)
          .populateTransaction.latestRoundData(),
      )
      const tx = await _tx.wait(1)
      expect(tx.cumulativeGasUsed).to.equal(28523)
    })

    it('should consume a known amount of gas for latestAnswer() @skip-coverage', async () => {
      // Initialise a round
      await arbitrumSequencerUptimeFeed
        .connect(l2Messenger)
        .updateStatus(true, now())

      const _tx = await l2Messenger.sendTransaction(
        await arbitrumSequencerUptimeFeed
          .connect(l2Messenger)
          .populateTransaction.latestAnswer(),
      )
      const tx = await _tx.wait(1)
      expect(tx.cumulativeGasUsed).to.equal(28329)
    })

    it('should consume a known amount of gas for latestTimestamp() @skip-coverage', async () => {
      // Initialise a round
      await arbitrumSequencerUptimeFeed
        .connect(l2Messenger)
        .updateStatus(true, now())

      const _tx = await l2Messenger.sendTransaction(
        await arbitrumSequencerUptimeFeed
          .connect(l2Messenger)
          .populateTransaction.latestTimestamp(),
      )
      const tx = await _tx.wait(1)
      expect(tx.cumulativeGasUsed).to.equal(28229)
    })

    it('should consume a known amount of gas for latestRound() @skip-coverage', async () => {
      // Initialise a round
      await arbitrumSequencerUptimeFeed
        .connect(l2Messenger)
        .updateStatus(true, now())

      const _tx = await l2Messenger.sendTransaction(
        await arbitrumSequencerUptimeFeed
          .connect(l2Messenger)
          .populateTransaction.latestRound(),
      )
      const tx = await _tx.wait(1)
      expect(tx.cumulativeGasUsed).to.equal(28245)
    })

    it('should consume a known amount of gas for getAnswer(roundId) @skip-coverage', async () => {
      // Initialise a round
      await arbitrumSequencerUptimeFeed
        .connect(l2Messenger)
        .updateStatus(true, now())

      const _tx = await l2Messenger.sendTransaction(
        await arbitrumSequencerUptimeFeed
          .connect(l2Messenger)
          .populateTransaction.getAnswer(1),
      )
      const tx = await _tx.wait(1)
      expect(tx.cumulativeGasUsed).to.equal(30865)
    })

    it('should consume a known amount of gas for getTimestamp(roundId) @skip-coverage', async () => {
      // Initialise a round
      await arbitrumSequencerUptimeFeed
        .connect(l2Messenger)
        .updateStatus(true, now())

      const _tx = await l2Messenger.sendTransaction(
        await arbitrumSequencerUptimeFeed
          .connect(l2Messenger)
          .populateTransaction.getTimestamp(1),
      )
      const tx = await _tx.wait(1)
      expect(tx.cumulativeGasUsed).to.equal(30731)
    })
  })
})
