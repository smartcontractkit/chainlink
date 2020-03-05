import { assert } from 'chai'
import fluxMonitorJob from '../fixtures/flux-monitor-job'
import {
  assertEventually,
  getArgs,
  wait,
  createProvider,
  fundAddress,
} from '../test-helpers/common'
import * as clClient from '../test-helpers/chainlink-cli'
import { contract, helpers as h, matchers } from '@chainlink/test-helpers'
import { FluxAggregatorFactory } from '@chainlink/contracts/ethers/v0.6/FluxAggregatorFactory'
import { JobSpec } from '../../../operator_ui/@types/operator_ui'
import 'isomorphic-unfetch'
import { ethers } from 'ethers'

const provider = createProvider()
const carol = ethers.Wallet.createRandom().connect(provider)
const linkTokenFactory = new contract.LinkTokenFactory(carol)
const fluxAggregatorFactory = new FluxAggregatorFactory(carol)
const { EXTERNAL_ADAPTER_URL } = getArgs(['EXTERNAL_ADAPTER_URL'])
const adapterURL = new URL('result', EXTERNAL_ADAPTER_URL).href
const deposit = h.toWei('1000')

async function changePriceFeed(value: number) {
  const response = await fetch(adapterURL, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ result: value }),
  })
  assert(response.ok)
}

describe('flux monitor eth client integration', () => {
  let linkToken: contract.Instance<contract.LinkTokenFactory>
  let fluxAggregator: contract.Instance<FluxAggregatorFactory>
  let job: JobSpec
  let node1Address: string

  beforeAll(async () => {
    await clClient.login()
    node1Address = (await clClient.getAdminInfo())[0].address
    await fundAddress(carol.address)
    await fundAddress(node1Address)
    linkToken = await linkTokenFactory.deploy()
    await linkToken.deployed()
  })

  afterEach(async () => {
    await clClient.archiveJob(job.id)
    await changePriceFeed(100) // original price
  })

  it('updates the price with a single node', async () => {
    fluxAggregator = await fluxAggregatorFactory.deploy(
      linkToken.address,
      1,
      3,
      1,
      ethers.utils.formatBytes32String('ETH/USD'),
    )
    await fluxAggregator.deployed()
    const tx1 = await fluxAggregator.addOracle(
      node1Address,
      node1Address,
      1,
      1,
      0,
    )
    await tx1.wait()
    const tx2 = await linkToken.transfer(fluxAggregator.address, deposit)
    await tx2.wait()
    const tx3 = await fluxAggregator.updateAvailableFunds()
    await tx3.wait()
    // await Promise.all([tx1.wait(), tx2.wait(), tx3.wait()])

    expect(await fluxAggregator.getOracles()).toEqual([node1Address])
    matchers.bigNum(
      await linkToken.balanceOf(fluxAggregator.address),
      deposit,
      'Unable to fund FluxAggregator',
    )

    const initialJobCount = (await clClient.getJobs()).length
    const initialRunCount = (await clClient.getRunResults()).length

    // create FM job
    fluxMonitorJob.initiators[0].params.address = fluxAggregator.address
    job = await clClient.createJob(JSON.stringify(fluxMonitorJob))
    assert.equal((await clClient.getJobs()).length, initialJobCount + 1)

    // Job should trigger initial FM run
    await assertEventually(async () => {
      return (await clClient.getRunResults()).length === initialRunCount + 1
    }, 'initial job never run')

    await assertEventually(async () => {
      return (await fluxAggregator.latestAnswer()).eq(10000)
    }, 'FluxAggregator latest answer not updated in round 0')

    // Nominally change price feed
    await changePriceFeed(101)
    await wait(10000)
    assert.equal(
      (await clClient.getRunResults()).length,
      initialRunCount + 1,
      'Flux Monitor should not run job after nominal price deviation',
    )

    // Significantly change price feed
    await changePriceFeed(110)
    await assertEventually(async () => {
      return (await clClient.getRunResults()).length === initialRunCount + 2
    }, 'second job never run')
    await assertEventually(async () => {
      return (await fluxAggregator.latestAnswer()).eq(11000)
    }, 'FluxAggregator latest answer not updated in round 1')
  })
})
