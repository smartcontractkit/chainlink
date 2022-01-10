import * as React from 'react'

import { render, screen } from 'support/test-utils'

import { buildTaskRun } from 'support/factories/gql/fetchJobRun'
import { TaskRunsCard } from './TaskRunsCard'

const { queryByText } = screen

describe('TaskRunsCard', () => {
  it('renders a job', () => {
    const observationSource = `fetch [type=bridge name="bridge-api0"]`
    const run = buildTaskRun({
      dotID: 'fetch',
    })

    render(
      <TaskRunsCard taskRuns={[run]} observationSource={observationSource} />,
    )

    expect(queryByText(run.dotID)).toBeInTheDocument()
    expect(queryByText(run.type)).toBeInTheDocument()
    expect(queryByText(': bridge-api0')).toBeInTheDocument()
  })
})
