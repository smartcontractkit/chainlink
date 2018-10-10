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
    expect.assertions(6)
    const response = { data: {
      id: 'tallbridge',
      type: 'bridges',
      attributes: {
        name: 'Tall Bridge',
        url: 'https://localhost.com:712/endpoint',
        confirmations: 9,
        incomingToken: 'incomingToken',
        outgoingToken: 'outgoingToken'
      }
    }}

    global.fetch.getOnce(`/v2/bridge_types/tallbridge`, response)

    const props = {match: {params: {bridgeId: 'tallbridge'}}}
    const wrapper = mountBridgeSpec(props)

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('Tall Bridge')
    expect(wrapper.text()).toContain('Confirmations')
    expect(wrapper.text()).toContain('https://localhost.com:712/endpoint')
    expect(wrapper.text()).toContain('incomingToken')
    expect(wrapper.text()).toContain('outgoingToken')
    expect(wrapper.text()).toContain('9')
  })
})
