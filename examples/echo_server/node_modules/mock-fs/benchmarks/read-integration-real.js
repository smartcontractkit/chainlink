var assert = require('assert');
var fs = require('fs');
var path = require('path');

var rimraf = require('rimraf');

var tmpPath = '.tmp';

/**
 * Timed test.  This includes the setup and teardown as part of the overall
 * test time.
 * @param {function(Error)} done Callback.
 */
exports.test = function(done) {
  fs.mkdir(tmpPath, function(mkdirErr) {
    assert.ifError(mkdirErr);
    var tmpFile = path.join(tmpPath, 'foo-real.txt');
    fs.writeFile(tmpFile, 'foo', function(writeErr) {
      assert.ifError(writeErr);
      fs.readFile(tmpFile, 'utf8', function(readErr, str) {
        assert.ifError(readErr);
        assert.equal(str, 'foo');
        rimraf(tmpPath, done);
      });
    });
  });
};
