import { act } from 'react-dom/test-utils'

export const syncFetch = async (wrapper) => {
  await act(async () => {
    await global.fetch
      .flush()
      .then(() => wrapper.update()) // Render after AJAX request changes state
      .then(() => wrapper.update()) // Bug in enzyme, can't query conditional fragments without another sync
  })

  // We want to update the wrapper after the fetch is done.
  // This is because if you wrap global.fetch.flush() in an
  // act() block they get executed in the same event loop cycle.
  wrapper.update()
}

export default (wrapper) => {
  return global.fetch
    .flush()
    .then(() => wrapper.update()) // Render after AJAX request changes state
    .then(() => wrapper.update()) // Bug in enzyme, can't query conditional fragments without another sync
}
