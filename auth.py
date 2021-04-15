# sample request headers
# GET /freews HTTP/1.1

import requests
import base64
import json
import websocket
import pdb
import sys

# sample cookies extracted from  requests
cookies = {
    "csrf-token" : "AABI8W/UnhdTileRmvWcDoIrqs7p9aO8OV9SzQ",
    "browser_info" : "FIREFOX:85:COMPUTER:SUPPORTED:FREEMIUM:MAC_OS_X:MAC_OS_X",
    "funnelType" : "free",
    "gnar_containerId": "aaukbtnoho4o302",
    "grauth" : "AABI8f0HwwHIIGWAMnnEx5503qNwPXfJSkDOV1H4XJZcxBKSzjATqAVh_Kn13lxM-IAo34dEN55NNN_j",
    "redirect_location": "eyJ0eXBlIjoiIiwibG9jYXRpb24iOiJodHRwczovL2FwcC5ncmFtbWFybHkuY29tL2Rkb2NzLzEwNTI0MTk1MDUifQ=="
    }

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
        'Cookie' : 'grauth=AABI8f0HwwHIIGWAMnnEx5503qNwPXfJSkDOV1H4XJZcxBKSzjATqAVh_Kn13lxM-IAo34dEN55NNN_j; csrf-token=AABI8W/UnhdTileRmvWcDoIrqs7p9aO8OV9SzQ; gnar_containerId=aaukbtnoho4o302; funnelType=free; browser_info=FIREFOX:85:COMPUTER:SUPPORTED:FREEMIUM:MAC_OS_X:MAC_OS_X;',
        'Host' : 'auth.grammarly.com', 
        'Origin' : 'moz-extension://6adb0179-68f0-aa4f-8666-ae91f500210b',
        'Pragma': 'no-cache',
        'X-Container-Id': containerId,
        'X-Client-Version': '8.852.2307',
        'X-Client-Type': 'extension-firefox',
        'User-Agent':
        'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36'
        }

def genReqStdHeaders( origin, host ):
  return {
    "Host": "capi.grammarly.com",
    "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.16; rv:85.0) Gecko/20100101 Firefox/85.0",
    "Accept": "*/*",
    "Accept-Language": "en-US,en;q=0.5",
    "Accept-Encoding": "gzip, deflate, br",
    "Origin": "moz-extension://6adb0179-68f0-aa4f-8666-ae91f500210b",
    "Host": "auth.grammarly.com",
    "Connection": "keep-alive, Upgrade",
    "Pragma": "no-cache",
    "Cache-Control": "no-cache",
    'X-Client-Version': '8.852.2307',
    'X-Client-Type': 'extension-firefox',
  }  

