import {
  Column,
  Connection,
  Entity,
  OneToOne,
  OneToMany,
  PrimaryGeneratedColumn
} from 'typeorm'
import { TaskRun } from './TaskRun'
import { Client } from './Client'

@Entity()
export class JobRun {
  @PrimaryGeneratedColumn()
  id: number

  @Column()
  clientId: number

  @Column()
  runId: string

  @Column()
  jobId: string

  @Column()
  status: string

  @Column()
  type: string

  @Column({ nullable: true })
  requestId: string

  @Column({ nullable: true })
  txHash: string

  @Column({ nullable: true })
  requester: string

  @Column({ nullable: true })
  error: string

  @Column()
  createdAt: Date

  @Column({ nullable: true })
  completedAt: Date

  @OneToMany(type => TaskRun, taskRun => taskRun.jobRun, {
    onDelete: 'CASCADE'
  })
  taskRuns: Array<TaskRun>

  @OneToOne(type => Client, client => client.jobRuns)
  client: Client
}

export const fromString = (str: string): JobRun => {
  const json = JSON.parse(str)
  const jr = new JobRun()
  jr.runId = json.runId
  jr.jobId = json.jobId
  jr.status = json.status
  jr.createdAt = new Date(json.createdAt)
  jr.completedAt = json.completedAt && new Date(json.completedAt)

  jr.type = json.initiator.type
  jr.requestId = json.initiator.requestId
  jr.txHash = json.initiator.txHash
  jr.requester = json.initiator.requester

  jr.taskRuns = json.tasks.map((trstr: any, index: number) => {
    const tr = new TaskRun()
    tr.index = index
    tr.type = trstr.type
    tr.status = trstr.status
    tr.error = trstr.error

    return tr
  })

  return jr
}

export interface ISearchParams {
  searchQuery?: string
  page?: number
  limit?: number
}

export const search = async (
  db: Connection,
  params: ISearchParams
): Promise<Array<JobRun>> => {
  let query = db.getRepository(JobRun).createQueryBuilder('job_run')

  if (params.searchQuery != null) {
    const searchTokens = params.searchQuery.split(/\s+/)
    query = query
      .where('job_run.runId IN(:...searchTokens)', { searchTokens })
      .orWhere('job_run.jobId IN(:...searchTokens)', { searchTokens })
  }

  if (params.limit != null) {
    query = query.limit(params.limit)
  }

  if (params.page !== undefined) {
    const offset = (params.page - 1) * params.limit
    query = query.offset(offset)
  }

  return query.orderBy('job_run.createdAt', 'DESC').getMany()
}

export const saveJobRunTree = async (db: Connection, jobRun: JobRun) => {
  await db.manager.transaction(async manager => {
    const builder = manager.createQueryBuilder()

    const response = await builder
      .insert()
      .into(JobRun)
      .values(jobRun)
      .onConflict(
        `("runId") DO UPDATE SET
        "status" = :status
        ,"error" = :error
        ,"completedAt" = :completedAt
      `
      )
      .setParameter('status', jobRun.status)
      .setParameter('error', jobRun.error)
      .setParameter('completedAt', jobRun.completedAt)
      .execute()

    await Promise.all(
      jobRun.taskRuns.map(tr => {
        // new builder since execute stmnt above seems to mutate.
        const builder = manager.createQueryBuilder()
        tr.jobRun = jobRun
        return builder
          .insert()
          .into(TaskRun)
          .values(tr)
          .onConflict(
            `("index", "jobRunId") DO UPDATE SET
              "status" = :status
              ,"error" = :error
              `
          )
          .setParameter('status', tr.status)
          .setParameter('error', tr.error)
          .execute()
      })
    )
  })
}
