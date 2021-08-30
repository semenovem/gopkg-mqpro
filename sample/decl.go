package main

import "regexp"

// Регулярное выражение для валидации MsgID / CorrelID
var regMsgId = regexp.MustCompile(`(?i)^[\da-f]{48}$`)
