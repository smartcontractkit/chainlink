import React from 'react'

import { render, screen } from '@testing-library/react'

import { TaskRunStatusIcon } from './TaskRunStatusIcon'
import { TaskRunStatus } from 'src/utils/taskRunStatus'

const { getByTestId } = screen

describe('TaskRunStatusIcon', () => {
  // Some default values for testing
  const dimensions = { width: 20, height: 20 }

  it('renders the completed icon', () => {
    render(
      <TaskRunStatusIcon {...dimensions} status={TaskRunStatus.COMPLETE} />,
    )
    expect(getByTestId('complete-run-icon')).toBeInTheDocument()
  })

  it('renders the errored icon', () => {
    render(<TaskRunStatusIcon {...dimensions} status={TaskRunStatus.ERROR} />)
    expect(getByTestId('error-run-icon')).toBeInTheDocument()
  })

  it('renders the pending icon', () => {
    render(<TaskRunStatusIcon {...dimensions} status={TaskRunStatus.PENDING} />)
    expect(getByTestId('pending-run-icon')).toBeInTheDocument()
  })

  it('renders the default icon', () => {
    render(<TaskRunStatusIcon {...dimensions} status={TaskRunStatus.UNKNOWN} />)
    expect(getByTestId('default-run-icon')).toBeInTheDocument()
  })
})
