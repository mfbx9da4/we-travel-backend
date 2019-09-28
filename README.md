# we-travel-backend

Kiwi.com hackathon WeTravel API

## Example API

Install

```
go get -u github.com/gorilla/mux
go get -u github.com/cheekybits/genny/generic
```

To build

```
go build example-api.go
```

To run

```
go run example-api.go
```

Tutorial:
[API tutorial](https://medium.com/the-andela-way/build-a-restful-json-api-with-golang-85a83420c9da)

    curl -o london.osm https://api.openstreetmap.org/api/0.6/map\?bbox\=-0.489,51.680,0.486,51.686

    curl -XPOST -H "Content-type: application/json" -d '{"fromLocation": [-0.100573,51.543750],"toLocation": [-0.097091,51.545466]}' 'http://localhost:8080/findpath'
