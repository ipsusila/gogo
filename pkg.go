package gogo

import "github.com/ipsusila/gogo/http"

//NewHTTPFormUploader creates form uploder instance.
func NewHTTPFormUploader() http.FormUploader {
	return http.NewFormUploader()
}
