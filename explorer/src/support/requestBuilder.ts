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

export function requestBuilder(server: Server): RequestBuilder {
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
    username: string,
    password: string,
    data?: object,
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
      .set('Accept', 'application/json')
      .set('Content-Type', 'application/json')
      .set(ADMIN_USERNAME_HEADER, username)
      .set(ADMIN_PASSWORD_HEADER, password)
  }

  return {
    sendGet,
    sendPost,
    sendDelete,
  }
}
