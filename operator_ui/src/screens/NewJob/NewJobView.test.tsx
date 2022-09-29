import React from 'react'

import { Route } from 'react-router-dom'
import { renderWithRouter, screen, waitFor } from 'test-utils'
import userEvent from '@testing-library/user-event'

import * as storage from 'utils/local-storage'
import { PERSIST_SPEC } from './NewJobFormCard/NewJobFormCard'
import { NewJobView } from './NewJobView'

const { findByText, getByRole, getByText } = screen

describe('NewJobView', () => {
  let handleOnSubmit: jest.Mock

  beforeEach(() => {
    handleOnSubmit = jest.fn()
  })

  afterEach(() => {
    // Clear the input here because the input will be saved into local storage
    // and affect the initial values of the form in subsequent tests.
    storage.remove(PERSIST_SPEC)
  })

  function renderComponent() {
    renderWithRouter(
      <Route path="/jobs/new">
        <NewJobView onSubmit={handleOnSubmit} />
      </Route>,
      { initialEntries: ['/jobs/new'] },
    )
  }

  it('renders the view', async () => {
    renderComponent()

    expect(getByText('New Job')).toBeInTheDocument()
    expect(getByText('Task List')).toBeInTheDocument()
    expect(getByText('No Task Graph Found')).toBeInTheDocument()
  })

  it('renders the task list preview from the current toml input', async () => {
    renderComponent()

    userEvent.paste(
      getByRole('textbox', { name: /job spec \(toml\) \*/i }),
      `observationSource="ds1 [type=bridge name=voter_turnout];"`,
    )

    expect(await findByText('ds1')).toBeInTheDocument()
  })

  it('submits the form', async () => {
    renderComponent()

    userEvent.paste(
      getByRole('textbox', { name: /job spec \(toml\) \*/i }),
      `observationSource="ds1 [type=bridge name=voter_turnout];"`,
    )

    userEvent.click(getByRole('button', { name: /create job/i }))

    await waitFor(() => {
      expect(handleOnSubmit).toHaveBeenCalled()
    })
  })
})
