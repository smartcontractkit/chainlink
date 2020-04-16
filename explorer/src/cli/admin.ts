import { createAdmin } from '../support/admin'
import { bootstrap } from './bootstrap'

export const seed = async (username: string, password: string) => {
  return bootstrap(async () => {
    const admin = await createAdmin(username, password)

    console.log('created new chainlink admin')
    console.log('username: ', admin.username)
    console.log('password: ', password)
  })
}
