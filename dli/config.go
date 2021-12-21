package dli

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Config struct {
	username string
	password string
	addr     string
}

func (c Config) setDigestAuth(r *http.Request, resp *http.Response) {
	headers := resp.Header["Www-Authenticate"]
	opts := make(map[string]string)
	for _, header := range headers {
		parts2 := strings.SplitN(header, " ", 2)
		parts2 = strings.Split(parts2[1], ", ")
		for _, part := range parts2 {
			vals := strings.SplitN(part, "=", 2)
			key := vals[0]
			val := strings.Trim(vals[1], "\",")
			opts[key] = val
		}
	}
	realm := opts["realm"]
	nonce := opts["nonce"]
	opaque := opts["opaque"]
	a1 := c.username + ":" + realm + ":" + c.password
	h := md5.New()
	io.WriteString(h, a1)
	ha1 := hex.EncodeToString(h.Sum(nil))
	h = md5.New()
	a2 := "PUT" + ":" + r.URL.Path
	io.WriteString(h, a2)
	ha2 := hex.EncodeToString(h.Sum(nil))
	nc_str := fmt.Sprintf("%08x", 1)
	hnc := "OWE4NmEwZGFkNjgzN2NiMjFiZmRmNzg5YjQwYzk5ZTA="
	respdig := fmt.Sprintf("%s:%s:%s:%s:%s:%s", ha1, nonce, nc_str, hnc, "auth", ha2)
	h = md5.New()
	io.WriteString(h, respdig)
	respdig = hex.EncodeToString(h.Sum(nil))
	format := `username="%s", realm="%s", nonce="%s", uri="%s", cnonce="%s", nc=%s, qop=auth, response="%s", opaque="%s", algorithm="MD5"`
	digest := fmt.Sprintf(format,
		c.username,
		realm,
		nonce,
		r.URL.Path,
		hnc,
		nc_str,
		respdig,
		opaque,
	)
	auth_str := "Digest " + digest
	r.Header.Add("Authorization", auth_str)
}
