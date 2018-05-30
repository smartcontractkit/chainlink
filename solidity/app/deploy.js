let compile = require('./compile.js')
let personal = require('./personal.js')
let utils = require('./utils.js')
let wallet = require('./wallet.js')

module.exports = async function deploy (filename) {
  let bytecode = compile(filename)

  await personal.send({to: wallet.address, value: utils.toWei(1)})
  let txHash = await wallet.send({
    gas: 2000000,
    data: '0x' + bytecode + '0000000000000000000000004b274dfcd56656742A55ad54549b3770c392aA87'
  })
  return utils.getTxReceipt(txHash).then((receipt) => {
    let address = receipt.contractAddress
    return address
  })
}
