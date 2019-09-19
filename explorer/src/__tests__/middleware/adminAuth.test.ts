import express from 'express'
import request from 'supertest'
import http from 'http'
import httpStatus from 'http-status-codes'
import cookieSession from 'cookie-session'
import { Connection } from 'typeorm'
import { closeDbConnection, getDb } from '../../database'
import { clearDb } from '../testdatabase'
import { createAdmin } from '../../support/admin'
import adminAuth from '../../middleware/adminAuth'
import {
  ADMIN_USERNAME_HEADER,
  ADMIN_PASSWORD_HEADER,
} from '../../utils/constants'

const USERNAME = 'myadmin'
const PASSWORD = 'validpassword'
const ADMIN_PATH = '/api/v1/admin'
const ROUTE_PATH = `${ADMIN_PATH}/test-route`

const app = express()
app.use(
  cookieSession({
    name: 'explorer',
    maxAge: 60_000,
    keys: ['key1', 'key2'],
  }),
)
app.use(ROUTE_PATH, adminAuth)
app.use(ROUTE_PATH, (req, res) => res.sendStatus(200))

let server: http.Server
let db: Connection

beforeAll(async () => {
  db = await getDb()
  server = app.listen(null)
})
afterAll(async done => {
  if (server) {
    server.close(done)
    await closeDbConnection()
  }
})
beforeEach(async () => {
  await clearDb()
  await createAdmin(db, USERNAME, PASSWORD)
})

function sendPost(
  path: string,
  data: object,
  username: string,
  password: string,
) {
  return request(server)
    .post(path)
    .send(data)
    .set('Accept', 'application/json')
    .set('Content-Type', 'application/json')
    .set(ADMIN_USERNAME_HEADER, username)
    .set(ADMIN_PASSWORD_HEADER, password)
}

describe('adminAuth middleware', () => {
  it('executes the next middleware when authentication is successful', done => {
    sendPost(ROUTE_PATH, {}, USERNAME, PASSWORD)
      .expect(httpStatus.OK)
      .end(done)
  })

  it('responds with a 401 when there is no username or password', done => {
    sendPost(ROUTE_PATH, {}, '', PASSWORD)
      .expect(httpStatus.UNAUTHORIZED)
      .end(done)

    sendPost(ROUTE_PATH, {}, USERNAME, '')
      .expect(httpStatus.UNAUTHORIZED)
      .end(done)
  })

  it('responds with a 401 when the password is invalid', done => {
    sendPost(ROUTE_PATH, {}, USERNAME, 'invalidpassword')
      .expect(httpStatus.UNAUTHORIZED)
      .end(done)
  })
})
