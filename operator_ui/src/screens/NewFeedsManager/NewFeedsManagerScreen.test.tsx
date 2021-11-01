import '@testing-library/jest-dom'

import * as React from 'react'
import {
  render,
  screen,
  waitForElementToBeRemoved,
} from '@testing-library/react'
import { MockedProvider } from '@apollo/client/testing'

import { FETCH_FEEDS_MANAGERS } from 'src/hooks/useFetchFeedsManager'
import { NewFeedsManagerScreen } from './NewFeedsManagerScreen'

import { createMemoryHistory } from 'history'
import { Route, Router } from 'react-router-dom'
import { Provider } from 'react-redux'
import createStore from 'createStore'

test('renders the page', async () => {
  const mocks: any = [
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

  const history = createMemoryHistory()

  render(
    <MockedProvider
      mocks={mocks}
      addTypename={false}
      defaultOptions={{
        watchQuery: { fetchPolicy: 'no-cache' },
        query: { fetchPolicy: 'no-cache' },
      }}
    >
      <Provider store={createStore()}>
        <Router history={history}>
          <NewFeedsManagerScreen />
        </Router>
      </Provider>
    </MockedProvider>,
  )

  await waitForElementToBeRemoved(() => screen.queryByText('Loading...'))

  expect(screen.queryByText('Register Feeds Manager')).toBeInTheDocument()
  expect(screen.queryByTestId('feeds-manager-form')).toBeInTheDocument()
})

test('redirects when a manager exists', async () => {
  const mocks: any = [
    {
      request: {
        query: FETCH_FEEDS_MANAGERS,
      },
      result: {
        data: {
          feedsManagers: {
            results: [
              {
                id: '1',
                name: 'Chainlink Feeds Manager',
                uri: 'localhost:8080',
                publicKey: '11111111',
                jobTypes: ['FLUX_MONITOR'],
                isConnectionActive: false,
                isBootstrapPeer: false,
                bootstrapPeerMultiaddr: undefined,
                createdAt: new Date(), // We don't display this so it doesn't matter
              },
            ],
          },
        },
      },
    },
  ]

  const history = createMemoryHistory()

  render(
    <MockedProvider
      mocks={mocks}
      addTypename={false}
      defaultOptions={{
        watchQuery: { fetchPolicy: 'no-cache' },
        query: { fetchPolicy: 'no-cache' },
      }}
    >
      <Provider store={createStore()}>
        <Router history={history}>
          <Route exact path="/">
            <NewFeedsManagerScreen />
          </Route>
          <Route path="/feeds_manager">Feeds Manager Page</Route>
        </Router>
      </Provider>
    </MockedProvider>,
  )

  await waitForElementToBeRemoved(() => screen.queryByText('Loading...'))

  expect(history.location.pathname).toEqual('/feeds_manager')
  expect(screen.queryByText('Feeds Manager Page')).toBeInTheDocument()
})
