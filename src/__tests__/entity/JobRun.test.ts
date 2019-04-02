import { createDbConnection, closeDbConnection, getDb } from '../../database'
import { clearDb } from '../testdatabase'
import { fromString, search } from '../../entity/JobRun'
import { JOB_RUN_A_ID, JOB_RUN_B_ID } from '../../seed'
import fixture from './JobRun.fixture.json'
import { Option } from 'prelude-ts'

beforeAll(async () => createDbConnection())
afterAll(async () => closeDbConnection())
beforeEach(async () => clearDb())

describe('fromString', () => {
  it('successfully creates a run and tasks from json', async () => {
    const jr = fromString(JSON.stringify(fixture))
    expect(jr.runId).toEqual(JOB_RUN_A_ID)
    expect(jr.jobId).toEqual('aeb2861d306645b1ba012079aeb2e53a')
    expect(jr.createdAt).toEqual(new Date('2019-04-01T22:07:04Z'))
    expect(jr.status).toEqual('in_progress')
    expect(jr.completedAt).toEqual(new Date('2018-04-01T22:07:04Z'))
    expect(jr.initiatorType).toEqual(fixture.initiator.type)
    expect(jr.taskRuns.length).toEqual(1)
    expect(jr.taskRuns[0].id).toEqual(fixture.taskRuns[0].id)
    expect(jr.taskRuns[0].index).toEqual(0)
    expect(jr.taskRuns[0].type).toEqual(fixture.taskRuns[0].task.type)
    expect(jr.taskRuns[0].status).toEqual(fixture.taskRuns[0].status)
    expect(jr.taskRuns[0].error).toEqual(fixture.taskRuns[0].result.error)

    const r = await getDb().manager.save(jr)
    expect(r.id).toBeDefined()
    expect(r.taskRuns.length).toEqual(1)
    expect(r.taskRuns[0].id).toBeDefined()
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
      const jr = fromString(`{"absolute":garbage`)
    } catch (err) {
      expect(err).toBeDefined()
    }
  })
})

describe('search', () => {
  beforeEach(async () => {
    const jrA = fromString(JSON.stringify(fixture))
    await getDb().manager.save(jrA)

    const fixtureB = Object.assign({}, fixture, {
      runId: JOB_RUN_B_ID,
      jobId: JOB_RUN_B_ID
    })
    const jrB = fromString(JSON.stringify(fixtureB))
    await getDb().manager.save(jrB)
  })

  it('returns all results when no query is supplied', async () => {
    const results = await search(getDb(), Option.none())
    expect(results).toHaveLength(2)
  })

  it('returns no results for blank search', async () => {
    const results = await search(getDb(), Option.of(''))
    expect(results).toHaveLength(0)
  })

  it('returns one result for an exact match on jobId', async () => {
    const results = await search(getDb(), Option.of(JOB_RUN_A_ID))
    expect(results).toHaveLength(1)
  })

  it('returns one result for an exact match on jobId and runId', async () => {
    const results = await search(
      getDb(),
      Option.of(`${JOB_RUN_A_ID} aeb2861d306645b1ba012079aeb2e53a`)
    )
    expect(results).toHaveLength(1)
  })

  it('returns two results when two matching runIds are supplied', async () => {
    const results = await search(
      getDb(),
      Option.of(`${JOB_RUN_A_ID} ${JOB_RUN_B_ID}`)
    )
    expect(results).toHaveLength(2)
  })
})
