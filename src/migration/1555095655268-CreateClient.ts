import { MigrationInterface, QueryRunner } from 'typeorm'

export class AddClient1555095655268 implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`CREATE TABLE "client" (
      "id" BIGSERIAL PRIMARY KEY,
      "createdAt" TIMESTAMP NOT NULL DEFAULT now(),
      "name" CHARACTER VARYING UNIQUE,
      "accessKey" VARCHAR(32) NOT NULL,
      "hashedSecret" VARCHAR(64) NOT NULL,
      "salt" VARCHAR(64) NOT NULL
    )`)
  }

  public async down(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`DROP TABLE "client"`)
  }
}
