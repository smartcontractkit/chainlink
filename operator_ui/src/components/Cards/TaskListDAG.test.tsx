import React from 'react'

import { render, screen } from '@testing-library/react'

import { TaskListDAG } from './TaskListDAG'
import { parseDot } from 'src/utils/parseDot'
import { TaskRunStatus } from 'src/utils/taskRunStatus'

const { queryByTestId, queryByText } = screen

describe('TaskListDAG', () => {
  it('renders the task list DAG', () => {
    const graph = parseDot('digraph { ds1 [type=bridge name=voter_turnout]; }')

    render(<TaskListDAG stratify={graph} />)

    expect(queryByTestId('default-run-icon')).toBeInTheDocument()
    expect(queryByText('ds1')).toBeInTheDocument()
  })

  it('renders the task list DAG with status icon', () => {
    const graph = parseDot('digraph { ds1 [type=bridge name=voter_turnout]; }')
    const node = graph[0]
    if (node && node.attributes) {
      node.attributes.status = TaskRunStatus.COMPLETE
    }

    render(<TaskListDAG stratify={graph} />)

    expect(queryByTestId('complete-run-icon')).toBeInTheDocument()
    expect(queryByText('ds1')).toBeInTheDocument()
  })
})
