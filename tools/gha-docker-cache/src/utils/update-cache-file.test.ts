import { includes } from './update-cache-file'

describe('update-cache-file tests', () => {
  describe('includes', () => {
    // we use the multiline flag here
    // to simulate 'sed' which will use this regex per line
    const regex = includes('matchme', 'm')
    const cases = [
      [`asdfmatchmeasd`, true],
      [
        `
        heyooo
        matchme asdfsdf
        `,
        true,
      ],
      ['foobar', false],
      [`dontmatchm`, false],
      [
        `
        yomatchme
        andmatchmetoo
        asdfasd`,
        true,
      ],
    ] as const

    it.each(cases)('%s should be matched? %s', (s, expected) => {
      const actual = s.match(regex)
      console.log(actual)
      expect(!!actual).toEqual(expected)
    })
  })
})
