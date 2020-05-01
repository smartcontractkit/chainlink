import { Config, Environment } from '../config'
import { randomBytes } from 'crypto'

describe('config', () => {
  const originalEnv = { ...process.env }
  beforeEach(() => {
    process.env = originalEnv
  })
  afterAll(() => {
    process.env = originalEnv
  })

  it('returns port key from the process env', () => {
    process.env.EXPLORER_SERVER_PORT = '3000'
    expect(Config.port()).toEqual(3000)
  })

  it('returns default port key from the process env', () => {
    process.env.EXPLORER_SERVER_PORT = undefined
    expect(Config.port()).toEqual(8080)
  })

  it('returns development Environment key from the process env', () => {
    process.env.NODE_ENV = 'development'
    expect(Config.env()).toEqual(Environment.DEV)
  })

  it('returns production Environment key from the process env', () => {
    process.env.NODE_ENV = 'production'
    expect(Config.env()).toEqual(Environment.PROD)
  })

  it('returns Test Environment key from the process env', () => {
    process.env.NODE_ENV = 'test'
    expect(Config.env()).toEqual(Environment.TEST)
  })

  it('returns nodeEnv key from the process env', () => {
    process.env.NODE_ENV = 'test'
    expect(Config.nodeEnv()).toEqual('test')
  })

  it('returns clientOrigin key from the process env', () => {
    process.env.EXPLORER_CLIENT_ORIGIN = 'clientOrigin'
    expect(Config.clientOrigin()).toEqual('clientOrigin')
  })

  it('returns default clientOrigin key from the process env', () => {
    process.env.EXPLORER_CLIENT_ORIGIN = undefined
    expect(Config.clientOrigin()).toEqual('')
  })

  it('returns cookieSecret key from the process env', () => {
    const secret = randomBytes(32).toString('hex')
    process.env.EXPLORER_COOKIE_SECRET = secret
    expect(Config.cookieSecret()).toEqual(secret)
  })

  it('returns default dev cookieSecret key from the process env', () => {
    process.env.EXPLORER_COOKIE_SECRET = undefined
    process.env.NODE_ENV = 'development'
    expect(Config.cookieSecret()).toEqual(
      'secret-sauce-secret-sauce-secret-sauce-secret-sauce-secret-sauce',
    )
  })

  it('returns cookieExpirationMs', () => {
    expect(Config.cookieExpirationMs()).toEqual(86_400_000)
  })

  it('returns typeOrmName key from the process env', () => {
    process.env.TYPEORM_NAME = 'typeOrmName'
    expect(Config.typeorm()).toEqual('typeOrmName')
  })

  it('returns nodeEnv key when TYPEORM_NAME is not defined', () => {
    process.env.TYPEORM_NAME = undefined
    process.env.NODE_ENV = 'test'
    expect(Config.typeorm()).toEqual('test')
  })

  it('returns default key when TYPEORM_NAME and nodeEnv is not defined', () => {
    process.env.TYPEORM_NAME = undefined
    process.env.NODE_ENV = undefined
    expect(Config.typeorm()).toEqual('development')
  })

  it('returns composeMode key from the process env', () => {
    process.env.COMPOSE_MODE = 'composeMode'
    expect(Config.composeMode()).toEqual('composeMode')
  })

  it('returns baseUrl key from the process env', () => {
    process.env.EXPLORER_BASE_URL = 'baseUrl'
    expect(Config.baseUrl()).toEqual('baseUrl')
  })

  it('returns default baseUrl key from the process env', () => {
    process.env.EXPLORER_BASE_URL = undefined
    expect(Config.baseUrl()).toEqual('http://localhost:8080')
  })

  it('returns adminUsername key from the process env', () => {
    process.env.EXPLORER_ADMIN_USERNAME = 'adminUsername'
    expect(Config.adminUsername()).toEqual('adminUsername')
  })

  it('returns adminPassword key from the process env', () => {
    process.env.EXPLORER_ADMIN_PASSWORD = 'adminPassword'
    expect(Config.adminPassword()).toEqual('adminPassword')
  })

  it('returns etherscanHost key from the process env', () => {
    process.env.ETHERSCAN_HOST = 'etherscan'
    expect(Config.etherscanHost()).toEqual('etherscan')
  })

  it('returns default etherscanHost key from the process env', () => {
    process.env.ETHERSCAN_HOST = undefined
    expect(Config.etherscanHost()).toEqual('ropsten.etherscan.io')
  })

  it('sets the process env', () => {
    Config.setEnv('KEY', 'value')
    expect(process.env.KEY).toEqual('value')
  })
})
