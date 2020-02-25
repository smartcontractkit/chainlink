import { MigrationInterface, QueryRunner } from 'typeorm'

export class AddIngestionTables1582233446000 implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`
CREATE TABLE ethereum_head (
  "id" BIGSERIAL PRIMARY KEY,
  "parent_hash" bytea NOT NULL,
  "uncle_hash" bytea NOT NULL,
  "coinbase" bytea NOT NULL,
  "root" bytea NOT NULL,
  "tx_hash" bytea NOT NULL,
  "receipt_hash" bytea NOT NULL,
  "bloom" bytea NOT NULL,
  "difficulty" numeric NOT NULL,
  "number" numeric NOT NULL,
  "gas_limit" bigint NOT NULL,
  "gas_used" bigint NOT NULL,
  "time" bigint NOT NULL,
  "extra" bytea NOT NULL,
  "mix_digest" bytea NOT NULL,
  "nonce" bytea NOT NULL
);
    `)

    await queryRunner.query(`
CREATE TABLE ethereum_log (
  "address" bytea PRIMARY KEY,
  "topics" text NOT NULL,
  "data" bytea NOT NULL,
  "block_number" bigint NOT NULL,
  "tx_hash" bytea NOT NULL,
  "tx_index" bytea NOT NULL,
  "block_hash" bytea NOT NULL,
  "index" bigint NOT NULL,
  "removed" bool NOT NULL DEFAULT FALSE
);
    `)
  }

  public async down(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`
DROP TABLE "ethereum_log";
DROP TABLE "ethereum_head";
    `)
  }
}
