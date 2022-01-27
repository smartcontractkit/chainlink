import * as React from 'react'

import { GraphQLError } from 'graphql'
import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'
import { MockedProvider, MockedResponse } from '@apollo/client/testing'
import userEvent from '@testing-library/user-event'

import {
  BridgeScreen,
  BRIDGE_QUERY,
  DELETE_BRIDGE_MUTATION,
} from './BridgeScreen'
import { buildBridgePayloadFields } from 'support/factories/gql/fetchBridge'
import Notifications from 'pages/Notifications'
import { waitForLoading } from 'support/test-helpers/wait'

const { findAllByText, findByTestId, findByText, getByRole } = screen

function renderComponent(mocks: MockedResponse[]) {
  renderWithRouter(
    <>
      <Notifications />
      <Route exact path="/bridges/:id">
        <MockedProvider mocks={mocks} addTypename={false}>
          <BridgeScreen />
        </MockedProvider>
      </Route>

      <Route exact path="/bridges">
        Redirect Success
      </Route>
    </>,
    { initialEntries: ['/bridges/bridge-api'] },
  )
}

function fetchBridgeQuery(bridge: BridgePayload_Fields) {
  return {
    request: {
      query: BRIDGE_QUERY,
      variables: { id: 'bridge-api' },
    },
    result: {
      data: {
        bridge,
      },
    },
  }
}

describe('BridgeScreen', () => {
  it('renders the page', async () => {
    const payload = buildBridgePayloadFields()
    const mocks: MockedResponse[] = [fetchBridgeQuery(payload)]

    renderComponent(mocks)

    await waitForLoading()

    expect(await findAllByText('bridge-api')).toHaveLength(2)
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

  it('renders GQL query errors', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: BRIDGE_QUERY,
          variables: { id: 'bridge-api' },
        },
        result: {
          errors: [new GraphQLError('Error!')],
        },
      },
    ]

    renderComponent(mocks)

    expect(await findByText('Error: Error!')).toBeInTheDocument()
  })

  it('deletes the bridge', async () => {
    const payload = buildBridgePayloadFields()

    const mocks: MockedResponse[] = [
      fetchBridgeQuery(payload),
      {
        request: {
          query: DELETE_BRIDGE_MUTATION,
          variables: { id: 'bridge-api' },
        },
        result: {
          data: {
            deleteBridge: {
              __typename: 'DeleteBridgeSuccess',
              bridge: payload,
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

    expect(await findByText('Redirect Success')).toBeInTheDocument()
  })

  it('bridge not found when deleting', async () => {
    const payload = buildBridgePayloadFields()
    const mocks: MockedResponse[] = [
      fetchBridgeQuery(payload),
      {
        request: {
          query: DELETE_BRIDGE_MUTATION,
          variables: { id: 'bridge-api' },
        },
        result: {
          data: {
            deleteBridge: {
              __typename: 'NotFoundError',
              message: 'bridge not found',
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

    expect(await findByText('bridge not found')).toBeInTheDocument()
  })

  it('delete bridge invalid name', async () => {
    const payload = buildBridgePayloadFields()

    const mocks: MockedResponse[] = [
      fetchBridgeQuery(payload),
      {
        request: {
          query: DELETE_BRIDGE_MUTATION,
          variables: { id: 'bridge-api' },
        },
        result: {
          data: {
            deleteBridge: {
              __typename: 'DeleteBridgeInvalidNameError',
              message: 'invalid bridge name',
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

    expect(await findByText('invalid bridge name')).toBeInTheDocument()
  })

  it('delete bridge conflict', async () => {
    const payload = buildBridgePayloadFields()

    const mocks: MockedResponse[] = [
      fetchBridgeQuery(payload),
      {
        request: {
          query: DELETE_BRIDGE_MUTATION,
          variables: { id: 'bridge-api' },
        },
        result: {
          data: {
            deleteBridge: {
              __typename: 'DeleteBridgeConflictError',
              message: 'conflict error',
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

    expect(await findByText('conflict error')).toBeInTheDocument()
  })
})
