import re
from html import unescape

from .client import get_html

_BASE = "https://news.hada.io"


def _clean(text: str) -> str:
    text = re.sub(r"<br\s*/?>", "\n", text, flags=re.IGNORECASE)
    text = re.sub(r"<[^>]+>", "", text)
    text = re.sub(r"[ \t]+", " ", text)
    lines = [ln.strip() for ln in text.splitlines()]
    text = "\n".join(lines)
    text = re.sub(r"\n{3,}", "\n\n", text)
    return unescape(text.strip())


def _parse_item(chunk: str) -> dict | None:
    id_m = re.search(r"data-topic-state-id='(\d+)'", chunk)
    title_m = re.search(r"<h1>(.*?)</h1>", chunk, re.DOTALL)
    if not id_m or not title_m:
        return None

    tid = id_m.group(1)
    title = unescape(title_m.group(1).strip())

    url_m = re.search(r"href='([^']+)' rel='nofollow'", chunk)
    url = url_m.group(1) if url_m else f"topic?id={tid}"
    if url.startswith("topic?") or url.startswith("/topic?"):
        url = f"{_BASE}/{url.lstrip('/')}"

    domain_m = re.search(r"<span class=topicurl>\(([^)]+)\)</span>", chunk)
    domain = domain_m.group(1) if domain_m else ""

    desc_m = re.search(r"class='c99 breakall'>(.*?)</a>", chunk, re.DOTALL)
    desc = _clean(desc_m.group(1)) if desc_m else ""

    pts_m = re.search(r"<span id='tp\d+'>(\d+)</span>", chunk)
    points = int(pts_m.group(1)) if pts_m else 0

    author_m = re.search(r"href='/@([^']+)'>[^<]+</a>", chunk)
    author = author_m.group(1) if author_m else ""

    time_m = re.search(r"href='/@[^']*'>[^<]+</a>\s*([^<]+?)\s*<span id='unvote", chunk)
    time = time_m.group(1).strip() if time_m else ""

    cmt_m = re.search(r"go=comments[^>]*>([^<]+)</a>", chunk)
    comments = cmt_m.group(1).strip() if cmt_m else ""

    return {
        "id": tid,
        "title": title,
        "url": url,
        "domain": domain,
        "desc": desc,
        "points": points,
        "author": author,
        "time": time,
        "comments": comments,
        "topic_url": f"{_BASE}/topic?id={tid}",
    }


def topics(section: str = "/", page: int = 1) -> list[dict]:
    url = f"{_BASE}{section}"
    if page > 1:
        sep = "&" if "?" in url else "?"
        url += f"{sep}page={page}"
    html = get_html(url)
    result = []
    for chunk in html.split("<div class='topic_row'")[1:]:
        item = _parse_item(chunk)
        if item:
            result.append(item)
    return result


def topic(tid: str) -> dict:
    html = get_html(f"{_BASE}/topic?id={tid}")

    title_m = re.search(r"class='bold ud'><h1>(.*?)</h1>", html, re.DOTALL)
    title = unescape(title_m.group(1).strip()) if title_m else ""

    url_m = re.search(r"href='([^']+)' class='bold ud'", html)
    url = url_m.group(1) if url_m else f"{_BASE}/topic?id={tid}"
    if url.startswith("topic?") or url.startswith("/topic?"):
        url = f"{_BASE}/{url.lstrip('/')}"

    pts_m = re.search(r"<span id='tp\d+'>(\d+)</span>", html)
    points = int(pts_m.group(1)) if pts_m else 0

    author_m = re.search(r"href='/@([^']+)'>[^<]+</a>", html)
    author = author_m.group(1) if author_m else ""

    time_m = re.search(r"<span title='([^']+)'>", html)
    time = time_m.group(1) if time_m else ""

    body = ""
    marker = "id='topic_contents'>"
    body_start = html.find(marker)
    if body_start >= 0:
        body_start += len(marker)
        end = html.find("</div></div></div>", body_start)
        raw = html[body_start:end] if end > 0 else html[body_start:body_start + 10000]
        body = _clean(raw)

    comments = []
    for chunk in html.split("<div class=comment_row")[1:]:
        cid_m = re.search(r"id=cid(\d+)", chunk)
        if not cid_m:
            continue
        depth_m = re.search(r"style=--depth:(\d+)", chunk)
        cauthor_m = re.search(r"href='/@([^']+)'>[^<]+</a>", chunk)
        ctime_m = re.search(r"href='comment\?id=\d+'>([^<]+)</a>", chunk)
        content_m = re.search(r"class='comment_contents'>(.*?)</span>", chunk, re.DOTALL)
        comments.append({
            "id": cid_m.group(1),
            "depth": int(depth_m.group(1)) if depth_m else 0,
            "author": cauthor_m.group(1) if cauthor_m else "",
            "time": ctime_m.group(1).strip() if ctime_m else "",
            "text": _clean(content_m.group(1)) if content_m else "",
        })

    return {
        "id": tid,
        "title": title,
        "url": url,
        "points": points,
        "author": author,
        "time": time,
        "body": body,
        "comments": comments,
        "topic_url": f"{_BASE}/topic?id={tid}",
    }
