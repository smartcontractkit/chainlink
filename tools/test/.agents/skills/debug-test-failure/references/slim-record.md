<slim-record>
```json
{
  "jira_key":            "KEY-NNN",
  "title":               "string",
  "description":         "string",
  "package":             "github.com/owner/repo/path | null",
  "test_name":           "TestFoo | TestFoo/subtest_name",
  "previous_attempts":   [
    {
      "outcome":               "MISMATCH | ABANDONED | FIXED",
      "date":                  "YYYY-MM-DD",
      "summary":               "string",
      "excluded_approaches":   ["string"],
      "rejection_reasons":     ["string"],
      "recommended_next_step": "string | null",
      "full_text":             "string"
    }
  ],
  "original_assignee":  "string | null",
  "skip_reason":        "string | null"
}
```

<field-rules>
- `package`: `customfield_13009`. Null if absent.
- `test_name`: `customfield_13007` (full path including subtest). If absent, longest `TestXxx`/`testXxx` token from title.
- `previous_attempts`: parsed per `investigation-comment` parsing rules.
- `skip_reason`: set by subagent, when claiming.
</field-rules>
</slim-record>