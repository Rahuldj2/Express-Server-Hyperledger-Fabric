// app.js

const express = require('express');
const bodyParser = require('body-parser');
const cors = require('cors');

const insuranceController = require('./controllers/insuranceController');
const claimsController = require('./controllers/claimsController');

const app = express();

// Middleware
app.use(bodyParser.json());
app.use(cors({ origin: 'http://localhost:3000' })); // Adjust the origin as per your frontend

// Health Check Endpoint
app.get('/health', (req, res) => {
    res.status(200).json({ status: 'OK', message: 'Server is running' });
});

// Insurance (Registration) Routes
app.post('/insurance/definePolicy', insuranceController.definePolicy);//tested
app.get('/insurance/queryPolicy/:policyID', insuranceController.queryPolicy);//tested
app.post('/insurance/registerForPolicy', insuranceController.registerForPolicy);//tested
app.get('/insurance/queryRegistration/:userId/:policyId', insuranceController.queryRegistration);//tested
app.post('/insurance/uploadHealthRecords', insuranceController.uploadHealthRecords);//tested
app.get('/insurance/queryHealthRecords/:id', insuranceController.queryHealthRecords);//tested
app.get('/insurance/queryHealthRecordsorg1/:id', insuranceController.queryHealthRecordsOrg1);//tested
app.get('/insurance/queryAllPolicies', insuranceController.queryAllPolicies);

// Claims Routes
app.post('/claims/uploadPatientDetails', claimsController.uploadPatientDetails);//tested
app.post('/claims/processClaim', claimsController.processClaim);
app.get('/claims/queryClaim/:userID', claimsController.queryClaim);

// Start the server
const PORT = 3001;
app.listen(PORT, '0.0.0.0', () => {
    console.log(`Server running on port ${PORT}`);
});
