import React from 'react'

import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'test-utils'

import { JobRunView } from './JobRunView'
import { buildRun, buildTaskRun } from 'support/factories/gql/fetchJobRun'

const { queryByRole, queryByText } = screen

describe('JobView', () => {
  function renderComponent(run: JobRunPayload_Fields) {
    renderWithRouter(
      <Route path="/runs/:id">
        <JobRunView run={run} />)
      </Route>,
      { initialEntries: [`/runs/${run.id}`] },
    )
  }

  it('renders the job run view with errored run ', async () => {
    const run = buildRun({
      status: 'ERRORED',
      allErrors: ['error 1'],
      job: {
        id: '10',
        name: 'job 1',
        observationSource: `fetch [type=bridge name="bridge-api0"]`,
      },
      taskRuns: [buildTaskRun({ dotID: 'fetch' })],
    })

    renderComponent(run)

    // Heading
    expect(queryByText(`Job Run #${run.id}`)).toBeInTheDocument()

    // Job Run Card
    expect(queryByText(run.id)).toBeInTheDocument()

    // Tabs
    expect(queryByRole('tab', { name: 'Overview' })).toBeInTheDocument()
    expect(queryByRole('tab', { name: 'JSON' })).toBeInTheDocument()

    // Task List Card
    expect(queryByText(/task list/i)).toBeInTheDocument()

    // Status Card
    expect(queryByText(/errored/i)).toBeInTheDocument()

    // Errors Card
    expect(queryByText(/errors/i)).toBeInTheDocument()

    // Task Runs Card
    expect(queryByText(/jsonparse/i)).toBeInTheDocument()
  })
})
