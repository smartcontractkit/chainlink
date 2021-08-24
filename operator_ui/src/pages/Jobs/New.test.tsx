import React from 'react'
import { syncFetch } from 'test-helpers/syncFetch'
import globPath from 'test-helpers/globPath'
import { mountWithProviders } from 'test-helpers/mountWithTheme'
import { ENDPOINT as TOML_CREATE_ENDPOINT } from 'api/v2/jobs'
import { Route } from 'react-router-dom'
import * as storage from 'utils/local-storage'
import New, { validate, SELECTED_FORMAT, PERSIST_SPEC } from './New'
import { ReactWrapper } from 'enzyme'
import { act } from 'react-dom/test-utils'

const fillTextarea = (component: ReactWrapper, jobSpec: string) => {
  return component.find(`textarea[name="jobSpec"]`).simulate('change', {
    target: { value: jobSpec, name: 'jobSpec' },
  })
}

describe('pages/Jobs/New', () => {
  beforeEach(() => {
    storage.remove(`${PERSIST_SPEC}`)
    storage.remove(SELECTED_FORMAT)
  })

  it('submits TOML job spec form', async () => {
    const expectedSubmit = 'foo="bar"'
    const submit = global.fetch.postOnce(globPath(TOML_CREATE_ENDPOINT), {})

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

    expect(submit.lastCall()[1].body).toEqual(
      JSON.stringify({ toml: expectedSubmit }),
    )
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
    expect(wrapper.text()).toContain(`TOML blob`)
  })

  it('saves the spec in a storage', () => {
    const wrapper = mountWithProviders(
      <Route path="/jobs/new" component={New} />,
      {
        initialEntries: [`/jobs/new`],
      },
    )

    expect(wrapper.text()).toContain(`TOML blob`)
  })

  describe('validate', () => {
    it('TOML format', () => {
      expect(
        validate({
          value: '"foo"="bar"',
        }),
      ).toEqual(true)

      expect(
        validate({
          value: '"foo""bar"',
        }),
      ).toEqual(false)
    })
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
    const tomlSpec =
      'observationSource = "" ds [type=ds]; ds_parse [type=ds_parse];  """'

    const wrapper = mountWithProviders(
      <Route path="/jobs/new" component={New} />,
      {
        initialEntries: [`/jobs/new`],
      },
    )

    fillTextarea(wrapper, tomlSpec)
    act(() => {
      jest.runAllTimers()
    })
    wrapper.update()
    expect(wrapper.text()).toContain('Tasks not found')
  })

  it('shows "Tasks not found" on empty job spec', async () => {
    const tomlSpec = 'observationSource = """  """'

    const wrapper = mountWithProviders(
      <Route path="/jobs/new" component={New} />,
      {
        initialEntries: [`/jobs/new`],
      },
    )

    fillTextarea(wrapper, tomlSpec)
    act(() => {
      jest.runAllTimers()
    })
    wrapper.update()

    expect(wrapper.text()).toContain('Tasks not found')
  })
})
