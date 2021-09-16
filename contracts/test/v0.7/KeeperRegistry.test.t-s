const KeeperRegistry = artifacts.require('KeeperRegistry')
const UpkeepMock = artifacts.require('UpkeepMock')
const UpkeepReverter = artifacts.require('UpkeepReverter')
const { LinkToken } = require('@chainlink/contracts/truffle/v0.4/LinkToken')
const { MockV3Aggregator } = require('@chainlink/contracts/truffle/v0.6/MockV3Aggregator')
const { BN, constants, ether, expectEvent, expectRevert, time } = require('@openzeppelin/test-helpers')

// -----------------------------------------------------------------------------------------------
// DEV: these *should* match the perform/check gas overhead values in the contract and on the node
const PERFORM_GAS_OVERHEAD = new BN('90000')
const CHECK_GAS_OVERHEAD = new BN('170000')
// -----------------------------------------------------------------------------------------------

contract('KeeperRegistry', (accounts) => {
  const owner = accounts[0]
  const keeper1 = accounts[1]
  const keeper2 = accounts[2]
  const keeper3 = accounts[3]
  const nonkeeper = accounts[4]
  const admin = accounts[5]
  const payee1 = accounts[6]
  const payee2 = accounts[7]
  const payee3 = accounts[8]
  const keepers = [keeper1, keeper2, keeper3]
  const payees = [payee1, payee2, payee3]
  const linkEth = new BN(300000000)
  const gasWei = new BN(100)
  const linkDivisibility = new BN("1000000000000000000")
  const executeGas = new BN('100000')
  const paymentPremiumBase = new BN('1000000000')
  const paymentPremiumPPB =  new BN('250000000')
  const blockCountPerTurn = new BN(3)
  const emptyBytes = '0x00'
  const zeroAddress = constants.ZERO_ADDRESS
  const extraGas = new BN('250000')
  const registryGasOverhead = new BN('80000')
  const stalenessSeconds = new BN(43820)
  const gasCeilingMultiplier = new BN(1)
  const maxCheckGas = new BN(20000000)
  const fallbackGasPrice = new BN(200)
  const fallbackLinkPrice = new BN(200000000)
  let linkToken, linkEthFeed, gasPriceFeed, registry, mock, id

  linkForGas = (upkeepGasSpent) => {
    const gasSpent = registryGasOverhead.add(new BN(upkeepGasSpent))
    const base = gasWei.mul(gasSpent).mul(linkDivisibility).div(linkEth)
    const premium = base.mul(paymentPremiumPPB).div(paymentPremiumBase)
    return base.add(premium)
  }

  beforeEach(async () => {
    LinkToken.setProvider(web3.currentProvider)
    MockV3Aggregator.setProvider(web3.currentProvider)
    linkToken = await LinkToken.new({ from: owner })
    gasPriceFeed = await MockV3Aggregator.new(0, gasWei, { from: owner })
    linkEthFeed = await MockV3Aggregator.new(9, linkEth, { from: owner })
    registry = await KeeperRegistry.new(
      linkToken.address,
      linkEthFeed.address,
      gasPriceFeed.address,
      paymentPremiumPPB,
      blockCountPerTurn,
      maxCheckGas,
      stalenessSeconds,
      gasCeilingMultiplier,
      fallbackGasPrice,
      fallbackLinkPrice,
      { from: owner }
    )
    mock = await UpkeepMock.new()
    await linkToken.transfer(keeper1, ether('1000'), { from: owner })
    await linkToken.transfer(keeper2, ether('1000'), { from: owner })
    await linkToken.transfer(keeper3, ether('1000'), { from: owner })

    await registry.setKeepers(keepers, payees, {from: owner})
    const { receipt } = await registry.registerUpkeep(
      mock.address,
      executeGas,
      admin,
      emptyBytes,
      { from: owner }
    )
    id = receipt.logs[0].args.id
  })

  describe('#setKeepers', () => {
    const IGNORE_ADDRESS = "0xFFfFfFffFFfffFFfFFfFFFFFffFFFffffFfFFFfF";
    it('reverts when not called by the owner', async () => {
      await expectRevert(
        registry.setKeepers([], [], {from: keeper1}),
        "Only callable by owner"
      )
    })

    it('reverts when adding the same keeper twice', async () => {
      await expectRevert(
        registry.setKeepers([keeper1, keeper1], [payee1, payee1], {from: owner}),
        "cannot add keeper twice"
      )
    })

    it('reverts with different numbers of keepers/payees', async () => {
      await expectRevert(
        registry.setKeepers([keeper1, keeper2], [payee1], {from: owner}),
        "address lists not the same length"
      )
      await expectRevert(
        registry.setKeepers([keeper1], [payee1, payee2], {from: owner}),
        "address lists not the same length"
      )
    })

    it('reverts if the payee is the zero address', async () => {
      await expectRevert(
        registry.setKeepers([keeper1, keeper2], [payee1, "0x0000000000000000000000000000000000000000"], {from: owner}),
        "cannot set payee to the zero address"
      )
    })

    it('emits events for every keeper added and removed', async () => {
      const oldKeepers = [keeper1, keeper2]
      const oldPayees = [payee1, payee2]
      await registry.setKeepers(oldKeepers, oldPayees, {from: owner})
      assert.deepEqual(oldKeepers, await registry.getKeeperList())

      // remove keepers
      const newKeepers = [keeper2, keeper3]
      const newPayees = [payee2, payee3]
      const { receipt } = await registry.setKeepers(newKeepers, newPayees, {from: owner})
      assert.deepEqual(newKeepers, await registry.getKeeperList())

      expectEvent(receipt, 'KeepersUpdated', {
        keepers: newKeepers,
        payees: newPayees
      })
    })

    it('updates the keeper to inactive when removed', async () => {
      await registry.setKeepers(keepers, payees, {from: owner})
      await registry.setKeepers([keeper1, keeper3], [payee1, payee3], {from: owner})
      const added = await registry.getKeeperInfo(keeper1)
      assert.isTrue(added.active)
      const removed = await registry.getKeeperInfo(keeper2)
      assert.isFalse(removed.active)
    })

    it('does not change the payee if IGNORE_ADDRESS is used as payee', async () => {
      const oldKeepers = [keeper1, keeper2]
      const oldPayees = [payee1, payee2]
      await registry.setKeepers(oldKeepers, oldPayees, {from: owner})
      assert.deepEqual(oldKeepers, await registry.getKeeperList())

      const newKeepers = [keeper2, keeper3]
      const newPayees = [IGNORE_ADDRESS, payee3]
      const { receipt } = await registry.setKeepers(newKeepers, newPayees, {from: owner})
      assert.deepEqual(newKeepers, await registry.getKeeperList())

      const ignored = await registry.getKeeperInfo(keeper2)
      assert.equal(payee2, ignored.payee)
      assert.equal(true, ignored.active)

      expectEvent(receipt, 'KeepersUpdated', {
        keepers: newKeepers,
        payees: newPayees
      })
    })

    it('reverts if the owner changes the payee', async () => {
      await registry.setKeepers(keepers, payees, {from: owner})
      await expectRevert(
        registry.setKeepers(keepers, [payee1, payee2, owner], {from: owner}),
        "cannot change payee"
      )
    })
  })

  describe('#registerUpkeep', () => {
    it('reverts if the target is not a contract', async () => {
      await expectRevert(
        registry.registerUpkeep(
          zeroAddress,
          executeGas,
          admin,
          emptyBytes,
          { from: owner }
        ),
        'target is not a contract'
      )
    })

    it('reverts if called by a non-owner', async () => {
      await expectRevert(
        registry.registerUpkeep(
          mock.address,
          executeGas,
          admin,
          emptyBytes,
          { from: keeper1 }
        ),
        'Only callable by owner or registrar'
      )
    })

    it('reverts if execute gas is too low', async () => {
      await expectRevert(
        registry.registerUpkeep(
          mock.address,
          2299,
          admin,
          emptyBytes,
          { from: owner }
        ),
        'min gas is 2300'
      )
    })


    it('reverts if execute gas is too high', async () => {
      await expectRevert(
        registry.registerUpkeep(
          mock.address,
          5000001,
          admin,
          emptyBytes,
          { from: owner }
        ),
        'max gas is 5000000'
      )
    })

    it('creates a record of the registration', async () => {
      const { receipt } = await registry.registerUpkeep(
        mock.address,
        executeGas,
        admin,
        emptyBytes,
        { from: owner }
      )
      id = receipt.logs[0].args.id
      expectEvent(receipt, 'UpkeepRegistered', {
        id: id,
        executeGas: executeGas
      })
      const registration = await registry.getUpkeep(id)
      assert.equal(mock.address, registration.target)
      assert.equal(0, registration.balance)
      assert.equal(emptyBytes, registration.checkData)
      assert.equal(0xffffffffffffffff, registration.maxValidBlocknumber)
    })
  })

  describe('#addFunds', () => {
    const amount = ether('1')

    beforeEach(async () => {
      await linkToken.approve(registry.address, ether('100'), { from: keeper1 })
    })

    it('reverts if the registration does not exist', async () => {
      await expectRevert(
        registry.addFunds(id + 1, amount, { from: keeper1 }),
        'upkeep must be active'
      )
    })

    it('adds to the balance of the registration', async () => {
      await registry.addFunds(id, amount, { from: keeper1 })
      const registration = await registry.getUpkeep(id)
      assert.isTrue(amount.eq(registration.balance))
    })

    it('emits a log', async () => {
      const { receipt } = await registry.addFunds(id, amount, { from: keeper1 })

      expectEvent(receipt, 'FundsAdded', {
        id: id,
        from: keeper1,
        amount: amount
      })
    })

    it('reverts if the upkeep is canceled', async () => {
      await registry.cancelUpkeep(id, { from: admin })
      await expectRevert(
        registry.addFunds(id, amount, { from: keeper1 }),
        "upkeep must be active",
      )
    })
  })

  describe('#checkUpkeep', () => {
    it('reverts if the upkeep is not funded', async () => {
      await mock.setCanPerform(true)
      await mock.setCanCheck(true)
      await expectRevert(
        registry.checkUpkeep.call(id, keeper1, {from: zeroAddress}),
        "insufficient funds"
      )
    })

    context('when the registration is funded', () => {
      beforeEach(async () => {
        await linkToken.approve(registry.address, ether('100'), { from: keeper1 })
        await registry.addFunds(id, ether('100'), { from: keeper1 })
      })

      it('reverts if executed', async () => {
        await mock.setCanPerform(true)
        await mock.setCanCheck(true)
        await expectRevert(
          registry.checkUpkeep(id, keeper1),
          'only for simulated backend'
        )
      })

      it('reverts if the specified keeper is not valid', async () => {
        await mock.setCanPerform(true)
        await mock.setCanCheck(true)
        await expectRevert(
          registry.checkUpkeep(id, owner),
          'only for simulated backend'
        )
      })

      context('and upkeep is not needed', () => {
        beforeEach(async () => {
          await mock.setCanCheck(false)
        })

        it('reverts', async () => {
          await expectRevert(
            registry.checkUpkeep.call(id, keeper1, {from: zeroAddress}),
            'upkeep not needed'
          )
        })
      })

      context('and the upkeep check fails', () => {
        beforeEach(async () => {
          const reverter = await UpkeepReverter.new()
          const { receipt } = await registry.registerUpkeep(
            reverter.address,
            2500000,
            admin,
            emptyBytes,
            { from: owner }
          )
          id = receipt.logs[0].args.id
          await linkToken.approve(registry.address, ether('100'), { from: keeper1 })
          await registry.addFunds(id, ether('100'), { from: keeper1 })
        })

        it('reverts', async () => {
          await expectRevert(
            registry.checkUpkeep.call(id, keeper1, {from: zeroAddress}),
            'call to check target failed'
          )
        })
      })

      context('and performing the upkeep simulation fails', () => {
        beforeEach(async () => {
          await mock.setCanCheck(true)
          await mock.setCanPerform(false)
        })

        it('reverts', async () => {
          await expectRevert(
            registry.checkUpkeep.call(id, keeper1, {from: zeroAddress}),
            'call to perform upkeep failed'
          )
        })
      })

      context('and upkeep check and perform simulations succeeds', () => {
        beforeEach(async () => {
          await mock.setCanCheck(true)
          await mock.setCanPerform(true)
        })

        context('and the registry is paused', () => {
          beforeEach(async () => {
            await registry.pause({from: owner})
          })

          it('reverts', async () => {
            await expectRevert(
              registry.checkUpkeep.call(id, keeper1, {from: zeroAddress}),
              'Pausable: paused'
            )

            await registry.unpause({from: owner})

            await registry.checkUpkeep.call(id, keeper1, {from: zeroAddress})
          })
        })

        it('returns true with pricing info if the target can execute', async () => {
          const newGasMultiplier = new BN(10)
          await registry.setConfig(
            paymentPremiumPPB,
            blockCountPerTurn,
            maxCheckGas,
            stalenessSeconds,
            newGasMultiplier,
            fallbackGasPrice,
            fallbackLinkPrice,
            { from: owner }
          )
          const response = await registry.checkUpkeep.call(id, keeper1, {from: zeroAddress})
          assert.isTrue(response.gasLimit.eq(executeGas))
          assert.isTrue(response.linkEth.eq(linkEth))
          assert.isTrue(response.adjustedGasWei.eq(gasWei.mul(newGasMultiplier)))
          assert.isTrue(response.maxLinkPayment.eq(linkForGas(executeGas).mul(newGasMultiplier)))
        })

        it('has a large enough gas overhead to cover upkeeps that use all their gas', async () => {
          await mock.setCheckGasToBurn(maxCheckGas)
          await mock.setPerformGasToBurn(executeGas)
          const gas = maxCheckGas.add(executeGas).add(PERFORM_GAS_OVERHEAD).add(CHECK_GAS_OVERHEAD)
          await registry.checkUpkeep.call(id, keeper1, { from: zeroAddress, gas: gas })
        })
      })
    })
  })

  describe('#performUpkeep', () => {
    let _lastKeeper = keeper1
    async function getPerformPaymentAmount(){
      _lastKeeper = _lastKeeper === keeper1 ? keeper2 : keeper1
      const before = (await registry.getKeeperInfo(_lastKeeper)).balance
      await registry.performUpkeep(id, "0x", {from: _lastKeeper})
      const after = (await registry.getKeeperInfo(_lastKeeper)).balance
      const difference = after.sub(before)
      return difference
    }

    it('reverts if the registration is not funded', async () => {
      await expectRevert(
        registry.performUpkeep(id, "0x", { from: keeper2 }),
        'insufficient funds'
      )
    })

    context('when the registration is funded', () => {
      beforeEach(async () => {
        await linkToken.approve(registry.address, ether('100'), { from: owner })
        await registry.addFunds(id, ether('100'), { from: owner })
      })

      it('does not revert if the target cannot execute', async () => {
        const mockResponse = await mock.checkUpkeep.call("0x", { from: zeroAddress })
        assert.isFalse(mockResponse.callable)

        await registry.performUpkeep(id, "0x", { from: keeper3 })
      })

      it('returns false if the target cannot execute', async () => {
        const mockResponse = await mock.checkUpkeep.call("0x", { from: zeroAddress })
        assert.isFalse(mockResponse.callable)

        assert.isFalse(await registry.performUpkeep.call(id, "0x", { from: keeper1 }))
      })

      it('returns true if called', async () => {
        await mock.setCanPerform(true)

        const response = await registry.performUpkeep.call(id, "0x", {from: keeper1})
        assert.isTrue(response)
      })

      it('reverts if not enough gas supplied', async () => {
        await mock.setCanPerform(true)

        await expectRevert.unspecified(
          registry.performUpkeep(id, "0x", { from: keeper1, gas: new BN('120000') })
        )
      })

      it('executes the data passed to the registry', async () => {
        await mock.setCanPerform(true)

        const performData = "0xc0ffeec0ffee"
        const tx = await registry.performUpkeep(id, performData, { from: keeper1, gas: extraGas })
        expectEvent(tx.receipt, 'UpkeepPerformed', {
          success: true,
          from: keeper1,
          performData: performData
        })
      })

      it('updates payment balances', async () => {
        const keeperBefore = await registry.getKeeperInfo(keeper1)
        const registrationBefore = await registry.getUpkeep(id)
        const keeperLinkBefore = await linkToken.balanceOf(keeper1)
        const registryLinkBefore = await linkToken.balanceOf(registry.address)

        //// Do the thing
        await registry.performUpkeep(id, "0x", { from: keeper1 })

        const keeperAfter = await registry.getKeeperInfo(keeper1)
        const registrationAfter = await registry.getUpkeep(id)
        const keeperLinkAfter = await linkToken.balanceOf(keeper1)
        const registryLinkAfter = await linkToken.balanceOf(registry.address)

        assert.isTrue(keeperAfter.balance.gt(keeperBefore.balance))
        assert.isTrue(registrationBefore.balance.gt(registrationAfter.balance))
        assert.isTrue(keeperLinkAfter.eq(keeperLinkBefore))
        assert.isTrue(registryLinkBefore.eq(registryLinkAfter))
      })

      it('only pays for gas used', async () => {
        const before = (await registry.getKeeperInfo(keeper1)).balance
        const { receipt } = await registry.performUpkeep(id, "0x", { from: keeper1 })
        const after = (await registry.getKeeperInfo(keeper1)).balance

        const max = linkForGas(executeGas)
        const totalTx = linkForGas(receipt.gasUsed)
        const difference = after.sub(before)
        assert.isTrue(max.gt(totalTx))
        assert.isTrue(totalTx.gt(difference))
        assert.isTrue(linkForGas(5700).lt(difference)) // exact number is flaky
        assert.isTrue(linkForGas(6000).gt(difference)) // instead test a range
      })

      it('only pays at a rate up to the gas ceiling', async () => {
        const multiplier = new BN(10)
        const gasPrice = new BN('1000000000') // 10M x the gas feed's rate
        await registry.setConfig(
          paymentPremiumPPB,
          blockCountPerTurn,
          maxCheckGas,
          stalenessSeconds,
          multiplier,
          fallbackGasPrice,
          fallbackLinkPrice,
          { from: owner }
        )

        const before = (await registry.getKeeperInfo(keeper1)).balance
        const { receipt } = await registry.performUpkeep(id, "0x", { from: keeper1, gasPrice })
        const after = (await registry.getKeeperInfo(keeper1)).balance

        const max = linkForGas(executeGas).mul(multiplier)
        const totalTx = linkForGas(receipt.gasUsed).mul(multiplier)
        const difference = after.sub(before)
        assert.isTrue(max.gt(totalTx))
        assert.isTrue(totalTx.gt(difference))
        assert.isTrue(linkForGas(5700).mul(multiplier).lt(difference))
        assert.isTrue(linkForGas(6000).mul(multiplier).gt(difference))
      })

      it('only pays as much as the node spent', async () => {
        const multiplier = new BN(10)
        const gasPrice = new BN(200) // 2X the gas feed's rate
        const effectiveMultiplier = new BN(2)
        await registry.setConfig(
          paymentPremiumPPB,
          blockCountPerTurn,
          maxCheckGas,
          stalenessSeconds,
          multiplier,
          fallbackGasPrice,
          fallbackLinkPrice,
          { from: owner }
        )

        const before = (await registry.getKeeperInfo(keeper1)).balance
        const { receipt } = await registry.performUpkeep(id, "0x", { from: keeper1, gasPrice })
        const after = (await registry.getKeeperInfo(keeper1)).balance

        const max = linkForGas(executeGas).mul(effectiveMultiplier)
        const totalTx = linkForGas(receipt.gasUsed).mul(effectiveMultiplier)
        const difference = after.sub(before)
        assert.isTrue(max.gt(totalTx))
        assert.isTrue(totalTx.gt(difference))
        assert.isTrue(linkForGas(5700).mul(effectiveMultiplier).lt(difference))
        assert.isTrue(linkForGas(6000).mul(effectiveMultiplier).gt(difference))
      })

      it('pays the caller even if the target function fails', async () => {
        const { receipt } = await registry.registerUpkeep(
          mock.address,
          executeGas,
          admin,
          emptyBytes,
          { from: owner }
        )
        const id = receipt.logs[0].args.id
        await linkToken.approve(registry.address, ether('100'), { from: owner })
        await registry.addFunds(id, ether('100'), { from: owner })
        const keeperBalanceBefore = (await registry.getKeeperInfo(keeper1)).balance

        // Do the thing
        const tx = await registry.performUpkeep(id, "0x", { from: keeper1 })

        const keeperBalanceAfter = (await registry.getKeeperInfo(keeper1)).balance
        assert.isTrue(keeperBalanceAfter.gt(keeperBalanceBefore))
      })

      it('reverts if called by a non-keeper', async () => {
        await expectRevert(
          registry.performUpkeep(id, "0x", { from: nonkeeper }),
          'only active keepers'
        )
      })

      it('reverts if the upkeep has been canceled', async () => {
        await mock.setCanPerform(true)

        await registry.cancelUpkeep(id, { from: owner })

        await expectRevert(
          registry.performUpkeep(id, "0x", { from: keeper1 }),
          'invalid upkeep id'
        )
      })

      it("uses the fallback gas price if the feed price is stale", async () => {
        const normalAmount = await getPerformPaymentAmount()
        const roundId = 99
        const answer = 100
        const updatedAt = 946684800 // New Years 2000 🥳
        const startedAt = 946684799
        await gasPriceFeed.updateRoundData(roundId, answer, updatedAt, startedAt, {from: owner})
        const amountWithStaleFeed = await getPerformPaymentAmount()
        assert.isTrue(normalAmount.lt(amountWithStaleFeed))
      })

      it("uses the fallback gas price if the feed price is non-sensical", async () => {
        const normalAmount = await getPerformPaymentAmount()
        const roundId = 99
        const updatedAt = Math.floor(Date.now() / 1000)
        const startedAt = 946684799
        await gasPriceFeed.updateRoundData(roundId, -100, updatedAt, startedAt, {from: owner})
        const amountWithNegativeFeed = await getPerformPaymentAmount()
        await gasPriceFeed.updateRoundData(roundId, 0, updatedAt, startedAt, {from: owner})
        const amountWithZeroFeed = await getPerformPaymentAmount()
        assert.isTrue(normalAmount.lt(amountWithNegativeFeed))
        assert.isTrue(normalAmount.lt(amountWithZeroFeed))
      })

      it("uses the fallback if the link price feed is stale", async () => {
        const normalAmount = await getPerformPaymentAmount()
        const roundId = 99
        const answer = 100
        const updatedAt = 946684800 // New Years 2000 🥳
        const startedAt = 946684799
        await linkEthFeed.updateRoundData(roundId, answer, updatedAt, startedAt, {from: owner})
        const amountWithStaleFeed = await getPerformPaymentAmount()
        assert.isTrue(normalAmount.lt(amountWithStaleFeed))
      })

      it("uses the fallback link price if the feed price is non-sensical", async () => {
        const normalAmount = await getPerformPaymentAmount()
        const roundId = 99
        const updatedAt = Math.floor(Date.now() / 1000)
        const startedAt = 946684799
        await linkEthFeed.updateRoundData(roundId, -100, updatedAt, startedAt, {from: owner})
        const amountWithNegativeFeed = await getPerformPaymentAmount()
        await linkEthFeed.updateRoundData(roundId, 0, updatedAt, startedAt, {from: owner})
        const amountWithZeroFeed = await getPerformPaymentAmount()
        assert.isTrue(normalAmount.lt(amountWithNegativeFeed))
        assert.isTrue(normalAmount.lt(amountWithZeroFeed))
      })

      it('reverts if the same caller calls twice in a row', async () => {
        await registry.performUpkeep(id, "0x", { from: keeper1 }),
        await expectRevert(
          registry.performUpkeep(id, "0x", { from: keeper1 }),
          'keepers must take turns'
        )
        await registry.performUpkeep(id, "0x", { from: keeper2 })
        await expectRevert(
          registry.performUpkeep(id, "0x", { from: keeper2 }),
          'keepers must take turns'
        )
        await registry.performUpkeep(id, "0x", { from: keeper1 })
      })

      it('has a large enough gas overhead to cover upkeeps that use all their gas', async () => {
        await mock.setPerformGasToBurn(executeGas)
        await mock.setCanPerform(true)
        const gas = executeGas.add(PERFORM_GAS_OVERHEAD)
        const performData = "0xc0ffeec0ffee"
        const { receipt } = await registry.performUpkeep(id, performData, { from: keeper1, gas: gas })
        expectEvent(receipt, 'UpkeepPerformed', {
          success: true,
          from: keeper1,
          performData: performData
        })
      })
    })
  })

  describe('#withdrawFunds', () => {
    beforeEach(async () => {
      await linkToken.approve(registry.address, ether('100'), { from: keeper1 })
      await registry.addFunds(id, ether('1'), { from: keeper1 })
    })

    it('reverts if called by anyone but the admin', async () => {
      await expectRevert(
        registry.withdrawFunds(id + 1, payee1, { from: owner }),
        'only callable by admin'
      )
    })

    it('reverts if called on an uncanceled upkeep', async () => {
      await expectRevert(
        registry.withdrawFunds(id, payee1, { from: admin }),
        'upkeep must be canceled'
      )
    })

    it('reverts if called with the 0 address', async () => {
      await expectRevert(
        registry.withdrawFunds(id, zeroAddress, { from: admin }),
        'cannot send to zero address'
      )
    })

    describe("after the registration is cancelled", () => {
      beforeEach(async () => {
        await registry.cancelUpkeep(id, { from: owner })
      })

      it('moves the funds out and updates the balance', async () => {
        const payee1Before = await linkToken.balanceOf(payee1)
        const registryBefore = await linkToken.balanceOf(registry.address)

        let registration = await registry.getUpkeep(id)
        assert.isTrue(ether('1').eq(registration.balance))

        await registry.withdrawFunds(id, payee1, { from: admin })

        const payee1After = await linkToken.balanceOf(payee1)
        const registryAfter = await linkToken.balanceOf(registry.address)

        assert.isTrue(payee1Before.add(ether('1')).eq(payee1After))
        assert.isTrue(registryBefore.sub(ether('1')).eq(registryAfter))

        registration = await registry.getUpkeep(id)
        assert.equal(0, registration.balance)
      })
    })
  })

  describe('#cancelUpkeep', () => {
    it('reverts if the ID is not valid', async () => {
      await expectRevert(
        registry.cancelUpkeep(id + 1, { from: owner }),
        'too late to cancel upkeep'
      )
    })

    it('reverts if called by a non-owner/non-admin', async () => {
      await expectRevert(
        registry.cancelUpkeep(id, { from: keeper1 }),
        'only owner or admin'
      )
    })

    describe("when called by the owner", async () => {
      it('sets the registration to invalid immediately', async () => {
        const { receipt } = await registry.cancelUpkeep(id, { from: owner })

        const registration = await registry.getUpkeep(id)
        assert.equal(registration.maxValidBlocknumber.toNumber(), receipt.blockNumber)
      })

      it('emits an event', async () => {
        const { receipt } = await registry.cancelUpkeep(id, { from: owner })

        expectEvent(receipt, 'UpkeepCanceled', {
          id: id,
          atBlockHeight: new BN(receipt.blockNumber)
        })
      })

      it('updates the canceled registrations list', async () => {
        let canceled = await registry.getCanceledUpkeepList.call()
        assert.deepEqual([], canceled)

        await registry.cancelUpkeep(id, { from: owner })

        canceled = await registry.getCanceledUpkeepList.call()
        assert.deepEqual([id], canceled)
      })

      it('immediately prevents upkeep', async () => {
        await registry.cancelUpkeep(id, { from: owner })

        await expectRevert(
          registry.performUpkeep(id, "0x", { from: keeper2 }),
          'invalid upkeep id'
        )
      })

      it('does not revert if reverts if called multiple times', async () => {
        await registry.cancelUpkeep(id, { from: owner })
        await expectRevert(
          registry.cancelUpkeep(id, { from: owner }),
          'too late to cancel upkeep'
        )
      })

      describe("when called by the owner when the admin has just canceled", () => {
        let oldExpiration

        beforeEach(async () => {
          await registry.cancelUpkeep(id, { from: admin })
          const registration = await registry.getUpkeep(id)
          oldExpiration = registration.maxValidBlocknumber
        })

        it('allows the owner to cancel it more quickly', async () => {
          await registry.cancelUpkeep(id, { from: owner })

          const registration = await registry.getUpkeep(id)
          const newExpiration = registration.maxValidBlocknumber
          assert.isTrue(newExpiration.lt(oldExpiration))
        })
      })
    })

    describe("when called by the admin", async () => {
      const delay = 50

      it('sets the registration to invalid in 50 blocks', async () => {
        const { receipt } = await registry.cancelUpkeep(id, { from: admin })
        const registration = await registry.getUpkeep(id)
        assert.isFalse(registration.maxValidBlocknumber.eq(receipt.blockNumber + 50))
      })

      it('emits an event', async () => {
        const { receipt } = await registry.cancelUpkeep(id, { from: admin })
        expectEvent(receipt, 'UpkeepCanceled', {
          id: id,
          atBlockHeight: new BN(receipt.blockNumber + delay)
        })
      })

      it('updates the canceled registrations list', async () => {
        let canceled = await registry.getCanceledUpkeepList.call()
        assert.deepEqual([], canceled)

        await registry.cancelUpkeep(id, { from: admin })

        canceled = await registry.getCanceledUpkeepList.call()
        assert.deepEqual([id], canceled)
      })

      it('immediately prevents upkeep', async () => {
        await linkToken.approve(registry.address, ether('100'), { from: owner })
        await registry.addFunds(id, ether('100'), { from: owner })
        await registry.cancelUpkeep(id, { from: admin })
        await registry.performUpkeep(id, "0x", { from: keeper2 }) // still works

        for (let i = 0; i < delay; i++) {
          await time.advanceBlock()
        }

        await expectRevert(
          registry.performUpkeep(id, "0x", { from: keeper2 }),
          'invalid upkeep id'
        )
      })

      it('reverts if called again by the admin', async () => {
        await registry.cancelUpkeep(id, { from: admin })

        await expectRevert(
          registry.cancelUpkeep(id, { from: admin }),
          'too late to cancel upkeep'
        )
      })

      it('does not revert or double add the cancellaition record if called by the owner immediately after', async () => {
        await registry.cancelUpkeep(id, { from: admin })

        await registry.cancelUpkeep(id, { from: owner })

        let canceled = await registry.getCanceledUpkeepList.call()
        assert.deepEqual([id], canceled)
      })

      it('reverts if called by the owner after the timeout', async () => {
        await registry.cancelUpkeep(id, { from: admin })

        for (let i = 0; i < delay; i++) {
          await time.advanceBlock()
        }

        await expectRevert(
          registry.cancelUpkeep(id, { from: owner }),
          'too late to cancel upkeep'
        )
      })
    })
  })

  describe('#withdrawPayment', () => {
    beforeEach(async () => {
      await linkToken.approve(registry.address, ether('100'), { from: owner })
      await registry.addFunds(id, ether('100'), { from: owner })
      await registry.performUpkeep(id, "0x", { from: keeper1 })
    })

    it("reverts if called by anyone but the payee", async () => {
      await expectRevert(
        registry.withdrawPayment(keeper1, nonkeeper, { from: payee2 }),
        "only callable by payee"
      )
    })

    it('reverts if called with the 0 address', async () => {
      await expectRevert(
        registry.withdrawPayment(keeper1, zeroAddress, { from: payee2 }),
        'cannot send to zero address'
      )
    })

    it("updates the balances", async () => {
      const to = nonkeeper
      const keeperBefore = (await registry.getKeeperInfo(keeper1)).balance
      const registrationBefore = (await registry.getUpkeep(id)).balance
      const toLinkBefore = await linkToken.balanceOf(to)
      const registryLinkBefore = await linkToken.balanceOf(registry.address)

      //// Do the thing
      await registry.withdrawPayment(keeper1, nonkeeper, { from: payee1 })

      const keeperAfter = (await registry.getKeeperInfo(keeper1)).balance
      const registrationAfter = (await registry.getUpkeep(id)).balance
      const toLinkAfter = await linkToken.balanceOf(to)
      const registryLinkAfter = await linkToken.balanceOf(registry.address)

      assert.isTrue(keeperAfter.eq(new BN(0)))
      assert.isTrue(registrationBefore.eq(registrationAfter))
      assert.isTrue(toLinkBefore.add(keeperBefore).eq(toLinkAfter))
      assert.isTrue(registryLinkBefore.sub(keeperBefore).eq(registryLinkAfter))
    })

    it("emits a log announcing the withdrawal", async () => {
      const balance = (await registry.getKeeperInfo(keeper1)).balance
      const { receipt } = await registry.withdrawPayment(
        keeper1, nonkeeper, { from: payee1 }
      )

      expectEvent(receipt, 'PaymentWithdrawn', {
        keeper: keeper1,
        amount: balance,
        to: nonkeeper,
        payee: payee1,
      })
    })
  })

  describe('#transferPayeeship', () => {
    it("reverts when called by anyone but the current payee", async () => {
      await expectRevert(
        registry.transferPayeeship(keeper1, payee2, { from: payee2 }),
        "only callable by payee"
      )
    })

    it("reverts when transferring to self", async () => {
      await expectRevert(
        registry.transferPayeeship(keeper1, payee1, { from: payee1 }),
        "cannot transfer to self"
      )
    })

    it("does not change the payee", async () => {
      await registry.transferPayeeship(keeper1, payee2, { from: payee1 })

      const info = await registry.getKeeperInfo(keeper1)
      assert.equal(payee1, info.payee)
    })

    it("emits an event announcing the new payee", async () => {
      const { receipt } = await registry.transferPayeeship(keeper1, payee2, { from: payee1 })

      expectEvent(receipt, 'PayeeshipTransferRequested', {
        keeper: keeper1,
        from: payee1,
        to: payee2,
      })
    })

    it("does not emit an event when called with the same proposal", async () => {
      await registry.transferPayeeship(keeper1, payee2, { from: payee1 })

      const { receipt } = await registry.transferPayeeship(keeper1, payee2, { from: payee1 })

      assert.equal(0, receipt.logs.length)
    })
  })

  describe('#acceptPayeeship', () => {
    beforeEach(async () => {
      await registry.transferPayeeship(keeper1, payee2, { from: payee1 })
    })

    it("reverts when called by anyone but the proposed payee", async () => {
      await expectRevert(
        registry.acceptPayeeship(keeper1, { from: payee1 }),
        "only callable by proposed payee"
      )
    })

    it("emits an event announcing the new payee", async () => {
      const { receipt } = await registry.acceptPayeeship(keeper1, { from: payee2 })

      expectEvent(receipt, 'PayeeshipTransferred', {
        keeper: keeper1,
        from: payee1,
        to: payee2,
      })
    })

    it("does change the payee", async () => {
      await registry.acceptPayeeship(keeper1, { from: payee2 })

      const info = await registry.getKeeperInfo(keeper1)
      assert.equal(payee2, info.payee)
    })
  })

  describe('#setConfig', () => {
    const payment = new BN(1)
    const checks = new BN(2)
    const staleness = new BN(3)
    const ceiling = new BN(10)
    const maxGas = new BN(4)
    const fbGasEth = new BN(5)
    const fbLinkEth = new BN(6)

    it("reverts when called by anyone but the proposed owner", async () => {
      await expectRevert(
        registry.setConfig(
          payment,
          checks,
          maxGas,
          staleness,
          gasCeilingMultiplier,
          fbGasEth,
          fbLinkEth,
          { from: payee1 }
        ),
        "Only callable by owner"
      )
    })

    it("updates the config", async () => {
      const old = await registry.getConfig()
      assert.isTrue(paymentPremiumPPB.eq(old.paymentPremiumPPB))
      assert.isTrue(blockCountPerTurn.eq(old.blockCountPerTurn))
      assert.isTrue(stalenessSeconds.eq(old.stalenessSeconds))
      assert.isTrue(gasCeilingMultiplier.eq(old.gasCeilingMultiplier))

      await registry.setConfig(
        payment,
        checks,
        maxGas,
        staleness,
        ceiling,
        fbGasEth,
        fbLinkEth,
        { from: owner }
      )

      const updated = await registry.getConfig()
      assert.isTrue(updated.paymentPremiumPPB.eq(payment))
      assert.isTrue(updated.blockCountPerTurn.eq(checks))
      assert.isTrue(updated.stalenessSeconds.eq(staleness))
      assert.isTrue(updated.gasCeilingMultiplier.eq(ceiling))
      assert.isTrue(updated.checkGasLimit.eq(maxGas))
      assert.isTrue(updated.fallbackGasPrice.eq(fbGasEth))
      assert.isTrue(updated.fallbackLinkPrice.eq(fbLinkEth))
    })

    it("emits an event", async () => {
      const { receipt } = await registry.setConfig(
        payment,
        checks,
        maxGas,
        staleness,
        ceiling,
        fbGasEth,
        fbLinkEth,
        { from: owner }
      )
      expectEvent(receipt, 'ConfigSet', {
        paymentPremiumPPB: payment,
        blockCountPerTurn: checks,
        checkGasLimit: maxGas,
        stalenessSeconds: staleness,
        gasCeilingMultiplier: ceiling,
        fallbackGasPrice: fbGasEth,
        fallbackLinkPrice: fbLinkEth,
      })
    })
  })

  describe('#onTokenTransfer', () => {
    const amount = ether('1')

    it("reverts if not called by the LINK token", async () => {
      const data = web3.eth.abi.encodeParameter('uint256', id.toNumber().toString())

      await expectRevert(
        registry.onTokenTransfer(keeper1, amount, data, {from: keeper1}),
        "only callable through LINK"
      )
    })

    it("reverts if not called with more or less than 32 bytes", async () => {
      const longData = web3.eth.abi.encodeParameters(['uint256', 'uint256'], ['33', '34'])
      const shortData = "0x12345678"

      await expectRevert.unspecified(
        linkToken.transferAndCall(registry.address, amount, longData, {from: owner})
      )
      await expectRevert.unspecified(
        linkToken.transferAndCall(registry.address, amount, shortData, {from: owner})
      )
    })

    it('reverts if the upkeep is canceled', async () => {
      await registry.cancelUpkeep(id, { from: admin })
      await expectRevert(
        registry.addFunds(id, amount, { from: keeper1 }),
        "upkeep must be active",
      )
    })

    it('updates the funds of the job id passed', async () => {
      const data = web3.eth.abi.encodeParameter('uint256', id.toNumber().toString())

      const before = (await registry.getUpkeep(id)).balance
      await linkToken.transferAndCall(registry.address, amount, data, { from: owner})
      const after = (await registry.getUpkeep(id)).balance

      assert.isTrue(before.add(amount).eq(after))
    })
  })

  describe('#recoverFunds', () => {
    const sent = ether('7')

    beforeEach(async () => {
      await linkToken.approve(registry.address, ether('100'), { from: keeper1 })

      // add funds to upkeep 1 and perform and withdraw some payment
      let { receipt } = await registry.registerUpkeep(
        mock.address,
        executeGas,
        admin,
        emptyBytes,
        { from: owner }
      )
      const id1 = receipt.logs[0].args.id
      await registry.addFunds(id1, ether('5'), { from: keeper1 })
      await registry.performUpkeep(id1, "0x", { from: keeper1 })
      await registry.performUpkeep(id1, "0x", { from: keeper2 })
      await registry.performUpkeep(id1, "0x", { from: keeper3 })
      await registry.withdrawPayment(keeper1, nonkeeper, { from: payee1 })

      // transfer funds directly to the registry
      await linkToken.transfer(registry.address, sent, { from: keeper1 })

      // add funds to upkeep 2 and perform and withdraw some payment
      const tx2 = await registry.registerUpkeep(
        mock.address,
        executeGas,
        admin,
        emptyBytes,
        { from: owner }
      )
      const id2 = tx2.receipt.logs[0].args.id
      await registry.addFunds(id2, ether('5'), { from: keeper1 })
      await registry.performUpkeep(id2, "0x", { from: keeper1 })
      await registry.performUpkeep(id2, "0x", { from: keeper2 })
      await registry.performUpkeep(id2, "0x", { from: keeper3 })
      await registry.withdrawPayment(keeper2, nonkeeper, { from: payee2 })

      // transfer funds using onTokenTransfer
      const data = web3.eth.abi.encodeParameter('uint256', id2.toNumber().toString())
      await linkToken.transferAndCall(registry.address, ether('1'), data, { from: owner})

      // remove a keeper
      await registry.setKeepers([keeper1, keeper2], [payee1, payee2], {from: owner})

      // withdraw some funds
      await registry.cancelUpkeep(id1, { from: owner })
      await registry.withdrawFunds(id1, admin, { from: admin })
    })

    it('reverts if not called by owner', async () => {
      await expectRevert(
        registry.recoverFunds({from: keeper1}),
        "Only callable by owner"
      )
    })

    it('allows any funds that have been accidentally transfered to be moved', async () => {
      const balanceBefore = await linkToken.balanceOf(registry.address)

      await linkToken.balanceOf(registry.address)

      const tx = await registry.recoverFunds({from: owner})
      const balanceAfter = await linkToken.balanceOf(registry.address)
      assert.isTrue(balanceBefore.eq(balanceAfter.add(sent)))
    })
  })

  describe('#pause', () => {
    it('reverts if called by a non-owner', async () => {
      await expectRevert(
        registry.pause({from: keeper1}),
        "Only callable by owner"
      )
    })

    it('marks the contract as paused', async () => {
      assert.isFalse(await registry.paused())

      await registry.pause({from: owner})

      assert.isTrue(await registry.paused())
    })
  })

  describe('#unpause', () => {
    beforeEach(async () => {
      await registry.pause({from: owner})
    })

    it('reverts if called by a non-owner', async () => {
      await expectRevert(
        registry.unpause({from: keeper1}),
        "Only callable by owner"
      )
    })

    it('marks the contract as not paused', async () => {
      assert.isTrue(await registry.paused())

      await registry.unpause({from: owner})

      assert.isFalse(await registry.paused())
    })
  })

  describe('#checkUpkeep / #performUpkeep', () => {
    const performData = "0xc0ffeec0ffee"
    const multiplier = new BN(10)
    const callGasPrice = 1

    it('uses the same minimum balance calculation', async () => {
      await registry.setConfig(
        paymentPremiumPPB,
        blockCountPerTurn,
        maxCheckGas,
        stalenessSeconds,
        multiplier,
        fallbackGasPrice,
        fallbackLinkPrice,
        { from: owner }
      )
      await linkToken.approve(registry.address, ether('100'), { from: owner })

      // max payment is .75 eth for this config - this spread will yield some eligible and some ineligible
      const balances = ['0', '0.01', '0.1', '0.4', '0.7', '0.8', '1', '2', '10']
      let revertCount = 0

      for (let idx = 0; idx < balances.length; idx++) {
        const balance = balances[idx];
        const { receipt } = await registry.registerUpkeep(
          mock.address,
          executeGas,
          admin,
          emptyBytes,
          { from: owner }
        )
        const upkeepID = receipt.logs[0].args.id
        await mock.setCanCheck(true)
        await mock.setCanPerform(true)
        await registry.addFunds(upkeepID, ether(balance), { from: owner })

        try {
          // try checkUpkeep
          await registry.checkUpkeep.call(upkeepID, keeper1, {from: zeroAddress, gasPrice: callGasPrice})
        } catch (err) {
        // if checkUpkeep reverts, we expect performUpkeep to revert as well
          revertCount++;
          await expectRevert(
            registry.performUpkeep(upkeepID, performData, { from: keeper1, gas: extraGas }),
            'insufficient funds'
          )
          continue
        }
        // if checkUpkeep succeeds, we expect performUpkeep to succeed as well
        try {
          await registry.performUpkeep(upkeepID, performData, { from: keeper1, gas: extraGas })
        } catch (err) {
          assert(false, `expected performUpkeep to have succeeded with balance ${balance} ETH, but it did not. err: ${err}`)
        }
      }

      // make sure _both_ scenarios are covered - future-proofs the test against contract / config changes
      assert.isTrue(
        revertCount > 0 && revertCount < balances.length,
        `expected 0 < revertCount < ${balances.length}, but revertCount was ${revertCount}`
      )
    })
  })

  describe('#getMinBalanceForUpkeep / #checkUpkeep', () => {
    it('calculates the minimum balance appropriately', async () => {
      const oneWei = new BN('1')
      await linkToken.approve(registry.address, ether('100'), { from: keeper1 })
      await mock.setCanCheck(true)
      await mock.setCanPerform(true)
      const minBalance = await registry.getMinBalanceForUpkeep(id)
      const tooLow = minBalance.sub(oneWei)
      await registry.addFunds(id, tooLow, { from: keeper1 })
      await expectRevert(
        registry.checkUpkeep.call(id, keeper1, {from: zeroAddress}),
        'insufficient funds'
      )
      await registry.addFunds(id, oneWei, { from: keeper1 })
      await registry.checkUpkeep.call(id, keeper1, {from: zeroAddress})
    })
  })
})
