import {
  Column,
  Connection,
  Entity,
  ManyToOne,
  OneToMany,
  PrimaryGeneratedColumn
} from 'typeorm'
import { JobRun } from './JobRun'
import { sha256 } from 'js-sha256'
import { randomBytes } from 'crypto'

@Entity()
export class Client {
  constructor(name: string, secret: string) {
    this.name = name
    this.accessKey = generateRandomString(16)
    this.salt = generateRandomString(32)

    const hashInput = `v0-${this.accessKey}-${secret}-${this.salt}`
    this.hashedSecret = sha256(hashInput)
  }

  @PrimaryGeneratedColumn()
  id: number

  @Column()
  name: string

  @Column()
  accessKey: string

  @Column()
  hashedSecret: string

  @Column()
  salt: string

  @OneToMany(type => JobRun, jobRun => jobRun.client, {
    onDelete: 'CASCADE'
  })
  jobRuns: Array<JobRun>
}

const generateRandomString = (size: number): string => {
  return randomBytes(size)
    .toString('base64')
    .replace(/[/+=]/g, '')
}

export const createClient = async (
  db: Connection,
  name: string
): Promise<[Client, string]> => {
  const secret = generateRandomString(16)
  const client = new Client(name, secret)
  return [await db.manager.save(client), secret]
}

export const deleteClient = (db: Connection, name: string) => {
  db.manager
    .createQueryBuilder()
    .delete()
    .from(Client)
    .where('name = :name', {
      name: name
    })
    .execute()
}
