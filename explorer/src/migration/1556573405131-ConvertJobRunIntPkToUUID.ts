import { MigrationInterface, QueryRunner } from 'typeorm'

export class ConvertJobRunIntPkToUUID1556573405131
  implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query('CREATE EXTENSION IF NOT EXISTS "uuid-ossp";')
    await queryRunner.query(
      'ALTER TABLE "task_run" DROP CONSTRAINT "task_run_jobRunId_fkey";'
    )
    await queryRunner.query(`
      ALTER TABLE job_run ALTER COLUMN id TYPE CHARACTER VARYING;
      ALTER TABLE task_run ALTER COLUMN "jobRunId" TYPE CHARACTER VARYING;
      ALTER TABLE job_run ALTER COLUMN id SET DEFAULT uuid_generate_v4();
    `)
    await queryRunner.query(`
      ALTER TABLE "task_run"
      ADD CONSTRAINT "task_run_jobRunId_fkey"
      FOREIGN KEY ("jobRunId")
      REFERENCES job_run(id);
    `)
  }

  public async down(queryRunner: QueryRunner): Promise<any> {}
}
