<div align="center">

  # meowtime
  an utility for gaining insights into the development of your projects
  <br>

  [why](#why)
</div>

---

> [!IMPORTANT]
> meowtime hasn't been shipped yet, and is being built right now, keep an eye out for updates

> [!NOTE]
> any questions? contact me at santiago [at] hackclub [dot] app and i'll try to read your message and reply to it as soon as possible.

# why?
i'm a person that participates on [hack club](https://hackclub.com/) events, and a requirement for entering on them and earning stuff is using something called hackatime

and don't get me wrong, hackatime infrastructure **it's actually really great!** and i took inspiration from it on how some stuff could be done, so look it up here: https://github.com/hackclub/hackatime

however, i'm not really talking about hackatime itself, i'm talking about how wakatime works. wakatime works by **tracking time** and doesn't **track what you really do**, for example, let's say i took about 15-20 tries to compile my program just to get it running, wakatime most cases wouldn't even pick it up as it's not intended for that + if it does, it will just grab all the compilation time instead, but not the actual effort **spent** on it.

and that's the issue i found! time is not really a proper way to metric if someone really gave all their effort on what they're doing, not only that, but is really simple to manipulate by the user due to how wakatime heartbeats work, generating a zero trust environment where people trying to analyze metrics don't really trust them and up enforcing other regulations for determining if someone worked on their project or not.

that's where meowtime comes in. it focuses on actions done in programs by allowing extension developers to add their respectful metadata for specific programs, so those actions can be tracked and analyzed, making it way harder to manipulate and making it way easier to notice when someone is actually working on their project.

# what's the difference with wakatime?
wakatime **tracks time**, meowtime **tracks what you do exactly**, that's the general difference.

meowtime intention is not to track you on the IDE, but rather to track where are you exactly, and what are you doing exactly. this means meowtime doesn't close up only for the development environment, but also opens up for any other type of stuff that might be done on a project such as: composing music, designing a PCB, researching and whatever floats your boat.

# how does it work?
meowtime uses a system called **sonar** which stores metadata about actions done by users, extensions from programs will call an API endpoint to send the metadata to the server, which will then be stored on the server when needed to calculate an project when requested.

we rely on Gofiber for the backend being a really lightweight framework for creating REST APIs and handling multiple requests efficiently. additionally, we use ScyllaDB as our database for storing data, which is a highly scalable and performant NoSQL database that can handle large amounts of data and provide fast read and write operations, making it an ideal choice for the project.
