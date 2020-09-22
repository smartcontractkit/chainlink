import mock from 'mock-fs'
import { App, load } from './config'

describe('config.load', () => {
  const strify = JSON.stringify

  function getDefaultConf(): App {
    return {
      contractsDir: 'src',
      artifactsDir: 'abi',
      contractAbstractionDir: '.',
      useDockerisedSolc: true,
      compilerSettings: {
        versions: {
          'v0.4': '0.4.24',
          'v0.5': '0.5.0',
          'v0.6': '0.6.2',
        },
      },
      publicVersions: ['0.4.24', '0.5.0'],
    }
  }

  afterEach(() => {
    mock.restore()
  })

  it('should throw on a missing config', () => {
    mock({
      src: '',
    })

    expect(() => load('./doesnte')).toThrowError('Could not load config')
  })

  it('should throw on a missing contracts directory', () => {
    mock({
      conf: strify(getDefaultConf()),
    })

    expect(() => load('./conf')).toThrowError(
      'Expected value of config.contractsDir to be a directory',
    )
  })

  it('should throw on a non-string artifacts directory value', () => {
    mock({
      src: {},
      conf: strify({ ...getDefaultConf(), artifactsDir: 5 }),
    })

    expect(() => load('./conf')).toThrowError(
      'Expected value of config.artifactsDir to be a string',
    )
  })

  it('should throw on a non-boolean useDockerisedSolc value', () => {
    mock({
      src: {},
      conf: strify({ ...getDefaultConf(), useDockerisedSolc: '' }),
    })

    expect(() => load('./conf')).toThrowError(
      'Expected value of config.useDockerisedSolc to be a boolean',
    )
  })

  it('should throw on an invalid contractAbstractionDir value', () => {
    mock({
      src: {},
      conf: strify({ ...getDefaultConf(), contractAbstractionDir: 5 }),
    })

    expect(() => load('./conf')).toThrowError(
      'Expected value of config.contractAbstractionDir to be a string',
    )
  })

  it('should throw on a non-valid compilerSettings value', () => {
    mock({
      src: {},
      conf: strify({ ...getDefaultConf(), compilerSettings: undefined }),
    })

    expect(() => load('./conf')).toThrowError(
      'Expected value of config.compilerSettings to be an object',
    )
  })

  it('should throw on a non-valid compilerSettings value', () => {
    mock({
      src: {},
      conf: strify({ ...getDefaultConf(), compilerSettings: {} }),
    })

    expect(() => load('./conf')).toThrowError(
      'Expected value of config.compilerSettings.versions to be a dictionary',
    )
  })

  it('should load a config correctly', () => {
    mock({
      src: {},
      conf: strify(getDefaultConf()),
    })

    expect(load('./conf')).toStrictEqual({ ...getDefaultConf() })
  })
})
