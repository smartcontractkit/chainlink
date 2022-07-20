import * as React from 'react'

import { GraphQLError } from 'graphql'
import { render, screen } from 'support/test-utils'
import { MockedProvider, MockedResponse } from '@apollo/client/testing'

import { ConfigurationCard, CONFIG_QUERY } from './ConfigurationCard'
import { buildConfigItem } from 'support/factories/gql/fetchConfig'

const { findByText } = screen

function renderComponent(mocks: MockedResponse[]) {
  render(
    <MockedProvider mocks={mocks} addTypename={false}>
      <ConfigurationCard />
    </MockedProvider>,
  )
}

function fetchConfigItemsQuery(items: ReadonlyArray<Config_ItemsFields>) {
  return {
    request: {
      query: CONFIG_QUERY,
    },
    result: {
      data: {
        config: {
          items,
        },
      },
    },
  }
}

describe('ConfigurationCard', () => {
  it('renders the configuration', async () => {
    const payload = [buildConfigItem()]
    const mocks: MockedResponse[] = [fetchConfigItemsQuery(payload)]

    renderComponent(mocks)

    expect(await findByText(payload[0].key)).toBeInTheDocument()
    expect(await findByText(payload[0].value)).toBeInTheDocument()
  })

  it('renders GQL query errors', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: CONFIG_QUERY,
        },
        result: {
          errors: [new GraphQLError('Error!')],
        },
      },
    ]

    renderComponent(mocks)

    expect(await findByText('Error!')).toBeInTheDocument()
  })
})
