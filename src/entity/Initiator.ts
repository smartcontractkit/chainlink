import { Column, Entity, OneToOne, PrimaryGeneratedColumn } from 'typeorm'
import { JobRun } from './JobRun'

@Entity()
export class Initiator {
  @PrimaryGeneratedColumn()
  id: number

  @OneToOne(type => JobRun, jobRun => jobRun.initiator)
  jobRun: JobRun

  @Column()
  type: string

  @Column()
  requestId: string

  @Column()
  txHash: string

  @Column()
  requester: string

  @Column()
  createdAt: Date
}
