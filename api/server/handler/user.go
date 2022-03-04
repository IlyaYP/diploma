package handler

import (
	"errors"
	"github.com/IlyaYP/diploma/model"
	"github.com/IlyaYP/diploma/pkg"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log"
	"net/http"
)

func (h *Handler) user(router chi.Router) {
	router.Post("/register", h.UserRegister)
	router.Post("/login", h.UserLogin)
}

// UserRegister register new user
func (h *Handler) UserRegister(w http.ResponseWriter, r *http.Request) {
	input := &model.User{}
	if err := render.Bind(r, input); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	user, err := h.userSvc.Register(r.Context(), input.Login, input.Password)
	if err != nil {
		//log.Println(err)
		if errors.Is(err, pkg.ErrAlreadyExists) {
			render.Render(w, r, ErrAlreadyExists)
			return
		}
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Render(w, r, &user) // TODO: Remove

}

// UserLogin authenticates user
func (h *Handler) UserLogin(w http.ResponseWriter, r *http.Request) {
	input := &model.User{}
	if err := render.Bind(r, input); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	user, err := h.userSvc.Login(r.Context(), input.Login, input.Password)
	if err != nil {
		log.Println(err)
		if errors.Is(err, pkg.ErrInvalidLogin) {
			render.Render(w, r, ErrInvalidLogin)
			return
		}
	}

	render.Render(w, r, &user) // TODO: Remove

	/*


		data := &ArticleRequest{}
		if err := render.Bind(r, data); err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		article := data.Article
		dbNewArticle(article)

		render.Status(r, http.StatusCreated)
		render.Render(w, r, NewArticleResponse(article))

	*/
}

func (h *Handler) UserContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//var article *Article
		//var err error
		//
		//if articleID := chi.URLParam(r, "articleID"); articleID != "" {
		//	article, err = dbGetArticle(articleID)
		//} else if articleSlug := chi.URLParam(r, "articleSlug"); articleSlug != "" {
		//	article, err = dbGetArticleBySlug(articleSlug)
		//} else {
		//	render.Render(w, r, ErrNotFound)
		//	return
		//}
		//if err != nil {
		//	render.Render(w, r, ErrNotFound)
		//	return
		//}
		//
		//ctx := context.WithValue(r.Context(), "article", article)
		//next.ServeHTTP(w, r.WithContext(ctx))
		log.Println("UserContext")
		next.ServeHTTP(w, r)
	})
}
