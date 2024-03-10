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

type commentController struct {
	commentService service.CommentService
}

func NewCommentController(commentService service.CommentService) *commentController {
	return &commentController{commentService}
}

// CommentCreate godoc
// @Summary create a new comment
// @Tags Comment
// @Accept json
// @Produce json
// @Security BearerToken
// @Param request body dto.CommentCreate true "required body"
// @Param photoID path int true "photo ID"
// @Success 201 {object} response.Response[dto.CommentCreateResponse]
// @Failure 400 {object} response.Response[any]
// @Failure 401 {object} response.Response[any]
// @Failure 500 {object} response.Response[any]
// @Router /photos/{photoID}/comments [post]
func (c *commentController) Create(w http.ResponseWriter, r *http.Request) {
	var (
		data dto.CommentRequest
		resp = response.New[dto.CommentCreateResponse](response.CommentCreate)
	)
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		resp.Error(err).Code(http.StatusBadRequest).Send(w)
		return
	}

	photoIDStr := r.PathValue("photoID")
	data.PhotoID, err = strconv.ParseUint(photoIDStr, 10, 64)
	if err != nil {
		resp.Error(helper.ErrInvalidID).Code(http.StatusBadRequest).Send(w)
		return
	}

	err = data.ValidateCreate()
	if err != nil {
		resp.Error(err).Code(http.StatusBadRequest).Send(w)
		return
	}

	comment, err := c.commentService.Create(r.Context(), data)
	if err != nil {
		respErr := new(helper.ResponseError)
		if errors.As(err, &respErr) {
			resp.Error(respErr).Code(respErr.Code()).Send(w)
			return
		}
		resp.Error(err).Code(http.StatusInternalServerError).Send(w)
		return
	}

	resp.Success(true).Data(comment).Code(http.StatusCreated).Send(w)
}

// CommentGetAll godoc
// @Summary get all comments
// @Tags Comment
// @Produce json
// @Security BearerToken
// @Success 200 {object} response.Response[[]dto.CommentResponse]
// @Failure 401 {object} response.Response[any]
// @Failure 500 {object} response.Response[any]
// @Router /comments [get]
func (c *commentController) GetAll(w http.ResponseWriter, r *http.Request) {
	var resp = response.New[[]dto.CommentResponse](response.CommentGetAll)

	comments, err := c.commentService.GetAll(r.Context())
	if err != nil {
		respErr := new(helper.ResponseError)
		if errors.As(err, &respErr) {
			resp.Error(respErr).Code(respErr.Code()).Send(w)
			return
		}
		resp.Error(err).Code(http.StatusInternalServerError).Send(w)
		return
	}

	resp.Data(comments).Success(true).Code(http.StatusOK).Send(w)
}

// CommentUpdate godoc
// @Summary update a comment
// @Tags Comment
// @Accept json
// @Produce json
// @Security BearerToken
// @Param commentID path int true "comment ID"
// @Param request body dto.CommentUpdate true "required body"
// @Success 200 {object} response.Response[dto.CommentUpdateResponse]
// @Failure 400 {object} response.Response[any]
// @Failure 401 {object} response.Response[any]
// @Failure 404 {object} response.Response[any]
// @Failure 500 {object} response.Response[any]
// @Router /comments/{commentID} [put]
func (c *commentController) Update(w http.ResponseWriter, r *http.Request) {
	var (
		data dto.CommentRequest
		resp = response.New[dto.CommentUpdateResponse](response.CommentUpdate)
	)

	commentIDStr := r.PathValue("commentID")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 64)
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

	comment, err := c.commentService.Update(r.Context(), commentID, data)
	if err != nil {
		respErr := new(helper.ResponseError)
		if errors.As(err, &respErr) {
			resp.Error(respErr).Code(respErr.Code()).Send(w)
			return
		}
		resp.Error(err).Code(http.StatusInternalServerError).Send(w)
		return
	}

	resp.Success(true).Data(comment).Code(http.StatusOK).Send(w)
}

