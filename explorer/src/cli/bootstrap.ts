import { Connection } from 'typeorm'
import { openDbConnection } from '../database'

interface Callback {
  (db: Connection): Promise<void>
}

export async function bootstrap(cb: Callback) {
  const db = await openDbConnection()
  try {
    return await cb(db)
  } catch (err) {
    console.error(err)
    throw err
  } finally {
    db.close()
  }
}
