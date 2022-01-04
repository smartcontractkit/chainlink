import * as React from 'react'

import { renderWithRouter, screen } from 'support/test-utils'

import { ActivityCard, Props as ActivityCardProps } from './ActivityCard'
import { buildRuns } from 'support/factories/gql/fetchRecentJobRuns'

const { queryByRole, queryByText } = screen

function renderComponent(cardProps: Omit<ActivityCardProps, 'classes'>) {
  renderWithRouter(<ActivityCard {...cardProps} />)
}

describe('ActivityCard', () => {
  it('renders the jobs', () => {
    const runs = buildRuns()

    renderComponent({
      loading: false,
      data: {
        jobRuns: {
          results: runs,
          metadata: {
            total: runs.length,
          },
        },
      },
      maxRunsSize: 2,
    })

    expect(queryByText(`Job: ${runs[0].job.id}`)).toBeInTheDocument()
    expect(queryByText(`Job: ${runs[1].job.id}`)).toBeInTheDocument()
    expect(queryByRole('link', { name: /view more/i })).toBeNull()
  })

  it('shows the view more button', () => {
    const runs = buildRuns()

    renderComponent({
      loading: false,
      data: {
        jobRuns: {
          results: runs,
          metadata: {
            total: runs.length,
          },
        },
      },
      maxRunsSize: 1,
    })

    expect(queryByRole('link', { name: /view more/i })).toBeInTheDocument()
  })

  it('renders no content', () => {
    renderComponent({
      loading: false,
      data: {
        jobRuns: {
          results: [],
          metadata: {
            total: 0,
          },
        },
      },
      maxRunsSize: 1,
    })

    expect(queryByText('No recent activity')).toBeInTheDocument()
  })

  it('renders a loading spinner', () => {
    renderComponent({
      loading: true,
      maxRunsSize: 1,
    })

    expect(queryByRole('progressbar')).toBeInTheDocument()
  })

  it('renders an error message', () => {
    renderComponent({
      loading: false,
      errorMsg: 'error message',
      maxRunsSize: 1,
    })

    expect(queryByText('error message')).toBeInTheDocument()
  })
})
