import { groups } from 'state/ducks/listing/selectors'
import feeds from 'feeds.json'

const answersMock = [
  {
    answer: 'answer',
    config: feeds[0],
  },
  {
    answer: 'answer2',
    config: feeds[1],
  },
]

describe('state/ducks/listing/selectors', () => {
  it('should build two groups', () => {
    const selected = groups.resultFunc([])
    expect(selected).toHaveLength(2)
    expect(selected[0].name).toMatch('USD')
    expect(selected[1].name).toMatch('ETH')
  })

  it('should return contract configs and answers', () => {
    const selected = groups.resultFunc(answersMock)
    expect(selected[0].list[0].config).toEqual(feeds[0])
    expect(selected[0].list[1].config).toEqual(feeds[1])
    expect(selected[0].list[0].answer).toMatch(answersMock[0].answer)
    expect(selected[0].list[1].answer).toMatch(answersMock[1].answer)
  })

  it('should return contract configs without answers', () => {
    const selected = groups.resultFunc([
      {
        config: feeds[0],
      },
    ])
    expect(selected[0].list[0].config).toEqual(feeds[0])
    expect(selected[0].list[0]).not.toHaveProperty('answer')
  })
})
