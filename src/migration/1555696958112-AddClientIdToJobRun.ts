import { MigrationInterface, QueryRunner } from 'typeorm'

export class AddClientIdToJobRun1555696958112 implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`ALTER TABLE "job_run"
      ADD COLUMN "clientId" BIGINT REFERENCES client (id) NOT NULL
    `)
    await queryRunner.query(
      `CREATE INDEX job_id_client_id_index ON "job_run" ("id", "clientId")`
    )
  }

  public async down(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`DROP INDEX job_id_client_id_index`)
  }
}
