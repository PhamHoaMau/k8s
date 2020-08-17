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
if [ -d "jwclab_orderer_org_crypto" ]; then
  rm -Rf jwclab_orderer_org_crypto/
fi

sleep 1
echo
echo "##########################################################"
echo "Generate certificates using Fabric CA's for Orderer"
echo "##########################################################"

#=================================STEP1=================================
echo
echo "---------> Step1: Enroll the CA admin"
echo

mkdir -p ${PWD}/jwclab_orderer_org_crypto
export FABRIC_CA_CLIENT_HOME=${PWD}/jwclab_orderer_org_crypto/

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
  OrganizationalUnitIdentifier: orderer' > ${PWD}/jwclab_orderer_org_crypto/msp/config.yaml

#=================================STEP2=================================
echo
echo "---------> Step2: Register orderer"
echo
set -x
fabric-ca-client register --id.name orderer --id.secret ordererpw --id.type orderer
set +x

#=================================STEP3=================================
echo
echo "---------> Step3: Register the orderer admin"
echo
set -x
fabric-ca-client register --id.name ordererAdmin --id.secret ordererAdminpw --id.type admin
set +x

#=================================STEP4=================================
mkdir -p jwclab_orderer_org_crypto/orderers
mkdir -p jwclab_orderer_org_crypto/orderers/jwclab.com
mkdir -p jwclab_orderer_org_crypto/orderers/orderer.jwclab.com

echo
echo "---------> Step4: Generate the orderer msp"
echo
set -x
fabric-ca-client enroll -u http://orderer:ordererpw@128.199.151.251:32769 -M ${PWD}/jwclab_orderer_org_crypto/orderers/orderer.jwclab.com/msp --csr.hosts orderer.jwclab.com --csr.hosts localhost
set +x

cp ${PWD}/jwclab_orderer_org_crypto/msp/config.yaml ${PWD}/jwclab_orderer_org_crypto/orderers/orderer.jwclab.com/msp/config.yaml

#=================================STEP5=================================
echo
echo "---------> Step5: Generate the orderer-tls certificates"
echo
set -x
fabric-ca-client enroll -u http://orderer:ordererpw@128.199.151.251:32769 -M ${PWD}/jwclab_orderer_org_crypto/orderers/orderer.jwclab.com/tls --enrollment.profile tls --csr.hosts orderer.jwclab.com --csr.hosts localhost
set +x

cp ${PWD}/jwclab_orderer_org_crypto/orderers/orderer.jwclab.com/tls/tlscacerts/* ${PWD}/jwclab_orderer_org_crypto/orderers/orderer.jwclab.com/tls/ca.crt
cp ${PWD}/jwclab_orderer_org_crypto/orderers/orderer.jwclab.com/tls/signcerts/* ${PWD}/jwclab_orderer_org_crypto/orderers/orderer.jwclab.com/tls/server.crt
cp ${PWD}/jwclab_orderer_org_crypto/orderers/orderer.jwclab.com/tls/keystore/* ${PWD}/jwclab_orderer_org_crypto/orderers/orderer.jwclab.com/tls/server.key

mkdir -p ${PWD}/jwclab_orderer_org_crypto/orderers/orderer.jwclab.com/msp/tlscacerts
cp ${PWD}/jwclab_orderer_org_crypto/orderers/orderer.jwclab.com/tls/tlscacerts/* ${PWD}/jwclab_orderer_org_crypto/orderers/orderer.jwclab.com/msp/tlscacerts/tlsca.jwclab.com-cert.pem

mkdir -p ${PWD}/jwclab_orderer_org_crypto/msp/tlscacerts
cp ${PWD}/jwclab_orderer_org_crypto/orderers/orderer.jwclab.com/tls/tlscacerts/* ${PWD}/jwclab_orderer_org_crypto/msp/tlscacerts/tlsca.jwclab.com-cert.pem

#=================================STEP6=================================
mkdir -p jwclab_orderer_org_crypto/users
mkdir -p jwclab_orderer_org_crypto/users/admin@jwclab.com

echo
echo "---------> Step6: Generate the admin msp"
echo
set -x
fabric-ca-client enroll -u http://ordererAdmin:ordererAdminpw@128.199.151.251:32769 -M ${PWD}/jwclab_orderer_org_crypto/users/admin@jwclab.com/msp
set +x

cp ${PWD}/jwclab_orderer_org_crypto/msp/config.yaml ${PWD}/jwclab_orderer_org_crypto/users/admin@jwclab.com/msp/config.yaml
