package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rgrs-x/service/api/factory"
	"github.com/rgrs-x/service/api/models"
	"github.com/rgrs-x/service/api/models/code"
	"github.com/rgrs-x/service/api/models/message"
	"github.com/rgrs-x/service/api/repository"
	u "github.com/rgrs-x/service/api/utils"
)

// Handle Http StatusCode with repository.Status input
func getStatusCode(status repository.Status) int {
	switch status {
	case repository.Created:
		return http.StatusCreated

	case repository.Success:
		return http.StatusOK

	case repository.Accepted:
		return http.StatusAccepted

	case repository.Unauthorized:
		return http.StatusUnauthorized

	case repository.Forbidden:
		return http.StatusForbidden

	case repository.NotFound:
		return http.StatusNotFound

	case repository.GetError:
		return http.StatusRequestTimeout

	default:
		return http.StatusInternalServerError
	}
}

// Set Status, UserMode for return Response
func handlerStatus(result repository.RepoResponse, status repository.Status, mode models.UserMode, ctx *gin.Context) {
	switch status {
	case repository.Created:
		handlerCreated(result, mode, ctx)
		return

	case repository.Success:
		handlerSuccess(result, mode, ctx)
		return

	case repository.Liked:
		handlerLiked(result, mode, ctx)
		return

	case repository.Uploaded:
		handlerUploaded(result, mode, ctx)
		return

	case repository.Deleted:
		handlerDeleted(result, mode, ctx)
		return

	case repository.Existed:
		handlerExisted(result, mode, ctx)
		return

	case repository.NotFound:
		handlerNotFound(result, mode, ctx)
		return

	case repository.GetError:
		handlerGetError(result, mode, ctx)
		return

	default:
		handlerInternalError(result, mode, ctx)
		return
	}
}

