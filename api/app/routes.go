package app

import (
	"github.com/gin-contrib/location"
	"github.com/rgrs-x/service/api/controllers"
	"github.com/rgrs-x/service/api/models"

	"github.com/gin-gonic/gin"
)

// SetupRoutes ...
func SetupRoutes() *gin.Engine {
	router := gin.Default()
	router.Use(location.Default())

	apiV1 := router.Group("api")
	{
		apiV1.GET("/user/avatar/:name", controllers.Render)
		apiV1.GET("/partner/avatar/:name", controllers.Render)
		apiV1.GET("/file/:name", controllers.Render)
		//@ api for user version 1.0.0
		apiV1.Use(APIAuthentication())

		// Public Infomations a Client
		apiV1.GET("/partner/info/:id", controllers.PublicPartnerInfo)
		apiV1.GET("/user/info/:id", controllers.PublicUserInfo)

		auth := apiV1.Group("/auth")

		// generation
		{
			auth.POST("/user/sign_up", controllers.CreateUserAccount)
			auth.POST("/user/sign_in", controllers.AuthenticateUser)
			auth.POST("/partner/sign_up", controllers.CreatePartnerAccount)
			auth.POST("/partner/sign_in", controllers.AuthenticatePartner)

			// for only refresh token
			auth.POST("/backend/get-access-token/user", controllers.UserToken)
			auth.POST("/backend/get-access-token/partner", controllers.PartnerToken)
		}

		// Use for all
		tracking := apiV1.Group("/tracking")
		{
			tracking.POST("/read-post", controllers.ReadPost)
		}

		contents := apiV1.Group("/contents")
		{
			contents.GET("/", controllers.Pagination)
			contents.GET("/filter", controllers.Filter)
			contents.PUT("/:id/like", controllers.LikePost)
			contents.GET("/post/:id", controllers.GetPost)

			contents.GET("/partner/:id", controllers.GetPartnerContents)
			contents.GET("/company/:id", controllers.GetCompanyContents)
		}

		locationsService := apiV1.Group("/location")
		{
			locationsService.GET("/:id", controllers.FindLocation)
		}

		company := apiV1.Group("/company")
		{
			company.POST("/", controllers.CreateCompany)
			company.GET("/", controllers.SwitchGetCompany)
		}

		// User and Partner
		tags := apiV1.Group("/post")
		{
			tags.GET("/tags", controllers.GetAllTags)
		}

		mentor := apiV1.Group("/mentor")
		{
			mentor.PUT("/:id/like", controllers.LikeMentor)
		}

		// Only user
		user := apiV1.Group("/user")
		{
			user.Use(UserAuthentication())

			controllers.UploadPool = make(chan controllers.WorkerMessage, 10)
			go controllers.InitWorker(controllers.UploadPool)

			user.PUT("", controllers.UpdateUserInfo)
			user.GET("/", controllers.GetAuthUserInfo)
			user.POST("/avatar", controllers.UpdateAvatarUser)
			user.POST("/cover", controllers.UpdateUserCover)
			user.POST("/time-line", controllers.CreateUserTimeLine)
			user.PUT("/time-line/:id", controllers.UpdateUserTimeLine)
			user.DELETE("/time-line/:id", controllers.DeleteUserTimeLine)

		}

		// contents_recommand ...
		contents_recommand := apiV1.Group("/recommand")
		{
			contents_recommand.Use(UserAuthentication())
			contents_recommand.GET("/", controllers.Pagination)
		}

		// Only Partner
		partnerStandby := apiV1.Group("/partner")
		{
			partnerStandby.Use(PartnerAuthentication(models.PartnerStandby))
			{
				partnerStandby.PUT("/company/requests/", controllers.JoinRequest)
				partnerStandby.PATCH("/company/requests/", controllers.CancelRequest)
			}
		}
		partner := apiV1.Group("/partner")
		{
			partner.Use(PartnerAuthentication(models.PartnerNormal))

			controllers.UploadPool = make(chan controllers.WorkerMessage, 10)
			go controllers.InitWorker(controllers.UploadPool)

			partner.PUT("", controllers.UpdatePartnerInfo)
			partner.GET("/", controllers.GetPartnerInfo)
			partner.POST("/avatar", controllers.UpdateAvatarPartner)
			partner.POST("/cover", controllers.UpdatePartnerCover)

			partner.POST("/file", controllers.UploadFile)

			infoCompany := partner.Group("introduction")
			{
				infoCompany.POST("", controllers.CreateIntroductionPost)
				infoCompany.GET("/:id", controllers.GetIntroductionPost)
			}

			post := partner.Group("/post")
			{
				post.POST("/", controllers.CreatePost)
				post.PUT("/:id", controllers.UpdatePost)
				post.DELETE("/:id", controllers.DeletePost)
			}

			company := partner.Group("/company")
			{
				company.GET("/requests/", controllers.GetRequestList)
				company.POST("/requests/", controllers.AcceptMemberRequest)
				company.DELETE("/requests/", controllers.DeclineMemberRequest)
			}
		}

		// Only admin
		admin := apiV1.Group("/admin")
		{
			admin.POST("/sign_in", controllers.AdminSignIn)
		}

	}

	// render json document for api
	//router.NotFoundHandler = app.NotFoundHandler
	return router
}
