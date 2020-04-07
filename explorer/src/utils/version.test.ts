import { fetchMeta } from './version'

describe('version tests', () => {
  it('should list the hash and branch', async () => {
    console.log(await fetchMeta())
  })
})
