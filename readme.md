# What is this?
This is nothing more then a project for me to get good at golang.

## Goals

+ Pure Go: I want this bot to be pure go. Every package and library used must be pure go as well. No outsourcing things to the ffmpeg command, using `os/exec` to run youtube-dl. Currently we're *commiting treason* which is making this extremely diffucult.

+ Improvement in style: There is no style for this project, and for good reason. I'm using this project to force myself to learn/write lots of packages and learn how to interact with and understand a API, etc, which will show me a lot of different style. Right now I like diaonds though.

## What does the bot do
Butch has 4 post-testing modules, one alpha module. The 4 post-testing ones are: `basic` *the basic things* - help - prefix: `profile` *Get to know those in your server* - profile - profiles: `appoint` *Setup appointments* - newbool - removebool - rsvp - editbool - pickedup - bool - bools.

#FQA - Frequently Questioned Answers

+ Why does this use [arikawa](https://github.com/diamondburned/arikawa) and not [discordgo](https://github.com/bwmarrin/discordgo) or [disgord](https://github.com/andersfylling/disgord)?

	• Because arikawa is a lot simpler, more elegant, and just generally a better structured api.

+ Why do you use [your fork of godesu](https://github.com/lordrusk/godesu) instead of the original [godesu](https://github.com/mtarnawa/godesu)?

	• There are a few problems that need to be fixed, and some improvements that I'm planning on implementing into my fork. At some point it will become it's own project, seperate from the original. But until I have time to develope my own, I'll just use this simple fork.
