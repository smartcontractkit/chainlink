import { MinLength } from 'class-validator'
import { randomBytes } from 'crypto'
import { sha256 } from 'js-sha256'
import {
  Column,
  createQueryBuilder,
  Entity,
  getRepository,
  OneToMany,
  PrimaryGeneratedColumn,
  Unique,
} from 'typeorm'
import { JobRun } from './JobRun'
import { Session } from './Session'

export interface ChainlinkNodePresenter {
  id: number
  name: string
}

@Entity()
@Unique(['name'])
export class ChainlinkNode {
  public static build({
    name,
    url,
    secret,
  }: {
    name: string
    url?: string
    secret: string
  }): ChainlinkNode {
    const cl = new ChainlinkNode()
    cl.name = name
    cl.url = url
    cl.accessKey = generateRandomString(16)
    cl.salt = generateRandomString(32)
    cl.hashedSecret = hashCredentials(cl.accessKey, secret, cl.salt)
    return cl
  }

  @PrimaryGeneratedColumn()
  id: number

  @MinLength(3, { message: 'must be at least 3 characters' })
  @Column()
  name: string

  @Column({ nullable: true })
  url: string

  @Column()
  accessKey: string

  @Column()
  hashedSecret: string

  @Column()
  salt: string

  @Column()
  createdAt: Date

  @OneToMany(
    () => JobRun,
    jobRun => jobRun.chainlinkNode,
    {
      onDelete: 'CASCADE',
    },
  )
  jobRuns: JobRun[]

  public present(): ChainlinkNodePresenter {
    return {
      id: this.id,
      name: this.name,
    }
  }
}

function generateRandomString(size: number): string {
  return randomBytes(size)
    .toString('base64')
    .replace(/[/+=]/g, '')
    .substring(0, size)
}

export const buildChainlinkNode = (
  name: string,
  url?: string,
): [ChainlinkNode, string] => {
  const secret = generateRandomString(64)
  const node = ChainlinkNode.build({ name, url, secret })

  return [node, secret]
}

export const createChainlinkNode = async (
  name: string,
  url?: string,
): Promise<[ChainlinkNode, string]> => {
  const secret = generateRandomString(64)
  const chainlinkNode = ChainlinkNode.build({ name, url, secret })
  const repo = getRepository(ChainlinkNode)
  return [await repo.save(chainlinkNode), secret]
}

export const deleteChainlinkNode = async (name: string) => {
  return createQueryBuilder()
    .delete()
    .from(ChainlinkNode)
    .where('name = :name', {
      name,
    })
    .execute()
}

export function hashCredentials(
  accessKey: string,
  secret: string,
  salt: string,
): string {
  return sha256(`v0-${accessKey}-${secret}-${salt}`)
}

export async function find(id: number): Promise<ChainlinkNode> {
  return getRepository(ChainlinkNode).findOne({ id })
}

export interface JobCountReport {
  completed: number
  errored: number
  in_progress: number
  total: number
}

export async function jobCountReport(
  node: ChainlinkNode | number,
): Promise<JobCountReport> {
  const id = node instanceof ChainlinkNode ? node.id : node

  const initialReport: JobCountReport = {
    completed: 0,
    errored: 0,
    in_progress: 0, // eslint-disable-line
    total: 0,
  }

  const jobCountQueryResult = await getRepository(JobRun)
    .createQueryBuilder()
    .select('COUNT(*), status')
    .where({ chainlinkNodeId: id })
    .groupBy('status')
    .getRawMany()

  const report = jobCountQueryResult.reduce((result, { count, status }) => {
    result[status] = parseInt(count, 10)
    result.total = result.total + result[status]
    return result
  }, initialReport)

  return report
}

// calculating uptime by diffing createdAt and finishedAt columns
// using strategy described here: http://www.sqlines.com/postgresql/how-to/datediff
// typeORM missing UNION function, so must do in two separate queries or write entire
// query in raw SQL
export async function uptime(node: ChainlinkNode | number) {
  const id = node instanceof ChainlinkNode ? node.id : node
  return (await historicUptime(id)) + (await currentUptime(id))
}

// uptime from completed sessions
async function historicUptime(id: number): Promise<number> {
  const queryResult = await createQueryBuilder()
    .select(
      `EXTRACT(EPOCH FROM session."finishedAt" - session."createdAt") as seconds`,
    )
    .from(Session, 'session')
    .where({ chainlinkNodeId: id })
    .andWhere('session.finishedAt is not null')
    .getRawOne()
  // NOTE: If there are no sessions, SELECT EXTRACT... returns null
  const seconds = queryResult?.seconds ?? 0
  return Math.max(0, seconds)
}

// uptime from current open session
async function currentUptime(id: number): Promise<number> {
  const queryResult = await createQueryBuilder()
    .select(
      `FLOOR(EXTRACT(EPOCH FROM (now() - session."createdAt"))) as seconds`,
    )
    .from(Session, 'session')
    .where({ chainlinkNodeId: id })
    .andWhere('session."finishedAt" is null')
    .getRawOne()
  // NOTE: If there are no sessions, SELECT EXTRACT... returns null
  const seconds = queryResult?.seconds ?? 0
  return Math.max(0, seconds)
}
