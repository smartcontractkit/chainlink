import '@testing-library/jest-dom/extend-expect'
import { uIntFrom } from 'utils/uIntFrom'

describe('utils/uIntFrom/', () => {
  it('should reject fractions ', () => {
    expect(uIntFrom(0.123)).toBeNaN()
    expect(uIntFrom(0.0)).toEqual(0)
  })

  it('should reject negative numbers', () => {
    expect(uIntFrom(-1)).toBeNaN()
    expect(uIntFrom(-0.1)).toBeNaN()
  })

  it('should reject numbers with characters', () => {
    expect(uIntFrom('abc')).toBeNaN()
    expect(uIntFrom('123abc')).toBeNaN()
  })

  it('should not reject numbers', () => {
    expect(uIntFrom('123')).toEqual(123)
    expect(uIntFrom(123)).toEqual(123)
  })
})
