import * as React from 'react'

import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'
import userEvent from '@testing-library/user-event'

import { buildRuns } from 'support/factories/gql/fetchJobRuns'
import { JobRunsView, Props as JobRunsViewProps } from './JobRunsView'

const { getAllByRole, findByText, queryByText, getByRole, queryByRole } = screen

function renderComponent(viewProps: JobRunsViewProps) {
  renderWithRouter(
    <>
      <Route exact path="/runs">
        <JobRunsView {...viewProps} />
      </Route>
      <Route exact path="/runs/:runID">
        Run Page
      </Route>
    </>,
    { initialEntries: ['/runs'] },
  )
}

describe('JobRunsView', () => {
  it('renders the job runs table', () => {
    const runs = buildRuns()

    renderComponent({
      loading: false,
      data: {
        jobRuns: {
          results: runs,
          metadata: { total: runs.length },
        },
      },
      page: 1,
      pageSize: 10,
    })

    expect(queryByText('Job Runs')).toBeInTheDocument()

    expect(getAllByRole('row')).toHaveLength(2)

    expect(queryByText(runs[0].id)).toBeInTheDocument()
    expect(queryByText('1 minute ago')).toBeInTheDocument()
    expect(queryByText(/completed/i)).toBeInTheDocument()

    expect(queryByText(runs[1].id)).toBeInTheDocument()
    expect(queryByText('2 minutes ago')).toBeInTheDocument()
    expect(queryByText(/errored/i)).toBeInTheDocument()

    expect(queryByText('1-2 of 2'))
  })

  it('renders the loading spinner', async () => {
    renderComponent({
      loading: true,
      page: 1,
      pageSize: 10,
    })

    expect(queryByRole('progressbar')).toBeInTheDocument()
  })

  it('navigates to the job run details page', async () => {
    const runs = buildRuns()

    renderComponent({
      loading: false,
      data: {
        jobRuns: {
          results: runs,
          metadata: { total: runs.length },
        },
      },
      page: 1,
      pageSize: 10,
    })

    userEvent.click(getByRole('cell', { name: runs[0].id }))

    expect(await findByText('Run Page')).toBeInTheDocument()
  })
})
