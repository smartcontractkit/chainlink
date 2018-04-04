const Sha3 = require('js-sha3')

const hashLengths = [ 224, 256, 384, 512 ]

var hash = function (bitcount) {
  if (bitcount !== undefined && hashLengths.indexOf(bitcount) == -1)
    throw new Error('Unsupported hash length')
  this.content = []
  this.bitcount = bitcount ? 'keccak_' + bitcount : 'keccak_512'
}

hash.prototype.update = function (i) {
  if (Buffer.isBuffer(i))
    this.content.push(i)
  else if (typeof i === 'string')
    this.content.push(new Buffer(i))
  else
    throw new Error('Unsupported argument to update')
  return this
}

hash.prototype.digest = function (encoding) {
  var result = Sha3[this.bitcount](Buffer.concat(this.content))
  if (encoding === 'hex')
    return result
  else if (encoding === 'binary' || encoding === undefined)
    return new Buffer(result, 'hex').toString('binary')
  else
    throw new Error('Unsupported encoding for digest: ' + encoding)
}

module.exports = {
  SHA3Hash: hash
}
