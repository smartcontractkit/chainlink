import '@testing-library/jest-dom'

import * as React from 'react'
import { fireEvent, render, screen } from '@testing-library/react'

import { FeedsManagerCard } from './FeedsManagerCard'
import { FeedsManager } from 'types/generated/graphql'

import { createMemoryHistory, History } from 'history'
import { Router } from 'react-router-dom'

const manager: FeedsManager = {
  id: '1',
  name: 'Chainlink Feeds Manager',
  uri: 'localhost:8080',
  publicKey: '11111111',
  jobTypes: ['FLUX_MONITOR'],
  isConnectionActive: false,
  isBootstrapPeer: false,
  bootstrapPeerMultiaddr: '/dns4/blah',
  createdAt: new Date(), // We don't display this so it doesn't matter
}

test('renders a disconnected Feeds Manager', () => {
  const history = createMemoryHistory()

  renderFeedsManagerCard(history, manager)

  expect(screen.queryByText('Chainlink Feeds Manager')).toBeInTheDocument()
  expect(screen.queryByText('localhost:8080')).toBeInTheDocument()
  expect(screen.queryByText('11111111')).toBeInTheDocument()
  expect(screen.queryByText('Disconnected')).toBeInTheDocument()
  expect(screen.queryByText('/dns4/blah')).toBeNull()
})

test('renders a connected boostrapper Feeds Manager', () => {
  // Create a new manager with connected bootstrap values
  const mgr: FeedsManager = {
    ...manager,
    jobTypes: [],
    isConnectionActive: true,
    isBootstrapPeer: true,
  }

  const history = createMemoryHistory()

  renderFeedsManagerCard(history, mgr)

  expect(screen.queryByText('Chainlink Feeds Manager')).toBeInTheDocument()
  expect(screen.queryByText('localhost:8080')).toBeInTheDocument()
  expect(screen.queryByText('11111111')).toBeInTheDocument()
  expect(screen.queryByText('FLUX_MONITOR')).toBeNull()
  expect(screen.queryByText('Connected')).toBeInTheDocument()
  expect(screen.queryByText('/dns4/blah')).toBeInTheDocument()
})

test('navigates to edit', () => {
  const history = createMemoryHistory()

  renderFeedsManagerCard(history, manager)

  fireEvent.click(screen.getByTestId('edit'))

  expect(history.location.pathname).toEqual('/feeds_manager/edit')
})

// renderFeedsManagerCard renders the component
function renderFeedsManagerCard(history: History<any>, manager: FeedsManager) {
  render(
    <Router history={history}>
      <FeedsManagerCard manager={manager} />
    </Router>,
  )
}
