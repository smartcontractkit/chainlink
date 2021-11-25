import * as React from 'react'

import userEvent from '@testing-library/user-event'
import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'

import { buildBridges } from 'support/factories/gql/fetchBridges'
import { BridgesView, Props as BridgesViewProps } from './BridgesView'

const { getAllByRole, getByText, findByText, queryByText } = screen

function renderComponent(viewProps: BridgesViewProps) {
  renderWithRouter(
    <>
      <Route exact path="/bridges">
        <BridgesView {...viewProps} />
      </Route>
      <Route exact path="/bridges/new">
        New Bridge Page
      </Route>
    </>,
    { initialEntries: ['/bridges'] },
  )
}

describe('BridgesView', () => {
  it('renders the bridges table', () => {
    const bridges = buildBridges()

    renderComponent({
      bridges,
      total: bridges.length,
      page: 1,
      pageSize: 10,
    })

    expect(getAllByRole('row')).toHaveLength(3)

    expect(queryByText('Name')).toBeInTheDocument()
    expect(queryByText('URL')).toBeInTheDocument()
    expect(queryByText('Default Confirmations')).toBeInTheDocument()
    expect(queryByText('Minimum Contract Payment')).toBeInTheDocument()

    expect(queryByText('bridge-api1')).toBeInTheDocument()
    expect(queryByText('http://bridge1.com')).toBeInTheDocument()
    expect(queryByText('1')).toBeInTheDocument()
    expect(queryByText('100')).toBeInTheDocument()

    expect(queryByText('bridge-api2')).toBeInTheDocument()
    expect(queryByText('http://bridge2.com')).toBeInTheDocument()
    expect(queryByText('2')).toBeInTheDocument()
    expect(queryByText('200')).toBeInTheDocument()

    expect(queryByText('1-2 of 2'))
  })

  it('navigates to the new bridge page', async () => {
    const bridges = buildBridges()

    renderComponent({
      bridges,
      total: bridges.length,
      page: 1,
      pageSize: 10,
    })

    userEvent.click(getByText(/new bridge/i))

    expect(await findByText('New Bridge Page')).toBeInTheDocument()
  })
})
