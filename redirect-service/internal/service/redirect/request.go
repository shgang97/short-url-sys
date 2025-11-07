package redirect

// RedirectRequest 重定向请求
type RedirectRequest struct {
	OriginalURL string
	IPAddress   string
	UserAgent   string
	Referer     string
	Username    string
}
