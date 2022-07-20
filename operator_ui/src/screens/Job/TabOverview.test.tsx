import * as React from 'react'

import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'
import userEvent from '@testing-library/user-event'

import { buildJob, buildRun } from 'support/factories/gql/fetchJob'
import { TabOverview } from './TabOverview'

const { getAllByRole, getByRole, queryByText } = screen

describe('TabOverview', () => {
  function renderComponent(job: JobPayload_Fields) {
    renderWithRouter(
      <>
        <Route exact path="/jobs/:id">
          <TabOverview job={job} />
        </Route>

        <Route exact path="/jobs/:id/runs">
          Runs Tab
        </Route>

        <Route exact path="/runs/:runId">
          Run Page
        </Route>
      </>,
      { initialEntries: ['/jobs/1'] },
    )
  }

  it('renders the recent job runs', () => {
    const job = buildJob()

    renderComponent(job)

    expect(getAllByRole('row')).toHaveLength(1)

    expect(queryByText(job.id)).toBeInTheDocument()
    expect(queryByText('1 minute ago')).toBeInTheDocument()
    expect(queryByText(/completed/i)).toBeInTheDocument()
  })

  it('navigates to the run page', () => {
    const job = buildJob()

    renderComponent(job)

    expect(getAllByRole('row')).toHaveLength(1)

    userEvent.click(getByRole('cell', { name: job.runs.results[0].id }))

    expect(queryByText('Run Page')).toBeInTheDocument()
  })

  it('renders the view more button and navigates to the runs tab', () => {
    const job = buildJob({
      runs: {
        results: [
          buildRun({ id: '1' }),
          buildRun({ id: '2' }),
          buildRun({ id: '3' }),
          buildRun({ id: '4' }),
          buildRun({ id: '5' }),
          buildRun({ id: '6' }),
        ],
        metadata: { total: 6 },
      },
    })

    renderComponent(job)

    userEvent.click(getByRole('link', { name: /view more/i }))

    expect(queryByText('Runs Tab')).toBeInTheDocument()
  })

  it('renders job tasks visualisation', async () => {
    const taskNames = ['testFetch', 'testParse', 'testMultiply']

    const job = buildJob({
      observationSource: `   ${taskNames[0]}    [type=http method=POST url="http://localhost:8001" requestData="{\\"hi\\": \\"hello\\"}"];\n    ${taskNames[1]}    [type=jsonparse path="data,result"];\n    ${taskNames[2]} [type=multiply times=100];\n    ${taskNames[0]} -\u003e ${taskNames[1]} -\u003e ${taskNames[2]};\n`,
    })

    renderComponent(job)

    expect(queryByText(taskNames[0])).toBeInTheDocument()
    expect(queryByText(taskNames[1])).toBeInTheDocument()
    expect(queryByText(taskNames[2])).toBeInTheDocument()
  })

  it('renders with an empty observation source', async () => {
    const job = buildJob({
      observationSource: '',
    })

    renderComponent(job)

    expect(queryByText('No Task Graph Found')).toBeInTheDocument()
  })

  it('renders with an invalid observation source', async () => {
    const job = buildJob({
      observationSource: 'this is totally invalid<!!@#!>',
    })

    renderComponent(job)

    expect(queryByText('Failed to parse task graph')).toBeInTheDocument()
  })
})
