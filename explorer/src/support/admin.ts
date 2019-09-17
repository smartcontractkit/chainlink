import { Connection } from 'typeorm'
import bcrypt from 'bcrypt'
import { Admin } from '../entity/Admin'

export async function createAdmin(
  db: Connection,
  username: string,
  password: string,
) {
  const hashedPassword = await bcrypt.hash(password, 1)
  const admin = new Admin()
  admin.username = username
  admin.hashedPassword = hashedPassword

  return db.manager.save(admin)
}
