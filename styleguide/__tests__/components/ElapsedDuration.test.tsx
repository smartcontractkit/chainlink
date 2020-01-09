import React from 'react'
import '@testing-library/jest-dom/extend-expect'
import { render } from '@testing-library/react'
import { ElapsedDuration } from '../../src/components/ElapsedDuration'

const START = '2020-01-03T23:45:28.613635Z'
const END = '2020-01-03T23:45:30.166261Z'

describe('ElapsedDuration', () => {
  test('displays the duration between 2 dates in a human readable foramt', () => {
    const { container } = render(<ElapsedDuration start={START} end={END} />)

    expect(container).toHaveTextContent('2s')
  })
})
