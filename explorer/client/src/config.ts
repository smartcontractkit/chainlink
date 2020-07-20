/**
 * The envionment the explorer client application is running in
 */

export class Config {
  static baseUrl(env = process.env): string | undefined {
    return env.REACT_APP_EXPLORER_BASEURL
  }

  static gaId(env = process.env): string {
    return env.REACT_APP_EXPLORER_GA_ID ?? __REACT_APP_EXPLORER_GA_ID__ ?? ''
  }
}
