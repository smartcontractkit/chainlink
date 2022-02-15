import { extractBuildInfo } from './extractBuildInfo'

describe('extractBuildInfo', () => {
  const originalEnv = process.env

  describe('valid values', () => {
    beforeEach(() => {
      jest.resetModules()
      process.env = {
        ...originalEnv,
        CHAINLINK_VERSION: '1.0.0@6989a388ef26d981e771fec6710dc65bcc8fb5af',
      }
    })

    afterEach(() => {
      process.env = originalEnv
    })

    it('extracts the build info', () => {
      const { version, sha } = extractBuildInfo()

      expect(version).toEqual('1.0.0')
      expect(sha).toEqual('6989a388ef26d981e771fec6710dc65bcc8fb5af')
    })
  })

  describe('invalid format', () => {
    beforeEach(() => {
      jest.resetModules()
      process.env = {
        ...originalEnv,
        CHAINLINK_VERSION: '',
      }
    })

    afterEach(() => {
      process.env = originalEnv
    })

    it('has unknown values', () => {
      const { version, sha } = extractBuildInfo()

      expect(version).toEqual('unknown')
      expect(sha).toEqual('unknown')
    })
  })

  it('handles undefined env', () => {
    const { version, sha } = extractBuildInfo()

    expect(version).toEqual('unknown')
    expect(sha).toEqual('unknown')
  })
})
