describe('containers/BridgeSpec', () => {
  it('renders the details of the job spec and its latest runs', async () => {
    global.fetch.getOnce(`/v2/specs/${jobSpecId}`, jobSpecResponse)
  })
})
