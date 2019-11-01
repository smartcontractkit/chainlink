import WebSocket from 'ws'
import jayson from 'jayson'
import { DEFAULT_TEST_PORT } from './server'
import { ACCESS_KEY_HEADER, SECRET_HEADER } from '../utils/constants'

export const ENDPOINT = `ws://localhost:${DEFAULT_TEST_PORT}`

export const newChainlinkNode = (
  url: string,
  accessKey: string,
  secret: string,
): Promise<WebSocket> => {
  const ws = new WebSocket(ENDPOINT, {
    headers: {
      [ACCESS_KEY_HEADER]: accessKey,
      [SECRET_HEADER]: secret,
    },
  })

  return new Promise((resolve: (arg0: WebSocket) => void, reject) => {
    ws.on('error', (error: Error) => {
      reject(error)
    })

    ws.on('open', () => resolve(ws))
  })
}

const jsonClient = new jayson.Client(null, null)
export const createRPCRequest = (method: string, params?: any) =>
  jsonClient.request(method, params)
