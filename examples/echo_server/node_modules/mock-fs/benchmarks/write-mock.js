var fs = require('fs');

var mock = require('..');

/**
 * Test setup.  Not timed.
 */
exports.beforeEach = function() {
  mock();
};

/**
 * Timed test.
 * @param {function(Error)} done Callback.
 */
exports.test = function(done) {
  fs.writeFile('foo-mock.txt', 'foo', done);
};

/**
 * Test teardown.  Not timed.
 */
exports.afterEach = function() {
  mock.restore();
};
