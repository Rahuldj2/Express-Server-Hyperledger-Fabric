peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C mychannel -n registration $PEER_CONN_PARMS -c '{"function":"UploadHealthRecords","Args":["user123","true","false"]}' --waitForEvent


peer chaincode query -C mychannel -n registration --peerAddresses localhost:9051 --tlsRootCertFiles /home/rahul/hyperledger-fabric/fabric-samples/test-network/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"Args":["QueryHealthRecords","user123"]}'
