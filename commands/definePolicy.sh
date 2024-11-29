peer chaincode invoke -o localhost:7050 \
--ordererTLSHostnameOverride orderer.example.com \
--tls \
--cafile $ORDERER_CA \
-C mychannel \
-n registration \
$PEER_CONN_PARMS \
-c '{"Args":["DefinePolicy","Policy123","Health","100000","5000","2024-01-01","2024-12-31","{\"isNonSmoker\":true}","[\"Cancer\",\"Diabetes\"]"]}' \
--waitForEvent
