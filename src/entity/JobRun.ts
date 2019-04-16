import {
  Column,
  Connection,
  Entity,
  OneToOne,
  OneToMany,
  PrimaryGeneratedColumn
} from 'typeorm'
import { TaskRun } from './TaskRun'
import { Initiator } from './Initiator'

@Entity()
export class JobRun {
  @PrimaryGeneratedColumn()
  id: number

  @Column()
  runId: string

  @Column()
  jobId: string

  @Column()
  status: string

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

  @OneToOne(type => Initiator, initiator => initiator.jobRun, {
    onDelete: 'CASCADE'
  })
  initiator: Initiator
}

export const fromString = (str: string): JobRun => {
  const json = JSON.parse(str)
  const jr = new JobRun()
  jr.runId = json.runId
  jr.jobId = json.jobId
  jr.status = json.status
  jr.createdAt = new Date(json.createdAt)
  jr.completedAt = json.completedAt && new Date(json.completedAt)
  jr.initiator = new Initiator()
  jr.initiator.type = json.initiator.type
  jr.initiator.requestId = json.initiator.requestId
  jr.initiator.txHash = json.initiator.txHash
  jr.initiator.requester = json.initiator.requester
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
