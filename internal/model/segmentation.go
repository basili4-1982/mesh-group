package model

import (
	"github.com/jmoiron/sqlx"
)

type Segmentation struct {
	ID           int    `db:"id"`
	AddressSapID string `db:"address_sap_id"`
	AdrSegment   string `db:"adr_segment"`
	SegmentID    int64  `db:"segment_id"`
}

func (s *Segmentation) Insert(db *sqlx.DB) error {
	query := `
		INSERT INTO segmentation (address_sap_id, adr_segment, segment_id)
		VALUES (:address_sap_id, :adr_segment, :segment_id)
		ON CONFLICT (address_sap_id) DO UPDATE
		SET adr_segment = EXCLUDED.adr_segment, segment_id = EXCLUDED.segment_id
	`
	_, err := db.NamedExec(query, s)
	return err
}
