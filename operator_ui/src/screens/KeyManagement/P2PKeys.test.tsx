import * as React from 'react'

import { GraphQLError } from 'graphql'
import {
  renderWithRouter,
  screen,
  waitForElementToBeRemoved,
} from 'support/test-utils'
import { MockedProvider, MockedResponse } from '@apollo/client/testing'
import userEvent from '@testing-library/user-event'
import { waitForLoading } from 'support/test-helpers/wait'

import {
  P2PKeys,
  CREATE_P2P_KEY_MUTATION,
  DELETE_P2P_KEY_MUTATION,
} from './P2PKeys'
import { buildP2PKey, buildP2PKeys } from 'support/factories/gql/fetchP2PKeys'
import Notifications from 'pages/Notifications'
import { P2P_KEYS_QUERY } from 'src/hooks/queries/useP2PKeysQuery'

const { findByText, getByRole, queryByText } = screen

function renderComponent(mocks: MockedResponse[]) {
  renderWithRouter(
    <>
      <Notifications />
      <MockedProvider mocks={mocks} addTypename={false}>
        <P2PKeys />
      </MockedProvider>
    </>,
  )
}

function fetchP2PKeysQuery(
  bundles: ReadonlyArray<P2PKeysPayload_ResultsFields>,
) {
  return {
    request: {
      query: P2P_KEYS_QUERY,
    },
    result: {
      data: {
        p2pKeys: {
          results: bundles,
        },
      },
    },
  }
}

describe('P2PKeys', () => {
  it('renders the page', async () => {
    const payload = buildP2PKeys()
    const mocks: MockedResponse[] = [fetchP2PKeysQuery(payload)]

    renderComponent(mocks)

    await waitForLoading()

    expect(
      await findByText(`Peer ID: ${payload[0].peerID}`),
    ).toBeInTheDocument()
    expect(
      await findByText(`Peer ID: ${payload[1].peerID}`),
    ).toBeInTheDocument()
  })

  it('renders GQL query errors', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: P2P_KEYS_QUERY,
        },
        result: {
          errors: [new GraphQLError('Error!')],
        },
      },
    ]

    renderComponent(mocks)

    expect(await findByText('Error!')).toBeInTheDocument()
  })

  it('creates an P2P Key', async () => {
    const payload = buildP2PKey()

    const mocks: MockedResponse[] = [
      fetchP2PKeysQuery([]),
      {
        request: {
          query: CREATE_P2P_KEY_MUTATION,
        },
        result: {
          data: {
            createP2PKey: {
              __typename: 'CreateP2PKeySuccess',
              p2pKey: payload,
            },
          },
        },
      },
      fetchP2PKeysQuery([payload]),
    ]

    renderComponent(mocks)

    await waitForLoading()

    userEvent.click(getByRole('button', { name: /new p2p key/i }))

    expect(
      await findByText(`Successfully created P2P Key: ${payload.id}`),
    ).toBeInTheDocument()
    expect(await findByText(`Peer ID: ${payload.peerID}`)).toBeInTheDocument()
  })

  it('errors on create', async () => {
    const mocks: MockedResponse[] = [
      fetchP2PKeysQuery([]),
      {
        request: {
          query: CREATE_P2P_KEY_MUTATION,
        },
        result: {
          errors: [new GraphQLError('Error!')],
        },
      },
    ]

    renderComponent(mocks)

    await waitForLoading()

    userEvent.click(getByRole('button', { name: /new p2p key/i }))

    expect(await findByText('Error!')).toBeInTheDocument()
  })

  it('deletes an OCR Key Bundle', async () => {
    const payload = buildP2PKey()

    const mocks: MockedResponse[] = [
      fetchP2PKeysQuery([payload]),
      {
        request: {
          query: DELETE_P2P_KEY_MUTATION,
          variables: { id: payload.id },
        },
        result: {
          data: {
            deleteP2PKey: {
              __typename: 'DeleteP2PKeySuccess',
              p2pKey: payload,
            },
          },
        },
      },
      fetchP2PKeysQuery([]),
    ]

    renderComponent(mocks)

    expect(await findByText(`Peer ID: ${payload.peerID}`)).toBeInTheDocument()

    userEvent.click(getByRole('button', { name: /delete/i }))
    userEvent.click(getByRole('button', { name: /confirm/i }))

    await waitForElementToBeRemoved(getByRole('dialog'))

    expect(await findByText('Successfully deleted P2P Key')).toBeInTheDocument()

    expect(queryByText(`Peer ID: ${payload.peerID}`)).toBeNull()
  })

  it('errors on delete', async () => {
    const payload = buildP2PKey()

    const mocks: MockedResponse[] = [
      fetchP2PKeysQuery([payload]),
      {
        request: {
          query: DELETE_P2P_KEY_MUTATION,
          variables: { id: payload.id },
        },
        result: {
          errors: [new GraphQLError('Error!')],
        },
      },
    ]

    renderComponent(mocks)

    expect(await findByText(`Peer ID: ${payload.peerID}`)).toBeInTheDocument()

    userEvent.click(getByRole('button', { name: /delete/i }))
    userEvent.click(getByRole('button', { name: /confirm/i }))

    await waitForElementToBeRemoved(getByRole('dialog'))

    expect(await findByText('Error!')).toBeInTheDocument()
  })
})
