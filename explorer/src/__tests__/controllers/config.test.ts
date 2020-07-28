import http from 'http'
import request from 'supertest'
import { start, stop } from '../../support/server'
import { Config } from '../../config'

let server: http.Server
beforeAll(async () => {
  server = await start()
})
afterAll(done => stop(server, done))

describe('GET /api/v1/config', () => {
  fit('can return ga id', async () => {
    const gaId = 'GA-ABC'
    Config.setEnv('GA_ID', gaId)
    const response = await request(server).get('/api/v1/config')
    expect(response.status).toEqual(200)
    expect(response.body.gaId).toEqual(gaId)
  })
})
