import { Connection } from 'typeorm'
import { openDbConnection } from '../database'

export async function bootstrap(cb: any) {
  openDbConnection().then(async (db: Connection) => {
    try {
      await cb(db)
    } catch (err) {
      console.error(err)
    }
  })
}

