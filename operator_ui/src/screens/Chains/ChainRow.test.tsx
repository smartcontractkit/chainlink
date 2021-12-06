import * as React from 'react'

import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'
import userEvent from '@testing-library/user-event'

import { ChainRow } from './ChainRow'
import { buildChain } from 'support/factories/gql/fetchChains'

const { findByText, getByRole, queryByText } = screen

function renderComponent(chain: ChainsPayload_ResultsFields) {
  renderWithRouter(
    <>
      <Route exact path="/">
        <table>
          <tbody>
            <ChainRow chain={chain} />
          </tbody>
        </table>
      </Route>
      <Route exact path="/chains/:id">
        Link Success
      </Route>
    </>,
  )
}

describe('ChainRow', () => {
  it('renders the row', () => {
    const chain = buildChain()

    renderComponent(chain)

    expect(queryByText('5')).toBeInTheDocument()
    expect(queryByText('true')).toBeInTheDocument()
    expect(queryByText('1 minute ago')).toBeInTheDocument()
  })

  it('links to the row details', async () => {
    const chain = buildChain()

    renderComponent(chain)

    const link = getByRole('link', { name: /5/i })
    expect(link).toHaveAttribute('href', '/chains/5')

    userEvent.click(link)

    expect(await findByText('Link Success')).toBeInTheDocument()
  })
})
