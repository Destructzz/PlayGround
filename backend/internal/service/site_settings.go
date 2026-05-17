package service

import (
	"backend/internal/domain"
	"backend/internal/repo/sqlc"
	"backend/pkg"
	"context"

	"github.com/jackc/pgx/v5"
)

type SiteSettingsService struct {
	queries *sqlc.Queries
}

func NewSiteSettings(queries *sqlc.Queries) *SiteSettingsService {
	return &SiteSettingsService{queries: queries}
}

func (s *SiteSettingsService) Get(ctx context.Context) (domain.SiteSettings, error) {
	row, err := s.queries.GetSiteSettings(ctx)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.SiteSettings{
				ID:               1,
				SettingsJSON:     []byte("[]"),
				GalleryItemsJSON: []byte("[]"),
			}, nil
		}
		return domain.SiteSettings{}, err
	}
	return domain.SiteSettings{
		ID:               row.ID,
		SettingsJSON:     pkg.JSONOrArray(row.SettingsJson),
		GalleryItemsJSON: pkg.JSONOrArray(row.GalleryItemsJson),
	}, nil
}

func (s *SiteSettingsService) Upsert(ctx context.Context, dto domain.UpdateSiteSettingsRequest) (domain.SiteSettings, error) {
	settingsJSON := dto.SettingsJSON
	if len(settingsJSON) == 0 {
		settingsJSON = []byte("[]")
	}

	galleryItemsJSON := dto.GalleryItemsJSON
	if len(galleryItemsJSON) == 0 {
		galleryItemsJSON = []byte("[]")
	}

	row, err := s.queries.UpsertSiteSettings(ctx, sqlc.UpsertSiteSettingsParams{
		SettingsJson:     settingsJSON,
		GalleryItemsJson: galleryItemsJSON,
	})
	if err != nil {
		return domain.SiteSettings{}, err
	}

	return domain.SiteSettings{
		ID:               row.ID,
		SettingsJSON:     pkg.JSONOrArray(row.SettingsJson),
		GalleryItemsJSON: pkg.JSONOrArray(row.GalleryItemsJson),
	}, nil
}
