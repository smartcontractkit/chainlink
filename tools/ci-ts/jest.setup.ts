import { getArgs, waitForService } from './test-helpers/common'

const { CHAINLINK_URL, EXTERNAL_ADAPTER_URL } = getArgs([
  'CHAINLINK_URL',
  'EXTERNAL_ADAPTER_URL',
])

beforeAll(async () => {
  await waitForService(CHAINLINK_URL)
  await waitForService(EXTERNAL_ADAPTER_URL)
})
