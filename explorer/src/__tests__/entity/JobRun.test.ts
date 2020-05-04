import { getRepository } from 'typeorm'
import { createChainlinkNode } from '../../entity/ChainlinkNode'
import { fromString, JobRun, saveJobRunTree } from '../../entity/JobRun'
import ethtxFixture from '../fixtures/JobRun.ethtx.fixture.json'
import fixture from '../fixtures/JobRun.fixture.json'

describe('entity/jobRun/fromString', () => {
  it('successfully creates a run and tasks from json', async () => {
    const jr = fromString(JSON.stringify(fixture))
    expect(jr.id).toBeUndefined()
    expect(jr.runId).toEqual('f1xtureAaaaaaaaaaaaaaaaaaaaaaaaa')
    expect(jr.jobId).toEqual('aeb2861d306645b1ba012079aeb2e53a')
    expect(jr.createdAt).toEqual(new Date('2019-04-01T22:07:04Z'))
    expect(jr.status).toEqual('in_progress')
    expect(jr.finishedAt).toBeNull()

    expect(jr.type).toEqual('runlog')
    expect(jr.requestId).toEqual('RequestID')
    expect(jr.txHash).toEqual(
      '0x00000000000000000000000000000000000000000000000000000000deadbeef',
    )
    expect(jr.requester).toEqual('0x9FBDa871d559710256a2502A2517b794B482Db40')

    expect(jr.taskRuns.length).toEqual(1)
    expect(jr.taskRuns[0].id).toBeUndefined()
    expect(jr.taskRuns[0].index).toEqual(0)
    expect(jr.taskRuns[0].type).toEqual('httpget')
    expect(jr.taskRuns[0].status).toEqual('')
    expect(jr.taskRuns[0].confirmations).toEqual('0')
    expect(jr.taskRuns[0].minimumConfirmations).toEqual('3')
    expect(jr.taskRuns[0].error).toEqual(null)

    const [chainlinkNode] = await createChainlinkNode(
      'job-run-fromString-chainlink-node',
    )
    jr.chainlinkNodeId = chainlinkNode.id
    const r = await getRepository(JobRun).save(jr)
    expect(r.id).toBeDefined()
    expect(r.type).toEqual('runlog')
    expect(r.taskRuns.length).toEqual(1)
  })

  it('successfully creates an ethtx tasks with transaction info', async () => {
    const jr = fromString(JSON.stringify(ethtxFixture))

    expect(jr.taskRuns.length).toEqual(4)
    const ethtxTask = jr.taskRuns[3]
    expect(ethtxTask.id).toBeUndefined()
    expect(ethtxTask.index).toEqual(3)
    expect(ethtxTask.type).toEqual('ethtx')
    expect(ethtxTask.status).toEqual('completed')
    expect(ethtxTask.error).toEqual(null)
    expect(ethtxTask.transactionHash).toEqual(
      '0x1111111111111111111111111111111111111111111111111111111111111111',
    )
    expect(ethtxTask.transactionStatus).toEqual('fulfilledRunLog')
  })

  it('creates when finishedAt is null', () => {
    const fixtureWithoutFinishedAt = Object.assign({}, fixture, {
      finishedAt: null,
    })
    const jr = fromString(JSON.stringify(fixtureWithoutFinishedAt))
    expect(jr.runId).toEqual('f1xtureAaaaaaaaaaaaaaaaaaaaaaaaa')
    expect(jr.finishedAt).toEqual(null)
  })

  it('errors on a malformed string', async () => {
    try {
      fromString('{"absolute":garbage')
    } catch (err) {
      expect(err).toBeDefined()
    }
  })
})

describe('entity/jobRun/saveJobRunTree', () => {
  it('updates jobRun error', async () => {
    const [chainlinkNode] = await createChainlinkNode(
      'testOverwriteJobRunsErrorOnConflict',
    )

    const jr = fromString(JSON.stringify(fixture))
    jr.chainlinkNodeId = chainlinkNode.id
    await saveJobRunTree(jr)

    const initial = await getRepository(JobRun).findOne()
    expect(initial.status).toEqual('in_progress')
    expect(initial.finishedAt).toBeNull()

    jr.status = 'errored'
    jr.error = 'something bad happened'
    jr.finishedAt = new Date('2018-04-01T22:07:04Z')
    await saveJobRunTree(jr)

    const actual = await getRepository(JobRun).findOne()
    expect(actual.status).toEqual(jr.status)
    expect(actual.finishedAt).toEqual(jr.finishedAt)
    expect(actual.error).toEqual(jr.error)
  })

  it('overwrites taskRun values on conflict', async () => {
    const [chainlinkNode] = await createChainlinkNode(
      'testOverwriteTaskRunsOnConflict',
    )

    const jr = fromString(JSON.stringify(fixture))
    jr.chainlinkNodeId = chainlinkNode.id
    await saveJobRunTree(jr)

    const modifications = {
      confirmations: '2',
      error: 'something bad happened',
      minimumConfirmations: '3',
      status: 'errored',
      transactionHash:
        '0x2222222222222222222222222222222222222222222222222222222222222222',
      transactionStatus: 'fulfilledRunLog',
    }

    const initial = await getRepository(JobRun).findOne()
    const initialTask = initial.taskRuns[0]
    expect(initialTask).not.toMatchObject(modifications)

    Object.assign(jr.taskRuns[0], modifications)
    await saveJobRunTree(jr)

    const actual = await getRepository(JobRun).findOne()
    const actualTask = actual.taskRuns[0]
    expect(actualTask).toMatchObject(modifications)
  })
})
