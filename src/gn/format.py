import json


def _tg(text: str) -> str:
    return f"```\n{text}\n```"


def fmt_topics(data: list, fmt: str) -> str:
    if fmt == "json":
        return json.dumps(data, ensure_ascii=False, indent=2)
    if not data:
        return _tg("토픽이 없습니다.")
    lines = []
    for item in data:
        pts = item.get("points", 0)
        cmt = item.get("comments", "")
        meta = f"{pts}P  {item['author']}  {item['time']}"
        if cmt:
            meta += f"  {cmt}"
        lines.append(f"[{item['id']}] {item['title']}")
        lines.append(f"  {item['url']}")
        lines.append(f"  {meta}")
        desc = item.get("desc", "")
        if desc:
            lines.append(f"  {desc[:120]}{'...' if len(desc) > 120 else ''}")
        lines.append("")
    return _tg("\n".join(lines).rstrip())


def fmt_topic(data: dict, fmt: str) -> str:
    if fmt == "json":
        return json.dumps(data, ensure_ascii=False, indent=2)

    title = data.get("title", "")
    url = data.get("url", "")
    pts = data.get("points", 0)
    author = data.get("author", "")
    time = data.get("time", "")
    body = data.get("body", "")
    comments = data.get("comments", [])
    topic_url = data.get("topic_url", "")

    lines = [
        title,
        url,
        f"{pts}P  {author}  {time}",
        f"토픽: {topic_url}",
        "─" * 50,
    ]
    if body:
        lines.append("")
        lines.append(body)

    if comments:
        lines.append("")
        lines.append(f"── 댓글 {len(comments)}개 " + "─" * 30)
        for c in comments:
            indent = "  " * c["depth"]
            lines.append("")
            lines.append(f"{indent}[{c['author']}] {c['time']}")
            for ln in c["text"].splitlines():
                lines.append(f"{indent}  {ln}")

    return _tg("\n".join(lines))
