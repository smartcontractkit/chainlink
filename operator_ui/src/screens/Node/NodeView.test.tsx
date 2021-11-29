import React from 'react'

import { render, screen } from 'test-utils'
import userEvent from '@testing-library/user-event'

import { NodeView } from './NodeView'
import { buildNodePayloadFields } from 'support/factories/gql/fetchNode'

const { getByRole, getByText } = screen

describe('NodeView', () => {
  let handleDelete: jest.Mock

  beforeEach(() => {
    handleDelete = jest.fn()
  })

  function renderComponent(node: NodePayload_Fields) {
    render(<NodeView node={node} onDelete={handleDelete} />)
  }

  it('renders the view', async () => {
    const node = buildNodePayloadFields()
    renderComponent(node)

    expect(getByRole('heading', { name: node.name })).toBeInTheDocument()
    expect(getByText('ID')).toBeInTheDocument()
    expect(getByText(node.id)).toBeInTheDocument()
  })

  it('handles opens a confirmation modal and handles delete', async () => {
    const node = buildNodePayloadFields()
    renderComponent(node)

    userEvent.click(getByRole('button', { name: /open-menu/i }))
    userEvent.click(getByRole('menuitem', { name: /delete/i }))

    expect(getByText(`Delete ${node.name}`)).toBeInTheDocument()
    expect(
      getByText(
        `This action cannot be undone and access to this page will be lost`,
      ),
    ).toBeInTheDocument()

    userEvent.click(getByRole('button', { name: /confirm/i }))

    expect(handleDelete).toHaveBeenCalled()
  })
})
