# obdi-saltconfigserver
Obdi plugin to allow configuration of servers

## Screenshot

![](images/saltconfigserver.png?raw=true)

## What is it?

This Obdi plugin allows the user to:
* view all servers in an environment.
* send a state.highstate to one or more servers.
* add or remove Salt States from individual servers.
  An External Node Classifier, or ENC, is provided as a REST end-point to
  accomplish this, along with a Python script that needs to be added to
  the Salt configuration file.
* assign GIT versions to one or more servers.
* view the Grains information for servers.
 
## Theory

#### Salt States and Formulas

Salt States and Formulas are referred to as Classes in all areas of the Web Interface.
This term has been borrowed from Puppet just because it's easier to fit into
various areas of the GUI.

When choosing states to add to a server only the states available for the current GIT
version are shown, along with a description of the state/formula. To supply a description
that will be shown in the Web Interface, the administrator adds comments at the top of
each state file containing the description. The format is as follows:

```
# INFO: A description for this state file
# INFO: that can span multiple lines
```

#### Datacentres and environments

Saved on each server as the grains: dc and env.

#### GIT Versioning

Saved on each server as the grain: version.

environment_\<semantic-version\>

For example, test_0.1.123 or prod_1.0.27

## Installation

#### Installing the plugin

Refer to [obdi-salt-repository](https://github.com/mclarkson/obdi-salt-repository) and [Obdi Plugins](https://github.com/mclarkson/obdi/blob/master/doc/plugins.md) to install this plugin.

## System configuration


