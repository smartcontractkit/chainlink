import mock from 'mock-fs'
import { RuntimeConfig, RuntimeConfigParser as Config } from './runtimeConfig'

describe('RuntimeConfig', () => {
  const strify = JSON.stringify

  function getSampleConf(): RuntimeConfig {
    return {
      chainId: 4,
      mnemonic:
        'raise clutch area heavy horn course filter farm deny solid finger sudden',
      infuraProjectId: 'fdf38d85d15e434e9b2ca152b7b1bc6f',
      etherscanAPIKey: 'US123YA2UIC73Q58PU6YZYBEXK99FWSJB3',
      gasPrice: 40000000000, // 40 gwei
      gasLimit: 8000000,
    }
  }

  afterEach(() => {
    mock.restore()
  })

  it('should throw on a missing .beltrc', () => {
    const conf = new Config('test-dir')
    mock({
      'test-dir': {},
    })

    expect(() => conf.load()).toThrowError('Could not load .beltrc')
  })

  // TODO: additional test cases for validation

  it('should load .beltrc successfully', () => {
    const conf = new Config('test-dir')
    mock({
      'test-dir': {
        '.beltrc': strify(getSampleConf()),
      },
    })

    expect(conf.load()).toStrictEqual({ ...getSampleConf() })
  })
})
