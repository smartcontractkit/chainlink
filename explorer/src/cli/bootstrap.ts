import { Connection } from 'typeorm'
import { closeDbConnection, getDb } from '../database'

export async function bootstrap(callback: (db: Connection) => void) {
  const db = await getDb()
  try {
    await callback(db)
  } catch (e) {
    console.error(e)
  }
  try {
    await closeDbConnection()
  } catch (e) {
    console.error(e)
  }
}
