export default wrapper => {
  wrapper = wrapper.update() // Fixes enzyme finder bug
  const nextPage = wrapper.find('button[aria-label="Next Page"]')
  nextPage.simulate('click')
}
