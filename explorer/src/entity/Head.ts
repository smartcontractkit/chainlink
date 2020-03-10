import { Column, Entity, PrimaryColumn, PrimaryGeneratedColumn } from 'typeorm'
import { BigNumber } from 'bignumber.js'

const transformer = {
  from: (value: string): BigNumber => new BigNumber(value),
  to: (value: BigNumber): string => value.toString(),
}

@Entity({ name: 'ethereum_head' })
export class Head {
  @PrimaryGeneratedColumn()
  id: number

  @Column('bytea')
  public parentHash: Buffer

  @Column('bytea')
  public uncleHash: Buffer

  @Column('bytea')
  public coinbase: Buffer

  @Column('bytea')
  public root: Buffer

  @Column('bytea')
  public txHash: Buffer

  @Column('bytea')
  public receiptHash: Buffer

  @Column('bytea')
  public bloom: Buffer

  @Column('bigint', { transformer })
  public difficulty: BigNumber

  @Column('bigint', { transformer })
  public number: BigNumber

  @Column('bigint', { transformer })
  public gasLimit: BigNumber

  @Column('bigint', { transformer })
  public gasUsed: BigNumber

  @Column('bigint', { transformer })
  public time: BigNumber

  @Column('bytea')
  public extra: Buffer

  @Column('bytea')
  public mixDigest: Buffer

  @Column('bytea')
  public nonce: Buffer

  @Column()
  createdAt: Date
}

@Entity({ name: 'ethereum_log' })
export class Log {
  @PrimaryColumn('bytea')
  public address: Buffer

  @Column('bytea')
  public topics: Buffer

  @Column('bytea')
  public data: Buffer

  @Column('bytea')
  public blockNumber: Buffer

  @Column('bytea')
  public txHash: Buffer

  @Column('bytea')
  public txIndex: Buffer

  @Column('bytea')
  public blockHash: Buffer

  @Column('bytea')
  public index: Buffer

  @Column()
  public removed: boolean

  @Column()
  createdAt: Date
}
