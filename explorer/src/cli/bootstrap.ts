import { Connection } from 'typeorm'
import { closeDbConnection, getDb } from '../database'

export async function bootstrap(callback: (db: Connection) => Promise<void>) {
  const db = await getDb()

  await callback(db).catch(console.error)
  await closeDbConnection().catch(console.error)
}
