import request from 'supertest'
import { Server } from 'http'
import {
  ADMIN_USERNAME_HEADER,
  ADMIN_PASSWORD_HEADER,
} from '../utils/constants'

type RequestFunction = (
  path: string,
  username: string,
  password: string,
  data?: object,
) => request.Test

export interface RequestBuilder {
  sendGet: RequestFunction
  sendPost: RequestFunction
  sendDelete: RequestFunction
}

type method = 'get' | 'post' | 'delete'

export function requestBuilder(server: Server): RequestBuilder {
  function sendRequest(
    method: method,
    path: string,
    username: string,
    password: string,
    data?: object,
  ) {
    return request(server)
      [method](path)
      .send(data)
      .set('Accept', 'application/json')
      .set('Content-Type', 'application/json')
      .set(ADMIN_USERNAME_HEADER, username)
      .set(ADMIN_PASSWORD_HEADER, password)
  }

  const buildRequestMethod = (method: method) => (
    path: string,
    username: string,
    password: string,
    data?: object,
  ) => sendRequest(method, path, username, password, data)

  return {
    sendGet: buildRequestMethod('get'),
    sendPost: buildRequestMethod('post'),
    sendDelete: buildRequestMethod('delete'),
  }
}
