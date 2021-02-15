import requests
import pdb

data = """TOM!"
    No answer.
    "TOM!"
    No answer.
    "What's gone with that boy,  I wonder? You TOM!"
    No answer.
    The old lady pulled her spectacles down and looked over them about the
    room; then she put them up and looked out under them. She seldom or
    never looked THROUGH them for so small a thing as a boy; they were her
    state pair, the pride of her heart, and were built for "style," not
    service--she could have seen through a pair of stove-lids just as well.
    She looked perplexed for a moment, and then said, not fiercely, but
    still loud enough for the furniture to hear:
    "Well, I lay if I get hold of you I'll--"
    She did not finish, for by this time she was bending down and punching
    under the bedwith the broom, and so she needed breath to punctuate the
    punches with. She resurrected nothing but the cat.
    "I never did see the beat of that boy!"
    She went to the open door and stood in it and looked out among the
    tomato vines and "jimpson" weeds that constituted the garden. No Tom.
    So she lifted up her voice at an angle calculated for distance and
    shouted:
    "Y-o-u-u TOM!"""; 

def genHeaders():
  return {
      "Host": "capi.grammarly.com",
      "User-Agent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.12; rv:66.0) Gecko/20100101 Firefox/66.0",
      "Accept":"*/*",
      "Accept-Language":"en-US,en;q=0.5",
      "Accept-Encoding":"gzip, deflate, br",
      "Sec-WebSocket-Version":"13",
      "Origin":"moz-extension://1543dcbd-db1d-c043-be06-53f1ac20e6ef",
      "Sec-WebSocket-Extensions":"permessage-deflate",
      "Sec-WebSocket-Key":"SS3zJbe2IjlUBUy/jCuvkA==",
      "Connection":"keep-alive, Upgrade",
      "Cookie":"gnar_containerId=lfky8hb42qii082; grauth=AABGUBpUbqpxaM7cLVcQhkdY2FEb4xNLIb45XI1XZrgwqAHHtJf0pK8WZlbC3ol7w9NHYQqLUi5NEFsy; csrf-token=AABGUHL0N3ZMobFpPOTvz5dTk8cCcbrvFI1h8g; funnelType=free; browser_info=FIREFOX:66:COMPUTER:SUPPORTED:FREEMIUM:MAC_OS_X:MAC_OS_X; redirect_location=eyJ0eXBlIjoiIiwibG9jYXRpb24iOiJodHRwczovL3d3dy5ncmFtbWFybHkuY29tL3NpZ251cD9icmVhZGNydW1icz10cnVlJnV0bV9zb3VyY2U9ZmlyZWZveCZwYWdlPWZyZWUmZXh0ZW5zaW9uX2luc3RhbGw9dHJ1ZSZ1dG1fbWVkaXVtPXN0b3JlIn0=; firefox_freemium=true; _ga=GA1.2.650690555.1555234454; _gid=GA1.2.539082346.1555234454; ga_clientId=650690555.1555234454; _gcl_au=1.1.1340797246.1555234454; _fbp=fb.1.1555234461054.265805239; experiment_groups=",
      "Pragma":"no-cache",
      "Cache-Control":"no-cache",
      "Upgrade":"websocket",
      }

if __name__ == "__main__":
  headers = genHeaders()
  res = requests.post("https://capi.grammarly.com/api/check", headers=headers, data=data )
  print(res)
