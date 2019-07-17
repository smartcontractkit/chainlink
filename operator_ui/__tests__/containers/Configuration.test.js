/* eslint-env jest */
import React from 'react'
import { MemoryRouter } from 'react-router-dom'
import { Provider } from 'react-redux'
import { mount } from 'enzyme'
import createStore from 'connectors/redux'
import configurationFactory from 'factories/configuration'
import syncFetch from 'test-helpers/syncFetch'
import { ConnectedConfiguration as Configuration } from 'containers/Configuration'

const classes = {}
const mountConfiguration = props =>
  mount(
    <Provider store={createStore()}>
      <MemoryRouter>
        <Configuration classes={classes} {...props} />
      </MemoryRouter>
    </Provider>
  )

describe('containers/Configuration', () => {
  it('renders the list of configuration options', async () => {
    expect.assertions(4)

    const configurationResponse = configurationFactory({
      band: 'Major Lazer',
      singer: 'Bob Marley'
    })
    global.fetch.getOnce('/v2/config', configurationResponse)

    const wrapper = mountConfiguration()

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('BAND')
    expect(wrapper.text()).toContain('Major Lazer')
    expect(wrapper.text()).toContain('SINGER')
    expect(wrapper.text()).toContain('Bob Marley')
  })
})
