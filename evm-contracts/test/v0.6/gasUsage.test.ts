import {
  contract,
  matchers,
  helpers as h,
  setup,
} from '@chainlink/test-helpers'
import { EACAggregatorProxy__factory } from '../../ethers/v0.6/factories/EACAggregatorProxy__factory'
import { AccessControlledAggregator__factory } from '../../ethers/v0.6/factories/AccessControlledAggregator__factory'
import { FluxAggregatorTestHelper__factory } from '../../ethers/v0.6/factories/FluxAggregatorTestHelper__factory'

let personas: setup.Personas

const provider = setup.provider()
const linkTokenFactory = new contract.LinkToken__factory()
const aggregatorFactory = new AccessControlledAggregator__factory()
const testHelperFactory = new FluxAggregatorTestHelper__factory()
const proxyFactory = new EACAggregatorProxy__factory()
const emptyAddress = '0x0000000000000000000000000000000000000000'
const decimals = 18
const phaseBase = h.bigNum(2).pow(64)

beforeAll(async () => {
  const users = await setup.users(provider)

  personas = users.personas
})

describe('gas usage', () => {
  let aggregator: contract.Instance<AccessControlledAggregator__factory>
  let proxy: contract.Instance<EACAggregatorProxy__factory>
  let testHelper: contract.Instance<FluxAggregatorTestHelper__factory>

  describe('EACAggreagtorProxy => AccessControlledAggreagtor', () => {
    beforeEach(async () => {
      await setup.snapshot(provider, async () => {
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
          .deploy(aggregator.address, emptyAddress)

        await aggregator.changeOracles(
          [],
          [personas.Neil.address],
          [personas.Neil.address],
          1,
          1,
          0,
        )
        await aggregator.connect(personas.Neil).submit(1, 100)
        await aggregator.disableAccessCheck()
        await aggregator.addAccess(proxy.address)
      })()
    })

    it('#latestAnswer', async () => {
      const tx1 = await testHelper.readLatestAnswer(aggregator.address)
      const tx2 = await testHelper.readLatestAnswer(proxy.address)

      matchers.gasDiffLessThan(3000, await tx1.wait(), await tx2.wait())
    })

    it('#latestRound', async () => {
      const tx1 = await testHelper.readLatestRound(aggregator.address)
      const tx2 = await testHelper.readLatestRound(proxy.address)

      matchers.gasDiffLessThan(3000, await tx1.wait(), await tx2.wait())
    })

    it('#latestTimestamp', async () => {
      const tx1 = await testHelper.readLatestTimestamp(aggregator.address)
      const tx2 = await testHelper.readLatestTimestamp(proxy.address)

      matchers.gasDiffLessThan(3000, await tx1.wait(), await tx2.wait())
    })

    it('#getAnswer', async () => {
      const aggId = 1
      const proxyId = phaseBase.add(aggId)
      const tx1 = await testHelper.readGetAnswer(aggregator.address, aggId)
      const tx2 = await testHelper.readGetAnswer(proxy.address, proxyId)

      matchers.gasDiffLessThan(4000, await tx1.wait(), await tx2.wait())
    })

    it('#getTimestamp', async () => {
      const aggId = 1
      const proxyId = phaseBase.add(aggId)
      const tx1 = await testHelper.readGetTimestamp(aggregator.address, aggId)
      const tx2 = await testHelper.readGetTimestamp(proxy.address, proxyId)

      matchers.gasDiffLessThan(4000, await tx1.wait(), await tx2.wait())
    })

    it('#latestRoundData', async () => {
      const tx1 = await testHelper.readLatestRoundData(aggregator.address)
      const tx2 = await testHelper.readLatestRoundData(proxy.address)

      matchers.gasDiffLessThan(3000, await tx1.wait(), await tx2.wait())
    })
  })
})
