module github.com/ddddddO/tag-mng

go 1.13

require (
	github.com/PuerkitoBio/goquery v1.5.1
	github.com/antonlindstrom/pgstore v0.0.0-20200229204646-b08ebf1105e0
	github.com/go-chi/chi v4.1.1+incompatible
	github.com/gorilla/sessions v1.2.0
	github.com/lib/pq v1.4.0
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rs/cors v1.7.0
)

// // リモートではなく、ローカルのリポジトリを見に行くための措置
// replace (
// 	github.com/ddddddO/tag-mng/internal/api/domain => ./internal/api/domain
// 	github.com/ddddddO/tag-mng/internal/api/domain/model => ./internal/api/domain/model
// 	github.com/ddddddO/tag-mng/internal/api/handlers => ./internal/api/handlers
// 	github.com/ddddddO/tag-mng/internal/api/infra => ./internal/api/infra
// 	github.com/ddddddO/tag-mng/internal/api/usecase => ./internal/api/usecase
// )
