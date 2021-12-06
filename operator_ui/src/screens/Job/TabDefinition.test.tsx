import * as React from 'react'

import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'
import userEvent from '@testing-library/user-event'

import { buildJob } from 'support/factories/gql/fetchJob'
import { TabDefinition } from './TabDefinition'

const { getByRole, getByTestId, queryByText } = screen

describe('TabOverview', () => {
  function renderComponent(job: JobPayload_Fields) {
    renderWithRouter(
      <>
        <Route exact path="/jobs/:id/definition">
          <TabDefinition job={job} />
        </Route>
      </>,
      { initialEntries: ['/jobs/1/definition'] },
    )
  }

  it('renders the definition', () => {
    const job = buildJob()

    renderComponent(job)

    expect(getByTestId('definition').textContent).toMatchSnapshot()
  })

  it('can copy the definition', () => {
    // The copy package used window.prompt to copy to clipboard
    window.prompt = jest.fn()

    const job = buildJob()

    renderComponent(job)

    userEvent.click(getByRole('button', { name: /copy/i }))

    expect(queryByText(/copied!/i)).toBeInTheDocument()
    expect(window.prompt).toHaveBeenCalled()
  })
})
