# Go SAML / OIDC Single Sign On examples

## About

* The repo demonstrates SAML and OIDC Single Singon with Keycloak as a SAML IdP
* It uses small Go programs as SAML/OIDC SP

## Keycloak

### How to configure

* Setup Keycloak on Linux (change timemachine to your host name)

```bash
openssl genrsa 2048 > ca.key
openssl req -x509 -new -nodes -key ca.key -subj "/CN=rootca" -days 10000 -out ca.crt
# create certs for keycloak
openssl genrsa 2048 > server.key
openssl req -new -key server.key -subj "/CN=timemachine" > server.csr
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -days 10000 -out server.crt
# for Keycloak
mkdir certs
cp server.key certs/tls.key
cp server.crt certs/tls.crt
# for Grafana
cp ca.crt grafana_oidc/ca.crt
cp server.key grafana_oidc/server.key
cp server.crt grafana_oidc/server.crt
cp server.key grafana_oidc/client.key
cp server.crt grafana_oidc/client.crt
```

### How to run

* Run Keycloak
* It didn't start successfully in Docker Preview on Apple Silicon Mac. I used Ubuntu 20.04 (timemachine) on Windows Hyper-V

```bash
docker-compose up
```

## SAML

### Keycloak config

* Configure Keycloak. The default user/pass is admin/admin
  * _Clients_ -> _Create_
  * ClientID = 'http://scottmm.local:8000/saml/metadata' (change `scottmm.local` to your host name)
  * _Client Protocol_ -> saml
  * Settings tab
    * Client Signature Required: OFF
    * _Root URL_ -> `http://scottmm.local:8000`
    * _Valid Redirect URIs_ -> `http://scottmm.local:8000/*`
  * Mappers -> Crate User Property
    * Name: username
	* Friendly Name: username
	* SAML Attribute Name: username
	* SAML Attribute Name Format: Basic
  * SAML Keys
    * Copy Private key and Certificate into go_saml/service.key and go_saml/service.crt by adding header/footer as below
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

* Create 'scott' user in Users menu with a password with _Temporary_ false

### How to run SAML SP

* Change `go_saml.go`
  * Change `scottmm.local` to your host name
  * Change `timemachine:8443` to your keycloak host name
* Run SAML SP - This worked on Apple Silicon Mac, too

```bash
cd go_saml
export GODEBUG="x509ignoreCN=0"
go run go_saml.go
```

* Access URL http://$yourhost:$port/hello in a browser
* It'll be redirected to Keycloak
* Logon as 'scott'
![Keycloak](docs/keycloak.png)
* You'll get 'Hello, scott!'

## OIDC

### Keycloak config

* Configure Keycloak. The default user/pass is admin/admin
  * _Clients_ -> _Create_
  * ClientID = 'http://scottmm.local:8000' (change `scottmm.local` to your host name)
  * _Client Protocol_ -> openid-connect
	* _Root URL_ -> `http://scottmm.local:8000`
  * Settings tab
    * Client Signature Required: OFF
    * _Root URL_ -> `http://scottmm.local:8000`
    * _Valid Redirect URIs_ -> `http://scottmm.local:8000/*`
    * _Access Type_ -> confidential
  * Credentials tab
    * Copy secret into go_oidc.go -> secret in getConfig()
* Create 'scott' user in Users menu with a password with _Temporary_ false

### How to run OIDC SP - Go

* Change `go_oidc.go`
  * Change `scottmm.local` to your host name
  * Change `timemachine:8443` to your keycloak host name
* Run OIDC SP - This worked on Apple Silicon Mac, too

```bash
cd go_oidc
export GODEBUG="x509ignoreCN=0"
go run go_oidc.go
```

* Access URL http://$yourhost:$port/hello in a browser
* It'll be redirected to Keycloak
* Logon as 'scott'
* You'll get something like below

```text
{
    "OAuth2Token": {
        "access_token": "eyJhbGciOiJSUz...",
        "token_type": "Bearer",
        "refresh_token": "eyJhbGciOiJIUzI1N...",
        "expiry": "2021-02-13T11:54:15.344809+09:00"
    },
    "UserInfo": {
        "sub": "6bcf488d-...",
        "profile": "",
        "email": "",
        "email_verified": false
    }
}
```

### How to run OIDC SP - Grafana

* Configure Keycloak similarly
* Change `grafana_oidc/grafana.ini`
  * Change `timemachine` to your Docker host which runs grafana
  * Change `timemachine:8443` to your keycloak host name
  * Change `allowed_domains` to your Keycloak users' email domain
  * Change `client_id` and `client_secret`
* Run grafana
```bash
docker run -v /home/scott/repo/SAML_OIDC/grafana_oidc:/var/tmp/grafana -e GF_PATHS_CONFIG=/var/tmp/grafana/grafana.ini -e GF_SERVER_ROOT_URL=https://timemachine:3000/ -p 3000:3000 --network host grafana/grafana
```
* Visit https://timemachine:3000 -> "Sing in with OAuth"
![Login](docs/grafana_login.png)
