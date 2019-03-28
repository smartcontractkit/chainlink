import { Entity, PrimaryColumn, Column, ManyToOne } from 'typeorm'
import { JobRun } from './JobRun'

@Entity()
export class TaskRun {
  @PrimaryColumn()
  id: string

  @ManyToOne(type => JobRun, jobRun => jobRun.taskRuns)
  jobRun: JobRun

  @Column()
  index: number

  @Column()
  type: string

  @Column()
  status: string

  @Column({ nullable: true })
  error: string
}
