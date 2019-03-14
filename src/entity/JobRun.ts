import { Entity, PrimaryGeneratedColumn, Column, CreateDateColumn } from "typeorm"

@Entity()
export class JobRun {

    @PrimaryGeneratedColumn()
    id: number;

    @Column()
    requestId: string;

    @CreateDateColumn()
    createdAt: Date;

}
