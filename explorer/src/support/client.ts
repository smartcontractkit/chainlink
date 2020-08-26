import jayson from 'jayson'
import WebSocket from 'ws'
import { Config } from '../config'
import {
  ACCESS_KEY_HEADER,
  CORE_SHA_HEADER,
  CORE_VERSION_HEADER,
  SECRET_HEADER,
} from '../utils/constants'

export const newChainlinkNode = (
  accessKey: string,
  secret: string,
  coreVersion?: string,
  coreSha?: string,
): Promise<WebSocket> => {
  const headers: any = {
    [ACCESS_KEY_HEADER]: accessKey,
    [SECRET_HEADER]: secret,
  }
  if (coreVersion) {
    headers[CORE_VERSION_HEADER] = coreVersion
  }
  if (coreSha) {
    headers[CORE_SHA_HEADER] = coreSha
  }

  const ws = new WebSocket(`ws://localhost:${Config.testPort()}`, {
    headers,
  })

  return new Promise((resolve, reject) => {
    ws.on('error', error => {
      error.message =
        '[newChainlinkNode] Error on opening websocket:' + error.message
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
