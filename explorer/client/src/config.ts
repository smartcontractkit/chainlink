/**
 * The envionment the explorer client application is running in
 */

export class Config {
  static baseUrl(env = process.env): string | undefined {
    return env.REACT_APP_EXPLORER_BASEURL
  }
}
