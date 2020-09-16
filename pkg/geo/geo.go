package geo

import (
	"log"
	"net"
	"strings"

	"github.com/gobuffalo/packr/v2"
	"github.com/oschwald/geoip2-golang"
)

// DB handle geo lookup reader.
type DB struct {
	reader *geoip2.Reader
}

func New() *DB {
	box := packr.New("Geo.Data", "./data")

	data, err := box.Find("GeoLite2-Country.mmdb")
	if err != nil {
		panic(err)
	}

	db, err := geoip2.FromBytes(data)
	if err != nil {
		log.Fatal(err)
	}

	instance := &DB{
		reader: db,
	}

	return instance
}

func (db *DB) Close() {
	db.Close()
}

// Resolve should
func (db DB) Resolve(ip string) (string, error) {
	record, err := db.reader.Country(net.ParseIP(ip))
	if err != nil {
		return "", err
	}

	return strings.ToLower(record.Country.IsoCode), nil
}
