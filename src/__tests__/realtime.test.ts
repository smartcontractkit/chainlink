import { Server } from 'http'
import WebSocket from 'ws'
import { start as startServer, DEFAULT_TEST_PORT } from '../support/server'
import { closeDbConnection, getDb } from '../database'
import fixture from './fixtures/JobRun.fixture.json'
import { JobRun } from '../entity/JobRun'

const ENDPOINT = `ws://localhost:${DEFAULT_TEST_PORT}`

describe('realtime', () => {
  let server: Server

  beforeAll(async () => {
    server = await startServer()
  })
  afterAll(async () => {
    return Promise.all([server.close(), closeDbConnection()])
  })

  it('can handle malformed JSON & create a job run for valid JSON', (done: any) => {
    expect.assertions(2)

    const wsA = new WebSocket(ENDPOINT)
    wsA.on('open', () => {
      wsA.send('{invalid json}')
    })
    wsA.on('message', (data: any) => {
      const result = JSON.parse(data)
      expect(result.status).toEqual(422)
    })

    const wsB = new WebSocket(ENDPOINT)
    wsB.on('open', () => {
      wsB.send(JSON.stringify(fixture))
    })
    wsB.on('message', (data: any) => {
      const result = JSON.parse(data)
      expect(result.status).toEqual(201)

      wsA.close()
      wsB.close()
      done()
    })
  })
})
