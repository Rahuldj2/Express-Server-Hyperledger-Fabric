// controllers/insuranceController.js

const connectToNetwork = require('../network/connect');

const CONTRACT_NAME = 'registration'; // Name of the Registration chaincode

module.exports = {
    definePolicy: async (req, res) => {
        const { policyID, policyType, coverAmount, premium, startDate, endDate, criteriaJSON, diseasesJSON } = req.body;
        try {
            const network = await connectToNetwork('org2', 'Admin@org2.example.com');
            const contract = network.getContract(CONTRACT_NAME);

            await contract.submitTransaction(
                'DefinePolicy',
                policyID,
                policyType,
                coverAmount.toString(),
                premium.toString(),
                startDate,
                endDate,
                criteriaJSON,
                diseasesJSON
            );

            res.status(200).send(`Policy ${policyID} defined successfully.`);
        } catch (error) {
            console.error('Error defining policy:', error);
            res.status(500).json({ error: error.message });
        }
    },

    queryPolicy: async (req, res) => {
        const { policyID } = req.params;
        try {
            const network = await connectToNetwork('org2', 'Admin@org2.example.com');
            const contract = network.getContract(CONTRACT_NAME);

            const result = await contract.evaluateTransaction('QueryPolicy', policyID);
            res.status(200).json(JSON.parse(result.toString()));
        } catch (error) {
            console.error('Error querying policy:', error);
            res.status(500).json({ error: error.message });
        }
    },

    registerForPolicy: async (req, res) => {
        const { userID, policyID, premiumPaid, isNonSmoker, hasDisease, consent } = req.body;
        try {
            const network = await connectToNetwork('org2', 'Admin@org2.example.com');
            const contract = network.getContract(CONTRACT_NAME);

            await contract.submitTransaction(
                'RegisterForPolicy',
                userID,
                policyID,
                premiumPaid.toString(),
                isNonSmoker.toString(),
                hasDisease.toString(),
                consent.toString()
            );

            res.status(200).send(`User ${userID} registered for policy ${policyID} successfully.`);
        } catch (error) {
            console.error('Error registering for policy:', error);
            res.status(500).json({ error: error.message });
        }
    },

    queryRegistration: async (req, res) => {
        const { userId, policyId } = req.params;
        try {
            const network = await connectToNetwork('org2', 'Admin@org2.example.com');
            const contract = network.getContract(CONTRACT_NAME);

            const result = await contract.evaluateTransaction('QueryRegistration', userId, policyId);
            res.status(200).json(JSON.parse(result.toString()));
        } catch (error) {
            console.error('Error querying registration:', error);
            res.status(500).json({ error: error.message });
        }
    },

    // Upload Health Records
    uploadHealthRecords: async (req, res) => {
        const { id, isNonSmoker, hasDisease } = req.body;

        try {
            const network = await connectToNetwork('org1', 'Admin@org1.example.com');
            const contract = network.getContract(CONTRACT_NAME);

            // Invoke the UploadHealthRecords function
            await contract.submitTransaction(
                'UploadHealthRecords',
                id,
                isNonSmoker.toString(),
                hasDisease.toString()
            );

            res.status(200).json({ message: `Health records uploaded successfully for ID: ${id}` });
        } catch (error) {
            console.error('Error uploading health records:', error);
            res.status(500).json({ error: error.message });
        }
    },

    // Query Health Records
    queryHealthRecords: async (req, res) => {
        const { id } = req.params;

        try {
            const network = await connectToNetwork('org2', 'Admin@org2.example.com');
            const contract = network.getContract(CONTRACT_NAME);

            // Invoke the QueryHealthRecords function
            const result = await contract.evaluateTransaction('QueryHealthRecords', id);
            const healthRecord = JSON.parse(result.toString());

            res.status(200).json(healthRecord);
        } catch (error) {
            console.error('Error querying health records:', error);
            res.status(500).json({ error: error.message });
        }
    },

    // Query Health Records
    queryHealthRecordsOrg1: async (req, res) => {
        const { id } = req.params;

        try {
            const network = await connectToNetwork('org1', 'Admin@org1.example.com');
            const contract = network.getContract(CONTRACT_NAME);

            // Invoke the QueryHealthRecords function
            const result = await contract.evaluateTransaction('QueryHealthRecords', id);
            const healthRecord = JSON.parse(result.toString());

            res.status(200).json(healthRecord);
        } catch (error) {
            console.error('Error querying health records:', error);
            res.status(500).json({ error: error.message });
        }
    }
};
