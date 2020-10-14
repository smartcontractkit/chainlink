export default (wrapper) => {
  // Fixes enzyme finder bug
  // https://github.com/airbnb/enzyme/issues/1233#issuecomment-385343903
  wrapper = wrapper.update()
  const nextPage = wrapper.find('button[aria-label="Next Page"]')
  nextPage.simulate('click')
}
