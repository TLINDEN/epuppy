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

## Installation

The tool does not have any dependencies.  Just download the binary for
your platform from the releases page and you're good to go.

### Installation using a pre-compiled binary

Go to the [latest release page](https://github.com/TLINDEN/epuppy/releases/latest)
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

Copyleft (c) 2024 Thomas von Dein
