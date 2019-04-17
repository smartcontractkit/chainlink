import build from 'redux-object'

export default (state, id) => build(state.jobs, 'items', id, { eager: true })
