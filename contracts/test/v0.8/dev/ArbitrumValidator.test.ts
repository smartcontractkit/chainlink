import { ethers } from 'hardhat'
import { BigNumber, BigNumberish, Contract, ContractFactory } from 'ethers'
import { expect } from 'chai'
import { SignerWithAddress } from '@nomiclabs/hardhat-ethers/signers'
import {
  deployMockContract,
  MockContract,
} from '@ethereum-waffle/mock-contract'
/// Pick ABIs from compilation
// @ts-ignore
import { abi as arbitrumSequencerStatusRecorderAbi } from '../../../artifacts/src/v0.8/dev/ArbitrumSequencerUptimeFeed.sol/ArbitrumSequencerUptimeFeed.json'
// @ts-ignore
import { abi as arbitrumInboxAbi } from '../../../artifacts/src/v0.8/dev/vendor/arb-bridge-eth/v0.8.0-custom/contracts/bridge/interfaces/IInbox.sol/IInbox.json'
// @ts-ignore
import { abi as aggregatorAbi } from '../../../artifacts/src/v0.8/interfaces/AggregatorV2V3Interface.sol/AggregatorV2V3Interface.json'

const truncateBigNumToAddress = (num: BigNumberish) => {
  // Pad, then slice off '0x' prefix
  const hexWithoutPrefix = BigNumber.from(num).toHexString().slice(2)
  // Ethereum address is 20B -> 40 hex chars w/o 0x prefix
  const truncated = hexWithoutPrefix
    .split('')
    .reverse()
    .slice(0, 40)
    .reverse()
    .join('')
  return '0x' + truncated
}

describe('ArbitrumValidator', () => {
  const MAX_GAS = BigNumber.from(1_000_000)
  const GAS_PRICE_BID = BigNumber.from(1_000_000)
  const BASE_FEE = BigNumber.from(14_000_000_000)
  /** Fake L2 target */
  const L2_SEQ_STATUS_RECORDER_ADDRESS =
    '0x491B1dDA0A8fa069bbC1125133A975BF4e85a91b'
  let arbitrumValidator: Contract
  let accessController: Contract
  let arbitrumSequencerStatusRecorderFactory: ContractFactory
  let mockArbitrumInbox: Contract
  let l1GasFeed: MockContract
  let deployer: SignerWithAddress
  let eoaValidator: SignerWithAddress
  let arbitrumValidatorL2Address: string
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
    arbitrumSequencerStatusRecorderFactory = await ethers.getContractFactory(
      'src/v0.8/dev/ArbitrumSequencerUptimeFeed.sol:ArbitrumSequencerUptimeFeed',
      deployer,
    )
    l1GasFeed = await deployMockContract(deployer as any, aggregatorAbi)
    await l1GasFeed.mock.latestRoundData.returns(
      '73786976294838220258' /** roundId */,
      '96800000000' /** answer */,
      '163826896' /** startedAt */,
      '1638268960' /** updatedAt */,
      '73786976294838220258' /** answeredInRound */,
    )
    // Arbitrum Inbox contract on L1
    const mockArbitrumInboxFactory = await ethers.getContractFactory(
      'src/v0.8/tests/MockArbitrumInbox.sol:MockArbitrumInbox',
    )
    mockArbitrumInbox = await mockArbitrumInboxFactory.deploy()

    // Contract under test
    const arbitrumValidatorFactory = await ethers.getContractFactory(
      'src/v0.8/dev/ArbitrumValidator.sol:ArbitrumValidator',
      deployer,
    )
    arbitrumValidator = await arbitrumValidatorFactory.deploy(
      mockArbitrumInbox.address,
      L2_SEQ_STATUS_RECORDER_ADDRESS,
      accessController.address,
      MAX_GAS /** L1 gas bid */,
      GAS_PRICE_BID /** L2 gas bid */,
      BASE_FEE,
      l1GasFeed.address,
      0,
    )
    // Transfer some ETH to the ArbitrumValidator contract
    await deployer.sendTransaction({
      to: arbitrumValidator.address,
      value: ethers.utils.parseEther('1.0'),
    })
    arbitrumValidatorL2Address = ethers.utils.getAddress(
      truncateBigNumToAddress(
        BigNumber.from(arbitrumValidator.address).add(
          '0x1111000000000000000000000000000000001111',
        ),
      ),
    )
  })

  describe('#validate', () => {
    it('post sequencer offline', async () => {
      await arbitrumValidator.addAccess(eoaValidator.address)

      const now = Math.ceil(Date.now() / 1000) + 1000
      await ethers.provider.send('evm_setNextBlockTimestamp', [now])
      const arbitrumSequencerStatusRecorderCallData =
        arbitrumSequencerStatusRecorderFactory.interface.encodeFunctionData(
          'updateStatus',
          [true, now],
        )
      await expect(arbitrumValidator.connect(eoaValidator).validate(0, 0, 1, 1))
        .to.emit(
          mockArbitrumInbox,
          'RetryableTicketNoRefundAliasRewriteCreated',
        )
        .withArgs(
          L2_SEQ_STATUS_RECORDER_ADDRESS,
          0,
          '25312000000000',
          arbitrumValidatorL2Address,
          arbitrumValidatorL2Address,
          MAX_GAS,
          GAS_PRICE_BID,
          arbitrumSequencerStatusRecorderCallData,
        )
    })
  })
})
