var SHA3 = require('sha3')
var BSHA3 = require('../index.js')
var tape = require('tape')
var crypto = require('crypto')

tape('basic test', function (t) {
  t.plan(2)
    // Generate 512-bit digest.
  var d = new SHA3.SHA3Hash()
  d.update('foo')

  // Generate 512-bit digest.
  var bd = new BSHA3.SHA3Hash()
  bd.update('foo')
  t.equal(d.digest('hex'), bd.digest('hex'))

  // Generate 224-bit digest.
  var d = new SHA3.SHA3Hash(224)
  d.update('foo')

  // Generate 224-bit digest.
  var bd = new BSHA3.SHA3Hash(224)
  bd.update('foo')
  t.equal(d.digest('hex'), bd.digest('hex'))
})

tape('encoding', function (t) {
  t.plan(1)

  // Generate 224-bit digest.
  var d = new SHA3.SHA3Hash(224)
  d.update('foo')

  // Generate 224-bit digest.
  var bd = new BSHA3.SHA3Hash(224)
  bd.update('foo')
  t.equal(d.digest().toString(), bd.digest().toString())
})

tape('random test', function (t) {
  t.plan(10)

  for (var i = 0; i < 10; i++) {
    var data = crypto.randomBytes(32)
      // Generate 512-bit digest.
    var d = new SHA3.SHA3Hash()
    d.update(data)

    // Generate 512-bit digest.
    var bd = new BSHA3.SHA3Hash()
    bd.update(data)
    t.equal(d.digest('hex'), bd.digest('hex'))
  }
})
