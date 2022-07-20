import * as React from 'react'

import { Route } from 'react-router'
import { renderWithRouter, screen } from 'support/test-utils'
import userEvent from '@testing-library/user-event'

import { buildRun } from 'support/factories/gql/fetchJobRun'
import { JobRunCard } from './JobRunCard'

const { findByText, getByRole, queryByText } = screen

describe('JobRunCard', () => {
  function renderComponent(run: JobRunPayload_Fields) {
    renderWithRouter(
      <>
        <Route exact path="/">
          <JobRunCard run={run} />
        </Route>

        <Route path="/jobs/:id">Job Page</Route>
      </>,
    )
  }

  it('renders the card', async () => {
    const run = buildRun()

    renderComponent(run)

    expect(queryByText('ID')).toBeInTheDocument()
    expect(queryByText('Job')).toBeInTheDocument()
    expect(queryByText('Started')).toBeInTheDocument()
    expect(queryByText('Finished')).toBeInTheDocument()

    expect(queryByText(run.id)).toBeInTheDocument()
    expect(queryByText(run.job.name)).toBeInTheDocument()
    expect(queryByText('2 minutes ago')).toBeInTheDocument()
    expect(queryByText('1 minute ago')).toBeInTheDocument()
  })

  it('navigates to the job page', async () => {
    const run = buildRun()

    renderComponent(run)

    userEvent.click(getByRole('link', { name: run.job.name }))

    expect(await findByText('Job Page')).toBeInTheDocument()
  })
})
