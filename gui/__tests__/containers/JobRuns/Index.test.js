import React from 'react'
import clickNextPage from 'test-helpers/clickNextPage'
import clickPreviousPage from 'test-helpers/clickPreviousPage'
import clickFirstPage from 'test-helpers/clickFirstPage'
import clickLastPage from 'test-helpers/clickLastPage'
import createStore from 'connectors/redux'
import syncFetch from 'test-helpers/syncFetch'
import jsonApiJobSpecRunFactory from 'factories/jsonApiJobSpecRuns'
import { mount } from 'enzyme'
import { Provider } from 'react-redux'
import { MemoryRouter } from 'react-router-dom'
import { ConnectedIndex as Index } from 'containers/JobRuns/Index'

const classes = {}
const mountIndex = (props) => (
  mount(
    <Provider store={createStore()}>
      <MemoryRouter>
        <Index classes={classes} {...props} />
      </MemoryRouter>
    </Provider>
  )
)

describe('containers/JobRuns/Index', () => {
  const jobSpecId = 'c60b9927eeae43168ddbe92584937b1b'

  it('renders the runs for the job spec', async () => {
    expect.assertions(4)

    const runsResponse = jsonApiJobSpecRunFactory([{
      jobId: jobSpecId
    }], jobSpecId)
    global.fetch.getOnce(
      `/v2/runs?jobSpecId=${jobSpecId}&sort=-createdAt&page=1&size=10`,
      runsResponse
    )

    const props = { match: { params: { jobSpecId: jobSpecId } } }
    const wrapper = mountIndex(props)

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain(jobSpecId)
    expect(wrapper.text()).toContain(runsResponse.data[0].id)
    expect(wrapper.text()).toContain('completed')
    expect(wrapper.text()).toContain('{"result":"value"}')
  })

  it('can page through the list of runs', async () => {
    expect.assertions(12)

    const pageOneResponse = jsonApiJobSpecRunFactory(
      [{ id: 'ID-ON-FIRST-PAGE' }],
      jobSpecId,
      3
    )
    global.fetch.getOnce(
      `/v2/runs?jobSpecId=${jobSpecId}&sort=-createdAt&page=1&size=1`,
      pageOneResponse
    )

    const props = { match: { params: { jobSpecId: jobSpecId } }, pageSize: 1 }
    const wrapper = mountIndex(props)

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('ID-ON-FIRST-PAGE')
    expect(wrapper.text()).not.toContain('ID-ON-SECOND-PAGE')

    const pageTwoResponse = jsonApiJobSpecRunFactory(
      [{ id: 'ID-ON-SECOND-PAGE' }],
      jobSpecId,
      3
    )
    global.fetch.getOnce(
      `/v2/runs?jobSpecId=${jobSpecId}&sort=-createdAt&page=2&size=1`,
      pageTwoResponse
    )
    clickNextPage(wrapper)

    await syncFetch(wrapper)
    expect(wrapper.text()).not.toContain('ID-ON-FIRST-PAGE')
    expect(wrapper.text()).toContain('ID-ON-SECOND-PAGE')

    global.fetch.getOnce(
      `/v2/runs?jobSpecId=${jobSpecId}&sort=-createdAt&page=1&size=1`,
      pageOneResponse
    )
    clickPreviousPage(wrapper)

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('ID-ON-FIRST-PAGE')
    expect(wrapper.text()).not.toContain('ID-ON-SECOND-PAGE')

    const pageThreeResponse = jsonApiJobSpecRunFactory(
      [{ id: 'ID-ON-THIRD-PAGE' }],
      jobSpecId,
      3
    )
    global.fetch.getOnce(
      `/v2/runs?jobSpecId=${jobSpecId}&sort=-createdAt&page=3&size=1`,
      pageThreeResponse
    )
    clickLastPage(wrapper)

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('ID-ON-THIRD-PAGE')
    expect(wrapper.text()).not.toContain('ID-ON-FIRST-PAGE')
    expect(wrapper.text()).not.toContain('ID-ON-SECOND-PAGE')

    global.fetch.getOnce(
      `/v2/runs?jobSpecId=${jobSpecId}&sort=-createdAt&page=1&size=1`,
      pageOneResponse
    )
    clickFirstPage(wrapper)

    await syncFetch(wrapper)
    expect(wrapper.text()).not.toContain('ID-ON-SECOND-PAGE')
    expect(wrapper.text()).not.toContain('ID-ON-THIRD-PAGE')
    expect(wrapper.text()).toContain('ID-ON-FIRST-PAGE')
  })
})
