import * as React from 'react'

import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'
import userEvent from '@testing-library/user-event'

import { JobRow } from './JobRow'
import { buildJob } from 'support/factories/gql/fetchJobs'

const { findByText, getByRole, queryByText } = screen

function renderComponent(job: JobsPayload_ResultsFields) {
  renderWithRouter(
    <>
      <Route exact path="/">
        <table>
          <tbody>
            <JobRow job={job} />
          </tbody>
        </table>
      </Route>
      <Route exact path="/jobs/:id">
        Link Success
      </Route>
    </>,
  )
}

describe('JobRow', () => {
  it('renders the row', () => {
    const job = buildJob()

    renderComponent(job)

    expect(queryByText('1')).toBeInTheDocument()
    expect(queryByText('job 1')).toBeInTheDocument()
    expect(
      queryByText('00000000-0000-0000-0000-000000000001'),
    ).toBeInTheDocument()
    expect(queryByText('Flux Monitor')).toBeInTheDocument()
    expect(queryByText('1 minute ago')).toBeInTheDocument()
  })

  it('links to the job details', async () => {
    const job = buildJob()

    renderComponent(job)

    const link = getByRole('link', { name: /1/i })
    expect(link).toHaveAttribute('href', '/jobs/1')

    userEvent.click(link)

    expect(await findByText('Link Success')).toBeInTheDocument()
  })
})
