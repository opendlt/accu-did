package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// TodoItem represents a single TODO/FIXME/etc finding
type TodoItem struct {
	Path          string `json:"path"`
	Line          int    `json:"line"`
	Tag           string `json:"tag"`
	Text          string `json:"text"`
	ContextBefore string `json:"contextBefore,omitempty"`
	ContextAfter  string `json:"contextAfter,omitempty"`
	GitSHAShort   string `json:"gitShaShort,omitempty"`
}

// TodoSummary represents aggregated statistics
type TodoSummary struct {
	CountsByTag  map[string]int `json:"countsByTag"`
	CountsByDir  map[string]int `json:"countsByDir"`
	CountsByFile map[string]int `json:"countsByFile"`
}

// TodoReport represents the complete scan report
type TodoReport struct {
	GeneratedAt time.Time   `json:"generatedAt"`
	RepoPath    string      `json:"repoPath"`
	GitSHAShort string      `json:"gitShaShort,omitempty"`
	Summary     TodoSummary `json:"summary"`
	Items       []TodoItem  `json:"items"`
	TotalCount  int         `json:"totalCount"`
}

var (
	// Patterns to search for technical debt markers - focus on actionable items
	todoPatterns = []string{
		// Comment-led markers (various comment styles)
		`(?i)(^|\s)(//|#|/\*|\*|--)\s*(TODO|FIXME|XXX|HACK|STUB|TBA|TBD|NOTIMPL|NOTIMPLEMENTED)\b`,
		// Line-leading markers (start of line)
		`(?i)^\s*(TODO|FIXME|XXX|HACK|STUB|TBA|TBD|NOTIMPL|NOTIMPLEMENTED)\b`,
		// PANIC("TODO") special case
		`(?i)PANIC\s*\(\s*["']TODO`,
	}

	// File extensions to include
	includeExtensions = []string{
		".go", ".md", ".yaml", ".yml", ".json", ".ps1", ".sh", ".mmd", ".txt",
	}

	// Special files to include (exact matches)
	includeSpecialFiles = []string{
		"Dockerfile", "Makefile",
	}

	// Directories to exclude
	excludeDirs = []string{
		".git", "dist", "node_modules", "bin", "vendor", ".idea", ".vscode",
		"reports", // prevent scanning its own generated reports
	}

	// Files to exclude (case-insensitive matching on base name and path suffix)
	excludeFiles = []string{
		"CLAUDE.md",
		"README.md",
		"PORTING.md",
		filepath.Join("sdks", "go", "accdid", "README.md"),
		"todo-report.json",
		"todo-report.md",
		"todo-report.csv",
		filepath.Join("reports", "todo-report.json"),
		filepath.Join("reports", "todo-report.md"),
		filepath.Join("reports", "todo-report.csv"),
	}

	// File extensions to exclude
	excludeExtensions = []string{
		".exe", ".tar.gz", ".zip", ".png", ".svg", ".ico", ".pdf", ".jpg", ".jpeg", ".gif",
	}
)

func main() {
	repoPath := "."
	if len(os.Args) > 1 {
		repoPath = os.Args[1]
	}

	fmt.Fprintf(os.Stderr, "Scanning repository: %s\n", repoPath)

	// Get git SHA if available
	gitSHA := getGitSHA(repoPath)

	// Scan the repository
	items, err := scanRepository(repoPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning repository: %v\n", err)
		os.Exit(1)
	}

	// Generate summary
	summary := generateSummary(items)

	// Create report
	report := TodoReport{
		GeneratedAt: time.Now().UTC(),
		RepoPath:    repoPath,
		GitSHAShort: gitSHA,
		Summary:     summary,
		Items:       items,
		TotalCount:  len(items),
	}

	// Ensure reports directory exists
	if err := os.MkdirAll("reports", 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating reports directory: %v\n", err)
		os.Exit(1)
	}

	// Generate outputs
	if err := generateJSONReport(report); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating JSON report: %v\n", err)
		os.Exit(1)
	}

	if err := generateMarkdownReport(report); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating Markdown report: %v\n", err)
		os.Exit(1)
	}

	if err := generateCSVReport(report); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating CSV report: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Scan complete. Found %d items.\n", len(items))
	fmt.Fprintf(os.Stderr, "Reports generated in ./reports/\n")
}

func getGitSHA(repoPath string) string {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func scanRepository(repoPath string) ([]TodoItem, error) {
	var items []TodoItem

	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories that should be excluded
		if info.IsDir() {
			for _, excludeDir := range excludeDirs {
				if info.Name() == excludeDir {
					return filepath.SkipDir
				}
			}
			return nil
		}

		// Check if file should be included
		if !shouldIncludeFile(path, info.Name()) {
			return nil
		}

		// Scan file for TODO items
		fileItems, err := scanFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: error scanning file %s: %v\n", path, err)
			return nil
		}

		items = append(items, fileItems...)
		return nil
	})

	return items, err
}

