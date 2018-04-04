// make sourcemaps work!
require('source-map-support/register')

var Provider = require("ganache-core/lib/provider");
var Server = require("ganache-core/lib/server");

// This interface exists so as not to cause breaking changes.
module.exports = {
  server: function(options) {
    return Server.create(options);
  },
  provider: function(options) {
    return new Provider(options);
  }
};
