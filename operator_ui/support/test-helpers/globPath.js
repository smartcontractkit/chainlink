export default (path) => {
  return `glob:http://localhost${process.env.CHAINLINK_PORT || ''}${path}*`
}
