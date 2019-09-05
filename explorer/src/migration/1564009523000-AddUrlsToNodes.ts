import { MigrationInterface, QueryRunner } from 'typeorm'

export class AddUrlsToNodes1564009523000 implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`
      ALTER TABLE "chainlink_node" ADD COLUMN "url" character varying DEFAULT NULL;
    `)
  }

  public async down(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`
      ALTER TABLE "chainlink_node" DROP COLUMN "url";
  `)
  }
}
