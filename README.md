Windows Agent
============
An agent which allows execution of shell and powershell commands on Windows written in Go. It also exposes a REST interface to consume its services using any REST client.

Installation
============
In order to install this agent, you need to clone this repo and run the command (in its root directory):
```
go run *.go
```
This command will start an instance of this agent on port 3000

Usage
==========
This agent basically supports creating, deleting and listing sessions which contain information like working directory, environment variables, etc so that the user can afterwards invoke commands directly on these sessions with these parameters already set. It also supports executing shell and/or powershell commands inside or outside the context of a session.

Python Client
===========
A Python client is shipped with this agent which can directly consume its services. It basically consumes its REST services in Python.

```python
from windowsagent_client import WindowsAgentClient
# Create an instance of the client
client = WindowsAgentClient("<host>", <port>)

# Create a new session
sessionId = client.createSession('test_session', 'c:\\', cmdType='ps', env=[])

# List available sessions
sessions = client.listSessions()

# Delete a session by its ID
result = client.deleteSession(sessionId)

# Execute a powershell command and print its output
output = client.executeCommand('ls', '<administrator_password>', 'c:\\', args=[], env=[], cmdType='ps', sessionId=None)
print output

# Execute a normal command in a session context
output = client.executeCommand('netstat', '<administrator_password>', '', args=['-ntp'], env=[], cmdType='', sessionId=<your_session_id>)
```
