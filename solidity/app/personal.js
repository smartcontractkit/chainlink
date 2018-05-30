const Eth = require('ethjs')

let eth = new Eth(new Eth.HttpProvider('http://localhost:18545'))
let from = '0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f'
let defaults = {
  data: '',
  from: from
}

module.exports = {
  send: async function send (params) {
    await eth.sendTransaction(Object.assign(defaults, params))
  }
}
