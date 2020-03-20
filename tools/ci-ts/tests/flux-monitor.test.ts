import { assert } from 'chai'
import fluxMonitorJob from '../fixtures/flux-monitor-job'
import {
  assertAsync,
  getArgs,
  wait,
  createProvider,
  fundAddress,
} from '../test-helpers/common'
import * as clClient from '../test-helpers/chainlink-cli'
import { contract, helpers as h, matchers } from '@chainlink/test-helpers'
import { FluxAggregatorFactory } from '../../../evm-contracts/ethers/v0.6/FluxAggregatorFactory'
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

async function assertJobRun(
  jobId: string,
  count: number,
  errorMessage: string,
) {
  await assertAsync(() => {
    const jobRuns = clClient.getJobRuns()
    const jobRun = jobRuns[jobRuns.length - 1]
    return (
      clClient.getJobRuns().length === count &&
      jobRun.status === 'completed' &&
      jobRun.jobId === jobId
    )
  }, errorMessage)
}

describe('flux monitor eth client integration', () => {
  let linkToken: contract.Instance<contract.LinkTokenFactory>
  let fluxAggregator: contract.Instance<FluxAggregatorFactory>
  let job: JobSpec
  let node1Address: string

  beforeAll(async () => {
    clClient.login()
    node1Address = clClient.getAdminInfo()[0].address
    await fundAddress(carol.address, 5)
    await fundAddress(node1Address, 5)
    linkToken = await linkTokenFactory.deploy()
    await linkToken.deployed()
  })

  afterEach(async () => {
    clClient.archiveJob(job.id)
    await changePriceFeed(100) // original price
  })

  it('updates the price with a single node', async () => {
    fluxAggregator = await fluxAggregatorFactory.deploy(
      linkToken.address,
      1,
      600,
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

    expect(await fluxAggregator.getOracles()).toEqual([node1Address])
    matchers.bigNum(
      await linkToken.balanceOf(fluxAggregator.address),
      deposit,
      'Unable to fund FluxAggregator',
    )

    const initialJobCount = clClient.getJobs().length
    const initialRunCount = clClient.getJobRuns().length

    // create FM job
    fluxMonitorJob.initiators[0].params.address = fluxAggregator.address
    job = clClient.createJob(JSON.stringify(fluxMonitorJob))
    assert.equal(clClient.getJobs().length, initialJobCount + 1)

    // Job should trigger initial FM run
    await assertJobRun(job.id, initialRunCount + 1, 'initial job never run')
    matchers.bigNum(10000, await fluxAggregator.latestAnswer())

    // Nominally change price feed
    await changePriceFeed(101)
    await wait(10000)
    assert.equal(
      clClient.getJobRuns().length,
      initialRunCount + 1,
      'Flux Monitor should not run job after nominal price deviation',
    )

    // Significantly change price feed
    await changePriceFeed(110)
    await assertJobRun(job.id, initialRunCount + 2, 'second job never run')
    matchers.bigNum(11000, await fluxAggregator.latestAnswer())
  })
})
