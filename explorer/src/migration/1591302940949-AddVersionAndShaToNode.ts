import { MigrationInterface, QueryRunner } from 'typeorm'

export class AddVersionAndShaToNode1591302940949 implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`
        ALTER TABLE "chainlink_node" ADD COLUMN "coreVersion" varchar NULL;
        ALTER TABLE "chainlink_node" ADD COLUMN "coreSHA" varchar NULL;
      `)
  }

  public async down(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`
        ALTER TABLE "chainlink_node" DROP COLUMN "coreVersion";
        ALTER TABLE "chainlink_node" DROP COLUMN "coreSHA";
      `)
  }
}
