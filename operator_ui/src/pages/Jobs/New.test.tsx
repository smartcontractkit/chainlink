import React from 'react'
import globPath from 'test-helpers/globPath'
import { ENDPOINT as TOML_CREATE_ENDPOINT } from 'api/v2/jobs'
import { Route } from 'react-router-dom'
import * as storage from 'utils/local-storage'
import New, { validate, SELECTED_FORMAT, PERSIST_SPEC } from './New'
import { renderWithRouter, screen } from 'support/test-utils'
import userEvent from '@testing-library/user-event'

const { getByPlaceholderText, getByRole, findByText } = screen

describe('pages/Jobs/New', () => {
  beforeEach(() => {
    storage.remove(`${PERSIST_SPEC}`)
    storage.remove(SELECTED_FORMAT)
  })

  it('submits TOML job spec form', async () => {
    const expectedSubmit = 'foo="bar"'
    const submit = global.fetch.postOnce(globPath(TOML_CREATE_ENDPOINT), {
      data: { id: 1 },
    })

    renderWithRouter(<Route path="/jobs/new" component={New} />, {
      initialEntries: [`/jobs/new`],
    })

    userEvent.paste(
      getByRole('textbox', { name: /TOML blob/i }),
      expectedSubmit,
    )

    userEvent.click(getByRole('button', { name: /Create/i }))

    expect(submit.lastCall()[1].body).toEqual(
      JSON.stringify({ toml: expectedSubmit }),
    )
  })

  it('loads TOML spec definition from search param', () => {
    const expected = '"foo"="bar"'

    renderWithRouter(<Route path="/jobs/new" component={New} />, {
      initialEntries: [`/jobs/new?definition=${expected}`],
    })

    expect(getByPlaceholderText('Paste TOML').textContent).toEqual(
      '"foo"="bar"',
    )
  })

  describe('validate', () => {
    it('TOML format', () => {
      expect(validate({ value: '"foo"="bar"' })).toEqual(true)
      expect(validate({ value: '"foo""bar"' })).toEqual(false)
    })
  })

  it('updates toml preview task list', async () => {
    const jobSpec1 =
      'observationSource = """ ds [type=ds]; ds_parse [type=ds_parse];  """'
    // const jobSpec2 =
    //   'observationSource = """ ds [type=ds]; ds_parse [type=ds_parse]; ds_multiply [type=ds_multiply]; """'

    renderWithRouter(<Route path="/jobs/new" component={New} />, {
      initialEntries: [`/jobs/new`],
    })

    userEvent.paste(getByRole('textbox', { name: /TOML blob/i }), jobSpec1)

    // Wait a second for the preview to update
    expect(
      await findByText('ds_parse', {}, { timeout: 1000 }),
    ).toBeInTheDocument()
    expect(await findByText('ds', {}, { timeout: 1000 })).toBeInTheDocument()

    // This tests the update of the textbox but it doesn't currently work
    // userEvent.paste(getByRole('textbox', { name: /TOML blob/i }), jobSpec2)

    // expect(
    //   await findByText('ds_multiple', {}, { timeout: 1000 }),
    // ).toBeInTheDocument()
    // expect(await findByText('ds_parse')).toBeInTheDocument()
    // expect(await findByText('ds')).toBeInTheDocument()
  })

  it('shows "Tasks not found" on job spec errors', async () => {
    const tomlSpec =
      'observationSource = "" ds [type=ds]; ds_parse [type=ds_parse];  """'

    renderWithRouter(<Route path="/jobs/new" component={New} />, {
      initialEntries: [`/jobs/new`],
    })

    userEvent.paste(getByRole('textbox', { name: /TOML blob/i }), tomlSpec)

    expect(await findByText('Tasks not found')).toBeInTheDocument()
  })

  it('shows "Tasks not found" on empty job spec', async () => {
    const tomlSpec = 'observationSource = """  """'

    renderWithRouter(<Route path="/jobs/new" component={New} />, {
      initialEntries: [`/jobs/new`],
    })

    userEvent.paste(getByRole('textbox', { name: /TOML blob/i }), tomlSpec)

    expect(await findByText('Tasks not found')).toBeInTheDocument()
  })
})
