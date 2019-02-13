import React from 'react'
import createStore from 'connectors/redux'
import syncFetch from 'test-helpers/syncFetch'
import { mount } from 'enzyme'
import { Provider } from 'react-redux'
import { MemoryRouter } from 'react-router-dom'
import { ConnectedShow as Show } from 'containers/Bridges/Show'

const classes = {}
const mountShow = (props) => (
  mount(
    <Provider store={createStore()}>
      <MemoryRouter>
        <Show classes={classes} {...props} />
      </MemoryRouter>
    </Provider>
  )
)

describe('containers/Bridges/Show', () => {
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
    } }

    global.fetch.getOnce(`/v2/bridge_types/tallbridge`, response)

    const props = { match: { params: { bridgeId: 'tallbridge' } } }
    const wrapper = mountShow(props)

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('Tall Bridge')
    expect(wrapper.text()).toContain('Confirmations')
    expect(wrapper.text()).toContain('https://localhost.com:712/endpoint')
    expect(wrapper.text()).toContain('incomingToken')
    expect(wrapper.text()).toContain('outgoingToken')
    expect(wrapper.text()).toContain('9')
  })
})
