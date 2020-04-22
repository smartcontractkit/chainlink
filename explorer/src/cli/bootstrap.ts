import { Connection } from 'typeorm'
import { openDbConnection } from '../database'

export async function bootstrap(cb: (db: Connection) => Promise<void>) {
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
