const { assert, web3, artifacts } = require('hardhat')
const {
  constants,
  expectRevert,
  expectEvent,
} = require('@openzeppelin/test-helpers')
const { publicAbi } = require('./helpers')

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

  it('has a limited public interface', async () => {
    publicAbi(validatorProxy, [
      // ConfirmedOwner functions
      'acceptOwnership',
      'owner',
      'transferOwnership',
      // ValidatorProxy functions
      'validate',
      'proposeNewAggregator',
      'upgradeAggregator',
      'getAggregators',
      'proposeNewValidator',
      'upgradeValidator',
      'getValidators',
    ])
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
    describe('failure', () => {
      it('reverts when not called by aggregator or proposed aggregator', async () => {
        let stranger = accounts[9]
        await expectRevert(
          validatorProxy.validate(99, 88, 77, 66, { from: stranger }),
          'Not a configured aggregator',
        )
      })

      it('reverts when there is no validator set', async () => {
        validatorProxy = await ValidatorProxyArtifact.new(
          aggregator,
          constants.ZERO_ADDRESS,
          { from: owner },
        )
        await expectRevert(
          validatorProxy.validate(99, 88, 77, 66, { from: aggregator }),
          'No validator set',
        )
      })
    })

    describe('success', () => {
      describe('from the aggregator', () => {
        let MockValidatorArtifact
        let mockValidator1
        let receipt

        beforeEach(async () => {
          MockValidatorArtifact = artifacts.require('MockAggregatorValidator')
          mockValidator1 = await MockValidatorArtifact.new(1)
          validatorProxy = await ValidatorProxyArtifact.new(
            aggregator,
            mockValidator1.address,
            { from: owner },
          )
        })

        describe('for a single validator', () => {
          beforeEach(async () => {
            receipt = await validatorProxy.validate(200, 300, 400, 500, {
              from: aggregator,
            })
          })

          it('calls validate on the validator', async () => {
            await expectEvent.inTransaction(
              receipt.tx,
              MockValidatorArtifact,
              'ValidateCalled',
              {
                id: '1',
                previousRoundId: '200',
                previousAnswer: '300',
                currentRoundId: '400',
                currentAnswer: '500',
              },
            )
          })

          it('uses a specific amount of gas', async () => {
            assert.equal(receipt.receipt.gasUsed, 34256)
          })
        })

        describe('for a validator and a proposed validator', () => {
          let mockValidator2

          beforeEach(async () => {
            mockValidator2 = await MockValidatorArtifact.new(2)
            await validatorProxy.proposeNewValidator(mockValidator2.address, {
              from: owner,
            })
            receipt = await validatorProxy.validate(2000, 3000, 4000, 5000, {
              from: aggregator,
            })
          })

          it('calls validate on the validator', async () => {
            await expectEvent.inTransaction(
              receipt.tx,
              MockValidatorArtifact,
              'ValidateCalled',
              {
                id: '1',
                previousRoundId: '2000',
                previousAnswer: '3000',
                currentRoundId: '4000',
                currentAnswer: '5000',
              },
            )
          })

          it('also calls validate on the proposed validator', async () => {
            await expectEvent.inTransaction(
              receipt.tx,
              MockValidatorArtifact,
              'ValidateCalled',
              {
                id: '2',
                previousRoundId: '2000',
                previousAnswer: '3000',
                currentRoundId: '4000',
                currentAnswer: '5000',
              },
            )
          })

          it('uses a specific amount of gas', async () => {
            assert.equal(receipt.receipt.gasUsed, 41958)
          })
        })
      })

      describe('from the proposed aggregator', () => {
        let newAggregator
        beforeEach(async () => {
          newAggregator = accounts[3]
          await validatorProxy.proposeNewAggregator(newAggregator, {
            from: owner,
          })
        })

        it('emits an event', async () => {
          const receipt = await validatorProxy.validate(555, 666, 777, 888, {
            from: newAggregator,
          })
          await expectEvent(receipt, 'ProposedAggregatorValidateCall', {
            proposed: newAggregator,
            previousRoundId: '555',
            previousAnswer: '666',
            currentRoundId: '777',
            currentAnswer: '888',
          })
        })
      })
    })
  })
})
