import {
  Column,
  Connection,
  Entity,
  OneToMany,
  PrimaryGeneratedColumn
} from 'typeorm'
import { JobRun } from './JobRun'
import { sha256 } from 'js-sha256'
import { randomBytes } from 'crypto'

export interface IChainlinkNodePresenter {
  id: number
  name: string
}

@Entity()
export class ChainlinkNode {
  public static build(name: string, secret: string): ChainlinkNode {
    const cl = new ChainlinkNode()
    cl.name = name
    cl.accessKey = generateRandomString(16)
    cl.salt = generateRandomString(32)
    cl.hashedSecret = hashCredentials(cl.accessKey, secret, cl.salt)
    return cl
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

  @OneToMany(() => JobRun, jobRun => jobRun.chainlinkNode, {
    onDelete: 'CASCADE'
  })
  jobRuns: Array<JobRun>

  public present(): IChainlinkNodePresenter {
    return {
      id: this.id,
      name: this.name
    }
  }
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
  const chainlinkNode = ChainlinkNode.build(name, secret)
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
