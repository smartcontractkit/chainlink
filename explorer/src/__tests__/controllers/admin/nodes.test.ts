import request, { Response } from 'supertest'
import http from 'http'
import { Connection } from 'typeorm'
import { closeDbConnection, getDb } from '../../../database'
import { clearDb } from '../../testdatabase'
import { createAdmin } from '../../../support/admin'
import {
  createChainlinkNode,
  find as findNode,
} from '../../../entity/ChainlinkNode'
import { start as testServer } from '../../../support/server'
import {
  ADMIN_USERNAME_HEADER,
  ADMIN_PASSWORD_HEADER,
} from '../../../utils/constants'

const USERNAME = 'myadmin'
const PASSWORD = 'validpassword'
const ADMIN_PATH = '/api/v1/admin'
const adminNodesPath = `${ADMIN_PATH}/nodes`

let server: http.Server
let db: Connection

beforeAll(async () => {
  db = await getDb()
  server = await testServer()
})
afterAll(async () => {
  if (server) {
    server.close()
    await closeDbConnection()
  }
})
beforeEach(async () => {
  await clearDb()
  await createAdmin(db, USERNAME, PASSWORD)
})

describe('POST /nodes', () => {
  it('can create a node and returns the generated information', done => {
    request(server)
      .post(adminNodesPath)
      .send({ name: 'nodeA', url: 'http://nodea.com' })
      .set('Accept', 'application/json')
      .set('Content-Type', 'application/json')
      .set(ADMIN_USERNAME_HEADER, USERNAME)
      .set(ADMIN_PASSWORD_HEADER, PASSWORD)
      .expect(201)
      .expect((res: Response) => {
        expect(res.body.id).toBeDefined()
        expect(res.body.accessKey).toBeDefined()
        expect(res.body.secret).toBeDefined()
      })
      .end(done)
  })

  it('returns an error with invalid params', done => {
    request(server)
      .post(adminNodesPath)
      .send({ url: 'http://nodea.com' })
      .set('Accept', 'application/json')
      .set('Content-Type', 'application/json')
      .set(ADMIN_USERNAME_HEADER, USERNAME)
      .set(ADMIN_PASSWORD_HEADER, PASSWORD)
      .expect(422)
      .expect((res: Response) => {
        const { errors } = res.body

        expect(errors).toBeDefined()
        expect(errors.name).toEqual({
          minLength: 'must be at least 3 characters',
        })
      })
      .end(done)
  })

  it('returns an error when the node already exists', async done => {
    const [node] = await createChainlinkNode(db, 'nodeA')

    request(server)
      .post(adminNodesPath)
      .send({ name: node.name })
      .set('Accept', 'application/json')
      .set('Content-Type', 'application/json')
      .set(ADMIN_USERNAME_HEADER, USERNAME)
      .set(ADMIN_PASSWORD_HEADER, PASSWORD)
      .expect(409)
      .end(done)
  })

  it('returns a 401 unauthorized with invalid admin credentials', done => {
    request(server)
      .post(adminNodesPath)
      .set('Content-Type', 'application/json')
      .set(ADMIN_USERNAME_HEADER, USERNAME)
      .set(ADMIN_PASSWORD_HEADER, 'invalidpassword')
      .expect(401)
      .end(done)
  })
})

describe('DELETE /nodes/:id', () => {
  function path(id: number): string {
    return `${adminNodesPath}/${id}`
  }

  it('can delete a node', async done => {
    const [node, _] = await createChainlinkNode(db, 'nodeA')

    request(server)
      .delete(path(node.id))
      .set('Content-Type', 'application/json')
      .set(ADMIN_USERNAME_HEADER, USERNAME)
      .set(ADMIN_PASSWORD_HEADER, PASSWORD)
      .expect(200)
      .expect(async () => {
        const nodeAfter = await findNode(db, node.id)
        expect(nodeAfter).not.toBeDefined()
      })
      .end(done)
  })

  it('returns a 401 unauthorized with invalid admin credentials', done => {
    request(server)
      .delete(path(10))
      .set('Content-Type', 'application/json')
      .set(ADMIN_USERNAME_HEADER, USERNAME)
      .set(ADMIN_PASSWORD_HEADER, 'invalidpassword')
      .expect(401)
      .end(done)
  })
})
