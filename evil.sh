#!/bin/sh

# Put this in /usr/local/etc/rc.d, give it root:wheel permissions and make it
# executable to create a FreeBSD service.

# PROVIDE: evil

. /etc/rc.subr

name="evil"
bin="evil-feed-reader"
rcvar=evil_enable

command="/home/freebsd/${bin}"
command_args="-feeds /home/freebsd/feeds.cfg -port 8000"

evil_user=freebsd
start_cmd="/usr/sbin/daemon -u ${evil_user} ${command} ${command_args}"

load_rc_config $name
run_rc_command "$1"
