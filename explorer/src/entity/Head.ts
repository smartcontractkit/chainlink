import {
  Column,
  Entity,
  UpdateResult,
} from 'typeorm'

@Entity()
export class Head {
  @Column()
  public address: Buffer

  @Column()
  public topics: Buffer

  @Column()
  public data: Buffer

  @Column()
  public blockNumber: Buffer

  @Column()
  public txHash: Buffer

  @Column()
  public txIndex: Buffer

  @Column()
  public blockHash: Buffer

  @Column()
  public index: Buffer
  
  @Column()
  public removed: boolean
}
