import fetch from 'node-fetch'
import httpStatus from 'http-status-codes'
import {
  ADMIN_USERNAME_HEADER,
  ADMIN_PASSWORD_HEADER,
} from '../utils/constants'
import { Config } from '../config'

interface CreateChainlinkNode {
  name: string
  url?: string
}

interface CreateChainlinkNodeOk {
  id: string
  accessKey: string
  secret: string
}

function logError(msg: string) {
  console.error(msg)
}

function logUnauthorized() {
  logError(
    'Invalid admin credentials. Please ensure the you have provided the correct admin username and password.',
  )
}

function logNotFound() {
  logError(
    `Error creating chainlink node. API endpoint not found. Have you set the correct EXPLORER_BASE_URL?`,
  )
}

export const add = async (name: string, url?: string) => {
  const createNodeUrl = `${Config.baseUrl()}/api/v1/admin/nodes`
  const data: CreateChainlinkNode = { name, url }
  const response = await fetch(createNodeUrl, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      [ADMIN_USERNAME_HEADER]: Config.adminUsername(),
      [ADMIN_PASSWORD_HEADER]: Config.adminPassword(),
    },
    body: JSON.stringify(data),
  })

  switch (response.status) {
    case httpStatus.CREATED: {
      const chainlinkNode: CreateChainlinkNodeOk = await response.json()
      console.log('Created new chainlink node with id %s', chainlinkNode.id)
      console.log('AccessKey', chainlinkNode.accessKey)
      console.log('Secret', chainlinkNode.secret)
      break
    }
    case httpStatus.NOT_FOUND:
      logNotFound()
      break
    case httpStatus.UNAUTHORIZED:
      logUnauthorized()
      break
    case httpStatus.CONFLICT:
      console.error(
        `Error creating chainlink node. A node with the name: "${name}" already exists.`,
      )
      break
    default:
      console.error(`Unhandled error. HTTP status: ${response.status}`)
      break
  }
}

export const remove = async (name: string) => {
  const deleteNodeUrl = `${Config.baseUrl()}/api/v1/admin/nodes/${name}`
  const response = await fetch(deleteNodeUrl, {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
      [ADMIN_USERNAME_HEADER]: Config.adminUsername(),
      [ADMIN_PASSWORD_HEADER]: Config.adminPassword(),
    },
  })

  switch (response.status) {
    case httpStatus.OK:
      console.log('Successfully deleted chainlink node with name %s', name)
      break
    case httpStatus.NOT_FOUND:
      logNotFound()
      break
    case httpStatus.UNAUTHORIZED:
      logUnauthorized()
      break
    default:
      console.error(`Unhandled error. HTTP status: ${response.status}`)
      break
  }
}
