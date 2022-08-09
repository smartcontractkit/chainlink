import { ethers } from 'hardhat'
import { assert } from 'chai'
import { evmRevert } from '../../test-helpers/matchers'
import { UpkeepTranscoderDev__factory as UpkeepTranscoderDevFactory } from '../../../typechain/factories/UpkeepTranscoderDev__factory'
import { UpkeepTranscoderDev } from '../../../typechain'

let upkeepTranscoderDevFactory: UpkeepTranscoderDevFactory
let transcoder: UpkeepTranscoderDev

before(async () => {
  upkeepTranscoderDevFactory = await ethers.getContractFactory(
    'UpkeepTranscoderDev',
  )
})

describe('UpkeepTranscoderDev', () => {
  const formatV1 = 0

  beforeEach(async () => {
    transcoder = await upkeepTranscoderDevFactory.deploy()
  })

  describe('#typeAndVersion', () => {
    it('uses the correct type and version', async () => {
      const typeAndVersion = await transcoder.typeAndVersion()
      assert.equal(typeAndVersion, 'UpkeepTranscoder 1.3.0')
    })
  })

  describe('#transcodeUpkeeps', () => {
    const encodedData = '0xc0ffee'

    it('reverts if the from type is not V1', async () => {
      await evmRevert(
        transcoder.transcodeUpkeeps(2, 1, encodedData),
        'function was called with incorrect parameters',
      )
    })

    context('when from and to versions equal V1', () => {
      it('returns the data that was passed in', async () => {
        const response = await transcoder.transcodeUpkeeps(
          formatV1,
          formatV1,
          encodedData,
        )
        assert.equal(encodedData, response)
      })
    })
  })
})
