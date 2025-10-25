[![Actions](https://github.com/tlinden/epuppy/actions/workflows/ci.yaml/badge.svg)](https://github.com/tlinden/epuppy/actions)
[![License](https://img.shields.io/badge/license-GPL-blue.svg)](https://github.com/tlinden/epuppy/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/tlinden/epuppy)](https://goreportcard.com/report/github.com/tlinden/epuppy)


# epuppy - terminal epub reader

This is a little TUI epub ebook reader. This is a work in progress and
may not work for all EPUB files yet. It uses a modified version of the
[epub module](https://github.com/kapmahc/epub/), which seems to be
unmaintained but the best I could find to parse EPUBs. Find it in the
`pkg/epub/` directory.

The idea behind this tool is to be able to just take a look into some
epub file without the need to leave the shell. And it had to be fast
enough to just peak into an ebook. However, it is possible to actually
read epub ebooks with epuppy but I'd encourage you to buy a hardware
ebook reader with an e-ink display. It's better for your eyes in the
long run.

## Screenshots

- Viewing an ebook in dark mode
![Screenshot](https://github.com/TLINDEN/epuppy/blob/main/.github/assets/darkmode.png)

- Viewing an ebook in light mode
![Screenshot](https://github.com/TLINDEN/epuppy/blob/main/.github/assets/light.png)

- You can interactively adjust text width
![Screenshot](https://github.com/TLINDEN/epuppy/blob/main/.github/assets/margin.png)

- Showing the help
![Screenshot](https://github.com/TLINDEN/epuppy/blob/main/.github/assets/help.png)

## Usage

To read an ebook, just give a filename as argument to `epuppy`.

Add  the option  `-s` to  store and  use a  previously stored  reading
progress.

Sometimes you may be unhappy with the colors. Depending on your
terminal style you can enable dark mode with `-D`, light mode is the
default. You can also configure custom colors in a config file in
`$HOME/.config/epuppy/confit.toml`:

```toml
# color setting for dark mode
colordark = {
  body = "#ffffff",
  title = "#7cfc00",
  chapter = "#ffff00"
}

# color setting for light mode
colorlight = {
  body = "#000000",
  title = "#8470ff",
  chapter = "#00008b"
}

# always use dark mode
dark = true
```

There are also cases where your current terminal just doesn't have the
capabilites for this stuff. I stumbled upon such a case during an SSH
session from my Android phone to a FreeBSD server. For this you can
either just disable colors with `-N` or by setting the environment
variable `$NO_COLOR` to 1. Or you can just dump the text of the ebook
and pipe it to some pager, e.g.:

```default
epuppy -t someebook.epub | less
```

There are also a couple of debug options etc, all options:

```default
Usage epuppy [options] <epub file>

Options:
-D --dark                enable dark mode
-s --store-progress      remember reading position
-n --line-numbers        add line numbers
-c --config <file>       use config <file>
-t --txt                 dump readable content to STDOUT
-x --xml                 dump source xml to STDOUT
-N --no-color            disable colors (or use $NO_COLOR env var)
-d --debug               enable debugging
-h --help                show help message
-v --version             show program version
```

## Installation

The tool does not have any dependencies.  Just download the binary for
your platform from the releases page and you're good to go.

### Installation using a pre-compiled binary

You can use [stew](https://github.com/marwanhawari/stew) to install epuppy:
```default
stew install tlinden/epuppy
```

Or go to the [latest release page](https://github.com/TLINDEN/epuppy/releases/latest)
and look for your OS and platform. There are two options to install the binary:

Directly     download     the     binary    for     your     platform,
e.g. `epuppy-linux-amd64-0.0.2`, rename it to `epuppy` (or whatever
you like more!)  and put it into  your bin dir (e.g. `$HOME/bin` or as
root to `/usr/local/bin`).

Be sure  to verify  the signature  of the binary  file. For  this also
download the matching `epuppy-linux-amd64-0.0.2.sha256` file and:

```shell
cat epuppy-linux-amd64-0.0.2.sha25 && sha256sum epuppy-linux-amd64-0.0.2
```
You should see the same SHA256 hash.

You  may  also download  a  binary  tarball  for your  platform,  e.g.
`epuppy-linux-amd64-0.0.2.tar.gz`,  unpack and  install it.  GNU Make  is
required for this:
   
```shell
tar xvfz epuppy-linux-amd64-0.0.2.tar.gz
cd epuppy-linux-amd64-0.0.2
sudo make install
```

### Installation from source

Check out the repository and execute `go build`, then copy the
compiled binary to your `$PATH`.

Or, if you have GNU Make installed, just execute:

```default
make
sudo make install
```

# Report bugs

[Please open an issue](https://github.com/TLINDEN/epuppy/issues). Thanks!

# License

This work is licensed under the terms of the General Public Licens
version 3.

# Author

Copyleft (c) 2025 Thomas von Dein
