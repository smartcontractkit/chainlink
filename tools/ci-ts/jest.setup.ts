import { getArgs, waitForService } from './test-helpers/common'

const { CHAINLINK_URL } = getArgs(['CHAINLINK_URL'])

beforeAll(async () => {
  await waitForService(CHAINLINK_URL)
})
