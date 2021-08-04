// import createStore from 'createStore'
// import { ConnectedIndex as Index } from 'pages/JobRuns/Index'
// import React from 'react'
// import { Provider } from 'react-redux'
// import { MemoryRouter } from 'react-router-dom'
// import clickFirstPage from 'test-helpers/clickFirstPage'
// import clickLastPage from 'test-helpers/clickLastPage'
// import clickNextPage from 'test-helpers/clickNextPage'
// import clickPreviousPage from 'test-helpers/clickPreviousPage'
// import mountWithTheme from 'test-helpers/mountWithTheme'
// import syncFetch from 'test-helpers/syncFetch'
// import globPath from 'test-helpers/globPath'

// const classes = {}
// const mountIndex = (props) =>
//   mountWithTheme(
//     <Provider store={createStore()}>
//       <MemoryRouter>
//         <Index
//           classes={classes}
//           pagePath="/jobs/:jobSpecId/runs/page"
//           {...props}
//         />
//       </MemoryRouter>
//     </Provider>,
//   )

describe('pages/JobRuns/Index', () => {
  it('should write a test when we implement this for v2 job runs', () => {})
  // const jobSpecId = 'c60b9927eeae43168ddbe92584937b1b'
  // it('renders the runs for the job spec', async () => {
  //   expect.assertions(2)
  //   const runsResponse = jsonApiJobSpecRunFactory([{ jobId: jobSpecId }])
  //   global.fetch.getOnce(globPath('/v2/runs'), runsResponse)
  //   const props = { match: { params: { jobSpecId } } }
  //   const wrapper = mountIndex(props)
  //   await syncFetch(wrapper)
  //   expect(wrapper.text()).toContain(runsResponse.data[0].id)
  //   expect(wrapper.text()).toContain('Complete')
  // })
  // it('can page through the list of runs', async () => {
  //   expect.assertions(12)
  //   const pageOneResponse = jsonApiJobSpecRunFactory(
  //     [{ id: 'ID-ON-FIRST-PAGE', jobId: jobSpecId }],
  //     3,
  //   )
  //   global.fetch.getOnce(globPath('/v2/runs'), pageOneResponse)
  //   const props = { match: { params: { jobSpecId } }, pageSize: 1 }
  //   const wrapper = mountIndex(props)
  //   await syncFetch(wrapper)
  //   expect(wrapper.text()).toContain('ID-ON-FIRST-PAGE')
  //   expect(wrapper.text()).not.toContain('ID-ON-SECOND-PAGE')
  //   const pageTwoResponse = jsonApiJobSpecRunFactory(
  //     [{ id: 'ID-ON-SECOND-PAGE', jobId: jobSpecId }],
  //     3,
  //   )
  //   global.fetch.getOnce(globPath('/v2/runs'), pageTwoResponse)
  //   clickNextPage(wrapper)
  //   await syncFetch(wrapper)
  //   expect(wrapper.text()).not.toContain('ID-ON-FIRST-PAGE')
  //   expect(wrapper.text()).toContain('ID-ON-SECOND-PAGE')
  //   global.fetch.getOnce(globPath('/v2/runs'), pageOneResponse)
  //   clickPreviousPage(wrapper)
  //   await syncFetch(wrapper)
  //   expect(wrapper.text()).toContain('ID-ON-FIRST-PAGE')
  //   expect(wrapper.text()).not.toContain('ID-ON-SECOND-PAGE')
  //   const pageThreeResponse = jsonApiJobSpecRunFactory(
  //     [{ id: 'ID-ON-THIRD-PAGE', jobId: jobSpecId }],
  //     3,
  //   )
  //   global.fetch.getOnce(globPath('/v2/runs'), pageThreeResponse)
  //   clickLastPage(wrapper)
  //   await syncFetch(wrapper)
  //   expect(wrapper.text()).toContain('ID-ON-THIRD-PAGE')
  //   expect(wrapper.text()).not.toContain('ID-ON-FIRST-PAGE')
  //   expect(wrapper.text()).not.toContain('ID-ON-SECOND-PAGE')
  //   global.fetch.getOnce(globPath('/v2/runs'), pageOneResponse)
  //   clickFirstPage(wrapper)
  //   await syncFetch(wrapper)
  //   expect(wrapper.text()).not.toContain('ID-ON-SECOND-PAGE')
  //   expect(wrapper.text()).not.toContain('ID-ON-THIRD-PAGE')
  //   expect(wrapper.text()).toContain('ID-ON-FIRST-PAGE')
  // })
  // it('displays an empty message', async () => {
  //   expect.assertions(1)
  //   const runsResponse = jsonApiJobSpecRunFactory([])
  //   await global.fetch.getOnce(globPath('/v2/runs'), runsResponse)
  //   const props = { match: { params: { jobSpecId } } }
  //   const wrapper = mountIndex(props)
  //   await syncFetch(wrapper)
  //   expect(wrapper.text()).toContain('No jobs have been run yet')
  // })
})
