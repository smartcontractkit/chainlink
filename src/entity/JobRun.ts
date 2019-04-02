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
  taskRuns: TaskRun[]
}

export const fromString = (str: string): JobRun => {
  const json = JSON.parse(str)
  const jr = new JobRun()
  jr.runId = json.runId
  jr.jobId = json.jobId
  jr.status = json.status
  jr.createdAt = new Date(json.createdAt)
  jr.completedAt = json.completedAt && new Date(json.completedAt)
  jr.taskRuns = json.taskRuns.map((trstr: any, index: number) => {
    const tr = new TaskRun()
    tr.id = trstr.id
    tr.index = index
    tr.type = trstr.task.type
    tr.status = trstr.status
    tr.error = trstr.result.error

    return tr
  })

  return jr
}

export const search = async (
  db: Connection,
  searchTokens: Array<string>
): Promise<Array<JobRun>> => {
  return db.getRepository(JobRun).find({
    where: [{ runId: In(searchTokens) }, { jobId: In(searchTokens) }]
  })
}
