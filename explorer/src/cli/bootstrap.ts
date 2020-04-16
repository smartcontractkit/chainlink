import { openDbConnection } from '../database'

export async function bootstrap(cb: any) {
  const db = await openDbConnection()
  try {
    await cb(db)
  } catch (err) {
    console.error(err)
  } finally {
    db.close()
  }
}

