export interface State {
  query?: string
}

export type Query = string | undefined

export type Action = { type: string; location?: Location }

const initialState = { query: undefined }

const parseQuery = (location: Location = document.location): Query => {
  const searchParams = new URL(location.toString()).searchParams
  const search = searchParams.get('search')

  if (search) {
    return search
  }
}

export default (state: State = initialState, action: Action) => {
  return Object.assign({}, state, { query: parseQuery(action.location) })
}
