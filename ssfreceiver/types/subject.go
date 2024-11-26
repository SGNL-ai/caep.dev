package types

import "github.com/sgnl-ai/caep.dev-receiver/secevent/pkg/subject"

const (
	DefaultSubjectsAll  = "ALL"
	DefaultSubjectsNone = "NONE"
)

type StreamSubjectRequest struct {
	StreamID string          `json:"stream_id"`
	Subject  subject.Subject `json:"subject"`
	Verified bool            `json:"verified,omitempty"`
}

func (r *StreamSubjectRequest) ValidateRequest() error {
	if r.StreamID == "" {
		return NewError(
			ErrInvalidConfiguration,
			"ValidateSubjectRequest",
			"stream_id is required",
		)
	}

	if r.Subject == nil {
		return NewError(
			ErrInvalidConfiguration,
			"ValidateSubjectRequest",
			"subject is required",
		)
	}

	if err := r.Subject.Validate(); err != nil {
		return NewError(
			err,
			"ValidateSubjectRequest",
			"invalid subject",
		)
	}

	return nil
}
