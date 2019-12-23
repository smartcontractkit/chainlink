import request from 'supertest'
import http from 'http'
import httpStatus from 'http-status-codes'
import { Connection } from 'typeorm'
import { getDb } from '../../../database'
import { clearDb } from '../../testdatabase'
import { createAdmin } from '../../../support/admin'
import {
  createChainlinkNode,
  find as findNode,
} from '../../../entity/ChainlinkNode'
import { start, stop } from '../../../support/server'
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
  server = await start()
})
afterAll(done => stop(server, done))
beforeEach(async () => {
  await clearDb()
  await createAdmin(db, USERNAME, PASSWORD)
})

function sendGet(path: string, username: string, password: string) {
  return request(server)
    .get(path)
    .set('Accept', 'application/json')
    .set('Content-Type', 'application/json')
    .set(ADMIN_USERNAME_HEADER, username)
    .set(ADMIN_PASSWORD_HEADER, password)
}

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

function sendDelete(path: string, username: string, password: string) {
  return request(server)
    .delete(path)
    .set('Content-Type', 'application/json')
    .set(ADMIN_USERNAME_HEADER, username)
    .set(ADMIN_PASSWORD_HEADER, password)
}

describe('POST /api/v1/admin/nodes', () => {
  it('can create a node and returns the generated information', done => {
    const data = { name: 'nodeA', url: 'http://nodea.com' }

    sendPost(adminNodesPath, data, USERNAME, PASSWORD)
      .expect(httpStatus.CREATED)
      .expect(res => {
        expect(res.body.id).toBeDefined()
        expect(res.body.accessKey).toBeDefined()
        expect(res.body.secret).toBeDefined()
      })
      .end(done)
  })

  it('returns an error with invalid params', done => {
    const data = { url: 'http://nodea.com' }

    sendPost(adminNodesPath, data, USERNAME, PASSWORD)
      .expect(httpStatus.UNPROCESSABLE_ENTITY)
      .expect(res => {
        const errors = res.body.errors

        expect(errors).toBeDefined()
        expect(errors.name).toEqual({
          minLength: 'must be at least 3 characters',
        })
      })
      .end(done)
  })

  it('returns an error when the node already exists', async done => {
    const [node] = await createChainlinkNode(db, 'nodeA')
    const data = { name: node.name }

    sendPost(adminNodesPath, data, USERNAME, PASSWORD)
      .expect(httpStatus.CONFLICT)
      .end(done)
  })

  it('returns a 401 unauthorized with invalid admin credentials', done => {
    sendPost(adminNodesPath, {}, USERNAME, 'invalidpassword')
      .expect(httpStatus.UNAUTHORIZED)
      .end(done)
  })
})

describe('DELETE /api/v1/admin/nodes/:name', () => {
  function path(name: string): string {
    return `${adminNodesPath}/${name}`
  }

  it('can delete a node', async done => {
    const [node] = await createChainlinkNode(db, 'nodeA')

    sendDelete(path(node.name), USERNAME, PASSWORD)
      .expect(httpStatus.OK)
      .expect(async () => {
        const nodeAfter = await findNode(db, node.id)
        expect(nodeAfter).not.toBeDefined()
      })
      .end(done)
  })

  it('returns a 401 unauthorized with invalid admin credentials', done => {
    sendDelete(path('idontexist'), USERNAME, 'invalidpassword')
      .expect(httpStatus.UNAUTHORIZED)
      .end(done)
  })
})

describe.only('GET /api/v1/admin/nodes/:id', () => {
  function path(id: number): string {
    return `${adminNodesPath}/${id}`
  }

  it('can get a node', async done => {
    const [node] = await createChainlinkNode(db, 'nodeA')

    sendGet(path(node.id), USERNAME, PASSWORD)
      .expect(httpStatus.OK)
      .expect(res => {
        console.log(res.body)
        expect(res.body.data.id).toBeDefined()
      })
      .end(done)
  })

  it('returns a 401 unauthorized with invalid admin credentials', async done => {
    const [node] = await createChainlinkNode(db, 'nodeA')
    const _nodePath = path(node.id)
    sendGet(_nodePath, USERNAME, 'invalidpassword')
      .expect(httpStatus.UNAUTHORIZED)
      .end(done)
  })
})
