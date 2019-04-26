import {
  Column,
  Entity,
  PrimaryGeneratedColumn
} from 'typeorm'

@Entity({ name: 'chainlink_node' })
export class PublicChainlinkNode {
  @PrimaryGeneratedColumn()
  id: number

  @Column()
  name: string
}
