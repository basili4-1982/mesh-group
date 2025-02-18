package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"

	"sap_segmentation/internal/config"
	"sap_segmentation/internal/model"
)

type ImportService struct {
	cfg *config.Config
	db  *sqlx.DB
}

func NewImportService(cfg *config.Config, db *sqlx.DB) *ImportService {
	return &ImportService{cfg: cfg, db: db}
}

func (s *ImportService) ImportData() error {
	offset := 0
	for {
		url := fmt.Sprintf("%s?p_limit=%d&p_offset=%d", s.cfg.ConnURI, s.cfg.ImportBatchSize, offset)

		log.Println(url)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}

		req.Header.Set("User-Agent", s.cfg.ConnUserAgent)
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(s.cfg.ConnAuthLoginPwd)))

		client := &http.Client{Timeout: s.cfg.ConnTimeout}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if len(body) == 0 {
			break
		}

		var segmentations []model.Segmentation
		if err := json.Unmarshal(body, &segmentations); err != nil {
			return err
		}

		for _, seg := range segmentations {
			if err := seg.Insert(s.db); err != nil {
				return err
			}
		}

		offset += s.cfg.ImportBatchSize
		time.Sleep(s.cfg.ConnInterval)
	}

	return nil
}
