import * as React from 'react'

import { Route, Switch } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'
import userEvent from '@testing-library/user-event'

import { JobRunTabs, Props as JobRunTabsProps } from './JobRunTabs'

const { getByRole, queryByRole } = screen

describe('JobRunsTabs', () => {
  function renderComponent(props: Omit<JobRunTabsProps, 'classes'>) {
    renderWithRouter(
      <>
        <Switch>
          <Route exact path="/runs/:id">
            <JobRunTabs {...props} />
          </Route>

          <Route path="/runs/:id/json">
            <JobRunTabs {...props} />
          </Route>
        </Switch>
      </>,
      { initialEntries: ['/runs/1'] },
    )
  }

  it('renders the tabs', () => {
    renderComponent({ id: '1' })

    expect(queryByRole('tab', { name: 'Overview' })).toBeInTheDocument()
    expect(queryByRole('tab', { name: 'JSON' })).toBeInTheDocument()
  })

  it('switches between tabs', () => {
    renderComponent({ id: '1' })

    expect(
      queryByRole('tab', { name: 'Overview', selected: true }),
    ).toBeInTheDocument()

    userEvent.click(getByRole('tab', { name: 'JSON' }))

    expect(
      queryByRole('tab', { name: 'JSON', selected: true }),
    ).toBeInTheDocument()
  })
})
