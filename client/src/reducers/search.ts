export interface IState {
  query?: string
}

export type Query = string | undefined

export type SearchAction =
  | { type: 'UPDATE_SEARCH_QUERY'; query?: string }
  | { type: '@@redux/INIT' }
  | { type: '@@INIT' }

const initialState = { query: undefined }

const initQuery = (): Query => {
  const searchParams = new URL(document.location.toString()).searchParams
  const search = searchParams.get('search')

  if (search) {
    return search
  }
}

export default (state: IState = initialState, action: SearchAction) => {
  switch (action.type) {
    case '@@redux/INIT':
    case '@@INIT':
      return Object.assign({}, state, { query: initQuery() })
    case 'UPDATE_SEARCH_QUERY':
      return Object.assign({}, state, { query: action.query })
    default:
      return state
  }
}
