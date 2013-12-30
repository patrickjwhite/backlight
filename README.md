backlight
=========

As of now, this is a crude tool built on Qt 5.2 to provide a simple slider for
controlling my system's LED backlight.

![screenshot](http://i.imgur.com/loAUO87.png)

Why?
----

My purpose is to have finer control over the backlight than what the native
Ubuntu slider can offer. I have a MacBook Pro turned into a
[TuxBook](http://cweiske.de/tagebuch/tuxbook.htm) and the minimum screen
brightness that I can set is still too bright when I'm working at night but the
hardware does allow a less bright screen, since on OS X I can have this, so this
is barely an user interface limitation.

How it works?
-------------

[I have found](http://askubuntu.com/questions/282201) that controlling the
backlight brightness is quite simple, I just have to fill in a number under
`/sys/class/backlight/intel_backlight/brightness`,
and the maximum brightness can be read from
`/sys/class/backlight/intel_backlight/max_brightness`.
Notice that these are system specific configuration files. I have an intel based
MacBook Pro, with no dedicated graphics card. The slider just fill in a value in
the `brightness` file from 0 up to the value read from `max_brightness`.
The `brightness` file is created by the system at each boot, and its permissions
prevent changes from non-privileged users. So, by default to change it one must
alter it under root. To avoid always have to run the tool under root I instead
created a simple script to run on each boot that relax the permissions for the
`brightness` file so any user can make changes.

Requiriments
------------

- Qt (>=5.2.0)

- If [_QML support for the Go language_](https://github.com/oblitum/qml) is not
already installed, a prior proper Qt 5 environment for building the dependency
is necessary. For Ubuntu I use the following:

        QT5PATH=<Qt 5 install location>
        QT5VERSION=<Qt 5 Version Number, e.g., 5.2.0>
        export PKG_CONFIG_PATH=$QT5PATH/lib/pkgconfig
        export CGO_CPPFLAGS=-I$QT5PATH/include/QtCore/$QT5VERSION/QtCore

    more instructions at the _QML support for Go_ repository.

Installation
------------

    go get github.com/oblitum/backlight

I don't provide an install package, so, the installation is sparse currently.

- the tool is the single generated executable installed under `$GOPATH/bin`.
- There's a `.desktop` file that should be copied to
`~/.local/share/applications/` and a corresponding svg icon that should be
copied to `~/.local/share/icons/backlight/`. Please, edit the `.desktop` file
with a text editor for setting the correct icon path. This must be a fully
expanded path.
- There's a `backlight.sh` script to relax file permissions of the backlight
system files. It should be copied to `/etc/init.d/` and after that the command
`sudo update-rc.d -f backlight.sh defaults 99 ` should be executed for
registering the script to execute on boot.

Usage
-------------

The tool is a simple slider, to close it press space or escape. If you install
the `.desktop` and icon you can find the tool in the Dash by typing `backlight`
and it can be dragged to the Launcher.

Other tidbits
-------------

For my machine I've altered boot parameters to avoid the system [having two
places for backlight changes](http://askubuntu.com/questions/134984). If I
didn't do that, backlight could be set both from
`/sys/class/backlight/intel_backlight/brightness` and from
`/sys/class/backlight/acpi_video0/brightness`. Having two sources of
modification cause the system brightness control not be in sync with this tool.
To make the system use only the single real source of backlight modification I
change the `/etd/default/grub` boot parameter line to the following, adding the
`acpi_backlight=vendor` to the end:

    GRUB_CMDLINE_LINUX_DEFAULT="quiet splash acpi_backlight=vendor"

After changing I run `sudo update-grub` for it to take effect.

By the way, I changed my Ubuntu to boot with EFI boot, which is much faster.

Notice
------

This is the Go port of the original C++ project. For the C++ version check the
[cpp](https://github.com/oblitum/backlight/tree/cpp) branch.
