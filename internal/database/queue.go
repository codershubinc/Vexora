package database

import (
	"database/sql"
	"time"
)

// QueueItem represents a single job in the system
type QueueItem struct {
	ID                   int64
	ProjectName          string
	RawNotes             string
	Status               string // PENDING, QUEUED, PROCESSING, WAITING_APPROVAL, APPROVED, PENDING-RETRY, FAILED
	Priority             string // NORMAL, HIGH
	AttemptCount         int
	CreatedAt            string
	GeneratedSubject     string
	GeneratedContent     string
	GeneratedTags        string
	ApprovalToken        string
	LastNotificationSent time.Time
	ErrorMsg             string
}

// --- Fetchers ---

// GetPendingIDs fetches IDs of jobs ready to be processed (Lane 1 -> Lane 2)
func GetPendingIDs() ([]int64, error) {
	rows, err := DB.Query("SELECT id FROM journal_entries WHERE status = 'PENDING'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// GetEntry fetches a single job by ID (Used by Worker)
func GetEntry(id int64) (*QueueItem, error) {
	var item QueueItem

	err := DB.QueryRow(`
		SELECT id, project_name, raw_notes, status, priority, attempt_count, created_at 
		FROM journal_entries WHERE id = ?`, id).
		Scan(&item.ID, &item.ProjectName, &item.RawNotes, &item.Status, &item.Priority, &item.AttemptCount, &item.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &item, nil
}

// GetJobsByStatus fetches full details for the Dashboard or Notifier
func GetJobsByStatus(status string) ([]QueueItem, error) {
	rows, err := DB.Query(`
		SELECT id, project_name, status, created_at, generated_subject, generated_content, approval_token 
		FROM journal_entries 
		WHERE status = ? 
		ORDER BY id DESC`, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []QueueItem
	for rows.Next() {
		var i QueueItem
		var subject, content, token sql.NullString // Handle NULLs safely

		if err := rows.Scan(&i.ID, &i.ProjectName, &i.Status, &i.CreatedAt, &subject, &content, &token); err != nil {
			return nil, err
		}
		i.GeneratedSubject = subject.String
		i.GeneratedContent = content.String
		i.ApprovalToken = token.String
		items = append(items, i)
	}
	return items, nil
}

// GetPendingRetryIDs fetches jobs that failed and need human confirmation
func GetPendingRetryIDs() ([]int64, error) {
	rows, err := DB.Query("SELECT id FROM journal_entries WHERE status = 'PENDING-RETRY'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// --- State Modifiers ---

// UpdateStatus moves a job to a new state
func UpdateStatus(id int64, status string) error {
	_, err := DB.Exec("UPDATE journal_entries SET status = ? WHERE id = ?", status, id)
	return err
}

// SetApprovalWait saves the LLM result and generates a token for human review
func SetApprovalWait(id int64, subject, content, tags, token string) error {
	_, err := DB.Exec(`
		UPDATE journal_entries 
		SET status = 'WAITING_APPROVAL', 
		    generated_subject = ?, 
		    generated_content = ?, 
		    generated_tags = ?, 
		    approval_token = ? 
		WHERE id = ?`,
		subject, content, tags, token, id)
	return err
}

// MarkRetry increments retry count and sets status to PENDING-RETRY (Human intervention needed)
func MarkRetry(id int64) error {
	_, err := DB.Exec(`
		UPDATE journal_entries 
		SET status = 'PENDING-RETRY', 
		    attempt_count = attempt_count + 1 
		WHERE id = ?`, id)
	return err
}

// ResetRetryCount is called when a human manually confirms a retry
func ResetRetryCount(id int64) error {
	_, err := DB.Exec("UPDATE journal_entries SET attempt_count = 0 WHERE id = ?", id)
	return err
}

// --- Auth / Security ---

// GetToken fetches the secure token to verify approval links
func GetToken(idStr string) (string, error) {
	var token string
	err := DB.QueryRow("SELECT approval_token FROM journal_entries WHERE id = ?", idStr).Scan(&token)
	if err != nil {
		return "", err
	}
	return token, nil
}
