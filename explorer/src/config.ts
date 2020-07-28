/**
 * The envionment the explorer application is running in
 */
export enum Environment {
  TEST,
  DEV,
  PROD,
}

export class Config {
  static port(env = process.env): number {
    return parseInt(env.EXPLORER_SERVER_PORT, 10) || 8080
  }

  static testPort(env = process.env): number {
    return parseInt(env.EXPLORER_TEST_SERVER_PORT, 10) || 8081
  }

  static env(env = process.env): Environment {
    switch (this.nodeEnv(env)) {
      case 'production':
        return Environment.PROD
      case 'test':
        return Environment.TEST
      default:
        return Environment.DEV
    }
  }

  static nodeEnv(env = process.env): string | undefined {
    return env.NODE_ENV
  }

  static clientOrigin(env = process.env): string {
    return env.EXPLORER_CLIENT_ORIGIN ?? ''
  }

  static cookieSecret(env = process.env): string | undefined {
    const cookieSecret =
      this.env() === Environment.DEV && !env.EXPLORER_COOKIE_SECRET
        ? 'secret-sauce-secret-sauce-secret-sauce-secret-sauce-secret-sauce'
        : env.EXPLORER_COOKIE_SECRET

    validateCookieSecret(cookieSecret)

    return cookieSecret
  }

  static cookieExpirationMs(): number {
    return 86_400_000 // 1 day in ms
  }

  static typeorm(env = process.env): string | Environment {
    return env.TYPEORM_NAME || this.nodeEnv() || 'development'
  }

  static composeMode(env = process.env): string | undefined {
    return env.COMPOSE_MODE
  }

  static baseUrl(env = process.env): string {
    return env.EXPLORER_BASE_URL || 'http://localhost:8080'
  }

  static adminUsername(env = process.env): string | undefined {
    return env.EXPLORER_ADMIN_USERNAME
  }

  static adminPassword(env = process.env): string | undefined {
    return env.EXPLORER_ADMIN_PASSWORD
  }

  static etherscanHost(env = process.env): string {
    return env.ETHERSCAN_HOST || 'ropsten.etherscan.io'
  }

  static gaId(env = process.env): string {
    return env.GA_ID
  }

  static setEnv(key: string, value: string | number) {
    Object.assign(process.env, {
      [key]: value,
    })
  }
}

/**
 * Assert that a cookie secret is at least 32 characters in length.
 *
 * @param secret The secret value to validate.
 */
function validateCookieSecret(secret?: string): asserts secret is string {
  if (!secret) {
    throw Error(
      'Cookie secret is not set! Set via environment variable EXPLORER_COOKIE_SECRET',
    )
  }

  if (secret.length < 32) {
    throw Error('Cookie secret must be at least 32 characters')
  }
}
