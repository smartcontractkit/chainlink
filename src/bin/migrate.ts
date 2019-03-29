import { createDbConnection, closeDbConnection, getDb } from '../database'

const runMigrations = async () => {
  await createDbConnection()
  const db = getDb()
  try {
    await db.runMigrations()
  } catch (err) {
    console.error(err)
  } finally {
    await closeDbConnection()
  }
}

runMigrations()
