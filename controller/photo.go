package controller

import (
	"encoding/json"
	"errors"
	"final-project/dto"
	"final-project/helper"
	"final-project/helper/response"
	"final-project/service"
	"net/http"
	"strconv"
)

type photoController struct {
	photoService service.PhotoService
}

func NewPhotoController(photoService service.PhotoService) *photoController {
	return &photoController{photoService}
}

// PhotoCreate godoc
// @Summary create a new photo
// @Tags Photo
// @Accept json
// @Produce json
// @Security BearerToken
// @Param request body dto.PhotoCreate true "required body"
// @Success 201 {object} response.Response[dto.PhotoCreateResponse]
// @Failure 400 {object} response.Response[any]
// @Failure 401 {object} response.Response[any]
// @Failure 500 {object} response.Response[any]
// @Router /photos [post]
func (c *photoController) Create(w http.ResponseWriter, r *http.Request) {
	var (
		data dto.PhotoRequest
		resp = response.New[dto.PhotoCreateResponse](response.PhotoCreate)
	)

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		resp.Error(err).Code(http.StatusBadRequest).Send(w)
		return
	}

	err = data.ValidateCreate()
	if err != nil {
		resp.Error(err).Code(http.StatusBadRequest).Send(w)
		return
	}

	photo, err := c.photoService.Create(r.Context(), data)
	if err != nil {
		respErr := new(helper.ResponseError)
		if errors.As(err, &respErr) {
			resp.Error(respErr).Code(respErr.Code()).Send(w)
			return
		}
		resp.Error(err).Code(http.StatusInternalServerError).Send(w)
		return
	}

	resp.Success(true).Data(photo).Code(http.StatusCreated).Send(w)
}

// PhotoGetAll godoc
// @Summary get all photos
// @Tags Photo
// @Accept json
// @Produce json
// @Security BearerToken
// @Success 200 {object} response.Response[[]dto.PhotoResponse]
// @Failure 401 {object} response.Response[any]
// @Failure 500 {object} response.Response[any]
// @Router /photos [get]
func (c *photoController) GetAll(w http.ResponseWriter, r *http.Request) {
	var resp = response.New[[]dto.PhotoResponse](response.PhotoGetAll)

	photos, err := c.photoService.GetAll(r.Context())
	if err != nil {
		respErr := new(helper.ResponseError)
		if errors.As(err, &respErr) {
			resp.Error(respErr).Code(respErr.Code()).Send(w)
			return
		}
		resp.Error(err).Code(http.StatusInternalServerError).Send(w)
		return
	}

	resp.Data(photos).Success(true).Code(http.StatusOK).Send(w)
}

// PhotoUpdate godoc
// @Summary update a photo
// @Tags Photo
// @Accept json
// @Produce json
// @Security BearerToken
// @Param photoID path int true "photo id"
// @Param request body dto.PhotoUpdate true "required body"
// @Success 200 {object} response.Response[dto.PhotoUpdateResponse]
// @Failure 400 {object} response.Response[any]
// @Failure 401 {object} response.Response[any]
// @Failure 404 {object} response.Response[any]
// @Failure 409 {object} response.Response[any]
// @Failure 500 {object} response.Response[any]
// @Router /photos/{photoID} [put]
func (c *photoController) Update(w http.ResponseWriter, r *http.Request) {
	var (
		data dto.PhotoRequest
		resp = response.New[dto.PhotoUpdateResponse](response.PhotoUpdate)
	)

	photoIDStr := r.PathValue("photoID")
	photoID, err := strconv.ParseUint(photoIDStr, 10, 64)
	if err != nil {
		resp.Error(helper.ErrInvalidID).Code(http.StatusBadRequest).Send(w)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		resp.Error(err).Code(http.StatusBadRequest).Send(w)
		return
	}

	err = data.ValidateUpdate()
	if err != nil {
		resp.Error(err).Code(http.StatusBadRequest).Send(w)
		return
	}

	photo, err := c.photoService.Update(r.Context(), photoID, data)
	if err != nil {
		respErr := new(helper.ResponseError)
		if errors.As(err, &respErr) {
			resp.Error(respErr).Code(respErr.Code()).Send(w)
			return
		}
		resp.Error(err).Code(http.StatusInternalServerError).Send(w)
		return
	}

	resp.Success(true).Data(photo).Code(http.StatusOK).Send(w)
}

