import { closeDbConnection, getDb } from '../database'

export async function bootstrap(cb: any) {
  const db = await getDb()
  try {
    await cb(db)
  } catch (e) {
    console.error(e)
  }
  try {
    await closeDbConnection()
  } catch (e) {
    console.error(e)
  }
}
