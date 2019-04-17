const initialState = {
  items: {},
  currentPage: null,
  recentlyCreated: null,
  count: 0
}

export const UPSERT_JOBS = 'UPSERT_JOBS'
export const UPSERT_RECENTLY_CREATED_JOBS = 'UPSERT_RECENTLY_CREATED_JOBS'
export const UPSERT_JOB = 'UPSERT_JOB'

export default (state = initialState, action = {}) => {
  switch (action.type) {
    case UPSERT_JOBS: {
      const { data } = action
      return Object.assign(
        {},
        state,
        { currentPage: data.meta.currentPageJobs.data.map(j => j.id) },
        { count: data.meta.currentPageJobs.meta.count },
        { items: Object.assign({}, state.items, action.data.specs) }
      )
    }
    case UPSERT_RECENTLY_CREATED_JOBS: {
      return Object.assign(
        {},
        state,
        {
          recentlyCreated: action.data.meta['recentlyCreatedJobs'].data.map(
            j => j.id
          )
        },
        { items: Object.assign({}, state.items, action.data.specs) }
      )
    }
    case UPSERT_JOB: {
      return Object.assign({}, state, {
        items: Object.assign({}, state.items, action.data.specs)
      })
    }
    default:
      return state
  }
}
