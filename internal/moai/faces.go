// Package moai provides functions for generating Moai feedback
// with different personalities and moods.
package moai

import (
	"math/rand"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	// Random number generator with time-based seed
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))

	// moaiFaces is a collection of Moai ASCII faces with different expressions
	moaiFaces = []string{
		"🗿  (ಠ_ಠ)",         // Disapproving
		"🗿  (¬_¬)",         // Suspicious
		"🗿  (⊙_⊙)",         // Surprised
		"🗿  (¯\\_(:/)_/¯)", // Confused shrug
		"🗿  (╯°□°）╯",       // Angry table flip
		"🗿  (◉_◉)",         // Surprised or shocked
		"🗿  (⊙﹏⊙)",         // Worried
		"🗿  (⚆_⚆)",         // Staring
		"🗿  ( ͡° ͜ʖ ͡°)",   // Smug
		"🗿  (ᵔᴥᵔ)",         // Happy
		"🗿  (•‿•)",         // Pleased
		"🗿  (≧◡≦)",         // Very happy
		"🗿  (─‿‿─)",        // Satisfied
	}

	// Feedback templates based on commit message patterns
	feedbackTemplates = map[string][]string{
		"fix": {
			"Ah, the classic 'fix' commit. What exactly needed fixing?",
			"Another fix? Your code must be quite the troublemaker.",
			"In a parallel universe, this code worked the first time.",
		},
		"update": {
			"Updates are like birthdays - everybody gets one, but nobody remembers why.",
			"Updating code: the software equivalent of rearranging furniture.",
			"Updated! But will anyone notice?",
		},
		"add": {
			"Feature creep intensifies.",
			"Another feature? Your codebase is becoming quite the collection.",
			"More code, more problems. Or solutions. We'll find out which soon.",
		},
		"remove": {
			"Deleting code feels better than writing it, doesn't it?",
			"Less is more. Unless it's test coverage.",
			"The best code is no code at all. You're getting closer!",
		},
		"initial": {
			"Every journey begins with a commit.",
			"Ah, the optimism of a fresh start. Cherish it while it lasts.",
			"The first commit is always the most innocent.",
		},
		"wip": {
			"WIP: Wisely Incomplete Progress.",
			"Work In Progress... or Wishful Thinking?",
			"Halfway there, or halfway to realizing you need to start over?",
		},
		"refactor": {
			"Refactoring: the art of moving furniture around while telling yourself it's cleaner.",
			"Your future self thanks you. Your current teammates curse you.",
			"Same same, but different. But still same.",
		},
	}

	// General feedback for when no specific pattern matches
	generalFeedback = []string{
		"I have no idea what that commit does. But then again, I'm just a Moai.",
		"Your commit is beyond my stone-faced comprehension.",
		"That's certainly a commit that was made. I'm sure of it.",
		"Intriguing commit. Very... human of you.",
		"I've been standing on Easter Island for centuries, and I still don't get that commit.",
		"The ancient wisdom of the Moai has no guidance for this commit.",
		"You've entered the 2AM hotfix arc. A legendary time.",
		"Future you will both thank and curse present you.",
		"This code now bears your fingerprints. No takebacks.",
		"The git blame will remember that.",
	}
)

// GetRandomFace returns a random Moai face
func GetRandomFace() string {
	return moaiFaces[rng.Intn(len(moaiFaces))]
}

// GetRandomFeedback generates feedback based on the commit message
func GetRandomFeedback(commitMsg string) string {
	commitMsg = strings.ToLower(commitMsg)

	// Check for specific patterns in the commit message
	for pattern, templates := range feedbackTemplates {
		if strings.Contains(commitMsg, pattern) {
			return color.YellowString(templates[rng.Intn(len(templates))])
		}
	}

	// If no specific pattern matched, return general feedback
	return color.YellowString(generalFeedback[rng.Intn(len(generalFeedback))])
}
