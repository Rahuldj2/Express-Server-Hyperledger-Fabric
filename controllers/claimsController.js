// controllers/claimsController.js

const connectToNetwork = require('../network/connect');

const REGISTRATION_CONTRACT = 'registration'; // Name of the Registration chaincode
const CLAIMS_CONTRACT = 'claims'; // Name of the Claims chaincode

module.exports = {
    uploadPatientDetails: async (req, res) => {
        const { userID, diseaseDiagnosis, treatmentPlan, hospitalName, admissionDate, dischargeDate } = req.body;
        try {
            const network = await connectToNetwork('org1', 'Admin@org1.example.com');
            const contract = network.getContract(CLAIMS_CONTRACT);

            await contract.submitTransaction(
                'UploadPatientDetails',
                userID,
                diseaseDiagnosis,
                treatmentPlan,
                hospitalName,
                admissionDate,
                dischargeDate
            );

            res.status(200).send(`Patient details for ${userID} uploaded successfully.`);
        } catch (error) {
            console.error('Error uploading patient details:', error);
            res.status(500).json({ error: error.message });
        }
    },

    processClaim: async (req, res) => {
        const { userID } = req.body;
        try {
            const network = await connectToNetwork('org2', 'Admin@org2.example.com');
            const contract = network.getContract(CLAIMS_CONTRACT);

            await contract.submitTransaction('ProcessClaim', userID);

            res.status(200).send(`Claim for user ${userID} processed successfully.`);
        } catch (error) {
            console.error('Error processing claim:', error);
            res.status(500).json({ error: error.message });
        }
    },

    queryClaim: async (req, res) => {
        const { userID } = req.params;
        try {
            const network = await connectToNetwork('org2', 'Admin@org2.example.com');
            const contract = network.getContract(CLAIMS_CONTRACT);

            const result = await contract.evaluateTransaction('QueryClaim', userID);
            res.status(200).json(JSON.parse(result.toString()));
        } catch (error) {
            console.error('Error querying claim:', error);
            res.status(500).json({ error: error.message });
        }
    }
};