def genUniqueAuthHeafers( contID ):
  return {
    'X-Container-Id': "aaukbtnoho4o302",
 "Cookie": "grauth=AABI8f0HwwHIIGWAMnnEx5503qNwPXfJSkDOV1H4XJZcxBKSzjATqAVh_Kn13lxM-IAo34dEN55NNN_j; csrf-token=AABI8W/UnhdTileRmvWcDoIrqs7p9aO8OV9SzQ; gnar_containerId=aaukbtnoho4o302; funnelType=free; browser_info=FIREFOX:85:COMPUTER:SUPPORTED:FREEMIUM:MAC_OS_X:MAC_OS_X; redirect_location=eyJ0eXBlIjoiIiwibG9jYXRpb24iOiJodHRwczovL2FwcC5ncmFtbWFybHkuY29tL2Rkb2NzLzEwNTI0MTk1MDUifQ==; _gcl_au=1.1.1041563917.1613428568; _ga=GA1.2.941863428.1613428568; _gid=GA1.2.1297436573.1613428568; ga_clientId=941863428.1613428568; _ga_CBK9K2ZWWE=GS1.1.1613428568.1.1.1613428616.0; grInstallSource=funnel; _hjTLDTest=1; _hjid=e1b366c2-0436-49e1-9ab9-1e8ae4e931c8; _hjFirstSeen=1; _hjAbsoluteSessionInProgress=0; funnel_firstTouchUtmSource=firefox; _uetsid=327c20b06fde11eb818ab1a860fb9dce; _uetvid=327c27206fde11ebbd41f1d2fb1528b4; tdi=kqixp3o90zfh9ddt; isGrammarlyUser=true; experiment_groups=gdpr_signup_enabled|iframe_integration_salesforce_rollout_enabled|gb_analytics_mvp_phase_one_30_day_enabled|auto_complete_correct_safari_enabled|officeaddin_ue_exp3_enabled|denali_capi_all_enabled|extension_new_rich_text_fields_enabled|officeaddin_upgrade_state_exp1_enabled1|safari_migration_inline_disabled_enabled|officeaddin_outcomes_ui_exp5_enabled1|premium_ungating_renewal_notification_enabled|small_hover_menus_existing_enabled|quarantine_messages_enabled|fsrw_in_assistant_all_enabled|iframe_integration_zendesk_rollout_enabled|emogenie_beta_enabled|extension_fluid_for_all_rollout_test_enabled|officeaddin_upgrade_state_exp2_enabled1|gb_in_editor_premium_Test1|apply_formatting_all_enabled|gb_analytics_mvp_phase_one_enabled|extension_assistant_experiment_all_enabled|denali_link_to_kaza_enabled|gdocs_for_all_safari_enabled|extension_assistant_all_enabled|additional_payment_verification_control_2|safari_migration_backup_notif1_enabled|auto_complete_correct_edge_enabled|extension_assistant_bundles_all_enabled|safari_migration_popup_editor_disabled_enabled|extension_plt_improvements_enabled|officeaddin_proofit_exp3_enabled|safari_migration_inline_warning_enabled|iframe_integration_facebook_rollout_enabled|gdocs_for_all_firefox_enabled|gdocs_new_mapping_enabled|officeaddin_muted_alerts_exp2_enabled1|officeaddin_perf_exp3_enabled",
    }

def genWebsocketStdHeader():
  return {

      }

def extractAuthCookies( cookies ):
  return {
    "gnar_containerId" : "aaukbtnoho4o302",
    "grauth" : cookies.get("grauth"),
    "csrf-token" : cookies.get("csrf-token"),
    "redirect-location" : genRedirectLocation(),   
      }

def genFurtherHeaders( headers ):
  hdrs = {
    'firefox_freemium': 'true',
    'funnelType': 'free',
    'browser_info': 'FIREFOX:67:COMPUTER:SUPPORTED:FREEMIUM:MAC_OS_X:MAC_OS_X'
      }
  headers.update(hdrs)
  return genCookieStr( headers )


def genCookieStr( cookies ):
   out = ""
   for key, val in cookies.items():
     out = out + key + '=' + val + '; '
   return out

def buildPlagHdrs( cook, contID, origin="moz-extension://6adb0179-68f0-aa4f-8666-ae91f500210b", host="auth.grammarly.com" ):
  return {
    'Accept-Encoding': 'gzip, deflate, br',
    'Accept-Language': 'en-GB,en-US;q=0.9,en;q=0.8',
    'Cache-Control': 'no-cache',
    'Cookie' : cook,
    'Host' : host,
    'Origin' : origin,
    'Pragma': 'no-cache',
    'X-Container-Id': contID,
    'X-Client-Version': '8.852.2307',
    'X-Client-Type': 'extension-firefox',
    'User-Agent':
      'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36'
}

def buildInitialMsg():
  return {
  "type": 'initial',
  "docid": "1234",
  "client": 'extension_chrome',
  "protocolVersion": '1.0',
  "clientSupports": [
    'free_clarity_alerts',
    'readability_check',
    'filler_words_check',
    'sentence_variety_check',
    'free_occasional_premium_alerts'
  ],
  "dialect": 'british',
  "clientVersion": '14.924.2437',
  "extDomain": 'keep.google.com',
  "action": 'start',
  "id": "0",
  "sid": "0",
}
def buildAuth( origin="moz-extension://6adb0179-68f0-aa4f-8666-ae91f500210b", host="capi.grammarly.com", authURL ="" ):
  stdHeaders = genReqStdHeaders( host, origin )
  print(stdHeaders)
  print(type(stdHeaders))
  uniqueAuthHeaders = genUniqueAuthHeafers("")
  stdHeaders.update(uniqueAuthHeaders)

  if authURL == "":
    authURL = genAuthURL("aaukbtnoho4o302")
  req = requests.get( authURL, headers=stdHeaders)

  return extractAuthCookies( req.cookies )

