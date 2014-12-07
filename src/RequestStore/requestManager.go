package RequestStore

import "strconv"

//This should all be replaced with code that accesses a database so that you can
// restart without losing all past requests
var requests = make(map[string]int)

func RequestUsed(request string) (count int) {
  count = requests[request]
  return
}

func RequestAdd(request string) {
  requests[request]++
}

type DomainNotRegisteredError string

func (e DomainNotRegisteredError) Error() string {
	return "invalid domain " + strconv.Quote(string(e))
}

var keys = make(map[string]string)

func RequestKey(domain string) (key []byte, err error) {
  keys["localhost:8000"] = "An excellent key"
  err = nil
  if value, ok := keys[domain]; ok {
    key = []byte(value)
  } else {
    err = DomainNotRegisteredError(domain)
  }
  return
}
