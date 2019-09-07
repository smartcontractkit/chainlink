import url from 'url'

export default (path, query = {}, options = {}) => {
  const formatOptions = {
    pathname: path,
    query: query,
  }

  if (options.port) {
    formatOptions['port'] = options.port
    formatOptions['hostname'] = options.hostname
  }
  return url.format(formatOptions)
}
