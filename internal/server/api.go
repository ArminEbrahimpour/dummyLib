package server

import "Library/pkg/db"

type API struct {
	Store *db.Store
}

func NewAPI(s *db.Store) *API {
	return &API{Store: s}
}

func (a *API) RegisterRoutes()
