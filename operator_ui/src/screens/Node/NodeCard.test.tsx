import * as React from 'react'

import { render, screen } from 'support/test-utils'
import userEvent from '@testing-library/user-event'

import { NodeCard } from './NodeCard'
import { buildNodePayloadFields } from 'support/factories/gql/fetchNode'

const { getByRole, queryByText } = screen

describe('NodeCard', () => {
  let handleDelete: jest.Mock

  function renderComponent(node: NodePayload_Fields) {
    render(<NodeCard node={node} onDelete={handleDelete} />)
  }

  beforeEach(() => {
    handleDelete = jest.fn()
  })

  it('renders a node', () => {
    const node = buildNodePayloadFields()

    renderComponent(node)

    expect(queryByText(node.id)).toBeInTheDocument()
    expect(queryByText(node.chain.id)).toBeInTheDocument()
    expect(queryByText(node.httpURL)).toBeInTheDocument()
    expect(queryByText(node.wsURL)).toBeInTheDocument()
    expect(queryByText('1 minute ago')).toBeInTheDocument()
  })

  it('calls delete', () => {
    const node = buildNodePayloadFields()

    renderComponent(node)

    userEvent.click(getByRole('button', { name: /open-menu/i }))
    userEvent.click(getByRole('menuitem', { name: /delete/i }))

    expect(handleDelete).toHaveBeenCalled()
  })
})
