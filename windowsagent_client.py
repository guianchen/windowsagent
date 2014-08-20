import requests
import json

class WindowsAgentClient(object):

    def __init__(self, host='localhost', port=3000):
        self.baseUrl = 'http://%s:%s' % (host, port)
        self.sessionUrl = self.baseUrl + '/session'

    def createSession(self, name, workingPath='', cmdType='ps', env=list()):
        data = {'name': name, 'workingPath': workingPath, 'cmdType': cmdType, 'env': json.dumps(env)}
        sessionId = requests.post(self.sessionUrl, data).json()
        return sessionId

    def deleteSession(self, sessionId):
        params = {'id': sessionId}
        result = requests.delete(self.sessionUrl, params=params).json()
        return result

    def listSessions(self):
        sessions = requests.get(self.sessionUrl).json()
        return sessions

    def executeCommand(self, cmd, password, workingPath=None, args=list(), env=list(), cmdType="ps", sessionId=None):
        url = self.baseUrl + '/execute'
        if sessionId:
            url += '/%s' % sessionId
        params = {'workingpath': workingPath, 'cmd': cmd, 'args': json.dumps(args),
                  'env': json.dumps(env), 'type': cmdType, 'password': password}
        output = requests.get(url, params=params).content
        return output


