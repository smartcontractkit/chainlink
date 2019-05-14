import fixture from '../fixtures/JobRun.fixture.json'
import { closeDbConnection, getDb } from '../../database'
import { Connection } from 'typeorm'
import { createChainlinkNode } from '../../entity/ChainlinkNode'
import ethtxFixture from '../fixtures/JobRun.ethtx.fixture.json'
import { fromString } from '../../entity/JobRun'

let db: Connection

beforeAll(async () => {
  db = await getDb()
})

afterAll(async () => closeDbConnection())

describe('fromString', () => {
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
