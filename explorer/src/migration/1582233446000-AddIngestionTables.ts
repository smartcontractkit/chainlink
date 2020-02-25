import { MigrationInterface, QueryRunner } from 'typeorm'

export class AddIngestionTables1582233446000 implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`
CREATE TABLE head (
  "address" bytea PRIMARY KEY,
  "topics" text NOT NULL,
  "data" bytea NOT NULL,
  "blockNumber" unsigned long NOT NULL,
  "txHash" bytea NOT NULL,
  "txIndex" bytea NOT NULL,
  "blockHash" bytea NOT NULL,
  "index" long NOT NULL,
  "removed" bool NOT NULL DEFAULT FALSE,
);
    `)
  }

  public async down(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query('DROP TABLE head;')
  }
}
