peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C mychannel -n registration $PEER_CONN_PARMS -c '{"function":"RegisterForPolicy","Args":["user123","policy123","1000.0","true","false","true"]}'