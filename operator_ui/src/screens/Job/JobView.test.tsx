import React from 'react'

import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'test-utils'
import userEvent from '@testing-library/user-event'

import { JobView, Props as JobViewProps } from './JobView'
import { buildJob, buildRun } from 'support/factories/gql/fetchJob'

const { getByRole, queryByRole, queryByText } = screen

describe('JobView', () => {
  let handleOnDelete: jest.Mock
  let handleOnRun: jest.Mock
  let handleRefetch: jest.Mock
  let handleRefetchRecentRuns: jest.Mock

  beforeEach(() => {
    handleOnDelete = jest.fn()
    handleOnRun = jest.fn()
    handleRefetch = jest.fn()
    handleRefetchRecentRuns = jest.fn()
  })

  function renderComponent(
    props: Pick<JobViewProps, 'job' | 'runsCount'>,
    initialEntry?: string,
  ) {
    const initialEntries = [`/jobs/1`]
    if (initialEntry) {
      initialEntries[0] = initialEntry
    }

    renderWithRouter(
      <>
        <Route path="/jobs/:id">
          <JobView
            onDelete={handleOnDelete}
            onRun={handleOnRun}
            refetch={handleRefetch}
            refetchRecentRuns={handleRefetchRecentRuns}
            {...props}
          />
        </Route>
      </>,
      { initialEntries },
    )
  }

  it('renders the job view', async () => {
    const job = buildJob()

    renderComponent({ job, runsCount: 1 })

    expect(queryByText(job.name)).toBeInTheDocument()
    expect(queryByText(/recent job runs/i)).toBeInTheDocument()
    expect(queryByText(/task list/i)).toBeInTheDocument()

    expect(queryByText(/overview/i)).toBeInTheDocument()
    expect(queryByText(/definition/i)).toBeInTheDocument()
    expect(queryByText(/errors/i)).toBeInTheDocument()
    expect(queryByText(/^runs/i)).toBeInTheDocument()
  })

  it('display -- for an empty job name', async () => {
    const job = buildJob({ name: '' })

    renderComponent({ job, runsCount: 1 })

    expect(queryByText('--')).toBeInTheDocument()
  })

  it('does not display the run button for other jobs', async () => {
    const job = buildJob({ name: '' })

    renderComponent({ job, runsCount: 1 })

    expect(queryByRole('button', { name: /run/i })).toBeNull()
  })

  it('runs a webhook job', async () => {
    const job = buildJob({
      type: 'webhook',
      name: 'webhook job',
      spec: {
        __typename: 'WebhookSpec',
      },
    })

    renderComponent({ job, runsCount: 1 })

    userEvent.click(getByRole('button', { name: /run/i }))
    userEvent.paste(getByRole('textbox'), '{someinput}')
    userEvent.click(getByRole('button', { name: /run job/i }))

    expect(handleOnRun).toBeCalledWith('{someinput}')
  })

  it('handles delete', async () => {
    const job = buildJob({})

    renderComponent({ job, runsCount: 1 })

    userEvent.click(getByRole('button', { name: /open-menu/i }))
    userEvent.click(getByRole('menuitem', { name: /delete/i }))
    userEvent.click(getByRole('button', { name: /confirm/i }))

    expect(handleOnDelete).toBeCalled()
  })

  it('handles refetch on next page click', async () => {
    const job = buildJob({
      runs: {
        results: [buildRun({ id: '1' })],
        metadata: {
          total: 4,
        },
      },
    })

    renderComponent({ job, runsCount: 1 }, '/jobs/1/runs?page=1&per=1')

    userEvent.click(getByRole('button', { name: /next-page/i }))

    expect(handleRefetch).toBeCalled()
  })

  it('handles refetchRecentRuns on tab change', async () => {
    const job = buildJob()

    renderComponent({ job, runsCount: 1 }, '/jobs/1/definition')

    userEvent.click(getByRole('tab', { name: 'Overview' }))

    expect(handleRefetchRecentRuns).toBeCalled()
  })
})
