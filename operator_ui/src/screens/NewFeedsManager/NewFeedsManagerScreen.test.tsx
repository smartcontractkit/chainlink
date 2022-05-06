import * as React from 'react'

import { GraphQLError } from 'graphql'
import { Route } from 'react-router-dom'
import {
  renderWithRouter,
  screen,
  waitForElementToBeRemoved,
} from 'support/test-utils'
import userEvent from '@testing-library/user-event'
import { MockedProvider, MockedResponse } from '@apollo/client/testing'

import { buildFeedsManager } from 'support/factories/gql/fetchFeedsManagers'
import { FEEDS_MANAGERS_QUERY } from 'src/hooks/queries/useFeedsManagersQuery'
import {
  CREATE_FEEDS_MANAGER_MUTATION,
  NewFeedsManagerScreen,
} from './NewFeedsManagerScreen'
import Notifications from 'pages/Notifications'

const { findByTestId, findByText, getByRole } = screen

function renderComponent(mocks: MockedResponse[]) {
  renderWithRouter(
    <>
      <Notifications />
      <Route exact path="/">
        <MockedProvider mocks={mocks} addTypename={false}>
          <NewFeedsManagerScreen />
        </MockedProvider>
      </Route>

      <Route path="/feeds_manager">Redirect Success</Route>
    </>,
  )
}

describe('NewFeedsManagerScreen', () => {
  it('renders the page', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: FEEDS_MANAGERS_QUERY,
        },
        result: {
          data: {
            feedsManagers: {
              results: [],
            },
          },
        },
      },
    ]

    renderComponent(mocks)

    await waitForElementToBeRemoved(() => screen.queryByRole('progressbar'))

    expect(await findByText('Register Feeds Manager')).toBeInTheDocument()
    expect(await findByTestId('feeds-manager-form')).toBeInTheDocument()
  })

  it('redirects when a manager exists', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: FEEDS_MANAGERS_QUERY,
        },
        result: {
          data: {
            feedsManagers: {
              results: [buildFeedsManager()],
            },
          },
        },
      },
    ]

    renderComponent(mocks)

    await waitForElementToBeRemoved(() => screen.queryByRole('progressbar'))

    expect(await findByText('Redirect Success')).toBeInTheDocument()
  })

  it('submits the form', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: FEEDS_MANAGERS_QUERY,
        },
        result: {
          data: {
            feedsManagers: {
              results: [],
            },
          },
        },
      },
      {
        request: {
          query: CREATE_FEEDS_MANAGER_MUTATION,
          variables: {
            input: {
              name: 'Chainlink Feeds Manager',
              uri: 'localhost:8080',
              publicKey: '1111',
            },
          },
        },
        result: {
          data: {
            createFeedsManager: {
              __typename: 'CreateFeedsManagerSuccess',
              feedsManager: buildFeedsManager(),
            },
          },
        },
      },
      {
        request: {
          query: FEEDS_MANAGERS_QUERY,
        },
        result: {
          data: {
            feedsManagers: {
              results: [buildFeedsManager()],
            },
          },
        },
      },
    ]

    renderComponent(mocks)

    await waitForElementToBeRemoved(() => screen.queryByRole('progressbar'))

    // Note: The name input has a default value so we don't have to set it
    userEvent.type(getByRole('textbox', { name: 'URI *' }), 'localhost:8080')
    userEvent.type(getByRole('textbox', { name: 'Public Key *' }), '1111')

    userEvent.click(getByRole('button', { name: /submit/i }))

    expect(await findByText('Feeds Manager Created')).toBeInTheDocument()
    expect(await findByText('Redirect Success')).toBeInTheDocument()
  })

  it('handles input errors', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: FEEDS_MANAGERS_QUERY,
        },
        result: {
          data: {
            feedsManagers: {
              results: [],
            },
          },
        },
      },
      {
        request: {
          query: CREATE_FEEDS_MANAGER_MUTATION,
          variables: {
            input: {
              name: 'Chainlink Feeds Manager',
              uri: 'localhost:8080',
              publicKey: '1111',
            },
          },
        },
        result: {
          data: {
            createFeedsManager: {
              __typename: 'InputErrors',
              errors: [
                {
                  code: 'INPUT_ERROR',
                  message: 'invalid hex value',
                  path: 'input/publicKey',
                },
              ],
            },
          },
        },
      },
      {
        request: {
          query: FEEDS_MANAGERS_QUERY,
        },
        result: {
          data: {
            feedsManagers: {
              results: [],
            },
          },
        },
      },
    ]

    renderComponent(mocks)

    await waitForElementToBeRemoved(() => screen.queryByRole('progressbar'))

    // Note: The name input has a default value so we don't have to set it
    userEvent.type(getByRole('textbox', { name: 'URI *' }), 'localhost:8080')
    userEvent.type(getByRole('textbox', { name: 'Public Key *' }), '1111')

    userEvent.click(getByRole('button', { name: /submit/i }))

    expect(await findByText('Invalid Input')).toBeInTheDocument()
    expect(await findByTestId('publicKey-helper-text')).toHaveTextContent(
      'invalid hex value',
    )
  })

  it('renders query GQL errors', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: FEEDS_MANAGERS_QUERY,
        },
        result: {
          errors: [new GraphQLError('Error!')],
        },
      },
    ]

    renderComponent(mocks)

    expect(await findByText('Error: Error!')).toBeInTheDocument()
  })

  it('renders mutation GQL errors', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: FEEDS_MANAGERS_QUERY,
        },
        result: {
          data: {
            feedsManagers: {
              results: [],
            },
          },
        },
      },
      {
        request: {
          query: CREATE_FEEDS_MANAGER_MUTATION,
          variables: {
            input: {
              name: 'Chainlink Feeds Manager',
              uri: 'localhost:8080',
              publicKey: '1111',
            },
          },
        },
        result: {
          errors: [new GraphQLError('Mutation Error!')],
        },
      },
    ]

    renderComponent(mocks)

    await waitForElementToBeRemoved(() => screen.queryByRole('progressbar'))

    userEvent.type(getByRole('textbox', { name: 'URI *' }), 'localhost:8080')
    userEvent.type(getByRole('textbox', { name: 'Public Key *' }), '1111')

    userEvent.click(getByRole('button', { name: /submit/i }))

    expect(await findByText('Mutation Error!')).toBeInTheDocument()
  })
})
