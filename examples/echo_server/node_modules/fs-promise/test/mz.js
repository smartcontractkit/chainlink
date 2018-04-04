// Tests from mz.  License MIT
// https://github.com/normalize/mz
var assert = require('assert')

describe('fs', function () {
  var fs = require('../')

  it('.stat()', function (done) {
    fs.stat(__filename).then(function (stats) {
      assert.equal(typeof stats.size, 'number')
      done()
    }).catch(done)
  })

  it('.mkdtemp()', function (done) {
    if (!require('fs').mkdtemp) this.skip()
    fs.mkdtemp('/tmp/').then(function (folder) {
      fs.rmdirSync(folder)
      done()
    }).catch(done)
  })

  it('.statSync()', function () {
    var stats = fs.statSync(__filename)
    assert.equal(typeof stats.size, 'number')
  })

  it('.exists()', function (done) {
    fs.exists(__filename).then(function (exists) {
      assert(exists)
      done()
    }).catch(done)
  })

  it('.existsSync()', function () {
    var exists = fs.existsSync(__filename)
    assert(exists)
  })

  describe('callback support', function () {
    it('.stat()', function (done) {
      fs.stat(__filename, function (err, stats) {
        assert(!err)
        assert.equal(typeof stats.size, 'number')
        done()
      })
    })

    it('.exists()', function (done) {
      fs.exists(__filename, function (err, exists) {
        assert(!err)
        assert(exists)
        done()
      })
    })
  })
})
