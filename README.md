# govenv
A simple and dirty virtual environment maker for golang.

Work with different versions of golang on the same computer.

Environment are fully isolated (virtual environment).

Each environment has its own packages

## Install
```shell
go get github.com/4383/govenv
go install github.com/4383/govenv/govenv
```

## Usages
```shell
$GOPATH/bin/govenv -dest <your-env-name> # virtualenv initialized with the latest version of golang
$GOPATH/bin/govenv -dest <your-env-name> -go-version 1.6 # virtualenv initialized with golang version 1.6 
```

## Warning
This a development version with few bugs inside and functionalities are not fully implemented
