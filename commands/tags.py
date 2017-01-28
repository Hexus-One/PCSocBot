from pony.orm import select

from commands.base import Command
from helpers import bold, CommandFailure
from models import Tag
from utils.embed_table import EmbedTable


class Tags(Command):
    db_required = True
    pprint = dict(platform="platform/game")
    desc = "Player tag storage for the UNSW PCSoc discord server."


class Add(Tags):
    desc = "Adds/changes a player tag with associated platform/game to the list"
    def eval(self, platform, tag):
        warning = ''
        if Tag.get(platform=platform.lower()) is None:
            platforms = ', '.join(sorted(select(t.platform for t in Tag)[:]))
            warning = bold('WARNING: creating a new platform. Please check that the platform doesn\'t already '
            'exist by another name.\n') + 'Current platforms are ' + platforms + '\n'
        Tag.create_or_update(user=self.user, platform=platform.lower(), tag=tag)
        return "%s%s added as %s tag for %s" % (warning, tag, platform.title(), self.name)


class Remove(Tags):
    desc = "Removes a user/player tag to the bot."
    def eval(self, platform):
        Tag.delete_or_err(user=self.user, platform=platform.lower())
        return "%s tag for %s removed" % (platform.title(), self.name)


class Get(Tags):
    desc = "Returns your own tag for a specified platform / game"
    def eval(self, platform):
        tag = Tag.get_or_err(user=self.user, platform=platform.lower())
        return "The %s tag of %s is %s" % (platform, self.name, tag.tag)


class List(Tags):
    desc = "Returns a list of user tags for a specified platform"
    def eval(self, platform):
        platform = platform.lower()
        tags = Tag.select_or_err(lambda x: x.platform == platform)
        return EmbedTable(fields=['User', 'Tag'],
                           table=[(self.from_id(tag.user).name, tag.tag) for tag in tags],
                           title="Showing tags for " + platform.title(), colour=self.EMBED_COLOR)


class View(Tags):
    desc = "Returns a list of tags for a specified user"
    def eval(self, name):
        user = self.from_name(name)
        if user is None:
            raise CommandFailure("User not found")
        tags = Tag.select_or_err(lambda x: x.user == int(user.id), "User has no tags")
        return EmbedTable(fields=['Platform', 'Tag'],
                          table=[(x.platform.title(), x.tag) for x in tags],
                          colour=self.EMBED_COLOR, user=user,
                          title="Tags for " + bold(user.name))
