import React from 'react'
import createStore from 'createStore'
import syncFetch from 'test-helpers/syncFetch'
import { mount } from 'enzyme'
import { Provider } from 'react-redux'
import { MemoryRouter } from 'react-router-dom'
import { ConnectedShow as Show } from 'containers/Bridges/Show'
import globPath from 'test-helpers/globPath'

const classes = {}
const mountShow = (props) =>
  mount(
    <Provider store={createStore()}>
      <MemoryRouter>
        <Show classes={classes} {...props} />
      </MemoryRouter>
    </Provider>,
  )

describe('containers/Bridges/Show', () => {
  it('renders the details of the bridge spec', async () => {
    expect.assertions(5)
    const response = {
      data: {
        id: 'tallbridge',
        type: 'bridges',
        attributes: {
          name: 'Tall Bridge',
          url: 'https://localhost.com:712/endpoint',
          confirmations: 9,
          outgoingToken: 'outgoingToken',
        },
      },
    }

    global.fetch.getOnce(globPath(`/v2/bridge_types/tallbridge`), response)

    const props = { match: { params: { bridgeId: 'tallbridge' } } }
    const wrapper = mountShow(props)

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('Tall Bridge')
    expect(wrapper.text()).toContain('Confirmations')
    expect(wrapper.text()).toContain('https://localhost.com:712/endpoint')
    expect(wrapper.text()).toContain('outgoingToken')
    expect(wrapper.text()).toContain('9')
  })
})
