export default (wrapper) => {
  const firstPage = wrapper.find('button[aria-label="First Page"]')
  firstPage.simulate('click')
}
