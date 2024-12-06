const { Wallets } = require('fabric-network');
const path = require('path');
const fs = require('fs');

async function importOrg1AdminIdentity() {
    try {
        // Paths to Org1 connection profile and wallet directory
        const ccpPath = path.resolve(__dirname, '..', 'organizations', 'peerOrganizations', 'org1.example.com', 'connection-org1.json');
        const walletPath = path.join(process.cwd(), 'wallet');

        // Load connection profile and initialize wallet
        const ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));
        const wallet = await Wallets.newFileSystemWallet(walletPath);

        // Paths to Admin identity files
        const identityPath = path.join(__dirname, '..', 'organizations', 'peerOrganizations', 'org1.example.com', 'users', 'Admin@org1.example.com', 'msp');
        const certPath = path.join(identityPath, 'signcerts', 'Admin@org1.example.com-cert.pem');
        const privateKeyPath = path.join(identityPath, 'keystore', 'priv_sk');

        // Ensure the certificate and private key files exist
        if (!fs.existsSync(certPath)) {
            throw new Error(`Certificate file not found at ${certPath}`);
        }
        if (!fs.existsSync(privateKeyPath)) {
            throw new Error(`Private key file not found at ${privateKeyPath}`);
        }

        // Read certificate and private key
        const certificate = fs.readFileSync(certPath).toString();
        const privateKey = fs.readFileSync(privateKeyPath).toString();

        // Define Admin identity
        const adminIdentity = {
            credentials: {
                certificate: certificate,
                privateKey: privateKey
            },
            mspId: 'Org1MSP', // MSP ID for Org1
            type: 'X.509'
        };

        // Check if identity already exists in wallet
        const identityExists = await wallet.get('Admin@org1.example.com');
        if (identityExists) {
            console.log('Identity Admin@org1.example.com already exists in the wallet.');
            return;
        }

        // Add Admin identity to wallet
        await wallet.put('Admin@org1.example.com', adminIdentity);
        console.log('Admin identity for Org1 added to the wallet successfully.');
    } catch (error) {
        console.error(`Error adding Admin@org1.example.com identity: ${error.message}`);
    }
}

// Call the function
importOrg1AdminIdentity();
