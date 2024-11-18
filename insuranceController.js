const { Gateway, Wallets } = require('fabric-network');
const path = require('path');
const fs = require('fs');

const ccpPath = path.resolve(__dirname, '..', 'organizations', 'peerOrganizations', 'org1.example.com', 'connection-org1.json');
const walletPath = path.join(process.cwd(), 'wallet');

async function connectToNetwork() {
    const ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));
    const wallet = await Wallets.newFileSystemWallet(walletPath);

    // Import Admin identity (if not already imported)
    const identityPath = path.join(__dirname, '..', 'organizations', 'peerOrganizations', 'org1.example.com', 'users', 'Admin@org1.example.com', 'msp');

    // Ensure the certificate and private key exist
    const certPath = path.join(identityPath, 'signcerts', 'Admin@org1.example.com-cert.pem');
    const privateKeyPath = path.join(identityPath, 'keystore', 'priv_sk');

    if (!fs.existsSync(certPath)) {
        throw new Error(`Certificate file not found at ${certPath}`);
    }
    if (!fs.existsSync(privateKeyPath)) {
        throw new Error(`Private key file not found at ${privateKeyPath}`);
    }

    const certificate = fs.readFileSync(certPath).toString();
    const privateKey = fs.readFileSync(privateKeyPath).toString();

    const adminIdentity = {
        credentials: {
            certificate: certificate,
            privateKey: privateKey
        },
        mspId: 'Org1MSP', // This is the default MSP ID for org1
        type: 'X.509' // Identity type should be X.509
    };

    // Check if wallet already has the identity
    const identityExists = await wallet.get('Admin@org1.example.com');
    if (!identityExists) {
        console.log('Adding Admin identity to wallet...');
        await wallet.put('Admin@org1.example.com', adminIdentity); // Add Admin identity to wallet
    }

    const gateway = new Gateway();
    await gateway.connect(ccp, {
        wallet,
        identity: 'Admin@org1.example.com',  // Use Admin identity
        discovery: { enabled: true, asLocalhost: true }
    });

    return gateway.getNetwork('mychannel');
}

module.exports = {
    registerPolicy: async (req, res) => {
        const { policyID, holderName } = req.body;
        try {
            console.log("Connecting to network...");
            const network = await connectToNetwork();
            console.log("Connected to network");

            const contract = network.getContract('insurecc');
            await contract.submitTransaction('RegisterPolicy', policyID, holderName);
            res.status(200).send(`Policy ${policyID} registered successfully`);
        } catch (error) {
            console.error("Error registering policy:", error);
            res.status(500).json({ error: error.message });
        }
    },

    queryPolicy: async (req, res) => {
        const { policyID } = req.params;
        try {
            const network = await connectToNetwork();
            const contract = network.getContract('insurecc');

            const result = await contract.evaluateTransaction('QueryPolicy', policyID);
            res.status(200).json(JSON.parse(result.toString()));
        } catch (error) {
            console.error("Error querying policy:", error);
            res.status(500).json({ error: error.message });
        }
    },

    processClaim: async (req, res) => {
        const { claimID, policyID, amount } = req.body;
        try {
            const network = await connectToNetwork();
            const contract = network.getContract('insurecc');

            await contract.submitTransaction('ProcessClaim', claimID, policyID, amount.toString());
            res.status(200).send(`Claim ${claimID} processed successfully`);
        } catch (error) {
            console.error("Error processing claim:", error);
            res.status(500).json({ error: error.message });
        }
    }
};
