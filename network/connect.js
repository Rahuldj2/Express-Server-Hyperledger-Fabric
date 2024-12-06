// network/connect.js

const { Gateway, Wallets } = require('fabric-network');
const path = require('path');
const fs = require('fs');

/**
 * Connects to the Hyperledger Fabric network.
 * @param {string} org - The organization (e.g., 'org1', 'org2').
 * @param {string} user - The identity label (e.g., 'Admin@org1.example.com').
 * @returns {Promise<Network>} - The Fabric network object.
 */
const connectToNetwork = async (org, user) => {
    // Load the connection profile for the specified organization
    const ccpPath = path.resolve(__dirname, '..', 'config', `connection-${org}.json`);
    const ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));

    // Create a new file system based wallet for managing identities
    const walletPath = path.resolve(__dirname, '..', 'wallet');
    const wallet = await Wallets.newFileSystemWallet(walletPath);
    console.log(`Wallet path: ${walletPath}`);

    // Check if the identity exists in the wallet
    const identity = await wallet.get(user);
    if (!identity) {
        throw new Error(`An identity for the user ${user} does not exist in the wallet.`);
    }

    // Create a new gateway for connecting to the peer node
    const gateway = new Gateway();
    await gateway.connect(ccp, {
        wallet,
        identity: user,
        discovery: { enabled: true, asLocalhost: true }
    });

    // Get the network (channel) your contract is deployed to
    const network = await gateway.getNetwork('mychannel');

    return network;
};

module.exports = connectToNetwork;
