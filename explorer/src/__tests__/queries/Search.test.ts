import fixture from '../fixtures/JobRun.fixture.json'
import { closeDbConnection, getDb } from '../../database'
import { Connection } from 'typeorm'
import { createChainlinkNode } from '../../entity/ChainlinkNode'
import { fromString } from '../../entity/JobRun'
import { search } from '../../queries/search'

let db: Connection

beforeAll(async () => {
  db = await getDb()
})

afterAll(async () => closeDbConnection())

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
      jobId: 'jobB',
      runId: 'runB'
    })
    const jrB = fromString(JSON.stringify(fixtureB))
    jrB.chainlinkNodeId = chainlinkNode.id
    jrB.createdAt = new Date('2019-04-09T01:00:00.000Z')
    jrB.txHash = 'fixtureBTxHash'
    jrB.requester = 'fixtureBRequester'
    jrB.requestId = 'fixtureBRequestID'
    await db.manager.save(jrB)

    const fixtureC = Object.assign({}, fixture, {
      jobId: 'jobB',
      runId: 'runC'
    })
    const jrC = fromString(JSON.stringify(fixtureC))
    jrC.chainlinkNodeId = chainlinkNode.id
    jrC.createdAt = new Date('2019-05-09T01:00:00.000Z')
    jrC.txHash = 'fixtureCTxHash'
    jrC.requester = 'fixtureCRequester'
    jrC.requestId = 'fixtureCRequestID'
    await db.manager.save(jrC)
  })

  it('returns no results when no query is supplied', async () => {
    let results
    results = await search(db, { searchQuery: undefined })
    expect(results).toHaveLength(0)
    results = await search(db, { searchQuery: null })
    expect(results).toHaveLength(0)
  })

  it('returns results in descending order by createdAt', async () => {
    const results = await search(db, { searchQuery: 'jobB' })
    expect(results.length).toEqual(2)
    expect(results[0].runId).toEqual('runC')
    expect(results[1].runId).toEqual('runB')
  })

  it('can set a limit', async () => {
    const results = await search(db, { searchQuery: 'jobB', limit: 1 })
    expect(results).toHaveLength(1)
    expect(results[0].runId).toEqual('runC')
  })

  it('can set a page with a 1 based index', async () => {
    let results

    results = await search(db, {
      limit: 1,
      page: 1,
      searchQuery: 'jobB'
    })
    expect(results[0].runId).toEqual('runC')

    results = await search(db, {
      limit: 1,
      page: 2,
      searchQuery: 'jobB'
    })
    expect(results[0].runId).toEqual('runB')
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
    const results = await search(db, { searchQuery: 'runB runC' })
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
