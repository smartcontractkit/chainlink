import {
  Column,
  Connection,
  CreateDateColumn,
  Entity,
  ManyToOne,
  OneToMany,
  PrimaryGeneratedColumn,
  UpdateDateColumn
} from 'typeorm'
import { ChainlinkNode } from './ChainlinkNode'
import { TaskRun } from './TaskRun'

@Entity()
export class Session {
  @Column()
  public chainlinkNodeId: number

  @Column({ nullable: true })
  public finishedAt: Date

  @PrimaryGeneratedColumn('uuid')
  public id: string

  @CreateDateColumn()
  private createdAt: Date

  @UpdateDateColumn()
  private updatedAt: Date
}
