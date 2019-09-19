import request from 'supertest'
import http from 'http'
import { Connection } from 'typeorm'
import { closeDbConnection, getDb } from '../../../database'
import { clearDb } from '../../testdatabase'
import { createAdmin } from '../../../support/admin'
import {
  ADMIN_USERNAME_HEADER,
  ADMIN_PASSWORD_HEADER,
} from '../../../utils/constants'
import { start as testServer } from '../../../support/server'

const USERNAME = 'myadmin'
const PASSWORD = 'validpassword'
const adminLoginPath = '/api/v1/admin/login'

let server: http.Server
let db: Connection

function sendPost(path: string, username: string, password: string) {
  return request(server)
    .post(adminLoginPath)
    .set('Accept', 'application/json')
    .set('Content-Type', 'application/json')
    .set(ADMIN_USERNAME_HEADER, username)
    .set(ADMIN_PASSWORD_HEADER, password)
}

beforeAll(async () => {
  db = await getDb()
  server = await testServer()
})
afterAll(async done => {
  if (server) {
    server.close(done)
    await closeDbConnection()
  }
})

describe('#index', () => {
  beforeEach(async () => {
    await clearDb()
    await createAdmin(db, USERNAME, PASSWORD)
  })

  it('returns a 200 with valid credentials', done => {
    sendPost(adminLoginPath, USERNAME, PASSWORD)
      .expect(200)
      .end(done)
  })

  it('returns a 401 unauthorized with invalid admin credentials', done => {
    request(server)
      .post(adminLoginPath)
      .set('Content-Type', 'application/json')
      .set('Explorer-Admin-Username', USERNAME)
      .set('Explorer-Admin-Password', 'invalidpassword')
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
