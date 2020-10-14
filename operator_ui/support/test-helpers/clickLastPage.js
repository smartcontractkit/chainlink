export default (wrapper) => {
  const lastPage = wrapper.find('button[aria-label="Last Page"]')
  lastPage.simulate('click')
}
