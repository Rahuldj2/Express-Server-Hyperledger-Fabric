peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C mychannel -n claims $PEER_CONN_PARMS -c '{"function":"UploadPatientDetails","Args":["user123","Cancer","Chemotherapy","Hospital A","2024-01-01","2024-01-15"]}'
