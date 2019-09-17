import { default as fetch, Response } from 'node-fetch'
import httpStatus from 'http-status-codes'
import {
  ADMIN_USERNAME_HEADER,
  ADMIN_PASSWORD_HEADER,
} from '../utils/constants'

const EXPLORER_BASE_URL =
  process.env.EXPLORER_BASE_URL || 'https://explorer.chain.link'
const EXPLORER_ADMIN_USERNAME = process.env.EXPLORER_ADMIN_USERNAME
const EXPLORER_ADMIN_PASSWORD = process.env.EXPLORER_ADMIN_PASSWORD

interface CreateChainlinkNode {
  name: string
  url?: string
}

interface CreateChainlinkNodeOk {
  id: string
  accessKey: string
  secret: string
}

const add = async (name: string, url?: string) => {
  const createNodeUrl = `${EXPLORER_BASE_URL}/api/v1/admin/nodes`
  const data: CreateChainlinkNode = { name: name }
  if (url) {
    data.url = url
  }
  const response: Response = await fetch(createNodeUrl, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      [ADMIN_USERNAME_HEADER]: EXPLORER_ADMIN_USERNAME,
      [ADMIN_PASSWORD_HEADER]: EXPLORER_ADMIN_PASSWORD,
    },
    body: JSON.stringify(data),
  })

  if (response.status === httpStatus.CREATED) {
    const chainlinkNode: CreateChainlinkNodeOk = await response.json()
    console.log('created new chainlink node with id %s', chainlinkNode.id)
    console.log('AccessKey', chainlinkNode.accessKey)
    console.log('Secret', chainlinkNode.secret)
  } else if (response.status === httpStatus.UNAUTHORIZED) {
    console.log(
      'Invalid admin credentials. Please ensure the you have provided the correct admin username and password.',
    )
  } else if (response.status === httpStatus.CONFLICT) {
    console.log(
      `Error creating chainlink node. A node with the name: "${name}" already exists.`,
    )
  } else {
    console.log(`Unhandled error. HTTP status: ${response.status}`)
  }
}

const remove = async (name: string) => {
  const deleteNodeUrl = `${EXPLORER_BASE_URL}/api/v1/admin/nodes/${name}`
  const response: Response = await fetch(deleteNodeUrl, {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
      [ADMIN_USERNAME_HEADER]: EXPLORER_ADMIN_USERNAME,
      [ADMIN_PASSWORD_HEADER]: EXPLORER_ADMIN_PASSWORD,
    },
  })

  if (response.status === httpStatus.OK) {
    console.log('successfully deleted chainlink node with name %s', name)
  } else if (response.status === httpStatus.UNAUTHORIZED) {
    console.log(
      'Invalid admin credentials. Please ensure the you have provided the correct admin username and password.',
    )
  } else {
    console.log(`Unhandled error. HTTP status: ${response.status}`)
  }
}

export { add, remove }
