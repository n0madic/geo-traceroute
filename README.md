# geo-traceroute

Traceroute with GeoIP lookup. Required MaxMind MMDB databases.

## Install

Download binaries from [release](https://github.com/n0madic/geo-traceroute/releases) page.

Or install from source:

```
go get -u github.com/n0madic/geo-traceroute
```

Requires *root* permission! Or set capabilities for binary:

```
sudo setcap cap_net_raw+ep geo-traceroute
```

## Help

```
Usage: geo-traceroute [--asn PATH] [--geo PATH] [--lang LANG] [--maxrtt RTT] [--maxttl TTL] [--path PATH] [--source IP] [--source6 IP] TARGET

Positional arguments:
  TARGET                 Target host

Options:
  --asn PATH, -a PATH    MaxMind ASN database path (optional) [default: GeoLite2-ASN.mmdb]
  --geo PATH, -g PATH    MaxMind GeoIP2 database path
  --lang LANG, -l LANG   Country/City locale [default: en]
  --maxrtt RTT, -r RTT   Maximum RTT [default: 5]
  --maxttl TTL, -m TTL   Maximal TTL (hops) [default: 30]
  --path PATH, -p PATH   Path prefix to MaxMind databases
  --source IP, -s IP     IPv4 source [default: 0.0.0.0]
  --source6 IP, -6 IP    IPv6 source [default: ::]
  --help, -h             display this help and exit
```

## Usage

```
$ geo-traceroute google.com
+----+-----------------+----------------------------+----------+-----------------------------+-------+
|  # | IP              | HOST                       |      RTT | COUNTRY/CITY                |   ASN |
+----+-----------------+----------------------------+----------+-----------------------------+-------+
|  1 | 192.168.31.1    | XiaoQiang                  |   1.25ms |                             |       |
|  2 | 88.81.236.150   |                            |  4.323ms | Ukraine/Heronymivka         |  3326 |
|  3 | 88.81.236.149   |                            |  5.024ms | Ukraine/Heronymivka         |  3326 |
|  4 | 88.81.244.161   |                            |  4.193ms | Ukraine/Kyiv                |  3326 |
|  5 | 88.81.240.182   | topnet-gw.google.com       |  5.372ms | Ukraine/Kyiv                |  3326 |
|  6 | 108.170.248.146 |                            | 14.152ms | United States               | 15169 |
|  7 | 108.170.234.225 |                            |  6.543ms | United States               | 15169 |
|  8 | 142.250.228.84  |                            | 18.475ms | United States               | 15169 |
|  9 | 74.125.242.241  |                            | 16.316ms | United States               | 15169 |
| 10 | 72.14.233.181   |                            | 16.423ms | United States               | 15169 |
| 11 | 216.58.214.206  | bud02s23-in-f206.1e100.net | 18.912ms | United States/Mountain View | 15169 |
+----+-----------------+----------------------------+----------+-----------------------------+-------+
```
