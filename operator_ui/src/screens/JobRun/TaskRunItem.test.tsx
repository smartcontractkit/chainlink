import * as React from 'react'

import { render, screen } from 'support/test-utils'

import { buildTaskRun } from 'support/factories/gql/fetchJobRun'
import { TaskRunItem, Props as TaskRunItemProps } from './TaskRunItem'

const { queryByText, queryByTestId } = screen

describe('TaskRunItemCard', () => {
  function renderComponent(itemProps: Omit<TaskRunItemProps, 'classes'>) {
    render(<TaskRunItem {...itemProps} />)
  }

  it('renders details of a completed task run', async () => {
    const run = buildTaskRun({
      error: null,
      output: "{'foo': 'bar'}",
    })

    renderComponent({ ...run, attrs: { type: run.type, path: 'result,data' } })

    // The completed icon
    expect(queryByTestId('complete-run-icon')).toBeInTheDocument()
    expect(queryByText(run.dotID)).toBeInTheDocument()
    expect(queryByText(run.type)).toBeInTheDocument()
    expect(queryByText(run.output)).toBeInTheDocument()
    expect(queryByText(': result,data')).toBeInTheDocument()
  })

  it('renders details of the errored task run', async () => {
    const run = buildTaskRun()

    renderComponent({ ...run, attrs: { type: run.type, path: 'result,data' } })

    // The error icon
    expect(queryByTestId('error-run-icon')).toBeInTheDocument()
    expect(queryByText(run.dotID)).toBeInTheDocument()
    expect(queryByText(run.type)).toBeInTheDocument()
    expect(queryByText(run.error as string)).toBeInTheDocument()
    expect(queryByText(': result,data')).toBeInTheDocument()
  })
})
