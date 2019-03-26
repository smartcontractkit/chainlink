import { Entity, PrimaryGeneratedColumn, Column, CreateDateColumn } from "typeorm"

@Entity()
export class JobRun {

    @PrimaryGeneratedColumn()
    id: number;

    @Column()
    jobRunId: string;

    @Column()
    jobId: string;

    @CreateDateColumn()
    createdAt: Date;

}
