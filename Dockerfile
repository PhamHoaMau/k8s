#Docker file for setting up the org1 peer
FROM hyperledger/fabric-orderer:2.2

LABEL  maintainer="Pham Hoa Mau <mauphamhoa@gmail.com>"

#11. Install the jq package - used in scripts
RUN apt-get update \ 
 && apt-get install git -y fabric 

#13. Launch the peer
CMD orderer
