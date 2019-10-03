import { createAdmin } from '../support/admin'
import { bootstrap } from './bootstrap'

export const seed = async (username: string, password: string) => {
  return bootstrap(async db => {
    const admin = await createAdmin(db, username, password)

    console.log('created new chainlink admin')
    console.log('username: ', admin.username)
    console.log('password: ', password)
  })
}
