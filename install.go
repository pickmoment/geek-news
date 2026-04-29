package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const skillName = "geek-news"

const skillContent = `---
name: geek-news
description: GeekNews(긱뉴스) 데이터를 조회하는 스킬. 인기/최신 토픽 목록, 토픽 본문과 댓글을 제공한다. 사용자가 긱뉴스 정보를 묻거나 /geek-news를 입력했을 때 사용.
---

# geek-news

news.hada.io(긱뉴스)를 스크래핑하는 CLI 도구 ` + "`geek-news`" + `을 실행해 토픽 데이터를 가져오는 스킬.

## CLI 위치

` + "`geek-news`" + ` 바이너리가 PATH에 설치되어 있어야 합니다.

모든 명령은 ` + "`geek-news <subcommand>`" + ` 형태로 실행합니다.

## 서브커맨드

### list — 인기 토픽 목록

` + "```bash" + `
geek-news list [-p N] [-n N] [-f json|telegram]
` + "```" + `

- 메인 페이지의 인기 토픽(포인트 순) 반환
- ` + "`-p`" + `: 페이지 번호 (기본 1, 페이지당 20개)
- ` + "`-n`" + `: 최대 항목 수 (기본: 전체)
- ` + "`-f`" + `: 출력 형식 (기본: json)
- 반환 필드: id, title, url, domain, desc, points, author, time, comments, topic_url

` + "```bash" + `
geek-news list
geek-news list -p 2
geek-news list -n 10
` + "```" + `

### new — 최신 토픽 목록

` + "```bash" + `
geek-news new [-p N] [-n N] [-f json|telegram]
` + "```" + `

- 등록 시간 순 최신 토픽 반환
- 옵션은 ` + "`list`" + `와 동일

` + "```bash" + `
geek-news new
geek-news new -n 5
` + "```" + `

### view — 토픽 상세 및 댓글

` + "```bash" + `
geek-news view <topic_id> [-f json|telegram]
` + "```" + `

- 토픽 본문 전체 + 댓글 트리 반환
- topic_id: 목록에서 ` + "`id`" + ` 필드의 숫자값
- 반환 필드: id, title, url, points, author, time, body, comments(id/depth/author/time/text), topic_url

` + "```bash" + `
geek-news view 28861
geek-news view 28873
` + "```" + `

## 출력 형식

- ` + "`-f json`" + ` (기본): JSON 출력
- ` + "`-f telegram`" + `: 코드블록(` + "``` ```" + `) 형태

사용자에게 보여줄 때는 ` + "`-f json`" + `으로 데이터를 받아 자연어로 정리한다.

## 사용 패턴

**"긱뉴스 인기 토픽 보여줘" / "요즘 핫한 개발 뉴스"**
` + "```bash" + `
geek-news list
` + "```" + `

**"긱뉴스 최신 글 뭐 있어?"**
` + "```bash" + `
geek-news new
` + "```" + `

**"상위 5개만 보여줘"**
` + "```bash" + `
geek-news list -n 5
` + "```" + `

**"이 토픽 본문이랑 댓글 보여줘" (ID 또는 URL에서 ID 추출)**
` + "```" + `
URL: https://news.hada.io/topic?id=28861
                                      ^^^^^
→ geek-news view 28861
` + "```" + `
` + "```bash" + `
geek-news view 28861
` + "```" + `

**"2페이지 토픽 보여줘"**
` + "```bash" + `
geek-news list -p 2
` + "```" + `

**토픽 본문 요약 요청 시 흐름**
` + "```bash" + `
# 1. 목록에서 ID 확인
geek-news list -n 20

# 2. 해당 ID로 상세 조회
geek-news view <id>
` + "```" + `

## 주의사항

- API 없이 HTML 직접 파싱 (구조 변경 시 동작 이상 가능)
- 댓글 depth는 0이 최상위, 숫자가 클수록 대댓글입니다.
`

var stdin = bufio.NewReader(os.Stdin)

func ask(question string, choices []string) string {
	for {
		fmt.Print(question)
		line, _ := stdin.ReadString('\n')
		ans := strings.TrimSpace(line)
		for _, c := range choices {
			if ans == c {
				return ans
			}
		}
		fmt.Printf("  %s 중 하나를 입력하세요.\n", strings.Join(choices, "/"))
	}
}

func confirm(question string) bool {
	fmt.Print(question)
	line, _ := stdin.ReadString('\n')
	ans := strings.TrimSpace(strings.ToLower(line))
	return ans == "" || ans == "y" || ans == "yes"
}

type installTarget struct {
	agent string
	scope string
	dir   string
}

var targets = []installTarget{
	{"Claude Code", "global", filepath.Join(os.Getenv("HOME"), ".claude", "skills")},
	{"Claude Code", "project", ".claude/skills"},
	{"Codex", "global", filepath.Join(os.Getenv("HOME"), ".agents", "skills")},
	{"Codex", "project", ".agents/skills"},
}

func runInstall(_ []string) {
	fmt.Println("에이전트를 선택하세요:")
	fmt.Println("  1) Claude Code")
	fmt.Println("  2) Codex")
	agent := ask("선택 [1/2]: ", []string{"1", "2"})

	fmt.Println("\n설치 범위를 선택하세요:")
	fmt.Println("  1) global  — 모든 프로젝트에서 사용")
	fmt.Println("  2) project — 현재 프로젝트에서만 사용")
	scope := ask("선택 [1/2]: ", []string{"1", "2"})

	idx := (atoi(agent)-1)*2 + (atoi(scope) - 1)
	t := targets[idx]

	dest := filepath.Join(t.dir, skillName, "SKILL.md")
	fmt.Printf("\n설치 위치: %s\n", dest)
	if !confirm("설치할까요? [Y/n]: ") {
		fmt.Println("취소했습니다.")
		return
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		errExit("디렉토리 생성 실패: " + err.Error())
	}
	if err := os.WriteFile(dest, []byte(skillContent), 0644); err != nil {
		errExit("파일 쓰기 실패: " + err.Error())
	}
	fmt.Printf("\n스킬 설치 완료: %s\n", dest)
	fmt.Printf("%s에서 /%s 으로 호출할 수 있습니다.\n", t.agent, skillName)
}

func atoi(s string) int {
	if s == "2" {
		return 2
	}
	return 1
}
