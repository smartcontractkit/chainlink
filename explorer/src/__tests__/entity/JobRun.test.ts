import { Connection } from 'typeorm'
import { closeDbConnection, getDb } from '../../database'
import { createChainlinkNode } from '../../entity/ChainlinkNode'
import { fromString, JobRun, saveJobRunTree } from '../../entity/JobRun'
import ethtxFixture from '../fixtures/JobRun.ethtx.fixture.json'
import fixture from '../fixtures/JobRun.fixture.json'
import updateFixture from '../fixtures/JobRunUpdate.fixture.json'

let db: Connection

beforeAll(async () => {
  db = await getDb()
})

afterAll(async () => closeDbConnection())

describe('entity/jobRun/fromString', () => {
  it('successfully creates a run and tasks from json', async () => {
    const jr = fromString(JSON.stringify(fixture))
    expect(jr.id).toBeUndefined()
    expect(jr.runId).toEqual('f1xtureAaaaaaaaaaaaaaaaaaaaaaaaa')
    expect(jr.jobId).toEqual('aeb2861d306645b1ba012079aeb2e53a')
    expect(jr.createdAt).toEqual(new Date('2019-04-01T22:07:04Z'))
    expect(jr.status).toEqual('in_progress')
    expect(jr.finishedAt).toEqual(new Date('2018-04-01T22:07:04Z'))

    expect(jr.type).toEqual('runlog')
    expect(jr.requestId).toEqual('RequestID')
    expect(jr.txHash).toEqual(
      '0x00000000000000000000000000000000000000000000000000000000deadbeef'
    )
    expect(jr.requester).toEqual('0x9FBDa871d559710256a2502A2517b794B482Db40')

    expect(jr.taskRuns.length).toEqual(1)
    expect(jr.taskRuns[0].id).toBeUndefined()
    expect(jr.taskRuns[0].index).toEqual(0)
    expect(jr.taskRuns[0].type).toEqual('httpget')
    expect(jr.taskRuns[0].status).toEqual('')
    expect(jr.taskRuns[0].confirmations).toEqual(0)
    expect(jr.taskRuns[0].minimumConfirmations).toEqual(3)
    expect(jr.taskRuns[0].error).toEqual(null)

    const [chainlinkNode, _] = await createChainlinkNode(
      db,
      'job-run-fromString-chainlink-node'
    )
    jr.chainlinkNodeId = chainlinkNode.id
    const r = await db.manager.save(jr)
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
      '0x1111111111111111111111111111111111111111111111111111111111111111'
    )
    expect(ethtxTask.transactionStatus).toEqual('fulfilledRunLog')
  })

  it('creates when finishedAt is null', () => {
    const fixtureWithoutFinishedAt = Object.assign({}, fixture, {
      finishedAt: null
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
  it('overwrites taskRun values on conflict', async () => {
    const [chainlinkNode, _] = await createChainlinkNode(
      db,
      'testOverwriteTaskRunsOnConflict'
    )

    const jr = fromString(JSON.stringify(fixture))
    jr.chainlinkNodeId = chainlinkNode.id
    await saveJobRunTree(db, jr)

    const initial = await db.manager.findOne(JobRun)
    expect(initial.taskRuns[0].confirmations).toEqual(0)

    const updatedJr = fromString(JSON.stringify(updateFixture))
    updatedJr.chainlinkNodeId = chainlinkNode.id
    await saveJobRunTree(db, updatedJr)

    const actual = await db.manager.findOne(JobRun)
    expect(actual.taskRuns[0].confirmations).toEqual(3)
  })
})
