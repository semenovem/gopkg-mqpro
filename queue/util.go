package queue

import (
  "fmt"
  "strings"
  "unicode/utf8"
)

func unionProps(dst map[string]interface{}, src map[string]interface{}) {
  for n, v := range src {
    dst[n] = v
  }
}

func unionPropsDeep(dst map[string]interface{}, src []map[string]interface{}) {
  for _, a := range src {
    for n, v := range a {
      dst[n] = v
    }
  }
}

func tailFour(n int) int {
  r := n % 4
  if r == 0 {
    return 0
  }
  return 4 - r
}

func tailFour32(n int32) int32 {
  r := n % 4
  if r == 0 {
    return 0
  }
  return 4 - r
}

var msgErrParseQueue = "строка '%s' не соответствует формату {name очереди}:{[PUT,GET,BROWSE]}"

func parseQueue(s string) (nameQue string, perm []permQueue, err error) {
  fnDuplicates := func(s permQueue) bool {
    for _, v := range perm {
      if v == s {
        return true
      }
    }
    return false
  }

  arg := strings.Split(s, ":")
  if len(arg) != 2 {
    err = fmt.Errorf(msgErrParseQueue, s)
    return
  }
  nameQue = strings.TrimSpace(arg[0])

  if nameQue == "" {
    err = fmt.Errorf("пустое значение названия очереди. "+msgErrParseQueue, s)
    return
  }

  a := strings.Split(arg[1], ",")

  for _, v := range a {
    vv, ok := permMapByVal[v]
    if !ok {
      err = fmt.Errorf("не валидное значение типа очереди. "+msgErrParseQueue, s)
      return
    }
    if fnDuplicates(vv) {
      err = fmt.Errorf("повторяющееся значение типа очереди. "+msgErrParseQueue, s)
      return
    }
    perm = append(perm, vv)
  }

  return
}

func PrintSetCli(set []map[string]string, p string) {
  var (
    max, l int
    k, v   string
    m      map[string]string
  )
  for _, m = range set {
    for k, v = range m {
      l = utf8.RuneCountInString(k)
      if max < l {
        max = l
      }
    }
  }
  max += 1

  for _, m = range set {
    for k, v = range m {
      v = strings.Repeat(" ", max-utf8.RuneCountInString(k)) + "= " + v
      if p != "" {
        k = p + "/" + k
      }
      fmt.Printf("%s%s\n", k, v)
    }
  }
}
