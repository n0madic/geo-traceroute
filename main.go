package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/kanocz/tracelib"
	"github.com/oschwald/geoip2-golang"
)

var (
	args struct {
		ASN     string   `arg:"-a" placeholder:"PATH" default:"GeoLite2-ASN.mmdb" help:"MaxMind ASN database path (optional)"`
		GeoDB   []string `arg:"-g,--geo" placeholder:"PATH" help:"MaxMind GeoIP2 database path"`
		Lang    string   `arg:"-l" default:"en" help:"Country/City locale"`
		MaxRTT  int64    `arg:"-r" default:"5" placeholder:"RTT" help:"Maximum RTT"`
		MaxTTL  int      `arg:"-m" default:"30" placeholder:"TTL" help:"Maximal TTL (hops)"`
		Path    string   `arg:"-p" help:"Path prefix to MaxMind databases"`
		Source  string   `arg:"-s" default:"0.0.0.0" placeholder:"IP" help:"IPv4 source"`
		Source6 string   `arg:"-6" default:"::" placeholder:"IP" help:"IPv6 source"`
		Target  string   `arg:"positional,required" help:"Target host"`
	}

	geoIP2Filenames = []string{
		"GeoIP2-City.mmdb",
		"GeoLite2-City.mmdb",
		"GeoIP2-Country.mmdb",
		"GeoLite2-Country.mmdb",
	}
)

type geoDB struct {
	GeoIP2 *geoip2.Reader
	ASN    *geoip2.Reader
}

func printStep(hop tracelib.Hop, num int, round int) {
	if hop.Error == nil {
		fmt.Print("!")
	} else {
		fmt.Print("*")
	}
}

func main() {
	var err error
	var db geoDB

	arg.MustParse(&args)

	// If the path is not specified, then search in standard directories
	if args.Path == "" && len(args.GeoDB) == 0 {
		searchPath := []string{
			".",
			"/usr/share/GeoIP/",
			"/usr/local/share/GeoIP/",
			"/var/lib/GeoIP/",
		}
		for _, path := range searchPath {
			matches, err := filepath.Glob(filepath.Join(path, "*.mmdb"))
			if len(matches) > 0 && err == nil {
				args.Path = path
				break
			}
		}
	}

	// Checking the slash at the end of the path
	if args.Path != "" && !strings.HasSuffix(args.Path, string(os.PathSeparator)) {
		args.Path += string(os.PathSeparator)
	}

	// Default list of database file names
	if len(args.GeoDB) == 0 {
		args.GeoDB = geoIP2Filenames
	}

	// Trying to open the specified database file names
	for _, file := range args.GeoDB {
		db.GeoIP2, err = geoip2.Open(args.Path + file)
		if err == nil {
			break
		}
	}
	if db.GeoIP2 == nil {
		fmt.Println("ERROR: GeoIP2 database not found!")
		os.Exit(1)
	}
	defer db.GeoIP2.Close()

	// Trying to open an optional ASN base file
	db.ASN, err = geoip2.Open(args.Path + args.ASN)
	if err == nil {
		defer db.ASN.Close()
	}

	// Run traceroute
	cache := tracelib.NewLookupCache()
	hops, err := tracelib.RunTrace(
		args.Target,
		args.Source,
		args.Source6,
		time.Second*time.Duration(args.MaxRTT),
		args.MaxTTL,
		cache,
		printStep)
	if err != nil {
		fmt.Println("Traceroute error:", err)
		os.Exit(2)
	}

	// Clear printing of trace steps
	clearLine := "\033[2K\r"
	if runtime.GOOS == "windows" {
		clearLine = "\r"
	}
	fmt.Print(clearLine)

	// Create a table of results
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	header := table.Row{"#", "IP", "Host", "RTT", "Country/City", "ASN"}
	if db.ASN != nil {
		header = append(header, "ASN Org")
	}
	t.AppendHeader(header)

	for i, hop := range hops {
		if hop.Error != nil {
			ip := "*"
			if hop.Addr != nil {
				ip = hop.Addr.String()
			}
			errmsg := ""
			if hop.Timeout {
				errmsg = "<timeout>"
			} else {
				errmsg = hop.Error.Error()
			}
			t.AppendRow([]interface{}{i + 1, ip, errmsg})
			continue
		}

		ip := net.ParseIP(hop.Addr.String())
		// Lookup GeoIP
		record, err := db.GeoIP2.City(ip)
		if err != nil {
			fmt.Println(err)
		}

		location := record.Country.Names[args.Lang]
		if record.City.Names[args.Lang] != "" {
			location = location + "/" + record.City.Names[args.Lang]
		}

		// Lookup ASN
		var asn *geoip2.ASN
		if db.ASN != nil {
			asn, err = db.ASN.ASN(ip)
			if err != nil {
				fmt.Println(err)
			}
			if asn.AutonomousSystemNumber == 0 && hop.AS > 0 {
				asn.AutonomousSystemNumber = uint(hop.AS)
			}
		}

		row := []interface{}{
			i + 1,
			hop.Addr,
			strings.TrimSuffix(hop.Host, "."),
			hop.RTT.Round(time.Microsecond),
			location,
		}
		if db.ASN != nil && asn.AutonomousSystemNumber > 0 {
			row = append(row, asn.AutonomousSystemNumber, asn.AutonomousSystemOrganization)
		} else if hop.AS > 0 {
			row = append(row, hop.AS)
		}

		t.AppendRow(row)
	}

	t.Render()
}
