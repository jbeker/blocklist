package main

import (
  "fmt"
  "net"
)

type BitNode struct {
  parent,zero,one *BitNode
  depth uint
  full bool
  value uint32
}

func main () {

  root := BitNode{parent: nil, zero :nil, one :nil, depth: 32, full: false, value: 0}


  addIP(&root,IPtoInt(net.ParseIP("128.239.1.1"))) 
  addIP(&root,IPtoInt(net.ParseIP("128.239.1.2"))) 
  addIP(&root,IPtoInt(net.ParseIP("128.239.1.3"))) 
  
  outputTree(&root)

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

func CheckBit(num uint32, bit uint) bool {
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


func addIP(node *BitNode, addr uint32) bool {

  if node.depth == 0 || node.full {
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
  
  addIP(child, addr)
  
  node.full = (node.one != nil && node.one.full) && (node.zero != nil && node.zero.full) 
  return node.full
}
