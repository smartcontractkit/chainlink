import build from 'redux-object'

export default ({ bridges }, id) => {
  return build(bridges, 'items', id)
}
