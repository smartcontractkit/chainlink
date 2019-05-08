import { Column, Entity, ManyToOne, PrimaryGeneratedColumn } from 'typeorm'
import { JobRun } from './JobRun'

type TransactionStatus = 'fulfilledRunLog' | 'noFulfilledRunLog'

@Entity()
export class TaskRun {
  @PrimaryGeneratedColumn()
  id: number

  @ManyToOne(type => JobRun, jobRun => jobRun.taskRuns)
  jobRun: JobRun

  @Column()
  index: number

  @Column()
  type: string

  @Column()
  status: string

  @Column({ nullable: true })
  error?: string

  @Column({ nullable: true })
  transactionHash?: string

  @Column({ nullable: true })
  transactionStatus?: TransactionStatus
}
