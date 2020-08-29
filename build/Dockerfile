FROM golang:1.15-buster AS build

WORKDIR $GOPATH/src/keycloak

# Fetch dependencies
COPY go.mod go.sum ./
RUN go mod download -x

# Copy golang code and build
COPY . . 
RUN CGO_ENABLED=0 go build -v -o /bin/keycloak-bridge

# Copy to scratch
FROM scratch
COPY --from=build /bin/keycloak-bridge /bin/keycloak-bridge
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT ["/bin/keycloak-bridge"]
