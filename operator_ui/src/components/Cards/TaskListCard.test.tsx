import React from 'react'

import { render, screen } from '@testing-library/react'

import { TaskListCard } from './TaskListCard'
import { TaskRunStatus } from 'src/utils/taskRunStatus'

const { queryByTestId, queryByText } = screen

describe('TaskListCard', () => {
  it('renders the task graph', () => {
    render(
      <TaskListCard observationSource="ds1 [type=bridge name=voter_turnout];" />,
    )

    expect(queryByTestId('default-run-icon')).toBeInTheDocument()
    expect(queryByText('ds1')).toBeInTheDocument()
  })

  it('renders the task graph with a completed status', () => {
    render(
      <TaskListCard
        observationSource="ds1 [type=bridge name=voter_turnout];"
        attributes={{ ds1: { status: TaskRunStatus.COMPLETE } }}
      />,
    )

    expect(queryByTestId('complete-run-icon')).toBeInTheDocument()
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
