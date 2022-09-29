import React from 'react'

import { MockedProvider, MockedResponse } from '@apollo/client/testing'
import { GraphQLError } from 'graphql'
import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'test-utils'
import userEvent from '@testing-library/user-event'

import { JobRunsScreen, JOB_RUNS_QUERY } from './JobRunsScreen'
import { buildRun, buildRuns } from 'support/factories/gql/fetchJobRuns'
import { waitForLoading } from 'support/test-helpers/wait'

const { findAllByRole, findByText, getByRole, queryByText } = screen

function renderComponent(mocks: MockedResponse[]) {
  renderWithRouter(
    <>
      <Route exact path="/runs">
        <MockedProvider mocks={mocks} addTypename={false}>
          <JobRunsScreen />
        </MockedProvider>
      </Route>
    </>,
    { initialEntries: ['/runs?per=2'] },
  )
}

describe('JobRunsScreen', () => {
  it('renders the list of runs', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: JOB_RUNS_QUERY,
          variables: { offset: 0, limit: 2 },
        },
        result: {
          data: {
            jobRuns: {
              results: buildRuns(),
              metadata: {
                total: 2,
              },
            },
          },
        },
      },
    ]

    renderComponent(mocks)

    await waitForLoading()

    const rows = await findAllByRole('row')

    expect(rows).toHaveLength(2)
  })

  it('can page through the list of runs', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: JOB_RUNS_QUERY,
          variables: { offset: 0, limit: 2 },
        },
        result: {
          data: {
            jobRuns: {
              results: buildRuns(),
              metadata: {
                total: 3,
              },
            },
          },
        },
      },
      {
        request: {
          query: JOB_RUNS_QUERY,
          variables: { offset: 2, limit: 2 },
        },
        result: {
          data: {
            jobRuns: {
              results: [
                buildRun({
                  id: '3',
                }),
              ],
              metadata: {
                total: 3,
              },
            },
          },
        },
      },
      {
        request: {
          query: JOB_RUNS_QUERY,
          variables: { offset: 0, limit: 2 },
        },
        result: {
          data: {
            jobRuns: {
              results: buildRuns(),
              metadata: {
                total: 3,
              },
            },
          },
        },
      },
    ]

    renderComponent(mocks)

    // Page 1
    await waitForLoading()

    expect(queryByText('1')).toBeInTheDocument()
    expect(queryByText('2')).toBeInTheDocument()
    expect(queryByText('1-2 of 3')).toBeInTheDocument()
    expect(getByRole('button', { name: /prev-page/i })).toBeDisabled()

    // Page 2
    userEvent.click(getByRole('button', { name: /next-page/i }))

    await waitForLoading()

    expect(queryByText('3')).toBeInTheDocument()
    expect(queryByText('3-3 of 3')).toBeInTheDocument()
    expect(getByRole('button', { name: /next-page/i })).toBeDisabled()

    // Page 1
    userEvent.click(getByRole('button', { name: /prev-page/i }))

    await waitForLoading()

    expect(queryByText('1')).toBeInTheDocument()
    expect(queryByText('2')).toBeInTheDocument()
    expect(queryByText('1-2 of 3')).toBeInTheDocument()
  })

  it('renders GQL errors', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: JOB_RUNS_QUERY,
          variables: { offset: 0, limit: 2 },
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
