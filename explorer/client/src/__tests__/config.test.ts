import { Config } from '../config'

describe('config', () => {
  const originalEnv = { ...process.env }
  beforeEach(() => {
    process.env = originalEnv
  })

  it('returns baseUrl key from the process env', () => {
    process.env.REACT_APP_EXPLORER_BASEURL = 'baseUrl'
    expect(Config.baseUrl()).toEqual('baseUrl')
  })
})
