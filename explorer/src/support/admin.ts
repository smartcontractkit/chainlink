import { Connection } from 'typeorm'
import { Admin } from '../entity/Admin'
import { hash } from '../services/password'

export async function createAdmin(
  db: Connection,
  username: string,
  password: string,
) {
  const hashedPassword = await hash(password, 1)
  const admin = new Admin()
  admin.username = username
  admin.hashedPassword = hashedPassword

  return db.manager.save(admin)
}
