import { InputErrors } from 'src/types/generated/graphql'

// reduceInputErrors reduces input errors into a simple object
// function reduceInputErrors(errors: InputErrors) {
//   return errors.reduce((obj, item) => {
//     const key = item['path'].replace(/^input\//, '')

//     return {
//       ...obj,
//       [key]: item.message,
//     }
//   }, {})
// }
