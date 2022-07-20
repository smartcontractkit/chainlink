import * as React from 'react'

import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'

import { buildError, buildJob } from 'support/factories/gql/fetchJob'
import { TabErrors } from './TabErrors'

const { getAllByRole, queryByText } = screen

describe('TabErrors', () => {
  function renderComponent(job: JobPayload_Fields) {
    renderWithRouter(
      <>
        <Route exact path="/jobs/:id/errors">
          <TabErrors job={job} />
        </Route>
      </>,
      { initialEntries: ['/jobs/1/errors'] },
    )
  }

  it('renders the errors', () => {
    const error = buildError()
    const job = buildJob({
      errors: [error],
    })

    renderComponent(job)

    expect(getAllByRole('row')).toHaveLength(2)

    expect(queryByText('Occurrences')).toBeInTheDocument()
    expect(queryByText('Created')).toBeInTheDocument()
    expect(queryByText('Last Seen')).toBeInTheDocument()
    expect(queryByText('Message')).toBeInTheDocument()

    expect(queryByText(error.occurrences)).toBeInTheDocument()
    expect(queryByText('1 minute ago')).toBeInTheDocument()
    expect(queryByText('2 minutes ago')).toBeInTheDocument()
    expect(queryByText(error.description)).toBeInTheDocument()
  })

  it('renders no errors', () => {
    const job = buildJob()

    renderComponent(job)

    expect(queryByText('No errors')).toBeInTheDocument()
  })
})
