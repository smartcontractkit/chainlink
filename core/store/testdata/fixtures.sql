-- Password for all encrypted keys is 'password'
-- Don't use any of these keys for anything outside of testing!
INSERT INTO keys (address, json, created_at, updated_at, next_nonce) VALUES (
    E'\\x3cb8e3FD9d27e39a5e9e6852b0e96160061fd4ea',
    '{"address":"3cb8e3fd9d27e39a5e9e6852b0e96160061fd4ea","crypto":{"cipher":"aes-128-ctr","ciphertext":"7515678239ccbeeaaaf0b103f0fba46a979bf6b2a52260015f35b9eb5fed5c17","cipherparams":{"iv":"87e5a5db334305e1e4fb8b3538ceea12"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"d89ac837b5dcdce5690af764762fe349d8162bb0086cea2bc3a4289c47853f96"},"mac":"57a7f4ada10d3d89644f541c91f89b5bde73e15e827ee40565e2d1f88bb0ac96"},"id":"c8cb9bc7-0a51-43bd-8348-8a67fd1ec52c","version":3}',
    '2019-01-01',
    '2019-01-01',
    0
);

INSERT INTO encrypted_ocr_key_bundles (id, on_chain_signing_address, off_chain_public_key, encrypted_private_keys, created_at, updated_at) VALUES (
    E'\\x1119b1811c471df8738ecaace1224eeaf2fecaf82edb1fa345e3ddc3f8882e80',
    E'\\xbae6d72ebc77e9da0d09f837b4bea44fd3692753',
    E'\\x5713036ef61523a7c82b2d062c8d7c46a03d867b69837e521bb1c074074ee3e6',
    '{"kdf": "scrypt", "mac": "7b2d9b26014537f60f1f26823746d0203561ac3f91504ed45caea4591689f6f6", "cipher": "aes-128-ctr", "kdfparams": {"n": 262144, "p": 1, "r": 8, "salt": "6dd1276355c6f9fa9a70e759300c62e63d8a22066b587c8ae9330c7e78df1921", "dklen": 32}, "ciphertext": "45e6c64df2bafc6bb1ab3d7342b0cd42e4aea1cb01256ca820ef96e3cac11919099fab67d1500d83459afc32b496a8e75886210acba1c58b10cf37ea479c9d344b666f674e1fd26b96076c7faa6a19993e2ceaf2daa4f5545a011476eaf1591d959fb6af0cd36802e49774942a699f781a3749a59e9484cf1528faea2d26c60c5211519965d7be347fde12732174b17f396f9faeabe6e52fadb346a571e4a23e44082358e8084dcabdcff450017931c07cb61a720b4db72186a02329fd720f59104c2041452768f93dad5d78a12d9cac5ecf03f2fd794d50a8fd1247aeb0b4fba33d1d021a6b89a316e7595a385d85844a5c9201f854e38d4d02b1ca8bc2505038cc665821d90925da39e1c53c1d252a94146d05f2a44a43ade0657c74e766a448066b844e3e67aeeaf32819d2cf8fd6153dea98c6d26221561b014ab45b03e293b51866b95a6489d136f44640", "cipherparams": {"iv": "46f15f8d634ef202df20e6174f59e28a"}}',
    '2020-10-07',
    '2020-10-07'
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
