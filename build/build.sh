export DOCKER_BUILDKIT=1
docker build -f ./build/Dockerfile -t registry.rijksapps.nl/cno/standaardplatform/basis/keycloak-bridge:local .
