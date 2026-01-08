import { HardhatUserConfig } from "hardhat/config";
import "@nomiclabs/hardhat-ethers";

const config: HardhatUserConfig = {
  solidity: "0.8.20",
  networks: {
    private: {
      url: "http://127.0.0.1:8545",
      chainId: 31337
    }
  }
};

export default config;
