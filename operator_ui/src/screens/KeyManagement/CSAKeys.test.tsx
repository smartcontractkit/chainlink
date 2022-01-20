import * as React from 'react'

import { GraphQLError } from 'graphql'
import { renderWithRouter, screen } from 'support/test-utils'
import { MockedProvider, MockedResponse } from '@apollo/client/testing'
import userEvent from '@testing-library/user-event'

import { CSAKeys, CSA_KEYS_QUERY, CREATE_CSA_KEY_MUTATION } from './CSAKeys'
import { buildCSAKey, buildCSAKeys } from 'support/factories/gql/fetchCSAKeys'
import Notifications from 'pages/Notifications'
import { waitForLoading } from 'support/test-helpers/wait'

const { findByText, getByRole } = screen

function renderComponent(mocks: MockedResponse[]) {
  renderWithRouter(
    <>
      <Notifications />
      <MockedProvider mocks={mocks} addTypename={false}>
        <CSAKeys />
      </MockedProvider>
    </>,
  )
}

function fetchCSAKeysQuery(
  csaKeys: ReadonlyArray<CsaKeysPayload_ResultsFields>,
) {
  return {
    request: {
      query: CSA_KEYS_QUERY,
    },
    result: {
      data: {
        csaKeys: {
          results: csaKeys,
        },
      },
    },
  }
}

describe('CSAKeys', () => {
  it('renders the page', async () => {
    const payload = buildCSAKeys()
    const mocks: MockedResponse[] = [fetchCSAKeysQuery(payload)]

    renderComponent(mocks)

    await waitForLoading()

    expect(await findByText(payload[0].publicKey)).toBeInTheDocument()
  })

  it('renders GQL query errors', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: CSA_KEYS_QUERY,
        },
        result: {
          errors: [new GraphQLError('Error!')],
        },
      },
    ]

    renderComponent(mocks)

    expect(await findByText('Error!')).toBeInTheDocument()
  })

  it('creates a CSA Key', async () => {
    const payload = buildCSAKey()

    const mocks: MockedResponse[] = [
      fetchCSAKeysQuery([]),
      {
        request: {
          query: CREATE_CSA_KEY_MUTATION,
        },
        result: {
          data: {
            createCSAKey: {
              __typename: 'CreateCSAKeySuccess',
              csaKey: payload,
            },
          },
        },
      },
      fetchCSAKeysQuery([payload]),
    ]

    renderComponent(mocks)

    await waitForLoading()

    userEvent.click(getByRole('button', { name: /new csa key/i }))

    expect(await findByText('CSA Key created')).toBeInTheDocument()
    expect(await findByText(payload.publicKey)).toBeInTheDocument()
  })

  it('errors on create', async () => {
    const mocks: MockedResponse[] = [
      fetchCSAKeysQuery([]),
      {
        request: {
          query: CREATE_CSA_KEY_MUTATION,
        },
        result: {
          errors: [new GraphQLError('Error!')],
        },
      },
    ]

    renderComponent(mocks)

    await waitForLoading()

    userEvent.click(getByRole('button', { name: /new csa key/i }))

    expect(await findByText('Error!')).toBeInTheDocument()
  })
})
