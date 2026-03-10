package api

import (
	"weblab/internal/bootstrap"
	"weblab/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, c *bootstrap.Container) {
	v1 := r.Group("/api/v1")

	auth := v1.Group("/auth")
	auth.POST("/register", c.AuthHandler.Register)
	auth.POST("/login", c.AuthHandler.Login)

	secured := v1.Group("")
	secured.Use(middleware.JWTAuth(c.JWTManager))

	secured.GET("/auth/me", c.AuthHandler.Me)

	secured.POST("/articles", c.ArticleHandler.Publish)
	secured.GET("/articles/feed", c.ArticleHandler.Feed)
	secured.GET("/articles/recommend", c.ArticleHandler.Recommend)
	secured.GET("/articles/:id", c.ArticleHandler.GetDetail)

	secured.POST("/questions", c.QuestionHandler.Ask)
	secured.GET("/questions/recommend", c.QuestionHandler.Recommend)
	secured.GET("/questions/following", c.QuestionHandler.FollowFeed)
	secured.GET("/questions/:id", c.QuestionHandler.Detail)
	secured.POST("/questions/:id/answers", c.QuestionHandler.Answer)
	secured.GET("/questions/:id/answers", c.QuestionHandler.ListAnswers)
	secured.POST("/questions/answers/:id/vote", c.QuestionHandler.VoteAnswer)
	secured.DELETE("/questions/answers/:id/vote", c.QuestionHandler.UnvoteAnswer)

	secured.POST("/social/follow/:id", c.SocialHandler.Follow)
	secured.GET("/social/discover", c.SocialHandler.Discover)
	secured.GET("/social/mutuals", c.SocialHandler.Mutuals)
	secured.POST("/social/messages", c.SocialHandler.SendMessage)
	secured.GET("/social/messages/:userID", c.SocialHandler.Conversation)

	secured.POST("/interactions/articles/:id/like", c.InteractionHandler.LikeArticle)
	secured.DELETE("/interactions/articles/:id/like", c.InteractionHandler.UnlikeArticle)
	secured.POST("/interactions/articles/:id/comments", c.InteractionHandler.AddComment)
	secured.GET("/interactions/articles/:id/comments", c.InteractionHandler.ListComments)
}
