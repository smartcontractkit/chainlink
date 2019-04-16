import * as h from './support/helpers'

contract('LinkEx', () => {
  let contract: any

  // Need to mine some blocks so that the check in update doesn't
  // fail when subtracting 25 from block.number.
  before(async () => {
    h.mineBlocks(25)
  })

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

      beforeEach(async () => {
        await contract.addOracle(h.oracleNode, { from: h.defaultAccount })
        await contract.update(expected, { from: h.oracleNode })
      })

      it('returns the historic rate', async () => {
        await h.sendToEvm('miner_stop') // Stop mining blocks
        const txData = h.createTxData('update(uint256)', ['uint256'], [updated])

        // Sends an update to the price without increasing the block. We need to
        // use sendTransaction here otherwise Truffle will wait indefinitely for
        // the block to be mined before proceeding.
        h.eth.sendTransaction({
          data: txData,
          from: h.oracleNode,
          to: contract.address
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
    const expected = 8616460799
    const expected2 = 8616460814
    const expected3 = 8616460681
    // Round down and discard any decimals, just like Solidity
    const expectedAvg = Math.trunc((expected + expected2 + expected3) / 3)

    context('when called by a stranger', () => {
      it('reverts', async () => {
        await h.assertActionThrows(async () => {
          await contract.update(expected, { from: h.stranger })
        })
        const rate = await contract.currentRate()
        assert.equal(rate, 0)
      })
    })

    context('when called by an authorized node', () => {
      beforeEach(async () => {
        await contract.addOracle(h.oracleNode, { from: h.defaultAccount })
      })

      it('updates the rate', async () => {
        await contract.update(expected, { from: h.oracleNode })
        const historicRate = await contract.currentRate.call()
        assert.equal(historicRate.toString(), expected.toString())
      })
    })

    context('when updated recently by oracles', () => {
      beforeEach(async () => {
        await contract.addOracle(h.oracleNode, { from: h.defaultAccount })
        await contract.addOracle(h.oracleNode2, { from: h.defaultAccount })
        await contract.addOracle(h.oracleNode3, { from: h.defaultAccount })

        await contract.update(expected, { from: h.oracleNode })
        await contract.update(expected2, { from: h.oracleNode2 })
        await contract.update(expected3, { from: h.oracleNode3 })
      })

      it('has an expected aggregate value', async () => {
        const rate = await contract.currentRate()
        assert.equal(rate, expectedAvg)
      })

      context('after adding more oracles', () => {
        const expected4 = 8616460198
        const expected5 = 8616460756
        const newExpectedAvg = Math.trunc(
          (expected + expected2 + expected3 + expected4 + expected5) / 5
        )

        beforeEach(async () => {
          await contract.addOracle(h.accounts[8], { from: h.defaultAccount })
          await contract.addOracle(h.accounts[9], { from: h.defaultAccount })

          await contract.update(expected4, { from: h.accounts[8] })
          await contract.update(expected5, { from: h.accounts[9] })
        })

        it('the new oracles contribute to the average', async () => {
          const rate = await contract.currentRate()
          assert.equal(rate, newExpectedAvg)
        })
      })

      context('after removing an oracle', () => {
        const updated = 8616460198

        context('in the middle of the array', () => {
          const newExpectedAvg = Math.trunc((expected2 + expected3) / 2)

          beforeEach(async () => {
            await contract.removeOracle(h.oracleNode, {
              from: h.defaultAccount
            })
          })

          it('the removed oracles do not contribute to the average', async () => {
            const rate = await contract.currentRate()
            assert.equal(rate.toString(), newExpectedAvg)
          })

          it('the removed oracle cannot update the price', async () => {
            await h.assertActionThrows(async () => {
              await contract.update(updated, { from: h.oracleNode })
            })
          })

          context('then adding an oracle', () => {
            const newExpectedAvg2 = Math.trunc(
              (updated + expected2 + expected3) / 3
            )

            beforeEach(async () => {
              await contract.addOracle(h.accounts[8], {
                from: h.defaultAccount
              })
            })

            it('the average is as expected', async () => {
              await contract.update(updated, { from: h.accounts[8] })
              const rate = await contract.currentRate()
              assert.equal(rate.toString(), newExpectedAvg2)
            })
          })
        })

        context('at the end of the array', () => {
          const newExpectedAvg = Math.trunc((updated + expected) / 2)

          beforeEach(async () => {
            await contract.removeOracle(h.oracleNode3, {
              from: h.defaultAccount
            })
            await contract.update(updated, { from: h.oracleNode2 })
          })

          it('the removed oracles do not contribute to the average', async () => {
            const rate = await contract.currentRate()
            assert.equal(rate.toString(), newExpectedAvg)
          })

          it('the removed oracle cannot update the price', async () => {
            await h.assertActionThrows(async () => {
              await contract.update(updated, { from: h.oracleNode3 })
            })
          })

          context('then adding an oracle', () => {
            const newUpdated = 8616460357
            const newExpectedAvg2 = Math.trunc(
              (updated + expected + newUpdated) / 3
            )

            beforeEach(async () => {
              await contract.addOracle(h.accounts[8], {
                from: h.defaultAccount
              })
            })

            it('the average is as expected', async () => {
              await contract.update(newUpdated, { from: h.accounts[8] })
              const rate = await contract.currentRate()
              assert.equal(rate.toString(), newExpectedAvg2)
            })
          })
        })
      })

      context('when updated by an oracle after 25 blocks', () => {
        beforeEach(async () => {
          h.mineBlocks(25)
          await contract.update(expected, { from: h.oracleNode })
        })

        it('adjusts the current rate', async () => {
          const rate = await contract.currentRate()
          assert.equal(rate, expected)
        })
      })
    })
  })

  describe('#addOracle', () => {
    context('when called by a stranger', () => {
      it('reverts', async () => {
        await h.assertActionThrows(async () => {
          await contract.addOracle(h.oracleNode, { from: h.stranger })
        })
        assert.isNotTrue(await contract.authorizedNodes.call(h.oracleNode))
      })
    })

    context('when called by the owner', () => {
      it('adds the oracle', async () => {
        await contract.addOracle(h.oracleNode, { from: h.defaultAccount })
        assert.isTrue(await contract.authorizedNodes.call(h.oracleNode))
      })

      context('when adding an oracle twice', () => {
        beforeEach(async () => {
          await contract.addOracle(h.oracleNode, { from: h.defaultAccount })
        })

        it('rejects the addition of the oracle', async () => {
          await h.assertActionThrows(async () => {
            await contract.addOracle(h.oracleNode, { from: h.defaultAccount })
          })
        })
      })
    })
  })

  describe('#removeOracle', () => {
    beforeEach(async () => {
      assert.isNotTrue(await contract.authorizedNodes.call(h.oracleNode))
      await contract.addOracle(h.oracleNode, { from: h.defaultAccount })
      assert.isTrue(await contract.authorizedNodes.call(h.oracleNode))
    })

    context('when called by a stranger', () => {
      it('does not remove the oracle', async () => {
        await h.assertActionThrows(async () => {
          await contract.removeOracle(h.oracleNode, { from: h.stranger })
        })
        assert.isTrue(await contract.authorizedNodes.call(h.oracleNode))
      })
    })

    context('when called by the owner', () => {
      context('if there was only one oracle', () => {
        it('removes the oracle', async () => {
          await contract.removeOracle(h.oracleNode, { from: h.defaultAccount })
          assert.isNotTrue(await contract.authorizedNodes.call(h.oracleNode))
        })
      })

      context('if there are multiple oracles', () => {
        const badData = 9999999999
        const expected2 = 8616460814
        const expected3 = 8616460681
        const badAvg = Math.trunc((badData + expected2 + expected3) / 3)
        const goodAvg = Math.trunc((expected2 + expected3) / 2)

        beforeEach(async () => {
          await contract.update(badData, { from: h.oracleNode })
          await contract.addOracle(h.oracleNode2, { from: h.defaultAccount })
          await contract.addOracle(h.oracleNode3, { from: h.defaultAccount })
          await contract.update(expected2, { from: h.oracleNode2 })
          await contract.update(expected3, { from: h.oracleNode3 })
        })

        it('updates the average', async () => {
          const beforeRate = await contract.currentRate()
          assert.equal(beforeRate.toString(), badAvg)
          await contract.removeOracle(h.oracleNode, { from: h.defaultAccount })
          const afterRate = await contract.currentRate()
          assert.equal(afterRate.toString(), goodAvg)
        })
      })

      context('if the oracle does not exist', () => {
        it('reverts', async () => {
          await h.assertActionThrows(async () => {
            await contract.removeOracle(h.oracleNode3, {
              from: h.defaultAccount
            })
          })
        })
      })
    })
  })
})
