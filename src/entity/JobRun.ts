import {
  Entity,
  PrimaryColumn,
  Column,
  CreateDateColumn,
  OneToMany
} from 'typeorm'
import { TaskRun } from './TaskRun'

@Entity()
export class JobRun {
  @PrimaryColumn()
  id: string

  @Column()
  jobId: string

  @Column()
  status: string

  @Column({ nullable: true })
  error: string

  @Column()
  initiatorType: string

  @CreateDateColumn()
  createdAt: Date

  @Column({ nullable: true })
  completedAt: Date

  @OneToMany(type => TaskRun, taskRun => taskRun.jobRun, {
    onDelete: 'CASCADE'
  })
  taskRuns: TaskRun[]
}

export const fromString = (str: any): JobRun => {
  const json = JSON.parse(str)
  const jr = new JobRun()
  jr.id = json.id
  jr.jobId = json.jobId
  jr.status = json.status
  jr.initiatorType = json.initiator.type
  jr.createdAt = new Date(json.createdAt)
  jr.completedAt = json.completedAt && new Date(json.completedAt)
  return jr
}
