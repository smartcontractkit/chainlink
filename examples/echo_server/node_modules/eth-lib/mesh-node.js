// Creates a simple, HTTP-based 
const servify = require("servify");
const array = require("./array");

const server = (port, config) => {
  const stablePeerCount = config.stablePeerCount || 8;
  const requestPeerAmount = config.requestPeerAmount || 3;
  const pingPeerDelay = config.pingPeerDelay || 10000;

  let peers = [];
  let peerByUrl = {};
  const Peer = url => {
    if (peerByUrl[url]) {
      return peerByUrl[url];
    }
    let lastPing = 0;
    let lastPong = 0;
    const rpc = servify.at(url);
    const rpcProxy = new Proxy({
      get: function(target, method) {
        return function() {
          return rpc[method].apply(this, arguments).then(result => {
            lastPong = Date.now();
            return result;
          });
        }
      }
    });
    const tick = () => {
      if (Date.now() - lastPing > pingPeerDelay) {
        
      }
    };
    return {
      url: url,
      rpc: rpcProxy,
    };
  };

  const getRandomPeers = (amount) =>
    array.generate(amount, () => state.peers[Math.random() * state.peers.length | 0]);

  const findNewPeers = (neighborAmount, requestedAmount) =>
    Promise
      .all(getRandomPeers(neighborAmount).map(peer => peer.rpc.getRandomPeers(requestedAmount)))
      .then(newPeers => ;

  let state = {
    peers: config.peerUrls.map(Peer),
    isPeer: config.peerUrls.reduce((isPeer, url) => (isPeer[url] = 1, isPeer), {});
  };

  return servify.api(8097, {
    registerPeer: (url) => {
      if (!isPeer[url]) {
        return servify.at(url).ping().then(pong => {
          if (pong === "pong") {
            isPeer[url] = 1;
            peers.push(Peer(url));
          }
        });
      }
      return "OK";
    },
    getRandomPeers: amount => getRandomPeers(amount).map(peer => peer.url),
    ping: () => "pong"
  
  });
};

