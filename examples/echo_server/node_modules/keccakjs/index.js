try {
  module.exports = require('sha3').SHA3Hash
} catch (err) {
  module.exports = require('./browser')
}
