import { Entity, PrimaryColumn, Column, ManyToOne } from 'typeorm'
import { JobRun } from './JobRun'

@Entity()
export class TaskRun {
  @PrimaryColumn()
  id: string

  @Column()
  type: string

  @ManyToOne(type => JobRun, jobRun => jobRun.taskRuns)
  jobRun: JobRun
}
