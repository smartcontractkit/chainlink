import build from 'redux-object'

export default (state, id) =>
  build(state.transactions, 'items', id, { eager: true })
