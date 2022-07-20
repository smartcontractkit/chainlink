import * as React from 'react'

import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'
import userEvent from '@testing-library/user-event'

import { BridgeRow } from './BridgeRow'
import { buildBridge } from 'support/factories/gql/fetchBridges'

const { findByText, getByRole, queryByText } = screen

function renderComponent(bridge: BridgesPayload_ResultsFields) {
  renderWithRouter(
    <>
      <Route exact path="/">
        <table>
          <tbody>
            <BridgeRow bridge={bridge} />
          </tbody>
        </table>
      </Route>
      <Route exact path="/bridges/:name">
        Link Success
      </Route>
    </>,
  )
}

describe('BridgeRow', () => {
  it('renders the row', () => {
    const bridge = buildBridge()

    renderComponent(bridge)

    expect(queryByText('bridge-api')).toBeInTheDocument()
    expect(queryByText('http://bridge.com')).toBeInTheDocument()
    expect(queryByText('1')).toBeInTheDocument()
    expect(queryByText('0')).toBeInTheDocument()
  })

  it('links to the row details', async () => {
    const bridge = buildBridge()

    renderComponent(bridge)

    const link = getByRole('link', { name: /bridge-api/i })
    expect(link).toHaveAttribute('href', '/bridges/bridge-api')

    userEvent.click(link)

    expect(await findByText('Link Success')).toBeInTheDocument()
  })
})
