# Terminal client for manpages.debian.org

This tool allows to view Debian manpages in your terminal without installing
the corresponding packages (and even without installing Debian).

```
usage: man [SECTION] MANPAGE

Examples:
        man man
        man 2 flock
        man bookworm/ddrescue
```

Dependencies:

  - groff (from `groff-base` in Debian)
  - less (or any other pager, will fall back to `cat` if nothing else found)

All manpages are published and maintained by Debian project.
This is just a simple client for <https://manpages.debian.org>


## Installation

  - Build and install with Go: `go install github.com/sio/man/cmd/man@latest`

  - Prebuilt binaries are also available in [GitHub Releases](https://github.com/sio/man/releases).

    The build process is reproducible.
    Use your own judgement to decide whether to trust the binaries publshed there.


## License and copyright

Copyright 2024 Vitaly Potyarkin

```
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
```


## Bugs and possible improvements

- Some manpages are not rendered correctly (try `groff_man`). This may or may
  not turn out to be an [upstream] bug. Needs further investigation.

- Relative include directives do not work.
  Try `ash`: <https://manpages.debian.org/bookworm/ash/ash.1.en.gz>

- Almost 5MB is a rather large size for this simple tool.
  It may have even been just a single pipe of `curl | groff | less`.
  Can this be reimplemented in something like Python (stdlib-only) to simplify
  builds and distribution?
  Will distributing a Python script be simpler compared to a (large-ish) static binary?

[upstream]: https://github.com/Debian/debiman/
