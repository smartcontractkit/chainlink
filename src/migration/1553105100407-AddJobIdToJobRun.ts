import { MigrationInterface, QueryRunner } from "typeorm"

export class AddJobIdToJobRun1553105100407 implements MigrationInterface {

    public async up(queryRunner: QueryRunner): Promise<any> {
      await queryRunner.query(`TRUNCATE TABLE "job_run"`)
      await queryRunner.query(`ALTER TABLE "job_run" ADD COLUMN "jobId" UUID NOT NULL`)
      await queryRunner.query(`CREATE INDEX job_id_idx ON "job_run" ("jobId")`)
    }

    public async down(queryRunner: QueryRunner): Promise<any> {
      await queryRunner.query(`ALTER TABLE "job_run" DROP COLUMN "jobId"`)
    }

}
