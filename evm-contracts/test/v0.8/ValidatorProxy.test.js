const { assert, web3 } = require('hardhat')
const { constants, expectRevert } = require('@openzeppelin/test-helpers')
const expectEvent = require('@openzeppelin/test-helpers/src/expectEvent')

describe('ValidatorProxy', () => {
  let accounts
  let owner, aggregator, validator
  let ValidatorProxyArtifact
  let validatorProxy

  beforeEach(async () => {
    ValidatorProxyArtifact = artifacts.require('ValidatorProxy')
    accounts = await web3.eth.getAccounts()
    owner = accounts[0]
    aggregator = accounts[1]
    validator = accounts[2]
    validatorProxy = await ValidatorProxyArtifact.new(aggregator, validator, {
      from: owner,
    })
  })

  describe('#constructor', () => {
    it('should set the aggregator addresses correctly', async () => {
      const response = await validatorProxy.getAggregators()
      assert.equal(response.current, aggregator)
      assert.equal(response.hasProposal, false)
      assert.equal(response.proposed, constants.ZERO_ADDRESS)
    })

    it('should set the validator addresses conrrectly', async () => {
      const response = await validatorProxy.getValidators()
      assert.equal(response.current, validator)
      assert.equal(response.hasProposal, false)
      assert.equal(response.proposed, constants.ZERO_ADDRESS)
    })

    it('should set the owner correctly', async () => {
      const response = await validatorProxy.owner()
      assert.equal(response, owner)
    })
  })

  describe('#proposeNewAggregator', () => {
    let newAggregator
    before(async () => {
      newAggregator = accounts[3]
    })

    it('should only be called by the owner', async () => {
      const stranger = accounts[4]
      await expectRevert(
        validatorProxy.proposeNewAggregator(newAggregator, {
          from: stranger,
        }),
        'Only callable by owner',
      )
    })

    describe('success', () => {
      let receipt

      beforeEach(async () => {
        receipt = await validatorProxy.proposeNewAggregator(newAggregator, {
          from: owner,
        })
      })

      it('should emit an event', async () => {
        await expectEvent(receipt, 'AggregatorProposed', {
          aggregator: newAggregator,
        })
      })

      it('should set the correct address and hasProposal is true', async () => {
        const response = await validatorProxy.getAggregators()
        assert.equal(response.current, aggregator)
        assert.equal(response.hasProposal, true)
        assert.equal(response.proposed, newAggregator)
      })

      it('should set a zero address and hasProposal is false', async () => {
        receipt = await validatorProxy.proposeNewAggregator(
          constants.ZERO_ADDRESS,
          {
            from: owner,
          },
        )
        const response = await validatorProxy.getAggregators()
        assert.equal(response.current, aggregator)
        assert.equal(response.hasProposal, false)
        assert.equal(response.proposed, constants.ZERO_ADDRESS)
      })
    })
  })

  describe('#upgradeAggregator', () => {
    describe('failure', () => {
      it('should only be called by the owner', async () => {
        const stranger = accounts[4]
        await expectRevert(
          validatorProxy.upgradeAggregator({
            from: stranger,
          }),
          'Only callable by owner',
        )
      })

      it('should revert if there is no proposal', async () => {
        await expectRevert(
          validatorProxy.upgradeAggregator({
            from: owner,
          }),
          'No proposal',
        )
      })
    })

    describe('success', () => {
      let newAggregator
      let receipt
      beforeEach(async () => {
        newAggregator = accounts[3]
        await validatorProxy.proposeNewAggregator(newAggregator, {
          from: owner,
        })
        receipt = await validatorProxy.upgradeAggregator({
          from: owner,
        })
      })

      it('should emit an event', async () => {
        await expectEvent(receipt, 'AggregatorUpgraded', {
          previous: aggregator,
          current: newAggregator,
        })
      })

      it('should upgrade the addresses', async () => {
        const response = await validatorProxy.getAggregators()
        assert.equal(response.current, newAggregator)
        assert.equal(response.hasProposal, false)
        assert.equal(response.proposed, constants.ZERO_ADDRESS)
      })
    })
  })

  describe('#proposeNewValidator', () => {
    let newValidator

    before(() => {
      newValidator = accounts[3]
    })

    it('should only be called by the owner', async () => {
      const stranger = accounts[4]
      await expectRevert(
        validatorProxy.proposeNewValidator(newValidator, {
          from: stranger,
        }),
        'Only callable by owner',
      )
    })

    describe('success', () => {
      let receipt

      beforeEach(async () => {
        receipt = await validatorProxy.proposeNewValidator(newValidator, {
          from: owner,
        })
      })

      it('should emit an event', async () => {
        await expectEvent(receipt, 'ValidatorProposed', {
          validator: newValidator,
        })
      })

      it('should set the correct address and hasProposal is true', async () => {
        const response = await validatorProxy.getValidators()
        assert.equal(response.current, validator)
        assert.equal(response.hasProposal, true)
        assert.equal(response.proposed, newValidator)
      })

      it('should set a zero address and hasProposal is false', async () => {
        receipt = await validatorProxy.proposeNewValidator(
          constants.ZERO_ADDRESS,
          {
            from: owner,
          },
        )
        const response = await validatorProxy.getValidators()
        assert.equal(response.current, validator)
        assert.equal(response.hasProposal, false)
        assert.equal(response.proposed, constants.ZERO_ADDRESS)
      })
    })
  })

  describe('#upgradeValidator', () => {
    describe('failure', () => {
      it('should only be called by the owner', async () => {
        const stranger = accounts[4]
        await expectRevert(
          validatorProxy.upgradeValidator({
            from: stranger,
          }),
          'Only callable by owner',
        )
      })

      it('should revert if there is no proposal', async () => {
        await expectRevert(
          validatorProxy.upgradeValidator({
            from: owner,
          }),
          'No proposal',
        )
      })
    })

    describe('success', () => {
      let newValidator
      let receipt
      beforeEach(async () => {
        newValidator = accounts[3]
        await validatorProxy.proposeNewValidator(newValidator, {
          from: owner,
        })
        receipt = await validatorProxy.upgradeValidator({
          from: owner,
        })
      })

      it('should emit an event', async () => {
        await expectEvent(receipt, 'ValidatorUpgraded', {
          previous: validator,
          current: newValidator,
        })
      })

      it('should upgrade the addresses', async () => {
        const response = await validatorProxy.getValidators()
        assert.equal(response.current, newValidator)
        assert.equal(response.hasProposal, false)
        assert.equal(response.proposed, constants.ZERO_ADDRESS)
      })
    })
  })

  describe('#validate', () => {
    // TODO
  })
})