func shouldIncludeFile(path, filename string) bool {
	// Exclude specific files (case-insensitive matching)
	for _, exclude := range excludeFiles {
		// Match either by exact base name or relative path fragment (case-insensitive)
		if strings.EqualFold(filename, exclude) ||
		   strings.EqualFold(filename, filepath.Base(exclude)) ||
		   strings.HasSuffix(strings.ToLower(path), strings.ToLower(exclude)) {
			return false
		}
	}

	// Check special files first
	for _, special := range includeSpecialFiles {
		if filename == special || strings.HasPrefix(filename, special) {
			return true
		}
	}

	// Check extensions
	ext := filepath.Ext(filename)

	// Exclude certain extensions
	for _, excludeExt := range excludeExtensions {
		if strings.HasSuffix(filename, excludeExt) {
			return false
		}
	}

	// Include certain extensions
	for _, includeExt := range includeExtensions {
		if ext == includeExt {
			return true
		}
	}

	return false
}

func scanFile(path string) ([]TodoItem, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var items []TodoItem
	var lines []string
	scanner := bufio.NewScanner(file)
	lineNum := 0

	// Read all lines first for context
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Compile regex patterns
	var compiledPatterns []*regexp.Regexp
	for _, pattern := range todoPatterns {
		compiled, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("error compiling pattern %s: %v", pattern, err)
		}
		compiledPatterns = append(compiledPatterns, compiled)
	}

	// Scan each line
	for i, line := range lines {
		lineNum = i + 1

		for _, pattern := range compiledPatterns {
			matches := pattern.FindAllStringSubmatch(line, -1)
			for _, match := range matches {
				tag := extractTag(match[0])

				// Skip scanner self-references and documentation about scanning
				if isScannerSelfReference(strings.TrimSpace(line), path) {
					continue
				}

				item := TodoItem{
					Path: path,
					Line: lineNum,
					Tag:  tag,
					Text: strings.TrimSpace(line),
				}

				// Add context lines
				if i > 0 {
					item.ContextBefore = strings.TrimSpace(lines[i-1])
				}
				if i < len(lines)-1 {
					item.ContextAfter = strings.TrimSpace(lines[i+1])
				}

				items = append(items, item)
			}
		}
	}

	return items, nil
}

func extractTag(match string) string {
	// Extract the actual tag from the match
	match = strings.ToUpper(match)

	if strings.Contains(match, "TODO") {
		return "TODO"
	}
	if strings.Contains(match, "FIXME") {
		return "FIXME"
	}
	if strings.Contains(match, "XXX") {
		return "XXX"
	}
	if strings.Contains(match, "HACK") {
		return "HACK"
	}
	if strings.Contains(match, "STUB") {
		return "STUB"
	}
	if strings.Contains(match, "TBA") {
		return "TBA"
	}
	if strings.Contains(match, "TBD") {
		return "TBD"
	}
	if strings.Contains(match, "NOTIMPL") {
		return "NOTIMPLEMENTED"
	}
	if strings.Contains(match, "PANIC") {
		return "PANIC-TODO"
	}

	return "OTHER"
}

// isScannerSelfReference checks if a line is a scanner self-reference that should be excluded
func isScannerSelfReference(line, path string) bool {
	lineUpper := strings.ToUpper(line)

	// Exclude scanner-related documentation and self-references
	scannerKeywords := []string{
		"TODO SCANNER",
		"TODO PATTERNS",
		"TODO MARKERS",
		"PANIC(\"TODO\")",
		"-SCAN:",
		"SCANNER -",
		"TODO/FIXME/XXX",
		"# TODO Scanner",
		"// TODO patterns",
		"PANIC(\"TODO\") special case",
	}

	for _, keyword := range scannerKeywords {
		if strings.Contains(lineUpper, strings.ToUpper(keyword)) {
			return true
		}
	}

	// Exclude Makefile targets and section headers
	if strings.Contains(path, "Makefile") || strings.HasSuffix(path, ".mk") {
		if strings.Contains(lineUpper, "TODO-SCAN:") ||
		   strings.Contains(lineUpper, "# TODO SCANNER") ||
		   (strings.HasPrefix(strings.TrimSpace(lineUpper), "TODO-SCAN:")) {
			return true
		}
	}

	// Exclude documentation table entries
	if strings.Contains(lineUpper, "|") && strings.Contains(lineUpper, "TODO") {
		return true
	}

	return false
}

