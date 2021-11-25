import * as React from 'react'

import { GraphQLError } from 'graphql'
import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'
import { MockedProvider, MockedResponse } from '@apollo/client/testing'
import userEvent from '@testing-library/user-event'

import { UPDATE_BRIDGE_MUTATION, EditBridgeScreen } from './EditBridgeScreen'
import Notifications from 'pages/Notifications'
import { BRIDGE_QUERY } from '../Bridge/BridgeScreen'
import { buildBridgePayloadFields } from 'support/factories/gql/fetchBridge'
import { waitForLoading } from 'support/test-helpers/wait'

const { findByText, findByTestId, getByRole, getByTestId, getByText } = screen

function renderComponent(mocks: MockedResponse[]) {
  renderWithRouter(
    <>
      <Notifications />
      <Route exact path="/bridges/:id/edit">
        <MockedProvider mocks={mocks} addTypename={false}>
          <EditBridgeScreen />
        </MockedProvider>
      </Route>
    </>,
    { initialEntries: ['/bridges/bridge-api/edit'] },
  )
}

function fetchBridgeQuery(bridge: BridgePayload_Fields) {
  return {
    request: {
      query: BRIDGE_QUERY,
      variables: { id: bridge.id },
    },
    result: {
      data: {
        bridge,
      },
    },
  }
}

describe('EditBridgeScreen', () => {
  it('updates a bridge', async () => {
    const bridge = buildBridgePayloadFields()

    const mocks: MockedResponse[] = [
      fetchBridgeQuery(bridge),
      {
        request: {
          query: UPDATE_BRIDGE_MUTATION,
          variables: {
            id: bridge.id,
            input: {
              name: 'bridge-api',
              url: 'https://www.test.com',
              minimumContractPayment: '1',
              confirmations: 10,
            },
          },
        },
        result: {
          data: {
            updateBridge: {
              __typename: 'UpdateBridgeSuccess',
              bridge: {
                id: 1,
              },
            },
          },
        },
      },
    ]

    renderComponent(mocks)

    await waitForLoading()

    expect(getByText('Edit Bridge')).toBeInTheDocument()
    expect(getByTestId('bridge-form')).toBeInTheDocument()

    const urlInput = getByRole('textbox', { name: 'Bridge URL *' })
    userEvent.clear(urlInput)
    userEvent.type(urlInput, 'https://www.test.com')

    const minConfInput = getByRole('textbox', {
      name: /minimum contract payment/i,
    })
    userEvent.clear(minConfInput)
    userEvent.type(minConfInput, '1')

    userEvent.type(getByRole('spinbutton', { name: /confirmations/i }), '0')

    userEvent.click(getByRole('button', { name: /save bridge/i }))

    expect(await findByText('Successfully updated')).toBeInTheDocument()
  })

  it('renders the not found page', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: BRIDGE_QUERY,
          variables: { id: 'bridge-api' },
        },
        result: {
          data: {
            bridge: {
              __typename: 'NotFoundError',
              message: 'bridge not found',
            },
          },
        },
      },
    ]

    renderComponent(mocks)

    await waitForLoading()

    expect(await findByTestId('not-found-page')).toBeInTheDocument()
  })

  it('handles a not found error from update', async () => {
    const bridge = buildBridgePayloadFields()

    const mocks: MockedResponse[] = [
      fetchBridgeQuery(bridge),
      {
        request: {
          query: UPDATE_BRIDGE_MUTATION,
          variables: {
            id: bridge.id,
            input: {
              name: bridge.name,
              url: bridge.url,
              minimumContractPayment: bridge.minimumContractPayment,
              confirmations: bridge.confirmations,
            },
          },
        },
        result: {
          data: {
            updateBridge: {
              __typename: 'NotFoundError',
              code: 'NOT_FOUND',
              message: 'bridge not found',
            },
          },
        },
      },
    ]

    renderComponent(mocks)

    await waitForLoading()

    userEvent.click(getByRole('button', { name: /save bridge/i }))

    expect(await findByText('bridge not found')).toBeInTheDocument()
  })

  it('renders mutation GQL errors', async () => {
    const bridge = buildBridgePayloadFields()

    const mocks: MockedResponse[] = [
      fetchBridgeQuery(bridge),
      {
        request: {
          query: UPDATE_BRIDGE_MUTATION,
          variables: {
            id: bridge.id,
            input: {
              name: bridge.name,
              url: bridge.url,
              minimumContractPayment: bridge.minimumContractPayment,
              confirmations: bridge.confirmations,
            },
          },
        },
        result: {
          errors: [new GraphQLError('Mutation Error!')],
        },
      },
    ]

    renderComponent(mocks)

    await waitForLoading()

    userEvent.click(getByRole('button', { name: /save bridge/i }))

    expect(await findByText('Mutation Error!')).toBeInTheDocument()
  })
})
