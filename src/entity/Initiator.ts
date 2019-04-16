import {
  Column,
  Entity,
  OneToOne,
  JoinColumn,
  PrimaryGeneratedColumn
} from 'typeorm'
import { JobRun } from './JobRun'

@Entity()
export class Initiator {
  @PrimaryGeneratedColumn()
  id: number

  @OneToOne(type => JobRun)
  @JoinColumn()
  jobRun: JobRun

  @Column()
  type: string

  @Column()
  requestId: string

  @Column()
  txHash: string

  @Column()
  requester: string
}
