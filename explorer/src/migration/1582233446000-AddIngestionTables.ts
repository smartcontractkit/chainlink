import { MigrationInterface, QueryRunner } from 'typeorm'

export class AddIngestionTables1582233446000 implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(`
CREATE TABLE ethereum_head (
  "id" BIGSERIAL PRIMARY KEY,
  "createdAt" timestamp without time zone DEFAULT now() NOT NULL,
  "blockHash" bytea NOT NULL,
  "parentHash" bytea NOT NULL,
  "uncleHash" bytea NOT NULL,
  "coinbase" bytea NOT NULL,
  "root" bytea NOT NULL,
  "txHash" bytea NOT NULL,
  "receiptHash" bytea NOT NULL,
  "bloom" bytea NOT NULL,
  "difficulty" numeric NOT NULL,
  "number" numeric NOT NULL,
  "gasLimit" bigint NOT NULL,
  "gasUsed" bigint NOT NULL,
  "time" bigint NOT NULL,
  "extra" bytea NOT NULL,
  "mixDigest" bytea NOT NULL,
  "nonce" bytea NOT NULL
);
    `)

    await queryRunner.query(`
CREATE TABLE ethereum_log (
  "id" BIGSERIAL PRIMARY KEY,
  "address" bytea NOT NULL,
  "createdAt" timestamp without time zone DEFAULT now() NOT NULL,
  "topics" bytea NOT NULL,
  "data" bytea NOT NULL,
  "blockNumber" bigint NOT NULL,
  "txHash" bytea NOT NULL,
  "txIndex" bytea NOT NULL,
  "blockHash" bytea NOT NULL,
  "index" bigint NOT NULL,
  "removed" bool NOT NULL DEFAULT FALSE
);
    `)
  }

  public async down(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(`
DROP TABLE "ethereum_log";
DROP TABLE "ethereum_head";
    `)
  }
}
