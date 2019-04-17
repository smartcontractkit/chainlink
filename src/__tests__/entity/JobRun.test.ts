import { createDbConnection, closeDbConnection, getDb } from '../../database'
import { fromString, search } from '../../entity/JobRun'
import { JOB_RUN_A_ID, JOB_RUN_B_ID } from '../../seed'
import fixture from '../fixtures/JobRun.fixture.json'

beforeAll(async () => createDbConnection())
afterAll(async () => closeDbConnection())

describe('fromString', () => {
  it('successfully creates a run and tasks from json', async () => {
    const jr = fromString(JSON.stringify(fixture))
    expect(jr.id).toBeUndefined()
    expect(jr.runId).toEqual(JOB_RUN_A_ID)
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

    const r = await getDb().manager.save(jr)
    expect(r.id).toBeDefined()
    expect(r.type).toEqual('runlog')
    expect(r.taskRuns.length).toEqual(1)
  })

  it('creates when completedAt is null', () => {
    const fixtureWithoutCompletedAt = Object.assign({}, fixture, {
      completedAt: null
    })
    const jr = fromString(JSON.stringify(fixtureWithoutCompletedAt))
    expect(jr.runId).toEqual(JOB_RUN_A_ID)
    expect(jr.completedAt).toEqual(null)
  })

  it('errors on a malformed string', async () => {
    try {
      fromString(`{"absolute":garbage`)
    } catch (err) {
      expect(err).toBeDefined()
    }
  })
})

describe('search', () => {
  beforeEach(async () => {
    const jrA = fromString(JSON.stringify(fixture))
    jrA.createdAt = new Date(Date.parse('2019-04-08T01:00:00.000Z'))
    await getDb().manager.save(jrA)

    const fixtureB = Object.assign({}, fixture, {
      runId: JOB_RUN_B_ID,
      jobId: JOB_RUN_B_ID
    })
    const jrB = fromString(JSON.stringify(fixtureB))
    jrB.createdAt = new Date(Date.parse('2019-04-09T01:00:00.000Z'))
    await getDb().manager.save(jrB)
  })

  it('returns all results when no query is supplied', async () => {
    let results
    results = await search(getDb(), { searchQuery: undefined })
    expect(results).toHaveLength(2)
    results = await search(getDb(), { searchQuery: null })
    expect(results).toHaveLength(2)
  })

  it('returns results in descending order by createdAt', async () => {
    const results = await search(getDb(), { searchQuery: undefined })
    expect(results[0].runId).toEqual(JOB_RUN_B_ID)
    expect(results[1].runId).toEqual(JOB_RUN_A_ID)
  })

  it('can set a limit', async () => {
    const results = await search(getDb(), { searchQuery: undefined, limit: 1 })
    expect(results[0].runId).toEqual(JOB_RUN_B_ID)
    expect(results).toHaveLength(1)
  })

  it('can set a page with a 1 based index', async () => {
    let results

    results = await search(getDb(), {
      searchQuery: undefined,
      page: 1,
      limit: 1
    })
    expect(results[0].runId).toEqual(JOB_RUN_B_ID)

    results = await search(getDb(), {
      searchQuery: undefined,
      page: 2,
      limit: 1
    })
    expect(results[0].runId).toEqual(JOB_RUN_A_ID)
  })

  it('returns no results for blank search', async () => {
    const results = await search(getDb(), { searchQuery: '' })
    expect(results).toHaveLength(0)
  })

  it('returns one result for an exact match on jobId', async () => {
    const results = await search(getDb(), { searchQuery: JOB_RUN_A_ID })
    expect(results).toHaveLength(1)
  })

  it('returns one result for an exact match on jobId and runId', async () => {
    const results = await search(getDb(), {
      searchQuery: `${JOB_RUN_A_ID} aeb2861d306645b1ba012079aeb2e53a`
    })
    expect(results).toHaveLength(1)
  })

  it('returns two results when two matching runIds are supplied', async () => {
    const results = await search(getDb(), {
      searchQuery: `${JOB_RUN_A_ID} ${JOB_RUN_B_ID}`
    })
    expect(results).toHaveLength(2)
  })
})