// handler Created Entity ...
func handlerCreated(result repository.RepoResponse, mode models.UserMode, ctx *gin.Context) {
	switch mode {
	case models.UserNormal:
		var userFactory factory.UserInfoFactory
		result.Data = userFactory.Create(result.Data)

		resp := u.BTResponse{Status: true, Message: message.UserCreated, Data: result.Data, Code: code.Created}
		ctx.JSON(http.StatusCreated, resp)
		return

	case models.TimeLineNormal:
		var userFactory factory.UserInfoFactory
		result.Data = userFactory.CreateDetail(result.Data)

		resp := u.BTResponse{Status: true, Message: message.TimeLineCreated, Data: result.Data, Code: code.Created}
		ctx.JSON(http.StatusCreated, resp)
		return

	case models.IntroductionNormal:
		var postFactory factory.PostInfoFactoty
		result.Data = postFactory.NewPost(result.Data)

		resp := u.BTResponse{Status: true, Message: message.IntroductionCreated, Data: result.Data, Code: code.Created}
		ctx.JSON(http.StatusCreated, resp)
		return

	default:
		resp := u.BTResponse{Status: false, Message: "Unknown Type of User", Data: []string{}, Code: code.InternalServerError}
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
}

// handler Success Resources Found ...
func handlerSuccess(result repository.RepoResponse, mode models.UserMode, ctx *gin.Context) {
	switch mode {
	case models.UserNormal:
		var userFactory factory.UserInfoFactory
		result.Data = userFactory.CreateDetail(result.Data)

		resp := u.BTResponse{Status: true, Message: message.UserSuccess, Data: result.Data, Code: code.Ok}
		ctx.JSON(http.StatusOK, resp)
		return

	case models.PartnerNormal:
		var partnerFactory factory.PartnerInfoFactory
		result.Data = partnerFactory.CreateDetail(result.Data)

		resp := u.BTResponse{Status: true, Message: message.PartnerSuccess, Data: result.Data, Code: code.Ok}
		ctx.JSON(http.StatusOK, resp)
		return

	case models.LocationNormal:
		var locationFactory factory.LocationInfoFactory
		result.Data = locationFactory.Create(result.Data)

		resp := u.BTResponse{Status: true, Message: message.LocationSuccess, Data: result.Data, Code: code.LocationSuccess}
		ctx.JSON(http.StatusOK, resp)
		return

	case models.TimeLineNormal:
		var userFactory factory.UserInfoFactory
		result.Data = userFactory.CreateDetail(result.Data)

		resp := u.BTResponse{Status: true, Message: message.Updated, Data: result.Data, Code: code.Ok}
		ctx.JSON(http.StatusOK, resp)
		return

	case models.IntroductionNormal:
		var postFactory factory.PostInfoFactoty
		result.Data = postFactory.NewPost(result.Data)

		resp := u.BTResponse{Status: true, Message: message.IntroductionSuccess, Data: result.Data, Code: code.Ok}
		ctx.JSON(http.StatusOK, resp)
		return

	default:
		resp := u.BTResponse{Status: false, Message: "Unknown Type of User", Data: []string{}, Code: code.InternalServerError}
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
}

// hander Liked for entity have Like Field ...
func handlerLiked(result repository.RepoResponse, mode models.UserMode, ctx *gin.Context) {
	switch mode {
	case models.PartnerNormal:
		var partnerFactory factory.PartnerInfoFactory
		result.Data = partnerFactory.CreateDetail(result.Data)

		resp := u.BTResponse{Status: true, Message: message.LikeSuccess, Data: result.Data, Code: code.LikeSuccess}
		ctx.JSON(http.StatusOK, resp)
		return

	default:
		resp := u.BTResponse{Status: false, Message: "Unknown Type of User", Data: []string{}, Code: code.InternalServerError}
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
}

// handlerUploaded each type ...
func handlerUploaded(result repository.RepoResponse, mode models.UserMode, ctx *gin.Context) {
	switch mode {
	case models.CoverNormal:
		resp := u.BTResponse{Status: true, Message: message.Uploaded, Data: result.Data, Code: code.Uploaded}
		ctx.JSON(http.StatusOK, resp)
		return

	default:
		resp := u.BTResponse{Status: false, Message: "Unknown Type of User", Data: []string{}, Code: code.InternalServerError}
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
}

// handler for Deleted Resource ...
func handlerDeleted(result repository.RepoResponse, mode models.UserMode, ctx *gin.Context) {
	switch mode {
	case models.TimeLineNormal:
		resp := u.BTResponse{Status: true, Message: message.TimeLineDeleted, Data: []string{}, Code: code.Deleted}
		ctx.JSON(http.StatusOK, resp)
		return

	default:
		resp := u.BTResponse{Status: false, Message: "Unknown Type of User", Data: []string{}, Code: code.InternalServerError}
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
}

// handler Existed Entity ...
func handlerExisted(result repository.RepoResponse, mode models.UserMode, ctx *gin.Context) {
	switch mode {
	case models.UserNormal:
		resp := u.BTResponse{Status: false, Message: message.EmailIsUsed, Data: []string{}, Code: code.EmailIsUsed}
		ctx.JSON(http.StatusForbidden, resp)
		return

	case models.IntroductionNormal:
		resp := u.BTResponse{Status: false, Message: message.IntroductionExited, Data: []string{}, Code: code.PostError}
		ctx.JSON(http.StatusForbidden, resp)
		return

	default:
		resp := u.BTResponse{Status: false, Message: "Unknown Type of User", Data: []string{}, Code: code.InternalServerError}
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
}

// handler Not Found Resources ...
func handlerNotFound(result repository.RepoResponse, mode models.UserMode, ctx *gin.Context) {
	switch mode {
	case models.UserNormal:
		resp := u.BTResponse{Status: false, Message: message.ResourceNotFound, Data: []string{}, Code: code.ResourceError}
		ctx.JSON(http.StatusNotFound, resp)
		return

	case models.PartnerNormal:
		resp := u.BTResponse{Status: false, Message: message.ResourceNotFound, Data: []string{}, Code: code.ResourceError}
		ctx.JSON(http.StatusNotFound, resp)
		return

	case models.LocationNormal:
		resp := u.BTResponse{Status: false, Message: message.LocationError, Data: []string{}, Code: code.LocationError}
		ctx.JSON(http.StatusNotFound, resp)
		return

	case models.CoverNormal:
		resp := u.BTResponse{Status: false, Message: message.NotFound, Data: []string{}, Code: code.UploadError}
		ctx.JSON(http.StatusNotFound, resp)
		return

	case models.IntroductionNormal:
		resp := u.BTResponse{Status: false, Message: message.NotFound, Data: []string{}, Code: code.ResourceError}
		ctx.JSON(http.StatusNotFound, resp)
		return

	case models.TimeLineNormal:
		resp := u.BTResponse{Status: false, Message: message.NotFound, Data: []string{}, Code: code.ResourceError}
		ctx.JSON(http.StatusNotFound, resp)
		return

	default:
		resp := u.BTResponse{Status: false, Message: "Unknown Type of User", Data: []string{}, Code: code.InternalServerError}
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
}

// handler Error Database Query ...
func handlerGetError(result repository.RepoResponse, mode models.UserMode, ctx *gin.Context) {
	switch mode {
	case models.UserNormal:
		resp := u.BTResponse{Status: false, Message: message.QueryError, Data: []string{}, Code: code.QueryError}
		ctx.JSON(http.StatusRequestTimeout, resp)
		return

	case models.PartnerNormal:
		resp := u.BTResponse{Status: false, Message: message.QueryError, Data: []string{}, Code: code.QueryError}
		ctx.JSON(http.StatusRequestTimeout, resp)
		return

	case models.LocationNormal:
		resp := u.BTResponse{Status: false, Message: message.QueryError, Data: []string{}, Code: code.QueryError}
		ctx.JSON(http.StatusRequestTimeout, resp)
		return

	case models.IntroductionNormal:
		resp := u.BTResponse{Status: false, Message: message.QueryError, Data: []string{}, Code: code.QueryError}
		ctx.JSON(http.StatusRequestTimeout, resp)
		return

	case models.CoverNormal:
		resp := u.BTResponse{Status: false, Message: message.QueryError, Data: []string{}, Code: code.QueryError}
		ctx.JSON(http.StatusRequestTimeout, resp)
		return

	case models.TimeLineNormal:
		resp := u.BTResponse{Status: false, Message: message.QueryError, Data: []string{}, Code: code.QueryError}
		ctx.JSON(http.StatusRequestTimeout, resp)
		return

	default:
		resp := u.BTResponse{Status: false, Message: "Unknown Type of User", Data: []string{}, Code: code.InternalServerError}
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
}

// handler Internal Server Error ...
func handlerInternalError(result repository.RepoResponse, mode models.UserMode, ctx *gin.Context) {
	switch mode {
	case models.UserNormal:
		resp := u.BTResponse{Status: false, Message: message.InternalServerError, Data: []string{}, Code: code.InternalServerError}
		ctx.JSON(http.StatusInternalServerError, resp)
		return

	case models.PartnerNormal:
		resp := u.BTResponse{Status: false, Message: message.InternalServerError, Data: []string{}, Code: code.InternalServerError}
		ctx.JSON(http.StatusInternalServerError, resp)
		return

	case models.LocationNormal:
		resp := u.BTResponse{Status: false, Message: message.QueryError, Data: []string{}, Code: code.QueryError}
		ctx.JSON(http.StatusRequestTimeout, resp)
		return

	case models.CoverNormal:
		resp := u.BTResponse{Status: false, Message: message.QueryError, Data: []string{}, Code: code.QueryError}
		ctx.JSON(http.StatusRequestTimeout, resp)
		return

	case models.TimeLineNormal:
		resp := u.BTResponse{Status: false, Message: message.QueryError, Data: []string{}, Code: code.QueryError}
		ctx.JSON(http.StatusRequestTimeout, resp)
		return

	default:
		resp := u.BTResponse{Status: false, Message: "Unknown Type of User", Data: []string{}, Code: code.InternalServerError}
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
}
