import * as React from 'react'

import { GraphQLError } from 'graphql'
import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'
import { MockedProvider, MockedResponse } from '@apollo/client/testing'

import { buildRun } from 'support/factories/gql/fetchJobRun'
import { JobRunScreen, JOB_RUN_QUERY } from './JobRunScreen'
import { waitForLoading } from 'support/test-helpers/wait'

const { findByTestId, findByText } = screen

function fetchJobRunQuery(run: JobRunPayload_Fields) {
  return {
    request: {
      query: JOB_RUN_QUERY,
      variables: { id: '1' },
    },
    result: {
      data: {
        jobRun: run,
      },
    },
  }
}

describe('JobScreen', () => {
  function renderComponent(mocks: MockedResponse[]) {
    renderWithRouter(
      <>
        <Route path="/runs/:id">
          <MockedProvider mocks={mocks} addTypename={false}>
            <JobRunScreen />
          </MockedProvider>
        </Route>
      </>,
      { initialEntries: [`/runs/1`] },
    )
  }

  it('renders the page', async () => {
    const payload = buildRun()
    const mocks: MockedResponse[] = [fetchJobRunQuery(payload)]

    renderComponent(mocks)

    await waitForLoading()

    expect(await findByText(`Job Run #${payload.id}`)).toBeInTheDocument()
  })

  it('renders the not found page', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: JOB_RUN_QUERY,
          variables: { id: '1' },
        },
        result: {
          data: {
            jobRun: {
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
          query: JOB_RUN_QUERY,
          variables: { id: '1' },
        },
        result: {
          errors: [new GraphQLError('Error!')],
        },
      },
    ]

    renderComponent(mocks)

    expect(await findByText('Error: Error!')).toBeInTheDocument()
  })
})
