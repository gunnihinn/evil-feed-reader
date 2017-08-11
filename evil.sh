#!/bin/sh

# Put this in /usr/local/etc/rc.d, give it root:wheel permissions and make it
# executable to create a FreeBSD service.

# PROVIDE: evil

. /etc/rc.subr

name="evil"
bin="evilfr"
rcvar=evil_enable

feedFile=/home/freebsd/evil.yaml
stateFile=/home/freebsd/.evil-state.json
logFile=/home/freebsd/evil.log
port=8000

command="/home/freebsd/${bin}"
command_args="-config $feedFile -state $stateFile -log $logFile -port $port"

evil_user=freebsd
start_cmd="/usr/sbin/daemon -u ${evil_user} ${command} ${command_args}"

load_rc_config $name
run_rc_command "$1"
