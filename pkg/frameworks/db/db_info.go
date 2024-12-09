// frameworks/db/db_info.go
package db

type DbInfo struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
    SSLMode  string
}

func NewDbInfo(host string, port string, user string, password string, name string, ssl_mode string) *DbInfo {
	if host == "" {
		panic("Host is not configured")
	}
	return &DbInfo{Host: host, Port: port, User: user, Password: password, DBName: name, SSLMode: ssl_mode}
}