INSERT INTO users (email, hashed_password, token_hashed_secret, role, created_at, updated_at) VALUES
(
    'apiuser@chainlink.test',
    '$2a$10$bUMgzjxp1Jtaq4nt5ICPB.fWsfVP6FpdxXB1ZOsI0t9je0JOIkpRW', -- hash of literal string '16charlengthp4SsW0rD1!@#_'
    '1eCP/w0llVkchejFaoBpfIGaLRxZK54lTXBCT22YLW+pdzE4Fafy/XO5LoJ2uwHi',
    'admin',
    '2019-01-01',
    '2019-01-01'
),
(
    'apiuser-edit@chainlink.test',
    '$2a$10$bUMgzjxp1Jtaq4nt5ICPB.fWsfVP6FpdxXB1ZOsI0t9je0JOIkpRW', -- hash of literal string '16charlengthp4SsW0rD1!@#_'
    '1eCP/w0llVkchejFaoBpfIGaLRxZK54lTXBCT22YLW+pdzE4Fafy/XO5LoJ2uwHi',
    'edit',
    '2018-01-01',
    '2018-01-01'
),
(
    'apiuser-edit-minimal@chainlink.test',
    '$2a$10$bUMgzjxp1Jtaq4nt5ICPB.fWsfVP6FpdxXB1ZOsI0t9je0JOIkpRW', -- hash of literal string '16charlengthp4SsW0rD1!@#_'
    '1eCP/w0llVkchejFaoBpfIGaLRxZK54lTXBCT22YLW+pdzE4Fafy/XO5LoJ2uwHi',
    'run',
    '2017-01-01',
    '2017-01-01'
),
(
    'apiuser-view-only@chainlink.test',
    '$2a$10$bUMgzjxp1Jtaq4nt5ICPB.fWsfVP6FpdxXB1ZOsI0t9je0JOIkpRW', -- hash of literal string '16charlengthp4SsW0rD1!@#_'
    '1eCP/w0llVkchejFaoBpfIGaLRxZK54lTXBCT22YLW+pdzE4Fafy/XO5LoJ2uwHi',
    'view',
    '2016-01-01',
    '2016-01-01'
);