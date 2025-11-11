package main

import (
	"log"

	"example.com/gin-app/src"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// -----------------------------------------------------------------
	// DB初期化 (※ router.Run() より前に実行する必要があります)
	// -----------------------------------------------------------------
	if err := src.InitDB(); err != nil {
		log.Fatal(err)
	}

	// ルータ初期化
	router := gin.Default()

	// -----------------------------------------------------------------
	// (任意) 開発環境でのプロキシ警告を消す場合
	// -----------------------------------------------------------------
	router.SetTrustedProxies([]string{"127.0.0.1"})

	// ミドルウェア (CORS)
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"}, // 必要に応じて "Authorization" など具体的に指定
		AllowCredentials: true,
	}))

	// ミドルウェア (Supabase認証)
	// ※ ルート登録より前に設定
	router.Use(src.SupabaseAuthMiddleware())

	// -----------------------------------------------------------------
	// ルート登録
	// -----------------------------------------------------------------
	src.RegisterAuthRoutes(router)
	src.RegisterMemoRoutes(router)

	// -----------------------------------------------------------------
	// 静的ファイル配信 (APIルートより後に定義するのが一般的)
	// -----------------------------------------------------------------
	router.GET("/", func(c *gin.Context) {
		c.File("./src/static/index.html")
	})
	router.NoRoute(func(c *gin.Context) {
		// 👇 ここのパスも修正が必要でした
		path := "./src/static" + c.Request.URL.Path
		c.File(path)
	})

	// -----------------------------------------------------------------
	// アプリケーションの実行
	// -----------------------------------------------------------------
	var port string = src.Config.ServerPort
	if port == "" {
		port = "8180" // Configで設定されていない場合のデフォルト
	}
	log.Println("Server started on http://localhost:" + port)

	// これが最後に実行されます (ここで処理がブロックされます)
	router.Run(":" + port)
}
