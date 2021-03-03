import React from 'react'
import { syncFetch } from 'test-helpers/syncFetch'
import globPath from 'test-helpers/globPath'
import { mountWithProviders } from 'test-helpers/mountWithTheme'
import { CREATE_ENDPOINT as JSON_CREATE_ENDPOINT } from 'api/v2/specs'
import { ENDPOINT as TOML_CREATE_ENDPOINT } from 'api/v2/jobs'
import { JobSpecFormats } from './utils'
import { Route } from 'react-router-dom'
import * as storage from 'utils/local-storage'
import New, { validate, SELECTED_FORMAT, PERSIST_SPEC } from './New'
import { ReactWrapper } from 'enzyme'
import { act } from 'react-dom/test-utils'

const switchTo = (component: ReactWrapper, format: JobSpecFormats) => {
  component.find(`input[value="${format}"]`).simulate('change', {
    target: { value: format },
  })
}

const fillTextarea = (component: ReactWrapper, jobSpec: string) => {
  return component.find(`textarea[name="jobSpec"]`).simulate('change', {
    target: { value: jobSpec, name: 'jobSpec' },
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
    fillTextarea(wrapper, expectedSubmit)
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
    fillTextarea(wrapper, expectedSubmit)
    wrapper
      .find('[data-testid="new-job-spec-submit"]')
      .first()
      .simulate('submit')

    await syncFetch(wrapper)

    expect(submit.lastCall()[1].body).toEqual(
      JSON.stringify({ toml: expectedSubmit }),
    )
  })

  it('loads and indents JSON spec definition from search param', () => {
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

  it('loads TOML spec definition from search param', () => {
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

  it('saves spec definition in storage', () => {
    const expected = '"foo"="bar"'
    const wrapper = mountWithProviders(
      <Route path="/jobs/new" component={New} />,
      {
        initialEntries: ['/jobs/new'],
      },
    )
    fillTextarea(wrapper, expected)
    expect(storage.get(`${PERSIST_SPEC}${JobSpecFormats.JSON}`)).toContain(
      expected,
    )
  })

  it('persists spec definitions when switch the format', () => {
    const expectedJson = '{"foo":"bar"}'
    const expectedToml = '"foo":"bar"'

    const wrapper = mountWithProviders(
      <Route path="/jobs/new" component={New} />,
      {
        initialEntries: [`/jobs/new`],
      },
    )
    const textArea = fillTextarea(wrapper, expectedJson)

    switchTo(wrapper, JobSpecFormats.TOML)
    expect(textArea.text()).not.toContain(expectedJson)

    textArea.simulate('change', {
      target: { value: expectedToml, name: 'jobSpec' },
    })

    expect(textArea.text()).toContain(expectedToml)
    switchTo(wrapper, JobSpecFormats.JSON)
    expect(textArea.text()).toContain(expectedJson)
  })

  it('saves last selected spec format in a storage', () => {
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

  it('updates json preview task list', () => {
    const jobSpec1 = '{"tasks":[{ "type": "Httpget"}, { "type": "Jsonparse"}]}'
    const jobSpec2 =
      '{"tasks":[{ "type": "Multiply"}, { "type": "Jsonparse"}, { "type": "Httpget"}]}'

    const wrapper = mountWithProviders(
      <Route path="/jobs/new" component={New} />,
      {
        initialEntries: [`/jobs/new`],
      },
    )
    jest.useFakeTimers()
    fillTextarea(wrapper, jobSpec1)
    act(() => {
      jest.runAllTimers()
    })
    wrapper.update()
    const taskList = wrapper.find('[data-testid="task-list-item"]')
    expect(taskList.map((task: ReactWrapper) => task.text())).toEqual([
      'Httpget',
      'Jsonparse',
    ])

    jest.runAllTimers()
    fillTextarea(wrapper, jobSpec2)
    act(() => {
      jest.runAllTimers()
    })
    wrapper.update()

    const taskList2 = wrapper.find('[data-testid="task-list-item"]')
    expect(taskList2.map((task: ReactWrapper) => task.text())).toEqual([
      'Multiply',
      'Jsonparse',
      'Httpget',
    ])
  })

  it('updates toml preview task list', () => {
    const jobSpec1 =
      'observationSource = """ ds [type=ds]; ds_parse [type=ds_parse];  """'
    const jobSpec2 =
      'observationSource = """ ds [type=ds]; ds_parse [type=ds_parse]; ds_multiply [type=ds_multiply]; """'

    const wrapper = mountWithProviders(
      <Route path="/jobs/new" component={New} />,
      {
        initialEntries: [`/jobs/new`],
      },
    )
    jest.useFakeTimers()
    fillTextarea(wrapper, jobSpec1)
    act(() => {
      jest.runAllTimers()
    })
    wrapper.update()

    const taskList = wrapper.find('[data-testid^="task-list-id-"]')
    expect(
      taskList.map((task: ReactWrapper) => task.prop('data-testid')),
    ).toEqual(['task-list-id-ds_parse', 'task-list-id-ds'])

    fillTextarea(wrapper, jobSpec2)
    act(() => {
      jest.runAllTimers()
    })
    wrapper.update()

    const taskList2 = wrapper.find('[data-testid^="task-list-id-"]')
    expect(
      taskList2.map((task: ReactWrapper) => task.prop('data-testid')),
    ).toEqual([
      'task-list-id-ds_multiply',
      'task-list-id-ds_parse',
      'task-list-id-ds',
    ])
  })

  it('shows "Tasks not found" on job spec errors', () => {
    const jsonSpec = '{"tasks":[{ "type": Httpget}, { "type": "Jsonparse"}]}'
    const tomlSpec =
      'observationSource = "" ds [type=ds]; ds_parse [type=ds_parse];  """'

    const wrapper = mountWithProviders(
      <Route path="/jobs/new" component={New} />,
      {
        initialEntries: [`/jobs/new`],
      },
    )
    jest.useFakeTimers()
    fillTextarea(wrapper, jsonSpec)
    act(() => {
      jest.runAllTimers()
    })
    wrapper.update()
    expect(wrapper.text()).toContain('Tasks not found')

    fillTextarea(wrapper, tomlSpec)
    act(() => {
      jest.runAllTimers()
    })
    wrapper.update()
    expect(wrapper.text()).toContain('Tasks not found')
  })

  it('shows "Tasks not found" on empty job spec', () => {
    const jsonSpec = '{"tasks":[]}'
    const tomlSpec = 'observationSource = ""  """'

    const wrapper = mountWithProviders(
      <Route path="/jobs/new" component={New} />,
      {
        initialEntries: [`/jobs/new`],
      },
    )
    fillTextarea(wrapper, jsonSpec)
    expect(wrapper.text()).toContain('Tasks not found')

    fillTextarea(wrapper, tomlSpec)
    expect(wrapper.text()).toContain('Tasks not found')
  })
})
