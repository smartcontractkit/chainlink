import React from 'react'
import { MemoryRouter } from 'react-router-dom'
import { render, screen } from '@testing-library/react'
import RecentlyCreated from 'components/Jobs/RecentlyCreated'
import isoDate, { MINUTE_MS, TWO_MINUTES_MS } from 'test-helpers/isoDate'

const { queryByText } = screen

const renderComponent = (jobs) =>
  render(
    <MemoryRouter>
      <RecentlyCreated jobs={jobs} />
    </MemoryRouter>,
  )

describe('components/Jobs/RecentlyCreated', () => {
  it('shows the id and creation date', () => {
    const minuteAgo = isoDate(Date.now() - MINUTE_MS)
    const jobB = {
      id: 'job_b',
      createdAt: minuteAgo,
    }
    const twoMinutesAgo = isoDate(Date.now() - TWO_MINUTES_MS)
    const jobA = {
      id: 'job_a',
      createdAt: twoMinutesAgo,
    }

    renderComponent([jobB, jobA])

    expect(queryByText('job_a')).toBeInTheDocument()
    expect(queryByText('2 minutes ago')).toBeInTheDocument()
    expect(queryByText('job_b')).toBeInTheDocument()
    expect(queryByText('1 minute ago')).toBeInTheDocument()
  })

  it('shows a loading indicator', () => {
    renderComponent(null)
    expect(queryByText('...')).toBeInTheDocument()
  })

  it('shows a message for no jobs', () => {
    renderComponent([])
    expect(queryByText('No recently created jobs')).toBeInTheDocument()
  })
})
