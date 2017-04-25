package sysmon

import "github.com/godbus/dbus/introspect"

const IntrospectDataString = `<node name="` + Path + `">
	<interface name="` + Name + `">
		<method name="OneMinLoadAvg">
			<arg direction="out" type="d"/>
		</method>
		<method name="FiveMinLoadAvg">
			<arg direction="out" type="d"/>
		</method>
		<method name="FifteenMinLoadAvg">
			<arg direction="out" type="d"/>
		</method>
		<method name="RunnableCount">
			<arg direction="out" type="i"/>
		</method>
		<method name="ExistsCount">
			<arg direction="out" type="i"/>
		</method>
		<method name="LastCreatedPID">
			<arg direction="out" type="i"/>
		</method>
	</interface>` + introspect.IntrospectDataString + `</node>`