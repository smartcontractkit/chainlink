contract('Basic Truffle Consumer', () => {
  const jobSpecId = '0xbadc0de5'
  const Oracle = artifacts.require('Oracle')
  const Consumer = artifacts.require('RunLog')
  const LinkToken = artifacts.require('LinkToken')
  let consumer, linkToken, oracle

  beforeEach(async () => {
    linkToken = await LinkToken.new()
    oracle = await Oracle.new(linkToken.address)
    consumer = await Consumer.new(linkToken.address, oracle.address, jobSpecId)
    await linkToken.transfer(consumer.address, web3.utils.toWei('1', 'ether'))
  })

  describe('#requestEthereumPrice', () => {
    it('updates the bytes32 value', async () => {
      const answer = web3.utils.toHex('Hi Mom!')
      const tx = await consumer.request()
      const request = decodeRunRequest(tx.receipt.rawLogs[3])

      await oracle.fulfillOracleRequest(request.id,
        request.payment,
        request.callbackAddress,
        request.callbackFunc,
        request.expiration,
        answer)

      assert.equal(request.id, await consumer.requestId.call())
      assert.equal(web3.utils.toUtf8(answer), web3.utils.toUtf8(await consumer.response.call()))
    })
  })
})

const Ox = value => (value.slice(0, 2) !== '0x') ? `0x${value}` : value

const decodeRunRequest = log => {
  const wordSize = 64
  const payload = log.data.slice(2)
  return {
    requester: payload.slice(wordSize * 0, wordSize * 1),
    id: Ox(payload.slice(wordSize * 1, wordSize * 2)),
    payment: Ox(payload.slice(wordSize * 2, wordSize * 3)),
    callbackAddress: Ox(payload.slice(wordSize * 3, wordSize * 4).slice(24)),
    callbackFunc: Ox(payload.slice(wordSize * 4, wordSize * 5)),
    expiration: Ox(payload.slice(wordSize * 5, wordSize * 6)),
    dataVersion: Ox(payload.slice(wordSize * 6, wordSize * 7)),
    data: Ox(payload.slice(wordSize * 7))
  }
}
