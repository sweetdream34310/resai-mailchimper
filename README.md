# About

`anymail-api` runs a gRPC server. All times must be in UTC timezone (including server's time).

## Development

Make sure the following exist in your sourced profile:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
export GOPATH=$(go env GOPATH)
```

* Go 1.15+
* `Docker` and `docker-compose`
* protoc v3.13.0

## Build

```bash
# create binary
make build

# create local docker image
make docker-build
make local-up

# stop service
make local-down
```
## How to create example Gmail token in development

go to directory `./testing` and make sure you have `credentials.json` file there.
If you don't have that file please contact your devops to generate that credential for you.

please run using the following command from your working directory:

```bash
go run gmail.go
```
The first time you run the sample, it prompts you to authorize access:

1. Copy the url that you get from your terminal / command-line prompt
2. Browse to the provided URL in your web browser.
3. If you're not already signed in to your Google account, you're prompted to sign in. If you're signed in to multiple Google accounts, you are asked to select one account to use for authorization.
4. Click the Accept button.
   Copy the code you're given, paste it into the command-line prompt, and press Enter.
5. If you get page not found with the following url `https://www.googleapis.com/oauth2/v3/token?state=state-token&code=4/0AX4XfWj4BtIUa7TbvVv2vWQAc9ykmE7V-JQf-Xba1ikAQS5yHBno5syA68xbWabRrTbvRg&scope=https://www.googleapis.com/auth/gmail.readonly` .
   Copy the value `code` from that url, paste it into the command-line prompt, and press Enter.
6. When the process is complete, it will generate `token.json` file, you can use access_token from that file to call `awaymail` API `POST /v2/session` to get service `token`.
7. Use service `token` to call all `awaymail` API

reference: https://developers.google.com/gmail/api/quickstart/go
## Protos

Contracts can be found in `proto/` directory and can be compiled with:

```bash
make protos
```