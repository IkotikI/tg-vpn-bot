package model

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"time"
)

type Protocol string

const (
	VMESS       Protocol = "vmess"
	VLESS       Protocol = "vless"
	DOKODEMO    Protocol = "dokodemo-door"
	HTTP        Protocol = "http"
	Trojan      Protocol = "trojan"
	Shadowsocks Protocol = "shadowsocks"
	Socks       Protocol = "socks"
	WireGuard   Protocol = "wireguard"
)

/* ---- Inbound ---- */
type Inbound struct {
	Id         int    `json:"id" form:"id" gorm:"primaryKey;autoIncrement"`
	UserId     int    `json:"-"`
	Up         int64  `json:"up" form:"up"`
	Down       int64  `json:"down" form:"down"`
	Total      int64  `json:"total" form:"total"`
	Remark     string `json:"remark" form:"remark"`
	Enable     bool   `json:"enable" form:"enable"`
	ExpiryTime int64  `json:"expiryTime" form:"expiryTime"`
	// ClientStats []xray.ClientTraffic `gorm:"foreignKey:InboundId;references:Id" json:"clientStats" form:"clientStats"`

	// config part
	Listen         string   `json:"listen" form:"listen"`
	Port           int      `json:"port" form:"port"`
	Protocol       Protocol `json:"protocol" form:"protocol"`
	Settings       string   `json:"settings" form:"settings"`
	StreamSettings string   `json:"streamSettings" form:"streamSettings"`
	Tag            string   `json:"tag" form:"tag" gorm:"unique"`
	Sniffing       string   `json:"sniffing" form:"sniffing"`
	Allocate       string   `json:"allocate" form:"allocate"`
}

type OutboundTraffics struct {
	Id    int    `json:"id" form:"id" gorm:"primaryKey;autoIncrement"`
	Tag   string `json:"tag" form:"tag" gorm:"unique"`
	Up    int64  `json:"up" form:"up" gorm:"default:0"`
	Down  int64  `json:"down" form:"down" gorm:"default:0"`
	Total int64  `json:"total" form:"total" gorm:"default:0"`
}

type InboundClientIps struct {
	Id          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	ClientEmail string `json:"clientEmail" form:"clientEmail" gorm:"unique"`
	Ips         string `json:"ips" form:"ips"`
}

type Setting struct {
	Id    int    `json:"id" form:"id" gorm:"primaryKey;autoIncrement"`
	Key   string `json:"key" form:"key"`
	Value string `json:"value" form:"value"`
}

type Client struct {
	ID         string `json:"id"` // uuid.UUID
	Security   string `json:"security"`
	Password   string `json:"password"`
	Flow       string `json:"flow"`
	Email      string `json:"email"`
	LimitIP    int    `json:"limitIp"`
	TotalGB    int64  `json:"totalGB" form:"totalGB"`
	ExpiryTime int64  `json:"expiryTime" form:"expiryTime"` // Milliseconds
	Enable     bool   `json:"enable" form:"enable"`
	TgID       int64  `json:"tgId" form:"tgId"`
	SubID      string `json:"subId" form:"subId"`
	Reset      int    `json:"reset" form:"reset"`
}

type ClientTraffic struct {
	Id         int    `json:"id" form:"id" gorm:"primaryKey;autoIncrement"`
	InboundId  int    `json:"inboundId" form:"inboundId"`
	Enable     bool   `json:"enable" form:"enable"`
	Email      string `json:"email" form:"email" gorm:"unique"`
	Up         int64  `json:"up" form:"up"`
	Down       int64  `json:"down" form:"down"`
	ExpiryTime int64  `json:"expiryTime" form:"expiryTime"`
	Total      int64  `json:"total" form:"total"`
	Reset      int    `json:"reset" form:"reset" gorm:"default:0"`
}

type Msg struct {
	Success bool        `json:"success"`
	Msg     string      `json:"msg"`
	Obj     interface{} `json:"obj"`
}

