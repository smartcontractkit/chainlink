import {
  Column,
  Entity,
  getRepository,
  getManager,
  PrimaryGeneratedColumn,
  UpdateDateColumn,
  UpdateResult,
} from 'typeorm'
import { ChainlinkNode } from './ChainlinkNode'

@Entity()
export class Session {
  @Column()
  public chainlinkNodeId: number

  @Column({ nullable: true })
  public finishedAt: Date

  @PrimaryGeneratedColumn('uuid')
  public id: string

  // @ts-ignore
  private createdAt: Date

  @UpdateDateColumn()
  // @ts-ignore
  private updatedAt: Date
}

export async function createSession(
  node: ChainlinkNode,
  manager = getManager(),
): Promise<Session> {
  // Close any other open sessions for this node
  await manager
    .getRepository(Session)
    .update(
      { chainlinkNodeId: node.id, finishedAt: null },
      { finishedAt: () => 'now()' },
    )

  const session = new Session()
  session.chainlinkNodeId = node.id
  return manager.save(session)
}

export async function retireSessions(): Promise<UpdateResult> {
  return getRepository(Session).update(
    { finishedAt: null },
    { finishedAt: () => 'now()' },
  )
}

export async function closeSession(session: Session): Promise<UpdateResult> {
  return getRepository(Session).update(session.id, {
    finishedAt: () => 'now()',
  })
}
