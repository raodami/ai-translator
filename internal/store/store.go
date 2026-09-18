package store

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

type TranslationRecord struct {
	ID           string    `json:"id"`
	Text         string    `json:"text"`
	SourceLang   string    `json:"source_lang"`
	TargetLang   string    `json:"target_lang"`
	Translated   string    `json:"translated"`
	CreatedAt    time.Time `json:"created_at"`
}

type Term struct {
	ID        string    `json:"id"`
	Source    string    `json:"source"`
	Target    string    `json:"target"`
	LangPair  string    `json:"lang_pair"`
	CreatedAt time.Time `json:"created_at"`
}

type Store struct {
	db *sql.DB
}

func New(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	s := &Store{db: db}
	if err := s.initSchema(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS translations (
		id TEXT PRIMARY KEY,
		text TEXT,
		source_lang TEXT,
		target_lang TEXT,
		translated TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS terms (
		id TEXT PRIMARY KEY,
		source TEXT,
		target TEXT,
		lang_pair TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := s.db.Exec(schema)
	return err
}

func (s *Store) SaveTranslation(r *TranslationRecord) error {
	_, err := s.db.Exec(
		"INSERT INTO translations (id, text, source_lang, target_lang, translated, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		r.ID, r.Text, r.SourceLang, r.TargetLang, r.Translated, r.CreatedAt,
	)
	return err
}

func (s *Store) GetTranslations(limit int) ([]*TranslationRecord, error) {
	rows, err := s.db.Query("SELECT id, text, source_lang, target_lang, translated, created_at FROM translations ORDER BY created_at DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*TranslationRecord
	for rows.Next() {
		var r TranslationRecord
		if err := rows.Scan(&r.ID, &r.Text, &r.SourceLang, &r.TargetLang, &r.Translated, &r.CreatedAt); err != nil {
			continue
		}
		records = append(records, &r)
	}
	return records, rows.Err()
}

func (s *Store) SaveTerm(t *Term) error {
	_, err := s.db.Exec(
		"INSERT INTO terms (id, source, target, lang_pair, created_at) VALUES (?, ?, ?, ?, ?)",
		t.ID, t.Source, t.Target, t.LangPair, t.CreatedAt,
	)
	return err
}

func (s *Store) GetTerms(langPair string) ([]*Term, error) {
	var rows *sql.Rows
	var err error

	if langPair != "" {
		rows, err = s.db.Query("SELECT id, source, target, lang_pair, created_at FROM terms WHERE lang_pair = ? ORDER BY created_at DESC", langPair)
	} else {
		rows, err = s.db.Query("SELECT id, source, target, lang_pair, created_at FROM terms ORDER BY created_at DESC")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var terms []*Term
	for rows.Next() {
		var t Term
		if err := rows.Scan(&t.ID, &t.Source, &t.Target, &t.LangPair, &t.CreatedAt); err != nil {
			continue
		}
		terms = append(terms, &t)
	}
	return terms, rows.Err()
}

func (s *Store) DeleteTerm(id string) error {
	_, err := s.db.Exec("DELETE FROM terms WHERE id = ?", id)
	return err
}

func (s *Store) GetStats() (map[string]interface{}, error) {
	var totalTrans, totalTerms int
	s.db.QueryRow("SELECT COUNT(*) FROM translations").Scan(&totalTrans)
	s.db.QueryRow("SELECT COUNT(*) FROM terms").Scan(&totalTerms)
	return map[string]interface{}{
		"total_translations": totalTrans,
		"total_terms":        totalTerms,
	}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}