def buildOTMessage(): 
  return {
  "ch": ['+0:0:{}:0'.format(data)],
  "rev": '0',
  "id": '0',
  "action": 'submit_ot'
  }

allowedGroups = [ 'Punctuation', 'Spelling' ]
def retrieveGrammarlyResponses( ws ):
   errorResponses = []
   while True:
      resp = ws.recv_data()
      jsonResp = json.loads(resp[1])
      if jsonResp.get( 'action' ) == 'finished':
         break
      if jsonResp.get( 'group' ) and jsonResp.get( 'group' ) in allowedGroups:
         if jsonResp.get( 'highlightText' ):
            errorResponses.append( jsonResp )
   return errorResponses

class GrammarlyError():
   def __init__(self, errorResponse):
      self.columnBegin = errorResponse.get( 'highlightBegin' )
      self.columnEnd = errorResponse.get( 'highlightEnd' )
      self.errorType = errorResponse.get( 'group' )
      self.replacements = errorResponse.get( 'replacements' )
      self.offendingText = errorResponse.get( 'text' )
   
   def __str__( self ):
      attrs = vars(self)
      return (', '.join("%s: '%s'" % item for item in attrs.items() ))

   def __setattr__( self, name, value ):
      if value is not None:
         self.__dict__[name] = str( value )

   # %f:%l:%sc:%ec: %m
   # f: file
   # l: line number
   # sc: start column
   # ec: end column
   # m: error message
   def toFormat( self ):
      def checkReplacement( replacement ):
         if replacement == '[]':
            return ' '
         else :
            return replacement

      formattedString = '' + file_name
      formattedString += ':' + str( line_number )
      formattedString += ':' + self.columnBegin
      formattedString += ':' + self.columnEnd
      errorMessage = ( ': ' + self.errorType + ' error: The offending text is \'' + 
                       self.offendingText + '\'.' )
      if self.replacements:
         errorMessage += ( ' Replace \'' + self.offendingText + '\' with \'' + 
                           checkReplacement( self.replacements ) + '\'.' )
      return formattedString + errorMessage

'''
data = ( 'A session holding a lock on the running datastore also locks the CLI '
    'configuration lock. If this session now faces an error during a '
    'transaction Commit, it might have impacted communication on the EAPI '
    'session. If the EAPI session closed due to an error, the CLI lock would '
    'be released, so it\'s important for the client to know that the lock was '
    'lost. In the absence of a way to tell this to the client, we close the '
    'connection.' )
'''

data = ( 'This is the begvinning of the text.\n' 
         'This line has two  spaces.\n'
         'I accidently typed this this twice.\n'
         'no uppercase\n' )

file_name = "testfile.py"
line_number = 34

if __name__ == "__main__":
  
  stdinData = ''
  for line in sys.stdin:
     stdinData += line
  data = stdinData

  auth = buildAuth()
  print("here")
  cookieStr = genFurtherHeaders( auth )
    
  plagHdrs = buildPlagHdrs( cookieStr, auth["gnar_containerId"], host='capi.grammarly.com')

  #req = requests.post('https://capi.grammarly.com/api/check', headers=plagHdrs, data=data )
  url = "wss://capi.grammarly.com/freews"
  ws = ws = websocket.WebSocket()
  ws.connect(url, header=plagHdrs, origin="moz-extension://6adb0179-68f0-aa4f-8666-ae91f500210b" )
  #ws.connect(url, origin="moz-extension://6adb0179-68f0-aa4f-8666-ae91f500210b" )
  ws.send(json.dumps(buildInitialMsg()))
  ws.send(json.dumps(buildOTMessage()))
  responses = retrieveGrammarlyResponses( ws )
  errorResponses = [ GrammarlyError( error ) for error in responses ]
  for error in errorResponses:
      print error.toFormat()
