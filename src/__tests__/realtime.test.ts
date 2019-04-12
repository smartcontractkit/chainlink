import { Server } from 'http'
import WebSocket from 'ws'
import { start as startServer, DEFAULT_TEST_PORT } from '../support/server'
import { closeDbConnection, getDb } from '../database'
import fixture from './fixtures/JobRun.fixture.json'
import { JobRun } from '../entity/JobRun'
import { TaskRun } from '../entity/TaskRun'

const ENDPOINT = `ws://localhost:${DEFAULT_TEST_PORT}`

describe('realtime', () => {
  let server: Server

  beforeAll(async () => {
    server = await startServer()
  })
  afterAll(async () => {
    return Promise.all([server.close(), closeDbConnection()])
  })

  it('can handle malformed JSON & create a job run for valid JSON', async (done: any) => {
    expect.assertions(5)

    const db = await getDb()

    const wsA = new WebSocket(ENDPOINT)
    wsA.on('open', () => {
      wsA.send('{invalid json}')
    })
    wsA.on('message', async (data: any) => {
      const result = JSON.parse(data)
      expect(result.status).toEqual(422)

      const count = await db.manager.count(JobRun)
      expect(count).toEqual(0)

      const wsB = new WebSocket(ENDPOINT)
      wsB.on('open', () => {
        wsB.send(JSON.stringify(fixture))
      })
      wsB.on('message', async (data: any) => {
        const result = JSON.parse(data)
        expect(result.status).toEqual(201)

        const jobRunCount = await db.manager.count(JobRun)
        expect(jobRunCount).toEqual(1)

        const taskRunCount = await db.manager.count(TaskRun)
        expect(taskRunCount).toEqual(1)

        wsA.close()
        wsB.close()
        done()
      })
    })
  })
})
