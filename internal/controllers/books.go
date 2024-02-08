package controllers

/*
import (
	"net/http"

	"github.com/leokuzmanovic/ai-bistro-chef/api"
	di "github.com/leokuzmanovic/ai-bistro-chef/internal/dependencyinjection"
	"github.com/leokuzmanovic/ai-bistro-chef/internal/services"
	"github.com/pkg/errors"

	"github.com/labstack/echo/v4"
)
*/
/*
type BooksController struct {
	booksService services.BooksService
}

func wireBooks(e *echo.Echo) {
	controller := &BooksController{
		booksService: di.Get[services.BooksService](),
	}

	e.POST(api.ENDPOINT_BOOKS, withValidJsonBody(controller.createNew), getMiddlewareFunction(constants.TimeoutRegular))
	e.GET(api.ENDPOINT_BOOKS_BY_ID, controller.getById, getMiddlewareFunction(constants.TimeoutQuick))
	e.DELETE(api.ENDPOINT_BOOKS_BY_ID, controller.delete, getMiddlewareFunction(constants.TimeoutRegular))
}

/*
// @Summary Create new book
// @Tags books
// @Router /books [POST]
// @Param body body api.CreateBookRequest true "CreateBookRequest"
// @Success 201
// @Failure 400 {object} errors.AppError "BadRequestError, InvalidJsonError, InvalidCredentialsError"
// @Failure 403 {object} errors.AppError "AuthorizationError"
// @Failure 404 {object} errors.AppError "NotFoundError"
*/
/*
func (s *BooksController) createNew(ctx echo.Context, body *api.CreateBookRequest) error {
	book, err := s.booksService.CreateBook(ctx.Request().Context(), body.Title, body.Author)
	if err != nil {
		return errors.Wrap(err, "books")
	}

	response := api.BookResponse{
		Id:        book.Id,
		CreatedAt: book.CreatedAt.String(),
		Title:     book.Title,
		Author:    book.Author,
	}
	return errors.Wrap(ctx.JSON(http.StatusCreated, response), "json")
}
*/

/*
// @Summary Get book by id
// @Tags books
// @Router /books/{bookId} [GET]
// @Param bookId path string true "book id"
// @Produce json
// @Success 200 {object} api.BookResponse
// @Success 400 {object} errors.AppError "BadRequestError"
// @Success 403 {object} errors.AppError "AuthorizationError"
// @Success 404 {object} errors.AppError "NotFoundError"
*/
/*
func (s *BooksController) getById(ctx echo.Context) error {
	bookId := ctx.Param("bookId")
	if bookId == "" {
		return &errs.BadRequestError{}
	}

	book, err := s.booksService.GetBookById(ctx.Request().Context(), bookId)
	if err != nil {
		return errors.Wrap(err, "books")
	}

	response := api.BookResponse{
		Id:        book.Id,
		CreatedAt: book.CreatedAt.String(),
		Title:     book.Title,
		Author:    book.Author,
	}
	return errors.Wrap(ctx.JSON(http.StatusOK, response), "json")
}
*/

/*
// @Summary Delete book
// @Tags books
// @Router /books/{bookId} [DELETE]
// @Param bookId path string true "book id"
// @Success 200
// @Failure 400 {object} errors.AppError "BadRequestError"
// @Failure 403 {object} errors.AppError "AuthorizationError"
// @Failure 404 {object} errors.AppError "NotFoundError"
*/
/*
func (s *BooksController) delete(ctx echo.Context) error {
	bookId := ctx.Param("bookId")
	if bookId == "" {
		return &errs.BadRequestError{}
	}

	err := s.booksService.DeleteBookById(ctx.Request().Context(), bookId)
	if err != nil {
		return errors.Wrap(err, "books")
	}

	return errors.Wrap(ctx.NoContent(http.StatusOK), "context")
}
*/
