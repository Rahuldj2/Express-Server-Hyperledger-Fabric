# PART 2 DEFINE POLICY
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C mychannel -n registration $PEER_CONN_PARMS -c '{"function":"DefinePolicy","Args":["policy123","HealthInsurance","100000.0","500.0","2024-01-01","2025-01-01","{\"IsNonSmoker\": true, \"HasDisease\": false}","[\"Cancer\", \"Diabetes\"]"]}'

#QUERY POLICY
peer chaincode query -C mychannel -n registration -c '{"function":"QueryPolicy","Args":["policy123"]}'


#QUERY ALL POLICIES
peer chaincode query -C mychannel -n registration --peerAddresses localhost:9051 --tlsRootCertFiles /home/rahul/hyperledger-fabric/fabric-samples/test-network/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"Args":["QueryAllPolicies"]}'
