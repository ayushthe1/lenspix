package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/ayushthe1/lenspix/context"
	"github.com/ayushthe1/lenspix/errors"
	"github.com/ayushthe1/lenspix/models"
	"github.com/go-chi/chi/v5"
)

type Galleries struct {
	Templates struct {
		// New template will be used to render the form that will be used to create a gallery
		New Template
		// Edit template will be used to render a page for editing the gallery
		Edit  Template
		Index Template
		Show  Template
	}
	// This will be used to process to that form
	GalleryService *models.GalleryService
}

func (g Galleries) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title string
	}

	data.Title = r.FormValue("title") // get the title for the url query-parameter if it's available or otherwise empty string
	g.Templates.New.Execute(w, r, data)
}

// Method(handler) to create a gallery
func (g Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var data struct {
		UserID int
		Title  string
	}
	// get the userID from the request's context
	data.UserID = context.User(r.Context()).ID
	// Title from the form coming in as post request
	data.Title = r.FormValue("title")

	gallery, err := g.GalleryService.Create(data.Title, data.UserID)
	if err != nil {
		// render the create new gallery form in case of any error
		g.Templates.New.Execute(w, r, data, err)
		return
	}
	// redirect the user to edit gallery page
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)

}

// Method(handler) to render the page(form) to edit gallery
func (g Galleries) Edit(w http.ResponseWriter, r *http.Request) {

	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}

	var data struct {
		ID    int
		Title string
	}
	data.ID = gallery.ID
	data.Title = gallery.Title

	g.Templates.Edit.Execute(w, r, data)

}

// Method(handler) to process the edit gallery form (once the update button is clicked)
func (g Galleries) Update(w http.ResponseWriter, r *http.Request) {

	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}

	gallery.Title = r.FormValue("title") // title value from form
	// update the gallery in db
	err = g.GalleryService.Update(gallery)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	// after the gallery is updated, redirect the user to the edit page
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

// handler to render index page where we'll see all the galleries that a user has access to
func (g Galleries) Index(w http.ResponseWriter, r *http.Request) {
	// take the current user ,look up at all of their galleries and then send them to the template to render them

	type Gallery struct {
		ID    int
		Title string
	}

	var data struct {
		Galleries []Gallery
	}

	user := context.User(r.Context())
	// query for all of the galleries a specific user owns
	galleries, err := g.GalleryService.ByUserID(user.ID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// convert the galleries returned from db into Gallery type which we can send to templates to render them
	for _, gallery := range galleries {
		data.Galleries = append(data.Galleries, Gallery{
			ID:    gallery.ID,
			Title: gallery.Title,
		})
	}

	g.Templates.Index.Execute(w, r, data)

}

// Handler function for deleting a gallery
func (g Galleries) Delete(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}

	err = g.GalleryService.Delete(gallery.ID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// Redirect to the galleries page
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

// Handler for showing a gallery. Anyone with a link to a gallery will be able to view it as we'll not restrict access to this page like we have done with other gallery pages
func (g Galleries) Show(w http.ResponseWriter, r *http.Request) {

	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	var data struct {
		ID     int
		Title  string
		Images []string
	}
	data.ID = gallery.ID
	data.Title = gallery.Title

	for i := 0; i < 20; i++ {
		// width and height are random values betwee 200 and 700
		w, h := rand.Intn(500)+200, rand.Intn(500)+200
		// using the width and height, we generate a URL
		catImageURL := fmt.Sprintf("http://placekitten.com/%d/%d", w, h)
		// Then we add the URL to our images.
		data.Images = append(data.Images, catImageURL)
	}

	// TODO: Render the gallery
	g.Templates.Show.Execute(w, r, data)
}

type galleryOpt func(http.ResponseWriter, *http.Request, *models.Gallery) error

// helper function to get the ID from the URL param, and then lookup the gallery.
// returns the gallery and the error
func (g Galleries) galleryByID(w http.ResponseWriter, r *http.Request, opts ...galleryOpt) (*models.Gallery, error) {
	// get the gallery id from the url query parameter
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusNotFound)
		return nil, err
	}
	// query for the gallery with the valid id
	gallery, err := g.GalleryService.ByID(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Gallery not found", http.StatusNotFound)
			return nil, err
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return nil, err
	}

	// for loop to iterate over all of our functional
	// options, calling each and returning if there is an error.
	for _, opt := range opts {
		err = opt(w, r, gallery)
		if err != nil {
			return nil, err
		}
	}

	return gallery, nil
}

func userMustOwnGallery(w http.ResponseWriter, r *http.Request, gallery *models.Gallery) error {
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "You are not authorized to edit this gallery", http.StatusForbidden)
		return fmt.Errorf("user doesn't have access to this gallery")
	}

	return nil
}
