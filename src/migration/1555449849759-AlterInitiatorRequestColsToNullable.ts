import { MigrationInterface, QueryRunner, TableColumn } from 'typeorm'

export class AlterInitiatorRequestColsToNullable1555449849759
  implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.changeColumns('initiator', [
      {
        oldColumn: new TableColumn({
          name: 'requestId',
          type: 'varchar',
          isNullable: false
        }),
        newColumn: new TableColumn({
          name: 'requestId',
          type: 'varchar',
          isNullable: true
        })
      },
      {
        oldColumn: new TableColumn({
          name: 'txHash',
          type: 'varchar',
          isNullable: false
        }),
        newColumn: new TableColumn({
          name: 'txHash',
          type: 'varchar',
          isNullable: true
        })
      },
      {
        oldColumn: new TableColumn({
          name: 'requester',
          type: 'varchar',
          isNullable: false
        }),
        newColumn: new TableColumn({
          name: 'requester',
          type: 'varchar',
          isNullable: true
        })
      }
    ])
  }

  public async down(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.changeColumns('initiator', [
      {
        oldColumn: new TableColumn({
          name: 'requestId',
          type: 'varchar',
          isNullable: true
        }),
        newColumn: new TableColumn({
          name: 'requestId',
          type: 'varchar',
          isNullable: false
        })
      },
      {
        oldColumn: new TableColumn({
          name: 'txHash',
          type: 'varchar',
          isNullable: true
        }),
        newColumn: new TableColumn({
          name: 'txHash',
          type: 'varchar',
          isNullable: false
        })
      },
      {
        oldColumn: new TableColumn({
          name: 'requester',
          type: 'varchar',
          isNullable: true
        }),
        newColumn: new TableColumn({
          name: 'requester',
          type: 'varchar',
          isNullable: false
        })
      }
    ])
  }
}
