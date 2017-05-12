package sysmon

import "github.com/godbus/dbus/introspect"

const IntrospectDataString = `<node name="` + Path + `">
	<interface name="` + InterfaceName + `">
		<method name="LoadAvgs">
			<arg direction="out" type="s"/>
		</method>
		<method name="Dev">
			<arg direction="out" type="ad"/>
		</method>
	</interface>` + introspect.IntrospectDataString + `</node>`

/*
	<method name="OneMinLoadAvg">
		<arg direction="out" type="ad"/>
	</method>
	<method name="FiveMinLoadAvg">
		<arg direction="out" type="s"/>
	</method>
	<method name="FifteenMinLoadAvg">
		<arg direction="out" type="d"/>
	</method>
	<method name="RunnableCount">
		<arg direction="out" type="t"/>
	</method>
	<method name="ExistsCount">
		<arg direction="out" type="t"/>
	</method>
	<method name="LastCreatedPID">
		<arg direction="out" type="t"/>
	</method>
*/
