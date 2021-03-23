import React from 'react'
import { mountWithProviders } from 'test-helpers/mountWithTheme'
import { syncFetch } from 'test-helpers/syncFetch'
import globPath from 'test-helpers/globPath'
import { logConfigFactory } from 'factories/jsonApiLogConfig'
import { LoggingCard } from './LoggingCard'
// import { Checkbox } from '@material-ui/core'

describe('pages/Configuration/LoggingCard', () => {
  it('renders the logging configuration card', async () => {
    const logConfig = logConfigFactory({
      level: 'info',
      sqlEnabled: false,
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
      level: 'info',
      sqlEnabled: false,
    })

    global.fetch.getOnce(globPath('/v2/log'), logConfig)
    const submit = global.fetch.patchOnce(globPath('/v2/log'), logConfig)

    const wrapper = mountWithProviders(<LoggingCard />)
    await syncFetch(wrapper)

    // Cannot figure out how to change the select and checkbox inputs for submit

    wrapper.find('form').simulate('submit')
    await syncFetch(wrapper)

    expect(submit.lastCall()[1].body).toEqual(
      '{"level":"info","sqlEnabled":false}',
    )
  })
})
