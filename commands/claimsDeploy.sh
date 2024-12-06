./network.sh deployCC -ccn claims -ccp ./chaincode/insurance-claims-processing -ccl go


 ./network.sh deployCC -ccn claims -ccp ./chaincode/insurance-claims-processing/ -ccl go -ccep "OR('Org1MSP.peer','Org2MSP.peer')" -cccg ./chaincode/insurance-claims-processing/collections_config.json
