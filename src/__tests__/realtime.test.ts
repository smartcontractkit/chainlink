import { Server } from 'http'
import WebSocket from 'ws'
import { start as startServer, DEFAULT_TEST_PORT } from '../support/server'
import { closeDbConnection, getDb } from '../database'
import createFixture from './fixtures/JobRun.fixture.json'
import updateFixture from './fixtures/JobRunUpdate.fixture.json'
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

  it('create a job run for valid JSON', async (done: any) => {
    expect.assertions(3)

    const db = await getDb()

    const ws = new WebSocket(ENDPOINT)
    ws.on('open', () => {
      ws.send(JSON.stringify(createFixture))
    })

    ws.on('message', async (data: any) => {
      const result = JSON.parse(data)
      expect(result.status).toEqual(201)

      const jobRunCount = await db.manager.count(JobRun)
      expect(jobRunCount).toEqual(1)

      const taskRunCount = await db.manager.count(TaskRun)
      expect(taskRunCount).toEqual(1)

      ws.close()
      done()
    })
  })

  it('can create and update a job run and task runs', async (done: any) => {
    expect.assertions(6)

    const db = await getDb()
    const assertionCallback = async () => {
      const jobRunCount = await db.manager.count(JobRun)
      expect(jobRunCount).toEqual(1)

      const taskRunCount = await db.manager.count(TaskRun)
      expect(taskRunCount).toEqual(1)

      const jr = await db.manager.findOne(JobRun, { relations: ['taskRuns'] })
      expect(jr.status).toEqual('completed')

      const tr = jr.taskRuns[0]
      expect(tr.status).toEqual('completed')
      done()
    }

    const ws = new WebSocket(ENDPOINT)
    let responses = 0
    ws.on('message', (data: any) => {
      responses += 1
      const result = JSON.parse(data)
      expect(result.status).toEqual(201)
      if (responses === 2) {
        ws.close()
        return assertionCallback()
      }
    })

    // send payloads
    ws.on('open', () => {
      ws.send(JSON.stringify(createFixture))
      ws.send(JSON.stringify(updateFixture))
    })
  })

  it('can handle malformed JSON', async (done: any) => {
    expect.assertions(2)

    const db = await getDb()

    const ws = new WebSocket(ENDPOINT)
    ws.on('open', () => {
      ws.send('{invalid json}')
    })
    ws.on('message', async (data: any) => {
      const result = JSON.parse(data)
      expect(result.status).toEqual(422)

      const count = await db.manager.count(JobRun)
      expect(count).toEqual(0)
      ws.close()

      const secondWs = new WebSocket(ENDPOINT)
      secondWs.on('open', () => {
        secondWs.close()
        done()
      })
    })
  })
})
