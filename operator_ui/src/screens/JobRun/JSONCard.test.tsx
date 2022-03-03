import * as React from 'react'
import { buildRun } from 'support/factories/gql/fetchJobRun'

import { render, screen } from 'support/test-utils'

import { JSONCard } from './JSONCard'

const { getByTestId } = screen

describe('JSONCard', () => {
  it('renders the run as json', () => {
    const start = '2020-01-03T22:45:00.166261Z'
    const end1m = '2020-01-03T22:46:00.166261Z'

    const run = buildRun({
      createdAt: start,
      finishedAt: end1m,
    })

    render(<JSONCard run={run} />)

    expect(getByTestId('pretty-json').textContent).toMatchSnapshot()
  })
})
