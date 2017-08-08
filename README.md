# sysmon

Shows a graph with the recent CPU usage history in your i3 status bar.

![demo](https://user-images.githubusercontent.com/9169414/29079574-d48fc87e-7c5d-11e7-895c-15c1fe500e86.gif)

# Installation

`
go get github.com/alcortesm/sysmon
`

# Usage

To use sysmon as part of your i3status bar:

1. Create an executable script
   with the following content
   and add it to your path: 

   ```
   #!/bin/bash
   i3status | while :
   do
       read line
       load=`sysmon`
       echo "sysmon: $load | $line" || exit 1
   done
   ```

   This will call the regular i3status command
   and the sysmon command
   and combine their outputs into a single line.

   For purpose of demonstration,
   we will call this script `i3status_with_sysmon`.

2. Now modify your `.i3/conf` file,
   to tell i3 to run your script,
   instead of the regular i3status command:

   ```
   bar {                                                                           
   #    status_command i3status                                               
       status_command i3status_with_sysmon                                               
   }
   ```
