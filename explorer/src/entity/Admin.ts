import { getRepository, Column, Entity, PrimaryGeneratedColumn } from 'typeorm'
import { compare as comparePassword } from '../services/password'

@Entity()
export class Admin {
  @PrimaryGeneratedColumn('uuid')
  id: string

  @Column()
  username: string

  @Column()
  hashedPassword: string

  @Column()
  createdAt: Date

  @Column()
  updatedAt: Date
}

export function find(username: string): Promise<Admin> {
  return getRepository(Admin).findOne({ username })
}

export async function isValidPassword(
  password: string,
  admin?: Admin,
): Promise<boolean> {
  if (!admin) {
    return false
  }

  return comparePassword(password, admin.hashedPassword)
}
