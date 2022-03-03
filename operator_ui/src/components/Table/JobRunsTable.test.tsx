import * as React from 'react'

import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'
import userEvent from '@testing-library/user-event'

import { JobRunsTable, Props as JobRunsTableProps } from './JobRunsTable'
import isoDate, { MINUTE_MS } from 'test-helpers/isoDate'

const { getByRole, queryByText } = screen

describe('JobRunsTable', () => {
  function renderComponent(props: Omit<JobRunsTableProps, 'classes'>) {
    renderWithRouter(
      <>
        <Route exact path="/">
          <JobRunsTable {...props} />
        </Route>
        <Route path="/runs/:id">Run Page</Route>
      </>,
    )
  }

  function buildRun(
    overrides?: Partial<JobRunsTableProps['runs'][0]>,
  ): JobRunsTableProps['runs'][0] {
    const minuteAgo = isoDate(Date.now() - MINUTE_MS)
    const twoMinutesAgo = isoDate(Date.now() - 2 * MINUTE_MS)

    return {
      id: '1',
      createdAt: twoMinutesAgo,
      errors: [],
      status: 'COMPLETED',
      finishedAt: minuteAgo,
      ...overrides,
    }
  }

  it('renders a completed run', () => {
    const run = buildRun()

    renderComponent({ runs: [run] })

    expect(queryByText(run.id)).toBeInTheDocument()
    expect(queryByText('2 minutes ago')).toBeInTheDocument()
    expect(queryByText(/completed/i)).toBeInTheDocument()
  })

  it('renders an errored run', () => {
    const run = buildRun({
      errors: ['some error'],
      status: 'ERRORED',
    })

    renderComponent({ runs: [run] })

    expect(queryByText(run.id)).toBeInTheDocument()
    expect(queryByText('2 minutes ago')).toBeInTheDocument()
    expect(queryByText(/errored/i)).toBeInTheDocument()
  })

  it('renders an suspended run', () => {
    const run = buildRun({
      status: 'SUSPENDED',
    })

    renderComponent({ runs: [run] })

    expect(queryByText(run.id)).toBeInTheDocument()
    expect(queryByText('2 minutes ago')).toBeInTheDocument()
    expect(queryByText(/suspended/i)).toBeInTheDocument()
  })

  it('renders an running run', () => {
    const run = buildRun({
      finishedAt: null,
      status: 'RUNNING',
    })

    renderComponent({ runs: [run] })

    expect(queryByText(run.id)).toBeInTheDocument()
    expect(queryByText('2 minutes ago')).toBeInTheDocument()
    expect(queryByText(/running/i)).toBeInTheDocument()
  })

  it('navigates to the run', () => {
    const run = buildRun()

    renderComponent({ runs: [run] })

    userEvent.click(getByRole('cell', { name: run.id }))

    expect(queryByText('Run Page')).toBeInTheDocument()
  })
})
