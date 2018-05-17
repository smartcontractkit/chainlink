export default wrapper => {
  return global.fetch
    .flush()
    .then(() => wrapper.update()) // Render after AJAX request changes state
    .then(() => wrapper.update()) // Bug in enzyme, can't query conditional fragments without another sync
}
