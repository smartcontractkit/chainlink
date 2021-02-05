import React from 'react'
import { syncFetch } from 'test-helpers/syncFetch'
import globPath from 'test-helpers/globPath'
import { mountWithProviders } from 'test-helpers/mountWithTheme'
import { CREATE_ENDPOINT as JSON_CREATE_ENDPOINT } from 'api/v2/specs'
import { ENDPOINT as TOML_CREATE_ENDPOINT } from 'api/v2/ocrSpecs'
import { JobSpecFormats } from './utils'
import { Route } from 'react-router-dom'
import * as storage from 'utils/local-storage'
import New, { validate, SELECTED_FORMAT, PERSIST_SPEC } from './New'
import { ReactWrapper } from 'enzyme'

const switchTo = (component: ReactWrapper, format: JobSpecFormats) => {
  component.find(`input[value="${format}"]`).simulate('change', {
    target: { value: format },
  })
}

describe('pages/Jobs/New', () => {
  beforeEach(() => {
    storage.remove(`${PERSIST_SPEC}${JobSpecFormats.JSON}`)
    storage.remove(`${PERSIST_SPEC}${JobSpecFormats.TOML}`)
    storage.remove(SELECTED_FORMAT)
  })
  it('submits JSON job spec form', async () => {
    const expectedSubmit = '{"foo":"bar"}'
    const submit = global.fetch.postOnce(globPath(JSON_CREATE_ENDPOINT), {})

    const wrapper = mountWithProviders(
      <Route path="/jobs/new" component={New} />,
      {
        initialEntries: [`/jobs/new`],
      },
    )

    wrapper.find('textarea[name="jobSpec"]').simulate('change', {
      target: { value: expectedSubmit, name: 'jobSpec' },
    })

    wrapper
      .find('[data-testid="new-job-spec-submit"]')
      .first()
      .simulate('submit')

    await syncFetch(wrapper)

    expect(submit.lastCall()[1].body).toEqual(expectedSubmit)
  })

  it('submits TOML job spec form', async () => {
    const expectedSubmit = 'foo="bar"'
    const submit = global.fetch.postOnce(globPath(TOML_CREATE_ENDPOINT), {})

    const wrapper = mountWithProviders(
      <Route path="/jobs/new" component={New} />,
      {
        initialEntries: [`/jobs/new?format=${JobSpecFormats.TOML}`],
      },
    )

    wrapper.find('textarea[name="jobSpec"]').simulate('change', {
      target: { value: expectedSubmit, name: 'jobSpec' },
    })

    wrapper
      .find('[data-testid="new-job-spec-submit"]')
      .first()
      .simulate('submit')

    await syncFetch(wrapper)

    expect(submit.lastCall()[1].body).toEqual(
      JSON.stringify({ toml: expectedSubmit }),
    )
  })

  it('loads and indents JSON spec definition from search param', async () => {
    const expected = { foo: 'bar' }

    const wrapper = mountWithProviders(
      <Route path="/jobs/new" component={New} />,
      {
        initialEntries: [`/jobs/new?definition=${JSON.stringify(expected)}`],
      },
    )

    expect(wrapper.text()).toContain('"foo":"bar"')
    expect(wrapper.text()).toContain(`${JobSpecFormats.JSON} blob`)
    expect(wrapper.text()).toContain(JSON.stringify(expected))
  })

  it('loads TOML spec definition from search param', async () => {
    const expected = '"foo"="bar"'

    const wrapper = mountWithProviders(
      <Route path="/jobs/new" component={New} />,
      {
        initialEntries: [`/jobs/new?definition=${expected}`],
      },
    )

    expect(wrapper.text()).toContain('"foo"="bar"')
    expect(wrapper.text()).toContain(`${JobSpecFormats.TOML} blob`)
  })

  it('saves spec definition in storage', async () => {
    const expected = '"foo"="bar"'
    const wrapper = mountWithProviders(
      <Route path="/jobs/new" component={New} />,
      {
        initialEntries: ['/jobs/new'],
      },
    )

    wrapper.find('textarea[name="jobSpec"]').simulate('change', {
      target: { value: expected, name: 'jobSpec' },
    })

    expect(storage.get(`${PERSIST_SPEC}${JobSpecFormats.JSON}`)).toContain(
      expected,
    )
  })

  it('persists spec definitions when switch the format', async () => {
    const expectedJson = '{"foo":"bar"}'
    const expectedToml = '"foo":"bar"'

    const wrapper = mountWithProviders(
      <Route path="/jobs/new" component={New} />,
      {
        initialEntries: [`/jobs/new`],
      },
    )

    const textArea = wrapper.find('textarea[name="jobSpec"]')

    textArea.simulate('change', {
      target: { value: expectedJson, name: 'jobSpec' },
    })

    switchTo(wrapper, JobSpecFormats.TOML)
    expect(textArea.text()).not.toContain(expectedJson)

    textArea.simulate('change', {
      target: { value: expectedToml, name: 'jobSpec' },
    })

    expect(textArea.text()).toContain(expectedToml)
    switchTo(wrapper, JobSpecFormats.JSON)
    expect(textArea.text()).toContain(expectedJson)
  })

  it('saves last selected spec format in a storage', async () => {
    const wrapper = mountWithProviders(
      <Route path="/jobs/new" component={New} />,
      {
        initialEntries: [`/jobs/new?format=${JobSpecFormats.JSON}`],
      },
    )
    expect(wrapper.text()).toContain(`${JobSpecFormats.JSON} blob`)

    switchTo(wrapper, JobSpecFormats.TOML)
    expect(wrapper.text()).toContain(`${JobSpecFormats.TOML} blob`)
    expect(storage.get(SELECTED_FORMAT)).toEqual(JobSpecFormats.TOML)

    switchTo(wrapper, JobSpecFormats.JSON)
    expect(wrapper.text()).toContain(`${JobSpecFormats.JSON} blob`)
    expect(storage.get(SELECTED_FORMAT)).toEqual(JobSpecFormats.JSON)
  })

  describe('validate', () => {
    it('string', () => {
      expect(
        validate({
          format: JobSpecFormats.JSON,
          value: '           ',
        }),
      ).toEqual(false)

      expect(
        validate({
          format: JobSpecFormats.JSON,
          value: '',
        }),
      ).toEqual(false)
    })

    it('JSON format', () => {
      expect(
        validate({
          format: JobSpecFormats.JSON,
          value: '{"foo":"bar"}',
        }),
      ).toEqual(true)

      expect(
        validate({
          format: JobSpecFormats.JSON,
          value: '{"foo":"bar"',
        }),
      ).toEqual(false)
    })

    it('TOML format', () => {
      expect(
        validate({
          format: JobSpecFormats.TOML,
          value: '"foo"="bar"',
        }),
      ).toEqual(true)

      expect(
        validate({
          format: JobSpecFormats.TOML,
          value: '"foo""bar"',
        }),
      ).toEqual(false)
    })
  })
})
