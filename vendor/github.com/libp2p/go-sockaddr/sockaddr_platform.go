package sockaddr

// import (
//   "unix"
//   "unsafe"
// )

// func sockaddrToAny(sa unix.Sockaddr) (*unix.RawSockaddrAny, Socklen, error) {
//   if sa == nil {
//     return nil, 0, unix.EINVAL
//   }

//   switch sa.(type) {
//   case *unix.SockaddrInet4:
//   case *unix.SockaddrInet6:
//   case *unix.SockaddrUnix:
//   case *unix.SockaddrDatalink:
//   }
//   return nil, 0, unix.EAFNOSUPPORT
// }

// func anyToSockaddr(rsa *unix.RawSockaddrAny) (unix.Sockaddr, error) {
//   if rsa == nil {
//     return nil, 0, unix.EINVAL
//   }

//   switch rsa.Addr.Family {
//   case unix.AF_NETLINK:
//   case unix.AF_PACKET:
//   case unix.AF_UNIX:
//   case unix.AF_INET:
//   case unix.AF_INET6:
//   }
//   return nil, unix.EAFNOSUPPORT
// }
