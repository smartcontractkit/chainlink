import fetch from 'node-fetch'
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

export const add = async (name: string, url?: string) => {
  const createNodeUrl = `${EXPLORER_BASE_URL}/api/v1/admin/nodes`
  const data: CreateChainlinkNode = { name, url }
  const response = await fetch(createNodeUrl, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      [ADMIN_USERNAME_HEADER]: EXPLORER_ADMIN_USERNAME,
      [ADMIN_PASSWORD_HEADER]: EXPLORER_ADMIN_PASSWORD,
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
      console.error(
        `Error creating chainlink node. API endpoint not found. Have you set the correct EXPLORER_BASE_URL?`,
      )
      break
    case httpStatus.UNAUTHORIZED:
      console.error(
        'Invalid admin credentials. Please ensure the you have provided the correct admin username and password.',
      )
      break
    case httpStatus.CONFLICT:
      console.error(
        `Error creating chainlink node. A node with the name: "${name}" already exists.`,
      )
      break
    default: {
      console.error(`Unhandled error. HTTP status: ${response.status}`)
      break
    }
  }
}

export const remove = async (name: string) => {
  const deleteNodeUrl = `${EXPLORER_BASE_URL}/api/v1/admin/nodes/${name}`
  const response = await fetch(deleteNodeUrl, {
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
    console.error(
      'Invalid admin credentials. Please ensure the you have provided the correct admin username and password.',
    )
  } else {
    console.error(`Unhandled error. HTTP status: ${response.status}`)
  }
}
