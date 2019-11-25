export default (wrapper, selector, value) => {
  wrapper.find(selector).simulate('change', { target: { value } })
}