// PhotoDelete godoc
// @Summary delete a photo
// @Tags Photo
// @Produce json
// @Security BearerToken
// @Param photoID path int true "photo id"
// @Success 200 {object} response.Response[any]
// @Failure 400 {object} response.Response[any]
// @Failure 401 {object} response.Response[any]
// @Failure 404 {object} response.Response[any]
// @Failure 500 {object} response.Response[any]
// @Router /photos/{photoID} [delete]
func (c *photoController) Delete(w http.ResponseWriter, r *http.Request) {
	var resp = response.New[any](response.PhotoDelete)

	photoIDStr := r.PathValue("photoID")
	photoID, err := strconv.ParseUint(photoIDStr, 10, 64)
	if err != nil {
		resp.Error(helper.ErrInvalidID).Code(http.StatusBadRequest).Send(w)
		return
	}

	err = c.photoService.Delete(r.Context(), photoID)
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

// PhotoGetByID godoc
// @Summary get a photo by id
// @Tags Photo
// @Produce json
// @Security BearerToken
// @Param photoID path int true "photo id"
// @Success 200 {object} response.Response[dto.PhotoResponse]
// @Failure 400 {object} response.Response[any]
// @Failure 401 {object} response.Response[any]
// @Failure 404 {object} response.Response[any]
// @Failure 500 {object} response.Response[any]
// @Router /photos/{photoID} [get]
func (c *photoController) GetByID(w http.ResponseWriter, r *http.Request) {
	var resp = response.New[dto.PhotoResponse](response.PhotoGetByID)

	photoIDStr := r.PathValue("photoID")
	photoID, err := strconv.ParseUint(photoIDStr, 10, 64)
	if err != nil {
		resp.Error(helper.ErrInvalidID).Code(http.StatusBadRequest).Send(w)
		return
	}

	photo, err := c.photoService.GetByID(r.Context(), photoID)
	if err != nil {
		respErr := new(helper.ResponseError)
		if errors.As(err, &respErr) {
			resp.Error(respErr).Code(respErr.Code()).Send(w)
			return
		}
		resp.Error(err).Code(http.StatusInternalServerError).Send(w)
		return
	}

	resp.Data(photo).Success(true).Code(http.StatusOK).Send(w)
}

// PhotoGetMine godoc
// @Summary get current user's photos
// @Tags Photo
// @Produce json
// @Security BearerToken
// @Success 200 {object} response.Response[[]dto.PhotoResponse]
// @Failure 401 {object} response.Response[any]
// @Failure 500 {object} response.Response[any]
// @Router /photos/my [get]
func (c *photoController) GetMine(w http.ResponseWriter, r *http.Request) {
	var resp = response.New[[]dto.PhotoResponse](response.PhotoGetMine)

	userID, ok := r.Context().Value(helper.UserIDKey).(float64)
	if !ok {
		resp.Error(helper.ErrInternal).Code(http.StatusInternalServerError).Send(w)
		return
	}

	photos, err := c.photoService.GetByUserID(r.Context(), uint64(userID))
	if err != nil {
		respErr := new(helper.ResponseError)
		if errors.As(err, &respErr) {
			resp.Error(respErr).Code(respErr.Code()).Send(w)
			return
		}
		resp.Error(err).Code(http.StatusInternalServerError).Send(w)
		return
	}

	resp.Data(photos).Success(true).Code(http.StatusOK).Send(w)
}

// PhotoGetByUsername godoc
// @Summary get all photos by username
// @Tags Photo
// @Produce json
// @Security BearerToken
// @Param username path string true "username"
// @Success 200 {object} response.Response[[]dto.PhotoResponse]
// @Failure 400 {object} response.Response[any]
// @Failure 401 {object} response.Response[any]
// @Failure 404 {object} response.Response[any]
// @Failure 500 {object} response.Response[any]
// @Router /users/{username}/photos [get]
func (c *photoController) GetByUsername(w http.ResponseWriter, r *http.Request) {
	var resp = response.New[[]dto.PhotoResponse](response.PhotoGetByUsername)

	username := r.PathValue("username")

	photos, err := c.photoService.GetByUsername(r.Context(), username)
	if err != nil {
		respErr := new(helper.ResponseError)
		if errors.As(err, &respErr) {
			resp.Error(respErr).Code(respErr.Code()).Send(w)
			return
		}
		resp.Error(err).Code(http.StatusInternalServerError).Send(w)
		return
	}

	resp.Data(photos).Success(true).Code(http.StatusOK).Send(w)
}