// CommentDelete godoc
// @Summary delete a comment
// @Tags Comment
// @Produce json
// @Security BearerToken
// @Param commentID path int true "comment ID"
// @Success 200 {object} response.Response[any]
// @Failure 400 {object} response.Response[any]
// @Failure 401 {object} response.Response[any]
// @Failure 404 {object} response.Response[any]
// @Failure 500 {object} response.Response[any]
// @Router /comments/{commentID} [delete]
func (c *commentController) Delete(w http.ResponseWriter, r *http.Request) {
	var resp = response.New[any](response.CommentDelete)

	commentIDStr := r.PathValue("commentID")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 64)
	if err != nil {
		resp.Error(helper.ErrInvalidID).Code(http.StatusBadRequest).Send(w)
		return
	}

	err = c.commentService.Delete(r.Context(), commentID)
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

// CommentGetByID godoc
// @Summary get a comment by ID
// @Tags Comment
// @Produce json
// @Security BearerToken
// @Param commentID path int true "comment ID"
// @Success 200 {object} response.Response[dto.CommentResponse]
// @Failure 400 {object} response.Response[any]
// @Failure 401 {object} response.Response[any]
// @Failure 404 {object} response.Response[any]
// @Failure 500 {object} response.Response[any]
// @Router /comments/{commentID} [get]
func (c *commentController) GetByID(w http.ResponseWriter, r *http.Request) {
	var resp = response.New[dto.CommentResponse](response.CommentGetByID)

	commentIDStr := r.PathValue("commentID")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 64)
	if err != nil {
		resp.Error(helper.ErrInvalidID).Code(http.StatusBadRequest).Send(w)
		return
	}

	comment, err := c.commentService.GetByID(r.Context(), commentID)
	if err != nil {
		respErr := new(helper.ResponseError)
		if errors.As(err, &respErr) {
			resp.Error(respErr).Code(respErr.Code()).Send(w)
			return
		}
		resp.Error(err).Code(http.StatusInternalServerError).Send(w)
		return
	}

	resp.Data(comment).Success(true).Code(http.StatusOK).Send(w)
}

// CommentGetByPhotoID godoc
// @Summary get all comments by photo ID
// @Tags Comment
// @Produce json
// @Security BearerToken
// @Param photoID path int true "photo ID"
// @Success 200 {object} response.Response[[]dto.CommentGetByPhotoIDResponse]
// @Failure 400 {object} response.Response[any]
// @Failure 401 {object} response.Response[any]
// @Failure 404 {object} response.Response[any]
// @Failure 500 {object} response.Response[any]
// @Router /photos/{photoID}/comments [get]
func (c *commentController) GetByPhotoID(w http.ResponseWriter, r *http.Request) {
	var (
		resp = response.New[[]dto.CommentGetByPhotoIDResponse](response.CommentGetByID)
		err  error
	)

	photoIDStr := r.PathValue("photoID")
	photoID, err := strconv.ParseUint(photoIDStr, 10, 64)
	if err != nil {
		resp.Error(err).Code(http.StatusBadRequest).Send(w)
		return
	}

	comments, err := c.commentService.GetByPhotoID(r.Context(), photoID)
	if err != nil {
		resp.Error(err).Code(http.StatusInternalServerError).Send(w)
		return
	}

	resp.Success(true).Data(comments).Code(http.StatusOK).Send(w)
}

// CommentGetMine godoc
// @Summary get current user's comments
// @Tags Comment
// @Produce json
// @Security BearerToken
// @Success 200 {object} response.Response[[]dto.CommentGetByUserIDResponse]
// @Failure 401 {object} response.Response[any]
// @Failure 500 {object} response.Response[any]
// @Router /comments/my [get]
func (c *commentController) GetMine(w http.ResponseWriter, r *http.Request) {
	var resp = response.New[[]dto.CommentGetByUserIDResponse](response.CommentGetMine)

	userID, ok := r.Context().Value(helper.UserIDKey).(float64)
	if !ok {
		resp.Error(helper.ErrInternal).Code(http.StatusInternalServerError).Send(w)
		return
	}

	comments, err := c.commentService.GetByUserID(r.Context(), uint64(userID))
	if err != nil {
		resp.Error(err).Code(http.StatusInternalServerError).Send(w)
		return
	}

	resp.Success(true).Data(comments).Code(http.StatusOK).Send(w)
}
