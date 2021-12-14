import * as React from 'react'

import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'

import { NewJobFormCard, PERSIST_SPEC } from './NewJobFormCard'
import * as storage from 'utils/local-storage'

const { getByRole } = screen

describe('TaskListPreviewCard', () => {
  let handleTOMLChange: jest.Mock
  let handleOnSubmit: jest.Mock

  beforeEach(() => {
    handleTOMLChange = jest.fn()
    handleOnSubmit = jest.fn()
  })

  function renderComponent(initialPath?: string) {
    const initialEntry = initialPath || '/jobs/new'

    renderWithRouter(
      <Route path="/jobs/new">
        <NewJobFormCard
          onSubmit={handleOnSubmit}
          onTOMLChange={handleTOMLChange}
        />
      </Route>,
      { initialEntries: [initialEntry] },
    )
  }

  it('renders the card with empty input', () => {
    renderComponent()

    expect(
      getByRole('textbox', { name: /job spec \(toml\) \*/i }),
    ).toHaveTextContent('')
  })

  it('renders the card with the value from the query param', () => {
    renderComponent('/jobs/new?definition=type%20%3D%20"webhook"')

    expect(
      getByRole('textbox', { name: /job spec \(toml\) \*/i }),
    ).toHaveTextContent('type = "webhook"')
  })

  it('stores the current TOML in local storage', () => {
    renderComponent('/jobs/new?definition=type%20%3D%20"webhook"')

    expect(storage.get(PERSIST_SPEC)).toEqual('type = "webhook"')
  })

  it('renders the card with spec from local storage', () => {
    storage.set(PERSIST_SPEC, 'type = "webhook"')

    renderComponent()

    expect(
      getByRole('textbox', { name: /job spec \(toml\) \*/i }),
    ).toHaveTextContent('type = "webhook"')
  })
})
