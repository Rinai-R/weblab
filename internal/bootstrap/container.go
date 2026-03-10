package bootstrap

import (
	"database/sql"
	"time"
	"weblab/internal/api/handler"
	"weblab/internal/dao"
	mysqldao "weblab/internal/dao/mysql"
	"weblab/internal/service"
	"weblab/internal/utils"
)

type Container struct {
	DB *sql.DB

	JWTManager *utils.JWTManager

	UserDAO        dao.UserDAO
	ArticleDAO     dao.ArticleDAO
	QuestionDAO    dao.QuestionDAO
	MessageDAO     dao.MessageDAO
	InteractionDAO dao.InteractionDAO

	AuthService        *service.AuthService
	ArticleService     *service.ArticleService
	QuestionService    *service.QuestionService
	SocialService      *service.SocialService
	InteractionService *service.InteractionService

	AuthHandler        *handler.AuthHandler
	ArticleHandler     *handler.ArticleHandler
	QuestionHandler    *handler.QuestionHandler
	SocialHandler      *handler.SocialHandler
	InteractionHandler *handler.InteractionHandler
}

func BuildContainer(jwtSecret string) (*Container, error) {
	jwtMgr := utils.NewJWTManager(jwtSecret, 24*time.Hour)

	db, err := mysqldao.OpenFromEnv()
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := mysqldao.EnsureSchema(db); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}

	userDAO := mysqldao.NewUserDAO(db)
	articleDAO := mysqldao.NewArticleDAO(db)
	questionDAO := mysqldao.NewQuestionDAO(db)
	messageDAO := mysqldao.NewMessageDAO(db)
	interactionDAO := mysqldao.NewInteractionDAO(db)

	authSvc := service.NewAuthService(userDAO, jwtMgr)
	articleSvc := service.NewArticleService(articleDAO, userDAO, interactionDAO)
	questionSvc := service.NewQuestionService(questionDAO, userDAO)
	socialSvc := service.NewSocialService(userDAO, messageDAO)
	interactionSvc := service.NewInteractionService(interactionDAO, articleDAO)

	return &Container{
		DB: sqlDB,

		JWTManager: jwtMgr,

		UserDAO:        userDAO,
		ArticleDAO:     articleDAO,
		QuestionDAO:    questionDAO,
		MessageDAO:     messageDAO,
		InteractionDAO: interactionDAO,

		AuthService:        authSvc,
		ArticleService:     articleSvc,
		QuestionService:    questionSvc,
		SocialService:      socialSvc,
		InteractionService: interactionSvc,

		AuthHandler:        handler.NewAuthHandler(authSvc),
		ArticleHandler:     handler.NewArticleHandler(articleSvc),
		QuestionHandler:    handler.NewQuestionHandler(questionSvc),
		SocialHandler:      handler.NewSocialHandler(socialSvc),
		InteractionHandler: handler.NewInteractionHandler(interactionSvc),
	}, nil
}
