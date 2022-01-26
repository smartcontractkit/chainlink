import React from 'react'

import { MockedProvider, MockedResponse } from '@apollo/client/testing'
import { GraphQLError } from 'graphql'
import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'test-utils'
import userEvent from '@testing-library/user-event'

import { JobsScreen, JOBS_QUERY } from './JobsScreen'
import { buildJob, buildJobs } from 'support/factories/gql/fetchJobs'
import { waitForLoading } from 'support/test-helpers/wait'

const { findAllByRole, findByText, getByRole, queryByText } = screen

function renderComponent(mocks: MockedResponse[]) {
  renderWithRouter(
    <>
      <Route exact path="/jobs">
        <MockedProvider mocks={mocks} addTypename={false}>
          <JobsScreen />
        </MockedProvider>
      </Route>

      <Route path="/jobs/new">Link Success</Route>
    </>,
    { initialEntries: ['/jobs?per=2'] },
  )
}

describe('JobsScreen', () => {
  it('renders the list of jobs', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: JOBS_QUERY,
          variables: { offset: 0, limit: 2 },
        },
        result: {
          data: {
            jobs: {
              results: buildJobs(),
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

    expect(rows).toHaveLength(3)
  })

  it('can page through the list of jobs', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: JOBS_QUERY,
          variables: { offset: 0, limit: 2 },
        },
        result: {
          data: {
            jobs: {
              results: buildJobs(),
              metadata: {
                total: 3,
              },
            },
          },
        },
      },
      {
        request: {
          query: JOBS_QUERY,
          variables: { offset: 2, limit: 2 },
        },
        result: {
          data: {
            jobs: {
              results: [
                buildJob({
                  name: 'job 3',
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
          query: JOBS_QUERY,
          variables: { offset: 0, limit: 2 },
        },
        result: {
          data: {
            jobs: {
              results: buildJobs(),
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

    expect(queryByText('job 1')).toBeInTheDocument()
    expect(queryByText('job 2')).toBeInTheDocument()
    expect(queryByText('1-2 of 3')).toBeInTheDocument()
    expect(getByRole('button', { name: /prev-page/i })).toBeDisabled()

    // Page 2
    userEvent.click(getByRole('button', { name: /next-page/i }))

    await waitForLoading()

    expect(queryByText('job 3')).toBeInTheDocument()
    expect(queryByText('3-3 of 3')).toBeInTheDocument()
    expect(getByRole('button', { name: /next-page/i })).toBeDisabled()

    // Page 1
    userEvent.click(getByRole('button', { name: /prev-page/i }))

    await waitForLoading()

    expect(queryByText('job 1')).toBeInTheDocument()
    expect(queryByText('job 2')).toBeInTheDocument()
    expect(queryByText('1-2 of 3')).toBeInTheDocument()
  })

  it('renders GQL errors', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: JOBS_QUERY,
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

  it('navigates to the new job page', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: JOBS_QUERY,
          variables: { offset: 0, limit: 2 },
        },
        result: {
          data: {
            jobs: {
              results: buildJobs(),
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

    userEvent.click(getByRole('link', { name: /new job/i }))

    expect(await findByText('Link Success')).toBeInTheDocument()
  })
})
