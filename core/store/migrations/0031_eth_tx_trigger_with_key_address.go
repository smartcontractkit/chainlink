package migrations

import (
	"gorm.io/gorm"
)

const up31 = `
CREATE OR REPLACE FUNCTION public.notifyethtxinsertion() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
        BEGIN
		PERFORM pg_notify('insert_on_eth_txes'::text, encode(NEW.from_address, 'hex'));
		RETURN NULL;
        END
        $$;

DROP TRIGGER IF EXISTS notify_eth_tx_insertion on public.eth_txes;
CREATE TRIGGER notify_eth_tx_insertion AFTER INSERT ON public.eth_txes FOR EACH ROW EXECUTE PROCEDURE public.notifyethtxinsertion();
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0031_eth_tx_trigger_with_key_address",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up31).Error
		},
		Rollback: func(db *gorm.DB) error {
			return nil
		},
	})
}
