INSERT INTO keys (address, json, created_at, updated_at, next_nonce) VALUES (
    E'\\x3cb8e3FD9d27e39a5e9e6852b0e96160061fd4ea',
    '{"address":"3cb8e3fd9d27e39a5e9e6852b0e96160061fd4ea","crypto":{"cipher":"aes-128-ctr","ciphertext":"7515678239ccbeeaaaf0b103f0fba46a979bf6b2a52260015f35b9eb5fed5c17","cipherparams":{"iv":"87e5a5db334305e1e4fb8b3538ceea12"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"d89ac837b5dcdce5690af764762fe349d8162bb0086cea2bc3a4289c47853f96"},"mac":"57a7f4ada10d3d89644f541c91f89b5bde73e15e827ee40565e2d1f88bb0ac96"},"id":"c8cb9bc7-0a51-43bd-8348-8a67fd1ec52c","version":3}',
    '2019-01-01',
    '2019-01-01',
    0
);
INSERT INTO users (email, hashed_password, token_secret, created_at, updated_at) VALUES (
    'apiuser@chainlink.test',
    '$2a$10$bbwErtZcZ6qQvRsfBiY2POvuY6D4lwj/Vxq/PcVAL6o64nRaPgaEa', -- hash of literal string 'password'
    '1eCP/w0llVkchejFaoBpfIGaLRxZK54lTXBCT22YLW+pdzE4Fafy/XO5LoJ2uwHi',
    '2019-01-01',
    '2019-01-01'
);
