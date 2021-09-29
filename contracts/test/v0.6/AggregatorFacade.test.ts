import { ethers } from 'hardhat'
import { numToBytes32, publicAbi } from '../test-helpers/helpers'
import { assert } from 'chai'
import { Contract, ContractFactory, Signer } from 'ethers'
import { getUsers } from '../test-helpers/setup'
import { convertFufillParams, decodeRunRequest } from '../test-helpers/oracle'
import { bigNumEquals, evmRevert } from '../test-helpers/matchers'

let defaultAccount: Signer

let linkTokenFactory: ContractFactory
let aggregatorFactory: ContractFactory
let oracleFactory: ContractFactory
let aggregatorFacadeFactory: ContractFactory

before(async () => {
  const users = await getUsers()

  defaultAccount = users.roles.defaultAccount
  linkTokenFactory = await ethers.getContractFactory(
    'src/v0.4/LinkToken.sol:LinkToken',
    defaultAccount,
  )
  aggregatorFactory = await ethers.getContractFactory(
    'src/v0.4/Aggregator.sol:Aggregator',
    defaultAccount,
  )
  oracleFactory = await ethers.getContractFactory(
    'src/v0.6/Oracle.sol:Oracle',
    defaultAccount,
  )
  aggregatorFacadeFactory = await ethers.getContractFactory(
    'src/v0.6/AggregatorFacade.sol:AggregatorFacade',
    defaultAccount,
  )
})

describe('AggregatorFacade', () => {
  const jobId1 =
    '0x4c7b7ffb66b344fbaa64995af81e355a00000000000000000000000000000001'
  const previousResponse = numToBytes32(54321)
  const response = numToBytes32(67890)
  const decimals = 18
  const description = 'LINK / USD: Historic Aggregator Facade'

  let link: Contract
  let aggregator: Contract
  let oc1: Contract
  let facade: Contract

  beforeEach(async () => {
    link = await linkTokenFactory.connect(defaultAccount).deploy()
    oc1 = await oracleFactory.connect(defaultAccount).deploy(link.address)
    aggregator = await aggregatorFactory
      .connect(defaultAccount)
      .deploy(link.address, 0, 1, [oc1.address], [jobId1])
    facade = await aggregatorFacadeFactory
      .connect(defaultAccount)
      .deploy(aggregator.address, decimals, description)

    let requestTx = await aggregator.requestRateUpdate()
    let receipt = await requestTx.wait()
    let request = decodeRunRequest(receipt.logs?.[3])
    await oc1.fulfillOracleRequest(
      ...convertFufillParams(request, previousResponse),
    )
    requestTx = await aggregator.requestRateUpdate()
    receipt = await requestTx.wait()
    request = decodeRunRequest(receipt.logs?.[3])
    await oc1.fulfillOracleRequest(...convertFufillParams(request, response))
  })

  it('has a limited public interface [ @skip-coverage ]', () => {
    publicAbi(facade, [
      'aggregator',
      'decimals',
      'description',
      'getAnswer',
      'getRoundData',
      'getTimestamp',
      'latestAnswer',
      'latestRound',
      'latestRoundData',
      'latestTimestamp',
      'version',
    ])
  })

  describe('#constructor', () => {
    it('uses the decimals set in the constructor', async () => {
      bigNumEquals(decimals, await facade.decimals())
    })

    it('uses the description set in the constructor', async () => {
      assert.equal(description, await facade.description())
    })

    it('sets the version to 2', async () => {
      bigNumEquals(2, await facade.version())
    })
  })

  describe('#getAnswer/latestAnswer', () => {
    it('pulls the rate from the aggregator', async () => {
      bigNumEquals(response, await facade.latestAnswer())
      const latestRound = await facade.latestRound()
      bigNumEquals(response, await facade.getAnswer(latestRound))
    })
  })

  describe('#getTimestamp/latestTimestamp', () => {
    it('pulls the timestamp from the aggregator', async () => {
      const height = await aggregator.latestTimestamp()
      assert.notEqual('0', height.toString())
      bigNumEquals(height, await facade.latestTimestamp())
      const latestRound = await facade.latestRound()
      bigNumEquals(
        await aggregator.latestTimestamp(),
        await facade.getTimestamp(latestRound),
      )
    })
  })

  describe('#getRoundData', () => {
    it('assembles the requested round data', async () => {
      const previousId = (await facade.latestRound()).sub(1)
      const round = await facade.getRoundData(previousId)
      bigNumEquals(previousId, round.roundId)
      bigNumEquals(previousResponse, round.answer)
      bigNumEquals(await facade.getTimestamp(previousId), round.startedAt)
      bigNumEquals(await facade.getTimestamp(previousId), round.updatedAt)
      bigNumEquals(previousId, round.answeredInRound)
    })

    it('returns zero data for non-existing rounds', async () => {
      const roundId = 13371337
      await evmRevert(facade.getRoundData(roundId), 'No data present')
    })
  })

  describe('#latestRoundData', () => {
    it('assembles the requested round data', async () => {
      const latestId = await facade.latestRound()
      const round = await facade.latestRoundData()
      bigNumEquals(latestId, round.roundId)
      bigNumEquals(response, round.answer)
      bigNumEquals(await facade.getTimestamp(latestId), round.startedAt)
      bigNumEquals(await facade.getTimestamp(latestId), round.updatedAt)
      bigNumEquals(latestId, round.answeredInRound)
    })

    describe('when there is no latest round', () => {
      beforeEach(async () => {
        aggregator = await aggregatorFactory
          .connect(defaultAccount)
          .deploy(link.address, 0, 1, [oc1.address], [jobId1])
        facade = await aggregatorFacadeFactory
          .connect(defaultAccount)
          .deploy(aggregator.address, decimals, description)
      })

      it('assembles the requested round data', async () => {
        await evmRevert(facade.latestRoundData(), 'No data present')
      })
    })
  })
})
