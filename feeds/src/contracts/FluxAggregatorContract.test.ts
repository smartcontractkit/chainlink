import { ethers } from 'ethers'
import { mocked } from 'ts-jest/utils'
import { partialAsFull } from '@chainlink/ts-helpers'
import FluxAggregatorAbi from '../contracts/FluxAggregatorAbi.json'
import FluxAggregatorContract from '../contracts/FluxAggregatorContract'
import { FeedConfig } from '../config'
import { createContract } from '../contracts/utils/createContract'

jest.mock('../contracts/utils/createContract', () => ({
  createContract: jest.fn(),
}))
const mockedCreateContract = mocked(createContract, true)

describe('FluxAggregatorContract', () => {
  const latestSubmission = ethers.utils.bigNumberify(5)
  const config = partialAsFull<FeedConfig>({
    contractAddress: '0x1000000000000000000000000000000000000000',
  })

  describe('#reportingRound', () => {
    it('returns the round id from the oracle round state when started at > 0', async () => {
      const startedAt = ethers.utils.bigNumberify(1)
      const ethersContract = partialAsFull<ethers.Contract>({
        // Our test only needs the first 4 items in the tuple. This test ignores the trailing params
        oracleRoundState: () =>
          Promise.resolve([false, 10, latestSubmission, startedAt]),
      })
      mockedCreateContract.mockReturnValue(ethersContract)
      const contract = new FluxAggregatorContract(config, FluxAggregatorAbi)

      const round = await contract.reportingRound()
      expect(round).toEqual(10)
    })

    it('subtracts 1 from the round id when started at = 0', async () => {
      const startedAt = ethers.utils.bigNumberify(0)
      const ethersContract = partialAsFull<ethers.Contract>({
        // Our test only needs the first 4 items in the tuple. This test ignores the trailing params
        oracleRoundState: () =>
          Promise.resolve([false, 100, latestSubmission, startedAt]),
      })
      mockedCreateContract.mockReturnValue(ethersContract)
      const contract = new FluxAggregatorContract(config, FluxAggregatorAbi)

      const round = await contract.reportingRound()
      expect(round).toEqual(99)
    })
  })
})
