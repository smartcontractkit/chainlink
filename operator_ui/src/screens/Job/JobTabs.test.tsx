import * as React from 'react'

import { Route, Switch } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'
import userEvent from '@testing-library/user-event'

import { JobTabs, Props as JobTabsProps } from './JobTabs'

const { getByRole, queryByRole } = screen

describe('JobTabs', () => {
  let handleRefetchRecentRuns: jest.Mock

  function renderComponent(props: Omit<JobTabsProps, 'classes'>) {
    renderWithRouter(
      <>
        <Switch>
          <Route exact path="/jobs/:id">
            <JobTabs {...props} />
          </Route>

          <Route path="/jobs/:id/definition">
            <JobTabs {...props} />
          </Route>

          <Route path="/jobs/:id/errors">
            <JobTabs {...props} />
          </Route>

          <Route path="/jobs/:id/runs">
            <JobTabs {...props} />
          </Route>
        </Switch>
      </>,
      { initialEntries: ['/jobs/1'] },
    )
  }

  beforeEach(() => {
    handleRefetchRecentRuns = jest.fn()
  })

  it('renders the tabs', () => {
    renderComponent({
      id: '1',
      errorsCount: 1,
      runsCount: 2,
      refetchRecentRuns: handleRefetchRecentRuns,
    })

    expect(queryByRole('tab', { name: 'Overview' })).toBeInTheDocument()
    expect(queryByRole('tab', { name: 'Definition' })).toBeInTheDocument()
    expect(queryByRole('tab', { name: 'Errors 1' })).toBeInTheDocument()
    expect(queryByRole('tab', { name: 'Runs 2' })).toBeInTheDocument()
  })

  it('switches between tabs', () => {
    renderComponent({
      id: '1',
      errorsCount: 0,
      runsCount: 0,
      refetchRecentRuns: handleRefetchRecentRuns,
    })

    expect(
      queryByRole('tab', { name: 'Overview', selected: true }),
    ).toBeInTheDocument()

    userEvent.click(getByRole('tab', { name: 'Definition' }))

    expect(
      queryByRole('tab', { name: 'Definition', selected: true }),
    ).toBeInTheDocument()

    userEvent.click(getByRole('tab', { name: 'Errors 0' }))

    expect(
      queryByRole('tab', { name: 'Errors 0', selected: true }),
    ).toBeInTheDocument()

    userEvent.click(getByRole('tab', { name: 'Runs 0' }))

    expect(
      queryByRole('tab', { name: 'Runs 0', selected: true }),
    ).toBeInTheDocument()
  })

  it('switches to overview fetches the recent runs', () => {
    renderComponent({
      id: '1',
      errorsCount: 0,
      runsCount: 0,
      refetchRecentRuns: handleRefetchRecentRuns,
    })

    userEvent.click(getByRole('tab', { name: 'Definition' }))
    userEvent.click(getByRole('tab', { name: 'Overview' }))

    expect(handleRefetchRecentRuns).toHaveBeenCalled()
  })
})
