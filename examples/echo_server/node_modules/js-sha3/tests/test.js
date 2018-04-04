(function (sha3_512, sha3_384, sha3_256, sha3_224, keccak_512, keccak_384, keccak_256, keccak_224) {
  Array.prototype.toHexString = ArrayBuffer.prototype.toHexString = function () {
    var array = new Uint8Array(this);
    var hex = '';
    for (var i = 0;i < array.length;++i) {
      var c = array[i].toString('16');
      hex += c.length == 1 ? '0' + c : c;
    }
    return hex;
  };

  function runTestCases(methods, testCases) {
    methods.forEach(function (method) {
      describe('#' + method.name, function () {
        var methodTestCases = testCases[method.name];
        for (var testCaseName in methodTestCases) {
          (function (testCaseName) {
            var testCase = methodTestCases[testCaseName];
            context('when ' + testCaseName, function () {
              for (var hash in testCase) {
                (function (message, hash) {
                  it('should be equal', function () {
                    expect(method.call(message)).to.be(hash);
                  });
                })(testCase[hash], hash);
              }
            });
          })(testCaseName);
        }
      });
    });
  }

  var methods = [
    {
      name: 'sha3_512',
      call: sha3_512
    },
    {
      name: 'sha3_384',
      call: sha3_384
    },
    {
      name: 'sha3_256',
      call: sha3_256
    },
    {
      name: 'sha3_224',
      call: sha3_224
    },
    {
      name: 'keccak_512',
      call: keccak_512
    },
    {
      name: 'keccak_384',
      call: keccak_384
    },
    {
      name: 'keccak_256',
      call: keccak_256
    },
    {
      name: 'keccak_224',
      call: keccak_224
    }
  ];

  var testCases = {
    sha3_512: {
      'ascii': {
        'a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26': '',
        '01dedd5de4ef14642445ba5f5b97c15e47b9ad931326e4b0727cd94cefc44fff23f07bf543139939b49128caf436dc1bdee54fcb24023a08d9403f9b4bf0d450': 'The quick brown fox jumps over the lazy dog',
        '18f4f4bd419603f95538837003d9d254c26c23765565162247483f65c50303597bc9ce4d289f21d1c2f1f458828e33dc442100331b35e7eb031b5d38ba6460f8': 'The quick brown fox jumps over the lazy dog.'
      },
      'ascii more than 128 bytes': {
        '4f8bcf3a60d3ee56a0bd405c3e6bb37dac44b6781c41bf76c91a5d8e621d472b7b13b8806d88914af3d97585df996363ebe17566d5dfeb6f4884a7949ba8263d': 'The MD5 message-digest algorithm is a widely used cryptographic hash function producing a 128-bit (16-byte) hash value, typically expressed in text format as a 32 digit hexadecimal number. MD5 has been utilized in a wide variety of cryptographic applications, and is also commonly used to verify data integrity.'
      },
      'UTF-8': {
        '059bbe2efc50cc30e4d8ec5a96be697e2108fcbf9193e1296192eddabc13b143c0120d059399a13d0d42651efe23a6c1ce2d1efb576c5b207fa2516050505af7': '中文',
        '35dfaf82d2ce4be79393dc90e327b4dd15b1c150d8a30f59d8d1b42ca4fc3c87f50b77da36acccf9dc76494d07fc57cfcc9470e627c38f95bce4deab311b87e0': 'aécio',
        '33ef254289f36527c93cd203ef1973aec1eff7475c23fa842c3092b0d30965d13b0805c61d0aa92c51245c56bfe26978c35c00f26eb558a043982043ee8b178c': '𠜎'
      },
      'UTF-8 more than 128 bytes': {
        'accb127bb24b0ffbb7550dc637222d2f78538a8a186c98bc5efdad685b9b396639f34148bf0b94ed470f0e9c3665dc3b4c1cb321bacd32dd317a646295e073d9': '訊息摘要演算法第五版（英語：Message-Digest Algorithm 5，縮寫為MD5），是當前電腦領域用於確保資訊傳輸完整一致而廣泛使用的雜湊演算法之一',
        '5b70eaad083f1b86fd535b6812e02f5f2876a4bd8b43aede8d62ae71bb1743ebd919dc41be56d73ba45b67b2876ff215d0575788560e7b0c92b879f8a2fc3111': '訊息摘要演算法第五版（英語：Message-Digest Algorithm 5，縮寫為MD5），是當前電腦領域用於確保資訊傳輸完整一致而廣泛使用的雜湊演算法之一（又譯雜湊演算法、摘要演算法等），主流程式語言普遍已有MD5的實作。'
      },
      'special length': {
        'bce9da5b408846edd5bec9f26c2dee9bd835215c3f2b3876197067d87bc4d1af0cd97f94fda59761a0d804fe82383be2c6c4886fbb82e005fcf899449029f221' :'012345678901234567890123456789012345678901234567890123456789012345678901',
        '8bdcb85e6b52c29fafac0d3daf65492f2e3499e066da1a095a65eb1144849a26b2790a8b39c2a7fb747456f749391d953841a61cb13289f9806f04981c180a86' :'01234567890123456789012345678901234567890123456789012345678901234567890' 
      },
      'Array': {
        'a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26': [],
        '01dedd5de4ef14642445ba5f5b97c15e47b9ad931326e4b0727cd94cefc44fff23f07bf543139939b49128caf436dc1bdee54fcb24023a08d9403f9b4bf0d450': [84, 104, 101, 32, 113, 117, 105, 99, 107, 32, 98, 114, 111, 119, 110, 32, 102, 111, 120, 32, 106, 117, 109, 112, 115, 32, 111, 118, 101, 114, 32, 116, 104, 101, 32, 108, 97, 122, 121, 32, 100, 111, 103]
      },
      'Uint8Array': {
        'a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26': new Uint8Array([]),
        '01dedd5de4ef14642445ba5f5b97c15e47b9ad931326e4b0727cd94cefc44fff23f07bf543139939b49128caf436dc1bdee54fcb24023a08d9403f9b4bf0d450': new Uint8Array([84, 104, 101, 32, 113, 117, 105, 99, 107, 32, 98, 114, 111, 119, 110, 32, 102, 111, 120, 32, 106, 117, 109, 112, 115, 32, 111, 118, 101, 114, 32, 116, 104, 101, 32, 108, 97, 122, 121, 32, 100, 111, 103])
      },
      'ArrayBuffer': {
        'a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26': new ArrayBuffer(0)
      }
    },
    sha3_384: {
      'ascii': {
        '0c63a75b845e4f7d01107d852e4c2485c51a50aaaa94fc61995e71bbee983a2ac3713831264adb47fb6bd1e058d5f004': '',
        '7063465e08a93bce31cd89d2e3ca8f602498696e253592ed26f07bf7e703cf328581e1471a7ba7ab119b1a9ebdf8be41': 'The quick brown fox jumps over the lazy dog',
        '1a34d81695b622df178bc74df7124fe12fac0f64ba5250b78b99c1273d4b080168e10652894ecad5f1f4d5b965437fb9': 'The quick brown fox jumps over the lazy dog.'
      },
      'ascii more than 128 bytes': {
        'ca6b121a6060bc85de05e5a8d70577838fad2481b092c8263d6f7bcbe5148740f0c7f9c4dc27061339570496956aaef6': 'The MD5 message-digest algorithm is a widely used cryptographic hash function producing a 128-bit (16-byte) hash value, typically expressed in text format as a 32 digit hexadecimal number. MD5 has been utilized in a wide variety of cryptographic applications, and is also commonly used to verify data integrity.'
      },
      'UTF-8': {
        '9fb5b99e3c546f2738dcd50a14e9aef9c313800c1bf8cf76bc9b2c3a23307841364c5a2d0794702662c5796fb72f5432': '中文',
        '70b447f1bd5ce5a4753ccf7a3697eca0315954774374bc1042aff19582ccc32d5067f7da6c2bea9d6d344e11924cbe72': 'aécio',
        '7add8d544b0a7cf188b54b1697a046f77e49d5f292e7ffe56feeed90a500b0bf026b9b68892888a1bafb9f8cb89ed874': '𠜎'
      },
      'UTF-8 more than 128 bytes': {
        '7d0f80fe5c79a04a2a37a30a440e0cc068eb78fe6c3182246ede29645c144b5d33c44607cb2c3111ba77ffc66107f1cd': '訊息摘要演算法第五版（英語：Message-Digest Algorithm 5，縮寫為MD5），是當前電腦領域用於確保資訊傳輸完整一致而廣泛使用的雜湊演算法之一',
        'e344b95c6a961a27793eff00fa5103ef78b4180fe41c93fc60a31aff49b3b5e95a92c84fda9a6c80fa403b7df58db59f': '訊息摘要演算法第五版（英語：Message-Digest Algorithm 5，縮寫為MD5），是當前電腦領域用於確保資訊傳輸完整一致而廣泛使用的雜湊演算法之一（又譯雜湊演算法、摘要演算法等），主流程式語言普遍已有MD5的實作。'
      }
    },
    sha3_256: {
      'ascii': {
        'a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a': '',
        '69070dda01975c8c120c3aada1b282394e7f032fa9cf32f4cb2259a0897dfc04': 'The quick brown fox jumps over the lazy dog',
        'a80f839cd4f83f6c3dafc87feae470045e4eb0d366397d5c6ce34ba1739f734d': 'The quick brown fox jumps over the lazy dog.'
      },
      'ascii more than 128 bytes': {
        'fa198893674a0bf9fb35980504e8cefb250aabd2311a37e5d2205f07fb023d36': 'The MD5 message-digest algorithm is a widely used cryptographic hash function producing a 128-bit (16-byte) hash value, typically expressed in text format as a 32 digit hexadecimal number. MD5 has been utilized in a wide variety of cryptographic applications, and is also commonly used to verify data integrity.'
      },
      'UTF-8': {
        'ac5305da3d18be1aed44aa7c70ea548da243a59a5fd546f489348fd5718fb1a0': '中文',
        '65c756408eb6c35a1ffa2d7e09711bdc9f0b28716b1376223844a2b4c52b6718': 'aécio',
        'babe9afc555b0311700dfb0b5b6296d49347b3d770480baedfcdc47a4aea6e82': '𠜎'
      },
      'UTF-8 more than 128 bytes': {
        '4b2f36e4320b86e6ead0ad001e47e6d9e7fcf0044cd5a5fd65490a633c0372a4': '訊息摘要演算法第五版（英語：Message-Digest Algorithm 5，縮寫為MD5），是當前電腦領域用於確保資訊傳輸完整一致而廣泛使用的雜湊演算法之一',
        '558a7f843b1ac5e7a8bbef90357876bcce0612992d0dfa2907e95521612f507f': '訊息摘要演算法第五版（英語：Message-Digest Algorithm 5，縮寫為MD5），是當前電腦領域用於確保資訊傳輸完整一致而廣泛使用的雜湊演算法之一（又譯雜湊演算法、摘要演算法等），主流程式語言普遍已有MD5的實作。'
      }
    },
    sha3_224: {
      'ascii': {
        '6b4e03423667dbb73b6e15454f0eb1abd4597f9a1b078e3f5b5a6bc7': '',
        'd15dadceaa4d5d7bb3b48f446421d542e08ad8887305e28d58335795': 'The quick brown fox jumps over the lazy dog',
        '2d0708903833afabdd232a20201176e8b58c5be8a6fe74265ac54db0': 'The quick brown fox jumps over the lazy dog.'
      },
      'ascii more than 128 bytes': {
        '06885009a28e43e15bf1af718561ad211515a27b542eabc36764a0ca': 'The MD5 message-digest algorithm is a widely used cryptographic hash function producing a 128-bit (16-byte) hash value, typically expressed in text format as a 32 digit hexadecimal number. MD5 has been utilized in a wide variety of cryptographic applications, and is also commonly used to verify data integrity.'
      },
      'UTF-8': {
        '106d169e10b61c2a2a05554d3e631ec94467f8316640f29545d163ee': '中文',
        'b16bad54608dc01864a5d7510d4c19b09f3a0f39cfc4ba1e53aa952a': 'aécio',
        'f59253c41cb87e5cd953311656716cb5b64dbafc9e8155f0dd68123c': '𠜎'
      },
      'UTF-8 more than 128 bytes': {
        '135c13deb71fdf6fb77b52b720c43ddd6ce7467f9147a74248557114': '訊息摘要演算法第五版（英語：Message-Digest Algorithm 5，縮寫為MD5），是當前電腦領域用於確保資訊傳輸完整一致而廣泛使用的雜湊演算法之一',
        'bd05581e02445c53e05aad2014f6a3819d77a9dff918b8c6bf60bd06': '訊息摘要演算法第五版（英語：Message-Digest Algorithm 5，縮寫為MD5），是當前電腦領域用於確保資訊傳輸完整一致而廣泛使用的雜湊演算法之一（又譯雜湊演算法、摘要演算法等），主流程式語言普遍已有MD5的實作。'
      }
    },
    keccak_512: {
      'ascii': {
        '0eab42de4c3ceb9235fc91acffe746b29c29a8c366b7c60e4e67c466f36a4304c00fa9caf9d87976ba469bcbe06713b435f091ef2769fb160cdab33d3670680e': '',
        'd135bb84d0439dbac432247ee573a23ea7d3c9deb2a968eb31d47c4fb45f1ef4422d6c531b5b9bd6f449ebcc449ea94d0a8f05f62130fda612da53c79659f609': 'The quick brown fox jumps over the lazy dog',
        'ab7192d2b11f51c7dd744e7b3441febf397ca07bf812cceae122ca4ded6387889064f8db9230f173f6d1ab6e24b6e50f065b039f799f5592360a6558eb52d760': 'The quick brown fox jumps over the lazy dog.'
      },
      'ascii more than 128 bytes': {
        '10dcbf6389980ce3594547939bbc685363d28adbd6a05bc4abd7fc62e7693a1f6e33569fed5a380bfecb56ae811d25939b95823f39bb0f16a08740629d066d43': 'The MD5 message-digest algorithm is a widely used cryptographic hash function producing a 128-bit (16-byte) hash value, typically expressed in text format as a 32 digit hexadecimal number. MD5 has been utilized in a wide variety of cryptographic applications, and is also commonly used to verify data integrity.'
      },
      'UTF-8': {
        '2f6a1bd50562230229af34b0ccf46b8754b89d23ae2c5bf7840b4acfcef86f87395edc0a00b2bfef53bafebe3b79de2e3e01cbd8169ddbb08bde888dcc893524': '中文',
        'c452ec93e83d4795fcab62a76eed0d88f2231a995ce108ac8f130246f87c4a11cb18a2c1a688a5695906a6f863e71bbe8997c6610319ab97f12d2e5bf0afe458': 'aécio',
        '8a2d72022ce19d989dbe6a0017faccbf5dc2e22c162d1c5eb168864d32dd1a71e1b4782652c148cf6ca47b77a72c96fff682e72bdfef0566d4b7cca3c9ccc59d': '𠜎'
      },
      'UTF-8 more than 128 bytes': {
        '6a67c28aa1946ca1be8382b861aac4aaf20052f495db9b6902d13adfa603eaba5d169f8896b86d461b2949283eb98e503c3f0640188ea7d6731526fc06568d37': '訊息摘要演算法第五版（英語：Message-Digest Algorithm 5，縮寫為MD5），是當前電腦領域用於確保資訊傳輸完整一致而廣泛使用的雜湊演算法之一',
        'd04ff5b0e85e9968be2a4d4e133c15c7ccee7497198bb651599a97d11d00bca6048d329ab75aa454566cd532648fa1cb4551985d9d645de9fa43a311a9ee8e4d': '訊息摘要演算法第五版（英語：Message-Digest Algorithm 5，縮寫為MD5），是當前電腦領域用於確保資訊傳輸完整一致而廣泛使用的雜湊演算法之一（又譯雜湊演算法、摘要演算法等），主流程式語言普遍已有MD5的實作。'
      }
    },
    keccak_384: {
      'ascii': {
        '2c23146a63a29acf99e73b88f8c24eaa7dc60aa771780ccc006afbfa8fe2479b2dd2b21362337441ac12b515911957ff': '',
        '283990fa9d5fb731d786c5bbee94ea4db4910f18c62c03d173fc0a5e494422e8a0b3da7574dae7fa0baf005e504063b3': 'The quick brown fox jumps over the lazy dog',
        '9ad8e17325408eddb6edee6147f13856ad819bb7532668b605a24a2d958f88bd5c169e56dc4b2f89ffd325f6006d820b': 'The quick brown fox jumps over the lazy dog.'
      },
      'ascii more than 128 bytes': {
        'e7ec8976b4d96e43f50ae8ecdcf2d97a56236e6406e8dd00efd0d9abe885659db58a2f4b138a4ecfb1bd0052f6569516': 'The MD5 message-digest algorithm is a widely used cryptographic hash function producing a 128-bit (16-byte) hash value, typically expressed in text format as a 32 digit hexadecimal number. MD5 has been utilized in a wide variety of cryptographic applications, and is also commonly used to verify data integrity.'
      },
      'UTF-8': {
        '743f64bb7544c6ed923be4741b738dde18b7cee384a3a09c4e01acaaac9f19222cdee137702bd3aa05dc198373d87d6c': '中文',
        '08990555e131af8597687614309da4c5053ce866f348544da0a0c2c78c2cc79680ebb57cfbe238286e78ea133a037897': 'aécio',
        '2a80f59abf3111f38a35a3daa25123b495f90e9736bd300e35911d19abdd8806498c581333f198ccbbf2252b57c2925d': '𠜎'
      },
      'UTF-8 more than 128 bytes': {
        'a3b043a2f69e4326a05d478fa4c8aa2bd7612453d775af37665a0b96ef2207cdc74c50cdba1629796a5136fe77300b05': '訊息摘要演算法第五版（英語：Message-Digest Algorithm 5，縮寫為MD5），是當前電腦領域用於確保資訊傳輸完整一致而廣泛使用的雜湊演算法之一',
        '66414c090cc3fe9c396d313cbaa100aefd335e851838b29382568b7f57357ada7c54b8fa8c17f859945bba88b2c2e332': '訊息摘要演算法第五版（英語：Message-Digest Algorithm 5，縮寫為MD5），是當前電腦領域用於確保資訊傳輸完整一致而廣泛使用的雜湊演算法之一（又譯雜湊演算法、摘要演算法等），主流程式語言普遍已有MD5的實作。'
      }
    },
    keccak_256: {
      'ascii': {
        'c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470': '',
        '4d741b6f1eb29cb2a9b9911c82f56fa8d73b04959d3d9d222895df6c0b28aa15': 'The quick brown fox jumps over the lazy dog',
        '578951e24efd62a3d63a86f7cd19aaa53c898fe287d2552133220370240b572d': 'The quick brown fox jumps over the lazy dog.'
      },
      'ascii more than 128 bytes': {
        'af20018353ffb50d507f1555580f5272eca7fdab4f8295db4b1a9ad832c93f6d': 'The MD5 message-digest algorithm is a widely used cryptographic hash function producing a 128-bit (16-byte) hash value, typically expressed in text format as a 32 digit hexadecimal number. MD5 has been utilized in a wide variety of cryptographic applications, and is also commonly used to verify data integrity.'
      },
      'UTF-8': {
        '70a2b6579047f0a977fcb5e9120a4e07067bea9abb6916fbc2d13ffb9a4e4eee': '中文',
        'd7d569202f04daf90432810d6163112b2695d7820da979327ebd894efb0276dc': 'aécio',
        '16a7cc7a58444cbf7e939611910ddc82e7cba65a99d3e8e08cfcda53180a2180': '𠜎'
      },
      'UTF-8 more than 128 bytes': {
        'd1021d2d4c5c7e88098c40f422af68493b4b64c913cbd68220bf5e6127c37a88': '訊息摘要演算法第五版（英語：Message-Digest Algorithm 5，縮寫為MD5），是當前電腦領域用於確保資訊傳輸完整一致而廣泛使用的雜湊演算法之一',
        'ffabf9bba2127c4928d360c9905cb4911f0ec21b9c3b89f3b242bccc68389e36': '訊息摘要演算法第五版（英語：Message-Digest Algorithm 5，縮寫為MD5），是當前電腦領域用於確保資訊傳輸完整一致而廣泛使用的雜湊演算法之一（又譯雜湊演算法、摘要演算法等），主流程式語言普遍已有MD5的實作。'
      }
    },
    keccak_224: {
      'ascii': {
        'f71837502ba8e10837bdd8d365adb85591895602fc552b48b7390abd': '',
        '310aee6b30c47350576ac2873fa89fd190cdc488442f3ef654cf23fe': 'The quick brown fox jumps over the lazy dog',
        'c59d4eaeac728671c635ff645014e2afa935bebffdb5fbd207ffdeab': 'The quick brown fox jumps over the lazy dog.'
      },
      'ascii more than 128 bytes': {
        '8dd58b706e3a08ec4f1f202af39295b38c355a39b23308ade7218a21': 'The MD5 message-digest algorithm is a widely used cryptographic hash function producing a 128-bit (16-byte) hash value, typically expressed in text format as a 32 digit hexadecimal number. MD5 has been utilized in a wide variety of cryptographic applications, and is also commonly used to verify data integrity.'
      },
      'UTF-8': {
        '7bc2a0b6e7e0a055a61e4f731e2944b560f41ff98967dcbf4bbf77a5': '中文',
        '66f3db76bf8cb35726cb278bac412d187c3484ab2083dc50ef5ffb55': 'aécio',
        '3bfa94845726f4cd5cf17d19b5eacac17b3694790e13a76d5c81c7c2': '𠜎'
      },
      'UTF-8 more than 128 bytes': {
        'd59eef8f394ef7d96967bb0bde578785c033f7f0a21913d6ba41ed1b': '訊息摘要演算法第五版（英語：Message-Digest Algorithm 5，縮寫為MD5），是當前電腦領域用於確保資訊傳輸完整一致而廣泛使用的雜湊演算法之一',
        '27123a2a3860d1041d4769778c4b078732bf4300f7e1c56536ab2644': '訊息摘要演算法第五版（英語：Message-Digest Algorithm 5，縮寫為MD5），是當前電腦領域用於確保資訊傳輸完整一致而廣泛使用的雜湊演算法之一（又譯雜湊演算法、摘要演算法等），主流程式語言普遍已有MD5的實作。'
      }
    }
  };

  runTestCases(methods, testCases);

  describe('sha3_512', function () {
    context('#arrayBuffer', function () {
      it('should be equal', function () {
        expect(sha3_512.arrayBuffer('').toHexString()).to.be('a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26');
        expect(sha3_512.buffer('').toHexString()).to.be('a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26');
      });
    });

    context('#hex', function () {
      it('should be equal', function () {
        expect(sha3_512.hex('')).to.be('a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26');
      });
    });

    context('#update', function () {
      it('should be equal', function () {
        expect(sha3_512.update('').hex()).to.be('a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26');
        expect(sha3_512.update('The quick brown fox ').update('jumps over the lazy dog').hex()).to.be('01dedd5de4ef14642445ba5f5b97c15e47b9ad931326e4b0727cd94cefc44fff23f07bf543139939b49128caf436dc1bdee54fcb24023a08d9403f9b4bf0d450');
      });
    });

    context('#create', function () {
      it('should be equal', function () {
        var bytes = [84, 104, 101, 32, 113, 117, 105, 99, 107, 32, 98, 114, 111, 119, 110, 32, 102, 111, 120, 32, 106, 117, 109, 112, 115, 32, 111, 118, 101, 114, 32, 116, 104, 101, 32, 108, 97, 122, 121, 32, 100, 111, 103];
        var hash = sha3_512.create();
        for (var i = 0;i < bytes.length;++i) {
          hash.update([bytes[i]]);
        }
        expect(hash.hex()).to.be('01dedd5de4ef14642445ba5f5b97c15e47b9ad931326e4b0727cd94cefc44fff23f07bf543139939b49128caf436dc1bdee54fcb24023a08d9403f9b4bf0d450');
      });
    });
  });

  describe('#keccak_512', function () {
    context('when special length', function () {
      it('should be equal', function () {
        expect(keccak_512('012345678901234567890123456789012345678901234567890123456789012345678901')).to.be('90b1d032c3bf06dcc78a46fe52054bab1250600224bfc6dfbfb40a7877c55e89bb982799a2edf198568a4166f6736678b45e76b12fac813cfdf0a76714e5eae8');
        expect(keccak_512('01234567890123456789012345678901234567890123456789012345678901234567890')).to.be('3173e7abc754a0b2909410d78986428a9183e996864af02f421d273d9fa1b4e4a5b14e2998b20767712f53a01ff8f6ae2c3e71e51e2c0f24257b03e6da09eb77');
      });
    });
  
    context('when Array', function () {
      it('should be equal', function () {
        expect(keccak_512([])).to.be('0eab42de4c3ceb9235fc91acffe746b29c29a8c366b7c60e4e67c466f36a4304c00fa9caf9d87976ba469bcbe06713b435f091ef2769fb160cdab33d3670680e');
        expect(keccak_512([84, 104, 101, 32, 113, 117, 105, 99, 107, 32, 98, 114, 111, 119, 110, 32, 102, 111, 120, 32, 106, 117, 109, 112, 115, 32, 111, 118, 101, 114, 32, 116, 104, 101, 32, 108, 97, 122, 121, 32, 100, 111, 103])).to.be('d135bb84d0439dbac432247ee573a23ea7d3c9deb2a968eb31d47c4fb45f1ef4422d6c531b5b9bd6f449ebcc449ea94d0a8f05f62130fda612da53c79659f609');
      });
    });

    context('when Uint8Array', function () {
      it('should be equal', function () {
        expect(keccak_512(new Uint8Array([]))).to.be('0eab42de4c3ceb9235fc91acffe746b29c29a8c366b7c60e4e67c466f36a4304c00fa9caf9d87976ba469bcbe06713b435f091ef2769fb160cdab33d3670680e');
        expect(keccak_512(new Uint8Array([84, 104, 101, 32, 113, 117, 105, 99, 107, 32, 98, 114, 111, 119, 110, 32, 102, 111, 120, 32, 106, 117, 109, 112, 115, 32, 111, 118, 101, 114, 32, 116, 104, 101, 32, 108, 97, 122, 121, 32, 100, 111, 103]))).to.be('d135bb84d0439dbac432247ee573a23ea7d3c9deb2a968eb31d47c4fb45f1ef4422d6c531b5b9bd6f449ebcc449ea94d0a8f05f62130fda612da53c79659f609');
      });
    });

    context('when ArrayBuffer', function () {
      it('should be equal', function () {
        expect(keccak_512(new ArrayBuffer(0))).to.be('0eab42de4c3ceb9235fc91acffe746b29c29a8c366b7c60e4e67c466f36a4304c00fa9caf9d87976ba469bcbe06713b435f091ef2769fb160cdab33d3670680e');
      });
    });

    context('when output ArrayBuffer', function () {
      it('should be equal', function () {
        expect(keccak_512.arrayBuffer('').toHexString()).to.be('0eab42de4c3ceb9235fc91acffe746b29c29a8c366b7c60e4e67c466f36a4304c00fa9caf9d87976ba469bcbe06713b435f091ef2769fb160cdab33d3670680e');
        expect(keccak_512.buffer('').toHexString()).to.be('0eab42de4c3ceb9235fc91acffe746b29c29a8c366b7c60e4e67c466f36a4304c00fa9caf9d87976ba469bcbe06713b435f091ef2769fb160cdab33d3670680e');
      });
    });

    context('when output Array', function () {
      it('should be equal', function () {
        expect(keccak_512.array('').toHexString()).to.be('0eab42de4c3ceb9235fc91acffe746b29c29a8c366b7c60e4e67c466f36a4304c00fa9caf9d87976ba469bcbe06713b435f091ef2769fb160cdab33d3670680e');
      });
    });
  });
})(sha3_512, sha3_384, sha3_256, sha3_224, keccak_512, keccak_384, keccak_256, keccak_224);
