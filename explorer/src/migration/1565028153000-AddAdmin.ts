import { MigrationInterface, QueryRunner } from 'typeorm'

export class AddAdmin1565028153000 implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query(`
CREATE TABLE admin (
  "id" character varying DEFAULT uuid_generate_v4() PRIMARY KEY,
  "username" character varying UNIQUE NOT NULL,
  "hashedPassword" character varying NOT NULL,
  "createdAt" timestamp without time zone DEFAULT now() NOT NULL,
  "updatedAt" timestamp without time zone DEFAULT now() NOT NULL
);
    `)
  }

  public async down(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.query('DROP TABLE admin;')
  }
}
