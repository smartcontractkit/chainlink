/**
 * Application configuration for the explorer
 */
export interface ExplorerConfig {
  /**
   * The port to run the server on
   */
  port: number
  /**
   * Whether dev mode is enabled or not
   */
  dev: boolean
  /**
   * Whether production mode is enabled or not
   */
  prod: boolean
  /**
   * The origin of the client, used for CORS purposes
   */
  clientOrigin: string
  /**
   * The value of the secret used to sign cookies.
   * Must be at least 32 characters.
   *
   * For production usage, make sure this value is kept secret
   * and has sufficient entropy
   *
   */
  cookieSecret: string
  /**
   * The cookie expiration time in milliseconds
   */
  cookieExpirationMs: number
}

/**
 * Get application configuration for the explorer app
 */
export function getConfig(): ExplorerConfig {
  const { env } = process

  const conf: ExplorerConfig = {
    port: parseInt(env.EXPLORER_SERVER_PORT) || 8080,
    dev: env.NODE_ENV == 'development',
    prod: env.NODE_ENV === 'production',
    clientOrigin: env.EXPLORER_CLIENT_ORIGIN ?? '',
    cookieSecret: env.EXPLORER_COOKIE_SECRET,
    cookieExpirationMs: 86_400_000, // 1 day in ms
  }

  validateCookieSecret(conf.cookieSecret)

  for (const [k, v] of Object.entries(conf)) {
    if (v == undefined) {
      throw Error(
        `Expected environment variable for ${k} to be set. Got "${v}".`,
      )
    }
  }

  return conf
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
