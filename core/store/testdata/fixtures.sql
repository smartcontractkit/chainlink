-- Password for all encrypted keys is 'password'
-- Scrypt params are chosen to be completely insecure and very fast to decrypt
-- Don't use any of these keys for anything outside of testing!
INSERT INTO "public"."keys"("address","json","created_at","updated_at","next_nonce","last_used","is_funding")
VALUES
(DECODE('27548a32b9ad5d64c5945eae9da5337bc3169d15','hex'),E'{"id": "1ccf542e-8f4d-48a0-ad1d-b4e6a86d4c6d", "crypto": {"kdf": "scrypt", "mac": "7f31bd05768a184278c4e9f077bcfba7b2003fed585b99301374a1a4a9adff25", "cipher": "aes-128-ctr", "kdfparams": {"n": 2, "p": 1, "r": 8, "salt": "99e83bf0fdeba39bd29c343db9c52d9e0eae536fdaee472d3181eac1968aa1f9", "dklen": 32}, "ciphertext": "ac22fa788b53a5f62abda03cd432c7aee1f70053b97633e78f93709c383b2a46", "cipherparams": {"iv": "6699ba30f953728787e51a754d6f9566"}}, "address": "27548a32b9ad5d64c5945eae9da5337bc3169d15", "version": 3}',E'2020-10-29 10:29:34.553191+00',E'2020-10-29 10:29:34.553191+00',0,NULL,FALSE);

INSERT INTO "public"."encrypted_ocr_key_bundles"("id","on_chain_signing_address","off_chain_public_key","encrypted_private_keys","created_at","updated_at","config_public_key")
VALUES
(DECODE('54f02f2756952ee42874182c8a03d51f048b7fc245c05196af50f9266f8e444a','hex'),DECODE('c135508f4c9ada03e56bb6ad98d724e7f4c93323','hex'),DECODE('a91e8a88584c18ad895a259800fa768a63be8760dcc2924ffd6311833aefb8c5','hex'),E'{"kdf": "scrypt", "mac": "acbd1623b39799eedb1fc75698d8e2986599922930032c15a5a3721247c9b748", "cipher": "aes-128-ctr", "kdfparams": {"n": 2, "p": 1, "r": 8, "salt": "ea4f33d745169327d2cdf9f70945af1b67822282c9c01fc2278fa80d6d8e7795", "dklen": 32}, "ciphertext": "e92467755b4abadf162d5d450d963daebe5d2bed6450a77d7c22b705e4f01300a30714a5b4da9686255f569469dc0ed15b4a4fa0acc5439d4257315d7ba033e8c85b6d1a73e1cfc8d0e668e230d9a17117030851794e549dda99bdae7b06501d3d21762ff7b1f7fa494187effdb43cf611fd619d740bc310bb84ccaa449d65f23f1f264491a72b312d9061cea3d3de87168d835339621b38dbb3723b96a694fd86324d319948b4e061ceacb54ce44421f5bf914c158f4e95bf3da039bd0d257241c738488532d4b7fa5cd23d84a8e41ac6653e4b823a3f3f0eb37896d2efebcc3d6061e42a50703621130077e99b96186029661765c8baad9a1bab646a0a10331cc1caf3b9ab926bd39233f06677249bb7d5f5b0a8cb337a2bdce61f2a666128d7b310659e6b8d7dc3039fb876badc3fe961d46778ab905fed2134876cf82bde966b8fabebbc9629c23812b6c80952c06b032af6", "cipherparams": {"iv": "863123caab3f0ae5b3bff6a113a80095"}}',E'2020-10-29 10:34:25.960967+00',E'2020-10-29 10:34:25.960967+00',DECODE('69a2b241acdeee304040940c458f315e911a63d4d6ec16337b123326a00b951f','hex'));


INSERT INTO "public"."encrypted_p2p_keys"("peer_id","pub_key","encrypted_priv_key","created_at","updated_at")
VALUES
(E'12D3KooWCJUPKsYAnCRTQ7SUNULt4Z9qF8Uk1xadhCs7e9M711Lp',DECODE('24eaaa7f7f8cd6d91bc4a83becedf2bd3650c050d5b680683ae26f0f1e209fdd','hex'),E'{"kdf": "scrypt", "mac": "41957a416ab525a3d1409b0dc7ec2fdd4f14fed9082245c05ae42b71cb2d438b", "cipher": "aes-128-ctr", "kdfparams": {"n": 2, "p": 1, "r": 8, "salt": "032413f1267991b5f2b7d01d5bb912aa9bdf07e1b9b109c45bafb0caa75672bc", "dklen": 32}, "ciphertext": "724604086076ec161831f580a0fbd1c435cddc5a908f37a641c76f401c75f33cc09acefb579d03ca47874645c868515aa044de63e43cbbb19f13273490a7dea46fa421a5", "cipherparams": {"iv": "faed66382c086036966a80ed62cffb77"}}',E'2020-10-29 10:33:50.854527+00',E'2020-10-29 10:33:50.854527+00');

INSERT INTO users (email, hashed_password, token_secret, created_at, updated_at) VALUES (
    'apiuser@chainlink.test',
    '$2a$10$bbwErtZcZ6qQvRsfBiY2POvuY6D4lwj/Vxq/PcVAL6o64nRaPgaEa', -- hash of literal string 'password'
    '1eCP/w0llVkchejFaoBpfIGaLRxZK54lTXBCT22YLW+pdzE4Fafy/XO5LoJ2uwHi',
    '2019-01-01',
    '2019-01-01'
);
