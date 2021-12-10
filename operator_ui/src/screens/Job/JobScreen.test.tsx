import * as React from 'react'

import { GraphQLError } from 'graphql'
import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'
import { MockedProvider, MockedResponse } from '@apollo/client/testing'
import userEvent from '@testing-library/user-event'

import { buildJob } from 'support/factories/gql/fetchJob'
import { JobScreen, JOB_QUERY, DELETE_JOB_MUTATION } from './JobScreen'
import Notifications from 'pages/Notifications'
import { waitForLoading } from 'support/test-helpers/wait'

const { findByTestId, findByText, getByRole } = screen

function renderComponent(mocks: MockedResponse[], initialEntry?: string) {
  const initialEntries = [`/jobs/1`]
  if (initialEntry) {
    initialEntries[0] = initialEntry
  }

  renderWithRouter(
    <>
      <Notifications />
      <Route path="/jobs/:id">
        <MockedProvider mocks={mocks} addTypename={false}>
          <JobScreen />
        </MockedProvider>
      </Route>

      <Route exact path="/jobs">
        Jobs Page
      </Route>
    </>,
    { initialEntries },
  )
}

function fetchJobQuery(job: JobPayload_Fields) {
  return {
    request: {
      query: JOB_QUERY,
      variables: { id: '1', offset: 0, limit: 5 },
    },
    result: {
      data: {
        job,
      },
    },
  }
}

describe('JobScreen', () => {
  it('renders the page', async () => {
    const payload = buildJob()
    const mocks: MockedResponse[] = [fetchJobQuery(payload)]

    renderComponent(mocks)

    await waitForLoading()

    expect(await findByText(payload.name)).toBeInTheDocument()
  })

  it('renders the not found page', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: JOB_QUERY,
          variables: { id: '1', offset: 0, limit: 5 },
        },
        result: {
          data: {
            job: {
              __typename: 'NotFoundError',
              message: 'job not found',
            },
          },
        },
      },
    ]

    renderComponent(mocks)

    await waitForLoading()

    expect(await findByTestId('not-found-page')).toBeInTheDocument()
  })

  it('renders GQL query errors', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: JOB_QUERY,
          variables: { id: '1', offset: 0, limit: 5 },
        },
        result: {
          errors: [new GraphQLError('Error!')],
        },
      },
    ]

    renderComponent(mocks)

    expect(await findByText('Error: Error!')).toBeInTheDocument()
  })

  it('deletes the job', async () => {
    const payload = buildJob()
    const mocks: MockedResponse[] = [
      fetchJobQuery(payload),
      {
        request: {
          query: DELETE_JOB_MUTATION,
          variables: { id: payload.id },
        },
        result: {
          data: {
            deleteJob: {
              __typename: 'DeleteJobSuccess',
              job: {
                id: payload.id,
              },
            },
          },
        },
      },
    ]

    renderComponent(mocks)

    await waitForLoading()

    userEvent.click(getByRole('button', { name: /open-menu/i }))
    userEvent.click(getByRole('menuitem', { name: /delete/i }))
    userEvent.click(getByRole('button', { name: /confirm/i }))

    expect(await findByText('Jobs Page')).toBeInTheDocument()
  })

  it('handles not found on job delete', async () => {
    const payload = buildJob()
    const mocks: MockedResponse[] = [
      fetchJobQuery(payload),
      {
        request: {
          query: DELETE_JOB_MUTATION,
          variables: { id: payload.id },
        },
        result: {
          data: {
            deleteJob: {
              __typename: 'NotFoundError',
              message: 'job not found',
            },
          },
        },
      },
    ]

    renderComponent(mocks)

    await waitForLoading()

    userEvent.click(getByRole('button', { name: /open-menu/i }))
    userEvent.click(getByRole('menuitem', { name: /delete/i }))
    userEvent.click(getByRole('button', { name: /confirm/i }))

    expect(await findByText('job not found')).toBeInTheDocument()
  })

  it('runs a webhook job', async () => {
    // TODO - Add a Run test once we switch over to GQL
  })

  it('calls a refetch', async () => {
    const payload = buildJob()

    const mocks: MockedResponse[] = [
      fetchJobQuery(payload),
      {
        request: {
          query: JOB_QUERY,
          variables: { id: '1', offset: 0, limit: 1 },
        },
        result: {
          data: {
            job: payload,
          },
        },
      },
      {
        request: {
          query: JOB_QUERY,
          variables: { id: '1', offset: 1, limit: 25 },
        },
        result: {
          data: {
            job: payload,
          },
        },
      },
    ]

    renderComponent(mocks)

    await waitForLoading()

    userEvent.click(getByRole('tab', { name: /^runs/i }))

    // Default value of the rows per page select
    userEvent.click(getByRole('button', { name: /10/i }))
    userEvent.click(getByRole('option', { name: /25/i }))
  })

  it('calls a refetches recent runs', async () => {
    const payload = buildJob()

    const mocks: MockedResponse[] = [
      fetchJobQuery(payload),
      fetchJobQuery(payload),
    ]

    renderComponent(mocks, '/jobs/1/definition')

    await waitForLoading()

    userEvent.click(getByRole('tab', { name: /overview/i }))
  })
})
