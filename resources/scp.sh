#!/usr/bin/expect

#set passwd "bluemoon2016#"
set username [lindex $argv 0]
set host [lindex $argv 1]
set passwd [lindex $argv 2]
set port [lindex $argv 3]
#set remoteDir [lindex $argv 4]
set resources [lindex $argv 4]

#spawn scp /data/watchDog/watch-dog appadm@192.168.240.53:/home/appadm
spawn scp -P $port $resources  $username@$host:/data/

#spawn scp -P $port /data/watchDog/watch-dog $username@$host:$remoteDir/watch-dog.tmp

expect {
  "password"
	{
	  send "$passwd\n"
	}
  "密码："
        {
          send "$passwd\n"
        }
   "pass"
        {
          send "$passwd\n"
        }
   "yes/no"
        {
          send_user "send yes"
          send "yes\n"
        }
   eof
    {
        send_user "eof\n"
    }
}

set timeout 10000
send "exit\r"
expect eof
