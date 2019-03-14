import * as h from './support/helpers'

contract('LinkEx', () => {
  let contract: any

  beforeEach(async () => {
    contract = await h.deploy('LinkEx.sol')
  })

  describe('#currentRate', () => {
    it('returns 0', async () => {
      const rate = await contract.currentRate()
      assert.equal(rate, 0)
    })
  })

  describe('#updateRate', () => {
    it('returns last set value', async () => {
      await contract.update(8616460799)
      const historicRate = await contract.currentRate.call()
      assert.equal(historicRate.toString(), '8616460799')
    })
  })

  describe('#addOracle', () => {
    context('when called by a stranger', () => {
      it('reverts', async () => {
        await h.assertActionThrows(async () => {
          await contract.addOracle(h.oracleNode, {from: h.stranger})
        })
      })
    })

    context('when called by the owner', () => {
      it('adds the oracle', async () => {
        await contract.addOracle(h.oracleNode, {from: h.defaultAccount})
        assert.isTrue(await contract.authorizedNodes.call(h.oracleNode))
      })
    })
  })
})
