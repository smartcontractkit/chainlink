export default path => {
  return `glob:${process.env.CHAINLINK_PORT || ''}${path}*`
}
