package main

import (
	"flag"
	"fmt"
	"os"
)

func usage() {
	fmt.Fprintln(os.Stderr, `GeekNews (긱뉴스) CLI

사용법: geek-news <명령어> [옵션]

명령어:
  list     인기 토픽 목록 조회
  new      최신 토픽 목록 조회
  view     토픽 상세 및 댓글 조회
  install  AI 에이전트용 스킬 파일 설치

각 명령어에 -h 플래그로 도움말을 볼 수 있습니다.`)
}

func errExit(msg string) {
	fmt.Fprintf(os.Stderr, "오류: %s\n", msg)
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	switch os.Args[1] {
	case "list":
		runList(os.Args[2:])
	case "new":
		runNew(os.Args[2:])
	case "view":
		runView(os.Args[2:])
	case "install":
		runInstall(os.Args[2:])
	case "-h", "--help", "help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "알 수 없는 명령어: %s\n\n", os.Args[1])
		usage()
		os.Exit(1)
	}
}

func addCommonFlags(fs *flag.FlagSet) (*int, *int, *string) {
	page := fs.Int("p", 1, "페이지 번호")
	limit := fs.Int("n", 0, "최대 항목 수 (기본: 전체)")
	format := fs.String("f", "json", "출력 형식 (json|telegram)")
	return page, limit, format
}

func runList(args []string) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "사용법: geek-news list [-p 페이지] [-n 개수] [-f json|telegram]\n")
		fs.PrintDefaults()
	}
	page, limit, format := addCommonFlags(fs)
	_ = fs.Parse(args)

	data, err := topics("/", *page)
	if err != nil {
		errExit(err.Error())
	}
	if *limit > 0 && *limit < len(data) {
		data = data[:*limit]
	}
	fmt.Println(fmtTopics(data, *format))
}

func runNew(args []string) {
	fs := flag.NewFlagSet("new", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "사용법: geek-news new [-p 페이지] [-n 개수] [-f json|telegram]\n")
		fs.PrintDefaults()
	}
	page, limit, format := addCommonFlags(fs)
	_ = fs.Parse(args)

	data, err := topics("/new", *page)
	if err != nil {
		errExit(err.Error())
	}
	if *limit > 0 && *limit < len(data) {
		data = data[:*limit]
	}
	fmt.Println(fmtTopics(data, *format))
}

func runView(args []string) {
	fs := flag.NewFlagSet("view", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "사용법: geek-news view <topic_id> [-f json|telegram]\n")
		fs.PrintDefaults()
	}
	format := fs.String("f", "json", "출력 형식 (json|telegram)")
	_ = fs.Parse(args)

	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "오류: topic_id 가 필요합니다")
		fs.Usage()
		os.Exit(1)
	}
	data, err := topic(fs.Arg(0))
	if err != nil {
		errExit(err.Error())
	}
	fmt.Println(fmtTopic(data, *format))
}
