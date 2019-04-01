import { createDbConnection, closeDbConnection, getDb } from '../../database'
import { clearDb } from '../testdatabase'
import { fromString, search } from '../../entity/JobRun'
import fixture from './JobRun.fixture.json'

beforeAll(async () => createDbConnection())
afterAll(async () => closeDbConnection())
beforeEach(async () => clearDb())

describe('fromString', () => {
  it('successfully creates a run and tasks from json', async () => {
    const jr = fromString(JSON.stringify(fixture))
    expect(jr.id).toEqual('d19e1df47ecb40fa85e7f29b4c25cd6e')
    expect(jr.jobId).toEqual('b7dbc97018ce4652b79f3f17e20fce00')
    expect(jr.createdAt).toEqual(new Date('2019-03-25T12:33:34.956255-07:00'))
    expect(jr.status).toEqual('in_progresss')
    expect(jr.completedAt).toEqual(null)
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
    expect(jr.id).toEqual(fixture.id)
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
    const jr = fromString(JSON.stringify(fixture))
    await getDb().manager.save(jr)
  })

  it('returns no results for blank search', async () => {
    const results = await search(getDb(), [''])
    expect(results).toEqual([])
  })

  it('returns one result for an exact match on jobID', async () => {
    const jr = fromString(JSON.stringify(fixture))
    await getDb().manager.save(jr)

    const results = await search(getDb(), ['b7dbc97018ce4652b79f3f17e20fce00'])
    expect(results).toHaveLength(1)
  })
})
