import { Connection } from 'typeorm'
import { openDbConnection } from '../database'

export async function bootstrap(cb: (db: Connection) => void) {
  const db = await openDbConnection()
  try {
    return cb(db)
  } catch (err) {
    console.error(err)
  } finally {
    db.close()
  }
}
