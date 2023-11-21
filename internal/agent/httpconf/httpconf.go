package httpconf

import "net/http"

var HTTPClient http.Client

func Init() {
	HTTPClient = http.Client{}
}
