# TODO Scan Report

**Generated:** 2025-09-23T10:52:52Z
**Repository:** .
**Git SHA:** df90a6a
**Total Items:** 576

## Summary by Tag

- **DEPRECATED**: 4
- **FIXME**: 19
- **HACK**: 12
- **NOTIMPLEMENTED**: 10
- **STUB**: 16
- **TBA**: 5
- **TBD**: 7
- **TODO**: 490
- **XXX**: 13

## Summary by Directory

- **.claude**: 1
- **docs**: 60
- **drivers**: 1
- **registrar-go**: 1
- **resolver-go**: 2
- **root**: 64
- **scripts**: 43
- **spec**: 367
- **tools**: 37

## Top Files by Count

- **spec\PARITY-UNI-DRIVERS.md**: 156
- **spec\PARITY-RESOLVER-REGISTRAR.md**: 109
- **spec\PARITY-SPEC-RESOLVER.md**: 96
- **docs\ops\OPERATIONS.md**: 58
- **CLAUDE.md**: 52
- **tools\todoscan\main.go**: 37
- **scripts\todo-scan.ps1**: 24
- **scripts\todo-scan.sh**: 19
- **Makefile**: 12
- **spec\BACKLOG.md**: 6

## Detailed Items

### DEPRECATED (4 items)

#### docs\ops

**docs\ops\OPERATIONS.md:741**
```
740: | `PANIC("TODO")` | Critical unimplemented paths | Critical |
741: | `@deprecated` | Deprecated code | Medium |
```

#### root

**CLAUDE.md:229**
```
228: - **PANIC("TODO")**: Critical unimplemented paths
229: - **@deprecated/DEPRECATED**: Deprecated code
```

#### tools\todoscan

