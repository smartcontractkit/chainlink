function _toConsumableArray(arr) { if (Array.isArray(arr)) { for (var i = 0, arr2 = Array(arr.length); i < arr.length; i++) { arr2[i] = arr[i]; } return arr2; } else { return Array.from(arr); } }

var njsp = require("nano-json-stream-parser");
var request = require("xhr-request-promise");

var EthereumProvider = function EthereumProvider(url, intercept) {
  intercept = intercept || function () {};

  var api = {};
  var onResponse = {};
  var callbacks = {};
  var nextId = 0;
  var send = void 0;

  var makeSender = function makeSender(send) {
    var P = function P(fn) {
      return function () {
        for (var _len = arguments.length, args = Array(_len), _key = 0; _key < _len; _key++) {
          args[_key] = arguments[_key];
        }

        return new Promise(function (resolve, reject) {
          return fn.apply(undefined, _toConsumableArray(args.concat(function (err, res) {
            return err ? reject(err) : resolve(res);
          })));
        });
      };
    };
    var sender = function sender(intercept) {
      return function (method, params, callback) {
        var intercepted = intercept(method, params, P(sender(function () {})));
        if (intercepted) {
          intercepted.then(function (response) {
            return callback(null, response);
          });
        } else {
          send(method, params, callback);
        }
      };
    };
    return sender(intercept);
  };

  var parseResponse = njsp(function (json) {
    onResponse[json.id] && onResponse[json.id](null, json.result);
  });

  var genPayload = function genPayload(method, params) {
    return {
      jsonrpc: "2.0",
      id: ++nextId,
      method: method,
      params: params
    };
  };

  api.on = function (name, callback) {
    callbacks[name] = callback;
  };

  if (/^ws/.test(url)) {
    var WebSocket = require("w" + "s");
    var ws = new WebSocket(url);
    api.send = makeSender(function (method, params, callback) {
      var intercepted = intercept(method, params, P(send(function () {})));
      if (intercepted) {
        intercepted.then(function (response) {
          return callback(null, response);
        });
      } else {
        var payload = genPayload(method, params);
        onResponse[payload.id] = callback;
        ws.send(JSON.stringify(payload));
      }
    });
    ws.on("message", parseResponse);
    ws.on("open", function () {
      return callbacks.connect && callbacks.connect(eth);
    });
    ws.on("close", function () {
      return callbacks.disconnect && callbacks.disconnect();
    });
  } else if (/^http/.test(url)) {
    api.send = makeSender(function (method, params, callback) {
      request(url, {
        method: "POST",
        contentType: "application/json-rpc",
        body: JSON.stringify(genPayload(method, params)) }).then(function (answer) {
        var resp = JSON.parse(answer);
        if (resp.error) {
          callback(resp.error.message);
        } else {
          callback(null, resp.result);
        }
      }).catch(function (err) {
        return callback("Couldn't connect to Ethereum node.");
      });
    });

    setTimeout(function () {
      callbacks.connect && callbacks.connect();
    }, 1);
  } else {
    throw "IPC not supported yet.";
  }

  return api;
};

module.exports = EthereumProvider;