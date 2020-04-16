import { getRepository } from 'typeorm'
import { Admin } from '../entity/Admin'
import { hash } from '../services/password'

export async function createAdmin(username: string, password: string) {
  const hashedPassword = await hash(password)
  const admin = new Admin()
  admin.username = username
  admin.hashedPassword = hashedPassword

  return getRepository(Admin).save(admin)
}
