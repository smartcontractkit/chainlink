import {
  Column,
  Connection,
  CreateDateColumn,
  Entity,
  In,
  OneToMany,
  PrimaryGeneratedColumn
} from 'typeorm'
import { TaskRun } from './TaskRun'

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

  @CreateDateColumn()
  createdAt: Date

  @Column({ nullable: true })
  completedAt: Date

  @OneToMany(type => TaskRun, taskRun => taskRun.jobRun, {
    onDelete: 'CASCADE'
  })
  taskRuns: Array<TaskRun>
}

export const fromString = (str: string): JobRun => {
  const json = JSON.parse(str)
  const jr = new JobRun()
  jr.runId = json.runId
  jr.jobId = json.jobId
  jr.status = json.status
  jr.createdAt = new Date(json.createdAt)
  jr.completedAt = json.completedAt && new Date(json.completedAt)
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

export const search = async (
  db: Connection,
  searchQuery?: string
): Promise<Array<JobRun>> => {
  let params = {}

  if (searchQuery != null) {
    const searchTokens = searchQuery.split(/\s+/)
    params = {
      where: [{ runId: In(searchTokens) }, { jobId: In(searchTokens) }]
    }
  }

  return db.getRepository(JobRun).find(params)
}
