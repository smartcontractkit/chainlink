import build from 'redux-object'

export default (state) => {
  const address = Object.keys(state.accountBalances)[0]
  if (address) {
    return build(state, 'accountBalances', address)
  }
}
