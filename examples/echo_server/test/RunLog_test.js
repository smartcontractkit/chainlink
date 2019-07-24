const LinkToken = artifacts.require('LinkToken')
const Oracle = artifacts.require('Oracle')
const RunLog = artifacts.require('RunLog')

contract('RunLog', () => {
  const arbitraryJobID =
    '0x0000000000000000000000000000000000000000000000000000000000000001'
  let link, logger, oc

  beforeEach(async () => {
    link = await LinkToken.new()
    oc = await Oracle.new(link.address)
    logger = await RunLog.new(link.address, oc.address, arbitraryJobID)
    await link.transfer(logger.address, web3.utils.toWei('1'))
  })

  it('has a limited public interface', async () => {
    let tx = await logger.request()
    assert.equal(4, tx.receipt.rawLogs.length)
  })
})
