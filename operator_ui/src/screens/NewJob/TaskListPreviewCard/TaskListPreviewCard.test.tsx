import * as React from 'react'

import { render, screen } from 'support/test-utils'

import { TaskListPreviewCard } from './TaskListPreviewCard'

const { queryByTestId, queryByText } = screen

describe('TaskListPreviewCard', () => {
  function renderComponent(toml: string) {
    render(<TaskListPreviewCard toml={toml} />)
  }

  it('renders the card', () => {
    renderComponent(`observationSource="ds1 [type=bridge name=voter_turnout];"`)

    expect(queryByText('Task List')).toBeInTheDocument()
    expect(queryByTestId('default-run-icon')).toBeInTheDocument()
    expect(queryByText('ds1')).toBeInTheDocument()
  })

  it('has an empty TOML string', () => {
    renderComponent('')

    expect(queryByText('No Task Graph Found')).toBeInTheDocument()
  })

  it('has an invalid TOML string', () => {
    renderComponent('invalidstring')

    expect(queryByText('No Task Graph Found')).toBeInTheDocument()
  })
})
