const { Wallets } = require('fabric-network');
const path = require('path');
const fs = require('fs');

async function importOrg2AdminIdentity() {
    try {
        // Paths to Org2 connection profile and wallet directory
        const ccpPath = path.resolve(__dirname, '..', 'organizations', 'peerOrganizations', 'org2.example.com', 'connection-org2.json');
        const walletPath = path.join(process.cwd(), 'wallet');

        // Load connection profile and initialize wallet
        const ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));
        const wallet = await Wallets.newFileSystemWallet(walletPath);

        // Paths to Admin identity files
        const identityPath = path.join(__dirname, '..', 'organizations', 'peerOrganizations', 'org2.example.com', 'users', 'Admin@org2.example.com', 'msp');
        const certPath = path.join(identityPath, 'signcerts', 'Admin@org2.example.com-cert.pem');
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
            mspId: 'Org2MSP', // MSP ID for Org2
            type: 'X.509'
        };

        // Check if identity already exists in wallet
        const identityExists = await wallet.get('Admin@org2.example.com');
        if (identityExists) {
            console.log('Identity Admin@org2.example.com already exists in the wallet.');
            return;
        }

        // Add Admin identity to wallet
        await wallet.put('Admin@org2.example.com', adminIdentity);
        console.log('Admin identity for Org2 added to the wallet successfully.');
    } catch (error) {
        console.error(`Error adding Admin@org2.example.com identity: ${error.message}`);
    }
}

// Call the function
importOrg2AdminIdentity();
