import { createDbConnection, closeDbConnection, getDb } from '../../database'
import { clearDb } from '../testdatabase'
import { fromString } from '../../entity/JobRun'
import fixture from './JobRun.fixture.json'

beforeAll(async () => createDbConnection())
afterAll(async () => closeDbConnection())
beforeEach(async () => clearDb())

describe('fromString', () => {
  it('successfully creates a run and tasks from json', async () => {
    const jr = fromString(JSON.stringify(fixture))
    expect(jr.id).toEqual(fixture.id)
    expect(jr.jobId).toEqual(fixture.jobId)
    expect(jr.status).toEqual(fixture.status)
    expect(jr.initiatorType).toEqual(fixture.initiator.type)
    expect(jr.createdAt).toEqual(new Date(fixture.createdAt))
    expect(jr.completedAt).toEqual(new Date(fixture.completedAt))
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
