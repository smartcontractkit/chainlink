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
   * The origin of the client, used for CORS purposes
   */
  clientOrigin: string
  /**
   * The value of the secret used to sign cookies.
   *
   * For production usage, make sure this value is kept secret,
   * and is sufficiently secure.
   *
   * When used for development/testing purposes, it can be set
   * to some simple value like 'key1'.
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
    dev: !!env.EXPLORER_DEV,
    clientOrigin: env.EXPLORER_CLIENT_ORIGIN ?? '',
    cookieSecret: env.EXPLORER_COOKIE_SECRET,
    cookieExpirationMs: 86_400_000, // 1 day in ms
  }

  if (!conf.cookieSecret) {
    console.warn(
      'WARNING: Cookie secret is not set! Set via EXPLORER_COOKIE_SECRET',
    )
  }

  for (const [k, v] of Object.entries(conf)) {
    if (v == undefined) {
      throw Error(
        `Expected environment variable for ${k} to be set. Got "${v}".`,
      )
    }
  }

  return conf
}
