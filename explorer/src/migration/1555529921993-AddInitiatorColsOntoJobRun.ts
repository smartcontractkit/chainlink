import { MigrationInterface, QueryRunner } from 'typeorm'

export class AddInitiatorColsOntoJobRun1555529921993
  implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    // backfill job_runs without a type to web, it is corrected
    // in a later migration, where we remove the default.
    await queryRunner.query(`ALTER TABLE "job_run"
      ADD COLUMN type varchar NOT NULL DEFAULT 'web',
      ADD COLUMN "requestId" varchar,
      ADD COLUMN "txHash" varchar,
      ADD COLUMN requester varchar;
    `)
  }

  public async down(queryRunner: QueryRunner): Promise<any> {}
}
