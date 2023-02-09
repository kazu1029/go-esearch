module github.com/kazu1029/go-elastic/search-api

go 1.12

require (
	github.com/gin-gonic/gin v1.7.7
	github.com/kazu1029/go-elastic v0.1.2
	github.com/mailru/easyjson v0.7.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/olivere/elastic v6.2.22+incompatible
)

replace (
	github.com/kazu1029/go-elastic v0.1.2 => ./..
	github.com/kazu1029/go-elastic/esearch v0.1.2 => ./../esearch
)
