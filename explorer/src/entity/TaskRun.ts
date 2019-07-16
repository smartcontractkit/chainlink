import { Column, Entity, ManyToOne, PrimaryGeneratedColumn } from 'typeorm'
import { JobRun } from './JobRun'

type TransactionStatus = 'fulfilledRunLog' | 'noFulfilledRunLog'

@Entity()
export class TaskRun {
  @PrimaryGeneratedColumn()
  id: number

  @ManyToOne(() => JobRun, jobRun => jobRun.taskRuns)
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

  @Column({ nullable: true, name: 'minimumConfirmations_new1562419039813' })
  minimumConfirmations?: string

  @Column({ nullable: true, name: 'confirmations_new1562419039813' })
  confirmations?: string

  //
  // TODO: Old columns will be deleted once node operators have migrated to this or
  // later version.
  @Column({ nullable: true, name: 'minimumConfirmations' })
  minimumConfirmationsOld?: number

  @Column({ nullable: true, name: 'confirmations' })
  confirmationsOld?: number
}
