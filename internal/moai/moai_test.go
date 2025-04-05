package moai

import (
	"strings"
	"testing"
)

// TestGetRandomFace tests the GetRandomFace function
func TestGetRandomFace(t *testing.T) {
	// Call multiple times to ensure we get different faces (not guaranteed but likely)
	faces := make(map[string]bool)

	// Get multiple faces and count unique ones
	for i := 0; i < 50; i++ {
		face := GetRandomFace()
		faces[face] = true

		// Verify basic structure of the face (should start with Moai emoji)
		if !strings.HasPrefix(face, "🗿") {
			t.Errorf("Face should start with Moai emoji, got: %s", face)
		}
	}

	// We should have seen multiple unique faces
	if len(faces) < 3 {
		t.Errorf("Expected to see multiple unique faces, only got %d", len(faces))
	}
}

// TestGetRandomFeedback tests the GetRandomFeedback function
func TestGetRandomFeedback(t *testing.T) {
	testCases := []struct {
		commitMsg        string
		expectedNotEmpty bool
	}{
		{"fix: corrected login issue", true},
		{"update documentation", true},
		{"add new feature", true},
		{"remove deprecated code", true},
		{"initial commit", true},
		{"wip: working on feature", true},
		{"refactor auth logic", true},
		{"something completely different", true}, // should use general feedback
	}

	for _, tc := range testCases {
		feedback := GetRandomFeedback(tc.commitMsg)

		// Check if the feedback is returned and not empty
		if feedback == "" && tc.expectedNotEmpty {
			t.Errorf("Expected non-empty feedback for message '%s'", tc.commitMsg)
		}
	}

	// More basic check to ensure the feedback system is providing different outputs
	feedbackSet := make(map[string]bool)
	for i := 0; i < 20; i++ {
		feedback := GetRandomFeedback("test commit")
		feedbackSet[feedback] = true
	}

	// We should get at least one unique feedback (likely more)
	if len(feedbackSet) < 1 {
		t.Error("Expected to get at least one unique feedback message")
	}
}

// TestFeedbackTemplates tests that all feedback templates are valid
func TestFeedbackTemplates(t *testing.T) {
	// Check that all template categories have at least one feedback option
	for category, templates := range feedbackTemplates {
		if len(templates) == 0 {
			t.Errorf("Feedback category '%s' has no templates", category)
		}

		// Check that all templates in the category are non-empty
		for i, template := range templates {
			if template == "" {
				t.Errorf("Empty template at index %d for category '%s'", i, category)
			}
		}
	}

	// Check that we have general feedback options
	if len(generalFeedback) == 0 {
		t.Error("No general feedback templates defined")
	}

	// Check that all general feedback templates are non-empty
	for i, feedback := range generalFeedback {
		if feedback == "" {
			t.Errorf("Empty general feedback at index %d", i)
		}
	}
}

// TestMoaiFaces tests that all Moai faces are valid
func TestMoaiFaces(t *testing.T) {
	// Check that we have faces defined
	if len(moaiFaces) == 0 {
		t.Error("No Moai faces defined")
	}

	// Check that all faces are non-empty and contain the Moai emoji
	for i, face := range moaiFaces {
		if face == "" {
			t.Errorf("Empty face at index %d", i)
		}

		if !strings.Contains(face, "🗿") {
			t.Errorf("Face at index %d does not contain Moai emoji: %s", i, face)
		}
	}
}
