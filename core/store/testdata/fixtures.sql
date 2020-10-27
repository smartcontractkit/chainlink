-- Password for all encrypted keys is 'password'
-- Don't use any of these keys for anything outside of testing!
INSERT INTO keys (address, json, created_at, updated_at, next_nonce) VALUES (
    E'\\x3cb8e3FD9d27e39a5e9e6852b0e96160061fd4ea',
    '{"address":"3cb8e3fd9d27e39a5e9e6852b0e96160061fd4ea","crypto":{"cipher":"aes-128-ctr","ciphertext":"7515678239ccbeeaaaf0b103f0fba46a979bf6b2a52260015f35b9eb5fed5c17","cipherparams":{"iv":"87e5a5db334305e1e4fb8b3538ceea12"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"d89ac837b5dcdce5690af764762fe349d8162bb0086cea2bc3a4289c47853f96"},"mac":"57a7f4ada10d3d89644f541c91f89b5bde73e15e827ee40565e2d1f88bb0ac96"},"id":"c8cb9bc7-0a51-43bd-8348-8a67fd1ec52c","version":3}',
    '2019-01-01',
    '2019-01-01',
    0
);

INSERT INTO "public"."encrypted_ocr_key_bundles"("id","on_chain_signing_address","off_chain_public_key","encrypted_private_keys","created_at","updated_at","config_public_key")
VALUES
(
    DECODE('0e5abb138f8764de22558c9a9d99120c4864c314c3d1509d7775b88edb186cc6','hex'),
    DECODE('66e0c032bcb30a5f107eadf363b5ea09e75e03f4','hex'),
    DECODE('1062f3aac99f6e6da5c399afe52dad8a3c882505e1dbc78c6d910ec06939a01d','hex'),
    E'{"kdf": "scrypt", "mac": "cc17a398106217d5c7550b6bd67ea38431e952b1918872db76042e9914eead9e", "cipher": "aes-128-ctr", "kdfparams": {"n": 262144, "p": 1, "r": 8, "salt": "713f5dff97dc7ba21026052fedb4b6b6c4fb3a613b5fd5d8bedc242511fcb7aa", "dklen": 32}, "ciphertext": "dfb483cc9b78dcbb81aa82cda43e8df6c991d6b238ccd4fa3d5ebe8afe079e0e7d20cf5cc09ddf42d8c23c7ce80db86a4b102cd60ac1284634f17761582da5590eda5623f96f2208e740293b1737bda7154ae3fb26c7ce4b6e998e0c7ecfd7b18df2a01da14ec19bbd3c438fd01e1426719beb282a6ffba40a07832255df8e26cb7ce33033b3747128dbee4129e231b18392c08259e2458885cc08797d10a8a2e200ee68a9edcfb58a9fa8b09c6fb226c84d45a53301d92943cbccb0891ba4ffaccd6e6283460c3ffaeb6f0988025a261478cbc6523f0ec1b929657f9a0fbaf9322e4fbbff5c283b0f41794105b3c97fb11cebbc173351c08348a7f90d191ce10446bef7bb66f737427117770e4f57e46ee7ec4f1bb1b086c7550bd05391df4a6ee4860434feef21898fe2b790eba69aa6dbade75a79a5720969d296d897b23419c139a49e934385987fb6b2d697833f", "cipherparams": {"iv": "f29542d1652ef17fb386388cb9d7e8e5"}}',E'2020-10-27 10:18:58.35498+00',E'2020-10-27 10:18:58.35498+00',
    DECODE('376f20b9d55e2755f38412e494938502a8ad1f44df3302394be66f2e6eef003c','hex')
);

INSERT INTO encrypted_p2p_keys (peer_id, pub_key, encrypted_priv_key, created_at, updated_at) VALUES (
    '12D3KooWGKbY6ymznHbQtqYq28Bdk9eifkxskuYoG4bHgccRoLvc',
    E'\\x60a314a01fd0050e54c3130305b05ba2dbaab274e8fa304e0efa7fd08f8003a1',
    '{"kdf": "scrypt", "mac": "195438588713832af8d9575ba07a615f9768ab3b891d470debcc52d8914a51e2", "cipher": "aes-128-ctr", "kdfparams": {"n": 262144, "p": 1, "r": 8, "salt": "1617b812658e30629409fa315711f41009b312390749de74bcc460a3c49650f2", "dklen": 32}, "ciphertext": "230c45e546ad733c6f7c1dcaad71398e34f9228d5c876cd3f34478e65af270d154422bf577658a786ddf41a936327177aaff522f187697380d778370b4292e18fa100e52", "cipherparams": {"iv": "65414d058423354be330c80e8e89d59a"}}',
    '2020-10-07',
    '2020-10-07'
);

INSERT INTO users (email, hashed_password, token_secret, created_at, updated_at) VALUES (
    'apiuser@chainlink.test',
    '$2a$10$bbwErtZcZ6qQvRsfBiY2POvuY6D4lwj/Vxq/PcVAL6o64nRaPgaEa', -- hash of literal string 'password'
    '1eCP/w0llVkchejFaoBpfIGaLRxZK54lTXBCT22YLW+pdzE4Fafy/XO5LoJ2uwHi',
    '2019-01-01',
    '2019-01-01'
);
