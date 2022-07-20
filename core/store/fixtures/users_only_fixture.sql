INSERT INTO users (email, hashed_password, token_hashed_secret, role, created_at, updated_at) VALUES
(
    'apiuser@chainlink.test',
    '$2a$10$bUMgzjxp1Jtaq4nt5ICPB.fWsfVP6FpdxXB1ZOsI0t9je0JOIkpRW', -- hash of literal string '16charlengthp4SsW0rD1!@#_'
    '1eCP/w0llVkchejFaoBpfIGaLRxZK54lTXBCT22YLW+pdzE4Fafy/XO5LoJ2uwHi',
    'admin',
    '2019-01-01',
    '2019-01-01'
);