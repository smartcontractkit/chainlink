import jayson from 'jayson'
import WebSocket from 'ws'
import { ACCESS_KEY_HEADER, SECRET_HEADER } from '../utils/constants'
import { DEFAULT_TEST_PORT } from './server'

export const newChainlinkNode = (
  accessKey: string,
  secret: string,
): Promise<WebSocket> => {
  const ws = new WebSocket(`ws://localhost:${DEFAULT_TEST_PORT}`, {
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
export const createRPCRequest = (
  method: string,
  params?: jayson.RequestParamsLike,
) => jsonClient.request(method, params)

// helper function that sends a message and only resolves once the
// rsponse is received
export const sendSingleMessage = (
  ws: WebSocket,
  request: string | object,
): Promise<any> =>
  new Promise(resolve => {
    const requestData: string =
      typeof request === 'object' ? JSON.stringify(request) : request
    ws.send(requestData)
    ws.on('message', async (data: string) => {
      const response = JSON.parse(data)
      resolve(response)
    })
  })
