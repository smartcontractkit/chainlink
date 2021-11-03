import '@testing-library/jest-dom'

import * as React from 'react'
import {
  render,
  screen,
  waitFor,
  waitForElementToBeRemoved,
} from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MockedProvider, MockedResponse } from '@apollo/client/testing'

import { FETCH_FEEDS_MANAGERS } from 'src/hooks/useFetchFeedsManager'
import {
  UPDATE_FEEDS_MANAGER,
  EditFeedsManagerScreen,
} from './EditFeedsManagerScreen'
import { RedirectTestRoute } from 'support/test-helpers/RedirectTestRoute'

import { MemoryRouter, Route } from 'react-router-dom'
import { Provider } from 'react-redux'
import createStore from 'createStore'
import { buildFeedsManager } from 'support/test-helpers/factories/feedsManager'

test('renders the page', async () => {
  const mocks: MockedResponse[] = [
    {
      request: {
        query: FETCH_FEEDS_MANAGERS,
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

  render(
    <MockedProvider mocks={mocks} addTypename={false}>
      <Provider store={createStore()}>
        <MemoryRouter>
          <EditFeedsManagerScreen />
        </MemoryRouter>
      </Provider>
    </MockedProvider>,
  )

  await waitForElementToBeRemoved(() => screen.queryByRole('progressbar'))

  expect(screen.queryByText('Edit Feeds Manager')).toBeInTheDocument()
  expect(screen.queryByTestId('feeds-manager-form')).toBeInTheDocument()
})

test('redirects when a manager does not exist', async () => {
  const mocks: MockedResponse[] = [
    {
      request: {
        query: FETCH_FEEDS_MANAGERS,
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

  render(
    <MockedProvider mocks={mocks} addTypename={false}>
      <Provider store={createStore()}>
        <MemoryRouter>
          <Route exact path="/">
            <EditFeedsManagerScreen />
          </Route>

          <RedirectTestRoute
            path="/feeds_manager/new"
            message="Feeds Manager Page"
          />
        </MemoryRouter>
      </Provider>
    </MockedProvider>,
  )

  await waitForElementToBeRemoved(() => screen.queryByRole('progressbar'))

  expect(screen.queryByText('Feeds Manager Page')).toBeInTheDocument()
})

test('submits the form', async () => {
  const { getByRole, getByTestId } = screen
  const mocks: MockedResponse[] = [
    {
      request: {
        query: FETCH_FEEDS_MANAGERS,
      },
      result: {
        data: {
          feedsManagers: {
            results: [buildFeedsManager()],
          },
        },
      },
    },
    {
      request: {
        query: UPDATE_FEEDS_MANAGER,
        variables: {
          input: {
            name: 'Chainlink updated',
            uri: 'localhost:80812',
            publicKey: '11112',
            jobTypes: [],
            isBootstrapPeer: false,
            bootstrapPeerMultiaddr: undefined,
          },
        },
      },
      result: {
        data: {
          createFeedsManager: {
            feedsManager: buildFeedsManager({
              name: 'Chainlink updated',
              uri: 'localhost:80812',
              publicKey: '11112',
              jobTypes: [],
            }),
          },
        },
      },
    },
    {
      request: {
        query: FETCH_FEEDS_MANAGERS,
      },
      result: {
        data: {
          feedsManagers: {
            results: [
              buildFeedsManager({
                name: 'Chainlink updated',
                uri: 'localhost:80812',
                publicKey: '11112',
                jobTypes: [],
              }),
            ],
          },
        },
      },
    },
  ]

  render(
    <MockedProvider mocks={mocks} addTypename={false}>
      <Provider store={createStore()}>
        <MemoryRouter>
          <Route exact path="/">
            <EditFeedsManagerScreen />
          </Route>

          <RedirectTestRoute
            path="/feeds_manager"
            message="Feeds Manager Page"
          />
        </MemoryRouter>
      </Provider>
    </MockedProvider>,
  )

  await waitForElementToBeRemoved(() => screen.queryByRole('progressbar'))

  userEvent.type(getByRole('textbox', { name: 'Name *' }), ' updated')
  userEvent.type(getByRole('textbox', { name: 'URI *' }), 'localhost:80812')
  userEvent.type(getByRole('textbox', { name: 'Public Key *' }), '2')
  userEvent.click(getByRole('checkbox', { name: 'Flux Monitor' }))

  userEvent.click(getByTestId('create-submit'))

  await waitFor(() =>
    expect(screen.queryByText('Feeds Manager Page')).toBeInTheDocument(),
  )
})
