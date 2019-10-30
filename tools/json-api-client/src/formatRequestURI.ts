import url, { UrlObject } from 'url'

export interface Options {
  port?: string
  hostname?: string
}

export default (path: string, query = {}, options: Options = {}) => {
  const formatOptions: UrlObject = {
    pathname: path,
    query: query,
  }

  if (options.port) {
    formatOptions.port = options.port
    formatOptions.hostname = options.hostname
  }

  return url.format(formatOptions)
}
