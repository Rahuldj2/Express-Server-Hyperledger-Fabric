

peer chaincode invoke -o localhost:7050 \
--ordererTLSHostnameOverride orderer.example.com \
--tls \
--cafile $ORDERER_CA \
-C mychannel \
-n registration \
$PEER_CONN_PARMS \
-c '{"Args":["DefinePolicy","Policy123","Health","100000","5000","2024-01-01","2024-12-31","{\"isNonSmoker\":true}","[\"Cancer\",\"Diabetes\"]"]}' \
--waitForEvent




# PART 2 DEFINE POLICY
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C mychannel -n registration $PEER_CONN_PARMS -c '{"function":"DefinePolicy","Args":["policy123","HealthInsurance","100000.0","500.0","2024-01-01","2025-01-01","{\"IsNonSmoker\": true, \"HasDisease\": false}","[\"Cancer\", \"Diabetes\"]"]}'

#QUERY POLICY
peer chaincode query -C mychannel -n registration -c '{"function":"QueryPolicy","Args":["policy123"]}'




peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C mychannel -n claims $PEER_CONN_PARMS -c '{"function":"ProcessClaim","Args":["12345", "Diabetes", "5000"]}'


