export default (path) => {
  return `glob:${process.env.CHAINLINK_BASEURL || 'http://localhost'}${path}*`
}
