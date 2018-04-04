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
    fs.writeFile(path.join(tmpPath, 'foo-real.txt'), 'foo', function(err) {
      assert.ifError(err);
      rimraf(tmpPath, done);
    });
  });
};
