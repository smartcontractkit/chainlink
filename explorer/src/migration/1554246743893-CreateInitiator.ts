import { MigrationInterface, QueryRunner } from 'typeorm'

export class AddInitiators1554246743893 implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`CREATE TABLE "initiator" (
      "id" BIGSERIAL PRIMARY KEY,
      "jobRunId" BIGINT REFERENCES "job_run" ("id") ON DELETE CASCADE,
      "requestId" character varying NOT NULL,
      "txHash" character varying NOT NULL,
      "requester" character varying NOT NULL,
      "createdAt" TIMESTAMP NOT NULL
    )`)
    await queryRunner.query(
      `CREATE INDEX initiator_id_idx ON "initiator" ("jobRunId")`
    )
  }

  public async down(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`DROP TABLE "initiator"`)
  }
}
