import { assertBigNum } from './support/matchers'

contract('Median', accounts => {
  const Median = artifacts.require('ConcreteMedian.sol')
  let median

  beforeEach(async () => {
    median = await Median.new()
  })

  describe('testing various lists', () => {
    const tests = [
      {
        name: 'ordered ascending',
        responses: [0, 1, 2, 3, 4, 5, 6, 7],
        want: 3,
      },
      {
        name: 'ordered descending',
        responses: [7, 6, 5, 4, 3, 2, 1, 0],
        want: 3,
      },
      {
        name: 'unordered 1',
        responses: [1001, 1, 101, 10, 11, 0, 111],
        want: 11,
      },
      {
        name: 'unordered 2',
        responses: [8, 8, 4, 5, 5, 7, 9, 5, 9],
        want: 7,
      },
      {
        name: 'unordered 3',
        responses: [33, 44, 89, 101, 67, 7, 23, 55, 88, 324, 0, 88],
        want: 61, // 67 + 55 / 2
      },
      {
        name: 'long unordered',
        responses: [
          333121,
          323453,
          337654,
          345363,
          345363,
          333456,
          335477,
          333323,
          332352,
          354648,
          983260,
          333856,
          335468,
          376987,
          333253,
          388867,
          337879,
          333324,
          338678,
        ],
        want: 335477,
      },
    ]

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    for (const test of tests) {
      it(test.name, async () => {
        assertBigNum(test.want, await median.publicGet.call(test.responses))
      })
    }
  })
})
