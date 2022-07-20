import { InputErrors } from 'src/types/generated/graphql'

export const parseInputErrors = (payload: InputErrors) => {
  return payload.errors.reduce((obj, item) => {
    const key = item['path'].replace(/^input\//, '')

    return {
      ...obj,
      [key]: item.message,
    }
  }, {})
}
