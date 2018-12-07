import {resolve, join } from 'path'
import {
  deploy,
  linkContract,
  toWei
} from '../../../solidity/test/support/helpers'

contract('RunLog', () => {
  const arbitraryJobID =
        "0x0000000000000000000000000000000000000000000000000000000000000001"
  let link, logger, oc;

  beforeEach(async () => {
    link = await linkContract()
    oc = await deploy('Oracle.sol', link.address)
    logger = await deploy('RunLog.sol', link.address, oc.address, arbitraryJobID)
    await link.transfer(logger.address, toWei("1"));
  });

  it("has a limited public interface", async () => {
    let tx = await logger.request();
    assert.equal(4, tx.receipt.logs.length);
  });
});
