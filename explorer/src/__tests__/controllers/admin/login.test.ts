import request from 'supertest'
import http from 'http'
import express from 'express'
import { Connection } from 'typeorm'
import { closeDbConnection, getDb } from '../../../database'
import { clearDb } from '../../testdatabase'
import { createAdmin } from '../../../support/admin'
import adminLogin from '../../../controllers/admin/login'

const USERNAME = 'myadmin'
const PASSWORD = 'validpassword'

const app = express()
app.use(express.json())
app.use('/api/v1', adminLogin)

let server: http.Server
let db: Connection

beforeAll(async () => {
  db = await getDb()
  server = app.listen(null)
})
afterAll(async () => {
  if (server) {
    server.close()
    await closeDbConnection()
  }
})

describe('#index', () => {
  const adminLoginPath = '/api/v1/admin/login'

  beforeEach(async () => {
    await clearDb()
    await createAdmin(db, USERNAME, PASSWORD)
  })

  it('returns a 200 with valid credentials', done => {
    request(server)
      .post(adminLoginPath)
      .set('Content-Type', 'application/json')
      .set('Explorer-Admin-Username', USERNAME)
      .set('Explorer-Admin-Password', PASSWORD)
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
