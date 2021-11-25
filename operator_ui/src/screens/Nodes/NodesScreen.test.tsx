import React from 'react'

import { MockedProvider, MockedResponse } from '@apollo/client/testing'
import { GraphQLError } from 'graphql'
import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'test-utils'
import userEvent from '@testing-library/user-event'

import { NodesScreen, NODES_QUERY } from './NodesScreen'
import { buildNode, buildNodes } from 'support/factories/gql/fetchNodes'
import { waitForLoading } from 'support/test-helpers/wait'

const { findAllByRole, findByText, getByRole, queryByText } = screen

function renderComponent(mocks: MockedResponse[]) {
  renderWithRouter(
    <>
      <Route exact path="/nodes">
        <MockedProvider mocks={mocks} addTypename={false}>
          <NodesScreen />
        </MockedProvider>
      </Route>
    </>,
    { initialEntries: ['/nodes?per=2'] },
  )
}

describe('NodesScreen', () => {
  it('renders the list of nodes', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: NODES_QUERY,
          variables: { offset: 0, limit: 2 },
        },
        result: {
          data: {
            nodes: {
              results: buildNodes(),
              metadata: {
                total: 2,
              },
            },
          },
        },
      },
    ]

    renderComponent(mocks)

    expect(await findAllByRole('row')).toHaveLength(3)
  })

  it('can page through the list of nodes', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: NODES_QUERY,
          variables: { offset: 0, limit: 2 },
        },
        result: {
          data: {
            nodes: {
              results: buildNodes(),
              metadata: {
                total: 3,
              },
            },
          },
        },
      },
      {
        request: {
          query: NODES_QUERY,
          variables: { offset: 2, limit: 2 },
        },
        result: {
          data: {
            nodes: {
              results: [
                buildNode({
                  name: 'node3',
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
          query: NODES_QUERY,
          variables: { offset: 0, limit: 2 },
        },
        result: {
          data: {
            nodes: {
              results: buildNodes(),
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

    expect(queryByText('node1')).toBeInTheDocument()
    expect(queryByText('node2')).toBeInTheDocument()
    expect(queryByText('1-2 of 3')).toBeInTheDocument()
    expect(getByRole('button', { name: /prev-page/i })).toBeDisabled()

    // Page 2
    userEvent.click(getByRole('button', { name: /next-page/i }))

    await waitForLoading()

    expect(queryByText('node3')).toBeInTheDocument()
    expect(queryByText('3-3 of 3')).toBeInTheDocument()
    expect(getByRole('button', { name: /next-page/i })).toBeDisabled()

    // Page 1
    userEvent.click(getByRole('button', { name: /prev-page/i }))

    await waitForLoading()

    expect(queryByText('node1')).toBeInTheDocument()
    expect(queryByText('node2')).toBeInTheDocument()
    expect(queryByText('1-2 of 3')).toBeInTheDocument()
  })

  it('renders GQL errors', async () => {
    const mocks: MockedResponse[] = [
      {
        request: {
          query: NODES_QUERY,
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
