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
  const L2_TOKEN_BRIDGE_ADDR = '0x4200000000000000000000000000000000000010'
  const L2_ETH_TOKEN_ADDR = '0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000'
  let optimismValidator: Contract
  let accessController: Contract
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
    const accessControllerFactory = await ethers.getContractFactory(
      'src/v0.8/SimpleWriteAccessController.sol:SimpleWriteAccessController',
      deployer,
    )
    accessController = await accessControllerFactory.deploy()

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
      accessController.address,
      L2_TOKEN_BRIDGE_ADDR,
      L2_ETH_TOKEN_ADDR,
      GAS_LIMIT,
    )
    // Transfer some ETH to the OptimismValidator contract
    await deployer.sendTransaction({
      to: optimismValidator.address,
      value: ethers.utils.parseEther('1.0'),
    })
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

  describe('#withdrawFunds', () => {
    it('reverts if called by non owner', async () => {
      const refundAddr = eoaValidator.address
      await expect(
        optimismValidator.connect(eoaValidator).withdrawFundsTo(refundAddr),
      ).to.be.revertedWith('Only callable by owner')
    })

    it('successfully withdraws funds', async () => {
      const refundAddr = deployer.address
      const priorBalance = await ethers.provider.getBalance(refundAddr)
      await optimismValidator.connect(deployer).withdrawFunds()
      const currentBalance = await ethers.provider.getBalance(refundAddr)
      expect(currentBalance.gte(priorBalance)).to.equal(true)
    })
  })

  describe('#withdrawFundsTo', () => {
    it('reverts if called by non owner', async () => {
      const refundAddr = eoaValidator.address
      await expect(
        optimismValidator.connect(eoaValidator).withdrawFundsTo(refundAddr),
      ).to.be.revertedWith('Only callable by owner')
    })

    it('successfully withdraws funds', async () => {
      const refundAddr = deployer.address
      const priorBalance = await ethers.provider.getBalance(refundAddr)
      await optimismValidator.connect(deployer).withdrawFundsTo(refundAddr)
      const currentBalance = await ethers.provider.getBalance(refundAddr)
      expect(currentBalance.gte(priorBalance)).to.equal(true)
    })
  })

  describe('#withdrawFundsFromL2', () => {
    const amountToWithdraw = BigNumber.from(10).pow(18)

    it('reverts if called by non owner', async () => {
      const refundAddr = deployer.address
      await expect(
        optimismValidator
          .connect(eoaValidator)
          .withdrawFundsFromL2(amountToWithdraw, refundAddr),
      ).to.be.revertedWith('Only callable by owner')
    })

    it('successfully withdraws funds from L2', async () => {
      const refundAddr = deployer.address
      await expect(
        optimismValidator
          .connect(deployer)
          .withdrawFundsFromL2(amountToWithdraw, refundAddr),
      )
        .to.emit(optimismValidator, 'L2WithdrawalRequested')
        .withArgs(amountToWithdraw, refundAddr)
    })
  })

  describe('#validate', () => {
    it('does not update the status when there is no status change', async () => {
      await optimismValidator.addAccess(eoaValidator.address)
      await expect(
        optimismValidator.connect(eoaValidator).validate(0, 0, 0, 0),
      ).not.to.emit(mockOptimismL1CrossDomainMessenger, 'SentMessage')
    })

    it('reverts if called by account with no access', async () => {
      await expect(
        optimismValidator.connect(eoaValidator).validate(0, 0, 1, 1),
      ).to.be.revertedWith('No access')
    })

    it('post sequencer offline', async () => {
      await optimismValidator.addAccess(eoaValidator.address)

      const now = Math.ceil(Date.now() / 1000) + 1000
      await ethers.provider.send('evm_setNextBlockTimestamp', [now])
      const arbitrumSequencerStatusRecorderCallData =
        optimismUptimeFeedFactory.interface.encodeFunctionData('updateStatus', [
          true,
          now,
        ])

      await expect(optimismValidator.connect(eoaValidator).validate(0, 0, 1, 1))
        .to.emit(mockOptimismL1CrossDomainMessenger, 'SentMessage')
        .withArgs(
          L2_SEQ_STATUS_RECORDER_ADDRESS,
          optimismValidator.address,
          arbitrumSequencerStatusRecorderCallData,
          0,
          GAS_LIMIT,
        )
    })
  })
})
