import * as React from 'react'

import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'

import { RecentJobRow } from './RecentJobRow'
import { buildRecentJob } from 'support/factories/gql/fetchRecentJobs'
import userEvent from '@testing-library/user-event'

const { getByRole, queryByText } = screen

describe('RecentJobRow', () => {
  function renderComponent(job: RecentJobsPayload_ResultsFields) {
    renderWithRouter(
      <>
        <Route exact path="/">
          <table>
            <tbody>
              <RecentJobRow job={job} />
            </tbody>
          </table>
        </Route>

        <Route path="/jobs/:id">Job Page</Route>
      </>,
    )
  }

  it('renders a row', () => {
    const job = buildRecentJob()

    renderComponent(job)

    expect(queryByText(job.name)).toBeInTheDocument()
    expect(queryByText('1 minute ago')).toBeInTheDocument()
  })

  it('navigates to the job page', () => {
    const job = buildRecentJob()

    renderComponent(job)

    userEvent.click(getByRole('link', { name: job.name }))

    expect(queryByText('Job Page')).toBeInTheDocument()
  })
})
