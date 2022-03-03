import * as React from 'react'

import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'
import userEvent from '@testing-library/user-event'

import { ActivityRow } from './ActivityRow'
import { buildRun } from 'support/factories/gql/fetchRecentJobRuns'

const { getByRole, queryByRole, queryByTestId, queryByText } = screen

describe('ActivityRow', () => {
  function renderComponent(run: RecentJobRunsPayload_ResultsFields) {
    renderWithRouter(
      <>
        <Route exact path="/">
          <table>
            <tbody>
              <ActivityRow run={run} />
            </tbody>
          </table>
        </Route>

        <Route exact path="/jobs/:id">
          Job Page
        </Route>
        <Route exact path="/runs/:runID">
          Run Page
        </Route>
      </>,
    )
  }

  it('renders a row', () => {
    const run = buildRun()

    renderComponent(run)

    expect(
      queryByRole('link', { name: `Job: ${run.job.id}` }),
    ).toBeInTheDocument()
    expect(queryByRole('link', { name: `Run: ${run.id}` })).toBeInTheDocument()
    expect(queryByText('1 minute ago')).toBeInTheDocument()
    expect(queryByTestId(/completed/i)).toBeInTheDocument()
  })

  it('navigates to the job page', () => {
    const run = buildRun()

    renderComponent(run)

    userEvent.click(getByRole('link', { name: `Job: ${run.job.id}` }))

    expect(queryByText('Job Page')).toBeInTheDocument()
  })

  it('navigates to the run page', () => {
    const run = buildRun()

    renderComponent(run)

    userEvent.click(getByRole('link', { name: `Run: ${run.id}` }))

    expect(queryByText('Run Page')).toBeInTheDocument()
  })
})
