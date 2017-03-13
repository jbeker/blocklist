package main

import (
  "fmt"
  "net"
  "os"
  "bufio"
)

type BitNode struct {
  parent,zero,one *BitNode
  depth uint32
  full bool
  value uint32
}

func main () {

  root := BitNode{parent: nil, zero :nil, one :nil, depth: 32, full: false, value: 0}

  scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
    line := scanner.Text()
    ipnet :=  StringToIPNet(line)

    if ipnet != nil {
      var bits,_ = ipnet.Mask.Size()
      addIP(&root,IPtoInt(ipnet.IP),uint32(bits))
    }

	}
  
  outputTree(&root)

}

func StringToIPNet(text string) *net.IPNet {
  var ip,ipnet,error = net.ParseCIDR(text)
  
  if error != nil {
    ip = net.ParseIP(text)
    if ip != nil {
      ipnet = &net.IPNet{IP: ip, Mask: net.CIDRMask(32,32)}
    }
  }
  
  return ipnet
}

func IPtoInt(addr net.IP) uint32 {
  return uint32(addr.To4()[0]) << 24 |
         uint32(addr.To4()[1]) << 16 |
         uint32(addr.To4()[2]) << 8 |
         uint32(addr.To4()[3])
}

func IntToIP(ip uint32) net.IP {
  result := make(net.IP, 4)
           result[3] = byte(ip)
           result[2] = byte(ip >>8)
           result[1] = byte(ip >>16)
           result[0] = byte(ip >> 24)
  return result
}

func CheckBit(num uint32, bit uint32) bool {
  return (num & (1 << bit)) != 0
}

func outputTree(node *BitNode) {
  if node.full {
    var ip = IPFromNode(node)
    var mask = 32 - node.depth
    fmt.Printf("%v/%d\n",ip,mask) 
  } else {
    if node.zero != nil {
      outputTree(node.zero)
    }
  
    if node.one != nil {
      outputTree(node.one)
    }
  }
}

func IPFromNode(node *BitNode) net.IP {
  var cur = node
  var accumulate uint32 = 0
  
  for cur.parent != nil {
    accumulate |= cur.value << cur.depth
    cur = cur.parent
  }
  
  return IntToIP(accumulate)
}


func addIP(node *BitNode, addr uint32, mask uint32) bool {

  if node.depth == 0 || 32 - node.depth == mask || node.full {
    node.full = true
    return node.full
  }
  
  var child *BitNode
    
  if CheckBit(addr,node.depth - 1) {
    if node.one == nil {
      child = &BitNode{parent: node, zero: nil, one: nil, depth: node.depth -1, full: false, value: 1}    
      node.one = child
    } else {
      child = node.one
    }
  } else {
    if node.zero == nil {
      child = &BitNode{parent: node, zero: nil, one: nil, depth: node.depth -1, full: false, value: 0 }    
      node.zero = child
    } else {
      child = node.zero
    }
  }
  
  addIP(child, addr, mask)
  
  node.full = (node.one != nil && node.one.full) && (node.zero != nil && node.zero.full) 
  return node.full
}
