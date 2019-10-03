import { MigrationInterface, QueryRunner, Table, TableColumn } from 'typeorm'

const TABLE_NAME = 'admin'

const idColumm = new TableColumn({
  name: 'id',
  type: 'character varying',
  isPrimary: true,
  isNullable: false,
  default: 'uuid_generate_v4()',
})

const usernameColumn = new TableColumn({
  name: 'username',
  type: 'character varying',
  isUnique: true,
  isNullable: false,
})

const hashedPasswordColumn = new TableColumn({
  name: 'hashedPassword',
  type: 'character varying',
  isNullable: false,
})

const createdAtColumn = new TableColumn({
  name: 'createdAt',
  type: 'timestamp without time zone',
  isNullable: true,
  default: 'now()',
})

const updatedAtColumn = new TableColumn({
  name: 'updatedAt',
  type: 'timestamp without time zone',
  isNullable: true,
  default: 'now()',
})

export class AddAdmin1565028153000 implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    const options = {
      name: TABLE_NAME,
      columns: [
        idColumm,
        usernameColumn,
        hashedPasswordColumn,
        createdAtColumn,
        updatedAtColumn,
      ],
    }

    await queryRunner.createTable(new Table(options))
  }

  public async down(queryRunner: QueryRunner): Promise<any> {
    const options = { name: TABLE_NAME }
    await queryRunner.dropTable(new Table(options))
  }
}
