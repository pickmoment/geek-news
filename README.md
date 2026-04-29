# geek-news

[긱뉴스(GeekNews)](https://news.hada.io) CLI 도구. 인기/최신 토픽 목록과 토픽 본문 및 댓글을 터미널에서 조회합니다.

## 설치

```bash
go install github.com/pickmoment/geek-news@latest
```

또는 git으로 클론 후 빌드:

```bash
git clone https://github.com/pickmoment/geek-news.git
cd geek-news
go build -o geek-news .
sudo mv geek-news /usr/local/bin/   # PATH에 추가
```

## 사용법

```
geek-news <명령어> [옵션]
```

### list — 인기 토픽 목록

```bash
geek-news list [-p 페이지] [-n 개수] [-f json|telegram]
```

메인 페이지의 인기 토픽(포인트 순)을 반환합니다.

```bash
geek-news list           # 1페이지 전체
geek-news list -p 2      # 2페이지
geek-news list -n 10     # 상위 10개
```

### new — 최신 토픽 목록

```bash
geek-news new [-p 페이지] [-n 개수] [-f json|telegram]
```

등록 시간 순 최신 토픽을 반환합니다.

```bash
geek-news new
geek-news new -n 5
```

### view — 토픽 상세 및 댓글

```bash
geek-news view <topic_id> [-f json|telegram]
```

토픽 본문 전체와 댓글 트리를 반환합니다. `topic_id`는 목록의 `id` 필드 값 또는 `https://news.hada.io/topic?id=<id>` URL에서 추출합니다.

```bash
geek-news view 28861
```

### install — AI 에이전트 스킬 설치

```bash
geek-news install
```

Claude Code 또는 Codex에서 `/geek-news`로 호출할 수 있는 스킬 파일을 설치합니다. 인터랙티브하게 에이전트(Claude Code / Codex)와 설치 범위(global / project)를 선택합니다.

## 옵션

| 옵션 | 기본값 | 설명 |
|------|--------|------|
| `-p N` | 1 | 페이지 번호 |
| `-n N` | 전체 | 최대 항목 수 |
| `-f` | `json` | 출력 형식 (`json` \| `telegram`) |

## 출력 형식

**JSON** (기본): 구조화된 JSON 출력

```json
[
  {
    "id": "28861",
    "title": "...",
    "url": "https://...",
    "points": 42,
    "author": "username",
    "time": "3시간 전",
    "comments": "댓글 12개",
    "topic_url": "https://news.hada.io/topic?id=28861"
  }
]
```

**Telegram** (`-f telegram`): 코드블록 형태의 텍스트 출력

## 참고

- API 없이 HTML을 직접 파싱합니다. 사이트 구조 변경 시 동작이 달라질 수 있습니다.
- 댓글 `depth`는 0이 최상위, 숫자가 클수록 대댓글입니다.
