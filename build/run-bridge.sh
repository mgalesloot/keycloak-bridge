#!/bin/bash

TOOLBOX_IMG=registry.rijksapps.nl/cno/standaardplatform/basis/keycloak-bridge:local

docker run --rm -it \
    -v $KUBECONFIG:/.kube/config \
    -v $PWD/examples:/examples \
    $TOOLBOX_IMG 

