# Wenda

## Summary

- Fully online website
- Allows you to create (untimed) tasks, that you can see in weekly agenda
- Can send invites to friends over Discord (i.e. gaming sessions, meetings, etc) which would be reflected on the calendar and weekly agenda
- Integrates seamlessly with Discord to the point where you are able to do every feature you could on the website through Discord
  - This can be implemented through a Discord bot which users can pm
- Needs some form of friending to prevent spam (can be implemented later)

## Technical details

- Authentication should be done through Discord, which would let us pair Wenda accounts to Discord accounts directly
- Frontend written in React (typescript)
- Design done in Figma
- Backend written in Golang 
- PostgreSQL DB (perhaps Google's Firebase)
- Hosting done through GCP or another cloud provider
