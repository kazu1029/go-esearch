module github.com/kazu1029/go-elastic/search-api

go 1.12

replace github.com/kazu1029/go-elastic v0.0.0-20190917002351-9cedffcaf86e => ./

require (
	github.com/gin-gonic/gin v1.4.0
	github.com/kazu1029/go-elastic v0.1.0
	github.com/mailru/easyjson v0.7.0 // indirect
	github.com/olivere/elastic v6.2.22+incompatible
)
