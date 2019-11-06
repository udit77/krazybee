## How to run the server

This repository should equate to:
```bash
$GOPATH/src/github.com/krazybee
```

Make sure to dep ensure
```bash
dep ensure -v
```

Edit searchservice.json file inside files directory for your db user and password. Copy this file into '/etc/krazybee/'.
```bash
$ sudo cp files/searchservice.json /etc/krazybee/
```

Once done you can build and run app server:
```
$ go build && ./krazybee  --search-service
```
