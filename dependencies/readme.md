# Runttime dependencies

## TLS Certificates
In order to use SSL/TLS with any listener or server used by this project you will need a certificate. The instructions bellow are just to make it easy to understand how it works, you might use any certificate at hand or at any location.

```bash
# make sure the cert folder exists
$ mkdir dependencies/cert

# answer the questions with some fake data
$ openssl req -new \
    -newkey rsa:4096 \
    -x509 \
    -sha256 \
    -days 365 \
    -nodes \
    -keyout dependencies/cert/api1.key \
    -out dependencies/cert/api1.crt
...
Country Name (2 letter code) [AU]:PY
State or Province Name (full name) [Some-State]:Central
Locality Name (eg, city) []:San Loranzo
Organization Name (eg, company) [Internet Widgits Pty Ltd]:ACME Ltd
Organizational Unit Name (eg, section) []:
Common Name (e.g. server FQDN or YOUR name) []:
Email Address []:
...
```

## MQ server
For this example we are using RabbitMQ. `.env` file contains some settings, update as needed. We are using the amqp 1.0 plugin by mounting `enable_plugins` at the `/etc/rabbitmq/enabled_plugins` folder inside the container, following the documentation.

```bash
$ cd dependencies/rabbitmq
$ docker compose --env-file=.env up -d
```


