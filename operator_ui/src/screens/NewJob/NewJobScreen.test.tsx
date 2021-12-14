import * as React from 'react'

import { GraphQLError } from 'graphql'
import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'
import { MockedProvider, MockedResponse } from '@apollo/client/testing'
import userEvent from '@testing-library/user-event'

import * as storage from 'utils/local-storage'
import { NewJobScreen, CREATE_JOB_MUTATION } from './NewJobScreen'
import Notifications from 'pages/Notifications'
import { PERSIST_SPEC } from './NewJobFormCard/NewJobFormCard'

const { findByRole, findByText, getByRole, getByTestId, getByText } = screen

function renderComponent(mocks: MockedResponse[]) {
  renderWithRouter(
    <>
      <Notifications />
      <Route exact path="/jobs/new">
        <MockedProvider mocks={mocks} addTypename={false}>
          <NewJobScreen />
        </MockedProvider>
      </Route>
    </>,
    { initialEntries: ['/jobs/new'] },
  )
}

describe('NewJobScreen', () => {
  afterEach(() => {
    // Clear the spec cache
    storage.remove(PERSIST_SPEC)
  })

  it('creates a job', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: CREATE_JOB_MUTATION,
          variables: {
            input: {
              TOML: 'type = "webhook"',
            },
          },
        },
        result: {
          data: {
            createJob: {
              __typename: 'CreateJobSuccess',
              job: {
                id: 1,
              },
            },
          },
        },
      },
    ]

    renderComponent(mocks)

    expect(getByText('New Job')).toBeInTheDocument()
    expect(getByTestId('job-form')).toBeInTheDocument()

    userEvent.paste(
      getByRole('textbox', { name: /job spec \(toml\) \*/i }),
      'type = "webhook"',
    )

    userEvent.click(getByRole('button', { name: /create job/i }))

    expect(await findByText('Successfully created job')).toBeInTheDocument()
    expect(await findByRole('link', { name: '1' })).toBeInTheDocument()
  })

  it('renders mutation GQL errors', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: CREATE_JOB_MUTATION,
          variables: {
            input: {
              TOML: 'type = "webhook"',
            },
          },
        },
        result: {
          errors: [new GraphQLError('Mutation Error!')],
        },
      },
    ]

    renderComponent(mocks)

    userEvent.paste(
      getByRole('textbox', { name: /job spec \(toml\) \*/i }),
      'type = "webhook"',
    )

    userEvent.click(getByRole('button', { name: /create job/i }))

    expect(await findByText('Mutation Error!')).toBeInTheDocument()
  })
})
