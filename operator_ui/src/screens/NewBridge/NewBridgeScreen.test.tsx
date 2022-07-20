import * as React from 'react'

import { GraphQLError } from 'graphql'
import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'
import { MockedProvider, MockedResponse } from '@apollo/client/testing'
import userEvent from '@testing-library/user-event'

import { CREATE_BRIDGE_MUTATION, NewBridgeScreen } from './NewBridgeScreen'
import Notifications from 'pages/Notifications'

const { findByText, getByRole, getByTestId, getByText } = screen

function renderComponent(mocks: MockedResponse[]) {
  renderWithRouter(
    <>
      <Notifications />
      <Route exact path="/bridges/new">
        <MockedProvider mocks={mocks} addTypename={false}>
          <NewBridgeScreen />
        </MockedProvider>
      </Route>
    </>,
    { initialEntries: ['/bridges/new'] },
  )
}

describe('NewBridgeScreen', () => {
  it('creates a bridge', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: CREATE_BRIDGE_MUTATION,
          variables: {
            input: {
              name: 'bridge1',
              url: 'https://www.test.com',
              minimumContractPayment: '0',
              confirmations: 0,
            },
          },
        },
        result: {
          data: {
            createBridge: {
              __typename: 'CreateBridgeSuccess',
              bridge: {
                id: 1,
              },
              incomingToken: 'abcdef123456',
            },
          },
        },
      },
    ]

    renderComponent(mocks)

    expect(getByText('New Bridge')).toBeInTheDocument()
    expect(getByTestId('bridge-form')).toBeInTheDocument()

    userEvent.type(getByRole('textbox', { name: 'Name *' }), 'bridge1')
    userEvent.type(
      getByRole('textbox', { name: 'Bridge URL *' }),
      'https://www.test.com',
    )

    userEvent.click(getByRole('button', { name: /create bridge/i }))

    expect(await findByText('Successfully created bridge')).toBeInTheDocument()
    expect(
      await findByText('with incoming access token: abcdef123456'),
    ).toBeInTheDocument()
  })

  it('renders mutation GQL errors', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: CREATE_BRIDGE_MUTATION,
          variables: {
            input: {
              name: 'bridge1',
              url: 'https://www.test.com',
              minimumContractPayment: '0',
              confirmations: 0,
            },
          },
        },
        result: {
          errors: [new GraphQLError('Mutation Error!')],
        },
      },
    ]

    renderComponent(mocks)

    userEvent.type(getByRole('textbox', { name: 'Name *' }), 'bridge1')
    userEvent.type(
      getByRole('textbox', { name: 'Bridge URL *' }),
      'https://www.test.com',
    )

    userEvent.click(getByRole('button', { name: /create bridge/i }))

    expect(await findByText('Mutation Error!')).toBeInTheDocument()
  })
})
