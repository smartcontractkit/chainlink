INSERT INTO evm.key_states (address,disabled,created_at,updated_at,evm_chain_id) 
SELECT address,disabled,created_at,updated_at,'{{ .ChainID }}' FROM {{ .Schema }}.key_states;

DROP TABLE {{ .Schema }}.key_states;