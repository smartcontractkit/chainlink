import url from 'url'

export default (path, query = {}, options = {}) => {
  let formatOptions = {
    pathname: path,
    query
  }

  if (options.port) {
    formatOptions['port'] = options.port
    formatOptions['hostname'] = options.hostname
  }
  return url.format(formatOptions)
}
