peer chaincode query -C mychannel -n registration --peerAddresses localhost:9051 --tlsRootCertFiles /home/rahul/hyperledger-fabric/fabric-samples/test-network/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"Args":["QueryPolicyByUserID","user123"]}'
