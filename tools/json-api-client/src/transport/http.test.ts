import { createUrl } from './http'

describe('http tests', () => {
  describe('createUrl', () => {
    // each element is in the format of
    // [expected, base, path, query?]
    const cases = [
      ['http://explorer:3001/foo', 'http://explorer:3001', 'foo', undefined],
      ['http://explorer:3001/foo', 'http://explorer:3001', '/foo', undefined],
      [
        'http://explorer:3001/foo',
        'http://explorer:3001/ignore/this/path',
        '/foo',
        undefined,
      ],
      [
        'http://explorer:3001/foo?bar=baz&boing=boing',
        'http://explorer:3001',
        'foo',
        { bar: 'baz', boing: 'boing' },
      ],
      [
        'http://explorer:3001/foo?stinky=false',
        'http://explorer:3001',
        'foo',
        { stinky: false, shouldNotExist: undefined, shouldNotExist2: null },
      ],
      [
        'http://explorer:3001/jobs/170?page=1&size=10',
        'http://explorer:3001',
        'jobs/170',
        { page: 1, size: 10 },
      ],
    ]

    it.each(cases as any[])(
      '%s\nbase=%s\npath=%s\nquery=%o\n',
      (expected, b, p, q) => {
        expect(createUrl(b, p, q).toString()).toEqual(expected)
      },
    )
  })
})
