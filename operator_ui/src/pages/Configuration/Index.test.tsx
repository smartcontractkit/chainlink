import React from 'react'
import { Route } from 'react-router-dom'
import { mountWithProviders } from 'test-helpers/mountWithTheme'
import syncFetch from 'test-helpers/syncFetch'
import globPath from 'test-helpers/globPath'
import Configuration from 'pages/Configuration/Index'
import configurationFactory from 'factories/configuration'

describe('pages/Configuration', () => {
  it('renders the list of configuration options', async () => {
    expect.assertions(4)

    const configurationResponse = configurationFactory({
      BAND: 'Major Lazer',
      SINGER: 'Bob Marley',
    })
    global.fetch.getOnce(globPath('/v2/config'), configurationResponse)

    const wrapper = mountWithProviders(
      <Route path="/config" component={Configuration} />,
      {
        initialEntries: [`/config`],
      },
    )
    await syncFetch(wrapper)

    expect(wrapper.text()).toContain('BAND')
    expect(wrapper.text()).toContain('Major Lazer')
    expect(wrapper.text()).toContain('SINGER')
    expect(wrapper.text()).toContain('Bob Marley')
  })
})
