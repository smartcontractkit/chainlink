import React from 'react'
import { mountWithProviders } from 'test-helpers/mountWithTheme'
import { syncFetch } from 'test-helpers/syncFetch'
import globPath from 'test-helpers/globPath'
import { logConfigFactory } from 'factories/jsonApiLogConfig'
import { LoggingCard } from './LoggingCard'
import { act } from 'react-dom/test-utils'

describe('pages/Configuration/LoggingCard', () => {
  it('renders the logging configuration card', async () => {
    const logConfig = logConfigFactory({
      serviceName: ['Global', 'IsSqlEnabled'],
      logLevel: ['info', 'false'],
    })

    global.fetch.getOnce(globPath('/v2/log'), logConfig)

    const wrapper = mountWithProviders(<LoggingCard />)
    await syncFetch(wrapper)

    expect(wrapper.find('input[name="level"]').first().props().value).toEqual(
      'info',
    )
    expect(
      wrapper.find('input[name="sqlEnabled"]').first().props().checked,
    ).toEqual(false)
  })

  it('updates the logging configuration', async () => {
    const logConfig = logConfigFactory({
      serviceName: ['Global', 'IsSqlEnabled'],
      logLevel: ['info', 'false'],
      defaultLogLevel: 'warn',
    })

    global.fetch.getOnce(globPath('/v2/log'), logConfig)
    const submit = global.fetch.patchOnce(globPath('/v2/log'), logConfig)

    const wrapper = mountWithProviders(<LoggingCard />)
    await syncFetch(wrapper)

    act(() => {
      const selectInput = wrapper.find('#select-level').first()
      const selectOnChange = selectInput.prop('onChange')
      if (selectOnChange) {
        selectOnChange({
          target: { name: 'level', value: 'debug' },
        } as any)
      }
    })

    wrapper.find('form').simulate('submit')
    await syncFetch(wrapper)

    expect(submit.lastCall()[1].body).toEqual(
      '{"defaultLogLevel":"warn","level":"debug","sqlEnabled":false}',
    )
  })
})
