import * as url from 'url'

export interface Options {
  port?: string
  hostname?: string
}

export default function(
  pathname: string,
  query: Record<string, string> = {},
  options: Options = {},
) {
  const formatOptions: url.UrlObject = { pathname, query }

  if (options.port) {
    formatOptions.port = options.port
    formatOptions.hostname = options.hostname
  }

  return url.format(formatOptions)
}
