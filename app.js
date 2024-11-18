// app.js
const express = require('express');
const bodyParser = require('body-parser');
const insuranceController = require('./insuranceController');


const app = express();
app.use(bodyParser.json());

const cors = require('cors');
app.use(cors({ origin: 'http://localhost:3000' }));



app.post('/registerPolicy', insuranceController.registerPolicy);
app.get('/queryPolicy/:policyID', insuranceController.queryPolicy);
app.post('/processClaim', insuranceController.processClaim);

app.get('/health', (req, res) => {
    res.status(200).json({ status: 'OK', message: 'Server is running' });
});


const PORT = 3001;
app.listen(PORT, () => {
    console.log(`Server running on port ${PORT}`);
});
