import { contract, helpers as h, setup } from '@chainlink/test-helpers'
import { assert } from 'chai'
import { EACAggregatorProxyFactory } from '../../ethers/v0.6/EACAggregatorProxyFactory'
import { AccessControlledAggregatorFactory } from '../../ethers/v0.6/AccessControlledAggregatorFactory'
import { SimpleReadAccessControllerFactory } from '../../ethers/v0.6/SimpleReadAccessControllerFactory'
import { FluxAggregatorTestHelperFactory } from '../../ethers/v0.6/FluxAggregatorTestHelperFactory'

let personas: setup.Personas

const provider = setup.provider()
const linkTokenFactory = new contract.LinkTokenFactory()
const accessControlFactory = new SimpleReadAccessControllerFactory()
const aggregatorFactory = new AccessControlledAggregatorFactory()
const testHelperFactory = new FluxAggregatorTestHelperFactory()
const proxyFactory = new EACAggregatorProxyFactory()
const emptyAddress = '0x0000000000000000000000000000000000000000'
const decimals = 18
const phaseBase = h.bigNum(2).pow(64)

beforeAll(async () => {
  const users = await setup.users(provider)

  personas = users.personas
})

function gasDiffLessThan(max: number, receipt1: any, receipt2: any) {
  assert(receipt1, 'receipt1 is not present for gas comparison')
  assert(receipt2, 'receipt2 is not present for gas comparison')
  const diff = receipt2.gasUsed?.sub(receipt1.gasUsed || 0)
  assert.isAbove(max, diff?.toNumber() || Infinity)
}

describe('gas usage', () => {
  let controller: contract.Instance<SimpleReadAccessControllerFactory>
  let aggregator: contract.Instance<AccessControlledAggregatorFactory>
  let proxy: contract.Instance<EACAggregatorProxyFactory>
  let testHelper: contract.Instance<FluxAggregatorTestHelperFactory>

  describe('EACAggreagtorProxy => AccessControlledAggreagtor', () => {
    beforeEach(async () => {
      await setup.snapshot(provider, async () => {
        controller = await accessControlFactory
          .connect(personas.Default)
          .deploy()
        testHelper = await testHelperFactory.connect(personas.Default).deploy()
        const link = await linkTokenFactory.connect(personas.Default).deploy()
        aggregator = await (aggregatorFactory as any)
          .connect(personas.Default)
          .deploy(
            link.address,
            0,
            0,
            emptyAddress,
            0,
            h.bigNum(2).pow(254),
            decimals,
            h.toBytes32String('TEST/LINK'),
            { gasLimit: 8_000_000 },
          )
        proxy = await proxyFactory
          .connect(personas.Default)
          .deploy(aggregator.address, controller.address)

        await aggregator.changeOracles(
          [],
          [personas.Neil.address],
          [personas.Neil.address],
          1,
          1,
          0,
        )
        await aggregator.connect(personas.Neil).submit(1, 100)

        await proxy.connect(personas.Default).setController(emptyAddress)
        await aggregator.disableAccessCheck()
        await aggregator.addAccess(proxy.address)
      })()
    })

    it('#latestAnswer', async () => {
      const tx1 = await testHelper.readLatestAnswer(aggregator.address)
      const tx2 = await testHelper.readLatestAnswer(proxy.address)

      gasDiffLessThan(3000, await tx1.wait(), await tx2.wait())
    })

    it('#latestRound', async () => {
      const tx1 = await testHelper.readLatestRound(aggregator.address)
      const tx2 = await testHelper.readLatestRound(proxy.address)

      gasDiffLessThan(3000, await tx1.wait(), await tx2.wait())
    })

    it('#latestTimestamp', async () => {
      const tx1 = await testHelper.readLatestTimestamp(aggregator.address)
      const tx2 = await testHelper.readLatestTimestamp(proxy.address)

      gasDiffLessThan(3000, await tx1.wait(), await tx2.wait())
    })

    it('#getAnswer', async () => {
      const aggId = 1
      const proxyId = phaseBase.add(aggId)
      const tx1 = await testHelper.readGetAnswer(aggregator.address, aggId)
      const tx2 = await testHelper.readGetAnswer(proxy.address, proxyId)

      gasDiffLessThan(4000, await tx1.wait(), await tx2.wait())
    })

    it('#getTimestamp', async () => {
      const aggId = 1
      const proxyId = phaseBase.add(aggId)
      const tx1 = await testHelper.readGetTimestamp(aggregator.address, aggId)
      const tx2 = await testHelper.readGetTimestamp(proxy.address, proxyId)

      gasDiffLessThan(4000, await tx1.wait(), await tx2.wait())
    })

    it('#latestRoundData', async () => {
      const tx1 = await testHelper.readLatestRoundData(aggregator.address)
      const tx2 = await testHelper.readLatestRoundData(proxy.address)

      gasDiffLessThan(3000, await tx1.wait(), await tx2.wait())
    })
  })
})
