export default wrapper => {
  const nextPage = wrapper.find('button[aria-label="Next Page"]')
  nextPage.simulate('click')
}
