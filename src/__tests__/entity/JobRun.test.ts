import { createDbConnection, closeDbConnection, getDb } from '../../database'
import { clearDb } from '../testdatabase'
import { fromString } from '../../entity/JobRun'
import fixture from './JobRun.fixture.json'

beforeAll(async () => createDbConnection())
afterAll(async () => closeDbConnection())
beforeEach(async () => clearDb())

describe('fromString', () => {
  it('successfully creates a run from json', async () => {
    const jr = fromString(JSON.stringify(fixture))
    expect(jr.id).toBeUndefined()
    expect(jr.jobRunId).toEqual(fixture.id)
    expect(jr.jobId).toEqual(fixture.jobId)
    expect(jr.createdAt).toEqual(new Date(fixture.createdAt))

    const e = await getDb().manager.save(jr)
    expect(e.id).toBeDefined()
  })

  it('errors on a malformed string', async () => {
    try {
      const jr = fromString(`{"absolute":garbage`)
    } catch (err) {
      expect(err).toBeDefined()
    }
  })
})
