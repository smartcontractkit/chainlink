import {
  contract,
  helpers as h,
  matchers,
  oracle,
  setup,
} from '@chainlink/test-helpers'
// import { assert } from 'chai'
// import { ethers, utils } from 'ethers'
import { BasicConsumer__factory } from '../../ethers/v0.6/factories/BasicConsumer__factory'
import { Operator__factory } from '../../ethers/v0.7/factories/Operator__factory'
import { Oracle__factory } from '../../ethers/v0.6/factories/Oracle__factory'

const operatorFactory = new Operator__factory()
const oracleFactory = new Oracle__factory()
const basicConsumerFactory = new BasicConsumer__factory()
const linkTokenFactory = new contract.LinkToken__factory()

let roles: setup.Roles
const provider = setup.provider()

beforeAll(async () => {
  const users = await setup.users(provider)

  roles = users.roles
})

describe('Operator Gas Tests', () => {
  const specId =
    '0x4c7b7ffb66b344fbaa64995af81e355a00000000000000000000000000000000'
  let link: contract.Instance<contract.LinkToken__factory>
  let oracle1: contract.Instance<Oracle__factory>
  let operator1: contract.Instance<Operator__factory>
  let operator2: contract.Instance<Operator__factory>

  const deployment = setup.snapshot(provider, async () => {
    link = await linkTokenFactory.connect(roles.defaultAccount).deploy()

    operator1 = await operatorFactory
      .connect(roles.defaultAccount)
      .deploy(link.address, roles.defaultAccount.address)
    await operator1.setAuthorizedSenders([roles.oracleNode.address])

    operator2 = await operatorFactory
      .connect(roles.defaultAccount)
      .deploy(link.address, roles.defaultAccount.address)
    await operator2.setAuthorizedSenders([roles.oracleNode.address])

    oracle1 = await oracleFactory
      .connect(roles.defaultAccount)
      .deploy(link.address)
    await oracle1.setFulfillmentPermission(roles.oracleNode.address, true)
  })

  beforeEach(async () => {
    await deployment()
  })

  // Test Oracle.fulfillOracleRequest vs Operator.fulfillOracleRequest
  describe('v0.6/Oracle vs v0.7/Operator #fulfillOracleRequest', () => {
    const response = 'Hi Mom!'
    let basicConsumer1: contract.Instance<BasicConsumer__factory>
    let basicConsumer2: contract.Instance<BasicConsumer__factory>

    let request1: ReturnType<typeof oracle.decodeRunRequest>
    let request2: ReturnType<typeof oracle.decodeRunRequest>

    beforeEach(async () => {
      basicConsumer1 = await basicConsumerFactory
        .connect(roles.consumer)
        .deploy(link.address, oracle1.address, specId)
      basicConsumer2 = await basicConsumerFactory
        .connect(roles.consumer)
        .deploy(link.address, operator1.address, specId)

      const paymentAmount = h.toWei('1')
      const currency = 'USD'

      await link.transfer(basicConsumer1.address, paymentAmount)
      const tx1 = await basicConsumer1.requestEthereumPrice(
        currency,
        paymentAmount,
      )
      const receipt1 = await tx1.wait()
      request1 = oracle.decodeRunRequest(receipt1.logs?.[3])

      await link.transfer(basicConsumer2.address, paymentAmount)
      const tx2 = await basicConsumer2.requestEthereumPrice(
        currency,
        paymentAmount,
      )
      const receipt2 = await tx2.wait()
      request2 = oracle.decodeRunRequest(receipt2.logs?.[3])
    })

    it('uses acceptable gas', async () => {
      const tx1 = await oracle1
        .connect(roles.oracleNode)
        .fulfillOracleRequest(...oracle.convertFufillParams(request1, response))
      const tx2 = await operator1
        .connect(roles.oracleNode)
        .fulfillOracleRequest(...oracle.convertFufillParams(request2, response))
      const receipt1 = await tx1.wait()
      const receipt2 = await tx2.wait()
      // 38014 vs 40260
      matchers.gasDiffLessThan(2500, receipt1, receipt2)
    })
  })

  // Test Operator1.fulfillOracleRequest vs Operator2.fulfillOracleRequest2
  // with single word response
  describe('Operator #fulfillOracleRequest vs #fulfillOracleRequest2', () => {
    const response = 'Hi Mom!'
    let basicConsumer1: contract.Instance<BasicConsumer__factory>
    let basicConsumer2: contract.Instance<BasicConsumer__factory>

    let request1: ReturnType<typeof oracle.decodeRunRequest>
    let request2: ReturnType<typeof oracle.decodeRunRequest>

    beforeEach(async () => {
      basicConsumer1 = await basicConsumerFactory
        .connect(roles.consumer)
        .deploy(link.address, operator1.address, specId)
      basicConsumer2 = await basicConsumerFactory
        .connect(roles.consumer)
        .deploy(link.address, operator2.address, specId)

      const paymentAmount = h.toWei('1')
      const currency = 'USD'

      await link.transfer(basicConsumer1.address, paymentAmount)
      const tx1 = await basicConsumer1.requestEthereumPrice(
        currency,
        paymentAmount,
      )
      const receipt1 = await tx1.wait()
      request1 = oracle.decodeRunRequest(receipt1.logs?.[3])

      await link.transfer(basicConsumer2.address, paymentAmount)
      const tx2 = await basicConsumer2.requestEthereumPrice(
        currency,
        paymentAmount,
      )
      const receipt2 = await tx2.wait()
      request2 = oracle.decodeRunRequest(receipt2.logs?.[3])
    })

    it('uses acceptable gas', async () => {
      const tx1 = await operator1
        .connect(roles.oracleNode)
        .fulfillOracleRequest(...oracle.convertFufillParams(request1, response))

      const responseTypes = ['bytes32']
      const responseValues = [h.toBytes32String(response)]
      const tx2 = await operator2
        .connect(roles.oracleNode)
        .fulfillOracleRequest2(
          ...oracle.convertFulfill2Params(
            request2,
            responseTypes,
            responseValues,
          ),
        )

      const receipt1 = await tx1.wait()
      const receipt2 = await tx2.wait()
      // 40260 vs 41423
      matchers.gasDiffLessThan(1200, receipt1, receipt2)
    })
  })
})
