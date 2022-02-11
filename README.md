### mqm (manager queue messages)

Обертка над [mq-golang](https://github.com/ibm-messaging/mq-golang)

### Создание / инициализация / подключение
```
import "github.com/semenovem/mqm/v2"

var (
  rootCtx, rootCtxEsc = context.WithCancel(context.Background())
  log = logrus.NewEntry(logrus.New())
  
  mq = mqm.New(rootCtx, log)
  que = mq.NewPipe("aliasQueue") // 'aliasQueue' - название пары очередей в файле конфигурации
)

func init() {
  cfg, err := mqm.ParseCfgYaml('путь к файлу конфигурации')
  if err != nil {
    log.Errorf("Ошибка парсинга файла конфигурации MQM '%s'", err)
  } else {
    // Тут можно изменить/дополнить конфигурацию, например: 
    cfg.Pass = "password"
  
    err = mq.Cfg(cfg)
    if err != nil {
      log.Error("Ошибка при установке конфигурации: ", err)
    }
  }
}

func main() {
  go func() {
    err := mq.Connect()
    if err != nil {
      log.Panic("Не удалось запуститься")
    }
  }()
}
```


### Отправка сообщения
```
ctx, cancel := context.WithTimeout(rootCtx, time.Second*10)
defer cancel()

// Свойства сообщения
props := map[string]interface{}{
  "foo": "10101001110110",
  "BAR": "cb31e8610231",
}

payload := []byte(`{"HoldJetFuelPaymentMsg":{"id":"f021d4ec-27f5-41be-8af3-946e65686902","result":"OK"}}`)

msg := &queue.Msg{
  Payload: payload,
  Props:   props,
}

err := que.Put(ctx, msg)

fmt.Printf("%s\n", err)
fmt.Printf("%+v\n", msg)
```


### Получение очередного сообщения
```
ctx, cancel := context.WithTimeout(rootCtx, time.Second*10)
defer cancel()

var msg = &queue.Msg{}
err := que.Get(ctx, msg)

fmt.Printf("%s\n", err)
fmt.Printf("%+v\n", msg)
```



### Получение сообщения по CorrelId
```
ctx, cancel := context.WithTimeout(rootCtx, time.Second*10)
defer cancel()

var CorrelId = []byte("x234123412341234213")

var msg = &queue.Msg{ CorrelId: CorrelId}
err := que.Get(ctx, msg)

fmt.Printf("%s\n", err)
fmt.Printf("%+v\n", msg)
```

