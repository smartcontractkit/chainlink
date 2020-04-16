import http from 'http'
import request from 'supertest'
import { createAdmin } from '../../../support/admin'
import { requestBuilder, RequestBuilder } from '../../../support/requestBuilder'
import { start, stop } from '../../../support/server'

const USERNAME = 'myadmin'
const PASSWORD = 'validpassword'
const adminLoginPath = '/api/v1/admin/login'

let server: http.Server
let rb: RequestBuilder

beforeAll(async () => {
  server = await start()
  rb = requestBuilder(server)
})
afterAll(done => stop(server, done))

describe('POST /api/v1/admin/login', () => {
  beforeEach(async () => {
    await createAdmin(USERNAME, PASSWORD)
  })

  it('returns a 200 with valid credentials', done => {
    rb.sendPost(adminLoginPath, USERNAME, PASSWORD)
      .expect(200)
      .expect(res => {
        expect(res.body).toEqual({})
      })
      .end(done)
  })

  it('returns a 401 unauthorized with invalid admin credentials', done => {
    rb.sendPost(adminLoginPath, USERNAME, 'invalidpassword')
      .expect(401)
      .end(done)
  })

  it('returns a 401 unauthorized when the username does not exist', done => {
    request(server)
      .post(adminLoginPath)
      .set('Content-Type', 'application/json')
      .expect(401)
      .end(done)
  })
})
