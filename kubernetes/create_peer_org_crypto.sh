#!/bin/bash

# add binary file directory
export PATH=${PWD}/../bin:$PATH

# Check if your have the fabric-ca-client binaries
fabric-ca-client version > /dev/null 2>&1
if [ $? -ne 0 ]; then
  echo "ERROR! fabric-ca-client binary not found.."
  exit 1
fi

# clean previous generation
if [ -d "$1_peer_org_crypto" ]; then
  rm -Rf $1_peer_org_crypto/
fi

sleep 1
echo
echo "##########################################################"
echo "Organization name: $1"
echo "Generate certificates using Fabric CA's for $1"
echo "##########################################################"

#=================================STEP1=================================
echo
echo "---------> Step1: Enroll the CA admin"
echo

mkdir -p ${PWD}/$1_peer_org_crypto
export FABRIC_CA_CLIENT_HOME=${PWD}/$1_peer_org_crypto/

set -x
fabric-ca-client enroll -u http://admin:adminpw@128.199.151.251:32769
set +x

echo 'NodeOUs:
Enable: true
ClientOUIdentifier:
  Certificate: cacerts/128-199-151-251-32769.pem
  OrganizationalUnitIdentifier: client
PeerOUIdentifier:
  Certificate: cacerts/128-199-151-251-32769.pem
  OrganizationalUnitIdentifier: peer
AdminOUIdentifier:
  Certificate: cacerts/128-199-151-251-32769.pem
  OrganizationalUnitIdentifier: admin
OrdererOUIdentifier:
  Certificate: cacerts/128-199-151-251-32769.pem
  OrganizationalUnitIdentifier: orderer' > ${PWD}/$1_peer_org_crypto/msp/config.yaml

#=================================STEP2=================================
echo
echo "---------> Step2: Register peer0 of $1 Organization"
echo
set -x
fabric-ca-client register --id.name $1_peer0 --id.secret peer0pw --id.type peer
set +x

#=================================STEP3=================================
echo
echo "---------> Step3: Register user1 of $1 Organization"
echo
set -x
fabric-ca-client register --id.name $1_user1 --id.secret user1pw --id.type client
set +x

#=================================STEP4=================================
echo
echo "---------> Step4: Register Admin of $1 Organization"
echo
set -x
fabric-ca-client register --id.name $1_admin --id.secret $1adminpw --id.type admin
set +x

#=================================STEP5=================================
mkdir -p $1_peer_org_crypto/peers
mkdir -p $1_peer_org_crypto/peers/peer0.$1.jwclab.com

echo
echo "---------> Step5: Generate the peer0 msp of $1 Organization"
echo
set -x
fabric-ca-client enroll -u http://$1_peer0:peer0pw@128.199.151.251:32769 -M ${PWD}/$1_peer_org_crypto/peers/peer0.$1.jwclab.com/msp --csr.hosts peer0.$1.jwclab.com
set +x

cp ${PWD}/$1_peer_org_crypto/msp/config.yaml ${PWD}/$1_peer_org_crypto/peers/peer0.$1.jwclab.com/msp/config.yaml

#=================================STEP6=================================
echo
echo "---------> Step6: Generate the peer0-tls certificates of $1 Organization"
echo
set -x
fabric-ca-client enroll -u http://$1_peer0:peer0pw@128.199.151.251:32769 -M ${PWD}/$1_peer_org_crypto/peers/peer0.$1.jwclab.com/tls --enrollment.profile tls --csr.hosts peer0.$1.jwclab.com --csr.hosts localhost
set +x

cp ${PWD}/$1_peer_org_crypto/peers/peer0.$1.jwclab.com/tls/tlscacerts/* ${PWD}/$1_peer_org_crypto/peers/peer0.$1.jwclab.com/tls/ca.crt
cp ${PWD}/$1_peer_org_crypto/peers/peer0.$1.jwclab.com/tls/signcerts/* ${PWD}/$1_peer_org_crypto/peers/peer0.$1.jwclab.com/tls/server.crt
cp ${PWD}/$1_peer_org_crypto/peers/peer0.$1.jwclab.com/tls/keystore/* ${PWD}/$1_peer_org_crypto/peers/peer0.$1.jwclab.com/tls/server.key

mkdir -p ${PWD}/$1_peer_org_crypto/msp/tlscacerts
cp ${PWD}/$1_peer_org_crypto/peers/peer0.$1.jwclab.com/tls/tlscacerts/* ${PWD}/$1_peer_org_crypto/msp/tlscacerts/ca.crt

mkdir -p ${PWD}/$1_peer_org_crypto/tlsca
cp ${PWD}/$1_peer_org_crypto/peers/peer0.$1.jwclab.com/tls/tlscacerts/* ${PWD}/$1_peer_org_crypto/tlsca/tlsca.$1.jwclab.com-cert.pem

mkdir -p ${PWD}/$1_peer_org_crypto/ca
cp ${PWD}/$1_peer_org_crypto/peers/peer0.$1.jwclab.com/msp/cacerts/* ${PWD}/$1_peer_org_crypto/ca/ca.$1.jwclab.com-cert.pem

#=================================STEP7=================================
mkdir -p $1_peer_org_crypto/users
mkdir -p $1_peer_org_crypto/users/user1@$1.jwclab.com

echo
echo "---------> Step7: Generate the user1 msp of $1 Organization"
echo
set -x
fabric-ca-client enroll -u http://$1_user1:user1pw@128.199.151.251:32769 -M ${PWD}/$1_peer_org_crypto/users/user1@$1.jwclab.com/msp
set +x

cp ${PWD}/$1_peer_org_crypto/msp/config.yaml ${PWD}/$1_peer_org_crypto/users/user1@$1.jwclab.com/msp/config.yaml

#=================================STEP8=================================
mkdir -p $1_peer_org_crypto/users/admin@$1.jwclab.com

echo
echo "---------> Step8: Generate the org admin msp of $1 Organization"
echo
set -x
fabric-ca-client enroll -u http://$1_admin:$1adminpw@128.199.151.251:32769 -M ${PWD}/$1_peer_org_crypto/users/admin@$1.jwclab.com/msp
set +x

cp ${PWD}/$1_peer_org_crypto/msp/config.yaml ${PWD}/$1_peer_org_crypto/users/admin@$1.jwclab.com/msp/config.yaml