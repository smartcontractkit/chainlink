import React from 'react'
import createStore from 'connectors/redux'
import syncFetch from 'test-helpers/syncFetch'
import jsonApiJobSpecFactory from 'factories/jsonApiJobSpec'
import { mount } from 'enzyme'
import { Provider } from 'react-redux'
import { Router } from 'react-static'
import { ConnectedDefinition as Definition } from 'containers/Jobs/Definition'

const classes = {}
const mountDefinition = (props) => (
  mount(
    <Provider store={createStore()}>
      <Router>
        <Definition classes={classes} {...props} />
      </Router>
    </Provider>
  )
)

describe('containers/Jobs/Definition', () => {
  const jobSpecId = 'c60b9927eeae43168ddbe92584937b1b'

  it('displays the definition without altering the keys', async () => {
    const jobSpecResponse = jsonApiJobSpecFactory({
      id: jobSpecId,
      initiators: [{ 'type': 'web' }],
      tasks: [
        { type: 'noop', params: { asset_id_base: 'BTC' } }
      ]
    })
    global.fetch.getOnce(`/v2/specs/${jobSpecId}`, jobSpecResponse)

    const props = { match: { params: { jobSpecId: jobSpecId } } }
    const wrapper = mountDefinition(props)

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('asset_id_base')
  })
})
