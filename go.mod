module github.com/fabdem/go-crowdinv2

go 1.20

require github.com/mreiferson/go-httpclient v0.0.0-20201222173833-5e475fde3a4d

// Get version from perforce rather than github
replace github.com/mreiferson/go-httpclient => ../../mreiferson/go-httpclient