**tools\todoscan\main.go:51**
```
50: `(?i)PANIC\s*\(\s*["']TODO`,
51: `(?i)@deprecated`,
52: `(?i)\bDEPRECATE\b`,
```

**tools\todoscan\main.go:305**
```
304: }
305: if strings.Contains(match, "DEPRECATE") {
306: return "DEPRECATED"
```

### FIXME (19 items)

#### docs\ops

**docs\ops\OPERATIONS.md:734**
```
733: | `TODO` | General work items | Medium |
734: | `FIXME` | Known bugs/issues | High |
735: | `XXX` | Code requiring attention | High |
```

**docs\ops\OPERATIONS.md:759**
```
758: # 2. Review critical items
759: grep -E "(FIXME|XXX|PANIC)" reports/todo-report.md
```

**docs\ops\OPERATIONS.md:798**
```
797: **Escalation Criteria:**
798: - **FIXME/XXX items**: Convert to GitHub issues if affecting operations
799: - **NOTIMPL items**: Add to `spec/BACKLOG.md` if blocking features
```

**docs\ops\OPERATIONS.md:811**
```
810: run: |
811: if jq -e '.items[] | select(.tag == "PANIC" or .tag == "FIXME")' reports/todo-report.json > /dev/null; then
812: echo "::warning::Critical TODOs found - review required"
```

**docs\ops\OPERATIONS.md:813**
```
812: echo "::warning::Critical TODOs found - review required"
813: jq '.items[] | select(.tag == "PANIC" or .tag == "FIXME")' reports/todo-report.json
814: fi
```

**docs\ops\OPERATIONS.md:822**
```
821: - Reference issues when applicable: `// TODO(#123): implement batch resolution`
822: - Use appropriate tags: `FIXME` for bugs, `TODO` for features
823: - Avoid generic comments: prefer `TODO: validate DID format` over `TODO: fix this`
```

#### root

**CLAUDE.md:222**
```
221: - **TODO**: General work items
222: - **FIXME**: Bugs that need fixing
223: - **XXX**: Code that needs attention
```

**CLAUDE.md:243**
```
242: # High-priority items
243: grep -E "(FIXME|XXX|PANIC)" reports/todo-report.md
```

**CLAUDE.md:273**
```
272: # High-priority items with context
273: jq '.items[] | select(.tag == "FIXME" or .tag == "XXX")' reports/todo-report.json
274: ```
```

**CLAUDE.md:291**
```
290: run: |
291: if grep -q "PANIC\|FIXME" reports/todo-report.md; then
292: echo "⚠️ Critical TODOs found - review required"
```

**CLAUDE.md:300**
```
299: **TODO lifecycle management:**
300: 1. **New TODOs**: Use specific tags (TODO for features, FIXME for bugs)
301: 2. **Context required**: Include brief description of what needs to be done
```

**CLAUDE.md:307**
```
306: - Use `TODO` for planned features or improvements
307: - Use `FIXME` for known bugs or issues
308: - Use `HACK` for temporary workarounds that need proper solutions
```

**Makefile:282**
```
281: @echo "🔍 Code analysis:"
282: @echo "  todo-scan       - Scan repository for TODO/FIXME/XXX markers"
283: @echo ""
```

#### scripts

**scripts\todo-scan.ps1:2**
```
1: # TODO Scanner - Windows PowerShell Wrapper
2: # Scans the repository for TODO, FIXME, XXX, HACK, and other markers
3: # Generates reports in JSON, Markdown, and CSV formats
```

**scripts\todo-scan.sh:4**
```
3: # TODO Scanner - Linux/Docker Wrapper
4: # Scans the repository for TODO, FIXME, XXX, HACK, and other markers
5: # Generates reports in JSON, Markdown, and CSV formats
```

#### tools\todoscan

**tools\todoscan\main.go:18**
```
18: // TodoItem represents a single TODO/FIXME/etc finding
19: type TodoItem struct {
```

**tools\todoscan\main.go:49**
```
48: todoPatterns = []string{
49: `(?i)\b(TODO|FIXME|XXX|HACK|STUB|TBA|TBD|NOTIMPL|NOTIMPLEMENTED)\b`,
50: `(?i)PANIC\s*\(\s*["']TODO`,
```

**tools\todoscan\main.go:278**
```
277: }
278: if strings.Contains(match, "FIXME") {
279: return "FIXME"
```

**tools\todoscan\main.go:279**
```
278: if strings.Contains(match, "FIXME") {
279: return "FIXME"
280: }
```

### HACK (12 items)

#### docs\ops

**docs\ops\OPERATIONS.md:736**
```
735: | `XXX` | Code requiring attention | High |
736: | `HACK` | Temporary workarounds | Medium |
737: | `STUB` | Placeholder implementations | Medium |
```

**docs\ops\OPERATIONS.md:765**
```
764: # 4. Monitor technical debt
765: grep -E "(HACK|DEPRECATED)" reports/todo-report.md
766: ```
```

**docs\ops\OPERATIONS.md:801**
```
800: - **PANIC items**: Address immediately - these indicate critical gaps
801: - **HACK items**: Schedule proper implementation in upcoming sprints
```

**docs\ops\OPERATIONS.md:828**
```
827: 2. **Sprint planning**: Convert high-priority TODOs to formal tasks
828: 3. **Refactoring sprints**: Dedicate time to address HACK items
829: 4. **Release preparation**: Ensure no PANIC items in production code
```

#### root

**CLAUDE.md:224**
```
223: - **XXX**: Code that needs attention
224: - **HACK**: Temporary workarounds
225: - **STUB**: Placeholder implementations
```

**CLAUDE.md:249**
```
248: # Technical debt
249: grep -E "(HACK|DEPRECATED)" reports/todo-report.md
250: ```
```

**CLAUDE.md:308**
```
307: - Use `FIXME` for known bugs or issues
308: - Use `HACK` for temporary workarounds that need proper solutions
309: - Use `DEPRECATED` when phasing out code
```

#### scripts

**scripts\todo-scan.ps1:2**
```
1: # TODO Scanner - Windows PowerShell Wrapper
2: # Scans the repository for TODO, FIXME, XXX, HACK, and other markers
3: # Generates reports in JSON, Markdown, and CSV formats
```

**scripts\todo-scan.sh:4**
```
3: # TODO Scanner - Linux/Docker Wrapper
4: # Scans the repository for TODO, FIXME, XXX, HACK, and other markers
5: # Generates reports in JSON, Markdown, and CSV formats
```

#### tools\todoscan

**tools\todoscan\main.go:49**
```
48: todoPatterns = []string{
49: `(?i)\b(TODO|FIXME|XXX|HACK|STUB|TBA|TBD|NOTIMPL|NOTIMPLEMENTED)\b`,
50: `(?i)PANIC\s*\(\s*["']TODO`,
```

**tools\todoscan\main.go:284**
```
283: }
284: if strings.Contains(match, "HACK") {
285: return "HACK"
```

**tools\todoscan\main.go:285**
```
284: if strings.Contains(match, "HACK") {
285: return "HACK"
286: }
```

### NOTIMPLEMENTED (10 items)

#### docs\ops

**docs\ops\OPERATIONS.md:739**
```
738: | `TBA/TBD` | Items to be added/determined | Low |
739: | `NOTIMPL` | Missing implementations | High |
740: | `PANIC("TODO")` | Critical unimplemented paths | Critical |
```

**docs\ops\OPERATIONS.md:762**
```
761: # 3. Check implementation gaps
762: grep -E "(NOTIMPL|STUB)" reports/todo-report.md
```

**docs\ops\OPERATIONS.md:799**
```
798: - **FIXME/XXX items**: Convert to GitHub issues if affecting operations
799: - **NOTIMPL items**: Add to `spec/BACKLOG.md` if blocking features
800: - **PANIC items**: Address immediately - these indicate critical gaps
```

#### root

**CLAUDE.md:227**
```
226: - **TBA/TBD**: To be added/determined
227: - **NOTIMPL/NOTIMPLEMENTED**: Missing implementations
228: - **PANIC("TODO")**: Critical unimplemented paths
```

**CLAUDE.md:227**
```
226: - **TBA/TBD**: To be added/determined
227: - **NOTIMPL/NOTIMPLEMENTED**: Missing implementations
228: - **PANIC("TODO")**: Critical unimplemented paths
```

**CLAUDE.md:246**
```
245: # Implementation gaps
246: grep -E "(TODO|NOTIMPL|STUB)" reports/todo-report.md
```

#### tools\todoscan

**tools\todoscan\main.go:49**
```
48: todoPatterns = []string{
49: `(?i)\b(TODO|FIXME|XXX|HACK|STUB|TBA|TBD|NOTIMPL|NOTIMPLEMENTED)\b`,
50: `(?i)PANIC\s*\(\s*["']TODO`,
```

**tools\todoscan\main.go:49**
```
48: todoPatterns = []string{
49: `(?i)\b(TODO|FIXME|XXX|HACK|STUB|TBA|TBD|NOTIMPL|NOTIMPLEMENTED)\b`,
50: `(?i)PANIC\s*\(\s*["']TODO`,
```

**tools\todoscan\main.go:296**
```
295: }
296: if strings.Contains(match, "NOTIMPL") {
297: return "NOTIMPLEMENTED"
```

**tools\todoscan\main.go:297**
```
296: if strings.Contains(match, "NOTIMPL") {
297: return "NOTIMPLEMENTED"
298: }
```

### STUB (16 items)

#### .claude\memory

**.claude\memory\CLAUDE.md:265**
```
264: - [ ] Table-driven tests with goldens
265: - [ ] Offline Accumulate client stub
```

#### docs\ops

**docs\ops\OPERATIONS.md:737**
```
736: | `HACK` | Temporary workarounds | Medium |
737: | `STUB` | Placeholder implementations | Medium |
738: | `TBA/TBD` | Items to be added/determined | Low |
```

**docs\ops\OPERATIONS.md:762**
```
761: # 3. Check implementation gaps
762: grep -E "(NOTIMPL|STUB)" reports/todo-report.md
```

#### drivers\uniregistrar-go

**drivers\uniregistrar-go\smoke.ps1:109**
```
109: # Optional: Test update and deactivate with stub requests
110: Write-Host "`n=====================================" -ForegroundColor Cyan
```

#### resolver-go\internal\resolve

**resolver-go\internal\resolve\deterministic_test.go:138**
```
137: func (m *DeterministicMockClient) GetKeyPageState(u string) (acc.KeyPageState, error) {
138: // Minimal stub for tests
139: return acc.KeyPageState{URL: u, Threshold: 1}, nil
```

#### root

**CLAUDE.md:225**
```
224: - **HACK**: Temporary workarounds
225: - **STUB**: Placeholder implementations
226: - **TBA/TBD**: To be added/determined
```

**CLAUDE.md:246**
```
245: # Implementation gaps
246: grep -E "(TODO|NOTIMPL|STUB)" reports/todo-report.md
```

#### spec

**spec\BACKLOG.md:41**
```
40: | R-003 | Implement canonical JSON | resolver-go | 🟡 IN-PROGRESS | Resolver Builder | internal/canon/json.go exists, hash tests failing |
41: | R-004 | Create Accumulate client stub | resolver-go | ✅ DONE | Resolver Builder | internal/acc/client.go + mock.go implemented |
42: | R-005 | Implement DID resolution logic | resolver-go | ✅ DONE | Resolver Builder | Deterministic resolver with full algorithm |
```

**spec\BACKLOG.md:70**
```
69: | G-003 | Implement auth policy v1 | registrar-go | ✅ DONE | Registrar Builder | internal/policy/v1.go with acc://<adi>/book/1 |
70: | G-004 | Create Accumulate client stub | registrar-go | ✅ DONE | Registrar Builder | internal/acc/submit.go + mock.go |
71: | G-005 | Implement registration logic | registrar-go | ✅ DONE | Registrar Builder | Complete create/update/deactivate workflows |
```

**spec\BACKLOG.md:133**
```
132: |----|------|-----------|--------|----------|-------|
133: | D-006 | Write didcomm.md stub | docs/interop | ✅ DONE | Interop Engineer | docs/interop/didcomm.md exists |
134: | D-007 | Write sd-jwt.md stub | docs/interop | ✅ DONE | Interop Engineer | docs/interop/sd-jwt.md exists |
```

**spec\BACKLOG.md:134**
```
133: | D-006 | Write didcomm.md stub | docs/interop | ✅ DONE | Interop Engineer | docs/interop/didcomm.md exists |
134: | D-007 | Write sd-jwt.md stub | docs/interop | ✅ DONE | Interop Engineer | docs/interop/sd-jwt.md exists |
135: | D-008 | Write bbs.md stub | docs/interop | ✅ DONE | Interop Engineer | docs/interop/bbs.md exists |
```

**spec\BACKLOG.md:135**
```
134: | D-007 | Write sd-jwt.md stub | docs/interop | ✅ DONE | Interop Engineer | docs/interop/sd-jwt.md exists |
135: | D-008 | Write bbs.md stub | docs/interop | ✅ DONE | Interop Engineer | docs/interop/bbs.md exists |
```

**spec\BACKLOG.md:185**
```
184: ### Mitigation
185: - Use stub clients for offline development
186: - Implement multiple canonical JSON options
```

#### tools\todoscan

**tools\todoscan\main.go:49**
```
48: todoPatterns = []string{
49: `(?i)\b(TODO|FIXME|XXX|HACK|STUB|TBA|TBD|NOTIMPL|NOTIMPLEMENTED)\b`,
50: `(?i)PANIC\s*\(\s*["']TODO`,
```

**tools\todoscan\main.go:287**
```
286: }
287: if strings.Contains(match, "STUB") {
288: return "STUB"
```

**tools\todoscan\main.go:288**
```
287: if strings.Contains(match, "STUB") {
288: return "STUB"
289: }
```

### TBA (5 items)

#### docs\ops

**docs\ops\OPERATIONS.md:738**
```
737: | `STUB` | Placeholder implementations | Medium |
738: | `TBA/TBD` | Items to be added/determined | Low |
739: | `NOTIMPL` | Missing implementations | High |
```

#### root

**CLAUDE.md:226**
```
225: - **STUB**: Placeholder implementations
226: - **TBA/TBD**: To be added/determined
227: - **NOTIMPL/NOTIMPLEMENTED**: Missing implementations
```

#### tools\todoscan

**tools\todoscan\main.go:49**
```
48: todoPatterns = []string{
49: `(?i)\b(TODO|FIXME|XXX|HACK|STUB|TBA|TBD|NOTIMPL|NOTIMPLEMENTED)\b`,
50: `(?i)PANIC\s*\(\s*["']TODO`,
```

**tools\todoscan\main.go:290**
```
289: }
290: if strings.Contains(match, "TBA") {
291: return "TBA"
```

**tools\todoscan\main.go:291**
```
290: if strings.Contains(match, "TBA") {
291: return "TBA"
292: }
```

### TBD (7 items)

#### docs\ops

**docs\ops\OPERATIONS.md:738**
```
737: | `STUB` | Placeholder implementations | Medium |
738: | `TBA/TBD` | Items to be added/determined | Low |
739: | `NOTIMPL` | Missing implementations | High |
```

#### docs\site

**docs\site\method.md:319**
```
318: **Last Updated:** 2024-01-21
319: **Next Review:** TBD based on implementation feedback
```

#### docs\spec

**docs\spec\method.md:319**
```
318: **Last Updated:** 2024-01-21
319: **Next Review:** TBD based on implementation feedback
```

#### root

**CLAUDE.md:226**
```
225: - **STUB**: Placeholder implementations
226: - **TBA/TBD**: To be added/determined
227: - **NOTIMPL/NOTIMPLEMENTED**: Missing implementations
```

#### tools\todoscan

**tools\todoscan\main.go:49**
```
48: todoPatterns = []string{
49: `(?i)\b(TODO|FIXME|XXX|HACK|STUB|TBA|TBD|NOTIMPL|NOTIMPLEMENTED)\b`,
50: `(?i)PANIC\s*\(\s*["']TODO`,
```

**tools\todoscan\main.go:293**
```
292: }
293: if strings.Contains(match, "TBD") {
294: return "TBD"
```

**tools\todoscan\main.go:294**
```
293: if strings.Contains(match, "TBD") {
294: return "TBD"
295: }
```

### TODO (490 items)

#### docs\ops

**docs\ops\OPERATIONS.md:701**
```
701: ## 8. Backlog Triage via TODO Scanner
```

**docs\ops\OPERATIONS.md:703**
```
703: The repository includes a comprehensive TODO scanner for tracking technical debt and work items across the codebase.
```

**docs\ops\OPERATIONS.md:709**
```
708: ```bash
709: make todo-scan
710: ```
```

**docs\ops\OPERATIONS.md:718**
```
717: # Linux/Unix shell script
718: ./scripts/todo-scan.sh
```

**docs\ops\OPERATIONS.md:721**
```
720: # Windows PowerShell script
721: .\scripts\todo-scan.ps1
```

**docs\ops\OPERATIONS.md:727**
```
727: ### TODO Markers Detected
```

**docs\ops\OPERATIONS.md:733**
```
732: |-----|---------|----------|
733: | `TODO` | General work items | Medium |
734: | `FIXME` | Known bugs/issues | High |
```

**docs\ops\OPERATIONS.md:740**
```
739: | `NOTIMPL` | Missing implementations | High |
740: | `PANIC("TODO")` | Critical unimplemented paths | Critical |
741: | `@deprecated` | Deprecated code | Medium |
```

**docs\ops\OPERATIONS.md:740**
```
739: | `NOTIMPL` | Missing implementations | High |
740: | `PANIC("TODO")` | Critical unimplemented paths | Critical |
741: | `@deprecated` | Deprecated code | Medium |
```

**docs\ops\OPERATIONS.md:747**
```
747: - **`todo-report.json`** - Machine-readable data with full context
748: - **`todo-report.md`** - Human-readable report with code excerpts
```

**docs\ops\OPERATIONS.md:748**
```
747: - **`todo-report.json`** - Machine-readable data with full context
748: - **`todo-report.md`** - Human-readable report with code excerpts
749: - **`todo-report.csv`** - Spreadsheet-compatible format for analysis
```

**docs\ops\OPERATIONS.md:749**
```
748: - **`todo-report.md`** - Human-readable report with code excerpts
749: - **`todo-report.csv`** - Spreadsheet-compatible format for analysis
```

**docs\ops\OPERATIONS.md:756**
```
755: # 1. Run scan during sprint planning
756: make todo-scan
```

**docs\ops\OPERATIONS.md:759**
```
758: # 2. Review critical items
759: grep -E "(FIXME|XXX|PANIC)" reports/todo-report.md
```

**docs\ops\OPERATIONS.md:762**
```
761: # 3. Check implementation gaps
762: grep -E "(NOTIMPL|STUB)" reports/todo-report.md
```

**docs\ops\OPERATIONS.md:765**
```
764: # 4. Monitor technical debt
765: grep -E "(HACK|DEPRECATED)" reports/todo-report.md
766: ```
```

**docs\ops\OPERATIONS.md:771**
```
770: # Resolver issues
771: jq '.items[] | select(.path | startswith("resolver-go/"))' reports/todo-report.json
```

**docs\ops\OPERATIONS.md:774**
```
773: # Registrar issues
774: jq '.items[] | select(.path | startswith("registrar-go/"))' reports/todo-report.json
```

**docs\ops\OPERATIONS.md:777**
```
776: # Documentation tasks
777: jq '.items[] | select(.path | startswith("docs/"))' reports/todo-report.json
```

**docs\ops\OPERATIONS.md:780**
```
779: # Universal driver issues
780: jq '.items[] | select(.path | startswith("drivers/"))' reports/todo-report.json
781: ```
```

**docs\ops\OPERATIONS.md:786**
```
785: # Compare counts over time
786: echo "$(date): $(jq '.totalCount' reports/todo-report.json)" >> reports/todo-trend.log
```

**docs\ops\OPERATIONS.md:786**
```
785: # Compare counts over time
786: echo "$(date): $(jq '.totalCount' reports/todo-report.json)" >> reports/todo-trend.log
```

**docs\ops\OPERATIONS.md:789**
```
788: # Tag distribution
789: jq '.summary.countsByTag' reports/todo-report.json
```

**docs\ops\OPERATIONS.md:792**
```
791: # Directory breakdown
792: jq '.summary.countsByDir' reports/todo-report.json
793: ```
```

**docs\ops\OPERATIONS.md:806**
```
805: # Example CI check
806: - name: TODO Scanner
807: run: make todo-scan
```

**docs\ops\OPERATIONS.md:807**
```
806: - name: TODO Scanner
807: run: make todo-scan
```

**docs\ops\OPERATIONS.md:811**
```
810: run: |
811: if jq -e '.items[] | select(.tag == "PANIC" or .tag == "FIXME")' reports/todo-report.json > /dev/null; then
812: echo "::warning::Critical TODOs found - review required"
```

**docs\ops\OPERATIONS.md:813**
```
812: echo "::warning::Critical TODOs found - review required"
813: jq '.items[] | select(.tag == "PANIC" or .tag == "FIXME")' reports/todo-report.json
814: fi
```

**docs\ops\OPERATIONS.md:819**
```
819: **TODO Creation Guidelines:**
820: - Include brief context: `// TODO: add rate limiting for /resolve endpoint`
```

**docs\ops\OPERATIONS.md:820**
```
819: **TODO Creation Guidelines:**
820: - Include brief context: `// TODO: add rate limiting for /resolve endpoint`
821: - Reference issues when applicable: `// TODO(#123): implement batch resolution`
```

**docs\ops\OPERATIONS.md:821**
```
820: - Include brief context: `// TODO: add rate limiting for /resolve endpoint`
821: - Reference issues when applicable: `// TODO(#123): implement batch resolution`
822: - Use appropriate tags: `FIXME` for bugs, `TODO` for features
```

**docs\ops\OPERATIONS.md:822**
```
821: - Reference issues when applicable: `// TODO(#123): implement batch resolution`
822: - Use appropriate tags: `FIXME` for bugs, `TODO` for features
823: - Avoid generic comments: prefer `TODO: validate DID format` over `TODO: fix this`
```

**docs\ops\OPERATIONS.md:823**
```
822: - Use appropriate tags: `FIXME` for bugs, `TODO` for features
823: - Avoid generic comments: prefer `TODO: validate DID format` over `TODO: fix this`
```

**docs\ops\OPERATIONS.md:823**
```
822: - Use appropriate tags: `FIXME` for bugs, `TODO` for features
823: - Avoid generic comments: prefer `TODO: validate DID format` over `TODO: fix this`
```

**docs\ops\OPERATIONS.md:832**
```
831: **Reporting:**
832: - Include TODO count trends in sprint retrospectives
833: - Track resolution rate of TODO items over time
```

**docs\ops\OPERATIONS.md:833**
```
832: - Include TODO count trends in sprint retrospectives
833: - Track resolution rate of TODO items over time
834: - Use TODO density (items per KLOC) as code quality metric
```

**docs\ops\OPERATIONS.md:834**
```
833: - Track resolution rate of TODO items over time
834: - Use TODO density (items per KLOC) as code quality metric
```

#### registrar-go\internal\acc

**registrar-go\internal\acc\submit.go:681**
```
681: // TODO: Implement proper conversion based on:
682: // 1. ops.Envelope structure
```

#### resolver-go\internal\acc

**resolver-go\internal\acc\client.go:338**
```
337: // For now, return a basic envelope structure
338: // TODO: Implement proper record to envelope conversion based on actual API types
339: envelope := Envelope{
```

#### root

**CLAUDE.md:195**
```
195: ## Backlog Triage via TODO Scanner
```

**CLAUDE.md:197**
```
197: The repository includes a comprehensive TODO scanner that helps track technical debt and work items across the codebase.
```

**CLAUDE.md:199**
```
199: ### Running TODO Scans
```

**CLAUDE.md:203**
```
202: ```bash
203: make todo-scan
204: ```
```

**CLAUDE.md:209**
```
208: # Linux/Docker (recommended)
209: ./scripts/todo-scan.sh
```

**CLAUDE.md:212**
```
211: # Windows PowerShell
212: .\scripts\todo-scan.ps1
```

**CLAUDE.md:218**
```
218: ### TODO Patterns Detected
```

**CLAUDE.md:221**
```
220: The scanner looks for these markers (case-insensitive):
221: - **TODO**: General work items
222: - **FIXME**: Bugs that need fixing
```

**CLAUDE.md:228**
```
227: - **NOTIMPL/NOTIMPLEMENTED**: Missing implementations
228: - **PANIC("TODO")**: Critical unimplemented paths
229: - **@deprecated/DEPRECATED**: Deprecated code
```

**CLAUDE.md:228**
```
227: - **NOTIMPL/NOTIMPLEMENTED**: Missing implementations
228: - **PANIC("TODO")**: Critical unimplemented paths
229: - **@deprecated/DEPRECATED**: Deprecated code
```

**CLAUDE.md:234**
```
233: Reports are generated in `./reports/`:
234: - **todo-report.json**: Machine-readable data for automation
235: - **todo-report.md**: Human-readable report with code excerpts
```

**CLAUDE.md:235**
```
234: - **todo-report.json**: Machine-readable data for automation
235: - **todo-report.md**: Human-readable report with code excerpts
236: - **todo-report.csv**: Spreadsheet-compatible format
```

**CLAUDE.md:236**
```
235: - **todo-report.md**: Human-readable report with code excerpts
236: - **todo-report.csv**: Spreadsheet-compatible format
```

**CLAUDE.md:243**
```
242: # High-priority items
243: grep -E "(FIXME|XXX|PANIC)" reports/todo-report.md
```

**CLAUDE.md:246**
```
245: # Implementation gaps
246: grep -E "(TODO|NOTIMPL|STUB)" reports/todo-report.md
```

**CLAUDE.md:246**
```
245: # Implementation gaps
246: grep -E "(TODO|NOTIMPL|STUB)" reports/todo-report.md
```

**CLAUDE.md:249**
```
248: # Technical debt
249: grep -E "(HACK|DEPRECATED)" reports/todo-report.md
250: ```
```

**CLAUDE.md:255**
```
254: # Resolver issues
255: grep "resolver-go" reports/todo-report.md
```

**CLAUDE.md:258**
```
257: # Registrar issues
258: grep "registrar-go" reports/todo-report.md
```

**CLAUDE.md:261**
```
260: # Documentation tasks
261: grep "docs/" reports/todo-report.md
262: ```
```

**CLAUDE.md:267**
```
266: # Count by tag
267: jq '.summary.countsByTag' reports/todo-report.json
```

**CLAUDE.md:270**
```
269: # Items in specific directory
270: jq '.items[] | select(.path | startswith("resolver-go/"))' reports/todo-report.json
```

**CLAUDE.md:273**
```
272: # High-priority items with context
273: jq '.items[] | select(.tag == "FIXME" or .tag == "XXX")' reports/todo-report.json
274: ```
```

**CLAUDE.md:279**
```
278: **Recommended workflow:**
279: 1. **Weekly scans**: Run `make todo-scan` during sprint planning
280: 2. **Triage new items**: Review `todo-report.md` for high-priority issues
```

**CLAUDE.md:280**
```
279: 1. **Weekly scans**: Run `make todo-scan` during sprint planning
280: 2. **Triage new items**: Review `todo-report.md` for high-priority issues
281: 3. **Convert to issues**: Move critical TODOs into GitHub issues or `spec/BACKLOG.md`
```

**CLAUDE.md:288**
```
287: - name: Scan TODOs
288: run: make todo-scan
289: - name: Check for critical TODOs
```

**CLAUDE.md:291**
```
290: run: |
291: if grep -q "PANIC\|FIXME" reports/todo-report.md; then
292: echo "⚠️ Critical TODOs found - review required"
```

**CLAUDE.md:293**
```
292: echo "⚠️ Critical TODOs found - review required"
293: cat reports/todo-report.md
294: fi
```

**CLAUDE.md:299**
```
299: **TODO lifecycle management:**
300: 1. **New TODOs**: Use specific tags (TODO for features, FIXME for bugs)
```

**CLAUDE.md:300**
```
299: **TODO lifecycle management:**
300: 1. **New TODOs**: Use specific tags (TODO for features, FIXME for bugs)
301: 2. **Context required**: Include brief description of what needs to be done
```

**CLAUDE.md:306**
```
305: **Tag guidelines:**
306: - Use `TODO` for planned features or improvements
307: - Use `FIXME` for known bugs or issues
```

**Makefile:100**
```
100: .PHONY: dev-shell test-all ci-local dev-up dev-down check-imports conformance perf help lint test-race vet qa dist-clean binaries-local images-local sbom-local scan-local docs-archive release-local sdk-test sdk-merge-spec example-sdk todo-scan
```

**Makefile:228**
```
227: # ========================================================================
228: # TODO Scanner
229: # ========================================================================
```

**Makefile:231**
```
231: todo-scan:
232: @echo "🔍 Scanning repository for TODO markers..."
```

**Makefile:232**
```
231: todo-scan:
232: @echo "🔍 Scanning repository for TODO markers..."
233: @mkdir -p reports
```

**Makefile:235**
```
234: @if command -v docker >/dev/null 2>&1 && test -f docker-compose.dev.yml; then \
235: docker compose -f docker-compose.dev.yml run --rm dev 'bash -lc "go run tools/todoscan/main.go . && echo \"\" && echo \"📊 Quick Summary:\" && test -f reports/todo-report.json && jq -r \".totalCount,(.summary.countsByTag | to_entries[] | \\\"  - \\(.key): \\(.value)\\\")\" reports/todo-report.json || echo \"Reports generated in ./reports/\""'; \
236: elif command -v go >/dev/null 2>&1 && test -f tools/todoscan/main.go; then \
```

**Makefile:235**
```
234: @if command -v docker >/dev/null 2>&1 && test -f docker-compose.dev.yml; then \
235: docker compose -f docker-compose.dev.yml run --rm dev 'bash -lc "go run tools/todoscan/main.go . && echo \"\" && echo \"📊 Quick Summary:\" && test -f reports/todo-report.json && jq -r \".totalCount,(.summary.countsByTag | to_entries[] | \\\"  - \\(.key): \\(.value)\\\")\" reports/todo-report.json || echo \"Reports generated in ./reports/\""'; \
236: elif command -v go >/dev/null 2>&1 && test -f tools/todoscan/main.go; then \
```

**Makefile:240**
```
239: echo "📊 Quick Summary:"; \
240: if command -v jq >/dev/null 2>&1 && test -f reports/todo-report.json; then \
241: jq -r '.totalCount,(.summary.countsByTag | to_entries[] | "  - \(.key): \(.value)")' reports/todo-report.json; \
```

**Makefile:241**
```
240: if command -v jq >/dev/null 2>&1 && test -f reports/todo-report.json; then \
241: jq -r '.totalCount,(.summary.countsByTag | to_entries[] | "  - \(.key): \(.value)")' reports/todo-report.json; \
242: else \
```

**Makefile:282**
```
281: @echo "🔍 Code analysis:"
282: @echo "  todo-scan       - Scan repository for TODO/FIXME/XXX markers"
283: @echo ""
```

**Makefile:282**
```
281: @echo "🔍 Code analysis:"
282: @echo "  todo-scan       - Scan repository for TODO/FIXME/XXX markers"
283: @echo ""
```

#### scripts

**scripts\todo-scan.ps1:1**
```
1: # TODO Scanner - Windows PowerShell Wrapper
2: # Scans the repository for TODO, FIXME, XXX, HACK, and other markers
```

**scripts\todo-scan.ps1:2**
```
1: # TODO Scanner - Windows PowerShell Wrapper
2: # Scans the repository for TODO, FIXME, XXX, HACK, and other markers
3: # Generates reports in JSON, Markdown, and CSV formats
```

**scripts\todo-scan.ps1:23**
```
22: Write-Host @"
23: TODO Scanner for Accumulate DID Repository
```

**scripts\todo-scan.ps1:26**
```
25: USAGE:
26: .\todo-scan.ps1 [RepoPath] [OPTIONS]
```

**scripts\todo-scan.ps1:35**
```
34: EXAMPLES:
35: .\todo-scan.ps1                                    # Scan current directory
36: .\todo-scan.ps1 C:\path\to\repo                   # Scan specific repository
```

**scripts\todo-scan.ps1:36**
```
35: .\todo-scan.ps1                                    # Scan current directory
36: .\todo-scan.ps1 C:\path\to\repo                   # Scan specific repository
37: .\todo-scan.ps1 -UseDocker yes                    # Force Docker usage
```

**scripts\todo-scan.ps1:37**
```
36: .\todo-scan.ps1 C:\path\to\repo                   # Scan specific repository
37: .\todo-scan.ps1 -UseDocker yes                    # Force Docker usage
38: .\todo-scan.ps1 -UseDocker no                     # Force local Go usage
```

**scripts\todo-scan.ps1:38**
```
37: .\todo-scan.ps1 -UseDocker yes                    # Force Docker usage
38: .\todo-scan.ps1 -UseDocker no                     # Force local Go usage
```

**scripts\todo-scan.ps1:42**
```
41: Reports are generated in: .\reports\
42: - todo-report.json     # Machine-readable JSON
43: - todo-report.md       # Human-readable Markdown
```

**scripts\todo-scan.ps1:43**
```
42: - todo-report.json     # Machine-readable JSON
43: - todo-report.md       # Human-readable Markdown
44: - todo-report.csv      # Spreadsheet-compatible CSV
```

**scripts\todo-scan.ps1:44**
```
43: - todo-report.md       # Human-readable Markdown
44: - todo-report.csv      # Spreadsheet-compatible CSV
```

**scripts\todo-scan.ps1:100**
```
99: function Invoke-LocalScan {
100: Write-StatusMessage "🔍" "Running TODO scanner locally..." "Blue"
```

**scripts\todo-scan.ps1:133**
```
132: function Invoke-DockerScan {
133: Write-StatusMessage "🐳" "Running TODO scanner in Docker..." "Blue"
```

**scripts\todo-scan.ps1:148**
```
147: $command = @"
148: echo 'Running TODO scanner...'
149: if [[ -f tools/todoscan/main.go ]]; then
```

**scripts\todo-scan.ps1:175**
```
174: if [[ -f tools/todoscan/main.go ]]; then
175: echo 'Running TODO scanner...'
176: go run tools/todoscan/main.go .
```

**scripts\todo-scan.ps1:194**
```
193: $outputDir = Join-Path $RepoPath "reports"
194: $jsonFile = Join-Path $outputDir "todo-report.json"
195: $mdFile = Join-Path $outputDir "todo-report.md"
```

**scripts\todo-scan.ps1:195**
```
194: $jsonFile = Join-Path $outputDir "todo-report.json"
195: $mdFile = Join-Path $outputDir "todo-report.md"
196: $csvFile = Join-Path $outputDir "todo-report.csv"
```

**scripts\todo-scan.ps1:196**
```
195: $mdFile = Join-Path $outputDir "todo-report.md"
196: $csvFile = Join-Path $outputDir "todo-report.csv"
```

**scripts\todo-scan.ps1:206**
```
205: $totalCount = $report.totalCount
206: Write-StatusMessage "📊" "Found $totalCount TODO items" "Blue"
```

**scripts\todo-scan.ps1:239**
```
238: Write-Host "  - Process JSON programmatically: $jsonFile"
239: Write-Host "  - Filter by tag: Select-String 'TODO' $mdFile"
240: Write-Host "  - Filter by directory: Select-String 'resolver-go' $mdFile"
```

**scripts\todo-scan.ps1:249**
```
249: Write-StatusMessage "🔍" "TODO Scanner for Accumulate DID Repository" "Blue"
250: Write-StatusMessage "📂" "Repository: $RepoPath" "Blue"
```

**scripts\todo-scan.sh:3**
```
2: #
3: # TODO Scanner - Linux/Docker Wrapper
4: # Scans the repository for TODO, FIXME, XXX, HACK, and other markers
```

**scripts\todo-scan.sh:4**
```
3: # TODO Scanner - Linux/Docker Wrapper
4: # Scans the repository for TODO, FIXME, XXX, HACK, and other markers
5: # Generates reports in JSON, Markdown, and CSV formats
```

**scripts\todo-scan.sh:27**
```
27: Scans a repository for TODO markers and generates reports.
```

**scripts\todo-scan.sh:43**
```
42: Reports are generated in: ./reports/
43: - todo-report.json         # Machine-readable JSON
44: - todo-report.md          # Human-readable Markdown
```

**scripts\todo-scan.sh:44**
```
43: - todo-report.json         # Machine-readable JSON
44: - todo-report.md          # Human-readable Markdown
45: - todo-report.csv         # Spreadsheet-compatible CSV
```

**scripts\todo-scan.sh:45**
```
44: - todo-report.md          # Human-readable Markdown
45: - todo-report.csv         # Spreadsheet-compatible CSV
46: EOF
```

**scripts\todo-scan.sh:75**
```
74: run_local() {
75: echo -e "${BLUE}🔍${NC} Running TODO scanner locally..."
```

**scripts\todo-scan.sh:93**
```
92: run_docker() {
93: echo -e "${BLUE}🐳${NC} Running TODO scanner in Docker..."
```

**scripts\todo-scan.sh:103**
```
102: docker-compose -f docker-compose.dev.yml run --rm dev bash -c "
103: echo 'Running TODO scanner...'
104: if [[ -f tools/todoscan/main.go ]]; then
```

**scripts\todo-scan.sh:122**
```
121: if [[ -f tools/todoscan/main.go ]]; then
122: echo 'Running TODO scanner...'
123: go run tools/todoscan/main.go .
```

**scripts\todo-scan.sh:133**
```
132: print_results() {
133: local json_file="$OUTPUT_DIR/todo-report.json"
134: local md_file="$OUTPUT_DIR/todo-report.md"
```

**scripts\todo-scan.sh:134**
```
133: local json_file="$OUTPUT_DIR/todo-report.json"
134: local md_file="$OUTPUT_DIR/todo-report.md"
135: local csv_file="$OUTPUT_DIR/todo-report.csv"
```

**scripts\todo-scan.sh:135**
```
134: local md_file="$OUTPUT_DIR/todo-report.md"
135: local csv_file="$OUTPUT_DIR/todo-report.csv"
```

**scripts\todo-scan.sh:143**
```
142: local total_count=$(jq -r '.totalCount // 0' "$json_file" 2>/dev/null || echo "unknown")
143: echo -e "${BLUE}📊${NC} Found ${YELLOW}$total_count${NC} TODO items"
```

**scripts\todo-scan.sh:168**
```
167: echo "  - Process JSON programmatically: $json_file"
168: echo "  - Filter by tag: grep 'TODO' $md_file"
169: echo "  - Filter by directory: grep 'resolver-go' $md_file"
```

**scripts\todo-scan.sh:181**
```
181: echo -e "${BLUE}🔍${NC} TODO Scanner for Accumulate DID Repository"
182: echo -e "${BLUE}📂${NC} Repository: $REPO_PATH"
```

#### spec

**spec\PARITY-RESOLVER-REGISTRAR.md:12**
```
11: |---------|------------------|-------------------|-------------------|
12: | **@context** | Must validate and return | Must validate on input | ❌ TODO |
13: | **id** | Must match resolved DID | Must match registration DID | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:13**
```
12: | **@context** | Must validate and return | Must validate on input | ❌ TODO |
13: | **id** | Must match resolved DID | Must match registration DID | ❌ TODO |
14: | **controller** | Return as stored | Validate controller format | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:14**
```
13: | **id** | Must match resolved DID | Must match registration DID | ❌ TODO |
14: | **controller** | Return as stored | Validate controller format | ❌ TODO |
15: | **verificationMethod** | Return with full details | Validate VM structure | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:15**
```
14: | **controller** | Return as stored | Validate controller format | ❌ TODO |
15: | **verificationMethod** | Return with full details | Validate VM structure | ❌ TODO |
16: | **authentication** | Return reference/embed | Validate references exist | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:16**
```
15: | **verificationMethod** | Return with full details | Validate VM structure | ❌ TODO |
16: | **authentication** | Return reference/embed | Validate references exist | ❌ TODO |
17: | **assertionMethod** | Return reference/embed | Validate references exist | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:17**
```
16: | **authentication** | Return reference/embed | Validate references exist | ❌ TODO |
17: | **assertionMethod** | Return reference/embed | Validate references exist | ❌ TODO |
18: | **service** | Return service endpoints | Validate service format | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:18**
```
17: | **assertionMethod** | Return reference/embed | Validate references exist | ❌ TODO |
18: | **service** | Return service endpoints | Validate service format | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:24**
```
23: |----------|----------------|-----------------|-------------------|
24: | **type** | "AccumulateKeyPage" | Must be "AccumulateKeyPage" | ❌ TODO |
25: | **keyPageUrl** | Return as "acc://..." | Validate "acc://..." format | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:25**
```
24: | **type** | "AccumulateKeyPage" | Must be "AccumulateKeyPage" | ❌ TODO |
25: | **keyPageUrl** | Return as "acc://..." | Validate "acc://..." format | ❌ TODO |
26: | **threshold** | Return as number | Validate positive integer | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:26**
```
25: | **keyPageUrl** | Return as "acc://..." | Validate "acc://..." format | ❌ TODO |
26: | **threshold** | Return as number | Validate positive integer | ❌ TODO |
27: | **controller** | Return DID reference | Validate matches DID | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:27**
```
26: | **threshold** | Return as number | Validate positive integer | ❌ TODO |
27: | **controller** | Return DID reference | Validate matches DID | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:33**
```
32: |-------|----------------|---------------------|-------------------|
33: | **versionId** | From stored metadata | Generate on create/update | ❌ TODO |
34: | **created** | First version timestamp | Set on create | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:34**
```
33: | **versionId** | From stored metadata | Generate on create/update | ❌ TODO |
34: | **created** | First version timestamp | Set on create | ❌ TODO |
35: | **updated** | Latest version timestamp | Set on update | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:35**
```
34: | **created** | First version timestamp | Set on create | ❌ TODO |
35: | **updated** | Latest version timestamp | Set on update | ❌ TODO |
36: | **deactivated** | Boolean from document | Set on deactivate | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:36**
```
35: | **updated** | Latest version timestamp | Set on update | ❌ TODO |
36: | **deactivated** | Boolean from document | Set on deactivate | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:44**
```
43: |------|------------------------|-------------------------|--------|
44: | **Key Ordering** | Lexicographic sort | Lexicographic sort | ❌ TODO |
45: | **Whitespace** | No extra whitespace | No extra whitespace | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:45**
```
44: | **Key Ordering** | Lexicographic sort | Lexicographic sort | ❌ TODO |
45: | **Whitespace** | No extra whitespace | No extra whitespace | ❌ TODO |
46: | **Number Format** | No trailing zeros | No trailing zeros | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:46**
```
45: | **Whitespace** | No extra whitespace | No extra whitespace | ❌ TODO |
46: | **Number Format** | No trailing zeros | No trailing zeros | ❌ TODO |
47: | **String Escaping** | Minimal escaping | Minimal escaping | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:47**
```
46: | **Number Format** | No trailing zeros | No trailing zeros | ❌ TODO |
47: | **String Escaping** | Minimal escaping | Minimal escaping | ❌ TODO |
48: | **Duplicate Keys** | Reject on parse | Reject on validation | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:48**
```
47: | **String Escaping** | Minimal escaping | Minimal escaping | ❌ TODO |
48: | **Duplicate Keys** | Reject on parse | Reject on validation | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:54**
```
53: |--------|----------|-----------|-------------------|
54: | **Algorithm** | SHA-256 | SHA-256 | ❌ TODO |
55: | **Input Format** | Canonical JSON | Canonical JSON | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:55**
```
54: | **Algorithm** | SHA-256 | SHA-256 | ❌ TODO |
55: | **Input Format** | Canonical JSON | Canonical JSON | ❌ TODO |
56: | **Output Format** | "sha256:..." | "sha256:..." | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:56**
```
55: | **Input Format** | Canonical JSON | Canonical JSON | ❌ TODO |
56: | **Output Format** | "sha256:..." | "sha256:..." | ❌ TODO |
57: | **Verification** | Compare stored hash | Generate content hash | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:57**
```
56: | **Output Format** | "sha256:..." | "sha256:..." | ❌ TODO |
57: | **Verification** | Compare stored hash | Generate content hash | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:65**
```
64: |------|----------------------|------------------------|--------|
65: | **Case Sensitivity** | Convert to lowercase | Convert to lowercase | ❌ TODO |
66: | **Trailing Dots** | Remove trailing dots | Remove trailing dots | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:66**
```
65: | **Case Sensitivity** | Convert to lowercase | Convert to lowercase | ❌ TODO |
66: | **Trailing Dots** | Remove trailing dots | Remove trailing dots | ❌ TODO |
67: | **Query Parameters** | Preserve order | N/A (not used) | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:67**
```
66: | **Trailing Dots** | Remove trailing dots | Remove trailing dots | ❌ TODO |
67: | **Query Parameters** | Preserve order | N/A (not used) | ❌ TODO |
68: | **Fragments** | Preserve as-is | Validate if present | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:68**
```
67: | **Query Parameters** | Preserve order | N/A (not used) | ❌ TODO |
68: | **Fragments** | Preserve as-is | Validate if present | ❌ TODO |
69: | **Path Components** | Support dereferencing | Validate but don't use | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:69**
```
68: | **Fragments** | Preserve as-is | Validate if present | ❌ TODO |
69: | **Path Components** | Support dereferencing | Validate but don't use | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:75**
```
74: |----------------|----------|-----------|--------|
75: | **Character Set** | [a-zA-Z0-9.-_] | [a-zA-Z0-9.-_] | ❌ TODO |
76: | **Dot Placement** | No leading/trailing | No leading/trailing | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:76**
```
75: | **Character Set** | [a-zA-Z0-9.-_] | [a-zA-Z0-9.-_] | ❌ TODO |
76: | **Dot Placement** | No leading/trailing | No leading/trailing | ❌ TODO |
77: | **Length Limits** | Accumulate limits | Accumulate limits | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:77**
```
76: | **Dot Placement** | No leading/trailing | No leading/trailing | ❌ TODO |
77: | **Length Limits** | Accumulate limits | Accumulate limits | ❌ TODO |
78: | **Reserved Names** | Check reserved list | Check reserved list | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:78**
```
77: | **Length Limits** | Accumulate limits | Accumulate limits | ❌ TODO |
78: | **Reserved Names** | Check reserved list | Check reserved list | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:86**
```
85: |----------------|----------------|-----------------|--------|
86: | **DID Not Found** | `notFound` (404) | `notFound` (404) | ❌ TODO |
87: | **Invalid DID Syntax** | `invalidDid` (400) | `invalidDid` (400) | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:87**
```
86: | **DID Not Found** | `notFound` (404) | `notFound` (404) | ❌ TODO |
87: | **Invalid DID Syntax** | `invalidDid` (400) | `invalidDid` (400) | ❌ TODO |
88: | **Deactivated DID** | `deactivated` (410) | `conflict` (409) | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:88**
```
87: | **Invalid DID Syntax** | `invalidDid` (400) | `invalidDid` (400) | ❌ TODO |
88: | **Deactivated DID** | `deactivated` (410) | `conflict` (409) | ❌ TODO |
89: | **Unauthorized** | N/A | `unauthorized` (403) | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:89**
```
88: | **Deactivated DID** | `deactivated` (410) | `conflict` (409) | ❌ TODO |
89: | **Unauthorized** | N/A | `unauthorized` (403) | ❌ TODO |
90: | **Invalid Document** | N/A | `invalidDocument` (400) | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:90**
```
89: | **Unauthorized** | N/A | `unauthorized` (403) | ❌ TODO |
90: | **Invalid Document** | N/A | `invalidDocument` (400) | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:96**
```
95: |-------|----------|-----------|-------------------|
96: | **error** | Error code string | Error code string | ❌ TODO |
97: | **message** | Human readable | Human readable | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:97**
```
96: | **error** | Error code string | Error code string | ❌ TODO |
97: | **message** | Human readable | Human readable | ❌ TODO |
98: | **details** | Additional context | Additional context | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:98**
```
97: | **message** | Human readable | Human readable | ❌ TODO |
98: | **details** | Additional context | Additional context | ❌ TODO |
99: | **requestId** | Request identifier | Request identifier | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:99**
```
98: | **details** | Additional context | Additional context | ❌ TODO |
99: | **requestId** | Request identifier | Request identifier | ❌ TODO |
100: | **timestamp** | ISO 8601 | ISO 8601 | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:100**
```
99: | **requestId** | Request identifier | Request identifier | ❌ TODO |
100: | **timestamp** | ISO 8601 | ISO 8601 | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:108**
```
107: |------------|---------------------|---------------------|--------|
108: | **Required Fields** | Validate structure | Validate required | ❌ TODO |
109: | **Field Types** | Type checking | Type checking | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:109**
```
108: | **Required Fields** | Validate structure | Validate required | ❌ TODO |
109: | **Field Types** | Type checking | Type checking | ❌ TODO |
110: | **Value Constraints** | Range/format checks | Range/format checks | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:110**
```
109: | **Field Types** | Type checking | Type checking | ❌ TODO |
110: | **Value Constraints** | Range/format checks | Range/format checks | ❌ TODO |
111: | **Cross-field Validation** | Referential integrity | Referential integrity | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:111**
```
110: | **Value Constraints** | Range/format checks | Range/format checks | ❌ TODO |
111: | **Cross-field Validation** | Referential integrity | Referential integrity | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:117**
```
116: |-------|----------|-----------|--------|
117: | **ID Format** | Must be valid URI | Must be valid URI | ❌ TODO |
118: | **Type Support** | Support AccumulateKeyPage | Support AccumulateKeyPage | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:118**
```
117: | **ID Format** | Must be valid URI | Must be valid URI | ❌ TODO |
118: | **Type Support** | Support AccumulateKeyPage | Support AccumulateKeyPage | ❌ TODO |
119: | **Controller Match** | Must match DID | Must match DID | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:119**
```
118: | **Type Support** | Support AccumulateKeyPage | Support AccumulateKeyPage | ❌ TODO |
119: | **Controller Match** | Must match DID | Must match DID | ❌ TODO |
120: | **Required Properties** | Complete structure | Complete structure | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:120**
```
119: | **Controller Match** | Must match DID | Must match DID | ❌ TODO |
120: | **Required Properties** | Complete structure | Complete structure | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:128**
```
127: |--------|----------|-----------|--------|
128: | **Key Page URL** | Validate format | Enforce auth policy | ❌ TODO |
129: | **Expected Location** | `acc://<adi>/book/1` | `acc://<adi>/book/1` | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:129**
```
128: | **Key Page URL** | Validate format | Enforce auth policy | ❌ TODO |
129: | **Expected Location** | `acc://<adi>/book/1` | `acc://<adi>/book/1` | ❌ TODO |
130: | **Threshold Check** | N/A (read-only) | Verify signature threshold | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:130**
```
129: | **Expected Location** | `acc://<adi>/book/1` | `acc://<adi>/book/1` | ❌ TODO |
130: | **Threshold Check** | N/A (read-only) | Verify signature threshold | ❌ TODO |
131: | **Authorization** | N/A | Validate against policy | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:131**
```
130: | **Threshold Check** | N/A (read-only) | Verify signature threshold | ❌ TODO |
131: | **Authorization** | N/A | Validate against policy | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:137**
```
136: |-------|----------------------|---------------------|--------|
137: | **contentType** | Parse if present | Set to "application/did+ld+json" | ❌ TODO |
138: | **document** | Extract DID document | Wrap DID document | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:138**
```
137: | **contentType** | Parse if present | Set to "application/did+ld+json" | ❌ TODO |
138: | **document** | Extract DID document | Wrap DID document | ❌ TODO |
139: | **meta.versionId** | Use for metadata | Generate unique ID | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:139**
```
138: | **document** | Extract DID document | Wrap DID document | ❌ TODO |
139: | **meta.versionId** | Use for metadata | Generate unique ID | ❌ TODO |
140: | **meta.timestamp** | Use for updated field | Set current time | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:140**
```
139: | **meta.versionId** | Use for metadata | Generate unique ID | ❌ TODO |
140: | **meta.timestamp** | Use for updated field | Set current time | ❌ TODO |
141: | **meta.authorKeyPage** | Validate authority | Set from auth policy | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:141**
```
140: | **meta.timestamp** | Use for updated field | Set current time | ❌ TODO |
141: | **meta.authorKeyPage** | Validate authority | Set from auth policy | ❌ TODO |
142: | **meta.proof** | Verify integrity | Generate proof data | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:142**
```
141: | **meta.authorKeyPage** | Validate authority | Set from auth policy | ❌ TODO |
142: | **meta.proof** | Verify integrity | Generate proof data | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:150**
```
149: |-----------|----------|-----------|--------|
150: | **Format** | Parse timestamp-hash | Generate timestamp-hash | ❌ TODO |
151: | **Timestamp** | Extract Unix timestamp | Use current Unix timestamp | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:151**
```
150: | **Format** | Parse timestamp-hash | Generate timestamp-hash | ❌ TODO |
151: | **Timestamp** | Extract Unix timestamp | Use current Unix timestamp | ❌ TODO |
152: | **Hash Prefix** | Extract first 8 chars | Use first 8 chars of hash | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:152**
```
151: | **Timestamp** | Extract Unix timestamp | Use current Unix timestamp | ❌ TODO |
152: | **Hash Prefix** | Extract first 8 chars | Use first 8 chars of hash | ❌ TODO |
153: | **Uniqueness** | Assume unique | Ensure uniqueness | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:153**
```
152: | **Hash Prefix** | Extract first 8 chars | Use first 8 chars of hash | ❌ TODO |
153: | **Uniqueness** | Assume unique | Ensure uniqueness | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:159**
```
158: |--------|----------|-----------|--------|
159: | **Storage Model** | Read append-only | Write append-only | ❌ TODO |
160: | **Version Time** | Support ?versionTime query | N/A | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:160**
```
159: | **Storage Model** | Read append-only | Write append-only | ❌ TODO |
160: | **Version Time** | Support ?versionTime query | N/A | ❌ TODO |
161: | **Latest Version** | Default to latest | Create new version | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:161**
```
160: | **Version Time** | Support ?versionTime query | N/A | ❌ TODO |
161: | **Latest Version** | Default to latest | Create new version | ❌ TODO |
162: | **Previous Version** | Link to previous | Set previousVersionId | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:162**
```
161: | **Latest Version** | Default to latest | Create new version | ❌ TODO |
162: | **Previous Version** | Link to previous | Set previousVersionId | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:170**
```
169: |--------|----------|-----------|--------|
170: | **application/did+ld+json** | Default output | Default input | ❌ TODO |
171: | **application/ld+json** | Alternative output | Alternative input | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:171**
```
170: | **application/did+ld+json** | Default output | Default input | ❌ TODO |
171: | **application/ld+json** | Alternative output | Alternative input | ❌ TODO |
172: | **application/json** | Fallback output | Fallback input | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:172**
```
171: | **application/ld+json** | Alternative output | Alternative input | ❌ TODO |
172: | **application/json** | Fallback output | Fallback input | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:178**
```
177: |--------|------------------|-------------------|--------|
178: | **Accept** | Respect client preference | N/A | ❌ TODO |
179: | **Content-Type** | Set appropriate type | Validate input type | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:179**
```
178: | **Accept** | Respect client preference | N/A | ❌ TODO |
179: | **Content-Type** | Set appropriate type | Validate input type | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:187**
```
186: |---------------|----------------|-----------------|--------|
187: | **Valid Documents** | Use for resolution | Use for creation | ❌ TODO |
188: | **Invalid Documents** | Return errors | Reject creation | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:188**
```
187: | **Valid Documents** | Use for resolution | Use for creation | ❌ TODO |
188: | **Invalid Documents** | Return errors | Reject creation | ❌ TODO |
189: | **Edge Cases** | Handle gracefully | Validate properly | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:189**
```
188: | **Invalid Documents** | Return errors | Reject creation | ❌ TODO |
189: | **Edge Cases** | Handle gracefully | Validate properly | ❌ TODO |
190: | **Canonical JSON** | Parse correctly | Generate correctly | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:190**
```
189: | **Edge Cases** | Handle gracefully | Validate properly | ❌ TODO |
190: | **Canonical JSON** | Parse correctly | Generate correctly | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:196**
```
195: |------|-------------|--------|
196: | **Create → Resolve** | Register then resolve same document | ❌ TODO |
197: | **Update → Resolve** | Update then resolve latest version | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:197**
```
196: | **Create → Resolve** | Register then resolve same document | ❌ TODO |
197: | **Update → Resolve** | Update then resolve latest version | ❌ TODO |
198: | **Deactivate → Resolve** | Deactivate then resolve shows deactivated | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:198**
```
197: | **Update → Resolve** | Update then resolve latest version | ❌ TODO |
198: | **Deactivate → Resolve** | Deactivate then resolve shows deactivated | ❌ TODO |
199: | **Version History** | Create multiple versions, resolve each | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:199**
```
198: | **Deactivate → Resolve** | Deactivate then resolve shows deactivated | ❌ TODO |
199: | **Version History** | Create multiple versions, resolve each | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:207**
```
206: |---------|----------|-----------|--------|
207: | **API URL** | Same endpoint | Same endpoint | ❌ TODO |
208: | **Timeout** | Read timeout | Write timeout | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:208**
```
207: | **API URL** | Same endpoint | Same endpoint | ❌ TODO |
208: | **Timeout** | Read timeout | Write timeout | ❌ TODO |
209: | **Retry Policy** | Read retries | Write retries | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:209**
```
208: | **Timeout** | Read timeout | Write timeout | ❌ TODO |
209: | **Retry Policy** | Read retries | Write retries | ❌ TODO |
210: | **Authentication** | API credentials | API credentials | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:210**
```
209: | **Retry Policy** | Read retries | Write retries | ❌ TODO |
210: | **Authentication** | API credentials | API credentials | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:216**
```
215: |---------|----------|-----------|--------|
216: | **Max Document Size** | Same limit | Same limit | ❌ TODO |
217: | **Max Array Length** | Same limit | Same limit | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:217**
```
216: | **Max Document Size** | Same limit | Same limit | ❌ TODO |
217: | **Max Array Length** | Same limit | Same limit | ❌ TODO |
218: | **Allowed VM Types** | Same types | Same types | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:218**
```
217: | **Max Array Length** | Same limit | Same limit | ❌ TODO |
218: | **Allowed VM Types** | Same types | Same types | ❌ TODO |
219: | **Service Endpoint Limits** | Same limits | Same limits | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:219**
```
218: | **Allowed VM Types** | Same types | Same types | ❌ TODO |
219: | **Service Endpoint Limits** | Same limits | Same limits | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:227**
```
226: |--------|----------|-----------|--------|
227: | **Request Count** | Track resolutions | Track operations | ❌ TODO |
228: | **Error Rate** | Track resolution errors | Track operation errors | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:228**
```
227: | **Request Count** | Track resolutions | Track operations | ❌ TODO |
228: | **Error Rate** | Track resolution errors | Track operation errors | ❌ TODO |
229: | **Latency** | Resolution time | Operation time | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:229**
```
228: | **Error Rate** | Track resolution errors | Track operation errors | ❌ TODO |
229: | **Latency** | Resolution time | Operation time | ❌ TODO |
230: | **Accumulate Calls** | API call count | API call count | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:230**
```
229: | **Latency** | Resolution time | Operation time | ❌ TODO |
230: | **Accumulate Calls** | API call count | API call count | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:237**
```
236: | **Service Health** | Return 200 OK | Return 200 OK | ✅ DONE |
237: | **Accumulate Connectivity** | Test API connection | Test API connection | ❌ TODO |
238: | **Database Connectivity** | N/A (stateless) | N/A (stateless) | ✅ N/A |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:246**
```
245: |---------------|-------------|--------|
246: | **Create → Resolve** | Registrar creates, resolver resolves | ❌ TODO |
247: | **Update → Resolve** | Registrar updates, resolver gets latest | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:247**
```
246: | **Create → Resolve** | Registrar creates, resolver resolves | ❌ TODO |
247: | **Update → Resolve** | Registrar updates, resolver gets latest | ❌ TODO |
248: | **Deactivate → Resolve** | Registrar deactivates, resolver shows status | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:248**
```
247: | **Update → Resolve** | Registrar updates, resolver gets latest | ❌ TODO |
248: | **Deactivate → Resolve** | Registrar deactivates, resolver shows status | ❌ TODO |
249: | **Error Consistency** | Both services return same errors for same inputs | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:249**
```
248: | **Deactivate → Resolve** | Registrar deactivates, resolver shows status | ❌ TODO |
249: | **Error Consistency** | Both services return same errors for same inputs | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:255**
```
254: |------|-------------|--------|
255: | **Canonical Equivalence** | Same document canonicalizes identically | ❌ TODO |
256: | **Hash Verification** | Registrar hash matches resolver verification | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:256**
```
255: | **Canonical Equivalence** | Same document canonicalizes identically | ❌ TODO |
256: | **Hash Verification** | Registrar hash matches resolver verification | ❌ TODO |
257: | **Metadata Alignment** | Generated metadata matches resolved metadata | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:257**
```
256: | **Hash Verification** | Registrar hash matches resolver verification | ❌ TODO |
257: | **Metadata Alignment** | Generated metadata matches resolved metadata | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:265**
```
264: |---------|---------------|----------------|--------|
265: | **Error Codes** | List all errors | List all errors | ❌ TODO |
266: | **Request Format** | N/A | Document structure | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:266**
```
265: | **Error Codes** | List all errors | List all errors | ❌ TODO |
266: | **Request Format** | N/A | Document structure | ❌ TODO |
267: | **Response Format** | Document structure | Document structure | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:267**
```
266: | **Request Format** | N/A | Document structure | ❌ TODO |
267: | **Response Format** | Document structure | Document structure | ❌ TODO |
268: | **Examples** | Valid requests/responses | Valid requests/responses | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:268**
```
267: | **Response Format** | Document structure | Document structure | ❌ TODO |
268: | **Examples** | Valid requests/responses | Valid requests/responses | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:274**
```
273: |--------------|----------|-----------|--------|
274: | **Basic Usage** | Resolution example | Creation example | ❌ TODO |
275: | **Error Handling** | Error response example | Error response example | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:275**
```
274: | **Basic Usage** | Resolution example | Creation example | ❌ TODO |
275: | **Error Handling** | Error response example | Error response example | ❌ TODO |
276: | **Advanced Features** | Version time example | Update example | ❌ TODO |
```

**spec\PARITY-RESOLVER-REGISTRAR.md:276**
```
275: | **Error Handling** | Error response example | Error response example | ❌ TODO |
276: | **Advanced Features** | Version time example | Update example | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:12**
```
11: |-------------|---------------|----------------------|-------|
12: | **DID URL Dereferencing** | [DID Core §7](https://www.w3.org/TR/did-core/#did-url-dereferencing) | ❌ TODO | Must support path, query, fragment |
13: | **Resolution Result Format** | [DID Core §7.1](https://www.w3.org/TR/did-core/#did-resolution-result) | ❌ TODO | didDocument, didDocumentMetadata, didResolutionMetadata |
```

**spec\PARITY-SPEC-RESOLVER.md:13**
```
12: | **DID URL Dereferencing** | [DID Core §7](https://www.w3.org/TR/did-core/#did-url-dereferencing) | ❌ TODO | Must support path, query, fragment |
13: | **Resolution Result Format** | [DID Core §7.1](https://www.w3.org/TR/did-core/#did-resolution-result) | ❌ TODO | didDocument, didDocumentMetadata, didResolutionMetadata |
14: | **Content Type Handling** | [DID Core §7.1.2](https://www.w3.org/TR/did-core/#did-resolution-metadata) | ❌ TODO | application/did+ld+json default |
```

**spec\PARITY-SPEC-RESOLVER.md:14**
```
13: | **Resolution Result Format** | [DID Core §7.1](https://www.w3.org/TR/did-core/#did-resolution-result) | ❌ TODO | didDocument, didDocumentMetadata, didResolutionMetadata |
14: | **Content Type Handling** | [DID Core §7.1.2](https://www.w3.org/TR/did-core/#did-resolution-metadata) | ❌ TODO | application/did+ld+json default |
15: | **Error Handling** | [DID Core §7.1.2](https://www.w3.org/TR/did-core/#did-resolution-metadata) | ❌ TODO | Standard error codes |
```

**spec\PARITY-SPEC-RESOLVER.md:15**
```
14: | **Content Type Handling** | [DID Core §7.1.2](https://www.w3.org/TR/did-core/#did-resolution-metadata) | ❌ TODO | application/did+ld+json default |
15: | **Error Handling** | [DID Core §7.1.2](https://www.w3.org/TR/did-core/#did-resolution-metadata) | ❌ TODO | Standard error codes |
```

**spec\PARITY-SPEC-RESOLVER.md:21**
```
20: |-------------|---------------|----------------------|-------|
21: | **Created Timestamp** | [DID Core §7.3](https://www.w3.org/TR/did-core/#did-document-metadata) | ❌ TODO | ISO 8601 format |
22: | **Updated Timestamp** | [DID Core §7.3](https://www.w3.org/TR/did-core/#did-document-metadata) | ❌ TODO | ISO 8601 format |
```

**spec\PARITY-SPEC-RESOLVER.md:22**
```
21: | **Created Timestamp** | [DID Core §7.3](https://www.w3.org/TR/did-core/#did-document-metadata) | ❌ TODO | ISO 8601 format |
22: | **Updated Timestamp** | [DID Core §7.3](https://www.w3.org/TR/did-core/#did-document-metadata) | ❌ TODO | ISO 8601 format |
23: | **Version ID** | [DID Core §7.3](https://www.w3.org/TR/did-core/#did-document-metadata) | ❌ TODO | Unique version identifier |
```

**spec\PARITY-SPEC-RESOLVER.md:23**
```
22: | **Updated Timestamp** | [DID Core §7.3](https://www.w3.org/TR/did-core/#did-document-metadata) | ❌ TODO | ISO 8601 format |
23: | **Version ID** | [DID Core §7.3](https://www.w3.org/TR/did-core/#did-document-metadata) | ❌ TODO | Unique version identifier |
24: | **Deactivated Flag** | [DID Core §7.3](https://www.w3.org/TR/did-core/#did-document-metadata) | ❌ TODO | Boolean deactivation status |
```

**spec\PARITY-SPEC-RESOLVER.md:24**
```
23: | **Version ID** | [DID Core §7.3](https://www.w3.org/TR/did-core/#did-document-metadata) | ❌ TODO | Unique version identifier |
24: | **Deactivated Flag** | [DID Core §7.3](https://www.w3.org/TR/did-core/#did-document-metadata) | ❌ TODO | Boolean deactivation status |
```

**spec\PARITY-SPEC-RESOLVER.md:32**
```
31: |-------------|---------------|----------------------|--------------|
32: | **Method Name** | did-acc-method.md §3 | ❌ TODO | Must accept "did:acc:" prefix |
33: | **ADI Name Format** | did-acc-method.md §3.1 | ❌ TODO | Validate ADI syntax rules |
```

**spec\PARITY-SPEC-RESOLVER.md:33**
```
32: | **Method Name** | did-acc-method.md §3 | ❌ TODO | Must accept "did:acc:" prefix |
33: | **ADI Name Format** | did-acc-method.md §3.1 | ❌ TODO | Validate ADI syntax rules |
34: | **Case Insensitive** | did-acc-method.md §6.1 | ❌ TODO | Normalize to lowercase |
```

**spec\PARITY-SPEC-RESOLVER.md:34**
```
33: | **ADI Name Format** | did-acc-method.md §3.1 | ❌ TODO | Validate ADI syntax rules |
34: | **Case Insensitive** | did-acc-method.md §6.1 | ❌ TODO | Normalize to lowercase |
35: | **URL Components** | did-acc-method.md §3.2 | ❌ TODO | Support path, query, fragment |
```

**spec\PARITY-SPEC-RESOLVER.md:35**
```
34: | **Case Insensitive** | did-acc-method.md §6.1 | ❌ TODO | Normalize to lowercase |
35: | **URL Components** | did-acc-method.md §3.2 | ❌ TODO | Support path, query, fragment |
```

**spec\PARITY-SPEC-RESOLVER.md:41**
```
40: |-------------|---------------|----------------------|--------------|
41: | **Data Account Lookup** | did-acc-method.md §4.2 | ❌ TODO | Query acc://<adi>/data/did |
42: | **Version Time Support** | did-acc-method.md §4.2 | ❌ TODO | ?versionTime parameter |
```

**spec\PARITY-SPEC-RESOLVER.md:42**
```
41: | **Data Account Lookup** | did-acc-method.md §4.2 | ❌ TODO | Query acc://<adi>/data/did |
42: | **Version Time Support** | did-acc-method.md §4.2 | ❌ TODO | ?versionTime parameter |
43: | **Latest Version Default** | did-acc-method.md §4.2 | ❌ TODO | Return most recent if no versionTime |
```

**spec\PARITY-SPEC-RESOLVER.md:43**
```
42: | **Version Time Support** | did-acc-method.md §4.2 | ❌ TODO | ?versionTime parameter |
43: | **Latest Version Default** | did-acc-method.md §4.2 | ❌ TODO | Return most recent if no versionTime |
44: | **Deactivated Handling** | did-acc-method.md §4.4 | ❌ TODO | Check deactivated field |
```

**spec\PARITY-SPEC-RESOLVER.md:44**
```
43: | **Latest Version Default** | did-acc-method.md §4.2 | ❌ TODO | Return most recent if no versionTime |
44: | **Deactivated Handling** | did-acc-method.md §4.4 | ❌ TODO | Check deactivated field |
```

**spec\PARITY-SPEC-RESOLVER.md:50**
```
49: |-------------|---------------|----------------------|--------------|
50: | **Type Recognition** | did-acc-method.md §5.1 | ❌ TODO | Handle AccumulateKeyPage type |
51: | **Key Page URL** | did-acc-method.md §5.1 | ❌ TODO | Validate keyPageUrl format |
```

**spec\PARITY-SPEC-RESOLVER.md:51**
```
50: | **Type Recognition** | did-acc-method.md §5.1 | ❌ TODO | Handle AccumulateKeyPage type |
51: | **Key Page URL** | did-acc-method.md §5.1 | ❌ TODO | Validate keyPageUrl format |
52: | **Threshold Property** | did-acc-method.md §5.1 | ❌ TODO | Include threshold value |
```

**spec\PARITY-SPEC-RESOLVER.md:52**
```
51: | **Key Page URL** | did-acc-method.md §5.1 | ❌ TODO | Validate keyPageUrl format |
52: | **Threshold Property** | did-acc-method.md §5.1 | ❌ TODO | Include threshold value |
53: | **Controller Validation** | did-acc-method.md §5.1 | ❌ TODO | Verify controller matches DID |
```

**spec\PARITY-SPEC-RESOLVER.md:53**
```
52: | **Threshold Property** | did-acc-method.md §5.1 | ❌ TODO | Include threshold value |
53: | **Controller Validation** | did-acc-method.md §5.1 | ❌ TODO | Verify controller matches DID |
```

**spec\PARITY-SPEC-RESOLVER.md:61**
```
60: |-------------|---------------|----------------------|--------------|
61: | **Key Ordering** | Rules.md §2.2 | ❌ TODO | Lexicographic order |
62: | **No Whitespace** | Rules.md §2.2 | ❌ TODO | Compact representation |
```

**spec\PARITY-SPEC-RESOLVER.md:62**
```
61: | **Key Ordering** | Rules.md §2.2 | ❌ TODO | Lexicographic order |
62: | **No Whitespace** | Rules.md §2.2 | ❌ TODO | Compact representation |
63: | **Number Format** | Rules.md §2.2 | ❌ TODO | No trailing zeros |
```

**spec\PARITY-SPEC-RESOLVER.md:63**
```
62: | **No Whitespace** | Rules.md §2.2 | ❌ TODO | Compact representation |
63: | **Number Format** | Rules.md §2.2 | ❌ TODO | No trailing zeros |
64: | **String Escaping** | Rules.md §2.2 | ❌ TODO | Minimal escaping |
```

**spec\PARITY-SPEC-RESOLVER.md:64**
```
63: | **Number Format** | Rules.md §2.2 | ❌ TODO | No trailing zeros |
64: | **String Escaping** | Rules.md §2.2 | ❌ TODO | Minimal escaping |
```

**spec\PARITY-SPEC-RESOLVER.md:70**
```
69: |-------------|---------------|----------------------|--------------|
70: | **SHA-256 Algorithm** | Rules.md §3 | ❌ TODO | Use SHA-256 for all hashes |
71: | **Hash Format** | Rules.md §3.3 | ❌ TODO | "sha256:" prefix |
```

**spec\PARITY-SPEC-RESOLVER.md:71**
```
70: | **SHA-256 Algorithm** | Rules.md §3 | ❌ TODO | Use SHA-256 for all hashes |
71: | **Hash Format** | Rules.md §3.3 | ❌ TODO | "sha256:" prefix |
72: | **Content Integrity** | Rules.md §3 | ❌ TODO | Verify stored vs computed hash |
```

**spec\PARITY-SPEC-RESOLVER.md:72**
```
71: | **Hash Format** | Rules.md §3.3 | ❌ TODO | "sha256:" prefix |
72: | **Content Integrity** | Rules.md §3 | ❌ TODO | Verify stored vs computed hash |
```

**spec\PARITY-SPEC-RESOLVER.md:78**
```
77: |-------------|---------------|----------------------|-------------|
78: | **Case Normalization** | Rules.md §8.1 | ❌ TODO | `did:acc:ALICE` → `did:acc:alice` |
79: | **Trailing Dot Removal** | Rules.md §8.1 | ❌ TODO | `did:acc:alice.` → `did:acc:alice` |
```

**spec\PARITY-SPEC-RESOLVER.md:79**
```
78: | **Case Normalization** | Rules.md §8.1 | ❌ TODO | `did:acc:ALICE` → `did:acc:alice` |
79: | **Trailing Dot Removal** | Rules.md §8.1 | ❌ TODO | `did:acc:alice.` → `did:acc:alice` |
80: | **Query Preservation** | Rules.md §8.1 | ❌ TODO | Maintain parameter order |
```

**spec\PARITY-SPEC-RESOLVER.md:80**
```
79: | **Trailing Dot Removal** | Rules.md §8.1 | ❌ TODO | `did:acc:alice.` → `did:acc:alice` |
80: | **Query Preservation** | Rules.md §8.1 | ❌ TODO | Maintain parameter order |
81: | **Fragment Preservation** | Rules.md §8.1 | ❌ TODO | Keep fragment as-is |
```

**spec\PARITY-SPEC-RESOLVER.md:81**
```
80: | **Query Preservation** | Rules.md §8.1 | ❌ TODO | Maintain parameter order |
81: | **Fragment Preservation** | Rules.md §8.1 | ❌ TODO | Keep fragment as-is |
```

**spec\PARITY-SPEC-RESOLVER.md:89**
```
88: |------------|-------------|---------------|----------------------|
89: | `notFound` | 404 | did-acc-method.md §8.1 | ❌ TODO |
90: | `deactivated` | 410 | did-acc-method.md §8.1 | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:90**
```
89: | `notFound` | 404 | did-acc-method.md §8.1 | ❌ TODO |
90: | `deactivated` | 410 | did-acc-method.md §8.1 | ❌ TODO |
91: | `invalidDid` | 400 | did-acc-method.md §8.1 | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:91**
```
90: | `deactivated` | 410 | did-acc-method.md §8.1 | ❌ TODO |
91: | `invalidDid` | 400 | did-acc-method.md §8.1 | ❌ TODO |
92: | `versionNotFound` | 404 | did-acc-method.md §8.1 | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:92**
```
91: | `invalidDid` | 400 | did-acc-method.md §8.1 | ❌ TODO |
92: | `versionNotFound` | 404 | did-acc-method.md §8.1 | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:98**
```
97: |-------|------|----------|----------------------|
98: | `error` | string | ✅ | ❌ TODO |
99: | `message` | string | ✅ | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:99**
```
98: | `error` | string | ✅ | ❌ TODO |
99: | `message` | string | ✅ | ❌ TODO |
100: | `details` | object | ❌ | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:100**
```
99: | `message` | string | ✅ | ❌ TODO |
100: | `details` | object | ❌ | ❌ TODO |
101: | `requestId` | string | ❌ | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:101**
```
100: | `details` | object | ❌ | ❌ TODO |
101: | `requestId` | string | ❌ | ❌ TODO |
102: | `timestamp` | string | ❌ | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:102**
```
101: | `requestId` | string | ❌ | ❌ TODO |
102: | `timestamp` | string | ❌ | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:110**
```
109: |---------|-------------|----------------------|---------------|
110: | **GET /resolve** | Core endpoint | ❌ TODO | ❌ TODO |
111: | **DID Parameter** | ?did=did:acc:alice | ❌ TODO | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:110**
```
109: |---------|-------------|----------------------|---------------|
110: | **GET /resolve** | Core endpoint | ❌ TODO | ❌ TODO |
111: | **DID Parameter** | ?did=did:acc:alice | ❌ TODO | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:111**
```
110: | **GET /resolve** | Core endpoint | ❌ TODO | ❌ TODO |
111: | **DID Parameter** | ?did=did:acc:alice | ❌ TODO | ❌ TODO |
112: | **Version Time** | ?versionTime=2024-01-01T00:00:00Z | ❌ TODO | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:111**
```
110: | **GET /resolve** | Core endpoint | ❌ TODO | ❌ TODO |
111: | **DID Parameter** | ?did=did:acc:alice | ❌ TODO | ❌ TODO |
112: | **Version Time** | ?versionTime=2024-01-01T00:00:00Z | ❌ TODO | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:112**
```
111: | **DID Parameter** | ?did=did:acc:alice | ❌ TODO | ❌ TODO |
112: | **Version Time** | ?versionTime=2024-01-01T00:00:00Z | ❌ TODO | ❌ TODO |
113: | **Accept Header** | Content type negotiation | ❌ TODO | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:112**
```
111: | **DID Parameter** | ?did=did:acc:alice | ❌ TODO | ❌ TODO |
112: | **Version Time** | ?versionTime=2024-01-01T00:00:00Z | ❌ TODO | ❌ TODO |
113: | **Accept Header** | Content type negotiation | ❌ TODO | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:113**
```
112: | **Version Time** | ?versionTime=2024-01-01T00:00:00Z | ❌ TODO | ❌ TODO |
113: | **Accept Header** | Content type negotiation | ❌ TODO | ❌ TODO |
114: | **CORS Support** | Cross-origin requests | ❌ TODO | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:113**
```
112: | **Version Time** | ?versionTime=2024-01-01T00:00:00Z | ❌ TODO | ❌ TODO |
113: | **Accept Header** | Content type negotiation | ❌ TODO | ❌ TODO |
114: | **CORS Support** | Cross-origin requests | ❌ TODO | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:114**
```
113: | **Accept Header** | Content type negotiation | ❌ TODO | ❌ TODO |
114: | **CORS Support** | Cross-origin requests | ❌ TODO | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:114**
```
113: | **Accept Header** | Content type negotiation | ❌ TODO | ❌ TODO |
114: | **CORS Support** | Cross-origin requests | ❌ TODO | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:120**
```
119: |---------|-------------|----------------------|---------------|
120: | **GET /health** | Service health check | ✅ DONE | ❌ TODO |
121: | **Status Response** | JSON status format | ✅ DONE | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:121**
```
120: | **GET /health** | Service health check | ✅ DONE | ❌ TODO |
121: | **Status Response** | JSON status format | ✅ DONE | ❌ TODO |
122: | **Dependency Checks** | Accumulate connectivity | ❌ TODO | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:122**
```
121: | **Status Response** | JSON status format | ✅ DONE | ❌ TODO |
122: | **Dependency Checks** | Accumulate connectivity | ❌ TODO | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:122**
```
121: | **Status Response** | JSON status format | ✅ DONE | ❌ TODO |
122: | **Dependency Checks** | Accumulate connectivity | ❌ TODO | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:130**
```
129: |-----------|-------|----------------|----------------------|
130: | **Uppercase DID** | `did:acc:ALICE` | `did:acc:alice` | ❌ TODO |
131: | **Mixed Case** | `did:acc:Alice.Org` | `did:acc:alice.org` | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:131**
```
130: | **Uppercase DID** | `did:acc:ALICE` | `did:acc:alice` | ❌ TODO |
131: | **Mixed Case** | `did:acc:Alice.Org` | `did:acc:alice.org` | ❌ TODO |
132: | **Trailing Dot** | `did:acc:alice.` | `did:acc:alice` | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:132**
```
131: | **Mixed Case** | `did:acc:Alice.Org` | `did:acc:alice.org` | ❌ TODO |
132: | **Trailing Dot** | `did:acc:alice.` | `did:acc:alice` | ❌ TODO |
133: | **Query Parameters** | `did:acc:alice?versionTime=...` | Preserved | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:133**
```
132: | **Trailing Dot** | `did:acc:alice.` | `did:acc:alice` | ❌ TODO |
133: | **Query Parameters** | `did:acc:alice?versionTime=...` | Preserved | ❌ TODO |
134: | **Fragment** | `did:acc:alice#key-1` | Preserved | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:134**
```
133: | **Query Parameters** | `did:acc:alice?versionTime=...` | Preserved | ❌ TODO |
134: | **Fragment** | `did:acc:alice#key-1` | Preserved | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:140**
```
139: |-----------|-------------|----------------------|
140: | **Basic Resolution** | Resolve existing DID | ❌ TODO |
141: | **Version Time** | Resolve at specific time | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:141**
```
140: | **Basic Resolution** | Resolve existing DID | ❌ TODO |
141: | **Version Time** | Resolve at specific time | ❌ TODO |
142: | **Not Found** | Non-existent DID | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:142**
```
141: | **Version Time** | Resolve at specific time | ❌ TODO |
142: | **Not Found** | Non-existent DID | ❌ TODO |
143: | **Deactivated** | Deactivated DID | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:143**
```
142: | **Not Found** | Non-existent DID | ❌ TODO |
143: | **Deactivated** | Deactivated DID | ❌ TODO |
144: | **Invalid DID** | Malformed DID syntax | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:144**
```
143: | **Deactivated** | Deactivated DID | ❌ TODO |
144: | **Invalid DID** | Malformed DID syntax | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:150**
```
149: |--------|--------|--------------------|----------------------|
150: | **Resolution Latency** | <100ms (cached) | Benchmark tests | ❌ TODO |
151: | **Resolution Latency** | <500ms (uncached) | Benchmark tests | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:151**
```
150: | **Resolution Latency** | <100ms (cached) | Benchmark tests | ❌ TODO |
151: | **Resolution Latency** | <500ms (uncached) | Benchmark tests | ❌ TODO |
152: | **Concurrent Requests** | 1000 req/s | Load testing | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:152**
```
151: | **Resolution Latency** | <500ms (uncached) | Benchmark tests | ❌ TODO |
152: | **Concurrent Requests** | 1000 req/s | Load testing | ❌ TODO |
153: | **Memory Usage** | <100MB baseline | Profiling | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:153**
```
152: | **Concurrent Requests** | 1000 req/s | Load testing | ❌ TODO |
153: | **Memory Usage** | <100MB baseline | Profiling | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:159**
```
158: |-------------|----------------------|---------------------|
159: | **Input Validation** | ❌ TODO | Fuzzing tests |
160: | **XSS Prevention** | ❌ TODO | Security scan |
```

**spec\PARITY-SPEC-RESOLVER.md:160**
```
159: | **Input Validation** | ❌ TODO | Fuzzing tests |
160: | **XSS Prevention** | ❌ TODO | Security scan |
161: | **DoS Protection** | ❌ TODO | Rate limiting |
```

**spec\PARITY-SPEC-RESOLVER.md:161**
```
160: | **XSS Prevention** | ❌ TODO | Security scan |
161: | **DoS Protection** | ❌ TODO | Rate limiting |
162: | **Content Integrity** | ❌ TODO | Hash verification |
```

**spec\PARITY-SPEC-RESOLVER.md:162**
```
161: | **DoS Protection** | ❌ TODO | Rate limiting |
162: | **Content Integrity** | ❌ TODO | Hash verification |
```

**spec\PARITY-SPEC-RESOLVER.md:170**
```
169: |---------|----------------------|---------------|
170: | **API Client** | ❌ TODO | ❌ TODO |
171: | **Data Account Queries** | ❌ TODO | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:170**
```
169: |---------|----------------------|---------------|
170: | **API Client** | ❌ TODO | ❌ TODO |
171: | **Data Account Queries** | ❌ TODO | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:171**
```
170: | **API Client** | ❌ TODO | ❌ TODO |
171: | **Data Account Queries** | ❌ TODO | ❌ TODO |
172: | **Error Handling** | ❌ TODO | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:171**
```
170: | **API Client** | ❌ TODO | ❌ TODO |
171: | **Data Account Queries** | ❌ TODO | ❌ TODO |
172: | **Error Handling** | ❌ TODO | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:172**
```
171: | **Data Account Queries** | ❌ TODO | ❌ TODO |
172: | **Error Handling** | ❌ TODO | ❌ TODO |
173: | **Retry Logic** | ❌ TODO | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:172**
```
171: | **Data Account Queries** | ❌ TODO | ❌ TODO |
172: | **Error Handling** | ❌ TODO | ❌ TODO |
173: | **Retry Logic** | ❌ TODO | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:173**
```
172: | **Error Handling** | ❌ TODO | ❌ TODO |
173: | **Retry Logic** | ❌ TODO | ❌ TODO |
174: | **Circuit Breaker** | ❌ TODO | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:173**
```
172: | **Error Handling** | ❌ TODO | ❌ TODO |
173: | **Retry Logic** | ❌ TODO | ❌ TODO |
174: | **Circuit Breaker** | ❌ TODO | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:174**
```
173: | **Retry Logic** | ❌ TODO | ❌ TODO |
174: | **Circuit Breaker** | ❌ TODO | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:174**
```
173: | **Retry Logic** | ❌ TODO | ❌ TODO |
174: | **Circuit Breaker** | ❌ TODO | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:180**
```
179: |---------|----------------------|---------------|
180: | **Mock Client** | ❌ TODO | ❌ TODO |
181: | **Golden File Tests** | ❌ TODO | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:180**
```
179: |---------|----------------------|---------------|
180: | **Mock Client** | ❌ TODO | ❌ TODO |
181: | **Golden File Tests** | ❌ TODO | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:181**
```
180: | **Mock Client** | ❌ TODO | ❌ TODO |
181: | **Golden File Tests** | ❌ TODO | ❌ TODO |
182: | **Test Vector Validation** | ❌ TODO | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:181**
```
180: | **Mock Client** | ❌ TODO | ❌ TODO |
181: | **Golden File Tests** | ❌ TODO | ❌ TODO |
182: | **Test Vector Validation** | ❌ TODO | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:182**
```
181: | **Golden File Tests** | ❌ TODO | ❌ TODO |
182: | **Test Vector Validation** | ❌ TODO | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:182**
```
181: | **Golden File Tests** | ❌ TODO | ❌ TODO |
182: | **Test Vector Validation** | ❌ TODO | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:188**
```
187: |----------|----------------------|----------------|
188: | **API Documentation** | ❌ TODO | ❌ TODO |
189: | **Usage Examples** | ❌ TODO | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:188**
```
187: |----------|----------------------|----------------|
188: | **API Documentation** | ❌ TODO | ❌ TODO |
189: | **Usage Examples** | ❌ TODO | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:189**
```
188: | **API Documentation** | ❌ TODO | ❌ TODO |
189: | **Usage Examples** | ❌ TODO | ❌ TODO |
190: | **Error Responses** | ❌ TODO | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:189**
```
188: | **API Documentation** | ❌ TODO | ❌ TODO |
189: | **Usage Examples** | ❌ TODO | ❌ TODO |
190: | **Error Responses** | ❌ TODO | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:190**
```
189: | **Usage Examples** | ❌ TODO | ❌ TODO |
190: | **Error Responses** | ❌ TODO | ❌ TODO |
191: | **Configuration Guide** | ❌ TODO | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:190**
```
189: | **Usage Examples** | ❌ TODO | ❌ TODO |
190: | **Error Responses** | ❌ TODO | ❌ TODO |
191: | **Configuration Guide** | ❌ TODO | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:191**
```
190: | **Error Responses** | ❌ TODO | ❌ TODO |
191: | **Configuration Guide** | ❌ TODO | ❌ TODO |
```

**spec\PARITY-SPEC-RESOLVER.md:191**
```
190: | **Error Responses** | ❌ TODO | ❌ TODO |
191: | **Configuration Guide** | ❌ TODO | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:12**
```
11: |-------------|---------------|----------------------|-------|
12: | **Endpoint Path** | GET /1.0/identifiers/{did} | ❌ TODO | Must match exact path |
13: | **Method Support** | GET only | ❌ TODO | No other HTTP methods |
```

**spec\PARITY-UNI-DRIVERS.md:13**
```
12: | **Endpoint Path** | GET /1.0/identifiers/{did} | ❌ TODO | Must match exact path |
13: | **Method Support** | GET only | ❌ TODO | No other HTTP methods |
14: | **DID Parameter** | Path parameter | ❌ TODO | Extract from URL path |
```

**spec\PARITY-UNI-DRIVERS.md:14**
```
13: | **Method Support** | GET only | ❌ TODO | No other HTTP methods |
14: | **DID Parameter** | Path parameter | ❌ TODO | Extract from URL path |
15: | **Response Format** | Universal format | ❌ TODO | Must match Universal spec |
```

**spec\PARITY-UNI-DRIVERS.md:15**
```
14: | **DID Parameter** | Path parameter | ❌ TODO | Extract from URL path |
15: | **Response Format** | Universal format | ❌ TODO | Must match Universal spec |
16: | **Content Type** | application/did+resolution-result+json | ❌ TODO | Default content type |
```

**spec\PARITY-UNI-DRIVERS.md:16**
```
15: | **Response Format** | Universal format | ❌ TODO | Must match Universal spec |
16: | **Content Type** | application/did+resolution-result+json | ❌ TODO | Default content type |
```

**spec\PARITY-UNI-DRIVERS.md:22**
```
21: |---------|----------------|----------------------|---------------|
22: | **DID Validation** | Validate DID syntax | ❌ TODO | ❌ TODO |
23: | **Method Filtering** | Only handle did:acc | ❌ TODO | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:22**
```
21: |---------|----------------|----------------------|---------------|
22: | **DID Validation** | Validate DID syntax | ❌ TODO | ❌ TODO |
23: | **Method Filtering** | Only handle did:acc | ❌ TODO | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:23**
```
22: | **DID Validation** | Validate DID syntax | ❌ TODO | ❌ TODO |
23: | **Method Filtering** | Only handle did:acc | ❌ TODO | ❌ TODO |
24: | **Accept Header** | Support content negotiation | ❌ TODO | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:23**
```
22: | **DID Validation** | Validate DID syntax | ❌ TODO | ❌ TODO |
23: | **Method Filtering** | Only handle did:acc | ❌ TODO | ❌ TODO |
24: | **Accept Header** | Support content negotiation | ❌ TODO | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:24**
```
23: | **Method Filtering** | Only handle did:acc | ❌ TODO | ❌ TODO |
24: | **Accept Header** | Support content negotiation | ❌ TODO | ❌ TODO |
25: | **Query Parameters** | Pass through to core resolver | ❌ TODO | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:24**
```
23: | **Method Filtering** | Only handle did:acc | ❌ TODO | ❌ TODO |
24: | **Accept Header** | Support content negotiation | ❌ TODO | ❌ TODO |
25: | **Query Parameters** | Pass through to core resolver | ❌ TODO | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:25**
```
24: | **Accept Header** | Support content negotiation | ❌ TODO | ❌ TODO |
25: | **Query Parameters** | Pass through to core resolver | ❌ TODO | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:25**
```
24: | **Accept Header** | Support content negotiation | ❌ TODO | ❌ TODO |
25: | **Query Parameters** | Pass through to core resolver | ❌ TODO | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:31**
```
30: |-------|------------------|---------------------|----------------|
31: | **didDocument** | Direct inclusion | Same | ❌ TODO |
32: | **didDocumentMetadata** | Universal format | Same structure | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:32**
```
31: | **didDocument** | Direct inclusion | Same | ❌ TODO |
32: | **didDocumentMetadata** | Universal format | Same structure | ❌ TODO |
33: | **didResolutionMetadata** | Universal format | Compatible | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:33**
```
32: | **didDocumentMetadata** | Universal format | Same structure | ❌ TODO |
33: | **didResolutionMetadata** | Universal format | Compatible | ❌ TODO |
34: | **@context** | Universal context | Convert if needed | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:34**
```
33: | **didResolutionMetadata** | Universal format | Compatible | ❌ TODO |
34: | **@context** | Universal context | Convert if needed | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:42**
```
41: |----------|--------|----------------|----------------------|
42: | **/1.0/create** | POST | Create new DID | ❌ TODO |
43: | **/1.0/update** | POST | Update existing DID | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:43**
```
42: | **/1.0/create** | POST | Create new DID | ❌ TODO |
43: | **/1.0/update** | POST | Update existing DID | ❌ TODO |
44: | **/1.0/deactivate** | POST | Deactivate DID | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:44**
```
43: | **/1.0/update** | POST | Update existing DID | ❌ TODO |
44: | **/1.0/deactivate** | POST | Deactivate DID | ❌ TODO |
45: | **/1.0/resolve** | GET | Optional resolution | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:45**
```
44: | **/1.0/deactivate** | POST | Deactivate DID | ❌ TODO |
45: | **/1.0/resolve** | GET | Optional resolution | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:53**
```
52: |-------|------------------|----------------------|----------------|
53: | **method** | Query parameter "acc" | Internal routing | ❌ TODO |
54: | **options** | Universal options | Convert to internal | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:54**
```
53: | **method** | Query parameter "acc" | Internal routing | ❌ TODO |
54: | **options** | Universal options | Convert to internal | ❌ TODO |
55: | **secret** | Universal secret format | Map to auth | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:55**
```
54: | **options** | Universal options | Convert to internal | ❌ TODO |
55: | **secret** | Universal secret format | Map to auth | ❌ TODO |
56: | **didDocument** | Universal format | Same | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:56**
```
55: | **secret** | Universal secret format | Map to auth | ❌ TODO |
56: | **didDocument** | Universal format | Same | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:62**
```
61: |-------|------------------|----------------------|----------------|
62: | **did** | DID to update | Same | ❌ TODO |
63: | **options** | Universal options | Convert to internal | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:63**
```
62: | **did** | DID to update | Same | ❌ TODO |
63: | **options** | Universal options | Convert to internal | ❌ TODO |
64: | **secret** | Auth credentials | Map to auth | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:64**
```
63: | **options** | Universal options | Convert to internal | ❌ TODO |
64: | **secret** | Auth credentials | Map to auth | ❌ TODO |
65: | **didDocument** | Updated document | Same | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:65**
```
64: | **secret** | Auth credentials | Map to auth | ❌ TODO |
65: | **didDocument** | Updated document | Same | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:71**
```
70: |-------|------------------|----------------------|----------------|
71: | **did** | DID to deactivate | Same | ❌ TODO |
72: | **options** | Universal options | Convert to internal | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:72**
```
71: | **did** | DID to deactivate | Same | ❌ TODO |
72: | **options** | Universal options | Convert to internal | ❌ TODO |
73: | **secret** | Auth credentials | Map to auth | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:73**
```
72: | **options** | Universal options | Convert to internal | ❌ TODO |
73: | **secret** | Auth credentials | Map to auth | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:81**
```
80: |-------|------------------|----------------------|----------------|
81: | **jobId** | Operation tracking | Generate UUID | ❌ TODO |
82: | **didState** | Current DID state | Map from internal | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:82**
```
81: | **jobId** | Operation tracking | Generate UUID | ❌ TODO |
82: | **didState** | Current DID state | Map from internal | ❌ TODO |
83: | **didRegistrationMetadata** | Operation metadata | Convert metadata | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:83**
```
82: | **didState** | Current DID state | Map from internal | ❌ TODO |
83: | **didRegistrationMetadata** | Operation metadata | Convert metadata | ❌ TODO |
84: | **didDocumentMetadata** | Document metadata | Same structure | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:84**
```
83: | **didRegistrationMetadata** | Operation metadata | Convert metadata | ❌ TODO |
84: | **didDocumentMetadata** | Document metadata | Same structure | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:90**
```
89: |------------|------------------|-------------|----------------|
90: | **invalidRequest** | Standard error | Map from 400 | ❌ TODO |
91: | **unauthorized** | Standard error | Map from 403 | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:91**
```
90: | **invalidRequest** | Standard error | Map from 400 | ❌ TODO |
91: | **unauthorized** | Standard error | Map from 403 | ❌ TODO |
92: | **conflict** | Standard error | Map from 409 | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:92**
```
91: | **unauthorized** | Standard error | Map from 403 | ❌ TODO |
92: | **conflict** | Standard error | Map from 409 | ❌ TODO |
93: | **internalError** | Standard error | Map from 500 | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:93**
```
92: | **conflict** | Standard error | Map from 409 | ❌ TODO |
93: | **internalError** | Standard error | Map from 500 | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:101**
```
100: |---------|----------------------|---------------|-------|
101: | **Request Validation** | ❌ TODO | ❌ TODO | Validate Universal format |
102: | **DID Extraction** | ❌ TODO | ❌ TODO | Extract from URL path |
```

**spec\PARITY-UNI-DRIVERS.md:101**
```
100: |---------|----------------------|---------------|-------|
101: | **Request Validation** | ❌ TODO | ❌ TODO | Validate Universal format |
102: | **DID Extraction** | ❌ TODO | ❌ TODO | Extract from URL path |
```

**spec\PARITY-UNI-DRIVERS.md:102**
```
101: | **Request Validation** | ❌ TODO | ❌ TODO | Validate Universal format |
102: | **DID Extraction** | ❌ TODO | ❌ TODO | Extract from URL path |
103: | **Core Service Call** | ❌ TODO | ❌ TODO | HTTP client to resolver |
```

**spec\PARITY-UNI-DRIVERS.md:102**
```
101: | **Request Validation** | ❌ TODO | ❌ TODO | Validate Universal format |
102: | **DID Extraction** | ❌ TODO | ❌ TODO | Extract from URL path |
103: | **Core Service Call** | ❌ TODO | ❌ TODO | HTTP client to resolver |
```

**spec\PARITY-UNI-DRIVERS.md:103**
```
102: | **DID Extraction** | ❌ TODO | ❌ TODO | Extract from URL path |
103: | **Core Service Call** | ❌ TODO | ❌ TODO | HTTP client to resolver |
104: | **Response Mapping** | ❌ TODO | ❌ TODO | Convert to Universal format |
```

**spec\PARITY-UNI-DRIVERS.md:103**
```
102: | **DID Extraction** | ❌ TODO | ❌ TODO | Extract from URL path |
103: | **Core Service Call** | ❌ TODO | ❌ TODO | HTTP client to resolver |
104: | **Response Mapping** | ❌ TODO | ❌ TODO | Convert to Universal format |
```

**spec\PARITY-UNI-DRIVERS.md:104**
```
103: | **Core Service Call** | ❌ TODO | ❌ TODO | HTTP client to resolver |
104: | **Response Mapping** | ❌ TODO | ❌ TODO | Convert to Universal format |
105: | **Error Handling** | ❌ TODO | ❌ TODO | Map error codes |
```

**spec\PARITY-UNI-DRIVERS.md:104**
```
103: | **Core Service Call** | ❌ TODO | ❌ TODO | HTTP client to resolver |
104: | **Response Mapping** | ❌ TODO | ❌ TODO | Convert to Universal format |
105: | **Error Handling** | ❌ TODO | ❌ TODO | Map error codes |
```

**spec\PARITY-UNI-DRIVERS.md:105**
```
104: | **Response Mapping** | ❌ TODO | ❌ TODO | Convert to Universal format |
105: | **Error Handling** | ❌ TODO | ❌ TODO | Map error codes |
```

**spec\PARITY-UNI-DRIVERS.md:105**
```
104: | **Response Mapping** | ❌ TODO | ❌ TODO | Convert to Universal format |
105: | **Error Handling** | ❌ TODO | ❌ TODO | Map error codes |
```

**spec\PARITY-UNI-DRIVERS.md:111**
```
110: |---------|----------------------|---------------|-------|
111: | **Request Validation** | ❌ TODO | ❌ TODO | Validate Universal format |
112: | **Method Filtering** | ❌ TODO | ❌ TODO | Only accept method=acc |
```

**spec\PARITY-UNI-DRIVERS.md:111**
```
110: |---------|----------------------|---------------|-------|
111: | **Request Validation** | ❌ TODO | ❌ TODO | Validate Universal format |
112: | **Method Filtering** | ❌ TODO | ❌ TODO | Only accept method=acc |
```

**spec\PARITY-UNI-DRIVERS.md:112**
```
111: | **Request Validation** | ❌ TODO | ❌ TODO | Validate Universal format |
112: | **Method Filtering** | ❌ TODO | ❌ TODO | Only accept method=acc |
113: | **Request Mapping** | ❌ TODO | ❌ TODO | Convert to core format |
```

**spec\PARITY-UNI-DRIVERS.md:112**
```
111: | **Request Validation** | ❌ TODO | ❌ TODO | Validate Universal format |
112: | **Method Filtering** | ❌ TODO | ❌ TODO | Only accept method=acc |
113: | **Request Mapping** | ❌ TODO | ❌ TODO | Convert to core format |
```

**spec\PARITY-UNI-DRIVERS.md:113**
```
112: | **Method Filtering** | ❌ TODO | ❌ TODO | Only accept method=acc |
113: | **Request Mapping** | ❌ TODO | ❌ TODO | Convert to core format |
114: | **Core Service Call** | ❌ TODO | ❌ TODO | HTTP client to registrar |
```

**spec\PARITY-UNI-DRIVERS.md:113**
```
112: | **Method Filtering** | ❌ TODO | ❌ TODO | Only accept method=acc |
113: | **Request Mapping** | ❌ TODO | ❌ TODO | Convert to core format |
114: | **Core Service Call** | ❌ TODO | ❌ TODO | HTTP client to registrar |
```

**spec\PARITY-UNI-DRIVERS.md:114**
```
113: | **Request Mapping** | ❌ TODO | ❌ TODO | Convert to core format |
114: | **Core Service Call** | ❌ TODO | ❌ TODO | HTTP client to registrar |
115: | **Response Mapping** | ❌ TODO | ❌ TODO | Convert to Universal format |
```

**spec\PARITY-UNI-DRIVERS.md:114**
```
113: | **Request Mapping** | ❌ TODO | ❌ TODO | Convert to core format |
114: | **Core Service Call** | ❌ TODO | ❌ TODO | HTTP client to registrar |
115: | **Response Mapping** | ❌ TODO | ❌ TODO | Convert to Universal format |
```

**spec\PARITY-UNI-DRIVERS.md:115**
```
114: | **Core Service Call** | ❌ TODO | ❌ TODO | HTTP client to registrar |
115: | **Response Mapping** | ❌ TODO | ❌ TODO | Convert to Universal format |
```

**spec\PARITY-UNI-DRIVERS.md:115**
```
114: | **Core Service Call** | ❌ TODO | ❌ TODO | HTTP client to registrar |
115: | **Response Mapping** | ❌ TODO | ❌ TODO | Convert to Universal format |
```

**spec\PARITY-UNI-DRIVERS.md:126**
```
125: | **UNIRESOLVER_DRIVER_DID_ACC_POOLVERSIONS** | N/A | ❌ N/A | N/A |
126: | **CORE_RESOLVER_URL** | Custom | ❌ TODO | http://resolver:8080 |
127: | **CORE_REGISTRAR_URL** | Custom | ❌ TODO | http://registrar:8081 |
```

**spec\PARITY-UNI-DRIVERS.md:127**
```
126: | **CORE_RESOLVER_URL** | Custom | ❌ TODO | http://resolver:8080 |
127: | **CORE_REGISTRAR_URL** | Custom | ❌ TODO | http://registrar:8081 |
```

**spec\PARITY-UNI-DRIVERS.md:133**
```
132: |---------|-------------------|----------------------|-------|
133: | **Port Exposure** | 8080 (resolver), 8081 (registrar) | ❌ TODO | Standard ports |
134: | **Health Checks** | /health endpoint | ❌ TODO | Docker health probes |
```

**spec\PARITY-UNI-DRIVERS.md:134**
```
133: | **Port Exposure** | 8080 (resolver), 8081 (registrar) | ❌ TODO | Standard ports |
134: | **Health Checks** | /health endpoint | ❌ TODO | Docker health probes |
135: | **Labels** | Universal labels | ❌ TODO | Metadata labels |
```

**spec\PARITY-UNI-DRIVERS.md:135**
```
134: | **Health Checks** | /health endpoint | ❌ TODO | Docker health probes |
135: | **Labels** | Universal labels | ❌ TODO | Metadata labels |
136: | **Network** | uni-resolver network | ❌ TODO | Network configuration |
```

**spec\PARITY-UNI-DRIVERS.md:136**
```
135: | **Labels** | Universal labels | ❌ TODO | Metadata labels |
136: | **Network** | uni-resolver network | ❌ TODO | Network configuration |
```

**spec\PARITY-UNI-DRIVERS.md:144**
```
143: |-------------|-------------------|----------------------|--------------|
144: | **Base Image** | Lightweight (Alpine/scratch) | ❌ TODO | Image size check |
145: | **Multi-stage Build** | Build and runtime stages | ❌ TODO | Build optimization |
```

**spec\PARITY-UNI-DRIVERS.md:145**
```
144: | **Base Image** | Lightweight (Alpine/scratch) | ❌ TODO | Image size check |
145: | **Multi-stage Build** | Build and runtime stages | ❌ TODO | Build optimization |
146: | **Security** | Non-root user | ❌ TODO | Security scan |
```

**spec\PARITY-UNI-DRIVERS.md:146**
```
145: | **Multi-stage Build** | Build and runtime stages | ❌ TODO | Build optimization |
146: | **Security** | Non-root user | ❌ TODO | Security scan |
147: | **Labels** | Standard metadata | ❌ TODO | Label validation |
```

**spec\PARITY-UNI-DRIVERS.md:147**
```
146: | **Security** | Non-root user | ❌ TODO | Security scan |
147: | **Labels** | Standard metadata | ❌ TODO | Label validation |
```

**spec\PARITY-UNI-DRIVERS.md:153**
```
152: |---------|-------------------|----------------------|-------|
153: | **Service Names** | driver-did-acc-* | ❌ TODO | Naming convention |
154: | **Network** | uni-resolver | ❌ TODO | Shared network |
```

**spec\PARITY-UNI-DRIVERS.md:154**
```
153: | **Service Names** | driver-did-acc-* | ❌ TODO | Naming convention |
154: | **Network** | uni-resolver | ❌ TODO | Shared network |
155: | **Dependencies** | Core services | ❌ TODO | Service dependencies |
```

**spec\PARITY-UNI-DRIVERS.md:155**
```
154: | **Network** | uni-resolver | ❌ TODO | Shared network |
155: | **Dependencies** | Core services | ❌ TODO | Service dependencies |
156: | **Environment** | Configuration vars | ❌ TODO | Env var passing |
```

**spec\PARITY-UNI-DRIVERS.md:156**
```
155: | **Dependencies** | Core services | ❌ TODO | Service dependencies |
156: | **Environment** | Configuration vars | ❌ TODO | Env var passing |
```

**spec\PARITY-UNI-DRIVERS.md:164**
```
163: |-------------|--------|----------------|-------|
164: | **drivers.json** | ❌ TODO | Create entry | Driver metadata |
165: | **Pattern Matching** | ❌ TODO | did:acc:.* | DID pattern |
```

**spec\PARITY-UNI-DRIVERS.md:165**
```
164: | **drivers.json** | ❌ TODO | Create entry | Driver metadata |
165: | **Pattern Matching** | ❌ TODO | did:acc:.* | DID pattern |
166: | **URL Configuration** | ❌ TODO | Driver endpoint | Service URL |
```

**spec\PARITY-UNI-DRIVERS.md:166**
```
165: | **Pattern Matching** | ❌ TODO | did:acc:.* | DID pattern |
166: | **URL Configuration** | ❌ TODO | Driver endpoint | Service URL |
167: | **Test DID** | ❌ TODO | did:acc:alice | Sample for testing |
```

**spec\PARITY-UNI-DRIVERS.md:167**
```
166: | **URL Configuration** | ❌ TODO | Driver endpoint | Service URL |
167: | **Test DID** | ❌ TODO | did:acc:alice | Sample for testing |
```

**spec\PARITY-UNI-DRIVERS.md:173**
```
172: |-----------|-------------------|----------------------|----------|
173: | **Basic Resolution** | Standard test | ❌ TODO | ❌ TODO |
174: | **Error Handling** | Standard test | ❌ TODO | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:173**
```
172: |-----------|-------------------|----------------------|----------|
173: | **Basic Resolution** | Standard test | ❌ TODO | ❌ TODO |
174: | **Error Handling** | Standard test | ❌ TODO | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:174**
```
173: | **Basic Resolution** | Standard test | ❌ TODO | ❌ TODO |
174: | **Error Handling** | Standard test | ❌ TODO | ❌ TODO |
175: | **Performance** | Standard test | ❌ TODO | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:174**
```
173: | **Basic Resolution** | Standard test | ❌ TODO | ❌ TODO |
174: | **Error Handling** | Standard test | ❌ TODO | ❌ TODO |
175: | **Performance** | Standard test | ❌ TODO | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:175**
```
174: | **Error Handling** | Standard test | ❌ TODO | ❌ TODO |
175: | **Performance** | Standard test | ❌ TODO | ❌ TODO |
176: | **Spec Compliance** | Standard test | ❌ TODO | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:175**
```
174: | **Error Handling** | Standard test | ❌ TODO | ❌ TODO |
175: | **Performance** | Standard test | ❌ TODO | ❌ TODO |
176: | **Spec Compliance** | Standard test | ❌ TODO | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:176**
```
175: | **Performance** | Standard test | ❌ TODO | ❌ TODO |
176: | **Spec Compliance** | Standard test | ❌ TODO | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:176**
```
175: | **Performance** | Standard test | ❌ TODO | ❌ TODO |
176: | **Spec Compliance** | Standard test | ❌ TODO | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:184**
```
183: |-------------|--------|----------------|-------|
184: | **drivers.json** | ❌ TODO | Create entry | Driver metadata |
185: | **Method Support** | ❌ TODO | acc | Method identifier |
```

**spec\PARITY-UNI-DRIVERS.md:185**
```
184: | **drivers.json** | ❌ TODO | Create entry | Driver metadata |
185: | **Method Support** | ❌ TODO | acc | Method identifier |
186: | **Operations** | ❌ TODO | create,update,deactivate | Supported ops |
```

**spec\PARITY-UNI-DRIVERS.md:186**
```
185: | **Method Support** | ❌ TODO | acc | Method identifier |
186: | **Operations** | ❌ TODO | create,update,deactivate | Supported ops |
187: | **Test Configuration** | ❌ TODO | Sample requests | Testing setup |
```

**spec\PARITY-UNI-DRIVERS.md:187**
```
186: | **Operations** | ❌ TODO | create,update,deactivate | Supported ops |
187: | **Test Configuration** | ❌ TODO | Sample requests | Testing setup |
```

**spec\PARITY-UNI-DRIVERS.md:193**
```
192: |-----------|-------------------|----------------------|----------|
193: | **Create Operation** | Standard test | ❌ TODO | ❌ TODO |
194: | **Update Operation** | Standard test | ❌ TODO | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:193**
```
192: |-----------|-------------------|----------------------|----------|
193: | **Create Operation** | Standard test | ❌ TODO | ❌ TODO |
194: | **Update Operation** | Standard test | ❌ TODO | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:194**
```
193: | **Create Operation** | Standard test | ❌ TODO | ❌ TODO |
194: | **Update Operation** | Standard test | ❌ TODO | ❌ TODO |
195: | **Deactivate Operation** | Standard test | ❌ TODO | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:194**
```
193: | **Create Operation** | Standard test | ❌ TODO | ❌ TODO |
194: | **Update Operation** | Standard test | ❌ TODO | ❌ TODO |
195: | **Deactivate Operation** | Standard test | ❌ TODO | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:195**
```
194: | **Update Operation** | Standard test | ❌ TODO | ❌ TODO |
195: | **Deactivate Operation** | Standard test | ❌ TODO | ❌ TODO |
196: | **Error Scenarios** | Standard test | ❌ TODO | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:195**
```
194: | **Update Operation** | Standard test | ❌ TODO | ❌ TODO |
195: | **Deactivate Operation** | Standard test | ❌ TODO | ❌ TODO |
196: | **Error Scenarios** | Standard test | ❌ TODO | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:196**
```
195: | **Deactivate Operation** | Standard test | ❌ TODO | ❌ TODO |
196: | **Error Scenarios** | Standard test | ❌ TODO | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:196**
```
195: | **Deactivate Operation** | Standard test | ❌ TODO | ❌ TODO |
196: | **Error Scenarios** | Standard test | ❌ TODO | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:206**
```
205: | **didDocument** | W3C DID Document | W3C DID Document | ✅ Compatible |
206: | **didDocumentMetadata** | Universal metadata | Acc metadata | ❌ TODO |
207: | **didResolutionMetadata** | Universal metadata | Acc metadata | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:207**
```
206: | **didDocumentMetadata** | Universal metadata | Acc metadata | ❌ TODO |
207: | **didResolutionMetadata** | Universal metadata | Acc metadata | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:213**
```
212: |-------|------------------|------------|---------------------|
213: | **jobId** | UUID string | Generate UUID | ❌ TODO |
214: | **didState** | DID state object | Map from internal | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:214**
```
213: | **jobId** | UUID string | Generate UUID | ❌ TODO |
214: | **didState** | DID state object | Map from internal | ❌ TODO |
215: | **didRegistrationMetadata** | Universal metadata | Convert | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:215**
```
214: | **didState** | DID state object | Map from internal | ❌ TODO |
215: | **didRegistrationMetadata** | Universal metadata | Convert | ❌ TODO |
216: | **didDocumentMetadata** | Universal metadata | Same | ✅ Compatible |
```

**spec\PARITY-UNI-DRIVERS.md:224**
```
223: |------------|-----------------|-------------|----------------|
224: | `notFound` | `notFound` | 404 | ❌ TODO |
225: | `deactivated` | `deactivated` | 410 | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:225**
```
224: | `notFound` | `notFound` | 404 | ❌ TODO |
225: | `deactivated` | `deactivated` | 410 | ❌ TODO |
226: | `invalidDid` | `invalidDid` | 400 | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:226**
```
225: | `deactivated` | `deactivated` | 410 | ❌ TODO |
226: | `invalidDid` | `invalidDid` | 400 | ❌ TODO |
227: | `versionNotFound` | `versionNotFound` | 404 | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:227**
```
226: | `invalidDid` | `invalidDid` | 400 | ❌ TODO |
227: | `versionNotFound` | `versionNotFound` | 404 | ❌ TODO |
228: | `internalError` | `internalError` | 500 | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:228**
```
227: | `versionNotFound` | `versionNotFound` | 404 | ❌ TODO |
228: | `internalError` | `internalError` | 500 | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:234**
```
233: |------------|-----------------|-------------|----------------|
234: | `unauthorized` | `unauthorized` | 403 | ❌ TODO |
235: | `conflict` | `conflict` | 409 | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:235**
```
234: | `unauthorized` | `unauthorized` | 403 | ❌ TODO |
235: | `conflict` | `conflict` | 409 | ❌ TODO |
236: | `invalidDocument` | `invalidRequest` | 400 | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:236**
```
235: | `conflict` | `conflict` | 409 | ❌ TODO |
236: | `invalidDocument` | `invalidRequest` | 400 | ❌ TODO |
237: | `thresholdNotMet` | `unauthorized` | 403 | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:237**
```
236: | `invalidDocument` | `invalidRequest` | 400 | ❌ TODO |
237: | `thresholdNotMet` | `unauthorized` | 403 | ❌ TODO |
238: | `internalError` | `internalError` | 500 | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:238**
```
237: | `thresholdNotMet` | `unauthorized` | 403 | ❌ TODO |
238: | `internalError` | `internalError` | 500 | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:246**
```
245: |---------------|----------------|------------------|--------|
246: | **Request Parsing** | HTTP request handling | HTTP request handling | ❌ TODO |
247: | **Response Mapping** | Format conversion | Format conversion | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:247**
```
246: | **Request Parsing** | HTTP request handling | HTTP request handling | ❌ TODO |
247: | **Response Mapping** | Format conversion | Format conversion | ❌ TODO |
248: | **Error Handling** | Error scenarios | Error scenarios | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:248**
```
247: | **Response Mapping** | Format conversion | Format conversion | ❌ TODO |
248: | **Error Handling** | Error scenarios | Error scenarios | ❌ TODO |
249: | **Validation** | Input validation | Input validation | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:249**
```
248: | **Error Handling** | Error scenarios | Error scenarios | ❌ TODO |
249: | **Validation** | Input validation | Input validation | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:255**
```
254: |-----------|-------------|--------|
255: | **End-to-End** | Universal → Driver → Core → Driver → Universal | ❌ TODO |
256: | **Error Propagation** | Error handling through full stack | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:256**
```
255: | **End-to-End** | Universal → Driver → Core → Driver → Universal | ❌ TODO |
256: | **Error Propagation** | Error handling through full stack | ❌ TODO |
257: | **Performance** | Latency and throughput | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:257**
```
256: | **Error Propagation** | Error handling through full stack | ❌ TODO |
257: | **Performance** | Latency and throughput | ❌ TODO |
258: | **Compatibility** | Universal framework tests | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:258**
```
257: | **Performance** | Latency and throughput | ❌ TODO |
258: | **Compatibility** | Universal framework tests | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:264**
```
263: |------|-------------|----------|--------|
264: | **Basic Resolution** | Resolve test DID | Windows (PS1) | ❌ TODO |
265: | **Basic Resolution** | Resolve test DID | Unix (SH) | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:265**
```
264: | **Basic Resolution** | Resolve test DID | Windows (PS1) | ❌ TODO |
265: | **Basic Resolution** | Resolve test DID | Unix (SH) | ❌ TODO |
266: | **Create Operation** | Create test DID | Windows (PS1) | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:266**
```
265: | **Basic Resolution** | Resolve test DID | Unix (SH) | ❌ TODO |
266: | **Create Operation** | Create test DID | Windows (PS1) | ❌ TODO |
267: | **Create Operation** | Create test DID | Unix (SH) | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:267**
```
266: | **Create Operation** | Create test DID | Windows (PS1) | ❌ TODO |
267: | **Create Operation** | Create test DID | Unix (SH) | ❌ TODO |
268: | **Docker Health** | Container health checks | Both | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:268**
```
267: | **Create Operation** | Create test DID | Unix (SH) | ❌ TODO |
268: | **Docker Health** | Container health checks | Both | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:276**
```
275: |--------|-------------------|----------------------|-------|
276: | **Request Count** | HTTP requests/sec | ❌ TODO | Prometheus format |
277: | **Response Time** | Latency percentiles | ❌ TODO | Histogram |
```

**spec\PARITY-UNI-DRIVERS.md:277**
```
276: | **Request Count** | HTTP requests/sec | ❌ TODO | Prometheus format |
277: | **Response Time** | Latency percentiles | ❌ TODO | Histogram |
278: | **Error Rate** | Error percentage | ❌ TODO | By error type |
```

**spec\PARITY-UNI-DRIVERS.md:278**
```
277: | **Response Time** | Latency percentiles | ❌ TODO | Histogram |
278: | **Error Rate** | Error percentage | ❌ TODO | By error type |
279: | **Core Service Calls** | Upstream calls | ❌ TODO | Dependency tracking |
```

**spec\PARITY-UNI-DRIVERS.md:279**
```
278: | **Error Rate** | Error percentage | ❌ TODO | By error type |
279: | **Core Service Calls** | Upstream calls | ❌ TODO | Dependency tracking |
```

**spec\PARITY-UNI-DRIVERS.md:285**
```
284: |-------|-------------------|----------------------|-------|
285: | **Driver Health** | /health endpoint | ❌ TODO | Driver status |
286: | **Core Service Health** | Upstream health | ❌ TODO | Dependency check |
```

**spec\PARITY-UNI-DRIVERS.md:286**
```
285: | **Driver Health** | /health endpoint | ❌ TODO | Driver status |
286: | **Core Service Health** | Upstream health | ❌ TODO | Dependency check |
287: | **Docker Health** | Container health | ❌ TODO | Docker integration |
```

**spec\PARITY-UNI-DRIVERS.md:287**
```
286: | **Core Service Health** | Upstream health | ❌ TODO | Dependency check |
287: | **Docker Health** | Container health | ❌ TODO | Docker integration |
```

**spec\PARITY-UNI-DRIVERS.md:295**
```
294: |----------|-------------|--------|-------|
295: | **Driver README** | Setup instructions | ❌ TODO | How to run |
296: | **Configuration** | Environment variables | ❌ TODO | All options |
```

**spec\PARITY-UNI-DRIVERS.md:296**
```
295: | **Driver README** | Setup instructions | ❌ TODO | How to run |
296: | **Configuration** | Environment variables | ❌ TODO | All options |
297: | **API Examples** | Sample requests/responses | ❌ TODO | Working examples |
```

**spec\PARITY-UNI-DRIVERS.md:297**
```
296: | **Configuration** | Environment variables | ❌ TODO | All options |
297: | **API Examples** | Sample requests/responses | ❌ TODO | Working examples |
298: | **Troubleshooting** | Common issues | ❌ TODO | Debug guide |
```

**spec\PARITY-UNI-DRIVERS.md:298**
```
297: | **API Examples** | Sample requests/responses | ❌ TODO | Working examples |
298: | **Troubleshooting** | Common issues | ❌ TODO | Debug guide |
```

**spec\PARITY-UNI-DRIVERS.md:304**
```
303: |----------|-------------|--------|-------|
304: | **Driver README** | Setup instructions | ❌ TODO | How to run |
305: | **Configuration** | Environment variables | ❌ TODO | All options |
```

**spec\PARITY-UNI-DRIVERS.md:305**
```
304: | **Driver README** | Setup instructions | ❌ TODO | How to run |
305: | **Configuration** | Environment variables | ❌ TODO | All options |
306: | **API Examples** | Sample requests/responses | ❌ TODO | Working examples |
```

**spec\PARITY-UNI-DRIVERS.md:306**
```
305: | **Configuration** | Environment variables | ❌ TODO | All options |
306: | **API Examples** | Sample requests/responses | ❌ TODO | Working examples |
307: | **Auth Guide** | Secret/credential format | ❌ TODO | Authentication |
```

**spec\PARITY-UNI-DRIVERS.md:307**
```
306: | **API Examples** | Sample requests/responses | ❌ TODO | Working examples |
307: | **Auth Guide** | Secret/credential format | ❌ TODO | Authentication |
```

**spec\PARITY-UNI-DRIVERS.md:315**
```
314: |-----------|-------------------|--------|-------------------|
315: | **Resolve** | <500ms | <300ms (including core) | ❌ TODO |
316: | **Create** | <2s | <1s (including core) | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:316**
```
315: | **Resolve** | <500ms | <300ms (including core) | ❌ TODO |
316: | **Create** | <2s | <1s (including core) | ❌ TODO |
317: | **Update** | <2s | <1s (including core) | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:317**
```
316: | **Create** | <2s | <1s (including core) | ❌ TODO |
317: | **Update** | <2s | <1s (including core) | ❌ TODO |
318: | **Deactivate** | <2s | <1s (including core) | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:318**
```
317: | **Update** | <2s | <1s (including core) | ❌ TODO |
318: | **Deactivate** | <2s | <1s (including core) | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:324**
```
323: |--------|-------------------|--------|-------------------|
324: | **Concurrent Requests** | 100 req/s | 100 req/s | ❌ TODO |
325: | **Memory Usage** | <100MB | <50MB | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:325**
```
324: | **Concurrent Requests** | 100 req/s | 100 req/s | ❌ TODO |
325: | **Memory Usage** | <100MB | <50MB | ❌ TODO |
326: | **CPU Usage** | <50% | <25% | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:326**
```
325: | **Memory Usage** | <100MB | <50MB | ❌ TODO |
326: | **CPU Usage** | <50% | <25% | ❌ TODO |
```

**spec\PARITY-UNI-DRIVERS.md:334**
```
333: |-------------|--------|----------------|-------|
334: | **Input Validation** | ❌ TODO | Validate all inputs | Prevent injection |
335: | **Rate Limiting** | ❌ TODO | Implement rate limits | DoS protection |
```

**spec\PARITY-UNI-DRIVERS.md:335**
```
334: | **Input Validation** | ❌ TODO | Validate all inputs | Prevent injection |
335: | **Rate Limiting** | ❌ TODO | Implement rate limits | DoS protection |
336: | **CORS Headers** | ❌ TODO | Proper CORS setup | Browser security |
```

**spec\PARITY-UNI-DRIVERS.md:336**
```
335: | **Rate Limiting** | ❌ TODO | Implement rate limits | DoS protection |
336: | **CORS Headers** | ❌ TODO | Proper CORS setup | Browser security |
337: | **Security Headers** | ❌ TODO | Standard headers | HTTP security |
```

**spec\PARITY-UNI-DRIVERS.md:337**
```
336: | **CORS Headers** | ❌ TODO | Proper CORS setup | Browser security |
337: | **Security Headers** | ❌ TODO | Standard headers | HTTP security |
```

**spec\PARITY-UNI-DRIVERS.md:343**
```
342: |-------------|--------|----------------|-------|
343: | **Non-root User** | ❌ TODO | Run as non-root | Privilege escalation |
344: | **Minimal Base** | ❌ TODO | Distroless/Alpine | Attack surface |
```

**spec\PARITY-UNI-DRIVERS.md:344**
```
343: | **Non-root User** | ❌ TODO | Run as non-root | Privilege escalation |
344: | **Minimal Base** | ❌ TODO | Distroless/Alpine | Attack surface |
345: | **Vulnerability Scan** | ❌ TODO | Container scanning | CVE detection |
```

**spec\PARITY-UNI-DRIVERS.md:345**
```
344: | **Minimal Base** | ❌ TODO | Distroless/Alpine | Attack surface |
345: | **Vulnerability Scan** | ❌ TODO | Container scanning | CVE detection |
346: | **Secret Management** | ❌ TODO | External secrets | No hardcoded secrets |
```

**spec\PARITY-UNI-DRIVERS.md:346**
```
345: | **Vulnerability Scan** | ❌ TODO | Container scanning | CVE detection |
346: | **Secret Management** | ❌ TODO | External secrets | No hardcoded secrets |
```

#### tools\todoscan

**tools\todoscan\main.go:18**
```
18: // TodoItem represents a single TODO/FIXME/etc finding
19: type TodoItem struct {
```

**tools\todoscan\main.go:47**
```
46: var (
47: // TODO patterns to search for (case-insensitive)
48: todoPatterns = []string{
```

**tools\todoscan\main.go:49**
```
48: todoPatterns = []string{
49: `(?i)\b(TODO|FIXME|XXX|HACK|STUB|TBA|TBD|NOTIMPL|NOTIMPLEMENTED)\b`,
50: `(?i)PANIC\s*\(\s*["']TODO`,
```

**tools\todoscan\main.go:50**
```
49: `(?i)\b(TODO|FIXME|XXX|HACK|STUB|TBA|TBD|NOTIMPL|NOTIMPLEMENTED)\b`,
50: `(?i)PANIC\s*\(\s*["']TODO`,
51: `(?i)@deprecated`,
```

**tools\todoscan\main.go:166**
```
166: // Scan file for TODO items
167: fileItems, err := scanFile(path)
```

**tools\todoscan\main.go:275**
```
275: if strings.Contains(match, "TODO") {
276: return "TODO"
```

**tools\todoscan\main.go:276**
```
275: if strings.Contains(match, "TODO") {
276: return "TODO"
277: }
```

**tools\todoscan\main.go:300**
```
299: if strings.Contains(match, "PANIC") {
300: return "PANIC-TODO"
301: }
```

**tools\todoscan\main.go:346**
```
345: func generateJSONReport(report TodoReport) error {
346: file, err := os.Create("reports/todo-report.json")
347: if err != nil {
```

**tools\todoscan\main.go:358**
```
357: func generateMarkdownReport(report TodoReport) error {
358: file, err := os.Create("reports/todo-report.md")
359: if err != nil {
```

**tools\todoscan\main.go:364**
```
364: fmt.Fprintf(file, "# TODO Scan Report\n\n")
365: fmt.Fprintf(file, "**Generated:** %s\n", report.GeneratedAt.Format(time.RFC3339))
```

**tools\todoscan\main.go:478**
```
477: func generateCSVReport(report TodoReport) error {
478: file, err := os.Create("reports/todo-report.csv")
479: if err != nil {
```

### XXX (13 items)

#### docs\ops

**docs\ops\OPERATIONS.md:735**
```
734: | `FIXME` | Known bugs/issues | High |
735: | `XXX` | Code requiring attention | High |
736: | `HACK` | Temporary workarounds | Medium |
```

**docs\ops\OPERATIONS.md:759**
```
758: # 2. Review critical items
759: grep -E "(FIXME|XXX|PANIC)" reports/todo-report.md
```

**docs\ops\OPERATIONS.md:798**
```
797: **Escalation Criteria:**
798: - **FIXME/XXX items**: Convert to GitHub issues if affecting operations
799: - **NOTIMPL items**: Add to `spec/BACKLOG.md` if blocking features
```

#### root

**CLAUDE.md:223**
```
222: - **FIXME**: Bugs that need fixing
223: - **XXX**: Code that needs attention
224: - **HACK**: Temporary workarounds
```

**CLAUDE.md:243**
```
242: # High-priority items
243: grep -E "(FIXME|XXX|PANIC)" reports/todo-report.md
```

**CLAUDE.md:273**
```
272: # High-priority items with context
273: jq '.items[] | select(.tag == "FIXME" or .tag == "XXX")' reports/todo-report.json
274: ```
```

**CLAUDE.md:310**
```
309: - Use `DEPRECATED` when phasing out code
310: - Avoid `XXX` except for urgent attention-required items
```

**Makefile:282**
```
281: @echo "🔍 Code analysis:"
282: @echo "  todo-scan       - Scan repository for TODO/FIXME/XXX markers"
283: @echo ""
```

#### scripts

**scripts\todo-scan.ps1:2**
```
1: # TODO Scanner - Windows PowerShell Wrapper
2: # Scans the repository for TODO, FIXME, XXX, HACK, and other markers
3: # Generates reports in JSON, Markdown, and CSV formats
```

**scripts\todo-scan.sh:4**
```
3: # TODO Scanner - Linux/Docker Wrapper
4: # Scans the repository for TODO, FIXME, XXX, HACK, and other markers
5: # Generates reports in JSON, Markdown, and CSV formats
```

#### tools\todoscan

**tools\todoscan\main.go:49**
```
48: todoPatterns = []string{
49: `(?i)\b(TODO|FIXME|XXX|HACK|STUB|TBA|TBD|NOTIMPL|NOTIMPLEMENTED)\b`,
50: `(?i)PANIC\s*\(\s*["']TODO`,
```

**tools\todoscan\main.go:281**
```
280: }
281: if strings.Contains(match, "XXX") {
282: return "XXX"
```

**tools\todoscan\main.go:282**
```
281: if strings.Contains(match, "XXX") {
282: return "XXX"
283: }
```

