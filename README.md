# Go SAML / OIDC Single Sign On examples

## About

* The repo demonstrates SAML and OIDC Single Singon with Keycloak as a SAML IdP
* It uses small Go programs as SAML SP

## SAML
### How to configure

* Setup Keycloak (change timemachine to your host name)

```bash
openssl genrsa 2048 > ca.key
openssl req -x509 -new -nodes -key ca.key -subj "/CN=rootca" -days 10000 -out ca.crt
# create certs for keycloak
openssl req -new -key server.key -subj "/CN=servername" > server.csr
openssl genrsa 2048 > server.key
openssl req -new -key server.key -subj "/CN=timemachine" > server.csr
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -days 10000 -out server.crt
mkdir certs
cp server.key certs/tls.key
cp server.crt certs/tls.crt
```

### How to run

* Run Keycloak

```bash
docker-compose up
```

* Configure Keycloak. The default user/pass is admin/admin
  * _Clients_ -> _Create_
  * ClientID = 'http://scottmm.local:8000/saml/metadata' (change `scottmm.local` to your host name)
  * _Client Protocol_ -> saml
  * Settings tab
    * Client Signature Required: OFF
	* Root URL -> `http://scottmm.local:8000`
	* Valid Redirect URIs -> `http://scottmm.local:8000/*`
  * Mappers -> Crate User Property
    * Name: username
	* Friendly Name: username
	* SAML Attribute Name: username
	* SAML Attribute Name Format: Basic
  * SAML Keys
    * Copy Private key and Certificate into service.key and service.crt by adding header/footer as below
    * service.key

    ```sh
    -----BEGIN RSA PRIVATE KEY-----
    MIIEow...
    -----END RSA PRIVATE KEY-----
    ```

    * service.crt

    ```sh
    -----BEGIN CERTIFICATE-----
    MIICn...
    -----END CERTIFICATE-----
    ```

* Create 'scott' user with a password
* Change `go_saml.go`
  * Change `scottmm.local` to your host name
  * Change `timemachine:8443` to your keycloak host name
* Run SAML SP
```bash
export GODEBUG="x509ignoreCN=0"
go run main.go
```

* Access go_saml URL http://$yourhost:$port/hello by a browser
* It'll be redirected to Keycloak
* Logon as 'scott'
![Keycloak](docs/keycloak.png)
* You'll get 'Hello, scott!'

## OIDC

* Similart to SAML (TBD)