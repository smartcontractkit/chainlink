import formatRequestURI from '../src/formatRequestURI'

describe('formatRequestURI', () => {
  describe('port specified', () => {
    it('returns host and port in URI', () => {
      const uri = formatRequestURI(
        '/api',
        {},
        { hostname: 'localhost', port: '6689' },
      )

      expect(uri).toEqual('localhost:6689/api')
    })
  })

  describe('no port specified', () => {
    it('returns just the path in the URI', () => {
      const uri = formatRequestURI('/api')

      expect(uri).toEqual('/api')
    })
  })
})
