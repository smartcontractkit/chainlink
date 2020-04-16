import bodyParser from 'body-parser'
import cookieSession from 'cookie-session'
import { randomBytes } from 'crypto'
import express from 'express'
import http from 'http'
import httpStatus from 'http-status-codes'
import request from 'supertest'
import adminAuth from '../../middleware/adminAuth'
import { createAdmin } from '../../support/admin'
import { stop } from '../../support/server'
import {
  ADMIN_PASSWORD_HEADER,
  ADMIN_PASSWORD_PARAM,
  ADMIN_USERNAME_HEADER,
  ADMIN_USERNAME_PARAM,
} from '../../utils/constants'
import { clearDb } from '../testdatabase'

const USERNAME = 'myadmin'
const PASSWORD = 'validpassword'
const ADMIN_PATH = '/api/v1/admin'
const ROUTE_PATH = `${ADMIN_PATH}/test-route`

const app = express()
app.use(bodyParser.json())
app.use(
  cookieSession({
    name: 'explorer',
    maxAge: 60_000,
    secret: randomBytes(32).toString(),
  }),
)
app.use(ROUTE_PATH, adminAuth)
app.use(ROUTE_PATH, (_, res) => res.sendStatus(200))

let server: http.Server

beforeAll(async () => {
  server = app.listen(null)
})
afterAll(done => stop(server, done))

beforeEach(async () => {
  await clearDb()
  await createAdmin(USERNAME, PASSWORD)
})

function sendPostHeaders(path: string, username: string, password: string) {
  return request(server)
    .post(path)
    .send({})
    .set('Accept', 'application/json')
    .set('Content-Type', 'application/json')
    .set(ADMIN_USERNAME_HEADER, username)
    .set(ADMIN_PASSWORD_HEADER, password)
}

function sendPostBody(path: string, username: string, password: string) {
  return request(server)
    .post(path)
    .send({
      [ADMIN_USERNAME_PARAM]: username,
      [ADMIN_PASSWORD_PARAM]: password,
    })
    .set('Accept', 'application/json')
    .set('Content-Type', 'application/json')
}

describe('adminAuth middleware headers', () => {
  it('executes the next middleware when authentication is successful', done => {
    sendPostHeaders(ROUTE_PATH, USERNAME, PASSWORD)
      .expect(httpStatus.OK)
      .end(done)
  })

  it('responds with a 401 when there is no username or password', done => {
    sendPostHeaders(ROUTE_PATH, '', PASSWORD)
      .expect(httpStatus.UNAUTHORIZED)
      .end(done)

    sendPostHeaders(ROUTE_PATH, USERNAME, '')
      .expect(httpStatus.UNAUTHORIZED)
      .end(done)
  })

  it('responds with a 401 when the password is invalid', done => {
    sendPostHeaders(ROUTE_PATH, USERNAME, 'invalidpassword')
      .expect(httpStatus.UNAUTHORIZED)
      .end(done)
  })
})

describe('adminAuth middleware body', () => {
  it('executes the next middleware when authentication is successful', done => {
    sendPostBody(ROUTE_PATH, USERNAME, PASSWORD)
      .expect(httpStatus.OK)
      .end(done)
  })

  it('responds with a 401 when there is no username or password', done => {
    sendPostBody(ROUTE_PATH, '', PASSWORD)
      .expect(httpStatus.UNAUTHORIZED)
      .end(done)

    sendPostBody(ROUTE_PATH, USERNAME, '')
      .expect(httpStatus.UNAUTHORIZED)
      .end(done)
  })

  it('responds with a 401 when the password is invalid', done => {
    sendPostBody(ROUTE_PATH, USERNAME, 'invalidpassword')
      .expect(httpStatus.UNAUTHORIZED)
      .end(done)
  })
})
