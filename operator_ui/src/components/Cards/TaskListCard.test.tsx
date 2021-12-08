import React from 'react'

import { render, screen } from '@testing-library/react'

import { TaskListCard } from './TaskListCard'

const { queryByText } = screen

describe('TaskListCard', () => {
  it('renders the task graph', () => {
    render(
      <TaskListCard observationSource="ds1 [type=bridge name=voter_turnout];" />,
    )

    expect(queryByText('ds1')).toBeInTheDocument()
  })

  it('renders a not found message', () => {
    render(<TaskListCard observationSource="" />)

    expect(queryByText('No Task Graph Found')).toBeInTheDocument()
  })

  it('renders an error message', () => {
    render(<TaskListCard observationSource="<1231!!@#>" />)

    expect(queryByText('Failed to parse task graph')).toBeInTheDocument()
  })
})
