import { MigrationInterface, QueryRunner } from 'typeorm'

export class AddSession1566510989000 implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`
CREATE TABLE session (
  "id" character varying DEFAULT uuid_generate_v4() PRIMARY KEY,
  "chainlinkNodeId" bigint REFERENCES chainlink_node (id) NOT NULL,
  "createdAt" timestamp without time zone DEFAULT now() NOT NULL,
  "updatedAt" timestamp without time zone DEFAULT now() NOT NULL,
  "finishedAt" timestamp without time zone
);
    `)
  }

  public async down(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query('DROP TABLE sessions;')
  }
}
