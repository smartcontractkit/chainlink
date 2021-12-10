import * as React from 'react'

import { Route, Switch } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'
import userEvent from '@testing-library/user-event'

import { JobCard } from './JobCard'
import { buildJob } from 'support/factories/gql/fetchJob'

const { findByText, getByRole, queryByText } = screen

describe('JobCard', () => {
  let handleDelete: jest.Mock

  function renderComponent(job: JobPayload_Fields) {
    renderWithRouter(
      <>
        <Switch>
          <Route exact path="/jobs/new">
            Link success
          </Route>

          <Route path="/jobs/:id">
            <JobCard job={job} onDelete={handleDelete} />
          </Route>
        </Switch>
      </>,
      { initialEntries: ['/jobs/1'] },
    )
  }

  beforeEach(() => {
    handleDelete = jest.fn()
  })

  it('renders a job', () => {
    const job = buildJob()

    renderComponent(job)

    expect(queryByText(job.id)).toBeInTheDocument()
    expect(queryByText('Direct Request')).toBeInTheDocument()
    expect(queryByText(job.externalJobID)).toBeInTheDocument()
    expect(queryByText('1 minute ago')).toBeInTheDocument()
  })

  it('calls delete', () => {
    const job = buildJob()

    renderComponent(job)

    userEvent.click(getByRole('button', { name: /open-menu/i }))
    userEvent.click(getByRole('menuitem', { name: /delete/i }))
    userEvent.click(getByRole('button', { name: /confirm/i }))

    expect(handleDelete).toHaveBeenCalled()
  })

  it('navigates to the new job page on duplicate click', async () => {
    const job = buildJob()

    renderComponent(job)

    userEvent.click(getByRole('button', { name: /open-menu/i }))
    userEvent.click(getByRole('menuitem', { name: /duplicate/i }))

    expect(await findByText('Link success')).toBeInTheDocument()
  })
})
