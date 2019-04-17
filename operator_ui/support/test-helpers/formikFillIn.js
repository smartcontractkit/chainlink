export default (wrapper, selector, value, name) => {
  const field = wrapper.find(selector)
  field.simulate('change', { target: { value, name } })
}