func generateSummary(items []TodoItem) TodoSummary {
	countsByTag := make(map[string]int)
	countsByDir := make(map[string]int)
	countsByFile := make(map[string]int)

	for _, item := range items {
		// Count by tag
		countsByTag[item.Tag]++

		// Count by directory (top-level)
		dir := filepath.Dir(item.Path)
		if dir == "." {
			dir = "root"
		} else {
			// Get top-level directory
			parts := strings.Split(dir, string(filepath.Separator))
			if len(parts) > 0 {
				dir = parts[0]
			}
		}
		countsByDir[dir]++

		// Count by file
		countsByFile[item.Path]++
	}

	return TodoSummary{
		CountsByTag:  countsByTag,
		CountsByDir:  countsByDir,
		CountsByFile: countsByFile,
	}
}

func generateJSONReport(report TodoReport) error {
	file, err := os.Create("reports/todo-report.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}

func generateMarkdownReport(report TodoReport) error {
	file, err := os.Create("reports/todo-report.md")
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "# TODO Scan Report\n\n")
	fmt.Fprintf(file, "**Generated:** %s\n", report.GeneratedAt.Format(time.RFC3339))
	fmt.Fprintf(file, "**Repository:** %s\n", report.RepoPath)
	if report.GitSHAShort != "" {
		fmt.Fprintf(file, "**Git SHA:** %s\n", report.GitSHAShort)
	}
	fmt.Fprintf(file, "**Total Items:** %d\n\n", report.TotalCount)

	// Summary by tag
	fmt.Fprintf(file, "## Summary by Tag\n\n")
	var tags []string
	for tag := range report.Summary.CountsByTag {
		tags = append(tags, tag)
	}
	sort.Strings(tags)

	for _, tag := range tags {
		count := report.Summary.CountsByTag[tag]
		fmt.Fprintf(file, "- **%s**: %d\n", tag, count)
	}

	// Summary by directory
	fmt.Fprintf(file, "\n## Summary by Directory\n\n")
	var dirs []string
	for dir := range report.Summary.CountsByDir {
		dirs = append(dirs, dir)
	}
	sort.Strings(dirs)

	for _, dir := range dirs {
		count := report.Summary.CountsByDir[dir]
		fmt.Fprintf(file, "- **%s**: %d\n", dir, count)
	}

	// Top files
	fmt.Fprintf(file, "\n## Top Files by Count\n\n")
	type fileCount struct {
		file  string
		count int
	}
	var fileCounts []fileCount
	for file, count := range report.Summary.CountsByFile {
		fileCounts = append(fileCounts, fileCount{file, count})
	}
	sort.Slice(fileCounts, func(i, j int) bool {
		return fileCounts[i].count > fileCounts[j].count
	})

	// Show top 10
	limit := 10
	if len(fileCounts) < limit {
		limit = len(fileCounts)
	}
	for i := 0; i < limit; i++ {
		fc := fileCounts[i]
		fmt.Fprintf(file, "- **%s**: %d\n", fc.file, fc.count)
	}

	// Detailed items by tag
	fmt.Fprintf(file, "\n## Detailed Items\n\n")

	// Group items by tag
	itemsByTag := make(map[string][]TodoItem)
	for _, item := range report.Items {
		itemsByTag[item.Tag] = append(itemsByTag[item.Tag], item)
	}

	for _, tag := range tags {
		items := itemsByTag[tag]
		if len(items) == 0 {
			continue
		}

		fmt.Fprintf(file, "### %s (%d items)\n\n", tag, len(items))

		// Group by directory
		itemsByDir := make(map[string][]TodoItem)
		for _, item := range items {
			dir := filepath.Dir(item.Path)
			if dir == "." {
				dir = "root"
			}
			itemsByDir[dir] = append(itemsByDir[dir], item)
		}

		var itemDirs []string
		for dir := range itemsByDir {
			itemDirs = append(itemDirs, dir)
		}
		sort.Strings(itemDirs)

		for _, dir := range itemDirs {
			dirItems := itemsByDir[dir]
			fmt.Fprintf(file, "#### %s\n\n", dir)

			for _, item := range dirItems {
				fmt.Fprintf(file, "**%s:%d**\n", item.Path, item.Line)
				fmt.Fprintf(file, "```\n")
				if item.ContextBefore != "" {
					fmt.Fprintf(file, "%d: %s\n", item.Line-1, item.ContextBefore)
				}
				fmt.Fprintf(file, "%d: %s\n", item.Line, item.Text)
				if item.ContextAfter != "" {
					fmt.Fprintf(file, "%d: %s\n", item.Line+1, item.ContextAfter)
				}
				fmt.Fprintf(file, "```\n\n")
			}
		}
	}

	return nil
}

func generateCSVReport(report TodoReport) error {
	file, err := os.Create("reports/todo-report.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"Path", "Line", "Tag", "Text", "ContextBefore", "ContextAfter"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data
	for _, item := range report.Items {
		record := []string{
			item.Path,
			strconv.Itoa(item.Line),
			item.Tag,
			item.Text,
			item.ContextBefore,
			item.ContextAfter,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}
