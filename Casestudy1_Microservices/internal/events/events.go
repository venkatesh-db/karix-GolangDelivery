package events

// ResumeUploaded represents the event for a resume submission.
type ResumeUploaded struct {
	StudentID string   `json:"student_id"`
	Name      string   `json:"name"`
	CGPA      float64  `json:"cgpa"`
	Branch    string   `json:"branch"`
	Skills    []string `json:"skills"`
}

// ShortlistedCandidate represents the event for a shortlisted candidate.
type ShortlistedCandidate struct {
	StudentID string  `json:"student_id"`
	Name      string  `json:"name"`
	CGPA      float64 `json:"cgpa"`
	Branch    string  `json:"branch"`
}

// InterviewScheduled represents the event for an interview schedule.
type InterviewScheduled struct {
	StudentID     string `json:"student_id"`
	Name          string `json:"name"`
	InterviewSlot string `json:"interview_slot"`
}
