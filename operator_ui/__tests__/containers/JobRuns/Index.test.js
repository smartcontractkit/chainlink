import React from 'react'
import clickNextPage from 'test-helpers/clickNextPage'
import clickPreviousPage from 'test-helpers/clickPreviousPage'
import clickFirstPage from 'test-helpers/clickFirstPage'
import clickLastPage from 'test-helpers/clickLastPage'
import createStore from 'connectors/redux'
import syncFetch from 'test-helpers/syncFetch'
import jsonApiJobSpecRunFactory from 'factories/jsonApiJobSpecRuns'
import mountWithTheme from 'test-helpers/mountWithTheme'
import { Provider } from 'react-redux'
import { MemoryRouter } from 'react-router-dom'
import { ConnectedIndex as Index } from 'containers/JobRuns/Index'

const classes = {}
const mountIndex = props =>
  mountWithTheme(
    <Provider store={createStore()}>
      <MemoryRouter>
        <Index
          classes={classes}
          pagePath="/jobs/:jobSpecId/runs/page"
          {...props}
        />
      </MemoryRouter>
    </Provider>
  )

describe('containers/JobRuns/Index', () => {
  const jobSpecId = 'c60b9927eeae43168ddbe92584937b1b'

  it('renders the runs for the job spec', async () => {
    expect.assertions(2)

    const runsResponse = jsonApiJobSpecRunFactory(
      [
        {
          jobId: jobSpecId
        }
      ],
      jobSpecId
    )
    global.fetch.getOnce(
      `/v2/runs?sort=-createdAt&page=1&size=25&jobSpecId=${jobSpecId}`,
      runsResponse
    )

    const props = { match: { params: { jobSpecId: jobSpecId } } }
    const wrapper = mountIndex(props)

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain(runsResponse.data[0].id)
    expect(wrapper.text()).toContain('Complete')
  })

  it('can page through the list of runs', async () => {
    expect.assertions(12)

    const pageOneResponse = jsonApiJobSpecRunFactory(
      [{ id: 'ID-ON-FIRST-PAGE' }],
      jobSpecId,
      3
    )
    global.fetch.getOnce(
      `/v2/runs?sort=-createdAt&page=1&size=1&jobSpecId=${jobSpecId}`,
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
      `/v2/runs?sort=-createdAt&page=2&size=1&jobSpecId=${jobSpecId}`,
      pageTwoResponse
    )
    clickNextPage(wrapper)

    await syncFetch(wrapper)
    expect(wrapper.text()).not.toContain('ID-ON-FIRST-PAGE')
    expect(wrapper.text()).toContain('ID-ON-SECOND-PAGE')

    global.fetch.getOnce(
      `/v2/runs?sort=-createdAt&page=1&size=1&jobSpecId=${jobSpecId}`,
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
      `/v2/runs?sort=-createdAt&page=3&size=1&jobSpecId=${jobSpecId}`,
      pageThreeResponse
    )
    clickLastPage(wrapper)

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('ID-ON-THIRD-PAGE')
    expect(wrapper.text()).not.toContain('ID-ON-FIRST-PAGE')
    expect(wrapper.text()).not.toContain('ID-ON-SECOND-PAGE')

    global.fetch.getOnce(
      `/v2/runs?sort=-createdAt&page=1&size=1&jobSpecId=${jobSpecId}`,
      pageOneResponse
    )
    clickFirstPage(wrapper)

    await syncFetch(wrapper)
    expect(wrapper.text()).not.toContain('ID-ON-SECOND-PAGE')
    expect(wrapper.text()).not.toContain('ID-ON-THIRD-PAGE')
    expect(wrapper.text()).toContain('ID-ON-FIRST-PAGE')
  })

  it('displays an empty message', async () => {
    expect.assertions(1)

    const runsResponse = jsonApiJobSpecRunFactory([], jobSpecId)
    global.fetch.getOnce(
      `/v2/runs?sort=-createdAt&page=1&size=25&jobSpecId=${jobSpecId}`,
      runsResponse
    )

    const props = { match: { params: { jobSpecId: jobSpecId } } }
    const wrapper = mountIndex(props)

    await syncFetch(wrapper)
    expect(wrapper.text()).toContain('No jobs have been run yet')
  })
})
