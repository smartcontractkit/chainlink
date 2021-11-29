import * as React from 'react'

import { renderWithRouter, screen } from 'support/test-utils'

import { buildNodes } from 'support/factories/gql/fetchNodes'
import { NodesView, Props as NodesViewProps } from './NodesView'
import userEvent from '@testing-library/user-event'

const { getAllByRole, getByRole, queryByText } = screen

function renderComponent(viewProps: NodesViewProps) {
  renderWithRouter(<NodesView {...viewProps} />)
}

describe('NodesView', () => {
  it('renders the nodes table', () => {
    const nodes = buildNodes()

    renderComponent({
      nodes,
      total: nodes.length,
      page: 1,
      pageSize: 10,
    })

    expect(getAllByRole('row')).toHaveLength(3)

    expect(queryByText('ID')).toBeInTheDocument()
    expect(queryByText('Name')).toBeInTheDocument()
    expect(queryByText('EVM Chain ID')).toBeInTheDocument()
    expect(queryByText('Created')).toBeInTheDocument()

    expect(queryByText('1')).toBeInTheDocument()
    expect(queryByText('2')).toBeInTheDocument()

    expect(queryByText('1-2 of 2'))
  })

  it('searches the nodes table', () => {
    const nodes = buildNodes()

    renderComponent({
      nodes,
      total: nodes.length,
      page: 1,
      pageSize: 10,
    })

    expect(getAllByRole('row')).toHaveLength(3)

    const searchInput = getByRole('textbox')

    // No match
    userEvent.paste(searchInput, 'doesnotmatchanything')
    expect(getAllByRole('row')).toHaveLength(2)
    expect(queryByText('No nodes found')).toBeInTheDocument()

    // By name
    userEvent.clear(searchInput)
    userEvent.paste(searchInput, 'node1')

    expect(getAllByRole('row')).toHaveLength(2)
    expect(queryByText('1')).toBeInTheDocument()

    // By id
    userEvent.clear(searchInput)
    userEvent.paste(searchInput, '1')

    expect(getAllByRole('row')).toHaveLength(2)
    expect(queryByText('1')).toBeInTheDocument()

    // By EVM ID
    userEvent.clear(searchInput)
    userEvent.paste(searchInput, '5')

    expect(getAllByRole('row')).toHaveLength(2)
    expect(queryByText('2')).toBeInTheDocument()
  })
})
