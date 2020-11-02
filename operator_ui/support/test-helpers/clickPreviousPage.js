export default (wrapper) => {
  const nextPage = wrapper.find('button[aria-label="Previous Page"]')
  nextPage.simulate('click')
}
