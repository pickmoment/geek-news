import sys
import click
from . import api
from . import format as fmt


_FORMAT_OPTION = click.option(
    "--format", "-f",
    "output_fmt",
    type=click.Choice(["json", "telegram"]),
    default="json",
    show_default=True,
    help="출력 형식",
)

_PAGE_OPTION = click.option(
    "--page", "-p",
    default=1,
    show_default=True,
    help="페이지 번호",
)

_LIMIT_OPTION = click.option(
    "--limit", "-n",
    default=0,
    help="최대 항목 수 (기본: 전체)",
)


def _out(text: str):
    click.echo(text)


def _err(msg: str):
    click.echo(f"오류: {msg}", err=True)
    sys.exit(1)


@click.group()
def cli():
    """GeekNews (긱뉴스) CLI"""


@cli.command()
@_PAGE_OPTION
@_LIMIT_OPTION
@_FORMAT_OPTION
def list(page: int, limit: int, output_fmt: str):
    """인기 토픽 목록 조회

    예) gn list
        gn list -p 2 -f telegram
        gn list -n 10 -f telegram
    """
    try:
        data = api.topics("/", page)
        if limit > 0:
            data = data[:limit]
        _out(fmt.fmt_topics(data, output_fmt))
    except Exception as e:
        _err(str(e))


@cli.command()
@_PAGE_OPTION
@_LIMIT_OPTION
@_FORMAT_OPTION
def new(page: int, limit: int, output_fmt: str):
    """최신 토픽 목록 조회

    예) gn new
        gn new -p 2 -f telegram
    """
    try:
        data = api.topics("/new", page)
        if limit > 0:
            data = data[:limit]
        _out(fmt.fmt_topics(data, output_fmt))
    except Exception as e:
        _err(str(e))


@cli.command()
@click.argument("topic_id")
@_FORMAT_OPTION
def view(topic_id: str, output_fmt: str):
    """토픽 상세 및 댓글 조회

    TOPIC_ID: 토픽 ID (예: 28861)
              목록에서 [ID] 형태로 표시됨

    예) gn view 28861
        gn view 28861 -f telegram
    """
    try:
        data = api.topic(topic_id)
        _out(fmt.fmt_topic(data, output_fmt))
    except Exception as e:
        _err(str(e))
