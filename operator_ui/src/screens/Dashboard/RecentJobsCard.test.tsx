import * as React from 'react'

import { renderWithRouter, screen } from 'support/test-utils'

import { buildRecentJobs } from 'support/factories/gql/fetchRecentJobs'
import { RecentJobsCard, Props as RecentJobsCardProps } from './RecentJobsCard'

const { queryByRole, queryByText } = screen

function renderComponent(cardProps: Omit<RecentJobsCardProps, 'classes'>) {
  renderWithRouter(<RecentJobsCard {...cardProps} />)
}

describe('RecentJobsCard', () => {
  it('renders the jobs', () => {
    const jobs = buildRecentJobs()

    renderComponent({
      loading: false,
      data: {
        jobs: {
          results: jobs,
        },
      },
    })

    expect(queryByText(jobs[0].name)).toBeInTheDocument()
    expect(queryByText(jobs[1].name)).toBeInTheDocument()
  })

  it('renders no content', () => {
    renderComponent({
      loading: false,
      data: {
        jobs: {
          results: [],
        },
      },
    })

    expect(queryByText('No recently created jobs')).toBeInTheDocument()
  })

  it('renders a loading spinner', () => {
    renderComponent({
      loading: true,
    })

    expect(queryByRole('progressbar')).toBeInTheDocument()
  })

  it('renders an error message', () => {
    renderComponent({
      loading: false,
      errorMsg: 'error message',
    })

    expect(queryByText('error message')).toBeInTheDocument()
  })
})
