# Introduction

## What is Betxin?

Betxin is an information trading market platform where you can trade the most controversial topics in the world such as (viruses, politics, cryptocurrencies, facts, etc.). On Betxin, you can build investment forecast portfolios based on your own forecasts. If your investment forecasts are correct, you will be rewarded. When you decide to make forecast investments in the market, you are actually weighing your own knowledge, Mind and vision for the future. Market ratios reflect to some extent people's views on the future over a period of time in the past. Here you can buy according to your forecast, and you can also choose to sell if the forecast does not meet expectations.

You can try with [https://betxin.one](https://betxin.one)

If you want to run on your computer.

```
mkdir configs
touch config.ini
./scripts/run.sh
```

```

├── cmd
├── configs
├── docs
├── internal
│   ├── api
│   │   ├── sd
│   │   └── v1
│   │       ├── administrator
│   │       ├── bonuse
│   │       ├── category
│   │       ├── collect
│   │       ├── comment
│   │       ├── currency
│   │       ├── feedback
│   │       ├── handler.go
│   │       ├── message
│   │       ├── mixinorder
│   │       ├── mixpayorder
│   │       ├── oauth
│   │       ├── praisecomment
│   │       ├── sendback
│   │       ├── snapshot
│   │       ├── swaporder
│   │       ├── topic
│   │       ├── upload
│   │       ├── user
│   │       └── usertotopic
│   ├── model
│   ├── router
│   ├── service
│   │   ├── auth.go
│   │   ├── dailycurrency
│   │   ├── message.go
│   │   ├── mixinclient.go
│   │   ├── mixpay
│   │   ├── refund.go
│   │   ├── snapshots.go
│   │   ├── stoptopic.go
│   │   ├── transfer.go
│   │   └── worker.go
│   └── utils
│       ├── cors
│       ├── errmsg
│       ├── jwt
│       ├── logger
│       ├── redis
│       ├── session
│       ├── setting.go
│       └── upload
├── log
│   ├── log
├── main.go
├── main_test.go
├── pkg
│   ├── convert
│   ├── mq
│   └── timewheel
├── README.md
├── scripts
└── web
    ├── admin
    └── front
```

TODO