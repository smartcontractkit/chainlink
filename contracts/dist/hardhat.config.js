"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
require("@nomiclabs/hardhat-waffle");
/**
 * @type import('hardhat/config').HardhatUserConfig
 */
exports.default = {
    solidity: "0.8.4",
    paths: {
        artifacts: "./artifacts",
        cache: "./cache",
        sources: "./src",
        tests: "./test",
    },
    overrides: {
        "src/v0.8/*": {
            version: "0.8.4",
            optimizer: {
                enabled: true,
                runs: 200,
            },
        },
    },
};
//# sourceMappingURL=hardhat.config.js.map