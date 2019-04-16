import formatRequestURI from 'utils/formatRequestURI'

describe('formatRequestURI', () => {
  describe('port specified', () => {
    it('returns host and port in URI', () => {
      expect(
        formatRequestURI('/api', {}, { hostname: 'localhost', port: 6689 })
      ).toEqual('localhost:6689/api')
    })
  })

  describe('no port specified', () => {
    it('returns just the path in the URI', () => {
      expect(formatRequestURI('/api')).toEqual('/api')
    })
  })
})
