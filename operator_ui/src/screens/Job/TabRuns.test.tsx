import * as React from 'react'

import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'
import userEvent from '@testing-library/user-event'

import { buildJob, buildRun } from 'support/factories/gql/fetchJob'
import { TabRuns } from './TabRuns'

const { getAllByRole, getByRole, queryByText } = screen

describe('TabRuns', () => {
  let handleFetchMore: jest.Mock

  beforeEach(() => {
    handleFetchMore = jest.fn()
  })

  function renderComponent(job: JobPayload_Fields, page = 1, per = 10) {
    renderWithRouter(
      <>
        <Route exact path="/jobs/:id/runs">
          <TabRuns job={job} fetchMore={handleFetchMore} />
        </Route>

        <Route exact path="/runs/:runId">
          Run Page
        </Route>
      </>,
      { initialEntries: [`/jobs/1/runs?page=${page}&per=${per}`] },
    )
  }

  it('renders a page of job runs', () => {
    const runs = [
      buildRun({ id: '1' }),
      buildRun({ id: '2' }),
      buildRun({ id: '3' }),
      buildRun({ id: '4' }),
      buildRun({ id: '5' }),
    ]

    const job = buildJob({
      runs: {
        results: runs,
        metadata: { total: runs.length },
      },
    })

    renderComponent(job)

    expect(getAllByRole('row')).toHaveLength(5)

    for (const run of runs) {
      expect(queryByText(run.id)).toBeInTheDocument()
    }
  })

  it('navigates to the run page', () => {
    const job = buildJob()

    renderComponent(job)

    expect(getAllByRole('row')).toHaveLength(1)

    userEvent.click(getByRole('cell', { name: job.runs.results[0].id }))

    expect(queryByText('Run Page')).toBeInTheDocument()
  })

  it('pages next', () => {
    const runs = [buildRun({ id: '1' }), buildRun({ id: '2' })]

    const job = buildJob({
      runs: {
        results: runs,
        metadata: { total: 4 },
      },
    })

    renderComponent(job, 1, 2)

    expect(getAllByRole('row')).toHaveLength(2)
    expect(getByRole('button', { name: /prev-page/i })).toBeDisabled()

    userEvent.click(getByRole('button', { name: /next-page/i }))

    expect(handleFetchMore).toHaveBeenCalledWith(2, 2)
  })

  it('pages prev', () => {
    const runs = [buildRun({ id: '1' }), buildRun({ id: '2' })]

    const job = buildJob({
      runs: {
        results: runs,
        metadata: { total: 4 },
      },
    })

    renderComponent(job, 2, 2)

    expect(getAllByRole('row')).toHaveLength(2)
    expect(getByRole('button', { name: /next-page/i })).toBeDisabled()

    userEvent.click(getByRole('button', { name: /prev-page/i }))

    expect(handleFetchMore).toHaveBeenCalledWith(1, 2)
  })

  it('changes per', () => {
    const job = buildJob()

    renderComponent(job)

    // Default value of the rows per page select
    userEvent.click(getByRole('button', { name: /10/i }))

    userEvent.click(getByRole('option', { name: /25/i }))

    expect(handleFetchMore).toHaveBeenCalledWith(1, 25)
  })
})
