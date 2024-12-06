 ./network.sh deployCC -ccn registration -ccp ./chaincode/insurance-registration/ -ccl go -ccep "OR('Org1MSP.peer','Org2MSP.peer')" -cccg ./chaincode/insurance-registration/collections_config.json


  ./network.sh deployCC -ccn claims -ccp ./chaincode/insurance-claims-processing -ccl go