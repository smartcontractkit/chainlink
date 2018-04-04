const njsp = require("nano-json-stream-parser");
const request = require("xhr-request-promise");

const EthereumProvider = (url, intercept) => {
  intercept = intercept || (() => {});

  let api = {};
  let onResponse = {}; 
  let callbacks = {};
  let nextId = 0;
  let send;

  const makeSender = send => {
    const P = fn => (...args) => new Promise((resolve, reject) =>
      fn(...args.concat((err,res) => err ? reject(err) : resolve(res))));
    const sender = intercept => (method, params, callback) => {
      const intercepted = intercept(method, params, P(sender(() => {})));
      if (intercepted) {
        intercepted.then(response => callback(null, response));
      } else {
        send(method, params, callback);
      }
    }
    return sender(intercept);
  };

  const parseResponse = njsp(json => {
    onResponse[json.id] && onResponse[json.id](null, json.result);
  });

  const genPayload = (method, params) => ({
    jsonrpc: "2.0",
    id: ++nextId,
    method: method,
    params: params
  });

  api.on = (name, callback) => {
    callbacks[name] = callback;
  }

  if (/^ws/.test(url)) {
    const WebSocket = require("w"+"s");
    const ws = new WebSocket(url);
    api.send = makeSender((method, params, callback) => {
      const intercepted = intercept(method, params, P(send(() => {})));
      if (intercepted) {
        intercepted.then(response => callback(null, response));
      } else {
        const payload = genPayload(method, params);
        onResponse[payload.id] = callback;
        ws.send(JSON.stringify(payload));
      }
    });
    ws.on("message", parseResponse);
    ws.on("open", () => callbacks.connect && callbacks.connect(eth));
    ws.on("close", () => callbacks.disconnect && callbacks.disconnect());
    
  } else if (/^http/.test(url)) {
    api.send = makeSender((method, params, callback) => {
      request(url, {
        method: "POST",
        contentType: "application/json-rpc",
        body: JSON.stringify(genPayload(method,params))})
        .then(answer => {
          var resp = JSON.parse(answer);
          if (resp.error) {
            callback(resp.error.message);
          } else {
            callback(null, resp.result)
          }
        })
        .catch(err => callback("Couldn't connect to Ethereum node."));
    });

    setTimeout(() => {
      callbacks.connect && callbacks.connect();
    }, 1);

  } else {
    throw "IPC not supported yet.";
  }

  return api;
};

module.exports = EthereumProvider;
