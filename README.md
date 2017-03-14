## Introduction

The idea for this project is based on the work done in the [exabgp-edgerouter](https://github.com/infowolfe/exabgp-edgerouter) by [infowolfe](https://github.com/infowolfe). My thanks to him for the inspiration.

There are numerous lists of IP addresses and ranges that represent "Bad Actors" that you might want to block from exchanging packets with your systems. This can be done through firewalls and host packets filters. However, this is often hard to integrate with routers in an automated fashion and is also not the best from a performance standpoint when you have thousands (or tens of thousands) entries.

Implementing the blocklists as a BGP feed that is then Null-routed on your router is a great way to implement this solution. Obviously assuming you have a router that can support this. Home routers will probably not. 

## Assumptions

* You understand how BGP works (enough)
* You have a router capable of being configured to receive a BGP feed and Null-route networks from that feed
* A computer that can provide the BGP feed to the router (preferably behind the security perimeter of your network)
* A [Go](https://golang.org) compiler set up 

## Setup

* Download the code in this repository
* Open up `blocklist.go` and configure the elements at the head of the file for your situation.
  * The blocklists you want to subscribe to
  * The interval to refresh things (don't make it less than 30 minutes)
  * The proper route announcement and withdrawal syntax for your setup
* Compile the `blocklist` application `go build blocklist.go`
* Install and configure [ExaBGP](https://github.com/Exa-Networks/exabgp)
  * Get it peering with your router
  * Have it use the `blocklist` application to provide routes
* Fire it up

### Example `exabgp` Config File

```
group AS65332 {
        router-id 192.168.1.1;
        local-as 65332;
        local-address 192.168.1.2;
        peer-as 65256;

        neighbor 192.168.1.1 {
                family {
                        ipv4 unicast;
                        ipv4 multicast;
                }
        }

        family {
                ipv4 unicast;
                ipv4 multicast;
        }

        process droproutes {
                run /wherever/you/put/the/application/blocklist;
        }
}
```

## Motivation

While the [exabgp-edgerouter](https://github.com/infowolfe/exabgp-edgerouter) provided the functionality that I wanted, the performance was not ideal as the blocklists grew in size. For a list of approximately 2000 entries, it would take about 90 seconds to process, deduplicate, and consolidate into CIDR blocks. When I increased the lists I wanted to follow to ones that composed of approximated 45000 entries, the script was still running 90 minutes later. This wasn't going to work. So, I rewrote the algorithm to be more efficient. A 45000 long entry is now processed in under a second.