type AllSetting struct {
	WebListen        string `json:"webListen" form:"webListen"`
	WebDomain        string `json:"webDomain" form:"webDomain"`
	WebPort          int    `json:"webPort" form:"webPort"`
	WebCertFile      string `json:"webCertFile" form:"webCertFile"`
	WebKeyFile       string `json:"webKeyFile" form:"webKeyFile"`
	WebBasePath      string `json:"webBasePath" form:"webBasePath"`
	SessionMaxAge    int    `json:"sessionMaxAge" form:"sessionMaxAge"`
	PageSize         int    `json:"pageSize" form:"pageSize"`
	ExpireDiff       int    `json:"expireDiff" form:"expireDiff"`
	TrafficDiff      int    `json:"trafficDiff" form:"trafficDiff"`
	RemarkModel      string `json:"remarkModel" form:"remarkModel"`
	TgBotEnable      bool   `json:"tgBotEnable" form:"tgBotEnable"`
	TgBotToken       string `json:"tgBotToken" form:"tgBotToken"`
	TgBotProxy       string `json:"tgBotProxy" form:"tgBotProxy"`
	TgBotAPIServer   string `json:"tgBotAPIServer" form:"tgBotAPIServer"`
	TgBotChatId      string `json:"tgBotChatId" form:"tgBotChatId"`
	TgRunTime        string `json:"tgRunTime" form:"tgRunTime"`
	TgBotBackup      bool   `json:"tgBotBackup" form:"tgBotBackup"`
	TgBotLoginNotify bool   `json:"tgBotLoginNotify" form:"tgBotLoginNotify"`
	TgCpu            int    `json:"tgCpu" form:"tgCpu"`
	TgLang           string `json:"tgLang" form:"tgLang"`
	TimeLocation     string `json:"timeLocation" form:"timeLocation"`
	SecretEnable     bool   `json:"secretEnable" form:"secretEnable"`
	SubEnable        bool   `json:"subEnable" form:"subEnable"`
	SubListen        string `json:"subListen" form:"subListen"`
	SubPort          int    `json:"subPort" form:"subPort"`
	SubPath          string `json:"subPath" form:"subPath"`
	SubDomain        string `json:"subDomain" form:"subDomain"`
	SubCertFile      string `json:"subCertFile" form:"subCertFile"`
	SubKeyFile       string `json:"subKeyFile" form:"subKeyFile"`
	SubUpdates       int    `json:"subUpdates" form:"subUpdates"`
	SubEncrypt       bool   `json:"subEncrypt" form:"subEncrypt"`
	SubShowInfo      bool   `json:"subShowInfo" form:"subShowInfo"`
	SubURI           string `json:"subURI" form:"subURI"`
	SubJsonPath      string `json:"subJsonPath" form:"subJsonPath"`
	SubJsonURI       string `json:"subJsonURI" form:"subJsonURI"`
	SubJsonFragment  string `json:"subJsonFragment" form:"subJsonFragment"`
	SubJsonNoises    string `json:"subJsonNoises" form:"subJsonNoises"`
	SubJsonMux       string `json:"subJsonMux" form:"subJsonMux"`
	SubJsonRules     string `json:"subJsonRules" form:"subJsonRules"`
	Datepicker       string `json:"datepicker" form:"datepicker"`
}

func (s *AllSetting) CheckValid() error {
	if s.WebListen != "" {
		ip := net.ParseIP(s.WebListen)
		if ip == nil {
			return fmt.Errorf("web listen is not valid ip: %v", s.WebListen)
		}
	}

	if s.SubListen != "" {
		ip := net.ParseIP(s.SubListen)
		if ip == nil {
			return fmt.Errorf("Sub listen is not valid ip: %v", s.SubListen)
		}
	}

	if s.WebPort <= 0 || s.WebPort > 65535 {
		return fmt.Errorf("web port is not a valid port: %v", s.WebPort)
	}

	if s.SubPort <= 0 || s.SubPort > 65535 {
		return fmt.Errorf("Sub port is not a valid port: %v", s.SubPort)
	}

	if (s.SubPort == s.WebPort) && (s.WebListen == s.SubListen) {
		return fmt.Errorf("Sub and Web could not use same ip:port, ", s.SubListen, ":", s.SubPort, " & ", s.WebListen, ":", s.WebPort)
	}

	if s.WebCertFile != "" || s.WebKeyFile != "" {
		_, err := tls.LoadX509KeyPair(s.WebCertFile, s.WebKeyFile)
		if err != nil {
			return fmt.Errorf("cert file <%v> or key file <%v> invalid: %v", s.WebCertFile, s.WebKeyFile, err)
		}
	}

	if s.SubCertFile != "" || s.SubKeyFile != "" {
		_, err := tls.LoadX509KeyPair(s.SubCertFile, s.SubKeyFile)
		if err != nil {
			return fmt.Errorf("cert file <%v> or key file <%v> invalid: %v", s.SubCertFile, s.SubKeyFile, err)
		}
	}

	if !strings.HasPrefix(s.WebBasePath, "/") {
		s.WebBasePath = "/" + s.WebBasePath
	}
	if !strings.HasSuffix(s.WebBasePath, "/") {
		s.WebBasePath += "/"
	}
	if !strings.HasPrefix(s.SubPath, "/") {
		s.SubPath = "/" + s.SubPath
	}
	if !strings.HasSuffix(s.SubPath, "/") {
		s.SubPath += "/"
	}

	if !strings.HasPrefix(s.SubJsonPath, "/") {
		s.SubJsonPath = "/" + s.SubJsonPath
	}
	if !strings.HasSuffix(s.SubJsonPath, "/") {
		s.SubJsonPath += "/"
	}

	_, err := time.LoadLocation(s.TimeLocation)
	if err != nil {
		return fmt.Errorf("time location not exist: %v", s.TimeLocation)
	}

	return nil
}
