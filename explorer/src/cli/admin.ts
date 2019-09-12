import { Connection } from 'typeorm'
import { Admin } from '../entity/Admin'
import { createAdmin } from '../support/admin'
import { bootstrap } from './bootstrap'

const seed = async (username: string, password: string) => {
  return bootstrap(async (db: Connection) => {
    const admin: Admin = await createAdmin(db, username, password)

    console.log('created new chainlink admin')
    console.log('username: ', admin.username)
    console.log('password: ', password)
  })
}

export { seed }
