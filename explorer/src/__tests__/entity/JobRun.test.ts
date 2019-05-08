import fixture from '../fixtures/JobRun.fixture.json'
import { closeDbConnection, getDb } from '../../database'
import { Connection } from 'typeorm'
import { createChainlinkNode } from '../../entity/ChainlinkNode'
import ethtxFixture from '../fixtures/JobRun.ethtx.fixture.json'
import { fromString, search } from '../../entity/JobRun'

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
    expect(jr.completedAt).toEqual(new Date('2018-04-01T22:07:04Z'))

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

  it('creates when completedAt is null', () => {
    const fixtureWithoutCompletedAt = Object.assign({}, fixture, {
      completedAt: null
    })
    const jr = fromString(JSON.stringify(fixtureWithoutCompletedAt))
    expect(jr.runId).toEqual('f1xtureAaaaaaaaaaaaaaaaaaaaaaaaa')
    expect(jr.completedAt).toEqual(null)
  })

  it('errors on a malformed string', async () => {
    try {
      fromString('{"absolute":garbage')
    } catch (err) {
      expect(err).toBeDefined()
    }
  })
})

describe('search', () => {
  beforeEach(async () => {
    const [chainlinkNode, _] = await createChainlinkNode(
      db,
      'job-run-search-chainlink-node'
    )

    const jrA = fromString(JSON.stringify(fixture))
    jrA.chainlinkNodeId = chainlinkNode.id
    jrA.createdAt = new Date('2019-04-08T01:00:00.000Z')
    await db.manager.save(jrA)

    const fixtureB = Object.assign({}, fixture, {
      jobId: 'f1xtureBbbbbbbbbbbbbbbbbbbbbbbbb',
      runId: 'f1xtureBbbbbbbbbbbbbbbbbbbbbbbbb'
    })
    const jrB = fromString(JSON.stringify(fixtureB))
    jrB.chainlinkNodeId = chainlinkNode.id
    jrB.createdAt = new Date('2019-04-09T01:00:00.000Z')
    jrB.txHash = 'fixtureBTxHash'
    jrB.requester = 'fixtureBRequester'
    jrB.requestId = 'fixtureBRequestID'
    await db.manager.save(jrB)
  })

  it('returns all results when no query is supplied', async () => {
    let results
    results = await search(db, { searchQuery: undefined })
    expect(results).toHaveLength(2)
    results = await search(db, { searchQuery: null })
    expect(results).toHaveLength(2)
  })

  it('returns results in descending order by createdAt', async () => {
    const results = await search(db, { searchQuery: undefined })
    expect(results[0].runId).toEqual('f1xtureBbbbbbbbbbbbbbbbbbbbbbbbb')
    expect(results[1].runId).toEqual('f1xtureAaaaaaaaaaaaaaaaaaaaaaaaa')
  })

  it('can set a limit', async () => {
    const results = await search(db, { searchQuery: undefined, limit: 1 })
    expect(results[0].runId).toEqual('f1xtureBbbbbbbbbbbbbbbbbbbbbbbbb')
    expect(results).toHaveLength(1)
  })

  it('can set a page with a 1 based index', async () => {
    let results

    results = await search(db, {
      limit: 1,
      page: 1,
      searchQuery: undefined
    })
    expect(results[0].runId).toEqual('f1xtureBbbbbbbbbbbbbbbbbbbbbbbbb')

    results = await search(db, {
      limit: 1,
      page: 2,
      searchQuery: undefined
    })
    expect(results[0].runId).toEqual('f1xtureAaaaaaaaaaaaaaaaaaaaaaaaa')
  })

  it('returns no results for blank search', async () => {
    const results = await search(db, { searchQuery: '' })
    expect(results).toHaveLength(0)
  })

  it('returns one result for an exact match on jobId', async () => {
    const results = await search(db, {
      searchQuery: 'f1xtureAaaaaaaaaaaaaaaaaaaaaaaaa'
    })
    expect(results).toHaveLength(1)
  })

  it('returns one result for an exact match on jobId and runId', async () => {
    const results = await search(db, {
      searchQuery: `${'f1xtureAaaaaaaaaaaaaaaaaaaaaaaaa'} aeb2861d306645b1ba012079aeb2e53a`
    })
    expect(results).toHaveLength(1)
  })

  it('returns two results when two matching runIds are supplied', async () => {
    const results = await search(db, {
      searchQuery: `${'f1xtureAaaaaaaaaaaaaaaaaaaaaaaaa'} ${'f1xtureBbbbbbbbbbbbbbbbbbbbbbbbb'}`
    })
    expect(results).toHaveLength(2)
  })

  it('returns one result for an exact match on requester', async () => {
    const requester = 'fixtureBRequester'
    const results = await search(db, { searchQuery: requester })
    expect(results).toHaveLength(1)
  })

  it('returns one result for a case insensitive match on requester', async () => {
    const requester = 'FIXTUREBREQUESTER'
    const results = await search(db, { searchQuery: requester })
    expect(results).toHaveLength(1)
  })

  it('returns one result for an exact match on requestId', async () => {
    const requestId = 'fixtureBRequestID'
    const results = await search(db, { searchQuery: requestId })
    expect(results).toHaveLength(1)
  })

  it('returns one result for an exact match on txHash', async () => {
    const txHash = 'fixtureBTxHash'
    const results = await search(db, { searchQuery: txHash })
    expect(results).toHaveLength(1)
  })
})
