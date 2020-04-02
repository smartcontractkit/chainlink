import { splitOnColon } from './update-docker-files'

describe('update-docker-files tests', () => {
  describe('splitOnColon', () => {
    const cases = [
      ['split:on', ['split', 'on']],
      ['split:on:me', ['split', 'on:me']],
      ['splitonme', ['splitonme']],
    ] as const

    it.each(cases)(
      'it should parse the string "%s" into [%s]',
      (s, expected) => {
        expect(splitOnColon(s)).toEqual(expected)
      },
    )
  })
})
