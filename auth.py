import requests
import base64
import json

import pdb

def getInitHeaders():
   return {
         "Sec-Fetch-Mode":"navigate",
         "Sec-Fetch-Sit":"same-origin",
         "Sec-Fetch-User":"?1",
         "Upgrade-Insecure-Requests":"1",
         "Referer":"https://www.grammarly.com/",
         }

   # * Redirect locations aren't tied to a session. They are base64 encoded objects.
def genRedirectLocation():
   header = {
      "type" : "",
      "location": "https://www.grammarly.com/after_install_page?extension_install=true&utm_medium=store&utm_source=firefox"
         }
   msgBytes = json.dumps(header).encode('ascii') 
   return base64.b64encode(msgBytes)

def genAuthURL( contID ):
  user = 'oranonymous'
  app  = 'firefoxExt'
  return "https://auth.grammarly.com/v3/user/{}?app={}&containerId={}".format( user, app, contID )

def getCookies( reqCookies ):
   return {
         "funnelType:" : reqCookies['funnelType'],
         "gnar_containerId" : reqCookies['gnar_containerId'],
         "firefox_freemium?": 'true',
         "browser_info" : 'FIREFOX:67:COMPUTER:SUPPORTED:FREEMIUM:MAC_OS_X:MAC_OS_X',
         "redirect_location" : genRedirectLocation() }


def buildAuthHeaders( containerId, cookies ):
     return {
    'Accept-Encoding': 'gzip, deflate, br',
    'Accept-Language': 'en-GB,en-US;q=0.9,en;q=0.8',
    'Cache-Control': 'no-cache',
    'Cookie' : cookies,
    'Host' : 'auth.grammarly.com', 
    'Origin' : 'moz-extension://6adb0179-68f0-aa4f-8666-ae91f500210b',
    'Pragma': 'no-cache',
    'X-Container-Id': containerId,
    'X-Client-Version': '8.852.2307',
    'X-Client-Type': 'extension-firefox',
    'User-Agent':
      'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36'
  }

def genCookieStr( cookies ):
   out = ""
   for key, val in cookies.items():
      out = out + key + '=' + val + ';'
   return out

if __name__ == "__main__":

   initReq = requests.get("https://grammarly.com/signin", params=getInitHeaders())
   cookies = getCookies( initReq.cookies )
   url = genAuthURL( cookies['gnar_containerId'])
   headers = buildAuthHeaders(cookies['gnar_containerId'], genCookieStr(cookies))
   tryIt = requests.get(url, params=headers)
   print(tryIt.text)
   print(headers)

