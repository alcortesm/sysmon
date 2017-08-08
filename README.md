# sysmon

Shows a graph with the recent CPU usage history in your i3 status bar.

![demo](https://user-images.githubusercontent.com/9169414/29079574-d48fc87e-7c5d-11e7-895c-15c1fe500e86.gif)

# Installation

`
go get github.com/alcortesm/sysmon
`

# Usage

Normally,
you configure your i3 window manager
to run the i3status bar
by adding the following lines to the `i3/conf` file:

```
bar {                                                                           
    status_command i3status                                               
}
```

To use sysmon as part of your i3status,
create an executable script as follows
anywhere inside your path: 

```
#!/bin/bash
i3status | while :
do
    read line
    load=`sysmon`
    echo "sysmon: $load | $line" || exit 1
done
```

and call it from your `.i3/conf` file:

```
bar {                                                                           
#    status_command i3status                                               
    status_command myi3status                                               
}
```




