import * as React from 'react'

import { renderWithRouter, screen } from 'support/test-utils'

import { buildChains } from 'support/factories/gql/fetchChains'
import { ChainsView, Props as ChainsViewProps } from './ChainsView'
import userEvent from '@testing-library/user-event'

const { getAllByRole, getByRole, queryByText } = screen

function renderComponent(viewProps: ChainsViewProps) {
  renderWithRouter(<ChainsView {...viewProps} />)
}

describe('ChainsView', () => {
  it('renders the chains table', () => {
    const chains = buildChains()

    renderComponent({
      chains,
      total: chains.length,
      page: 1,
      pageSize: 10,
    })

    expect(getAllByRole('row')).toHaveLength(3)

    expect(queryByText('Chain ID')).toBeInTheDocument()
    expect(queryByText('Enabled')).toBeInTheDocument()
    expect(queryByText('Created')).toBeInTheDocument()

    expect(queryByText('5')).toBeInTheDocument()
    expect(queryByText('42')).toBeInTheDocument()

    expect(queryByText('1-2 of 2'))
  })

  it('searches the chains table', () => {
    const chains = buildChains()

    renderComponent({
      chains,
      total: chains.length,
      page: 1,
      pageSize: 10,
    })

    expect(getAllByRole('row')).toHaveLength(3)

    const searchInput = getByRole('textbox')

    // No match
    userEvent.paste(searchInput, 'doesnotmatchanything')
    expect(getAllByRole('row')).toHaveLength(2)
    expect(queryByText('No chains found')).toBeInTheDocument()

    // By id
    userEvent.clear(searchInput)
    userEvent.paste(searchInput, '42')

    expect(getAllByRole('row')).toHaveLength(2)
    expect(queryByText('42')).toBeInTheDocument()
  })
})
