import { MigrationInterface, QueryRunner } from 'typeorm'

export class CreateJobRun1553105100407 implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`CREATE TABLE "job_run" (
      "id" character varying NOT NULL,
      "jobId" character varying NOT NULL,
      "createdAt" TIMESTAMP NOT NULL DEFAULT now(),
      CONSTRAINT "PK_96fe0b041b8bc157dcec25bd8ef" PRIMARY KEY ("id")
    )`)
    await queryRunner.query(`CREATE INDEX job_id_idx ON "job_run" ("jobId")`)
  }

  public async down(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`DROP TABLE "job_run"`)
  }
}
