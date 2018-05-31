/* eslint-env jest */
import React from 'react'
import jobSpecFactory from 'factories/jobSpec'
import accountBalanceFactory from 'factories/accountBalance'
import syncFetch from 'test-helpers/syncFetch'
import createStore from 'connectors/redux'
import { mount } from 'enzyme'
import { Provider } from 'react-redux'
import { ConnectedJobs as Jobs } from 'containers/Jobs.js'

const classes = {}
const mountJobs = () => (
  mount(
    <Provider store={createStore()}>
      <Jobs classes={classes} />
    </Provider>
  )
)

describe('containers/Job', () => {
  it('renders the list of jobs and account balance', async () => {
    expect.assertions(6)

    const jobSpecsResponse = jobSpecFactory([{
      id: 'c60b9927eeae43168ddbe92584937b1b',
      initiators: [{'type': 'web'}],
      createdAt: '2018-05-10T00:41:54.531043837Z'
    }])
    global.fetch.getOnce('/v2/specs', jobSpecsResponse)
    const accountBalanceResponse = accountBalanceFactory(
      '10120000000000000000000',
      '7460000000000000000000'
    )
    global.fetch.getOnce('/v2/account_balance', accountBalanceResponse)

    const wrapper = mountJobs()

    await syncFetch(wrapper).then(() => {
      expect(wrapper.text()).toContain('c60b9927eeae43168ddbe92584937b1b')
      expect(wrapper.text()).toContain('web')
      expect(wrapper.text()).toContain('2018-05-10T00:41:54.531043837Z')

      expect(wrapper.text()).toContain('Ethereum10.12k')
      expect(wrapper.text()).toContain('Link7.46k')

      expect(wrapper.text()).toContain('Jobs1')
    })
  })

  it('displays an error message when the network requests fail', async () => {
    expect.assertions(3)

    global.fetch.catch(() => { throw new TypeError('Failed to fetch') })

    const wrapper = mountJobs()

    await syncFetch(wrapper).then(() => {
      expect(wrapper.text()).toContain(
        'There was an error fetching the jobs. Please reload the page.'
      )
      expect(wrapper.text()).toContain(
        'Ethereumerror fetching balance'
      )
      expect(wrapper.text()).toContain(
        'Linkerror fetching balance'
      )
    })
  })
})
