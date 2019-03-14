import * as h from './support/helpers'

contract('LinkEx', () => {
  let contract: any

  beforeEach(async () => {
    contract = await h.deploy('LinkEx.sol')
  })

  describe('#currentRate', () => {
    context('after an initial deployment', () => {
      it('returns 0', async () => {
        const rate = await contract.currentRate()
        assert.equal(rate, 0)
      })
    })

    context('when requested in the same block as an update', () => {
      const expected = 3542157117
      const updated = 8616460799

      it('returns the historic rate', async () => {
        await contract.update(expected) // Set an initial rate
        await h.sendToEvm('miner_stop') // Stop mining blocks
        const txData = h.createTxData('update(uint256)', ['uint256'], [updated])

        // Sends an update to the price without increasing the block. We need to
        // use sendTransaction here otherwise Truffle will wait indefinitely for
        // the block to be mined before proceeding.
        h.eth.sendTransaction({
          from: h.defaultAccount,
          to: contract.address,
          data: txData
        })
        const expectedRate = await contract.currentRate()

        await h.sendToEvm('miner_start') // Start mining again
        assert.equal(expectedRate.toString(), expected.toString())

        // After a block has been mined, the rate is updated
        const updatedRate = await contract.currentRate()
        assert.equal(updatedRate.toString(), updated.toString())
      })
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
        assert.isNotTrue(await contract.authorizedNodes.call(h.oracleNode))
      })
    })

    context('when called by the owner', () => {
      it('adds the oracle', async () => {
        await contract.addOracle(h.oracleNode, {from: h.defaultAccount})
        assert.isTrue(await contract.authorizedNodes.call(h.oracleNode))
      })
    })
  })

  describe('#removeOracle', () => {
    beforeEach(async () => {
      assert.isNotTrue(await contract.authorizedNodes.call(h.oracleNode))
      await contract.addOracle(h.oracleNode, {from: h.defaultAccount})
      assert.isTrue(await contract.authorizedNodes.call(h.oracleNode))
    })

    context('when called by a stranger', () => {
      it('does not remove the oracle', async () => {
        await h.assertActionThrows(async () => {
          await contract.removeOracle(h.oracleNode, {from: h.stranger})
        })
        assert.isTrue(await contract.authorizedNodes.call(h.oracleNode))
      })
    })

    context('when called by the owner', () => {
      it('removes the oracle', async () => {
        await contract.removeOracle(h.oracleNode, {from: h.defaultAccount})
        assert.isNotTrue(await contract.authorizedNodes.call(h.oracleNode))
      })
    })

  })
})
