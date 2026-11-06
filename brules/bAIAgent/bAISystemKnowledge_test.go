package baiagent

import (
	"regexp"
	"strings"
	"testing"
)

var aiPathPattern = regexp.MustCompile(`\[(/[^\]\s]+)\]`)

// TestAIKnowledgePaths fails when the knowledge links to a screen that is not
// in aiScreenPaths (the list that must match p3 assistantRoutes.js).
func TestAIKnowledgePaths(t *testing.T) {
	known := map[string]bool{}
	for _, path := range aiScreenPaths {
		known[path] = true
	}

	found := 0
	for _, match := range aiPathPattern.FindAllStringSubmatch(aiSystemKnowledge+aiAnswerRules, -1) {
		found++
		if !known[match[1]] {
			t.Errorf("the AI knowledge uses the unknown path [%s]; add it to aiScreenPaths and p3 assistantRoutes.js, or fix it", match[1])
		}
	}
	if found == 0 {
		t.Fatal("no [/path] found in the AI knowledge")
	}
}

// TestAIKnowledgeVersion makes sure the version is set and sent to Claude.
func TestAIKnowledgeVersion(t *testing.T) {
	if strings.TrimSpace(aiKnowledgeVersion) == "" {
		t.Fatal("aiKnowledgeVersion is empty")
	}
	prompt := buildSystemPrompt(AIRequestContext{Platform: "shao"})
	if !strings.Contains(prompt, "Knowledge version: "+aiKnowledgeVersion) {
		t.Fatal("the system prompt does not include the knowledge version")
	}
}
