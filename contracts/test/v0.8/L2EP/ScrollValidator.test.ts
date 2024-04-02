import { ethers } from 'hardhat'
import { BigNumber, Contract, ContractFactory } from 'ethers'
import { expect } from 'chai'
import { SignerWithAddress } from '@nomiclabs/hardhat-ethers/signers'

describe('ScrollValidator', () => {
  const GAS_LIMIT = BigNumber.from(1_900_000)
  /** Fake L2 target */
  const L2_SEQ_STATUS_RECORDER_ADDRESS =
    '0x491B1dDA0A8fa069bbC1125133A975BF4e85a91b'
  let scrollValidator: Contract
  let l1MessageQueue: Contract
  let scrollUptimeFeedFactory: ContractFactory
  let mockScrollL1CrossDomainMessenger: Contract
  let deployer: SignerWithAddress
  let eoaValidator: SignerWithAddress

  before(async () => {
    const accounts = await ethers.getSigners()
    deployer = accounts[0]
    eoaValidator = accounts[1]
  })

  beforeEach(async () => {
    // Required for building the calldata
    scrollUptimeFeedFactory = await ethers.getContractFactory(
      'src/v0.8/l2ep/dev/scroll/ScrollSequencerUptimeFeed.sol:ScrollSequencerUptimeFeed',
      deployer,
    )

    // Scroll Messenger contract on L1
    const mockScrollL1CrossDomainMessengerFactory =
      await ethers.getContractFactory(
        'src/v0.8/l2ep/test/mocks/scroll/MockScrollL1CrossDomainMessenger.sol:MockScrollL1CrossDomainMessenger',
      )
    mockScrollL1CrossDomainMessenger =
      await mockScrollL1CrossDomainMessengerFactory.deploy()

    // Scroll Message Queue contract on L1
    const l1MessageQueueFactory = await ethers.getContractFactory(
      'src/v0.8/l2ep/test/mocks/scroll/MockScrollL1MessageQueue.sol:MockScrollL1MessageQueue',
      deployer,
    )
    l1MessageQueue = await l1MessageQueueFactory.deploy()

    // Contract under test
    const scrollValidatorFactory = await ethers.getContractFactory(
      'src/v0.8/l2ep/dev/scroll/ScrollValidator.sol:ScrollValidator',
      deployer,
    )

    scrollValidator = await scrollValidatorFactory.deploy(
      mockScrollL1CrossDomainMessenger.address,
      L2_SEQ_STATUS_RECORDER_ADDRESS,
      l1MessageQueue.address,
      GAS_LIMIT,
    )
  })

  describe('#setGasLimit', () => {
    it('correctly updates the gas limit', async () => {
      const newGasLimit = BigNumber.from(2_000_000)
      const tx = await scrollValidator.setGasLimit(newGasLimit)
      await tx.wait()
      const currentGasLimit = await scrollValidator.getGasLimit()
      expect(currentGasLimit).to.equal(newGasLimit)
    })
  })

  describe('#validate', () => {
    it('reverts if called by account with no access', async () => {
      await expect(
        scrollValidator.connect(eoaValidator).validate(0, 0, 1, 1),
      ).to.be.revertedWith('No access')
    })

    it('posts sequencer status when there is not status change', async () => {
      await scrollValidator.addAccess(eoaValidator.address)

      const currentBlock = await ethers.provider.getBlock('latest')
      const futureTimestamp = currentBlock.timestamp + 5000

      await ethers.provider.send('evm_setNextBlockTimestamp', [futureTimestamp])
      const sequencerStatusRecorderCallData =
        scrollUptimeFeedFactory.interface.encodeFunctionData('updateStatus', [
          false,
          futureTimestamp,
        ])

      await expect(scrollValidator.connect(eoaValidator).validate(0, 0, 0, 0))
        .to.emit(mockScrollL1CrossDomainMessenger, 'SentMessage')
        .withArgs(
          scrollValidator.address, // sender
          L2_SEQ_STATUS_RECORDER_ADDRESS, // target
          0, // value
          0, // nonce
          GAS_LIMIT, // gas limit
          sequencerStatusRecorderCallData, // message
        )
    })

    it('post sequencer offline', async () => {
      await scrollValidator.addAccess(eoaValidator.address)

      const currentBlock = await ethers.provider.getBlock('latest')
      const futureTimestamp = currentBlock.timestamp + 10000

      await ethers.provider.send('evm_setNextBlockTimestamp', [futureTimestamp])
      const sequencerStatusRecorderCallData =
        scrollUptimeFeedFactory.interface.encodeFunctionData('updateStatus', [
          true,
          futureTimestamp,
        ])

      await expect(scrollValidator.connect(eoaValidator).validate(0, 0, 1, 1))
        .to.emit(mockScrollL1CrossDomainMessenger, 'SentMessage')
        .withArgs(
          scrollValidator.address, // sender
          L2_SEQ_STATUS_RECORDER_ADDRESS, // target
          0, // value
          0, // nonce
          GAS_LIMIT, // gas limit
          sequencerStatusRecorderCallData, // message
        )
    })
  })
})
