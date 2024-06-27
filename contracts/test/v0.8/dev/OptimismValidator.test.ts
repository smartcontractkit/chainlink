import { ethers } from 'hardhat'
import { BigNumber, Contract, ContractFactory } from 'ethers'
import { expect } from 'chai'
import { SignerWithAddress } from '@nomiclabs/hardhat-ethers/signers'
/// Pick ABIs from compilation
// @ts-ignore
import { abi as optimismSequencerStatusRecorderAbi } from '../../../artifacts/src/v0.8/dev/OptimismSequencerUptimeFeed.sol/OptimismSequencerUptimeFeed.json'
// @ts-ignore
import { abi as optimismL1CrossDomainMessengerAbi } from '@eth-optimism/contracts/artifacts/contracts/L1/messaging/L1CrossDomainMessenger.sol'
// @ts-ignore
import { abi as aggregatorAbi } from '../../../artifacts/src/v0.8/interfaces/AggregatorV2V3Interface.sol/AggregatorV2V3Interface.json'

describe('OptimismValidator', () => {
  const GAS_LIMIT = BigNumber.from(1_900_000)
  /** Fake L2 target */
  const L2_SEQ_STATUS_RECORDER_ADDRESS =
    '0x491B1dDA0A8fa069bbC1125133A975BF4e85a91b'
  let optimismValidator: Contract
  let optimismUptimeFeedFactory: ContractFactory
  let mockOptimismL1CrossDomainMessenger: Contract
  let deployer: SignerWithAddress
  let eoaValidator: SignerWithAddress

  before(async () => {
    const accounts = await ethers.getSigners()
    deployer = accounts[0]
    eoaValidator = accounts[1]
  })

  beforeEach(async () => {
    // Required for building the calldata
    optimismUptimeFeedFactory = await ethers.getContractFactory(
      'src/v0.8/dev/OptimismSequencerUptimeFeed.sol:OptimismSequencerUptimeFeed',
      deployer,
    )

    // Optimism Messenger contract on L1
    const mockOptimismL1CrossDomainMessengerFactory =
      await ethers.getContractFactory(
        'src/v0.8/tests/MockOptimismL1CrossDomainMessenger.sol:MockOptimismL1CrossDomainMessenger',
      )
    mockOptimismL1CrossDomainMessenger =
      await mockOptimismL1CrossDomainMessengerFactory.deploy()

    // Contract under test
    const optimismValidatorFactory = await ethers.getContractFactory(
      'src/v0.8/dev/OptimismValidator.sol:OptimismValidator',
      deployer,
    )

    optimismValidator = await optimismValidatorFactory.deploy(
      mockOptimismL1CrossDomainMessenger.address,
      L2_SEQ_STATUS_RECORDER_ADDRESS,
      GAS_LIMIT,
    )
  })

  describe('#setGasLimit', () => {
    it('correctly updates the gas limit', async () => {
      const newGasLimit = BigNumber.from(2_000_000)
      const tx = await optimismValidator.setGasLimit(newGasLimit)
      await tx.wait()
      const currentGasLimit = await optimismValidator.getGasLimit()
      expect(currentGasLimit).to.equal(newGasLimit)
    })
  })

  describe('#validate', () => {
    it('reverts if called by account with no access', async () => {
      await expect(
        optimismValidator.connect(eoaValidator).validate(0, 0, 1, 1),
      ).to.be.revertedWith('No access')
    })

    it('posts sequencer status when there is not status change', async () => {
      await optimismValidator.addAccess(eoaValidator.address)

      const currentBlockNum = await ethers.provider.getBlockNumber()
      const currentBlock = await ethers.provider.getBlock(currentBlockNum)
      const futureTimestamp = currentBlock.timestamp + 5000

      await ethers.provider.send('evm_setNextBlockTimestamp', [futureTimestamp])
      const sequencerStatusRecorderCallData =
        optimismUptimeFeedFactory.interface.encodeFunctionData('updateStatus', [
          false,
          futureTimestamp,
        ])

      await expect(optimismValidator.connect(eoaValidator).validate(0, 0, 0, 0))
        .to.emit(mockOptimismL1CrossDomainMessenger, 'SentMessage')
        .withArgs(
          L2_SEQ_STATUS_RECORDER_ADDRESS,
          optimismValidator.address,
          sequencerStatusRecorderCallData,
          0,
          GAS_LIMIT,
        )
    })

    it('post sequencer offline', async () => {
      await optimismValidator.addAccess(eoaValidator.address)

      const currentBlockNum = await ethers.provider.getBlockNumber()
      const currentBlock = await ethers.provider.getBlock(currentBlockNum)
      const futureTimestamp = currentBlock.timestamp + 10000

      await ethers.provider.send('evm_setNextBlockTimestamp', [futureTimestamp])
      const sequencerStatusRecorderCallData =
        optimismUptimeFeedFactory.interface.encodeFunctionData('updateStatus', [
          true,
          futureTimestamp,
        ])

      await expect(optimismValidator.connect(eoaValidator).validate(0, 0, 1, 1))
        .to.emit(mockOptimismL1CrossDomainMessenger, 'SentMessage')
        .withArgs(
          L2_SEQ_STATUS_RECORDER_ADDRESS,
          optimismValidator.address,
          sequencerStatusRecorderCallData,
          0,
          GAS_LIMIT,
        )
    })
  })
})
