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
export class ChainlinkNode {
  constructor(name: string, secret: string) {
    this.name = name
    this.accessKey = generateRandomString(16)
    this.salt = generateRandomString(32)

    this.hashedSecret = hashCredentials(this.accessKey, secret, this.salt)
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

  @OneToMany(type => JobRun, jobRun => jobRun.chainlinkNode, {
    onDelete: 'CASCADE'
  })
  jobRuns: Array<JobRun>
}

const generateRandomString = (size: number): string => {
  return randomBytes(size)
    .toString('base64')
    .replace(/[/+=]/g, '')
    .substring(0, size)
}

export const createChainlinkNode = async (
  db: Connection,
  name: string
): Promise<[ChainlinkNode, string]> => {
  const secret = generateRandomString(64)
  const chainlinkNode = new ChainlinkNode(name, secret)
  return [await db.manager.save(chainlinkNode), secret]
}

export const deleteChainlinkNode = async (db: Connection, name: string) => {
  return db.manager
    .createQueryBuilder()
    .delete()
    .from(ChainlinkNode)
    .where('name = :name', {
      name: name
    })
    .execute()
}

export const hashCredentials = (
  accessKey: string,
  secret: string,
  salt: string
): string => {
  return sha256(`v0-${accessKey}-${secret}-${salt}`)
}
