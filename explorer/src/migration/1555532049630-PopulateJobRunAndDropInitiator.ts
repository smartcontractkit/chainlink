import { MigrationInterface, QueryRunner } from 'typeorm'

export class PopulateJobRunAndDropInitiator1555532049630
  implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`INSERT INTO "job_run"
      (type, "requestId", "txHash", requester)
      SELECT initiator.type, "initiator"."requestId", "initiator"."txHash", initiator.requester
      FROM initiator
      JOIN job_run ON "job_run"."id" = initiator."jobRunId";
    `)
    await queryRunner.query(`
      ALTER TABLE "job_run"
      ALTER COLUMN type
      DROP DEFAULT;
    `)
    await queryRunner.query(`DROP TABLE "initiator"`)
  }

  public async down(queryRunner: QueryRunner): Promise<any> {}
}
