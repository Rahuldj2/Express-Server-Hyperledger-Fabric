  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C mychannel -n private $PEER_CONN_PARMS -c '{"Args":["DefinePolicy","Policy123","Health","100000","1000","2024-01-01","2024-12-31", "{\"NonSmoker\":true, \"NoPreExistingDisease\":true}"]}' --waitForEvent
  peer chaincode query -C mychannel -n private -c '{"Args":["QueryPolicy","Policy123"]}'
  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C mychannel -n private $PEER_CONN_PARMS -c '{"Args":["RegisterForPolicy","User123","Policy123","1000","true","false"]}' --waitForEvent
  peer chaincode query -C mychannel -n private -c '{"Args":["QueryRegistration","User123","Policy123"]}'
  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C mychannel -n private $PEER_CONN_PARMS -c '{"Args":["UploadHealthRecords","Patient is in good health"]}' --waitForEvent
  peer chaincode query -C mychannel -n private -c '{"Args":["QueryHealthRecords","1732712031"]}'


#TIMESTAMP COMMANDS
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C mychannel -n registration $PEER_CONN_PARMS -c '{"Args":["UploadHealthRecords","Record123","Patient is in good health"]}' --waitForEvent



  peer chaincode query -C mychannel -n registration --peerAddresses localhost:7051 --tlsRootCertFiles /home/rahul/hyperledger-fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt -c '{"Args":["GetPrivateData","123"]}'
  peer chaincode query -C mychannel -n registration --peerAddresses localhost:9051 --tlsRootCertFiles /home/rahul/hyperledger-fabric/fabric-samples/test-network/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"Args":["GetPrivateData","123"]}'

  peer chaincode query -C mychannel -n registration --peerAddresses localhost:9051 --tlsRootCertFiles /home/rahul/hyperledger-fabric/fabric-samples/test-network/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"Args":["QueryHealthRecords","Record123"]}'