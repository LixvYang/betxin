package utils

import (
	"log"
	"sync"

	"gopkg.in/ini.v1"
)

var (
	DomainName string
	IP         string
	AppMode    string
	HttpPort   string
	JwtKey     string

	Username string
	Password string

	Db         string
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassWord string
	DbName     string

	Pin        string
	ClientId   string
	SessionId  string
	PinToken   string
	PrivateKey string
	AppSecret  string

	Zone       string
	AccessKey  string
	SecretKey  string
	Bucket     string
	QiniuSever string

	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisSecret   string

	PUSD  string
	BTC   string
	BOX   string
	XIN   string
	ETH   string
	MOB   string
	USDC  string
	USDT  string
	EOS   string
	SOL   string
	UNI   string
	DOGE  string
	RUM   string
	DOT   string
	WOO   string
	ZEC   string
	LTC   string
	SHIB  string
	BCH   string
	MANA  string
	FIL   string
	BNB   string
	XRP   string
	SC    string
	MATIC string
	ETC   string
	XMR   string
	DCR   string
	TRX   string
	ATOM  string
	CKB   string
	LINK  string
	GTC   string
	HNS   string
	DASH  string
	XLM   string
	ZEN   string
)

var InitSetting sync.Once

func Init() {
	// start := time.Now()
	f, err := ini.Load("configs/config.ini")
	if err != nil {
		log.Printf("配置文件读取错误:%s", err)
	}

	LoadServer(f)
	LoadAdmin(f)
	LoadData(f)
	LoadMixinBot(f)
	LoadRedis(f)
	LoadQiniu(f)
	LoadMixinAssetId(f)
	// expired := time.Now().Sub(start)
	// fmt.Println(expired)
}

func LoadServer(file *ini.File) {
	DomainName = file.Section("server").Key("DomainName").MustString("")
	IP = file.Section("server").Key("IP").MustString("")
	AppMode = file.Section("server").Key("AppMode").MustString("")
	HttpPort = file.Section("server").Key("HttpPort").MustString(":8080")
	JwtKey = file.Section("server").Key("JwtKey").MustString("")
}

func LoadAdmin(file *ini.File) {
	Username = file.Section("administrator").Key("Username").MustString("")
	Password = file.Section("administrator").Key("Password").MustString("")
}

func LoadData(file *ini.File) {
	Db = file.Section("database").Key("Db").MustString("")
	DbHost = file.Section("database").Key("DbHost").MustString("")
	DbPort = file.Section("database").Key("DbPort").MustString("")
	DbUser = file.Section("database").Key("DbUser").MustString("")
	DbPassWord = file.Section("database").Key("DbPassWord").MustString("")
	DbName = file.Section("database").Key("DbName").MustString("")
}

func LoadRedis(file *ini.File) {
	RedisHost = file.Section("redis").Key("RedisHost").MustString("")
	RedisPort = file.Section("redis").Key("RedisPort").MustString("")
	RedisPassword = file.Section("redis").Key("RedisPassword").MustString("")
	RedisSecret = file.Section("redis").Key("RedisSecret").MustString("")
}

func LoadQiniu(file *ini.File) {
	Zone = file.Section("qiniu").Key("Zone").MustString("")
	AccessKey = file.Section("qiniu").Key("AccessKey").MustString("")
	SecretKey = file.Section("qiniu").Key("SecretKey").MustString("")
	Bucket = file.Section("qiniu").Key("Bucket").MustString("")
	QiniuSever = file.Section("qiniu").Key("QiniuSever").MustString("")
}

func LoadMixinBot(file *ini.File) {
	Pin = file.Section("mixinbot").Key("Pin").MustString("")
	ClientId = file.Section("mixinbot").Key("ClientId").MustString("")
	SessionId = file.Section("mixinbot").Key("SessionId").MustString("")
	PinToken = file.Section("mixinbot").Key("PinToken").MustString("")
	PrivateKey = file.Section("mixinbot").Key("PrivateKey").MustString("")
	AppSecret = file.Section("mixinbot").Key("AppSecret").MustString("")
}

func LoadMixinAssetId(file *ini.File) {
	PUSD = file.Section("MixinAssetId").Key("PUSD").MustString("")
	BTC = file.Section("MixinAssetId").Key("BTC").MustString("")
	BOX = file.Section("MixinAssetId").Key("BOX").MustString("")
	XIN = file.Section("MixinAssetId").Key("XIN").MustString("")
	ETH = file.Section("MixinAssetId").Key("ETH").MustString("")
	MOB = file.Section("MixinAssetId").Key("MOB").MustString("")
	USDC = file.Section("MixinAssetId").Key("USDC").MustString("")
	USDT = file.Section("MixinAssetId").Key("USDT").MustString("")
	EOS = file.Section("MixinAssetId").Key("EOS").MustString("")
	SOL = file.Section("MixinAssetId").Key("SOL").MustString("")
	UNI = file.Section("MixinAssetId").Key("UNI").MustString("")
	DOGE = file.Section("MixinAssetId").Key("DOGE").MustString("")
	RUM = file.Section("MixinAssetId").Key("RUM").MustString("")
	DOT = file.Section("MixinAssetId").Key("DOT").MustString("")
	WOO = file.Section("MixinAssetId").Key("WOO").MustString("")
	ZEC = file.Section("MixinAssetId").Key("ZEC").MustString("")
	LTC = file.Section("MixinAssetId").Key("LTC").MustString("")
	SHIB = file.Section("MixinAssetId").Key("SHIB").MustString("")
	BCH = file.Section("MixinAssetId").Key("BCH").MustString("")
	MANA = file.Section("MixinAssetId").Key("MANA").MustString("")
	FIL = file.Section("MixinAssetId").Key("FIL").MustString("")
	BNB = file.Section("MixinAssetId").Key("BNB").MustString("")
	XRP = file.Section("MixinAssetId").Key("XRP").MustString("")
	SC = file.Section("MixinAssetId").Key("SC").MustString("")
	MATIC = file.Section("MixinAssetId").Key("MATIC").MustString("")
	ETC = file.Section("MixinAssetId").Key("ETC").MustString("")
	XMR = file.Section("MixinAssetId").Key("XMR").MustString("")
	DCR = file.Section("MixinAssetId").Key("DCR").MustString("")
	TRX = file.Section("MixinAssetId").Key("TRX").MustString("")
	ATOM = file.Section("MixinAssetId").Key("ATOM").MustString("")
	CKB = file.Section("MixinAssetId").Key("CKB").MustString("")
	LINK = file.Section("MixinAssetId").Key("LINK").MustString("")
	GTC = file.Section("MixinAssetId").Key("GTC").MustString("")
	HNS = file.Section("MixinAssetId").Key("HNS").MustString("")
	DASH = file.Section("MixinAssetId").Key("DASH").MustString("")
	XLM = file.Section("MixinAssetId").Key("XLM").MustString("")
	ZEN = file.Section("MixinAssetId").Key("ZEN").MustString("")
}
