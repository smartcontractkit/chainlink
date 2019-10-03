import { Column, Entity, PrimaryGeneratedColumn } from 'typeorm'

@Entity()
export class Admin {
  @PrimaryGeneratedColumn('uuid')
  id: string

  @Column()
  username: string

  @Column()
  hashedPassword: string

  @Column()
  createdAt: Date

  @Column()
  updatedAt: Date
}
