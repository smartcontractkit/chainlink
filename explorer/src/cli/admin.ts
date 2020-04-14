import { createAdmin } from '../support/admin'
import { bootstrap } from './bootstrap'
import { Connection } from 'typeorm'

export const seed = async (username: string, password: string) => {
  return bootstrap(async (db: Connection) => {
    const admin = await createAdmin(db, username, password)

    console.log('created new chainlink admin')
    console.log('username: ', admin.username)
    console.log('password: ', password)
  })
}
