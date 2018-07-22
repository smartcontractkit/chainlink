import React from 'react'
import createStore from 'connectors/redux'
import syncFetch from 'test-helpers/syncFetch'
import { mount } from 'enzyme'
import { Provider } from 'react-redux'
import { Router } from 'react-static'
import { ConnectedBridgeSpec as BridgeSpec } from 'containers/BridgeSpec'

const classes = {}
const mountBridgeSpec = (props) => (
  mount(
    <Provider store={createStore()}>
      <Router>
        <BridgeSpec classes={classes} {...props} />
      </Router>
    </Provider>
  )
)

describe('containers/BridgeSpec', () => {
  it('renders the details of the bridge spec', async () => {
    expect.assertions(4)
    const bridgeName = 'TallBridge'
    const bridgeSpecResponse = {
      name: bridgeName,
      url: 'https://localhost.com:712/endpoint',
      defaultConfirmations: 9
    }

    global.fetch.getOnce(`/v2/bridge_types/${bridgeName}`, bridgeSpecResponse)

    const props = {match: {params: {bridgeName: bridgeName}}}
    const wrapper = mountBridgeSpec(props)

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('TallBridge')
    expect(wrapper.text()).toContain('Confirmations')
    expect(wrapper.text()).toContain('https://localhost.com:712/endpoint')
    expect(wrapper.text()).toContain('9')
  })
})
