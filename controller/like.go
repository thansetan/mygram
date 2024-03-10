package controller

import (
	"errors"
	"final-project/dto"
	"final-project/helper"
	"final-project/helper/response"
	"final-project/service"
	"net/http"
	"strconv"
)

type likeController struct {
	likeService service.LikeService
}

func NewLikeController(likeService service.LikeService) *likeController {
	return &likeController{likeService}
}

// Create godoc
// @Summary Create a like
// @Description Create a like
// @Tags Like
// @Produce json
// @Param photoID path int true "Photo ID"
// @Security BearerToken
// @Success 201 {object} dto.LikeCreateResponse
// @Failure 400 {object} helper.ResponseError
// @Failure 401 {object} helper.ResponseError
// @Failure 404 {object} helper.ResponseError
// @Failure 500 {object} helper.ResponseError
// @Router /photos/{photoID}/likes [post]
func (c *likeController) Create(w http.ResponseWriter, r *http.Request) {
	var (
		data dto.LikeRequest
		resp = response.New[dto.LikeCreateResponse](response.LikeCreate)
		err  error
	)

	photoIDStr := r.PathValue("photoID")
	data.PhotoID, err = strconv.ParseUint(photoIDStr, 10, 64)
	if err != nil {
		resp.Error(err).Code(http.StatusBadRequest).Send(w)
		return
	}

	like, err := c.likeService.Create(r.Context(), data)
	if err != nil {
		respErr := new(helper.ResponseError)
		if errors.As(err, &respErr) {
			resp.Error(respErr).Code(respErr.Code()).Send(w)
			return
		}
		resp.Error(err).Code(http.StatusInternalServerError).Send(w)
		return
	}

	resp.Success(true).Data(like).Code(http.StatusCreated).Send(w)
}

// FindByPhotoID godoc
// @Summary Get likes by photo ID
// @Description Get likes by photo ID
// @Tags Like
// @Produce json
// @Param photoID path int true "Photo ID"
// @Security BearerToken
// @Success 200 {array} dto.LikeResponse
// @Failure 400 {object} helper.ResponseError
// @Failure 401 {object} helper.ResponseError
// @Failure 404 {object} helper.ResponseError
// @Failure 500 {object} helper.ResponseError
// @Router /photos/{photoID}/likes [get]
func (c *likeController) FindByPhotoID(w http.ResponseWriter, r *http.Request) {
	var (
		resp = response.New[[]dto.LikeResponse](response.LikeFindByPhotoID)
		err  error
	)

	photoIDStr := r.PathValue("photoID")
	photoID, err := strconv.ParseUint(photoIDStr, 10, 64)
	if err != nil {
		resp.Error(err).Code(http.StatusBadRequest).Send(w)
		return
	}

	likes, err := c.likeService.GetByPhotoID(r.Context(), photoID)
	if err != nil {
		resp.Error(err).Code(http.StatusInternalServerError).Send(w)
		return
	}

	resp.Success(true).Data(likes).Code(http.StatusOK).Send(w)
}

// Delete godoc
// @Summary Delete a like
// @Description Delete a like
// @Tags Like
// @Produce json
// @Param photoID path int true "Photo ID"
// @Security BearerToken
// @Success 200 {object} bool
// @Failure 400 {object} helper.ResponseError
// @Failure 401 {object} helper.ResponseError
// @Failure 404 {object} helper.ResponseError
// @Failure 500 {object} helper.ResponseError
// @Router /photos/{photoID}/likes [delete]
func (c *likeController) Delete(w http.ResponseWriter, r *http.Request) {
	var resp = response.New[interface{}](response.LikeDelete)

	photoIDStr := r.PathValue("photoID")
	photoID, err := strconv.ParseUint(photoIDStr, 10, 64)
	if err != nil {
		resp.Error(err).Code(http.StatusBadRequest).Send(w)
		return
	}

	err = c.likeService.Delete(r.Context(), photoID)
	if err != nil {
		respErr := new(helper.ResponseError)
		if errors.As(err, &respErr) {
			resp.Error(respErr).Code(respErr.Code()).Send(w)
			return
		}
		resp.Error(err).Code(http.StatusInternalServerError).Send(w)
		return
	}

	resp.Success(true).Code(http.StatusOK).Send(w)
}

// LikeGetMine godoc
// @Summary Get list of photo that user liked
// @Description Get list of photo that user liked
// @Tags Like
// @Produce json
// @Security BearerToken
// @Success 200 {object} response.Response[[]dto.PhotoResponse]
// @Failure 400 {object} response.Response[any]
// @Failure 401 {object} response.Response[any]
// @Failure 500 {object} response.Response[any]
// @Router /likes/my [get]
func (c *likeController) GetMine(w http.ResponseWriter, r *http.Request) {
	var (
		resp = response.New[[]dto.GetLikeByUserIDResponse](response.LikeGetMine)
		err  error
	)

	userID, ok := r.Context().Value(helper.UserIDKey).(float64)
	if !ok {
		resp.Error(helper.ErrInternal).Code(http.StatusInternalServerError).Send(w)
		return
	}

	photos, err := c.likeService.GetByUserID(r.Context(), uint64(userID))
	if err != nil {
		resp.Error(err).Code(http.StatusInternalServerError).Send(w)
		return
	}

	resp.Success(true).Data(photos).Code(http.StatusOK).Send(w)
}
