/* eslint-env jest */
import { ConnectedConfiguration as Configuration } from 'containers/Configuration'
import configurationFactory from 'factories/configuration'
import React from 'react'
import mountWithinStoreAndRouter from 'test-helpers/mountWithinStoreAndRouter'
import syncFetch from 'test-helpers/syncFetch'

const classes = {}
const mount = props => {
  return mountWithinStoreAndRouter(
    <Configuration classes={classes} {...props} />
  )
}

describe('containers/Configuration', () => {
  it('renders the list of configuration options', async () => {
    expect.assertions(4)

    const configurationResponse = configurationFactory({
      band: 'Major Lazer',
      singer: 'Bob Marley'
    })
    global.fetch.getOnce('/v2/config', configurationResponse)

    const wrapper = mount()

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('BAND')
    expect(wrapper.text()).toContain('Major Lazer')
    expect(wrapper.text()).toContain('SINGER')
    expect(wrapper.text()).toContain('Bob Marley')
  })
})
