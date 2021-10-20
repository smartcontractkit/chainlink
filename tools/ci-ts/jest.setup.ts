import { getEnvVars, waitForService } from './test-helpers/common'

const {
  CLIENT_NODE_URL,
  CLIENT_NODE_2_URL,
  EXTERNAL_ADAPTER_URL,
  EXTERNAL_ADAPTER_2_URL,
} = getEnvVars([
  'CLIENT_NODE_URL',
  'CLIENT_NODE_2_URL',
  'EXTERNAL_ADAPTER_URL',
  'EXTERNAL_ADAPTER_2_URL',
])

beforeAll(async () => {
  await Promise.all([
    waitForService(CLIENT_NODE_URL),
    waitForService(CLIENT_NODE_2_URL),
    waitForService(EXTERNAL_ADAPTER_URL),
    waitForService(EXTERNAL_ADAPTER_2_URL),
  ])
})
